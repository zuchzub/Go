#!/bin/bash

# TG Go Music Bot - Auto Setup and Run Script
# This script will automatically install all required dependencies and run the bot

set -e  # Exit on error

echo "=================================================="
echo "  TG Go Music Bot - Auto Setup & Run"
echo "=================================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_info() {
    echo -e "${YELLOW}[!]${NC} $1"
}

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    SUDO=""
else 
    SUDO="sudo"
fi

# 1. Check and install FFmpeg
print_info "Checking FFmpeg..."
if ! command -v ffmpeg &> /dev/null; then
    print_info "FFmpeg not found. Installing..."
    $SUDO apt-get update -qq
    $SUDO apt-get install -y ffmpeg
    print_status "FFmpeg installed successfully"
else
    print_status "FFmpeg already installed"
fi

# 2. Check and install yt-dlp
print_info "Checking yt-dlp..."
if ! command -v yt-dlp &> /dev/null; then
    print_info "yt-dlp not found. Installing..."
    $SUDO wget -q -O /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux
    $SUDO chmod +x /usr/local/bin/yt-dlp
    print_status "yt-dlp installed successfully"
else
    print_status "yt-dlp already installed"
    # Update to latest version
    print_info "Updating yt-dlp to latest version..."
    $SUDO yt-dlp -U 2>/dev/null || $SUDO wget -q -O /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux
fi

# 3. Setup ntgcalls library
print_info "Checking ntgcalls library..."
if [ ! -f "pkg/vc/libntgcalls.so" ]; then
    print_info "ntgcalls library not found. Running setup..."
    go run setup_ntgcalls.go
    print_status "ntgcalls library downloaded"
fi

# Install libntgcalls.so to system
if [ ! -f "/usr/local/lib/libntgcalls.so" ]; then
    print_info "Installing libntgcalls.so to system..."
    $SUDO cp pkg/vc/libntgcalls.so /usr/local/lib/
    $SUDO ldconfig
    print_status "libntgcalls.so installed to system"
else
    print_status "libntgcalls.so already in system"
fi

# 4. Check Python and Pyrogram (for session generator)
print_info "Checking Python environment..."
if command -v python3 &> /dev/null; then
    if ! python3 -c "import pyrogram" 2>/dev/null; then
        print_info "Installing Pyrogram..."
        pip install -q pyrogram tgcrypto
        print_status "Pyrogram installed"
    else
        print_status "Pyrogram already installed"
    fi
else
    print_error "Python3 not found. Install Python3 to use session generator."
fi

# 5. Check .env file
print_info "Checking configuration..."
if [ ! -f ".env" ]; then
    if [ -f "sample.env" ]; then
        print_info "Creating .env from sample.env..."
        cp sample.env .env
        print_error "Please edit .env file with your credentials before running the bot!"
        exit 1
    else
        print_error ".env file not found. Please create one with required variables."
        exit 1
    fi
else
    print_status "Configuration file found"
fi

# 6. Create downloads directory
print_info "Creating downloads directory..."
mkdir -p database/music
print_status "Downloads directory ready"

# 7. Build the bot
print_info "Building the bot..."
if go build -o tgmusicbot main.go; then
    print_status "Bot built successfully"
else
    print_error "Failed to build the bot"
    exit 1
fi

# 8. Display system info
echo ""
echo "=================================================="
echo "  System Information"
echo "=================================================="
echo "FFmpeg version: $(ffmpeg -version | head -n1)"
echo "yt-dlp version: $(yt-dlp --version)"
echo "Go version: $(go version)"
echo "Python version: $(python3 --version 2>/dev/null || echo 'Not installed')"
echo "=================================================="
echo ""

# 9. Run the bot
print_status "All dependencies installed!"
print_info "Starting the bot..."
echo ""

# Run the bot with proper error handling
if [ -f "./tgmusicbot" ]; then
    ./tgmusicbot
else
    go run main.go
fi
