@echo off
setlocal EnableDelayedExpansion

:: Read user input
set /p INPUT_FILE=Enter input file or directory name: 
set /p OUTPUT_FILE=Enter output file name: 

:: Default parameters
set MODE=merge
set FPS=32
set RESOLUTION=
set BITRATE=0
set PRESET=p3
set CQ=16
@REM 0 is auto (If width and height are set >0, resolution will be ignored)
set WIDTH=0
set HEIGHT=0
@REM gpu, cpu
set ENCODER=gpu
@REM mp4, mkv, avi, flv, webm, mov, wmv, ts
set OUTPUT_EXTENSION=mp4

:: Show current parameters
echo ================================================
echo Current parameters:
echo Input file:           %INPUT_FILE%
echo Output file:          %OUTPUT_FILE%
echo Mode:                 %MODE%
echo Frame rate:           %FPS%
echo Resolution:           %RESOLUTION%
echo Bitrate:              %BITRATE%
echo Preset:               %PRESET%
echo Quality value (CQ):   %CQ%
echo Width:                %WIDTH%
echo Height:               %HEIGHT%
echo Encoder:              %ENCODER%
echo Output extension:     %OUTPUT_EXTENSION%
echo ================================================

:: Check if input file or directory exists
if not exist "%INPUT_FILE%" (
    echo Error: Input "%INPUT_FILE%" not found.
    pause
    exit /b 1
)

:: Create output directory if it does not exist
for %%I in ("%OUTPUT_FILE%") do set "OUTPUT_DIR=%%~dpI"
if not exist "%OUTPUT_DIR%" (
    mkdir "%OUTPUT_DIR%"
)

:: Run Go program
go run src\main.go ^
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

:: Show message based on result
if %ERRORLEVEL% equ 0 (
    echo Video compression completed!
) else (
    echo Video compression failed!
)

:: Pause and wait for user to press any key
pause
