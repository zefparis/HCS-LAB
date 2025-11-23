# Quick git update script for HCS Lab API

param(
    [string]$Message = "Update HCS Lab API"
)

Write-Host "Git Update for HCS Lab API" -ForegroundColor Cyan
Write-Host "============================" -ForegroundColor Cyan

# Check git status
Write-Host "`nChecking status..." -ForegroundColor Yellow
git status --short

# Add all changes
Write-Host "`nAdding changes..." -ForegroundColor Yellow
git add .

# Commit with message
Write-Host "`nCommitting with message: $Message" -ForegroundColor Yellow
git commit -m "$Message"

if ($LASTEXITCODE -eq 0) {
    # Push to origin
    Write-Host "`nPushing to GitHub..." -ForegroundColor Yellow
    git push origin main
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`nSuccess! Changes pushed to GitHub." -ForegroundColor Green
        Write-Host "Repository: https://github.com/zefparis/HCS-LAB" -ForegroundColor Cyan
    } else {
        Write-Host "`nPush failed!" -ForegroundColor Red
    }
} else {
    Write-Host "`nNo changes to commit or commit failed." -ForegroundColor Yellow
}
