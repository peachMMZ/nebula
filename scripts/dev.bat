@echo off
echo Starting Nebula Backend (Dev Mode)...

echo.
echo Starting backend server...
echo API will be available at http://localhost:9050/api
echo.

set SERVER_MODE=dev
go run ./cmd/nebula-server/main.go

