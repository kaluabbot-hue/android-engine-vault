#!/bin/bash
mkdir -p production_builds
echo "=================================================="
echo "    INITIALIZING PRODUCTION CROSS-COMPILATION     "
echo "=================================================="

# 1. Android Native Target Compilation (Mobile/Tablets)
echo "[BUILD] Compiling binary for Android Architecture (ARM64)..."
GOOS=android GOARCH=arm64 go build -ldflags="-s -w" -o production_builds/core_optimizer_android_arm64 main.go

# 2. Windows Desktop Target Compilation (Laptops/Workstations)
echo "[BUILD] Compiling binary for Windows Architecture (x64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o production_builds/core_optimizer_windows_x64.exe main.go

echo "--------------------------------------------------"
echo "[SUCCESS] Stripped production payloads generated inside 'production_builds/' folder."
echo "=================================================="
