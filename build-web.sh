#!/bin/bash

# Build script for WebAssembly deployment

set -e

echo "Building Minesweeper for WebAssembly..."

# Clean previous build
rm -f minesweeper.wasm wasm_exec.js

# Build the WebAssembly binary
GOOS=js GOARCH=wasm go build -o minesweeper.wasm .

# Copy the wasm_exec.js file from Go installation
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .

echo "Build complete! Files created:"
echo "  - minesweeper.wasm"
echo "  - wasm_exec.js"
echo ""
echo "To test locally, run: go run github.com/hajimehoshi/wasmserve@latest ."