#!/bin/bash

# Set color output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Set default parameters
INPUT_FILE="output_comparison_streaming_v1.mp4"
OUTPUT_FILE="video_compressed.mp4"
FPS=32
RESOLUTION="720p"
BITRATE=0
PRESET="p7"
CQ=32
WIDTH=0
HEIGHT=0

# Check if parameters are overridden
if [ ! -z "$1" ]; then INPUT_FILE="$1"; fi
if [ ! -z "$2" ]; then OUTPUT_FILE="$2"; fi
if [ ! -z "$3" ]; then FPS="$3"; fi
if [ ! -z "$4" ]; then RESOLUTION="$4"; fi
if [ ! -z "$5" ]; then BITRATE="$5"; fi
if [ ! -z "$6" ]; then PRESET="$6"; fi
if [ ! -z "$7" ]; then CQ="$7"; fi
if [ ! -z "$8" ]; then WIDTH="$8"; fi
if [ ! -z "$9" ]; then HEIGHT="$9"; fi

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

# Check if video_compressor exists in current directory
if [ ! -f "./video_compressor" ]; then
    echo -e "${RED}Error: video_compressor executable not found in current directory${NC}"
    exit 1
fi

# Run video compression
./video_compressor -input "$INPUT_FILE" -output "$OUTPUT_FILE" -fps "$FPS" -resolution "$RESOLUTION" -bitrate "$BITRATE" -preset "$PRESET" -cq "$CQ" -width "$WIDTH" -height "$HEIGHT"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Video compression completed!${NC}"
else
    echo -e "${RED}Video compression failed!${NC}"
    exit 1
fi 