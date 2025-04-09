# Video Compressor

A powerful video compression tool that supports both GPU (NVIDIA) and CPU encoding.

## Features

- Video compression with configurable quality settings
- Support for NVIDIA GPU acceleration (HEVC/H.265)
- Automatic fallback to CPU encoding if GPU is not available
- Multiple resolution presets (1080p, 720p, 480p)
- Customizable compression parameters
- Cross-platform support (Windows, Linux, macOS)
- Automatic FFmpeg installation if not present

## Requirements

- Go 1.22 or later
- FFmpeg (automatically downloaded if not present)
- NVIDIA GPU (optional, for hardware acceleration)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/video_compressor.git
cd video_compressor
```

2. Build the project:
```bash
# On Linux/macOS
./build.sh

# On Windows
build.bat
```

## Usage

### Basic Usage

```bash
# On Linux/macOS
./run.sh input.mp4 output.mp4

# On Windows
run.bat input.mp4 output.mp4
```

### Advanced Usage

```bash
./video_compressor -input input.mp4 -output output.mp4 -fps 30 -resolution 720p -bitrate 2500 -preset p7 -cq 32
```

### Parameters

- `-input`: Input video file path
- `-output`: Output video file path
- `-fps`: Frame rate (default: 32)
- `-resolution`: Video resolution (1080p, 720p, 480p)
- `-bitrate`: Custom bitrate in Kbps (0 for default)
- `-preset`: Encoder preset (p1=fastest, p7=best quality)
- `-cq`: Constant quality value (0-51, lower is better)
- `-width`: Custom width (0 for default)
- `-height`: Custom height (0 for default)

## Supported Video Formats

- Input: MP4, AVI, MKV, MOV, WMV, FLV
- Output: MP4 (H.265/HEVC)