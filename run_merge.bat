@echo off
setlocal enabledelayedexpansion

:: Get user input for input and output files
set /p INPUT_FILE="Enter input directory name: "
set /p OUTPUT_FILE="Enter output file name: "

:: Set default parameters
set MODE=merge
set FPS=32
set RESOLUTION=""
set BITRATE=0
set PRESET=p3
set CQ=32
:: If width and height are set, resolution will be ignored (0 is auto)
set WIDTH=0
set HEIGHT=0
set ENCODER=gpu

:: Display current parameters
echo ================================================
echo Current parameters:
echo Input file: %INPUT_FILE%
echo Output file: %OUTPUT_FILE%
echo Mode: %MODE%
echo Frame rate: %FPS%
echo Resolution: %RESOLUTION%
echo Bitrate: %BITRATE%
echo Preset: %PRESET%
echo Quality value: %CQ%
echo Width: %WIDTH%
echo Height: %HEIGHT%
echo Encoder: %ENCODER%
echo ================================================

:: Check if input file exists
if not exist "%INPUT_FILE%" (
    echo Error: Input file %INPUT_FILE% not found
    exit /b 1
)

:: Create output directory
for %%I in ("%OUTPUT_FILE%") do set "OUTPUT_DIR=%%~dpI"
if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

:: Run video compression
video_compressor.exe -input "%INPUT_FILE%" -output "%OUTPUT_FILE%" -mode %MODE% -fps %FPS% -resolution %RESOLUTION% -bitrate %BITRATE% -preset %PRESET% -cq %CQ% -width %WIDTH% -height %HEIGHT% -encoder %ENCODER%

if %ERRORLEVEL% equ 0 (
    echo Video compression completed!
    pause
) else (
    echo Video compression failed!
    pause
)