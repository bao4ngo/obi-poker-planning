# PostgreSQL Database Setup Script

Write-Host "Setting up PostgreSQL database for Poker Planning..." -ForegroundColor Green

# Database connection details
$env:PGPASSWORD = "1234"
$host = "localhost"
$port = "5432"
$user = "postgres"
$database = "poker_planning"

# Check if psql is available
if (!(Get-Command psql -ErrorAction SilentlyContinue)) {
    Write-Host "Error: psql command not found. Please install PostgreSQL client tools." -ForegroundColor Red
    exit 1
}

Write-Host "Creating database..." -ForegroundColor Yellow
psql -h $host -p $port -U $user -c "CREATE DATABASE $database;"

if ($LASTEXITCODE -eq 0) {
    Write-Host "Database created successfully!" -ForegroundColor Green
} else {
    Write-Host "Database might already exist or there was an error. Continuing..." -ForegroundColor Yellow
}

Write-Host "Running schema migration..." -ForegroundColor Yellow
psql -h $host -p $port -U $user -d $database -f "database/schema.sql"

if ($LASTEXITCODE -eq 0) {
    Write-Host "Schema created successfully!" -ForegroundColor Green
    Write-Host "Database setup complete!" -ForegroundColor Green
} else {
    Write-Host "Error running schema migration." -ForegroundColor Red
    exit 1
}

# Clear password from environment
Remove-Item Env:\PGPASSWORD
