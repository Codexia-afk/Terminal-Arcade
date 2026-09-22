# Terminal Arcade - Windows One-Line PowerShell Installer
# Usage: irm https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "Codexia-afk/Terminal-Arcade"
$BinaryName = "arcade.exe"
$InstallDir = "$HOME\.local\bin"

Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  Terminal Arcade - Windows Installer                       " -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan

# 1. Determine architecture
$Arch = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $Arch = "arm64"
}

$DownloadFile = "arcade_windows_$Arch.exe"
$ReleaseUrl = "https://github.com/$Repo/releases/latest/download/$DownloadFile"
$RawFallbackUrl = "https://raw.githubusercontent.com/$Repo/main/dist/$DownloadFile"

# 2. Ensure install directory exists
if (!(Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

$Destination = Join-Path $InstallDir $BinaryName
$Installed = $false

# 3. Check for Go
if (Get-Command go -ErrorAction SilentlyContinue) {
    Write-Host "==> Go compiler detected. Attempting go install..." -ForegroundColor Cyan
    try {
        go install "github.com/$Repo/cmd/arcade@latest"
        $GoPathBin = (go env GOPATH) + "\bin\arcade.exe"
        if (Test-Path $GoPathBin) {
            Copy-Item -Path $GoPathBin -Destination $Destination -Force
            $Installed = $true
            Write-Host "✔ Built and installed via Go toolchain." -ForegroundColor Green
        }
    } catch {
        Write-Host "⚠ Go install failed, falling back to pre-compiled binary download..." -ForegroundColor Yellow
    }
}

# 4. Download pre-compiled binary
if (-not $Installed) {
    Write-Host "==> Downloading $DownloadFile..." -ForegroundColor Cyan
    try {
        Invoke-WebRequest -Uri $ReleaseUrl -OutFile $Destination -UseBasicParsing
        $Installed = $true
    } catch {
        try {
            Invoke-WebRequest -Uri $RawFallbackUrl -OutFile $Destination -UseBasicParsing
            $Installed = $true
        } catch {
            Write-Error "Failed to download binary from GitHub. Please download manually from: https://github.com/$Repo/releases"
            exit 1
        }
    }
}

# 5. Check PATH environment variable
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "==> Adding $InstallDir to User PATH..." -ForegroundColor Cyan
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
}

Write-Host ""
Write-Host "============================================================" -ForegroundColor Green
Write-Host "🎉 Terminal Arcade successfully installed!" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Green
Write-Host ""
Write-Host "👉 Open a new PowerShell terminal and type: arcade" -ForegroundColor Cyan
Write-Host "Enjoy Snake, Pacman, and Breakout!" -ForegroundColor Yellow
