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

	"github.com/fvbommel/sortorder"
)

// CompressVideo compresses the video using ffmpeg
func CompressVideo(inputPath, outputPath string, cfg config.VideoConfig, verbose bool) error {
	// Validate input format
	if !ffmpeg.IsSupportedFormat(inputPath) {
		return fmt.Errorf(
			"unsupported input format; supported: MP4, AVI, MKV, MOV, WMV, FLV, WEBM",
		)
	}

	// Handle output file extension
	ext := strings.ToLower(cfg.OutputExtension)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	if !ffmpeg.SupportedFormats[ext] {
		return fmt.Errorf(
			"unsupported output extension %q; supported: %v",
			ext, ffmpeg.SupportedFormatsKeys(),
		)
	}
	if filepath.Ext(outputPath) != ext {
		outputPath += ext
	}

	// Get original file size
	origSize, err := utils.GetVideoSize(inputPath)
	if err != nil {
		return fmt.Errorf("failed to get input file size: %v", err)
	}

	// Auto-calculate width and height if needed
	if cfg.Resolution != config.ResolutionNone && cfg.Width == 0 && cfg.Height == 0 {
		ow, oh, e := utils.GetVideoDimensions(inputPath)
		if e != nil {
			fmt.Printf("Warning: cannot get dimensions: %v\n", e)
		}
		w, h, br := utils.GetRecommendedSettings(cfg.Resolution, ow, oh)
		cfg.Width, cfg.Height, cfg.Bitrate = w, h, br
	}

	// Build ffmpeg arguments
	// Set input file
	args := []string{"-i", inputPath}
	// Determine codec and bitrate
	args = append(args, ffmpeg.DetermineCodec(ext, cfg)...)
	// Set fps
	args = append(args, "-r", strconv.Itoa(cfg.Fps))
	// Scale if width and height are set
	if cfg.Width > 0 && cfg.Height > 0 {
		args = append(args, "-vf", fmt.Sprintf("scale=%d:%d", cfg.Width, cfg.Height))
	}
	// Set container
	// Get the muxer name by removing the leading dot from the extension (e.g., ".mkv" becomes "mkv")
	muxer := strings.TrimPrefix(ext, ".")
	switch muxer {
	case "mkv":
		muxer = "matroska"
	case "ts":
		muxer = "mpegts"
	case "wmv":
		muxer = "asf"
	}
	args = append(args, "-f", muxer)
	// Overwrite output file
	args = append(args, outputPath, "-y")

	// Run FFmpeg
	cmd := exec.Command(cfg.FfmpegPath, args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Println("FFmpeg command:", cmd.String())
	}
	if err := cmd.Run(); err != nil {
		// Show warning if GPU encoding fails
		if cfg.Encoder == "gpu" && strings.Contains(err.Error(), "hevc_nvenc") {
			fmt.Println("Warning: NVIDIA GPU encoding failed. Please retry with CPU encoder.")
			return fmt.Errorf("gpu encoder error: %v", err)
		}
		return fmt.Errorf("ffmpeg execution error: %v", err)
	}

	// Show statistics
	newSize, err := utils.GetVideoSize(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get output file size: %v", err)
	}
	if verbose {
		fmt.Println("Compression completed!")
		fmt.Printf(
			"Original: %.2fMB, Compressed: %.2fMB, Reduction: %.2f%%\n",
			float64(origSize)/1024/1024,
			float64(newSize)/1024/1024,
			(1-float64(newSize)/float64(origSize))*100,
		)
	}
	return nil
}

// MergeVideos reencodes and merges all .ts and .mp4 files in the given directory
func MergeVideos(inputDir, outputPath string, cfg config.VideoConfig) error {
	// Handle and validate output file extension
	ext := strings.ToLower(cfg.OutputExtension)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	if !ffmpeg.SupportedFormats[ext] {
		return fmt.Errorf("unsupported output extension %q; supported: %v", ext, ffmpeg.SupportedFormatsKeys())
	}
	if filepath.Ext(outputPath) != ext {
		outputPath += ext
	}

	// Analyze minimum width/height ratio
	ratio, err := utils.AnalyzeVideoRatios(inputDir, "min")
	if err != nil {
		return fmt.Errorf("failed to analyze video dimensions: %v", err)
	}
	if cfg.Resolution != config.ResolutionNone {
		cfg.Width, cfg.Height = utils.GetResolutionDimensionsRatio(cfg.Resolution, ratio)
	} else {
		fmt.Println("No resolution specified, defaulting to 1080p")
		cfg.Width, cfg.Height = utils.GetResolutionDimensionsRatio(config.Resolution1080p, ratio)
	}
	fmt.Printf("Using resolution: %dx%d (ratio %.3f)\n", cfg.Width, cfg.Height, ratio)

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "video_merge_*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Search for all supported formats and .ts input files
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		l := strings.ToLower(name)
		extIn := filepath.Ext(l)
		if ffmpeg.SupportedFormats[extIn] {
			files = append(files, name)
		}
	}
	if len(files) == 0 {
		return fmt.Errorf("no video files (%v supported containers) found in %s", ffmpeg.SupportedFormatsKeys(), inputDir)
	}

	// Sort files
	if cfg.Reverse {
		sort.Slice(files, func(i, j int) bool {
			return sortorder.NaturalLess(files[i], files[j])
		})
	} else {
		sort.Slice(files, func(i, j int) bool {
			return sortorder.NaturalLess(files[j], files[i])
		})
	}

	// Print first 20 files after sorting
	maxShow := 20
	fmt.Printf("First %d files after sorting:\n", min(len(files), maxShow))
	for i, name := range files[:min(len(files), maxShow)] {
		fmt.Printf("  %d: %s\n", i+1, name)
	}

	// Re-encode and generate list
	listFile := filepath.Join(tempDir, "files.txt")
	var sb strings.Builder
	fmt.Println("Step 1: Re-encoding individual files...")
	for i, name := range files {
		in := filepath.Join(inputDir, name)
		tempOut := filepath.Join(tempDir, fmt.Sprintf("seg_%03d.%s", i, cfg.OutputExtension))
		fmt.Printf("  [%d/%d] %s → %s\n", i+1, len(files), name, filepath.Base(tempOut))
		if err := CompressVideo(in, tempOut, cfg, false); err != nil {
			return fmt.Errorf("failed to reencode %s: %v", name, err)
		}
		sb.WriteString(fmt.Sprintf("file '%s'\n", tempOut))
	}
	if err := os.WriteFile(listFile, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write list file: %v", err)
	}

	// Merge re-encoded segments
	fmt.Println("Step 2: Merging re-encoded segments...")
	filter := fmt.Sprintf(
		"scale=%d:%d:force_original_aspect_ratio=decrease,"+
			"pad=%d:%d:(ow-iw)/2:(oh-ih)/2,setsar=1",
		cfg.Width, cfg.Height, cfg.Width, cfg.Height,
	)
	args := []string{
		"-f", "concat", "-safe", "0",
		"-i", listFile,
		"-vf", filter,
	}
	// Insert codec+bitrate parameters
	args = append(args, ffmpeg.DetermineCodec(ext, cfg)...)
	// Force container, overwrite
	args = append(args,
		"-f", strings.TrimPrefix(ext, "."), outputPath, "-y",
	)

	cmd := exec.Command(cfg.FfmpegPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to merge videos: %v", err)
	}

	fmt.Printf("Merge complete, output: %s\n", outputPath)
	return nil
}
