# Script to stop local PostgreSQL service on Windows
# Run this as Administrator

Write-Host "Stopping local PostgreSQL services..." -ForegroundColor Cyan

# Find all PostgreSQL services
$services = Get-Service | Where-Object {$_.Name -like '*postgresql*'}

if ($services.Count -eq 0) {
    Write-Host "No PostgreSQL services found." -ForegroundColor Yellow
    exit 0
}

Write-Host "Found the following PostgreSQL services:" -ForegroundColor Green
$services | Format-Table Name, Status, DisplayName

foreach ($service in $services) {
    if ($service.Status -eq 'Running') {
        Write-Host "Stopping service: $($service.Name)..." -ForegroundColor Yellow
        try {
            Stop-Service -Name $service.Name -Force
            Write-Host "[OK] Stopped: $($service.Name)" -ForegroundColor Green
        } catch {
            Write-Host "[ERROR] Failed to stop: $($service.Name)" -ForegroundColor Red
            Write-Host "  Error: $_" -ForegroundColor Red
            Write-Host "  You may need to run this script as Administrator" -ForegroundColor Yellow
        }
    } else {
        Write-Host "Service $($service.Name) is already stopped" -ForegroundColor Gray
    }
}

Write-Host ""
Write-Host "Done! You can now run: ./run.sh start" -ForegroundColor Green




