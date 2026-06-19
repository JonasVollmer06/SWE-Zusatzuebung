$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "postgres-tools.ps1")

Push-Location $composeDir
try {
    Write-Host "Setting up PostgreSQL for swe_zusatzuebung..."
    Write-Host "Using old volume names: pg_data, pg_tablespace, pg_init"

    Initialize-Volumes
    Initialize-TlsIfNeeded
    Start-Postgres
    Ensure-DatabaseExists

    Write-Host "PostgreSQL setup finished."
    Write-Host "Run .\reset-db.ps1 to load the CSV test data."
}
finally {
    Pop-Location
}
