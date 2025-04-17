#!/bin/bash

# Exit on error
set -e

# Create bin directory if it doesn't exist
mkdir -p bin

# Function to build a tool
build_tool() {
    local tool_name=$1
    local tool_dir=$2
    local output_name=$3

    echo "Building $tool_name..."
    if [ -f "$tool_dir/main.go" ] && [ -s "$tool_dir/main.go" ]; then
        go build -o "bin/$output_name" "./$tool_dir" || echo "Failed to build $tool_name"
        chmod +x "bin/$output_name"
    else
        echo "Skipping $tool_name - no valid main.go found"
    fi
}

# Build main binary
build_tool "email" "cmd/email" "email"

# Build individual tools
build_tool "list tool" "cmd/list" "list"
build_tool "send tool" "cmd/send" "send"
build_tool "login tool" "cmd/login" "login"

echo "Build complete! Binaries are in the bin/ directory" 