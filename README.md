# Video Compressor

A powerful video compression tool that supports both GPU (NVIDIA) and CPU encoding, significantly reducing video file sizes while maintaining good quality. The tool also includes batch processing and video merging capabilities.

## Features

- High-performance video compression with configurable quality settings
- Support for NVIDIA GPU acceleration (HEVC/H.265)
- Automatic fallback to CPU encoding if GPU is not available
- Multiple resolution presets (1080p, 720p, 480p)
- Customizable compression parameters
- Cross-platform support (Windows, Linux, macOS)
- Automatic FFmpeg installation if not present
- Batch processing support
- Video merging functionality
- Progress bar display during compression
- Estimated time remaining calculation

## Requirements

- Go 1.22 or later
- FFmpeg (automatically downloaded if not present)
- NVIDIA GPU (optional, for hardware acceleration)
- At least 4GB RAM
- 1GB free disk space for installation

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
```

## Usage

### Basic Usage (Single Video Compression)

Simply double-click `run.bat` or run it from command line:

```bash
run.bat
```

The script will prompt you for:
- Input video file name
- Output video file name

Default compression settings:
- Resolution: 1080p
- Frame rate: 32 fps
- Preset: p3 (balanced speed/quality)
- Quality value: 32
- GPU encoding (automatically falls back to CPU if GPU is unavailable)

### Video Merging and Batch Processing

The tool supports merging multiple videos from a directory. Use `run_merge.bat` (Windows) for this functionality:

1. Run the merge script:
```bash
run_merge.bat
```

2. When prompted:
   - Enter the input directory containing your videos
   - Enter the desired output filename

Default merge parameters:
- Frame rate: 32 fps
- Preset: p3 (balanced speed/quality)
- Quality value: 32
- Encoder: GPU (with automatic fallback to CPU)

### Advanced Usage

```bash
# Single video compression with custom parameters
./video_compressor -input input.mp4 -output output.mp4 -fps 30 -resolution 720p -bitrate 2500 -preset p7 -cq 32

# Video merging with custom parameters
./video_compressor -input "input_directory" -output "output.mp4" -mode merge -fps 30 -preset p4 -cq 28 -encoder gpu
```

### Parameters

#### Common Parameters
- `-input`: Input video file or directory path
- `-output`: Output video file path
- `-fps`: Frame rate (default: 32)
- `-resolution`: Video resolution (1080p, 720p, 480p)
- `-bitrate`: Custom bitrate in Kbps (0 for default)
- `-preset`: Encoder preset (p1=fastest, p7=best quality)
- `-cq`: Constant quality value (0-51, lower is better)
- `-width`: Custom width (0 for default)
- `-height`: Custom height (0 for default)
- `-encoder`: Encoding device (gpu/cpu)

#### Merge-specific Parameters
- `-mode`: Operation mode (merge for combining videos)

## Supported Video Formats

### Input Formats
- MP4 (H.264/AVC, H.265/HEVC)
- AVI
- MKV
- MOV
- WMV
- FLV

### Output Format
- MP4 (H.265/HEVC)

## Performance Tips

1. Use GPU encoding when available for best performance
2. Choose appropriate presets based on your needs:
   - p1-p3: Fast encoding, larger file size
   - p4-p5: Balanced
   - p6-p7: Best quality, slower encoding
3. For batch processing or merging:
   - Organize input videos in a dedicated directory
   - Ensure sufficient disk space for temporary files
   - Consider using lower quality presets for large batches
   - Videos will be processed in alphabetical order