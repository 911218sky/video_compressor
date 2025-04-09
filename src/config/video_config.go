package config

// VideoConfig holds all video compression parameters
type VideoConfig struct {
	Fps        int
	Resolution string
	Bitrate    int
	Preset     string
	Cq         int
	Width      int
	Height     int
}

// GetResolutionSettings returns width and height based on resolution string
func GetResolutionSettings(resolution string) (width, height, bitrate int) {
	switch resolution {
	case "1080p":
		return 1920, 1080, 5000
	case "720p":
		return 1280, 720, 2500
	case "480p":
		return 854, 480, 1000
	default:
		return 1280, 720, 2500 // default to 720p
	}
} 