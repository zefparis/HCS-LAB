# Build script for HCS Lab API

Write-Host "Building HCS Lab API..." -ForegroundColor Cyan

# Download dependencies
Write-Host "Downloading dependencies..." -ForegroundColor Yellow
go mod download

# Run tests
Write-Host "`nRunning tests..." -ForegroundColor Yellow
go test ./tests/... -v

if ($LASTEXITCODE -ne 0) {
    Write-Host "`nTests failed! Build aborted." -ForegroundColor Red
    exit 1
}

Write-Host "`nTests passed!" -ForegroundColor Green

# Build CLI tool
Write-Host "`nBuilding hcsgen CLI..." -ForegroundColor Yellow
go build -o hcsgen.exe ./cmd/hcsgen

if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build hcsgen!" -ForegroundColor Red
    exit 1
}

# Build API server
Write-Host "Building hcsapi server..." -ForegroundColor Yellow
go build -o hcsapi.exe ./cmd/hcsapi

if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build hcsapi!" -ForegroundColor Red
    exit 1
}

Write-Host "`nBuild successful!" -ForegroundColor Green
Write-Host "Executables created:" -ForegroundColor Cyan
Write-Host "  - hcsgen.exe (CLI tool)" -ForegroundColor White
Write-Host "  - hcsapi.exe (HTTP API server)" -ForegroundColor White
