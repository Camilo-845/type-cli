param(
    [string]$Dir = "$env:LOCALAPPDATA\Programs\tpg",
    [string]$Version = ""
)

$ErrorActionPreference = "Stop"
$Repo = "Camilo-845/typingame"
$Binary = "tpg"

if ($Version) {
    $DownloadUrl = "https://github.com/$Repo/releases/download/$Version/$Binary-windows"
} else {
    $DownloadUrl = "https://github.com/$Repo/releases/latest/download/$Binary-windows"
}

$Arch = [System.Runtime.InteropServices.RuntimeEnvironment]::GetRuntimeDirectory()
if ([Environment]::Is64BitOperatingSystem) {
    $Arch = "amd64"
} else {
    Write-Error "Unsupported architecture: 32-bit. Only 64-bit Windows is supported."
    exit 1
}

$DownloadUrl = "${DownloadUrl}-${Arch}.exe"

New-Item -ItemType Directory -Force -Path $Dir | Out-Null

Write-Host "Downloading $Binary windows/$Arch ..."
Invoke-WebRequest -Uri $DownloadUrl -OutFile "$Dir\$Binary.exe" -UseBasicParsing

Write-Host "Installed $Binary to $Dir\$Binary.exe" -ForegroundColor Green

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$Dir*") {
    Write-Host ""
    Write-Host "Warning: $Dir is not in your PATH." -ForegroundColor Yellow
    Write-Host "  Run this to add it:"
    Write-Host "  [Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$Dir', 'User')" -ForegroundColor Green
    Write-Host "  Then restart your terminal."
}
