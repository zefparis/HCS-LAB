# Run script for HCS Lab API Server

param(
    [int]$Port = 8080
)

Write-Host "Starting HCS Lab API Server..." -ForegroundColor Cyan
Write-Host "Port: $Port" -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop" -ForegroundColor Gray
Write-Host ""

# Set environment variable
$env:PORT = $Port

# Check if binary exists
if (-not (Test-Path "./hcsapi.exe")) {
    Write-Host "hcsapi.exe not found. Building..." -ForegroundColor Yellow
    go build -o hcsapi.exe ./cmd/hcsapi
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Build failed!" -ForegroundColor Red
        exit 1
    }
}

# Run the server
./hcsapi.exe
