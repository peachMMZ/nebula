#!/usr/bin/env bash
set -e
BINARY=nebula-server
OUT_DIR=dist
WEB_DIR=web

echo "Building frontend..."
(cd "$WEB_DIR" && pnpm run build)

[[ -d "$OUT_DIR" ]] || mkdir -p "$OUT_DIR"
SUFFIX=""
[[ "$GOOS" == "windows" ]] && SUFFIX=".exe"

echo "Building $BINARY..."
go build -o "$OUT_DIR/${BINARY}${SUFFIX}" ./cmd/nebula-server/
echo "Build done: $OUT_DIR/${BINARY}${SUFFIX}"
