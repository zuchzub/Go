#!/bin/bash

# Stop script for TG Go Music Bot

echo "🛑 Stopping TG Go Music Bot..."

# Kill the bot process
pkill -f "tgmusicbot"
pkill -f "go run main.go"

sleep 1

# Check if still running
if pgrep -f "tgmusicbot" > /dev/null || pgrep -f "go run main.go" > /dev/null; then
    echo "⚠️  Force stopping..."
    pkill -9 -f "tgmusicbot"
    pkill -9 -f "go run main.go"
    sleep 1
fi

if pgrep -f "tgmusicbot" > /dev/null || pgrep -f "go run main.go" > /dev/null; then
    echo "❌ Failed to stop bot"
    exit 1
else
    echo "✅ Bot stopped successfully"
fi
