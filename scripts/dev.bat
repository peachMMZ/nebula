@echo off
echo Starting Nebula development environment...

echo.
echo Starting backend server (dev mode)...
start "Nebula Backend" cmd /c "set SERVER_MODE=dev && go run ./cmd/nebula-server/main.go"

echo Waiting for backend to start...
timeout /t 3 /nobreak >nul

echo.
echo Starting frontend (Vite dev server)...
cd web
pnpm dev