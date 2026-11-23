# Installation scripts

The installers expect a built server binary in `binaries/` and must be run from the repository root. They do not download binaries or work when piped directly from GitHub.

## Build and install

Build the release binaries with Docker and Make:

```bash
git clone https://github.com/apimgr/jokes.git
cd jokes
make build
```

On a Linux system using systemd, run:

```bash
sudo ./scripts/install.sh
```

On macOS, run the same command from the checkout; the script uses `sudo` for installation operations. On Windows, open PowerShell as Administrator in the checkout and run:

```powershell
.\scripts\windows.ps1
```

The shell installer supports Linux with systemd and macOS. BSD service installation is not yet supported by these scripts; use a release binary with the platform's service manager.

## Canonical paths

These paths follow AI.md PART 4 for the privileged service installations. The application creates its configuration file as `server.yml` on first run.

| Platform | Binary | Configuration | Data | Logs | Service |
|---|---|---|---|---|---|
| Linux | `/usr/local/bin/jokes` | `/etc/apimgr/jokes/server.yml` | `/var/lib/apimgr/jokes/` | `/var/log/apimgr/jokes/` | `/etc/systemd/system/jokes.service` |
| macOS | `/usr/local/bin/jokes` | `/Library/Application Support/apimgr/jokes/server.yml` | `/Library/Application Support/apimgr/jokes/data/` | `/Library/Logs/apimgr/jokes/` | `/Library/LaunchDaemons/us.apimgr.jokes.plist` |
| Windows | `C:\Program Files\apimgr\jokes\jokes.exe` | `%ProgramData%\apimgr\jokes\server.yml` | `%ProgramData%\apimgr\jokes\data\` | `%ProgramData%\apimgr\jokes\logs\` | Windows Service `jokes` |

Other canonical paths, including cache and backup locations, are listed in [AI.md PART 4](../AI.md#part-4-os-specific-paths). User-mode installations use the user-specific paths in that section.

## Service operations

Linux service controls use `systemctl` with service name `jokes`; logs are available with `journalctl -u jokes`. macOS uses `launchctl` and the LaunchDaemon plist above. Windows uses `Get-Service`, `Start-Service`, `Stop-Service`, and `Restart-Service` with service name `jokes`.

The Windows installer creates the configured data and logs directories separately under `%ProgramData%\apimgr\jokes`.

## Troubleshooting

- **Binary not found:** build from the repository root and confirm the matching platform artifact exists, such as `binaries/jokes-linux-amd64`, `binaries/jokes-darwin-arm64`, or `binaries/jokes-windows-amd64.exe`.
- **Linux service start failure:** inspect `journalctl -u jokes` and verify the configured paths are writable by the `jokes` service user.
- **Port conflict:** edit the runtime-generated `server.yml` and restart the service.

For project issues, use [GitHub Issues](https://github.com/apimgr/jokes/issues). For deployment-specific operations, contact that deployment's operator.
