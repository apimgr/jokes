# Jokes

[![CI](https://github.com/apimgr/jokes/actions/workflows/ci.yml/badge.svg)](https://github.com/apimgr/jokes/actions/workflows/ci.yml) [![Release](https://img.shields.io/github/v/release/apimgr/jokes)](https://github.com/apimgr/jokes/releases) [![License](https://img.shields.io/github/license/apimgr/jokes)](LICENSE.md) [![Docs](https://readthedocs.org/projects/apimgr-jokes/badge/?version=latest)](https://apimgr-jokes.readthedocs.io/)

## About

Jokes is a Go service for browsing and retrieving a curated joke collection. It includes a public read-only joke API, a server-rendered web catalog, GraphQL, an authenticated MCP endpoint, and the companion `jokes-cli`. The server embeds its dataset and web assets in the binary and does not require an external runtime service or user account.

## Features

- Versioned REST API with random, category, ID, search, and paginated collection endpoints
- GraphQL query endpoint and browser query editor
- Authenticated, read-only Streamable HTTP MCP tools
- Server-rendered catalog, category, search, random, and joke detail pages
- Guest theme and language preferences, owner-token management, and preference sync
- Localized interface in six languages, including right-to-left Arabic
- Reporting, public privacy and security pages, health reporting, and API references
- Installable web interface with offline support for public pages
- Companion `jokes-cli` for shell use and structured output
- Embedded joke corpus with no runtime database dependency for joke data

## Production

### Docker Compose

Review the port and persistent volume settings in [docker/docker-compose.yml](docker/docker-compose.yml), then start the service:

```bash
docker compose -f docker/docker-compose.yml up -d
```

The supplied Compose configuration publishes the service at `172.17.0.1:64580`. Adjust the address and mounts for your host before deployment.

### Binary

Download a release from [GitHub Releases](https://github.com/apimgr/jokes/releases/latest). Releases include server and CLI binaries for Linux, macOS, Windows, and FreeBSD on amd64 and arm64.

| Platform | Server binary | CLI binary |
|---|---|---|
| Linux amd64 / arm64 | `jokes-linux-{amd64,arm64}` | `jokes-cli-linux-{amd64,arm64}` |
| macOS Intel / Apple Silicon | `jokes-darwin-{amd64,arm64}` | `jokes-cli-darwin-{amd64,arm64}` |
| Windows amd64 / arm64 | `jokes-windows-{amd64,arm64}.exe` | `jokes-cli-windows-{amd64,arm64}.exe` |
| FreeBSD amd64 / arm64 | `jokes-freebsd-{amd64,arm64}` | `jokes-cli-freebsd-{amd64,arm64}` |

For Linux amd64:

```bash
curl -q -LSsf -o jokes https://github.com/apimgr/jokes/releases/latest/download/jokes-linux-amd64
chmod +x jokes
sudo install -m 0755 jokes /usr/local/bin/jokes
jokes --help
```

On first start, the server creates its configuration and secret files in the configured directories, then initializes persistent application state. The operator token is printed once at first start; save it securely. For service installation and other platforms, see the [installation guide](docs/installation.md).

### Native services

The server binary supports service installation and lifecycle commands for the detected host service manager. Review the generated service definition and its user, paths, and environment before starting it. The repository installer scripts support systemd on Linux and LaunchDaemon on macOS; see [installation details](docs/installation.md).

## Client

Download the matching `jokes-cli` release binary for your platform. On Linux amd64:

```bash
curl -q -LSsf -o jokes-cli https://github.com/apimgr/jokes/releases/latest/download/jokes-cli-linux-amd64
chmod +x jokes-cli
sudo install -m 0755 jokes-cli /usr/local/bin/jokes-cli
```

Joke retrieval is public and does not require a token. Set the server URL before using the examples:

```bash
: "${JOKES_SERVER_URL:?Set JOKES_SERVER_URL to the deployed API base URL}"
jokes-cli --server "$JOKES_SERVER_URL" --search programming --count 5
jokes-cli --server "$JOKES_SERVER_URL" --category dadjokes
jokes-cli --server "$JOKES_SERVER_URL" --id "$JOKE_ID" --format json
jokes-cli --help
```

See the [CLI reference](docs/cli.md) for options, configuration, shell integration, and operator commands.

## Configuration

The server creates `server.yml`, secret files, and configuration drop-ins in its configuration directory. The `server.yml` file holds primary settings; additional feature settings use `yaml.d/`. See the [configuration guide](docs/configuration.md) for paths, keys, defaults, and secret handling.

## API

Joke retrieval routes are public and read-only. Public joke IDs are opaque UUIDs; numeric IDs from the embedded corpus are never exposed. The service provides REST under `/api/v1/`, GraphQL at `/api/v1/server/graphql`, the OpenAPI document at `/api/v1/server/swagger`, and browser references at `/server/docs/swagger` and `/server/docs/graphql`.

```bash
: "${JOKES_SERVER_URL:?Set JOKES_SERVER_URL to the deployed API base URL}"
curl -q -LSsf "$JOKES_SERVER_URL/api/v1/jokes/random"
curl -q -LSsf "$JOKES_SERVER_URL/api/v1/jokes/search?q=programming&limit=10"
curl -q -LSsf "$JOKES_SERVER_URL/api/v1/jokes/categories"
```

The [API reference](docs/api.md) lists pages, discovery paths, endpoint methods, response formats, and authentication requirements. See [integrations](docs/integrations.md) for GraphQL, OpenAPI, MCP, and CORS.

## Other

Documentation is published on [Read the Docs](https://apimgr-jokes.readthedocs.io/). The [security guide](docs/security.md) explains token handling and protected endpoints. See the [translation guide](docs/translating.md) to contribute a locale. The curated joke corpus is third-party content; verify its source and item-level terms before redistributing it.

## Development

Contributors need Docker and Make. Go builds and tests run in the configured Docker toolchain:

```bash
make dev
make local
make build
make test
```

`make test` runs Go vet, vulnerability scanning, unit tests, and the 60% coverage gate. Run `./tests/run_tests.sh` for the basic Docker integration tier. Add `INCUS=1 INCUS_MODE=systemd` (or `openrc` / `all`) or `E2E=1` to append those checks to `make test`; run `./tests/incus.sh` or `./tests/e2e.sh` directly when preferred. See the [development guide](docs/development.md) for commands and test tiers.

For a local-only image, use `docker buildx build --platform linux/amd64 --load -t jokes:latest -f docker/Dockerfile .`. `make docker` builds and pushes a multi-architecture image to the configured `REGISTRY`.

## Disclaimer

This software is provided as is without warranty. Evaluate it before production use. The application license does not establish reuse rights for the third-party joke corpus.

## License

The application code is licensed under MIT. See [LICENSE.md](LICENSE.md) for dependency attribution. The joke corpus is third-party content; verify its source and item-level terms before redistributing it.
