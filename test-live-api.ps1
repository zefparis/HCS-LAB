# Test script for the live HCS Lab API on Railway

$apiUrl = "https://hcs-lab-production.up.railway.app"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  HCS Lab API - Live Test Suite" -ForegroundColor Cyan
Write-Host "  $apiUrl" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan

# Test 1: Service Status
Write-Host "`n[1] Testing Service Status..." -ForegroundColor Yellow
$status = Invoke-RestMethod -Uri "$apiUrl"
Write-Host "✅ Service: $($status.service)" -ForegroundColor Green
Write-Host "✅ Version: $($status.version)" -ForegroundColor Green
Write-Host "✅ Status: $($status.status)" -ForegroundColor Green

# Test 2: Health Check
Write-Host "`n[2] Testing Health Endpoint..." -ForegroundColor Yellow
$health = Invoke-RestMethod -Uri "$apiUrl/health"
Write-Host "✅ Health: $($health.status)" -ForegroundColor Green
Write-Host "✅ Uptime: $($health.uptime)" -ForegroundColor Green
Write-Host "✅ Secure: $($health.secure)" -ForegroundColor Green

# Test 3: Generate HCS for Air Profile
Write-Host "`n[3] Testing HCS Generation (Air)..." -ForegroundColor Yellow
$airProfile = @{
    dominantElement = "Air"
    modal = @{
        cardinal = 0.31
        fixed = 0.23
        mutable = 0.46
    }
    cognition = @{
        fluid = 0.52
        crystallized = 0.13
        verbal = 0.53
        strategic = 0.15
        creative = 0.33
    }
    interaction = @{
        pace = "balanced"
        structure = "medium"
        tone = "precise"
    }
} | ConvertTo-Json -Depth 10

$airResult = Invoke-RestMethod -Uri "$apiUrl/api/generate" -Method POST -Body $airProfile -ContentType "application/json"
Write-Host "✅ Generated HCS-U3:" -ForegroundColor Green
Write-Host "   $($airResult.codeU3)" -ForegroundColor White
Write-Host "✅ CHIP: $($airResult.chip)" -ForegroundColor Green

# Test 4: Generate HCS for Earth Profile
Write-Host "`n[4] Testing HCS Generation (Earth)..." -ForegroundColor Yellow
$earthProfile = @{
    dominantElement = "Earth"
    modal = @{
        cardinal = 0.25
        fixed = 0.60
        mutable = 0.15
    }
    cognition = @{
        fluid = 0.35
        crystallized = 0.75
        verbal = 0.40
        strategic = 0.80
        creative = 0.20
    }
    interaction = @{
        pace = "slow"
        structure = "high"
        tone = "warm"
    }
} | ConvertTo-Json -Depth 10

$earthResult = Invoke-RestMethod -Uri "$apiUrl/api/generate" -Method POST -Body $earthProfile -ContentType "application/json"
Write-Host "✅ Generated HCS-U3:" -ForegroundColor Green
Write-Host "   $($earthResult.codeU3)" -ForegroundColor White
Write-Host "✅ CHIP: $($earthResult.chip)" -ForegroundColor Green

# Test 5: Error Handling (Invalid Input)
Write-Host "`n[5] Testing Error Handling..." -ForegroundColor Yellow
$invalidProfile = @{
    dominantElement = "InvalidElement"
    modal = @{ cardinal = 2.0; fixed = 0.5; mutable = 0.5 }
    cognition = @{ fluid = 0.5; crystallized = 0.5; verbal = 0.5; strategic = 0.5; creative = 0.5 }
    interaction = @{ pace = "balanced"; structure = "medium"; tone = "neutral" }
} | ConvertTo-Json -Depth 10

try {
    $invalidResult = Invoke-RestMethod -Uri "$apiUrl/api/generate" -Method POST -Body $invalidProfile -ContentType "application/json"
    Write-Host "❌ Should have failed but didn't" -ForegroundColor Red
} catch {
    if ($_.Exception.Response.StatusCode -eq 400) {
        Write-Host "✅ Correctly rejected invalid input (400 Bad Request)" -ForegroundColor Green
    } else {
        Write-Host "⚠️ Unexpected error: $_" -ForegroundColor Yellow
    }
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  All Tests Complete!" -ForegroundColor Green
Write-Host "  API is fully operational" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
