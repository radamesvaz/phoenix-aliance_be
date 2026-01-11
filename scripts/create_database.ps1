# Script to create PostgreSQL database for Phoenix Alliance
# This script will try to find PostgreSQL installation and create the database

param(
    [string]$DbName = "phoenix_alliance",
    [string]$DbUser = "postgres",
    [string]$DbPassword = "",
    [string]$DbHost = "localhost",
    [string]$DbPort = "5432"
)

Write-Host "Creating PostgreSQL database: $DbName" -ForegroundColor Cyan

# Common PostgreSQL installation paths
$postgresPaths = @(
    "C:\Program Files\PostgreSQL\16\bin\psql.exe",
    "C:\Program Files\PostgreSQL\15\bin\psql.exe",
    "C:\Program Files\PostgreSQL\14\bin\psql.exe",
    "C:\Program Files\PostgreSQL\13\bin\psql.exe",
    "C:\Program Files\PostgreSQL\12\bin\psql.exe",
    "$env:ProgramFiles\PostgreSQL\*\bin\psql.exe"
)

$psqlPath = $null

# Try to find psql in PATH first
$psqlInPath = Get-Command psql -ErrorAction SilentlyContinue
if ($psqlInPath) {
    $psqlPath = $psqlInPath.Source
    Write-Host "Found psql in PATH: $psqlPath" -ForegroundColor Green
} else {
    # Try common installation paths
    foreach ($path in $postgresPaths) {
        $resolvedPath = Resolve-Path $path -ErrorAction SilentlyContinue
        if ($resolvedPath -and (Test-Path $resolvedPath)) {
            $psqlPath = $resolvedPath.Path
            Write-Host "Found psql at: $psqlPath" -ForegroundColor Green
            break
        }
    }
}

if (-not $psqlPath) {
    Write-Host "`nERROR: PostgreSQL psql command not found!" -ForegroundColor Red
    Write-Host "`nPlease do one of the following:" -ForegroundColor Yellow
    Write-Host "1. Add PostgreSQL bin directory to your PATH" -ForegroundColor Yellow
    Write-Host "   (Usually: C:\Program Files\PostgreSQL\<version>\bin)" -ForegroundColor Yellow
    Write-Host "2. Install PostgreSQL from: https://www.postgresql.org/download/windows/" -ForegroundColor Yellow
    Write-Host "3. Use pgAdmin to create the database manually" -ForegroundColor Yellow
    Write-Host "`nOr provide the full path to psql.exe:" -ForegroundColor Yellow
    Write-Host "   .\scripts\create_database.ps1 -PsqlPath 'C:\Program Files\PostgreSQL\16\bin\psql.exe'" -ForegroundColor Yellow
    exit 1
}

# Build connection string
$env:PGPASSWORD = $DbPassword
$connectionString = "-h $DbHost -p $DbPort -U $DbUser"

# Check if database already exists
Write-Host "`nChecking if database exists..." -ForegroundColor Cyan
$checkDb = & $psqlPath $connectionString -tAc "SELECT 1 FROM pg_database WHERE datname='$DbName'" 2>&1

if ($checkDb -match "1") {
    Write-Host "Database '$DbName' already exists!" -ForegroundColor Yellow
    $response = Read-Host "Do you want to drop and recreate it? (y/N)"
    if ($response -eq "y" -or $response -eq "Y") {
        Write-Host "Dropping existing database..." -ForegroundColor Yellow
        & $psqlPath $connectionString -c "DROP DATABASE $DbName;" 2>&1 | Out-Null
    } else {
        Write-Host "Keeping existing database. Exiting." -ForegroundColor Green
        exit 0
    }
}

# Create database
Write-Host "`nCreating database '$DbName'..." -ForegroundColor Cyan
$result = & $psqlPath $connectionString -c "CREATE DATABASE $DbName;" 2>&1

if ($LASTEXITCODE -eq 0) {
    Write-Host "Database '$DbName' created successfully!" -ForegroundColor Green
    Write-Host "`nNext steps:" -ForegroundColor Cyan
    Write-Host "1. Run migrations using the commands in README.md" -ForegroundColor Yellow
    Write-Host "2. Or use: .\scripts\run_migrations.ps1" -ForegroundColor Yellow
} else {
    Write-Host "`nERROR: Failed to create database!" -ForegroundColor Red
    Write-Host $result -ForegroundColor Red
    Write-Host "`nCommon issues:" -ForegroundColor Yellow
    Write-Host "- PostgreSQL service might not be running" -ForegroundColor Yellow
    Write-Host "- Wrong password for user '$DbUser'" -ForegroundColor Yellow
    Write-Host "- User '$DbUser' doesn't have permission to create databases" -ForegroundColor Yellow
    exit 1
}














