## Project description

Jokes is a full-stack Go web application serving over 5,160 curated jokes across 16 categories — including Chuck Norris facts, dad jokes, programming humor, and anime jokes — through a versioned REST API, GraphQL endpoint, and a server-side rendered web UI. Jokes come in two formats: single one-liners and two-part setup/punchline pairs. All joke data is embedded in the binary at build time. A companion CLI tool lets users pull jokes directly from the terminal. Deployed as a single self-contained static binary.

## Project variables

project_name: jokes
project_org: apimgr
internal_name: jokes
internal_org: apimgr
app_name: Jokes API
repo: https://github.com/apimgr/jokes
license: MIT
binary: jokes
client_binary: jokes-cli

## Business logic

### Product scope & non-goals

**In scope:**
- 5,160+ jokes across 16 categories: anime, chucknorris, dadjokes, programming, and others
- Two joke formats: `single` (one-liner) and `twopart` (setup + delivery/punchline)
- Random joke retrieval (any category or filtered by category)
- Category-filtered retrieval
- Joke lookup by ID
- Paginated full list and per-category list
- Category listing with joke counts
- Keyword search across all jokes
- Full web frontend (server-side Go templates, dark/light/auto theme, PWA, mobile-first)
- Server pages: `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, `/server/terms`
- CLI client (`jokes-cli`) for shell-pipeline use
- OpenAPI/Swagger docs at `/api/{api_version}/server/swagger`
- GraphQL at `/graphql`

**Non-goals:**
- No user accounts, registration, or login of any kind
- No admin web panel (server configured via `server.yml` only)
- No user-submitted or community jokes (curated dataset only, updated via releases)
- No paid tiers, no API keys, no rate-limited access tiers
- No content rating or flagging system

### Roles & permissions

There are no user roles. All endpoints are public and require no authentication.

| Actor | Access |
|-------|--------|
| **Anonymous visitor (browser)** | Full read access to all web pages and API endpoints |
| **Anonymous API client (curl/CLI)** | Full read access to all API endpoints |
| **Server operator** | Configures server via `server.yml` only; no web management interface |

### Data model & sensitivity

**Joke record** (embedded at build time, no PII):

| Field | Type | Sensitivity |
|-------|------|-------------|
| `id` | string — unique identifier | Public |
| `category` | string — category name | Public |
| `type` | string — `single` or `twopart` | Public |
| `joke` | string — full joke text (single-type only) | Public |
| `setup` | string — setup line (twopart only) | Public |
| `delivery` | string — punchline (twopart only) | Public |

No PII stored or served.

### Trust boundaries & external services

| Boundary | Trust level | Notes |
|----------|-------------|-------|
| Joke dataset (embedded at build) | Fully trusted | Static, compiled into binary |
| Incoming HTTP requests | **Untrusted** | All query parameters validated |

No external services called at runtime.

### Threat model & abuse cases

**Primary assets:** service availability.

**Attacker/abuser goals:**
- DoS via high-rate requests
- Bulk scraping of the full dataset

**Defenses:**
- Rate limiting on all endpoints
- Request size limits on all inputs
- Paginated list endpoints limit per-request data volume
- No user accounts eliminates credential stuffing and privilege escalation entirely

### Security decisions & exceptions

- **No authentication on any endpoint**: intentional. Public read-only reference API.
- **All responses include `Access-Control-Allow-Origin: *`**: intentional. Public data API designed for cross-origin browser use.
