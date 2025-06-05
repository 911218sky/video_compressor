package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"video_compressor/src/config"
	"video_compressor/src/utils"
	"video_compressor/src/video"
)

func main() {
	// Parse command line arguments
	inputPath := flag.String("input", "", "Input video file path")
	outputPath := flag.String("output", "", "Output video file path (default: use input file name)")

	// Video compression parameters
	mode := flag.String("mode", "compress", "Mode (options: compress, merge)")
	fps := flag.Int("fps", 32, "Frame rate (default: 32)")
	resolution := flag.String("resolution", "", "Video resolution (options: 1080p, 720p, 480p)")
	bitrate := flag.Int("bitrate", 0, "Custom bitrate in Kbps (0 for default)")
	preset := flag.String("preset", "p7", "Encoder preset (p1=fastest, p7=best quality)")
	cq := flag.Int("cq", 32, "Constant quality value (0-51, lower is better)")
	width := flag.Int("width", 0, "Custom width (0 for default)")
	height := flag.Int("height", 0, "Custom height (0 for default)")
	encoder := flag.String("encoder", "gpu", "Encoder type (options: gpu, cpu)")

	flag.Parse()

	// Trim whitespace from input and output paths
	*inputPath = strings.TrimSpace(*inputPath)
	*outputPath = strings.TrimSpace(*outputPath)

	// Check if input file exists
	if _, err := os.Stat(*inputPath); os.IsNotExist(err) {
		fmt.Printf("Error: Input file not found: %s\n", *inputPath)
		return
	}

	// If output path is not specified, use input file name
	if *outputPath == "" {
		*outputPath = filepath.Base(*inputPath)
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
	resolutionStr, err := config.StringToResolution(*resolution)
	if err != nil {
		fmt.Printf("Error: Failed to convert resolution: %v\n", err)
		return
	}
	videoConfig := config.VideoConfig{
		Fps:        *fps,
		Resolution: resolutionStr,
		Bitrate:    *bitrate,
		Preset:     *preset,
		Cq:         *cq,
		Width:      *width,
		Height:     *height,
		Encoder:    *encoder,
	}

	// If custom width/height is specified, clear resolution to prevent override
	if *width != 0 && *height != 0 {
		videoConfig.Resolution = config.ResolutionNone
		fmt.Println("Custom width/height specified, ignoring resolution")
	}

	// If no custom width/height or resolution is specified, use 1080p
	if *width == 0 && *height == 0 && videoConfig.Resolution == config.ResolutionNone {
		videoConfig.Resolution = config.Resolution1080p
		fmt.Println("No custom width/height or resolution specified, using 1080p")
	}

	// If custom bitrate is specified, use it
	if *bitrate == 0 && videoConfig.Resolution != config.ResolutionNone {
		_, _, videoConfig.Bitrate = utils.GetRecommendedSettings(videoConfig.Resolution, 0, 0)
	}

	if *mode == "compress" {
		// Compress the video
		if err := video.CompressVideo(*inputPath, *outputPath, videoConfig, true); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	} else if *mode == "merge" {
		// Merge the video
		if err := video.MergeVideos(*inputPath, *outputPath, videoConfig); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}
}
