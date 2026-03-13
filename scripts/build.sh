#!/usr/bin/env bash
set -e
BINARY=nebula-server
OUT_DIR=dist

echo "Building Nebula Backend for production..."

[[ -d "$OUT_DIR" ]] || mkdir -p "$OUT_DIR"
SUFFIX=""
[[ "$GOOS" == "windows" ]] && SUFFIX=".exe"

echo "Building $BINARY..."
go build -ldflags="-s -w" -o "$OUT_DIR/${BINARY}${SUFFIX}" ./cmd/nebula-server/

echo ""
echo "========================================"
echo "Build completed successfully!"
echo "========================================"
echo "Output: $OUT_DIR/${BINARY}${SUFFIX}"
echo ""
echo "To run:"
echo "  SERVER_MODE=prod ./$OUT_DIR/${BINARY}${SUFFIX}"
