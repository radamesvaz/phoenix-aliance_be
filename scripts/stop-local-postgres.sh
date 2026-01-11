#!/bin/bash

# Script to help stop local PostgreSQL service on Windows
# This script provides instructions and attempts to stop the service

echo "Attempting to stop local PostgreSQL service..."
echo ""

# Try to find PostgreSQL service using sc command (Windows)
if command -v sc &> /dev/null; then
    echo "Found PostgreSQL services:"
    sc query | grep -i postgresql || echo "No PostgreSQL services found via 'sc query'"
    echo ""
    
    # Try common service names
    SERVICE_NAMES=(
        "postgresql-x64-16"
        "postgresql-x64-15"
        "postgresql-x64-14"
        "postgresql-x64-13"
        "postgresql-x64-12"
    )
    
    for service in "${SERVICE_NAMES[@]}"; do
        if sc query "$service" &> /dev/null; then
            echo "Found service: $service"
            echo "Attempting to stop..."
            if sc stop "$service" 2>/dev/null; then
                echo "✓ Stopped: $service"
            else
                echo "✗ Failed to stop: $service (may need Administrator privileges)"
            fi
        fi
    done
else
    echo "Note: 'sc' command not available in this shell"
    echo ""
fi

echo ""
echo "If the above didn't work, try one of these methods:"
echo ""
echo "1. PowerShell (as Administrator):"
echo "   Get-Service postgresql* | Stop-Service"
echo ""
echo "2. Command Prompt (as Administrator):"
echo "   net stop postgresql-x64-16"
echo "   (Replace '16' with your PostgreSQL version)"
echo ""
echo "3. Services GUI:"
echo "   - Press Win+R, type 'services.msc'"
echo "   - Find 'postgresql-x64-*' service"
echo "   - Right-click → Stop"
echo ""












