# Script to run all database migrations
# This script will try to find PostgreSQL installation and run all migration files

param(
    [string]$DbName = "phoenix_alliance",
    [string]$DbUser = "postgres",
    [string]$DbPassword = "",
    [string]$DbHost = "localhost",
    [string]$DbPort = "5432"
)

Write-Host "Running database migrations for: $DbName" -ForegroundColor Cyan

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
    Write-Host "Please add PostgreSQL bin directory to your PATH" -ForegroundColor Yellow
    exit 1
}

# Build connection string
$env:PGPASSWORD = $DbPassword
$connectionString = "-h $DbHost -p $DbPort -U $DbUser -d $DbName"

# Get migration files in order
$migrationFiles = Get-ChildItem -Path "migrations" -Filter "*_up.sql" | Sort-Object Name

if ($migrationFiles.Count -eq 0) {
    Write-Host "No migration files found in migrations/ directory!" -ForegroundColor Red
    exit 1
}

Write-Host "`nFound $($migrationFiles.Count) migration file(s) to run" -ForegroundColor Cyan

# Run each migration
$successCount = 0
$failCount = 0

foreach ($file in $migrationFiles) {
    Write-Host "`nRunning: $($file.Name)..." -ForegroundColor Yellow
    $result = & $psqlPath $connectionString -f $file.FullName 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Successfully applied: $($file.Name)" -ForegroundColor Green
        $successCount++
    } else {
        Write-Host "✗ Failed to apply: $($file.Name)" -ForegroundColor Red
        Write-Host $result -ForegroundColor Red
        $failCount++
    }
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Migration Summary:" -ForegroundColor Cyan
Write-Host "  Successful: $successCount" -ForegroundColor Green
Write-Host "  Failed: $failCount" -ForegroundColor $(if ($failCount -gt 0) { "Red" } else { "Green" })
Write-Host "========================================" -ForegroundColor Cyan

if ($failCount -gt 0) {
    exit 1
}














