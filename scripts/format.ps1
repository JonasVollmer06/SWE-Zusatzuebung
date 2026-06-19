$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "go-tools.ps1")

$gofmt = Resolve-GoTool "gofmt"
$goFiles = Get-ChildItem `
    -Path (Join-Path $PSScriptRoot "..\cmd"), (Join-Path $PSScriptRoot "..\internal") `
    -Recurse `
    -Filter "*.go" `
    -File |
    ForEach-Object { $_.FullName }

if ($goFiles.Count -eq 0) {
    Write-Host "Keine Go-Dateien gefunden."
    exit 0
}

Write-Host "Formatiere Go-Dateien mit gofmt..."
& $gofmt -w @goFiles

Write-Host "Formatierung abgeschlossen."
