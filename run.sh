#!/bin/bash

# Set color output
GREEN='\033[0;32m'
RED='\033[0;31m'
# No Color
NC='\033[0m'

# Get user input for input and output files
read -p "Enter input file name: " INPUT_FILE
read -p "Enter output file name: " OUTPUT_FILE

# Set default parameters
MODE="compress"
FPS=32
RESOLUTION="1080p"
BITRATE=0
PRESET="p3"
CQ=16
# If width and height are set, resolution will be ignored (0 is auto)
WIDTH=0
HEIGHT=0
ENCODER="gpu"
# mp4, mkv, avi, flv, webm, mov, wmv, ts
OUTPUT_EXTENSION="mp4"

# Display current parameters
echo ================================================
echo -e "${GREEN}Current parameters:${NC}"
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

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo -e "${RED}Error: Input file $INPUT_FILE not found${NC}"
    exit 1
fi

# Create output directory
OUTPUT_DIR=$(dirname "$OUTPUT_FILE")
if [ ! -d "$OUTPUT_DIR" ]; then
    mkdir -p "$OUTPUT_DIR"
fi

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
    echo -e "${GREEN}Video compression completed!${NC}"
else
    echo -e "${RED}Video compression failed!${NC}"
    exit 1
fi