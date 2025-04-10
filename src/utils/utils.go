package utils

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"video_compressor/src/config"
	"video_compressor/src/ffmpeg"
)

// GetVideoSize returns the size of the video file in bytes
func GetVideoSize(videoPath string) (int64, error) {
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// GetVideoDimensions returns the width and height of the video using ffprobe
func GetVideoDimensions(videoPath string) (width, height int, err error) {
	ffprobePath, err := ffmpeg.CheckFFprobe()
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe not found: %v", err)
	}

	cmd := exec.Command(ffprobePath,
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=p=0",
		videoPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe error: %v", err)
	}

	dimensions := strings.Split(strings.TrimSpace(string(output)), ",")
	if len(dimensions) != 2 {
		return 0, 0, fmt.Errorf("unexpected ffprobe output format")
	}

	width, err = strconv.Atoi(dimensions[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse width: %v", err)
	}

	height, err = strconv.Atoi(dimensions[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse height: %v", err)
	}

	return width, height, nil
}

// AnalyzeVideoRatios analyzes video aspect ratios in a directory and returns ratio based on specified mode
// mode: most_common, min, max, average
func AnalyzeVideoRatios(inputDir string, mode string) (ratio float64, err error) {
	// Get all video files in the directory
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read directory: %v", err)
	}

	// Filter supported video files
	var videoFiles []string
	for _, file := range files {
		if ffmpeg.IsSupportedFormat(file.Name()) {
			videoFiles = append(videoFiles, file.Name())
		}
	}

	if len(videoFiles) == 0 {
		return 0, fmt.Errorf("no valid videos found in directory")
	}

	// Sample size calculation - analyze at most 10 files or 20% of files, whichever is smaller
	sampleSize := 10
	if len(videoFiles) < sampleSize {
		sampleSize = len(videoFiles)
	} else if len(videoFiles) > 50 {
		sampleSize = len(videoFiles) * 3 / 10 // 30% sampling for large directories
	}

	// Shuffle and take sample
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(videoFiles), func(i, j int) {
		videoFiles[i], videoFiles[j] = videoFiles[j], videoFiles[i]
	})
	sample := videoFiles[:sampleSize]

	// Maps to store dimension frequencies and all ratios
	ratioFreq := make(map[float64]int)
	var ratios []float64

	// Analyze sampled video files
	for _, fileName := range sample {
		videoPath := filepath.Join(inputDir, fileName)
		w, h, err := GetVideoDimensions(videoPath)
		if err != nil {
			continue // Skip files that can't be analyzed
		}

		rawRatio := float64(w) / float64(h)
		var found bool
		for existingRatio := range ratioFreq {
			// if the ratio is within 0.2 of each other, group them
			if math.Abs(rawRatio-existingRatio) <= 0.2 {
				ratioFreq[existingRatio]++
				found = true
				break
			}
		}
		if !found {
			ratioFreq[rawRatio] = 1
		}
		ratios = append(ratios, rawRatio)
	}

	if len(ratios) == 0 {
		return 0, fmt.Errorf("no valid videos could be analyzed in sample")
	}

	switch mode {
	case "most_common":
		var maxFreq int
		for r, freq := range ratioFreq {
			if freq > maxFreq {
				maxFreq = freq
				ratio = r
			}
		}

	case "min":
		ratio = ratios[0]
		for _, r := range ratios {
			if r < ratio {
				ratio = r
			}
		}

	case "max":
		ratio = ratios[0]
		for _, r := range ratios {
			if r > ratio {
				ratio = r
			}
		}

	case "average":
		var sum float64
		for _, r := range ratios {
			sum += r
		}
		ratio = sum / float64(len(ratios))

	default:
		return 0, fmt.Errorf("invalid mode: %s. Supported modes: most_common, min, max, average", mode)
	}

	return ratio, nil
}

// GetRecommendedBitrate returns the recommended bitrate based on video dimensions
func GetRecommendedBitrate(width, height int) int {
	// Calculate total pixels
	pixels := width * height

	// Define bitrate thresholds based on pixel count
	switch {
	case pixels >= 8294400: // 3840x2160 (4K) or larger
		return 20000
	case pixels >= 5184000: // 2560x1440 (2K) or larger
		return 10000
	case pixels >= 2073600: // 1920x1080 or larger
		return 5000
	case pixels >= 921600: // 1280x720 or larger
		return 2500
	case pixels >= 460800: // 854x480 or larger
		return 1000
	case pixels >= 230400: // 640x360 or larger
		return 800
	case pixels >= 92160: // 320x288 or larger
		return 500
	default:
		return 250 // For very small resolutions
	}
}

// GetRecommendedSettings returns width and height based on resolution string and original dimensions
func GetRecommendedSettings(resolution config.Resolution, originalWidth, originalHeight int) (width, height, bitrate int) {
	switch resolution {
	case config.Resolution4K:
		height = 2160
		bitrate = 20000
	case config.Resolution2K:
		height = 1440
		bitrate = 10000
	case config.Resolution1080p:
		height = 1080
		bitrate = 5000
	case config.Resolution720p:
		height = 720
		bitrate = 2500
	case config.Resolution480p:
		height = 480
		bitrate = 1000
	default:
		height = 1080 // default to 1080p
		bitrate = 5000
	}

	// If original dimensions are not available, use default ratios
	if originalWidth == 0 || originalHeight == 0 {
		// Use 16:9 as default aspect ratio
		width = height * 16 / 9
		return width, height, bitrate
	}

	// Calculate new dimensions while maintaining aspect ratio
	aspectRatio := float64(originalWidth) / float64(originalHeight)

	if originalHeight > originalWidth {
		// Portrait video
		width = int(float64(height) * aspectRatio)
	} else {
		// Landscape video
		width = int(float64(height) * aspectRatio)
		// If width exceeds common maximums, scale down
		if width > 3840 { // Updated to support 4K
			width = 3840
			height = int(float64(width) / aspectRatio)
		}
	}

	// Ensure dimensions are even numbers (required by some codecs)
	width = width - (width % 2)
	height = height - (height % 2)

	return width, height, bitrate
}

// GetResolutionDimensions calculates width and height based on resolution and aspect ratio
func GetResolutionDimensionsRatio(resolution config.Resolution, ratio float64) (width, height int) {
	switch resolution {
	case config.Resolution4K:
		height = 2160
	case config.Resolution2K:
		height = 1440
	case config.Resolution1080p:
		height = 1080
	case config.Resolution720p:
		height = 720
	case config.Resolution480p:
		height = 480
	case config.Resolution360p:
		height = 360
	case config.Resolution240p:
		height = 240
	default:
		height = 1080 // default to 1080p
	}

	// Calculate width based on aspect ratio
	if ratio < 1 {
		// Portrait video
		width = int(float64(height) * ratio)
	} else {
		// Landscape video
		width = int(float64(height) * ratio)
		// If width exceeds common maximums, scale down
		if width > 3840 { // 4K max width
			width = 3840
			height = int(float64(width) / ratio)
		}
	}

	// Ensure dimensions are even numbers (required by some codecs)
	width = width - (width % 2)
	height = height - (height % 2)

	return width, height
}
