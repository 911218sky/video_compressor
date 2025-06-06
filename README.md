# 🎬 Video Compressor

> High-performance video compression tool with GPU/CPU encoding support, dramatically reducing file sizes while maintaining excellent quality

## ✨ Key Features

- 🚀 **GPU Acceleration** - NVIDIA hardware encoding for lightning-fast processing
- 📁 **Batch Processing** - Process entire directories of videos at once
- 🎯 **Smart Compression** - Auto-optimized settings for best quality
- 🎨 **Color Output** - Clear success/error message display
- 🔄 **Video Merging** - Combine multiple videos into a single file

---

## 🚀 Quick Start

### Windows Users

1. **Download** → [GitHub Releases](https://github.com/911218sky/video_compressor/releases)
2. **Extract** → Ensure these files are in the same folder:
   ```
   📁 video_compressor/
   ├── 🔧 video_compressor.exe
   ├── 📝 run.bat
   └── 📝 run_merge.bat
   ```
3. **Run** → Double-click `run.bat` to start

### Linux/macOS Users

```bash
# Set permissions
chmod +x video_compressor run.sh run_merge.sh

# Start using
./run.sh
```

---

## 💡 Usage

### 🎯 Single Video Compression

| Platform | Command |
|----------|---------|
| Windows | `run.bat` |
| Linux/macOS | `./run.sh` |

**Process:**
1. Enter input video filename
2. Enter output filename  
3. Wait for compression to complete ✅

### 📁 Batch Merging

| Platform | Command |
|----------|---------|
| Windows | `run_merge.bat` |
| Linux/macOS | `./run_merge.sh` |

**Process:**
1. Enter directory path containing videos
2. Enter merged output filename
3. Wait for processing to complete ✅

---

## ⚙️ Command-Line Parameters

### 🎛️ Essential Parameters

| Parameter | Description | Default | Options/Examples |
|-----------|-------------|---------|------------------|
| `-input` | Input video file/directory path | **Required** | `video.mp4`, `./videos/` |
| `-output` | Output video file path | Auto-generated | `output.mp4` |
| `-reverse` | Reverse the order of the files to be merged | `false` | `true`, `false` |
| `-mode` | Operation mode | `compress` | `compress`, `merge` |

### 📹 Video Settings

| Parameter | Description | Default | Options/Examples |
|-----------|-------------|---------|------------------|
| `-fps` | Frame rate | `32` | `24`, `30`, `60` |
| `-resolution` | Video resolution | Auto (1080p) | `1080p`, `720p`, `480p` |
| `-width` | Custom width (overrides resolution) | `0` (auto) | `1920`, `1280` |
| `-height` | Custom height (overrides resolution) | `0` (auto) | `1080`, `720` |

### 🎚️ Quality Settings

| Parameter | Description | Default | Range/Options |
|-----------|-------------|---------|---------------|
| `-preset` | Encoder preset | `p7` | `p1` (fastest) ~ `p7` (best quality) |
| `-cq` | Constant quality value | `16` | `0` (best) ~ `51` (worst) |
| `-bitrate` | Custom bitrate in Kbps | `0` (auto) | `2000`, `5000`, `10000` |

### 🔧 Technical Settings

| Parameter | Description | Default | Options |
|-----------|-------------|---------|---------|
| `-encoder` | Encoding device | `gpu` | `gpu`, `cpu` |
| `-output-extension` | Output file format | `.mp4` | `.mp4`, `.avi`, `.mkv`, `.mov`, `.wmv`, `.flv`, `.webm`, `.ts` |

### 🏃‍♂️ Speed vs Quality

| Use Case | Recommended Settings |
|----------|---------------------|
| **Speed Priority** | `preset p1` + `cq 35` |
| **Balanced** | `preset p3` + `cq 32` |
| **Quality Priority** | `preset p7` + `cq 20` |

### 💾 File Size Optimization

| File Size Target | CQ Value | Notes |
|------------------|----------|-------|
| **Small files** | `30-35` | Good for sharing |
| **Medium files** | `25-30` | Balanced quality |
| **Large files** | `18-25` | High quality |

---

### 📥 Input Formats
`MP4` `AVI` `MKV` `MOV` `WMV` `FLV` `TS` `WEBM`

### 📤 Output Formats  
| Format | Extension | Best For |
|--------|-----------|----------|
| **MP4** | `.mp4` | General use, compatibility |
| **MKV** | `.mkv` | High quality, multiple tracks |
| **AVI** | `.avi` | Legacy compatibility |
| **MOV** | `.mov` | Apple ecosystem |
| **WEBM** | `.webm` | Web streaming |
| **TS** | `.ts` | Broadcasting, streaming |
| **WMV** | `.wmv` | Windows compatibility |
| **FLV** | `.flv` | Flash video (legacy) |

---

## 🔧 Troubleshooting

| Issue | Solution |
|-------|----------|
| GPU encoding fails | Automatically switches to CPU |
| Input file not found | Check file path and permissions |
| FFmpeg missing | Tool auto-downloads FFmpeg |


**🚀 Start compressing! Enjoy efficient, high-quality video processing**