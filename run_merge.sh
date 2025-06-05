#!/usr/bin/env bash
set -euo pipefail

# Read user input
read -rp "Enter input file or directory name: " INPUT_FILE
read -rp "Enter output file name: " OUTPUT_FILE

# Default parameters
MODE="merge"
FPS=32
RESOLUTION=""
BITRATE=0
PRESET="p3"
CQ=32
WIDTH=0      # If width and height are set (>0), resolution will be ignored
HEIGHT=0
ENCODER="gpu"
# mp4, mkv, avi, flv, webm, mov, wmv, ts
OUTPUT_EXTENSION="mov"

# Show current parameters
echo ================================================
echo "Current parameters:"
echo "Input file: $INPUT_FILE"
echo "Output file: $OUTPUT_FILE"
echo "Mode: $MODE"
echo "Frame rate: $FPS"
echo "Resolution: $RESOLUTION"
echo "Bitrate: $BITRATE"
echo "Preset: $PRESET"
echo "Quality value: $CQ"
echo "Width: $WIDTH"
echo "Height: $HEIGHT"
echo "Encoder: $ENCODER"
echo "Output extension: $OUTPUT_EXTENSION"
echo ================================================

# Check if input file or directory exists
if [[ ! -e "$INPUT_FILE" ]]; then
  echo "Error: Input '$INPUT_FILE' not found."
  exit 1
fi

# Create output directory
OUTPUT_DIR=$(dirname "$OUTPUT_FILE")
mkdir -p "$OUTPUT_DIR"

# Run Go program
if ./video_compressor \
  -input      "$INPUT_FILE" \
  -output     "$OUTPUT_FILE" \
  -mode       "$MODE" \
  -fps        "$FPS" \
  -resolution "$RESOLUTION" \
  -bitrate    "$BITRATE" \
  -preset     "$PRESET" \
  -cq         "$CQ" \
  -width      "$WIDTH" \
  -height     "$HEIGHT" \
  -encoder    "$ENCODER" \
  -output-extension  "$OUTPUT_EXTENSION"
then
  echo "Video compression completed!"
else
  echo "Video compression failed!"
fi

# Wait for user to press any key before exiting
read -n1 -r -p "Press any key to continue..." _
echo
