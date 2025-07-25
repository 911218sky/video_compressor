@echo off
setlocal enabledelayedexpansion

:: Default settings
set "USE_UPX=true"
set "UPX_LEVEL=6"

:: Parse command line arguments
:parse_args
if "%~1"=="" goto start_build
if "%~1"=="--no-upx" (
    set "USE_UPX=false"
    shift
    goto parse_args
)
if "%~1"=="--upx-level" (
    set "UPX_LEVEL=%~2"
    shift
    shift
    goto parse_args
)
shift
goto parse_args

:start_build
echo Starting to compile video compressor...

:: Create build directory
if not exist build mkdir build

:: Check if UPX is installed if enabled
if "%USE_UPX%"=="true" (
    upx --version >nul 2>&1
    if errorlevel 1 (
        echo UPX is not installed. Please install UPX first or use --no-upx flag.
        exit /b 1
    )
)

:: Compile all versions
call :compile_and_compress "Windows" "video_compressor.exe" "windows"
if errorlevel 1 exit /b 1

call :compile_and_compress "Linux" "video_compressor" "linux"
if errorlevel 1 exit /b 1

call :compile_and_compress "macOS" "video_compressor_mac" "darwin"
if errorlevel 1 exit /b 1

echo All versions compiled successfully!
echo Binary sizes:
dir build /Q

if "%USE_UPX%"=="true" (
    echo UPX compressed binary sizes:
    dir build\upx /Q
)

goto :eof

:compile_and_compress
set "os=%~1"
set "output_name=%~2"
set "goos=%~3"

echo %output_name% | findstr /i "\.exe$" >nul
if errorlevel 1 (
    set "mini_output_name=%output_name%-mini"
) else (
    for %%f in (%output_name%) do set "mini_output_name=%%~nf-mini.exe"
)

echo Compiling %os% version...
set GOOS=%goos%
set GOARCH=amd64
go build -ldflags="-s -w" -o "build\%output_name%" .\src

if errorlevel 1 (
    echo %os% version compilation failed!
    exit /b 1
)

echo %os% version compiled successfully!

if "%USE_UPX%"=="true" (
    if not exist build\upx mkdir build\upx
    copy "build\%output_name%" "build\upx\%mini_output_name%" >nul
    echo Compressing %os% binary with UPX...
    upx --best --lzma -%UPX_LEVEL% "build\upx\%mini_output_name%"
)

goto :eof
