#!/bin/bash

echo "Building Nebula for production..."

echo ""
echo "Building frontend..."
cd web
pnpm build
if [ $? -ne 0 ]; then
    echo "Frontend build failed!"
    exit 1
fi

cd ..

echo ""
echo "Building backend..."
go build -o nebula-server ./cmd/nebula-server
if [ $? -ne 0 ]; then
    echo "Backend build failed!"
    exit 1
fi

echo ""
echo "========================================"
echo "Build completed successfully!"
echo "========================================"
echo ""
echo "To run in production mode:"
echo "  SERVER_MODE=prod ./nebula-server"
echo ""
echo "Or use the start script:"
echo "  ./scripts/start.sh"
