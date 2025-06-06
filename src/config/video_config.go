package config

import (
	"fmt"
	"strings"
)

// Resolution represents supported video resolutions
type Resolution string

const (
	Resolution4K    Resolution = "4k"
	Resolution2K    Resolution = "2k"
	Resolution1080p Resolution = "1080p"
	Resolution720p  Resolution = "720p"
	Resolution480p  Resolution = "480p"
	Resolution360p  Resolution = "360p"
	Resolution240p  Resolution = "240p"
	ResolutionNone  Resolution = ""
)

// StringToResolution converts a string or int to Resolution type
func StringToResolution(s interface{}) (Resolution, error) {
	switch v := s.(type) {
	case string:
		switch strings.ToLower(v) {
		case "4k", "3840p", "2160p":
			return Resolution4K, nil
		case "2k", "1440p":
			return Resolution2K, nil
		case "1080p", "1080":
			return Resolution1080p, nil
		case "720p", "720":
			return Resolution720p, nil
		case "480p", "480":
			return Resolution480p, nil
		case "360p", "360":
			return Resolution360p, nil
		case "240p", "240":
			return Resolution240p, nil
		default:
			return ResolutionNone, nil
		}
	case int:
		switch v {
		case 3840, 2160:
			return Resolution4K, nil
		case 2560, 1440:
			return Resolution2K, nil
		case 1080:
			return Resolution1080p, nil
		case 720:
			return Resolution720p, nil
		case 480:
			return Resolution480p, nil
		case 360:
			return Resolution360p, nil
		case 240:
			return Resolution240p, nil
		default:
			return ResolutionNone, nil
		}
	}
	return "", fmt.Errorf("unsupported resolution: %v", s)
}

// VideoConfig holds all video compression parameters
type VideoConfig struct {
	FfmpegPath      string
	Fps             int
	Resolution      Resolution
	Bitrate         int
	Preset          string
	Cq              int
	Width           int    // Target width (0 means auto). If set, Resolution will be ignored
	Height          int    // Target height (0 means auto). If set, Resolution will be ignored
	Encoder         string // "gpu" for NVIDIA HEVC or "cpu" for libx265
	OutputExtension string // ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv"

	// Reverse the order of the files to be merged
	Reverse bool
}
