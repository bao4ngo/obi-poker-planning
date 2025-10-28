#!/bin/bash
# Deployment Helper Script
# This script helps verify your setup before deployment

echo "🚀 Poker Planning - Deployment Readiness Check"
echo "=================================================="

# Check if required tools are installed
echo ""
echo "📦 Checking required tools..."

declare -A tools=(
    ["go"]="Go (Backend)"
    ["node"]="Node.js (Frontend)"
    ["npm"]="NPM (Package Manager)"
    ["fly"]="Fly CLI (Backend Deployment)"
    ["vercel"]="Vercel CLI (Frontend Deployment - Optional)"
    ["psql"]="PostgreSQL (Database - Local)"
)

missing_tools=()

for tool in "${!tools[@]}"; do
    if command -v "$tool" &> /dev/null; then
        echo "  ✅ ${tools[$tool]}"
    else
        echo "  ❌ ${tools[$tool]} - NOT FOUND"
        missing_tools+=("$tool")
    fi
done

# Check configuration files
echo ""
echo "📄 Checking configuration files..."

backend_files=(
    "back_end/Dockerfile"
    "back_end/fly.toml"
    "back_end/.env.example"
    "back_end/.gitignore"
)

frontend_files=(
    "front_end/vercel.json"
    "front_end/next.config.js"
    "front_end/.env.example"
)

missing_files=()

for file in "${backend_files[@]}" "${frontend_files[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✅ $file"
    else
        echo "  ❌ $file - NOT FOUND"
        missing_files+=("$file")
    fi
done

# Check documentation
echo ""
echo "📚 Checking documentation..."

docs=(
    "DEPLOYMENT.md"
    "QUICKSTART.md"
    "CHECKLIST.md"
    "READY_TO_DEPLOY.md"
)

for doc in "${docs[@]}"; do
    if [ -f "$doc" ]; then
        echo "  ✅ $doc"
    else
        echo "  ❌ $doc - NOT FOUND"
    fi
done

# Test backend build
echo ""
echo "🔨 Testing backend build..."
cd back_end
if go build -o test-build . 2>/dev/null; then
    echo "  ✅ Backend builds successfully"
    rm -f test-build
else
    echo "  ❌ Backend build failed"
fi
cd ..

# Test frontend dependencies
echo ""
echo "🔨 Testing frontend dependencies..."
if [ -d "front_end/node_modules" ]; then
    echo "  ✅ Frontend dependencies installed"
else
    echo "  ⚠️  Frontend dependencies not installed - run 'npm install'"
fi

# Summary
echo ""
echo "=================================================="
echo "📊 Summary"
echo "=================================================="

if [ ${#missing_tools[@]} -eq 0 ] && [ ${#missing_files[@]} -eq 0 ]; then
    echo "✅ All checks passed! You're ready to deploy."
    echo ""
    echo "Next steps:"
    echo "  1. Review CHECKLIST.md"
    echo "  2. Read DEPLOYMENT.md"
    echo "  3. Deploy backend: cd back_end && fly launch"
    echo "  4. Deploy frontend: cd front_end && vercel --prod"
else
    echo "⚠️  Some issues found:"
    
    if [ ${#missing_tools[@]} -gt 0 ]; then
        echo ""
        echo "Missing tools:"
        for tool in "${missing_tools[@]}"; do
            echo "  - ${tools[$tool]}"
        done
    fi
    
    if [ ${#missing_files[@]} -gt 0 ]; then
        echo ""
        echo "Missing files:"
        for file in "${missing_files[@]}"; do
            echo "  - $file"
        done
    fi
    
    echo ""
    echo "💡 See DEPLOYMENT.md for installation instructions"
fi

echo ""
echo "📖 Documentation:"
echo "  - READY_TO_DEPLOY.md - Quick deployment overview"
echo "  - DEPLOYMENT.md      - Detailed deployment guide"
echo "  - QUICKSTART.md      - Local development guide"
echo "  - CHECKLIST.md       - Pre-deployment checklist"
echo ""
