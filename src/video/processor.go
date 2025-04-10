package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"video_compressor/src/config"
	"video_compressor/src/ffmpeg"
	"video_compressor/src/utils"
)

// CompressVideo compresses the video using ffmpeg
func CompressVideo(inputPath, outputPath string, videoConfig config.VideoConfig, print bool) error {
	// Check for ffmpeg
	ffmpegPath, err := ffmpeg.CheckFFmpeg()
	if err != nil {
		fmt.Println("FFmpeg not found, attempting to download...")
		if err := ffmpeg.DownloadFFmpeg(); err != nil {
			return fmt.Errorf("failed to download FFmpeg: %v", err)
		}
		ffmpegPath, err = ffmpeg.CheckFFmpeg()
		if err != nil {
			return fmt.Errorf("FFmpeg still not found after download: %v", err)
		}
	}

	// Check if input file format is supported
	if !ffmpeg.IsSupportedFormat(inputPath) {
		return fmt.Errorf("unsupported input file format. Supported formats: MP4, AVI, MKV, MOV, WMV, FLV")
	}

	// Check if output path has valid extension
	if !strings.HasSuffix(strings.ToLower(outputPath), ".mp4") {
		return fmt.Errorf("output file must have .mp4 extension")
	}

	oldSize, err := utils.GetVideoSize(inputPath)
	if err != nil {
		return fmt.Errorf("error getting input file size: %v", err)
	}

	// If resolution is specified, override width, height and bitrate
	if videoConfig.Resolution != config.ResolutionNone && videoConfig.Width == 0 && videoConfig.Height == 0 {
		// Get original video dimensions
		fmt.Println("Getting video dimensions...")
		originalWidth, originalHeight, err := utils.GetVideoDimensions(inputPath)
		if err != nil {
			fmt.Printf("Warning: Could not get video dimensions: %v\n", err)
			originalWidth, originalHeight = 0, 0
		}
		videoConfig.Width, videoConfig.Height, videoConfig.Bitrate = utils.GetRecommendedSettings(videoConfig.Resolution, originalWidth, originalHeight)
	}

	var cmd *exec.Cmd
	if videoConfig.Encoder == "gpu" {
		// Try GPU encoder first if specified
		args := []string{
			"-i", inputPath,
			"-c:v", "hevc_nvenc", // Use NVIDIA GPU HEVC encoder
			"-preset", videoConfig.Preset,
			"-rc", "vbr", // Variable bitrate
			"-cq", strconv.Itoa(videoConfig.Cq),
			"-b:v", fmt.Sprintf("%dk", videoConfig.Bitrate),
			"-maxrate", fmt.Sprintf("%dk", videoConfig.Bitrate),
			"-bufsize", fmt.Sprintf("%dk", videoConfig.Bitrate*2),
			"-r", strconv.Itoa(videoConfig.Fps),
			"-vf", fmt.Sprintf("scale=%d:%d", videoConfig.Width, videoConfig.Height),
			outputPath,
			"-y",
		}

		cmd = exec.Command(ffmpegPath, args...)
	} else {
		// Use CPU encoder
		args := []string{
			"-i", inputPath,
			"-c:v", "libx265", // Use CPU HEVC encoder
			"-preset", videoConfig.Preset,
			"-crf", strconv.Itoa(videoConfig.Cq),
			"-b:v", fmt.Sprintf("%dk", videoConfig.Bitrate),
			"-maxrate", fmt.Sprintf("%dk", videoConfig.Bitrate),
			"-bufsize", fmt.Sprintf("%dk", videoConfig.Bitrate*2),
			"-r", strconv.Itoa(videoConfig.Fps),
			"-vf", fmt.Sprintf("scale=%d:%d", videoConfig.Width, videoConfig.Height),
			outputPath,
			"-y",
		}
		cmd = exec.Command(ffmpegPath, args...)
	}

	if print {
		// Create pipes for real-time output
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		fmt.Println("Command:", cmd.String())
		fmt.Printf("Starting video compression with %s encoder...\n",
			map[string]string{"gpu": "NVIDIA GPU", "cpu": "CPU"}[videoConfig.Encoder])
		if videoConfig.Width > 0 && videoConfig.Height > 0 {
			fmt.Printf("Output resolution: %dx%d\n", videoConfig.Width, videoConfig.Height)
		}
	}

	// Run the command
	err = cmd.Run()
	if err != nil {
		// If GPU encoder was selected but failed, offer to try CPU encoder
		if videoConfig.Encoder == "gpu" && strings.Contains(err.Error(), "hevc_nvenc") {
			fmt.Println("\nNVIDIA GPU encoder failed. Please try using the CPU encoder instead.")
			return fmt.Errorf("GPU encoder error: %v", err)
		}
		return fmt.Errorf("ffmpeg error: %v", err)
	}

	// Get new size and print statistics
	newSize, err := utils.GetVideoSize(outputPath)
	if err != nil {
		return fmt.Errorf("error getting output file size: %v", err)
	}

	if print {
		fmt.Println("Compression complete!")
		fmt.Printf("Original size: %.2fMB\n", float64(oldSize)/1024/1024)
		fmt.Printf("Compressed size: %.2fMB\n", float64(newSize)/1024/1024)
		fmt.Printf("Reduced by %.2f%%\n", (1-float64(newSize)/float64(oldSize))*100)
	}

	return nil
}

