$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "postgres-tools.ps1")

Push-Location $composeDir
try {
    Write-Host "Starting PostgreSQL for swe_zusatzuebung..."
    Start-Postgres
    Show-RowCount
}
finally {
    Pop-Location
}
