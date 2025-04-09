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
	"strings"
)

// Supported video formats
var supportedFormats = map[string]bool{
	".mp4": true,
	".avi": true,
	".mkv": true,
	".mov": true,
	".wmv": true,
	".flv": true,
}

// IsSupportedFormat checks if the given file format is supported
func IsSupportedFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return supportedFormats[ext]
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

// extractWindows extracts ffmpeg.exe from a zip archive
func extractWindows(filename string) error {
	archive, err := zip.OpenReader(filename)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		// Check if the file is ffmpeg.exe in the bin directory
		if strings.HasSuffix(f.Name, "bin/ffmpeg.exe") || strings.HasSuffix(f.Name, "bin\\ffmpeg.exe") {
			dstFile, err := os.OpenFile("ffmpeg.exe", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("failed to create ffmpeg.exe: %v", err)
			}
			defer dstFile.Close()

			// create new file in current directory
			srcFile, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open ffmpeg.exe from zip: %v", err)
			}
			defer srcFile.Close()

			// copy the file
			if _, err := io.Copy(dstFile, srcFile); err != nil {
				return fmt.Errorf("failed to extract ffmpeg.exe: %v", err)
			}
			return nil
		}
	}
	return fmt.Errorf("ffmpeg.exe not found in zip archive")
}

// extractLinux extracts ffmpeg from a tar.xz archive
func extractLinux(filename string) error {
	// Extract the tar.xz file
	cmd := exec.Command("tar", "xf", filename)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extraction failed: %v", err)
	}

	// Walk through extracted directory to find ffmpeg binary
	var ffmpegPath string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "ffmpeg" && strings.Contains(path, "bin/ffmpeg") {
			ffmpegPath = path
			return nil
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to find ffmpeg binary: %v", err)
	}
	if ffmpegPath == "" {
		return fmt.Errorf("ffmpeg binary not found in extracted files")
	}

	// Move ffmpeg binary to current directory
	if err := os.Rename(ffmpegPath, "ffmpeg"); err != nil {
		return fmt.Errorf("failed to move ffmpeg binary: %v", err)
	}

	// Make ffmpeg executable
	if err := os.Chmod("ffmpeg", 0755); err != nil {
		return fmt.Errorf("failed to make ffmpeg executable: %v", err)
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
