@echo off
setlocal EnableDelayedExpansion

REM Set color output
REM Define ESC (escape) character
for /F "delims=" %%A in ('echo prompt $E ^| cmd') do set "ESC=%%A"
set "GREEN=!ESC![0;32m"
set "RED=!ESC![0;31m"
REM No Color
set "NC=!ESC![0m"

REM Get user input for input and output files
set /p INPUT_FILE=Enter input file name: 
set /p OUTPUT_FILE=Enter output file name: 

REM Set default parameters
set "MODE=compress"
set "FPS=32"
set "RESOLUTION=1080p"
set "BITRATE=0"
set "PRESET=p3"
set "CQ=32"
REM If width and height are set, resolution will be ignored (0 is auto)
set "WIDTH=0"
set "HEIGHT=0"
set "ENCODER=gpu"
REM mp4, mkv, avi, flv, webm, mov, wmv, ts
set "OUTPUT_EXTENSION=mp4"

REM Display current parameters
echo ================================================
echo(!GREEN!Current parameters:!NC!
echo Input file:         %INPUT_FILE%
echo Output file:        %OUTPUT_FILE%
echo Mode:               %MODE%
echo Frame rate:         %FPS%
echo Resolution:         %RESOLUTION%
echo Bitrate:            %BITRATE%
echo Preset:             %PRESET%
echo Quality value (CQ): %CQ%
echo Width:              %WIDTH%
echo Height:             %HEIGHT%
echo Encoder:            %ENCODER%
echo Output extension:   %OUTPUT_EXTENSION%
echo ================================================

REM Check if input file exists
if not exist "%INPUT_FILE%" (
    echo(!RED!Error: Input file %INPUT_FILE% not found!%NC%
    pause
    exit /b 1
)

REM Create output directory
for %%I in ("%OUTPUT_FILE%") do set "OUTPUT_DIR=%%~dpI"
if not exist "%OUTPUT_DIR%" (
    mkdir "%OUTPUT_DIR%"
)

REM Run video compressor
video_compressor.exe ^
    -input "%INPUT_FILE%" ^
    -output "%OUTPUT_FILE%" ^
    -mode "%MODE%" ^
    -fps %FPS% ^
    -resolution "%RESOLUTION%" ^
    -bitrate %BITRATE% ^
    -preset "%PRESET%" ^
    -cq %CQ% ^
    -width %WIDTH% ^
    -height %HEIGHT% ^
    -encoder "%ENCODER%" ^
    -output-extension "%OUTPUT_EXTENSION%"

REM Check exit status and display result
if %ERRORLEVEL% equ 0 (
    echo(!GREEN!Video compression completed!%NC%
) else (
    echo(!RED!Video compression failed!%NC%
)

REM Wait for user to press any key before exiting
pause
