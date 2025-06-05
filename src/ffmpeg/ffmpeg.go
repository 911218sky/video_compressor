package ffmpeg

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"video_compressor/src/config"
)

// Supported video formats
var SupportedFormats = map[string]bool{
	".mp4":  true,
	".avi":  true,
	".mkv":  true,
	".mov":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
	".ts":   true,
}

// gpuSupportedExt lists container formats that work with NVIDIA HEVC GPU encoder.
var GpuSupportedExt = map[string]bool{
	".mp4": true,
	".mov": true,
	".mkv": true,
	".ts":  true,
}

// SupportedFormatsKeys returns a slice of the keys of a map
func SupportedFormatsKeys() []string {
	k := make([]string, 0, len(SupportedFormats))
	for x := range SupportedFormats {
		k = append(k, x)
	}
	return k
}

// DetermineCodec returns the codec-related FFmpeg args based on output extension.
func DetermineCodec(ext string, cfg config.VideoConfig) []string {
	// If GPU is requested but the container does not support HEVC_NVENC, fall back to CPU HEVC
	if cfg.Encoder == "gpu" && !GpuSupportedExt[ext] {
		fmt.Printf(
			"Warning: GPU encoding (hevc_nvenc) is not supported for %s, falling back to CPU libx265.\n",
			ext,
		)
		cfg.Encoder = "cpu"
	}

	if cfg.Encoder == "gpu" {
		// hevc_nvenc: supports -rc, -cq, -b:v, -maxrate, -bufsize
		return []string{
			"-c:v", "hevc_nvenc",
			"-rc", "vbr",
			"-cq", strconv.Itoa(cfg.Cq),
			"-b:v", fmt.Sprintf("%dk", cfg.Bitrate),
			"-maxrate", fmt.Sprintf("%dk", cfg.Bitrate),
			"-bufsize", fmt.Sprintf("%dk", cfg.Bitrate*2),
		}
	}

	switch ext {
	case ".mp4", ".mov", ".avi", ".flv", ".ts":
		// libx264: supports -preset, -crf, -b:v, -maxrate, -bufsize
		return []string{
			"-c:v", "libx264",
			"-preset", cfg.Preset,
			"-crf", strconv.Itoa(cfg.Cq),
			"-b:v", fmt.Sprintf("%dk", cfg.Bitrate),
			"-maxrate", fmt.Sprintf("%dk", cfg.Bitrate),
			"-bufsize", fmt.Sprintf("%dk", cfg.Bitrate*2),
		}

	case ".mkv":
		// libx265: supports -preset, -crf, -b:v, -maxrate, -bufsize
		return []string{
			"-c:v", "libx265",
			"-preset", cfg.Preset,
			"-crf", strconv.Itoa(cfg.Cq),
			"-b:v", fmt.Sprintf("%dk", cfg.Bitrate),
			"-maxrate", fmt.Sprintf("%dk", cfg.Bitrate),
			"-bufsize", fmt.Sprintf("%dk", cfg.Bitrate*2),
		}

	case ".webm":
		// libvpx-vp9: supports -b:v, -crf
		return []string{
			"-c:v", "libvpx-vp9",
			"-crf", strconv.Itoa(cfg.Cq),
			"-b:v", fmt.Sprintf("%dk", cfg.Bitrate),
		}

	case ".wmv":
		// wmv2: supports -b:v
		return []string{
			"-c:v", "wmv2",
			"-b:v", fmt.Sprintf("%dk", cfg.Bitrate),
		}

	default:
		// fallback to libx264
		return []string{
			"-c:v", "libx264",
			"-preset", cfg.Preset,
			"-crf", strconv.Itoa(cfg.Cq),
			"-b:v", fmt.Sprintf("%dk", cfg.Bitrate),
			"-maxrate", fmt.Sprintf("%dk", cfg.Bitrate),
			"-bufsize", fmt.Sprintf("%dk", cfg.Bitrate*2),
		}
	}
}

// IsSupportedFormat checks if the given file format is supported
func IsSupportedFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return SupportedFormats[ext]
}

// CheckFFmpeg checks if ffmpeg is available in PATH or current directory
func CheckFFmpeg() (string, error) {
	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Set executable suffix based on OS
	exeSuffix := ""
	if runtime.GOOS == "windows" {
		exeSuffix = ".exe"
	}

	// Check in current directory
	localFFmpeg := filepath.Join(currentDir, "ffmpeg"+exeSuffix)
	if _, err := os.Stat(localFFmpeg); err == nil {
		return localFFmpeg, nil
	}

	// Check in parent directory
	parentFFmpeg := filepath.Join(currentDir, "..", "ffmpeg"+exeSuffix)
	if _, err := os.Stat(parentFFmpeg); err == nil {
		return parentFFmpeg, nil
	}

	// Then check in PATH
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err == nil {
		return ffmpegPath, nil
	}

	return "", fmt.Errorf("ffmpeg not found")
}

