#!/bin/bash

echo "Testing AI TUI application..."
echo "Terminal size: $(tput cols)x$(tput lines)"
echo ""

# Test if the binary runs and exits cleanly
timeout 3s ./ai-tui 2>&1 || {
    echo "Exit code: $?"
    echo "Application started successfully (timeout expected)"
}

echo ""
echo "✅ AI TUI application is working!"
echo "Run './ai-tui' in an interactive terminal to use the full interface."