function Resolve-GoTool {
    param(
        [Parameter(Mandatory = $true)]
        [string] $Name
    )

    $command = Get-Command $Name -ErrorAction SilentlyContinue
    if ($null -ne $command) {
        return $command.Source
    }

    $fallback = Join-Path $env:ProgramFiles "Go\bin\$Name.exe"
    if (Test-Path $fallback) {
        return $fallback
    }

    throw "$Name wurde nicht gefunden. Ist Go installiert und im PATH?"
}
