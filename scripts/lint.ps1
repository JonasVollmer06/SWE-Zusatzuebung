$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "go-tools.ps1")

$go = Resolve-GoTool "go"
$projectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")

Push-Location $projectRoot
try {
    Write-Host "Fuehre go vet aus..."
    & $go vet ./...

    Write-Host "Linting abgeschlossen."
}
finally {
    Pop-Location
}
