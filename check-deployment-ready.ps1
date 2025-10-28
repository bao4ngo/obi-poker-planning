# Deployment Helper Script
# This script helps verify your setup before deployment

Write-Host "🚀 Poker Planning - Deployment Readiness Check" -ForegroundColor Cyan
Write-Host "=" * 50

# Check if required tools are installed
Write-Host "`n📦 Checking required tools..." -ForegroundColor Yellow

$tools = @{
    "go" = "Go (Backend)"
    "node" = "Node.js (Frontend)"
    "npm" = "NPM (Package Manager)"
    "fly" = "Fly CLI (Backend Deployment)"
    "vercel" = "Vercel CLI (Frontend Deployment - Optional)"
    "psql" = "PostgreSQL (Database - Local)"
}

$missingTools = @()

foreach ($tool in $tools.Keys) {
    try {
        $null = Get-Command $tool -ErrorAction Stop
        Write-Host "  ✅ $($tools[$tool])" -ForegroundColor Green
    } catch {
        Write-Host "  ❌ $($tools[$tool]) - NOT FOUND" -ForegroundColor Red
        $missingTools += $tool
    }
}

# Check configuration files
Write-Host "`n📄 Checking configuration files..." -ForegroundColor Yellow

$backendFiles = @(
    "back_end/Dockerfile",
    "back_end/fly.toml",
    "back_end/.env.example",
    "back_end/.gitignore"
)

$frontendFiles = @(
    "front_end/vercel.json",
    "front_end/next.config.js",
    "front_end/.env.example"
)

$missingFiles = @()

foreach ($file in ($backendFiles + $frontendFiles)) {
    if (Test-Path $file) {
        Write-Host "  ✅ $file" -ForegroundColor Green
    } else {
        Write-Host "  ❌ $file - NOT FOUND" -ForegroundColor Red
        $missingFiles += $file
    }
}

# Check documentation
Write-Host "`n📚 Checking documentation..." -ForegroundColor Yellow

$docs = @(
    "DEPLOYMENT.md",
    "QUICKSTART.md",
    "CHECKLIST.md",
    "READY_TO_DEPLOY.md"
)

foreach ($doc in $docs) {
    if (Test-Path $doc) {
        Write-Host "  ✅ $doc" -ForegroundColor Green
    } else {
        Write-Host "  ❌ $doc - NOT FOUND" -ForegroundColor Red
    }
}

# Test backend build
Write-Host "`n🔨 Testing backend build..." -ForegroundColor Yellow
Push-Location back_end
try {
    $buildOutput = go build -o test-build.exe 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✅ Backend builds successfully" -ForegroundColor Green
        Remove-Item test-build.exe -ErrorAction SilentlyContinue
    } else {
        Write-Host "  ❌ Backend build failed" -ForegroundColor Red
        Write-Host "  Error: $buildOutput" -ForegroundColor Red
    }
} catch {
    Write-Host "  ❌ Cannot test backend build" -ForegroundColor Red
}
Pop-Location

# Test frontend build
Write-Host "`n🔨 Testing frontend dependencies..." -ForegroundColor Yellow
Push-Location front_end
try {
    if (Test-Path "node_modules") {
        Write-Host "  ✅ Frontend dependencies installed" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️  Frontend dependencies not installed - run 'npm install'" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  ❌ Cannot check frontend dependencies" -ForegroundColor Red
}
Pop-Location

# Summary
Write-Host "`n" + ("=" * 50)
Write-Host "📊 Summary" -ForegroundColor Cyan
Write-Host ("=" * 50)

if ($missingTools.Count -eq 0 -and $missingFiles.Count -eq 0) {
    Write-Host "✅ All checks passed! You're ready to deploy." -ForegroundColor Green
    Write-Host "`nNext steps:" -ForegroundColor Yellow
    Write-Host "  1. Review CHECKLIST.md"
    Write-Host "  2. Read DEPLOYMENT.md"
    Write-Host "  3. Deploy backend: cd back_end && fly launch"
    Write-Host "  4. Deploy frontend: cd front_end && vercel --prod"
} else {
    Write-Host "⚠️  Some issues found:" -ForegroundColor Yellow
    
    if ($missingTools.Count -gt 0) {
        Write-Host "`nMissing tools:" -ForegroundColor Red
        foreach ($tool in $missingTools) {
            Write-Host "  - $($tools[$tool])"
        }
    }
    
    if ($missingFiles.Count -gt 0) {
        Write-Host "`nMissing files:" -ForegroundColor Red
        foreach ($file in $missingFiles) {
            Write-Host "  - $file"
        }
    }
    
    Write-Host "`n💡 See DEPLOYMENT.md for installation instructions" -ForegroundColor Yellow
}

Write-Host "`n📖 Documentation:"
Write-Host "  - READY_TO_DEPLOY.md - Quick deployment overview"
Write-Host "  - DEPLOYMENT.md      - Detailed deployment guide"
Write-Host "  - QUICKSTART.md      - Local development guide"
Write-Host "  - CHECKLIST.md       - Pre-deployment checklist"
Write-Host ""
