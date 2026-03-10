@echo off
echo Building Nebula for production...

echo.
echo Building frontend...
cd web
call pnpm build
if errorlevel 1 (
    echo Frontend build failed!
    exit /b 1
)

cd ..

echo.
echo Building backend...
go build -o nebula-server.exe ./cmd/nebula-server
if errorlevel 1 (
    echo Backend build failed!
    exit /b 1
)

echo.
echo ========================================
echo Build completed successfully!
echo ========================================
echo.
echo To run in production mode:
echo   set SERVER_MODE=prod
echo   .\nebula-server.exe
echo.
echo Or use the start script:
echo   .\scripts\start.bat
