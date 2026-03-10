@echo off
setlocal
set BINARY=nebula-server
set OUT_DIR=dist
set WEB_DIR=web

echo Building frontend...
cd %WEB_DIR%
call pnpm run build
if %ERRORLEVEL% neq 0 (
    echo Frontend build failed.
    exit /b 1
)
cd ..

if not exist %OUT_DIR% mkdir %OUT_DIR%
echo Building %BINARY%...
go build -o %OUT_DIR%\%BINARY%.exe ./cmd/nebula-server/
if %ERRORLEVEL% neq 0 (
    echo Build failed.
    exit /b 1
)
echo Build done: %OUT_DIR%\%BINARY%.exe
endlocal
