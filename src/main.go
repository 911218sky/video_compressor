package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"video_compressor/src/config"
	"video_compressor/src/video"
)

func main() {
	// Parse command line arguments
	inputPath := flag.String("input", "", "Input video file path")
	outputPath := flag.String("output", "", "Output video file path")

	// Video compression parameters
	fps := flag.Int("fps", 32, "Frame rate (default: 32)")
	resolution := flag.String("resolution", "720p", "Video resolution (options: 1080p, 720p, 480p)")
	bitrate := flag.Int("bitrate", 0, "Custom bitrate in Kbps (0 for default)")
	preset := flag.String("preset", "p7", "Encoder preset (p1=fastest, p7=best quality)")
	cq := flag.Int("cq", 32, "Constant quality value (0-51, lower is better)")
	width := flag.Int("width", 0, "Custom width (0 for default)")
	height := flag.Int("height", 0, "Custom height (0 for default)")

	flag.Parse()

	// Check if input file exists
	if _, err := os.Stat(*inputPath); os.IsNotExist(err) {
		fmt.Printf("Error: Input file not found: %s\n", *inputPath)
		return
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(*outputPath)
	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Printf("Error: Failed to create output directory: %v\n", err)
			return
		}
	}

	// Create video config
	videoConfig := config.VideoConfig{
		Fps:        *fps,
		Resolution: *resolution,
		Bitrate:    *bitrate,
		Preset:     *preset,
		Cq:         *cq,
		Width:      *width,
		Height:     *height,
	}

	// If custom width/height is specified, clear resolution to prevent override
	if *width != 0 && *height != 0 {
		videoConfig.Resolution = ""
	}

	// If custom bitrate is specified, use it
	if *bitrate == 0 {
		_, _, videoConfig.Bitrate = config.GetResolutionSettings(videoConfig.Resolution)
	}

	// Compress the video
	if err := video.CompressVideo(*inputPath, *outputPath, videoConfig); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}
