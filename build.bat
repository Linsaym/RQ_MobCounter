@echo off
REM Build script for RQ_MobCounter
REM Windows batch file for building the application

setlocal enabledelayedexpansion

echo.
echo ========================================
echo RQ_MobCounter - Build Script
echo ========================================
echo.

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo Error: Go is not installed or not in PATH
    exit /b 1
)

echo [1/3] Downloading dependencies...
go mod tidy
if errorlevel 1 (
    echo Error: Failed to download dependencies
    exit /b 1
)

echo [2/3] Running tests...
powershell -ExecutionPolicy Bypass -File test.ps1
if errorlevel 1 (
    echo Error: Tests failed
    exit /b 1
)

echo [3/3] Building application...
go build -o build/RQ_MobCounter.exe
if errorlevel 1 (
    echo Error: Build failed
    exit /b 1
)

echo.
echo ========================================
echo Build completed successfully!
echo ========================================
echo.
echo Executable: build/RQ_MobCounter.exe
echo Usage: build/RQ_MobCounter.exe --help
echo.
pause
