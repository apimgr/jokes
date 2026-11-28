# Development

## Prerequisites

Docker and Make are used for Go builds, tests, and analysis. The Makefile uses `casjaysdev/go:latest` and builds with CGO disabled.

## Build and checks

```bash
make dev
make local
make build
make test
```

`make dev` builds into a project-scoped temporary directory. `make local` builds the current platform; `make build` creates server and CLI binaries for Linux, macOS, Windows, and FreeBSD on amd64 and arm64. `make test` runs `go vet`, `staticcheck`, `govulncheck`, ShellCheck, `shfmt`, Go tests with coverage data conversion and the configured coverage floor, then the basic Docker binary and server-page smoke checks. Add `INCUS=1 INCUS_MODE=systemd` (or `openrc` / `all`) for service integration tests, or `E2E=1` for browser tests. Go checks run in the configured toolchain container.

## Docker integration smoke checks

`make test` calls `./tests/run_tests.sh` after the Go checks. The script runs only the basic Docker binary-validation tier: it builds the server and CLI binaries, starts the server in a disposable Alpine container, and checks command output, localized server help with `--lang`, fallback to English for an unsupported language, and exit status plus stderr for invalid CLI input. It also checks public API routes and server-rendered pages. The route checks include the joke detail page, themed 404 for an unknown joke, and plain-text joke-by-ID API endpoint. Runtime files and default build outputs are stored in a project-scoped directory in the OS temporary area and cleaned up on exit. Run `./tests/docker.sh` to run only the binary smoke tier. Both scripts require Docker.

## On-demand integration and browser checks

Run `./tests/incus.sh [systemd|openrc|all]` or `make test INCUS=1 INCUS_MODE=systemd` to exercise service installation, restart/start/stop, boot persistence, API and page routes, owner-token/MCP access, and Tor/I2P reporting in Debian systemd and Alpine OpenRC Incus containers. The default mode is `all`; Incus must be installed and available. The script builds test binaries and stores temporary runtime data under the project-scoped OS temporary directory, then removes them on exit.

Run `./tests/e2e.sh` or `make test E2E=1` for the three E2E tiers. It builds with the project Go image into a project-scoped temporary directory and runs tagged tests in Docker with a headless Chromium sidecar. Tier 1 checks SSR page content and document metadata, joke detail/API consistency, themed errors, and API/frontend content negotiation. Tier 2 disables page JavaScript while crawling registered pages and submits the native preferences theme form. Tier 3 crawls the pages in Chromium and checks navigation, service-worker registration, JavaScript exceptions, console errors, failed requests, and failed assets. These checks are optional and run only when requested. Set `E2E_SKIP_BUILD=1` to reuse an existing server binary in `BINDIR`; use it only after rebuilding that binary for the current source. Set `E2E_KEEP_ARTIFACTS=1` to retain the run’s private server logs and fixtures under the project-scoped temporary directory when diagnosing a failure.

`make docker` builds and loads a local image for the Docker host platform; it does not push to a registry. For an explicit local platform, run `docker buildx build --platform linux/amd64 --load -t jokes:latest -f docker/Dockerfile .`.

The code is organized under `src/`; web templates and static assets are embedded in the server binary. Avoid committing generated binaries or runtime configuration.

Server-rendered pages are declared in `src/server/page_routes.go`. The app mounts that registry, and the server conformance test renders every registered page plus the shared error pages. Page-specific state-changing endpoints remain in the app route setup and are not treated as pages.

Runtime dataset health checks are registered in `src/health`; health responses evaluate their cached status callbacks in registration order and expose only stable keys and `ok`, `degraded`, or `error` values. Checks for disabled datasets are removed from the registry.

## Contributions

Keep changes scoped, update user documentation when behavior changes, and include tests for behavior changes. Review the repository’s AI.md and IDEA.md before changing implementation or product behavior.
