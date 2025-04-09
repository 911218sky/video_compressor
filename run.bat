@echo off
setlocal enabledelayedexpansion

:: Set default parameters
set INPUT_FILE=output_comparison_streaming_v1.mp4
set OUTPUT_FILE=video_compressed.mp4
set FPS=32
set RESOLUTION=720p
set BITRATE=0
set PRESET=p7
set CQ=32
set WIDTH=0
set HEIGHT=0

:: Check if parameters are overridden
if not "%1"=="" set INPUT_FILE=%1
if not "%2"=="" set OUTPUT_FILE=%2
if not "%3"=="" set FPS=%3
if not "%4"=="" set RESOLUTION=%4
if not "%5"=="" set BITRATE=%5
if not "%6"=="" set PRESET=%6
if not "%7"=="" set CQ=%7
if not "%8"=="" set WIDTH=%8
if not "%9"=="" set HEIGHT=%9

:: Display current parameters
echo Current parameters:
echo Input file: %INPUT_FILE%
echo Output file: %OUTPUT_FILE%
echo Frame rate: %FPS%
echo Resolution: %RESOLUTION%
echo Bitrate: %BITRATE%
echo Preset: %PRESET%
echo Quality value: %CQ%
echo Width: %WIDTH%
echo Height: %HEIGHT%

:: Check if input file exists
if not exist "%INPUT_FILE%" (
    echo Error: Input file %INPUT_FILE% not found
    exit /b 1
)

:: Create output directory
for %%I in ("%OUTPUT_FILE%") do set "OUTPUT_DIR=%%~dpI"
if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

:: Run video compression
video_compressor.exe -input "%INPUT_FILE%" -output "%OUTPUT_FILE%" -fps %FPS% -resolution %RESOLUTION% -bitrate %BITRATE% -preset %PRESET% -cq %CQ% -width %WIDTH% -height %HEIGHT%

if %ERRORLEVEL% equ 0 (
    echo Video compression completed!
    pause
) else (
    echo Video compression failed!
    pause
)