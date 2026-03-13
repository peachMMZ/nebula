@echo off
setlocal
set BINARY=nebula-server
set OUT_DIR=dist

echo Building Nebula Backend for production...

if not exist %OUT_DIR% mkdir %OUT_DIR%
echo Building %BINARY%...
go build -ldflags="-s -w" -o %OUT_DIR%\%BINARY%.exe ./cmd/nebula-server/
if %ERRORLEVEL% neq 0 (
    echo Build failed.
    exit /b 1
)

echo.
echo ========================================
echo Build completed successfully!
echo ========================================
echo Output: %OUT_DIR%\%BINARY%.exe
echo.
echo To run:
echo   set SERVER_MODE=prod
echo   .\%OUT_DIR%\%BINARY%.exe
endlocal
