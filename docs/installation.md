# Installation

## Release binaries

Download the server or `jokes-cli` binary from [GitHub Releases](https://github.com/apimgr/jokes/releases/latest). Each release targets Linux, macOS (Darwin), Windows, and FreeBSD on amd64 and arm64. Windows artifacts use `.exe`; the other platforms do not.

For Linux amd64:

```bash
curl -q -LSsf -o jokes https://github.com/apimgr/jokes/releases/latest/download/jokes-linux-amd64
chmod +x jokes
sudo install -m 0755 jokes /usr/local/bin/jokes
jokes --help
jokes --shell completions bash > ~/.local/share/bash-completion/completions/jokes
```

For macOS, select `jokes-darwin-amd64` (Intel) or `jokes-darwin-arm64` (Apple Silicon), make it executable, and clear quarantine with `xattr -d com.apple.quarantine ./jokes 2>/dev/null || true`. On Windows, download the matching `.exe` and add its directory to `%PATH%`.

## Repository installer scripts

The repository also contains `scripts/install.sh`, which delegates installation to `scripts/linux.sh` or `scripts/macos.sh`. Run it from a complete repository checkout at the repository root so it can find `binaries/` and `src/data/`:

```bash
sudo bash scripts/install.sh
```

Linux installation requires root and installs a systemd service. The macOS script installs a LaunchDaemon and invokes `sudo` for privileged operations. Both scripts support amd64 and arm64 builds. FreeBSD and OpenBSD are detected by the wrapper but are not installed by these scripts; use a release binary and configure the native service manager separately. The installer has no command-line options and expects an already-built binary in `binaries/` or at the repository root.

## Docker Compose

Review [docker/docker-compose.yml](https://github.com/apimgr/jokes/blob/main/docker/docker-compose.yml), then run:

```bash
curl -q -LSsf -o docker-compose.yml https://raw.githubusercontent.com/apimgr/jokes/main/docker/docker-compose.yml
docker compose up -d
```

The compose configuration publishes the service on `172.17.0.1:64580` and defines its persistent mounts. Adjust these for your host before deployment.

## Run from source

See [Development](development.md) for the containerized build workflow. Start the resulting server binary with `jokes --help` to review options.

## Native service management

The server includes service install, disable, uninstall, start, stop, restart, reload, and help commands. `--service=--disable` stops the service and removes its boot activation while retaining the service definition and application data. `--service=--help` also reports the detected manager's current state. On a supported host, run the server binary with the appropriate `--service` command; for example:

```bash
sudo jokes --service=--install
sudo jokes --service=start
sudo jokes --service=stop
sudo jokes --service=--uninstall
```

The service manager is detected from the operating system. The server has installers for systemd, OpenRC, SysVinit, runit, s6, launchd, Windows services, and BSD rc. The repository installer scripts cover only systemd on Linux and LaunchDaemons on macOS; use the server binary's `--service` commands for the other supported managers. Review the generated service definition and its user, paths, and environment before starting it. Containers should use the supplied Docker image or Compose configuration instead of installing a host service.

On Linux, a root or system-service start creates or reuses the dedicated `jokes` system user and group, prepares and assigns the runtime directories, binds any privileged ports, then drops the server and its Tor/I2P provider processes to that account. The service manager starts the binary as root so the binary can complete that setup itself. Uninstall asks for confirmation, stops and removes the service, deletes application data and a server-created account, and leaves the installed binary in place for manual removal.
