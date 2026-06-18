param(
    [string]$Dir = "$env:LOCALAPPDATA\Programs\tpg",
    [string]$Version = "",
    [switch]$NoPath
)

$ErrorActionPreference = "Stop"
$Repo = "Camilo-845/type-cli"
$Binary = "tpg"

if ($Version) {
    $DownloadUrl = "https://github.com/$Repo/releases/download/$Version/$Binary-windows"
} else {
    $DownloadUrl = "https://github.com/$Repo/releases/latest/download/$Binary-windows"
}

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

$userPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User) ?? ""
if ($userPath -notlike "*$Dir*") {
    if (-not $NoPath) {
        [Environment]::SetEnvironmentVariable(
            "Path",
            [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User) + ";$Dir",
            [EnvironmentVariableTarget]::User
        )
        $env:Path += ";$Dir"
        Write-Host "Added $Dir to user PATH." -ForegroundColor Green
    } else {
        Write-Host ""
        Write-Host "Warning: $Dir is not in your PATH." -ForegroundColor Yellow
        Write-Host "  Run this to add it:"
        Write-Host "  [Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$Dir', [EnvironmentVariableTarget]::User)" -ForegroundColor Green
        Write-Host "  Then restart your terminal."
    }
}
