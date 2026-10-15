# Jokes

Jokes is a public, read-only Go service for browsing and retrieving a curated joke collection. The server embeds its dataset and web assets in the binary; it does not require an external runtime service, account, or API key.

## Features

- Versioned REST API with random, category, ID, search, and paginated collection endpoints
- GraphQL query endpoint and browser query form
- Server-rendered web catalog, categories, search, theme preferences, and PWA assets
- Companion `jokes-cli` for shell use
- Embedded, build-time joke data

## Run

Download a release binary or build from source using the repository Makefile. The server reads `server.yml` from its configuration directory, creates defaults when needed, and binds to the configured address and port. `CONFIG_DIR` can select the configuration directory; `--config` is also supported.

```sh
jokes --help
jokes --version
jokes --status
jokes --address 127.0.0.1 --port 8080
```

The site defaults to the configured port. Set `PORT` when running the container; the default container port is 80.

## API

All API endpoints are public and read-only.

```sh
curl 'https://jokes.apimgr.us/api/v1/jokes/random'
curl 'https://jokes.apimgr.us/api/v1/jokes/search?q=programming&limit=10'
curl 'https://jokes.apimgr.us/api/v1/jokes/categories'
curl 'https://jokes.apimgr.us/api/v1/jokes/all?page=1&limit=25'
```

- API reference: `/openapi` and `/api/v1/server/swagger`
- Swagger JSON: `/openapi.json`
- Swagger YAML: `/openapi.yaml`
- GraphQL query form: `/graphql`
- GraphQL endpoint: `POST /api/v1/graphql`
- Plain text: add `.txt` to supported joke GET endpoints or send `Accept: text/plain`

## CLI client

```sh
jokes-cli --help
jokes-cli --search programming --count 5
jokes-cli --category dadjokes
jokes-cli --id 42 --format json
jokes-cli --server https://jokes.example.org --count 3
```

The server can be selected with `--server` or `JOKES_SERVER_PRIMARY`. HTTPS is required except for localhost.

## Docker

Build the production image from the repository root:

```sh
docker build -f docker/Dockerfile -t jokes:latest .
```

The human-managed compose file is `docker/docker-compose.yml`. Compose resolves its `./volumes` mounts relative to `docker/`, stores configuration and data under `docker/volumes/`, and binds the service to `172.17.0.1:64580`.

## License

The application is licensed under MIT. Joke dataset source and item-level reuse terms require confirmation; see `TODO.AI.md` before redistributing modified content.
