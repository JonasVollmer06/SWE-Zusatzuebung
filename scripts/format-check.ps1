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

$unformatted = & $gofmt -l @goFiles

if ($unformatted.Count -gt 0) {
    Write-Host "Diese Go-Dateien sind nicht formatiert:"
    $unformatted | ForEach-Object { Write-Host $_ }
    exit 1
}

Write-Host "Formatierung ist korrekt."
