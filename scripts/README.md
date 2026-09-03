# Jokes API - Installation Scripts

This directory contains installation scripts for different operating systems.

## Quick Install

### Linux
```bash
curl -fsSL https://raw.githubusercontent.com/apimgr/jokes/main/scripts/install.sh | sudo bash
```

### macOS
```bash
curl -fsSL https://raw.githubusercontent.com/apimgr/jokes/main/scripts/install.sh | bash
```

### Windows (PowerShell as Administrator)
```powershell
irm https://raw.githubusercontent.com/apimgr/jokes/main/scripts/windows.ps1 | iex
```

## Manual Installation

### Prerequisites
- Go 1.21 or higher (for building from source)
- Git (optional, for cloning repository)

### Install from Binary

1. Download the appropriate binary for your platform from the [releases page](https://github.com/apimgr/jokes/releases)
2. Extract the archive
3. Run the installation script for your OS

### Build from Source

```bash
git clone https://github.com/apimgr/jokes.git
cd jokes
make build
sudo ./scripts/install.sh
```

## Installation Locations

### Linux
- Binary: `/usr/local/bin/jokes-api`
- Data: `/usr/share/apimgr/jokes/`
- Config: `/etc/apimgr/jokes/server.yaml`
- Logs: `/var/log/apimgr/jokes/`
- Service: `/etc/systemd/system/jokes-api.service`

### macOS
- Binary: `/usr/local/bin/jokes-api`
- Data: `/usr/local/share/apimgr/jokes/`
- Config: `/usr/local/etc/apimgr/jokes/server.yaml`
- Logs: `/usr/local/var/log/apimgr/jokes/`
- Service: `/Library/LaunchDaemons/com.apimgr.jokes.plist`

### Windows
- Binary: `C:\Program Files\APIMGR\jokes\jokes-api.exe`
- Data: `C:\ProgramData\APIMGR\jokes\`
- Config: `C:\ProgramData\APIMGR\jokes\server.yaml`
- Logs: `C:\ProgramData\APIMGR\jokes\logs\`
- Service: Windows Service named "JokesAPI"

## Post-Installation

After installation, the service will start automatically. You can access the API at:
- Web Interface: http://localhost:64xxx/ (random port in 64xxx range)
- API: http://localhost:64xxx/api/v1/
- Health Check: http://localhost:64xxx/healthz

To find the actual port, check the configuration file or service logs.

## Uninstall

### Linux
```bash
sudo systemctl stop jokes-api
sudo systemctl disable jokes-api
sudo rm /usr/local/bin/jokes-api
sudo rm /etc/systemd/system/jokes-api.service
sudo rm -rf /usr/share/apimgr/jokes
sudo rm -rf /etc/apimgr/jokes
```

### macOS
```bash
sudo launchctl unload /Library/LaunchDaemons/com.apimgr.jokes.plist
sudo rm /Library/LaunchDaemons/com.apimgr.jokes.plist
sudo rm /usr/local/bin/jokes-api
sudo rm -rf /usr/local/share/apimgr/jokes
sudo rm -rf /usr/local/etc/apimgr/jokes
```

### Windows (PowerShell as Administrator)
```powershell
Stop-Service -Name "JokesAPI"
sc.exe delete "JokesAPI"
Remove-Item "C:\Program Files\APIMGR\jokes" -Recurse -Force
Remove-Item "C:\ProgramData\APIMGR\jokes" -Recurse -Force
```

## Troubleshooting

### Service won't start
Check the logs:
- Linux: `journalctl -u jokes-api -f`
- macOS: `tail -f /usr/local/var/log/apimgr/jokes/jokes.log`
- Windows: Check Event Viewer or `C:\ProgramData\APIMGR\jokes\logs\`

### Port already in use
Edit the configuration file and change the port, then restart the service.

### Permission denied
Make sure you're running the installation script with appropriate permissions (sudo on Linux/macOS, Administrator on Windows).

## Support

For issues and questions:
- GitHub Issues: https://github.com/apimgr/jokes/issues
- Website: https://jokes.apimgr.us
