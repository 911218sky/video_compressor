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
FPS=32
RESOLUTION="1080p"
BITRATE=0
PRESET="p3"
CQ=32
WIDTH=0
HEIGHT=0
ENCODER="gpu"

# Display current parameters
echo -e "${GREEN}Current parameters:${NC}"
echo "Input file: $INPUT_FILE"
echo "Output file: $OUTPUT_FILE"
echo "Frame rate: $FPS"
echo "Resolution: $RESOLUTION"
echo "Bitrate: $BITRATE"
echo "Preset: $PRESET"
echo "Quality value: $CQ"
echo "Width: $WIDTH"
echo "Height: $HEIGHT"
echo "Encoder: $ENCODER"

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

# Run video compression
./video_compressor -input "$INPUT_FILE" -output "$OUTPUT_FILE" -fps "$FPS" -resolution "$RESOLUTION" -bitrate "$BITRATE" -preset "$PRESET" -cq "$CQ" -width "$WIDTH" -height "$HEIGHT" -encoder "$ENCODER"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Video compression completed!${NC}"
else
    echo -e "${RED}Video compression failed!${NC}"
    exit 1
fi