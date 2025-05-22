@echo off
echo VPNProxy Simulator Demo
echo =======================
echo.

echo Checking for certificates...
if not exist cert.pem (
    echo Generating certificates...
    go run main.go
    if %ERRORLEVEL% neq 0 (
        echo Failed to generate certificates!
        pause
        exit /b 1
    )
)

echo.
echo Starting proxy server in a new window...
start cmd /k "go run proxy/main.go --listen :8080"

echo Waiting for proxy to initialize...
timeout /t 3 /nobreak > nul

echo.
echo Starting client with visualization...
echo (The client will connect to example.com through the proxy)
echo.
go run client/main.go --target example.com

pause