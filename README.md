# Jokes

Jokes is a public, read-only Go service for browsing and retrieving a curated joke collection. The server embeds its dataset and web assets in the binary; it does not require an external runtime service or user account.

[![License](https://img.shields.io/github/license/apimgr/jokes)](LICENSE.md) [![Documentation](https://readthedocs.org/projects/apimgr-jokes/badge/?version=latest)](https://apimgr-jokes.readthedocs.io/)

## Install

Download a release from [GitHub Releases](https://github.com/apimgr/jokes/releases/latest). Releases include server and CLI binaries for Linux, macOS, Windows, and FreeBSD on amd64 and arm64.

| Platform | Server binary | CLI binary |
|---|---|---|
| Linux amd64 / arm64 | `jokes-linux-{amd64,arm64}` | `jokes-cli-linux-{amd64,arm64}` |
| macOS Intel / Apple Silicon | `jokes-darwin-{amd64,arm64}` | `jokes-cli-darwin-{amd64,arm64}` |
| Windows amd64 / arm64 | `jokes-windows-{amd64,arm64}.exe` | `jokes-cli-windows-{amd64,arm64}.exe` |
| FreeBSD amd64 / arm64 | `jokes-freebsd-{amd64,arm64}` | `jokes-cli-freebsd-{amd64,arm64}` |

On Linux amd64, install the server with:

```bash
curl -q -LSsf -o jokes https://github.com/apimgr/jokes/releases/latest/download/jokes-linux-amd64
chmod +x jokes
sudo install -m 0755 jokes /usr/local/bin/jokes
```

For macOS, choose the matching `darwin` binary and remove quarantine after download with `xattr -d com.apple.quarantine ./jokes 2>/dev/null || true`. Windows users should download the matching `.exe` and add its directory to `%PATH%`.

## Docker

The production compose file is [docker/docker-compose.yml](docker/docker-compose.yml). Review its port and volume settings, then run `docker compose -f docker/docker-compose.yml up -d`.

## Server

Run `jokes --help` for server options. On first start, the server creates its configuration and secret files in its configuration directory, then initializes its persistent SQLite state. The operator token is printed once at first start; save it securely. The web interface includes browse, search, category, random joke, API reference, owner-token management, guest preferences, reporting, and public server information pages.

## Client CLI

The `jokes-cli` companion supports shell use and structured output:

```bash
jokes-cli --help
jokes-cli --search programming --count 5
jokes-cli --category dadjokes
jokes-cli --id "$JOKE_ID" --format json
```

## API

Joke retrieval routes are public and read-only. Public joke IDs are opaque UUIDs; the embedded dataset's numeric source IDs are never exposed. The service exposes REST at `/api/v1/`, GraphQL at `/api/v1/server/graphql`, the OpenAPI document at `/api/v1/server/swagger`, and browser references at `/server/docs/swagger` and `/server/docs/graphql`.

```bash
: "${JOKES_SERVER_URL:?Set JOKES_SERVER_URL to the deployed API base URL}"
curl -q -LSsf "$JOKES_SERVER_URL/api/v1/jokes/random"
curl -q -LSsf "$JOKES_SERVER_URL/api/v1/jokes/search?q=programming&limit=10"
curl -q -LSsf "$JOKES_SERVER_URL/api/v1/jokes/categories"
```

See the [API reference](https://apimgr-jokes.readthedocs.io/en/latest/api/) for routes, pagination, and response formats.

## Documentation

User and operator guides are published at [Read the Docs](https://apimgr-jokes.readthedocs.io/), including the [translation guide](https://apimgr-jokes.readthedocs.io/en/latest/translating/).

## Features

- Versioned REST API with random, category, ID, search, and paginated collection endpoints
- GraphQL query endpoint and browser query form
- Authenticated Streamable HTTP MCP with read-only joke tools
- Server-rendered web catalog, categories, search, theme preferences, and PWA assets
- Server-rendered joke detail pages with themed not-found responses
- Persistent owner and linked tokens for MCP and owner-scoped API operations
- Companion `jokes-cli` for shell use
- Embedded, build-time joke data

## Development

Contributors need Docker and Make. All Go commands run in the configured Go toolchain container:

```bash
make dev
make local
make build
make test
```

`make test` runs Go vet, vulnerability scanning, unit tests, and the 60% coverage gate. Run the lint checks in the Go toolchain container with `go vet ./...`, `staticcheck ./...`, and `gofmt -l ./src`.

Run `./tests/run_tests.sh` for the basic Docker tier: it runs `make test`, builds the server and CLI binaries, then smoke-tests their command output, public API routes, and server-rendered pages in a disposable Alpine container. The smoke checks include joke detail pages, the themed 404 for an unknown joke, and the plain-text joke-by-ID endpoint. To run only the binary smoke tier, use `./tests/docker.sh`. Both scripts require Docker.

For on-demand system integration checks, `./tests/incus.sh [systemd|openrc|all]` exercises service installation and lifecycle, API and page routes, owner-token/MCP access, and Tor/I2P health reporting in Debian systemd and Alpine OpenRC containers. It requires Incus; `all` is the default. `./tests/e2e.sh` runs three browser-oriented tiers in Docker: SSR and content-negotiation checks, page/form checks with JavaScript disabled, and Chromium navigation/PWA registration checks. Neither on-demand script is part of `make test` or the commit gate.

### Docker build

`make docker` builds and pushes the multi-architecture image to the configured registry. Set `REGISTRY` before running it. For a local-only image, use `docker buildx build --platform linux/amd64 --load -t jokes:latest -f docker/Dockerfile .`.

## Disclaimer

This software is provided as is without warranty. Evaluate it before production use. The application license does not establish reuse rights for the third-party joke corpus.

## License

The application code is licensed under MIT. See [LICENSE.md](LICENSE.md) for dependency attribution. The joke corpus is third-party content; verify its source and item-level terms before redistributing it.
