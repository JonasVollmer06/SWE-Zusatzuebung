$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "go-tools.ps1")

$go = Resolve-GoTool "go"
$projectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")

& (Join-Path $PSScriptRoot "format.ps1")
& (Join-Path $PSScriptRoot "lint.ps1")

Push-Location $projectRoot
try {
    Write-Host "Fuehre Tests aus..."
    & $go test ./...

    Write-Host "Check abgeschlossen."
}
finally {
    Pop-Location
}
