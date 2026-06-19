$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "postgres-tools.ps1")

Push-Location $composeDir
try {
    Write-Host "Resetting PostgreSQL data for swe_zusatzuebung..."
    Start-Postgres
    Ensure-DatabaseExists
    Reset-DatabaseFromCsv
    Show-RowCount
}
finally {
    Pop-Location
}
