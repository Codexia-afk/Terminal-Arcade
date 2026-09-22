#!/usr/bin/env bash
set -e

DIST_DIR="dist"
CMD_DIR="./cmd/arcade"
BINARY_NAME="arcade"
VERSION="${VERSION:-v1.0.0}"
FLAGS="-buildvcs=false -ldflags -s -ldflags -w -ldflags -X=main.version=${VERSION}"

echo "Creating distribution directory: $DIST_DIR"
mkdir -p "$DIST_DIR"

echo "1/5 Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build $FLAGS -o "$DIST_DIR/${BINARY_NAME}_linux_amd64" "$CMD_DIR"

echo "2/5 Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build $FLAGS -o "$DIST_DIR/${BINARY_NAME}_linux_arm64" "$CMD_DIR"

echo "3/5 Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build $FLAGS -o "$DIST_DIR/${BINARY_NAME}_darwin_amd64" "$CMD_DIR"

echo "4/5 Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build $FLAGS -o "$DIST_DIR/${BINARY_NAME}_darwin_arm64" "$CMD_DIR"

echo "5/5 Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build $FLAGS -o "$DIST_DIR/${BINARY_NAME}_windows_amd64.exe" "$CMD_DIR"

echo "Build complete! Release artifacts in $DIST_DIR:"
ls -lh "$DIST_DIR"
