#!/bin/bash

echo "Starting Nebula development environment..."

echo ""
echo "Starting backend server (dev mode)..."
SERVER_MODE=dev go run ./cmd/nebula-server/main.go &
BACKEND_PID=$!

echo "Backend started with PID: $BACKEND_PID"
echo "Waiting for backend to start..."
sleep 3

echo ""
echo "Starting frontend (Vite dev server)..."
cd web
pnpm dev

# Cleanup function to kill backend when script exits
cleanup() {
    echo ""
    echo "Shutting down backend server..."
    kill $BACKEND_PID 2>/dev/null
    exit
}

trap cleanup INT TERM
