#!/bin/bash

# Quick restart script for TG Go Music Bot

echo "🔄 Restarting TG Go Music Bot..."

# Kill existing process
pkill -f "tgmusicbot" 2>/dev/null
pkill -f "go run main.go" 2>/dev/null

# Wait a moment
sleep 2

# Start the bot
if [ -f "./tgmusicbot" ]; then
    nohup ./tgmusicbot > bot.log 2>&1 &
    echo "✅ Bot started (using compiled binary)"
else
    nohup go run main.go > bot.log 2>&1 &
    echo "✅ Bot started (using go run)"
fi

echo "📝 Logs: tail -f bot.log"
echo "🛑 Stop: ./stop.sh"