// MergeVideos reencodes and merges all .ts and .mp4 files in the given directory
func MergeVideos(inputDir, outputPath string, videoConfig config.VideoConfig) error {
	// Check if output path has valid extension
	if !strings.HasSuffix(strings.ToLower(outputPath), ".mp4") {
		return fmt.Errorf("output file must have .mp4 extension")
	}

	// get the min ratio
	ratio, err := utils.AnalyzeVideoRatios(inputDir, "min")
	if err != nil {
		return fmt.Errorf("failed to analyze video dimensions: %v", err)
	}

	// Check if video is portrait (ratio < 1) or landscape (ratio > 1)
	if videoConfig.Resolution != config.ResolutionNone {
		videoConfig.Width, videoConfig.Height = utils.GetResolutionDimensionsRatio(videoConfig.Resolution, ratio)
	} else {
		fmt.Println("No resolution specified, using 1080p")
		videoConfig.Width, videoConfig.Height = utils.GetResolutionDimensionsRatio(config.Resolution1080p, ratio)
	}

	fmt.Printf("Ratio: %f\n", ratio)
	fmt.Printf("Using resolution: %dx%d\n", videoConfig.Width, videoConfig.Height)

	// Create temporary directory for reencoded files
	tempDir, err := os.MkdirTemp("", "video_merge_*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up temp directory when done

	// Get all .ts and .mp4 files
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	var videoFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := strings.ToLower(file.Name())
		if strings.HasSuffix(name, ".ts") || strings.HasSuffix(name, ".mp4") {
			videoFiles = append(videoFiles, file.Name())
		}
	}

	if len(videoFiles) == 0 {
		return fmt.Errorf("no .ts or .mp4 files found in directory")
	}

	// Sort files to ensure consistent order
	sort.Slice(videoFiles, func(i, j int) bool {
		// Extract numbers from the beginning of filenames
		numI := 0
		numJ := 0
		fmt.Sscanf(videoFiles[i], "%d", &numI)
		fmt.Sscanf(videoFiles[j], "%d", &numJ)

		// If both files start with numbers, compare numerically
		if numI > 0 && numJ > 0 {
			return numI < numJ
		}

		// If only one file starts with a number, put it first
		if numI > 0 {
			return true
		}
		if numJ > 0 {
			return false
		}

		// If neither file starts with a number, compare alphabetically
		return strings.Compare(videoFiles[i], videoFiles[j]) < 0
	})

	// Create a file list for ffmpeg
	listFile := filepath.Join(tempDir, "files.txt")
	var listContent strings.Builder

	// Reencode each video file
	fmt.Println("Step 1: Reencoding individual files...")
	for i, fileName := range videoFiles {
		inputPath := filepath.Join(inputDir, fileName)
		tempOutput := filepath.Join(tempDir, fmt.Sprintf("reencoded_%d.mp4", i))

		fmt.Printf("Reencoding file %d/%d (%.1f%%): %s\n", i+1, len(videoFiles), float64(i+1)/float64(len(videoFiles))*100, fileName)
		err := CompressVideo(inputPath, tempOutput, videoConfig, false)
		if err != nil {
			return fmt.Errorf("failed to reencode %s: %v", fileName, err)
		}

		// Add to file list
		listContent.WriteString(fmt.Sprintf("file '%s'\n", tempOutput))
	}

	// Write the file list
	err = os.WriteFile(listFile, []byte(listContent.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file list: %v", err)
	}

	// Check for ffmpeg
	ffmpegPath, err := ffmpeg.CheckFFmpeg()
	if err != nil {
		return fmt.Errorf("ffmpeg not found: %v", err)
	}

	fmt.Println("\nStep 2: Merging reencoded files...")

	var cmd *exec.Cmd
	if videoConfig.Encoder == "gpu" {
		// Merge all reencoded files with padding and aspect ratio preservation using GPU
		cmd = exec.Command(ffmpegPath,
			"-f", "concat",
			"-safe", "0",
			"-i", listFile,
			"-c:v", "hevc_nvenc", // Use NVIDIA GPU encoder
			"-preset", videoConfig.Preset,
			"-rc", "vbr", // Variable bitrate
			"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,setsar=1",
				videoConfig.Width, videoConfig.Height,
				videoConfig.Width, videoConfig.Height),
			outputPath,
			"-y",
		)
	} else {
		// Merge all reencoded files with padding and aspect ratio preservation using CPU
		cmd = exec.Command(ffmpegPath,
			"-f", "concat",
			"-safe", "0",
			"-i", listFile,
			"-c:v", "libx265", // Use CPU HEVC encoder
			"-preset", videoConfig.Preset,
			"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,setsar=1",
				videoConfig.Width, videoConfig.Height,
				videoConfig.Width, videoConfig.Height),
			outputPath,
			"-y",
		)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to merge videos: %v", err)
	}

	fmt.Println("Video merge complete!")
	fmt.Printf("Output saved to: %s\n", outputPath)

	return nil
}
