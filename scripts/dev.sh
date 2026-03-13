#!/bin/bash

echo "Starting Nebula Backend (Dev Mode)..."
echo ""
echo "API will be available at http://localhost:9050/api"
echo ""

SERVER_MODE=dev go run ./cmd/nebula-server/main.go
