#===============================================================================
# Jokes API - Windows Installation Script
# Run as Administrator
#===============================================================================

#Requires -RunAsAdministrator

$ErrorActionPreference = "Stop"

# Variables
$ProjectName = "jokes"
$BinaryName = "jokes-api.exe"
$ServiceName = "JokesAPI"
$InstallDir = "C:\Program Files\APIMGR\$ProjectName"
$DataDir = "C:\ProgramData\APIMGR\$ProjectName"
$LogDir = "$DataDir\logs"

# Check if running as Administrator
function Test-Administrator {
    $currentUser = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    return $currentUser.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

if (-not (Test-Administrator)) {
    Write-Host "❌ This script must be run as Administrator" -ForegroundColor Red
    exit 1
}

# Create directories
function Create-Directories {
    Write-Host "→ Creating directories..." -ForegroundColor Green

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    New-Item -ItemType Directory -Force -Path $DataDir | Out-Null
    New-Item -ItemType Directory -Force -Path $LogDir | Out-Null
}

# Install binary
function Install-Binary {
    Write-Host "→ Installing binary..." -ForegroundColor Green

    $BinaryPath = $null
    if (Test-Path ".\binaries\$BinaryName") {
        $BinaryPath = ".\binaries\$BinaryName"
    }
    elseif (Test-Path ".\$BinaryName") {
        $BinaryPath = ".\$BinaryName"
    }
    else {
        Write-Host "❌ Binary not found" -ForegroundColor Red
        exit 1
    }

    Copy-Item -Path $BinaryPath -Destination "$InstallDir\$BinaryName" -Force
}

# Install data files
function Install-Data {
    Write-Host "→ Installing data files..." -ForegroundColor Green

    if (Test-Path ".\src\data") {
        Copy-Item -Path ".\src\data\*" -Destination $DataDir -Recurse -Force
    }
}

# Create Windows Service
function Create-Service {
    Write-Host "→ Creating Windows Service..." -ForegroundColor Green

    # Stop and remove existing service if it exists
    $existingService = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    if ($existingService) {
        Stop-Service -Name $ServiceName -Force
        sc.exe delete $ServiceName | Out-Null
        Start-Sleep -Seconds 2
    }

    # Create new service
    $BinaryPath = "$InstallDir\$BinaryName --data `"$DataDir`""

    New-Service -Name $ServiceName `
        -BinaryPathName $BinaryPath `
        -DisplayName "Jokes API" `
        -Description "Jokes API Service - 5,160+ jokes across 16 categories" `
        -StartupType Automatic

    # Start service
    Start-Service -Name $ServiceName

    # Configure service recovery options
    sc.exe failure $ServiceName reset= 86400 actions= restart/60000/restart/60000/restart/60000 | Out-Null
}

# Create firewall rule
function Create-FirewallRule {
    Write-Host "→ Creating firewall rule..." -ForegroundColor Green

    $ruleName = "Jokes API"
    $existingRule = Get-NetFirewallRule -DisplayName $ruleName -ErrorAction SilentlyContinue

    if ($existingRule) {
        Remove-NetFirewallRule -DisplayName $ruleName
    }

    New-NetFirewallRule -DisplayName $ruleName `
        -Direction Inbound `
        -Program "$InstallDir\$BinaryName" `
        -Action Allow `
        -Profile Any `
        -Enabled True | Out-Null
}

# Main installation
function Main {
    Write-Host "🎭 Jokes API Installer for Windows" -ForegroundColor Cyan
    Write-Host "===================================`n" -ForegroundColor Cyan

    Create-Directories
    Install-Binary
    Install-Data
    Create-Service
    Create-FirewallRule

    Write-Host "`n✅ Installation complete!" -ForegroundColor Green
    Write-Host "`nService commands:" -ForegroundColor Yellow
    Write-Host "  Status:  Get-Service -Name $ServiceName"
    Write-Host "  Stop:    Stop-Service -Name $ServiceName"
    Write-Host "  Start:   Start-Service -Name $ServiceName"
    Write-Host "  Restart: Restart-Service -Name $ServiceName"
    Write-Host "  Logs:    Get-Content `"$LogDir\jokes.log`" -Wait"
    Write-Host "`nConfiguration: $DataDir\server.yaml" -ForegroundColor Yellow
}

Main
