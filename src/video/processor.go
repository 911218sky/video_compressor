package video

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

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

// CompressVideo compresses the video using ffmpeg
func CompressVideo(inputPath, outputPath string, videoConfig config.VideoConfig) error {
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

	oldSize, err := GetVideoSize(inputPath)
	if err != nil {
		return fmt.Errorf("error getting input file size: %v", err)
	}

	// If resolution is specified, override width, height and bitrate
	if videoConfig.Resolution != "" {
		videoConfig.Width, videoConfig.Height, videoConfig.Bitrate = config.GetResolutionSettings(videoConfig.Resolution)
	}
	// Construct ffmpeg command
	cmd := exec.Command(ffmpegPath,
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
	)

	// Create pipes for real-time output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Starting video compression with NVIDIA GPU encoder...")
	fmt.Println("Command:", cmd.String())

	// Run the command
	err = cmd.Run()
	if err != nil {
		// If NVIDIA encoder fails, try with CPU encoder
		if strings.Contains(err.Error(), "hevc_nvenc") {
			fmt.Println("\nNVIDIA GPU encoder not available, switching to CPU encoder...")
			cmd = exec.Command(ffmpegPath,
				"-i", inputPath,
				"-c:v", "libx265", // Use CPU HEVC encoder
				"-preset", "medium",
				"-crf", strconv.Itoa(videoConfig.Cq),
				"-b:v", fmt.Sprintf("%dk", videoConfig.Bitrate),
				"-maxrate", fmt.Sprintf("%dk", videoConfig.Bitrate),
				"-bufsize", fmt.Sprintf("%dk", videoConfig.Bitrate*2),
				"-r", strconv.Itoa(videoConfig.Fps),
				"-vf", fmt.Sprintf("scale=%d:%d", videoConfig.Width, videoConfig.Height),
				outputPath,
				"-y",
			)

			// Set up real-time output for CPU encoder
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Println("Command:", cmd.String())
			err = cmd.Run()
			if err != nil {
				return fmt.Errorf("ffmpeg error: %v", err)
			}
		} else {
			return fmt.Errorf("ffmpeg error: %v", err)
		}
	}

	// Get new size and print statistics
	newSize, err := GetVideoSize(outputPath)
	if err != nil {
		return fmt.Errorf("error getting output file size: %v", err)
	}

	fmt.Println("Compression complete!")
	fmt.Printf("Original size: %.2fMB\n", float64(oldSize)/1024/1024)
	fmt.Printf("Compressed size: %.2fMB\n", float64(newSize)/1024/1024)
	fmt.Printf("Reduced by %.2f%%\n", (1-float64(newSize)/float64(oldSize))*100)

	return nil
}
