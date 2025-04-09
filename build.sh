#!/bin/bash

# Set color output
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Default settings
USE_UPX=true
UPX_LEVEL=6
KEEP_UNCOMPRESSED=true

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --no-upx)
            USE_UPX=false
            shift
            ;;
        --upx-level)
            UPX_LEVEL=$2
            shift 2
            ;;
        --keep-uncompressed)
            KEEP_UNCOMPRESSED=true
            shift
            ;;
        *)
            shift
            ;;
    esac
done

echo -e "${GREEN}Starting to compile video compressor...${NC}"

# Create build directory
mkdir -p build

# Check if UPX is installed if enabled
if [ "$USE_UPX" = true ]; then
    if ! command -v upx &> /dev/null; then
        echo "UPX is not installed. Please install UPX first or use --no-upx flag."
        exit 1
    fi
fi

# Function to handle binary compilation and compression
compile_and_compress() {
    local os=$1
    local output_name=$2
    local goos=$3
    
    echo -e "${GREEN}Compiling ${os} version...${NC}"
    GOOS=$goos GOARCH=amd64 go build -ldflags="-s -w" -o "build/${output_name}" ./src
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}${os} version compiled successfully!${NC}"
        
        if [ "$USE_UPX" = true ]; then
            # Create upx directory if it doesn't exist
            mkdir -p build/upx
            
            # Copy to upx directory for compression
            cp "build/${output_name}" "build/upx/${output_name}"
            
            echo -e "${GREEN}Compressing ${os} binary with UPX...${NC}"
            upx --best --lzma -$UPX_LEVEL "build/upx/${output_name}"
        fi
    else
        echo "${os} version compilation failed!"
        exit 1
    fi
}

# Compile all versions
compile_and_compress "Windows" "video_compressor.exe" "windows"
compile_and_compress "Linux" "video_compressor" "linux"
compile_and_compress "macOS" "video_compressor_mac" "darwin"

echo -e "${GREEN}All versions compiled successfully!${NC}"
echo -e "${GREEN}Binary sizes:${NC}"
ls -lh build/
if [ "$USE_UPX" = true ]; then
    echo -e "${GREEN}UPX compressed binary sizes:${NC}"
    ls -lh build/upx/
fi