// CheckFFprobe checks if ffprobe is available in PATH or current directory
func CheckFFprobe() (string, error) {
	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Set executable suffix based on OS
	exeSuffix := ""
	if runtime.GOOS == "windows" {
		exeSuffix = ".exe"
	}

	// Check in current directory
	localFFprobe := filepath.Join(currentDir, "ffprobe"+exeSuffix)
	if _, err := os.Stat(localFFprobe); err == nil {
		return localFFprobe, nil
	}

	// Check in parent directory
	parentFFprobe := filepath.Join(currentDir, "..", "ffprobe"+exeSuffix)
	if _, err := os.Stat(parentFFprobe); err == nil {
		return parentFFprobe, nil
	}

	// Then check in PATH
	ffprobePath, err := exec.LookPath("ffprobe")
	if err == nil {
		return ffprobePath, nil
	}

	return "", fmt.Errorf("ffprobe not found")
}

// downloadFile downloads a file from a URL and saves it to disk
func downloadFile(url, filename string) error {
	fmt.Printf("Downloading FFmpeg from %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create download file: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("error during download: %v", err)
	}

	return nil
}

// extractWindows extracts ffmpeg.exe and ffprobe.exe from a zip archive
func extractWindows(filename string) error {
	archive, err := zip.OpenReader(filename)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer archive.Close()

	// Extract both ffmpeg.exe and ffprobe.exe
	for _, f := range archive.File {
		var dstName string
		if strings.HasSuffix(f.Name, "bin/ffmpeg.exe") || strings.HasSuffix(f.Name, "bin\\ffmpeg.exe") {
			dstName = "ffmpeg.exe"
		} else if strings.HasSuffix(f.Name, "bin/ffprobe.exe") || strings.HasSuffix(f.Name, "bin\\ffprobe.exe") {
			dstName = "ffprobe.exe"
		} else {
			continue
		}

		dstFile, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create %s: %v", dstName, err)
		}
		defer dstFile.Close()

		srcFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open %s from zip: %v", dstName, err)
		}
		defer srcFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return fmt.Errorf("failed to extract %s: %v", dstName, err)
		}
	}
	return nil
}

// extractLinux extracts ffmpeg and ffprobe from a tar.xz archive
func extractLinux(filename string) error {
	// Extract the tar.xz file
	cmd := exec.Command("tar", "xf", filename)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extraction failed: %v", err)
	}

	// Walk through extracted directory to find ffmpeg and ffprobe binaries
	binaries := map[string]string{
		"ffmpeg":  "",
		"ffprobe": "",
	}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "ffmpeg" && strings.Contains(path, "bin/ffmpeg") {
			binaries["ffmpeg"] = path
		} else if info.Name() == "ffprobe" && strings.Contains(path, "bin/ffprobe") {
			binaries["ffprobe"] = path
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to find binaries: %v", err)
	}

	// Move binaries to current directory
	for name, path := range binaries {
		if path == "" {
			return fmt.Errorf("%s binary not found in extracted files", name)
		}
		if err := os.Rename(path, name); err != nil {
			return fmt.Errorf("failed to move %s binary: %v", name, err)
		}
		if err := os.Chmod(name, 0755); err != nil {
			return fmt.Errorf("failed to make %s executable: %v", name, err)
		}
	}

	return nil
}

// DownloadFFmpeg downloads ffmpeg binary based on OS
func DownloadFFmpeg() error {
	var filename string
	var ffmpegDir string

	url := "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/"
	switch runtime.GOOS {
	case "windows":
		url += "ffmpeg-n7.1-latest-win64-gpl-7.1.zip"
		filename = "ffmpeg.zip"
		ffmpegDir = "ffmpeg-n7.1-latest-win64-gpl-7.1"
	case "linux":
		url += "ffmpeg-n7.1-latest-linux64-gpl-7.1.tar.xz"
		filename = "ffmpeg.tar.xz"
		ffmpegDir = "ffmpeg-n7.1-latest-linux64-gpl-7.1"
	case "darwin":
		return fmt.Errorf("macOS builds are no longer available")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Download and extract
	if err := downloadFile(url, filename); err != nil {
		return err
	}
	fmt.Println("Download complete, extracting...")

	defer os.Remove(filename)
	defer os.RemoveAll(ffmpegDir)

	// Extract based on OS
	if runtime.GOOS == "windows" {
		if err := extractWindows(filename); err != nil {
			return err
		}
	} else if runtime.GOOS == "linux" {
		if err := extractLinux(filename); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	fmt.Println("FFmpeg installation complete!")
	return nil
}
