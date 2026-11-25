## Project description

Jokes is a full-stack Go application serving over 5,160 curated jokes across 16 categories through a versioned REST API, GraphQL endpoint, server-rendered web UI, and companion CLI. Joke records use single one-liners or two-part setup/punchline pairs, and the static corpus is embedded in the server binary. The service also follows the shared AI.md operational contract for configuration, security, scheduled maintenance, backups, updates, service management, and supported optional integrations.

## Project variables

project_name: jokes
project_org: apimgr
internal_name: jokes
internal_org: apimgr
app_name: Jokes API
tagline: Curated jokes for web, API, and CLI.
owner_token: jokes_owner_token_wnggjttn
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
- Server-rendered joke detail pages by ID, with themed 404 responses for unknown IDs
- Paginated full list and per-category list
- Category listing with joke counts
- Keyword search across all jokes
- Full web frontend (server-side Go templates, dark/light/auto theme, PWA, mobile-first)
- Public server information pages for the service, help, health, privacy, and terms
- Registered help and preferences hubs and their configured sub-pages, including token and preference-sync flows
- Shared server/operator security, audit, privacy, backup, update, scheduler, TLS, metrics, and notification capabilities specified by AI.md
- Optional integrations and their configuration-controlled runtime behavior, including external release/feed sources, CAPTCHA/email providers, GeoIP data, and Tor/I2P services where enabled
- CLI client (`jokes-cli`) for shell-pipeline use
- OpenAPI/Swagger documentation and a GraphQL endpoint with a browser query form

**Non-goals:**
- No user accounts, registration, or login of any kind
- No admin web panel (server configured via `server.yml` only)
- No user-submitted or community jokes (curated dataset only, updated via releases)
- No paid tiers or API-key plans. The operator, anonymous-owner, and service tokens required by AI.md are infrastructure credentials, not paid access tiers or user accounts.
- No content rating or flagging system

### Roles & permissions

There are no user accounts or user roles. Public joke pages and retrieval endpoints are anonymous; system, operator, and credential-only owner-token endpoints follow AI.md.

| Actor | Access |
|-------|--------|
| **Anonymous visitor (browser)** | Full read access to public web pages and joke data |
| **Anonymous API client (curl/CLI)** | Full read access to public joke API endpoints |
| **Anonymous token holder** | Uses the credential-only owner token model required by AI.md; it does not create an account or profile |
| **Server operator** | Configures server via `server.yml` only; no web management interface |

### Data model & sensitivity

**Joke record** (embedded at build time, no PII):

| Field | Type | Sensitivity |
|-------|------|-------------|
| `id` | opaque UUID v4 persisted by the server for public lookup and links | Public |
| `source_id` | integer — identifier from the embedded curated dataset; internal only and never returned by a public surface | Internal |
| `categories` | string array — category names | Public |
| `type` | string — `single` or `twopart` | Public |
| `joke` | string — full joke text (single-type only) | Public |
| `setup` | string — setup line (twopart only) | Public |
| `delivery` | string — punchline (twopart only) | Public |

No PII stored or served.

**Public health stats extension:** `/server/healthz` and its API mirror may include `stats.jokes_total`, the aggregate number of jokes loaded by this application. The count contains no user or deployment-specific data.

### Trust boundaries & external services

| Boundary | Trust level | Notes |
|----------|-------------|-------|
| Joke dataset (embedded at build) | Trusted application data | Static corpus compiled into the binary; provenance and item-level redistribution evidence must be maintained separately from the application-code license. |
| Incoming HTTP requests | **Untrusted** | All query parameters validated |
| Configured external integrations | Untrusted | Network responses and provider inputs are validated; optional integrations run only under AI.md configuration and feature rules. |


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

- **Public joke data**: joke retrieval endpoints require no authentication. Server and token-management endpoints follow the authentication requirements in AI.md; no user accounts, registration, password login, or admin web panel are part of the product.
- **CORS**: Public endpoints may use the configured wildcard origin without credentials. Credentialed requests use explicitly allowed origins; wildcard origins are never combined with credentials, as required by AI.md.

### API routes

All project data is public and read-only; these routes require no authentication. API responses use the versioned REST envelope documented in `docs/api.md`.

| Method | Path | Auth | Purpose; request and response |
|---|---|---|---|
| GET | `/api/v1/jokes/random` | Public | Optional category and include/exclude filters; returns one random joke. |
| GET | `/api/v1/jokes/random.txt` | Public | Plain-text alias for one random joke; supports the same filters. |
| GET | `/api/v1/jokes/random/:count` | Public | Count from 1 to 100 plus optional filters; returns that many random jokes. |
| GET | `/api/v1/jokes/:id` | Public | Opaque UUID v4 joke ID; returns one joke or not found. |
| GET | `/api/v1/jokes/all` | Public | Optional search, category filters, page and limit; returns a paginated list. |
| GET | `/api/v1/jokes/all.txt` | Public | Plain-text alias for the collection endpoint. |
| GET | `/api/v1/jokes/search` | Public | Required query plus optional category, page and limit; returns paginated matches. |
| GET | `/api/v1/jokes/search.txt` | Public | Plain-text alias for search results. |
| GET | `/api/v1/jokes/categories` | Public | No request body; returns category names and joke counts. |
| GET | `/api/v1/jokes/categories.txt` | Public | Plain-text alias for category names and counts. |
| GET | `/api/v1/jokes/categories/:category` | Public | Category plus optional page and limit; returns paginated jokes. |
| GET | `/api/v1/jokes/count` | Public | No request body; returns total count and per-category counts. |
| GET | `/api/v1/jokes/count.txt` | Public | Plain-text alias for collection totals. |
| GET | `/browse` | Public | Server-rendered searchable, filterable joke catalog. |
| GET | `/search` | Public | Server-rendered results for the supplied search query and filters. |
| GET | `/random` | Public | Server-rendered random joke, optionally filtered by category. |
| GET | `/jokes/:id` | Public | Server-rendered joke detail page by opaque UUID v4; an unknown ID returns a themed 404 page. |
| GET | `/categories` | Public | Server-rendered category directory and counts. |
| GET | `/categories/:category` | Public | Server-rendered jokes for one category. |

Project data routes are read-only; there are no project-defined create, update, or delete operations.
