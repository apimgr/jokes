# API

The API uses `/api/{api_version}` (the current version is `v1`). Joke retrieval is public and read-only. API responses default to JSON. Supported GET routes return plain text when `.txt` is appended or the request accepts `text/plain`; `Accept: application/json` selects JSON unless the URL ends in `.txt`. API JSON uses two-space indentation and ends with one newline. JSON API responses use `{ "ok": true, "data": ... }` on success and `{ "ok": false, "error": "ERROR_CODE", "message": "...", "details": {} }` on failure. Health responses remain the bare health object. MCP uses its protocol-native JSON-RPC response shape.

## Public Pages

The server renders these pages as HTML. Pages under `/server/` use the shared document shell, locale catalog, and theme. All are public; the token-management page only exposes actions when its owner credential is present.

| Method | Path | Purpose |
|---|---|---|
| GET | `/` | Home page and featured joke |
| GET | `/browse` | Paginated joke collection |
| GET | `/search` | Search jokes |
| GET | `/random` | Random joke page |
| GET | `/jokes/{id}` | Joke detail by public UUID |
| GET | `/categories` | Category list |
| GET | `/categories/{category}` | Category collection |
| GET | `/server` | Permanent redirect to `/server/about` |
| GET | `/server/about` | About and attribution |
| GET | `/server/help` | Help index |
| GET | `/server/help/mcp` | MCP connection help |
| GET | `/server/healthz` | Human-readable service health |
| GET | `/server/privacy` | Privacy policy |
| GET | `/server/tracking.js` | Consent-gated same-origin analytics initialization; no cache and `404` when disabled, unconsented, or served through Tor/I2P |
| GET | `/server/terms` | Terms of service |
| GET | `/server/dpo` | Data protection contact |
| GET | `/server/contact` | Contact and report form |
| GET | `/server/security` | Security overview and reporting instructions |
| GET | `/server/security/policy` | Security policy |
| GET | `/server/security/thanks` | Security acknowledgments |
| GET | `/server/security/report/{tracking_id}` | Token-gated report status |
| GET | `/server/preferences` | Guest preferences |
| GET | `/server/preferences/tokens` | Owner-token management |
| GET | `/server/preferences/sync` | Preference export/import |
| GET | `/server/docs/swagger` | Interactive OpenAPI reference |
| GET | `/server/docs/graphql` | GraphiQL editor |

Preference, consent, CCPA, contact, and token actions use HTML forms and server-side validation. Preferences submit with POST `/server/preferences`; consent and CCPA requests use POST `/server/consent` and `/server/ccpa`. Owner-token actions use POST `/server/preferences/token/{create|restore|reveal|rotate|revoke|forget}`. Linked-token actions use POST `/server/preferences/tokens` and `/server/preferences/tokens/{id}/{reveal|rotate|update|revoke}`. The sync import form accepts GET `/server/preferences/sync/import` for its confirmation page and POST `/server/preferences/sync` to apply the import.

## Discovery and Static Routes

| Method | Path | Purpose |
|---|---|---|
| GET, HEAD | `/robots.txt` | Crawler rules |
| GET, HEAD | `/sitemap.xml` | Sitemap |
| GET, HEAD | `/feed.xml` | Atom feed; the current corpus has no item creation timestamps, so the feed is valid and empty |
| GET, HEAD | `/rss.xml` | RSS 2.0 feed; the current corpus has no item creation timestamps, so the feed is valid and empty |
| GET, HEAD | `/favicon.ico` | Configured browser favicon or embedded default |
| GET, HEAD | `/branding/assets/{filename}` | Validated local branding image |
| GET, HEAD | `/.well-known/security.txt` | Security contact and policy metadata, including Encryption while a key is published |
| GET, HEAD | `/.well-known/pgp-key.asc` | Active security-report public key; 404 when unavailable or publication is disabled |
| GET, HEAD | `/.well-known/llms.txt`, `/llms.txt` | Machine-readable project guidance |
| GET | `/manifest.json` | Web app manifest |
| GET | `/sw.js` | Service worker |
| GET | `/locales/{lang}` | Embedded locale catalog |
| GET | `/static/*` | Static assets |
| GET | `/api/autodiscover` | API and CLI release discovery |
| GET | `/cli/binaries/{filename}` | Published CLI binary download |

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/api/{api_version}/jokes/random` | One random joke; optional `category` |
| GET | `/api/{api_version}/jokes/random/{count}` | Multiple random jokes |
| GET | `/api/{api_version}/jokes/{id}` | Joke by opaque UUID v4; the embedded dataset's numeric source ID is never exposed |
| GET | `/api/{api_version}/jokes/all` | Collection; pagination/search filters |
| GET | `/api/{api_version}/jokes/search?q={query}` | Search the collection |
| GET | `/api/{api_version}/jokes/categories` | Category names and counts |
| GET | `/api/{api_version}/jokes/categories/{category}` | Jokes in a category |
| GET | `/api/{api_version}/jokes/count` | Collection totals |
| GET | `/jokes/{id}` | Server-rendered joke detail page by opaque UUID v4; unknown IDs return a themed 404 |
| GET | `/api/{api_version}/` | API name, version, collection counts, and route links |
| GET | `/api/{api_version}/server/healthz` | JSON service health |
| GET | `/healthz` | Optional direct alias to the server health handler; enabled with `server.healthz.root.enabled: true` (default: false) |
| GET | `/api/{api_version}/server/swagger` | OpenAPI document |
| POST | `/api/{api_version}/server/graphql` | GraphQL query |
| GET, POST | `/api/{api_version}/server/preferences` | Read or save guest preferences |
| GET | `/api/{api_version}/server/preferences/sync` | Export guest preferences as JSON |
| GET | `/api/{api_version}/server/preferences/sync/import` | Apply guest preferences from query/code |
| POST | `/api/{api_version}/server/preferences/sync` | Apply guest preferences |
| GET | `/api/{api_version}/server/{about,help,contact,help/mcp,privacy,terms,dpo,security,security/policy,security/thanks}` | JSON mirrors of public service pages |
| POST | `/api/{api_version}/server/contact` | Submit a contact or report form through its validation and review flow |
| GET | `/api/{api_version}/server/preferences/tokens` | List tokens for the authenticated primary owner token |
| POST | `/api/{api_version}/server/preferences/tokens/primary` | Issue the caller's primary owner token |
| DELETE | `/api/{api_version}/server/preferences/tokens/primary` | Revoke the caller's owner and linked tokens |
| POST | `/api/{api_version}/server/preferences/tokens` | Create a linked token (primary token required) |
| POST | `/api/{api_version}/server/preferences/tokens/{id}/reveal` | Reveal a linked token (primary token required) |
| POST | `/api/{api_version}/server/preferences/tokens/{id}/rotate` | Rotate a linked token (primary token required) |
| PATCH | `/api/{api_version}/server/preferences/tokens/{id}` | Rename or update a linked token (primary token required) |
| DELETE | `/api/{api_version}/server/preferences/tokens/{id}` | Revoke a linked token (primary token required) |
| GET | `/api/{api_version}/server/auth/check` | Check the presented token and return its owner/token kind |
| POST | `/api/{api_version}/server/mcp` | Stateless Streamable HTTP MCP endpoint; linked-token Bearer authentication required |
| GET | `/api/{api_version}/server/metrics[/{service}]` | Internal metrics services; each service requires its own Bearer token |
| GET | `/api/autodiscover` | Public API version, public origins, exposed overlay addresses, and verified CLI release metadata |

Autodiscovery returns a standard `{ "ok": true, "data": ... }` envelope. `data` contains the project `name`, running `version`, `api.current_version`, public origin URLs in `urls`, and verified `cli_versions` plus `cli_min_version`. On onion and I2P requests, `urls` is empty to avoid disclosing clearnet origins. The `onion` and `i2p` addresses appear only when the operator explicitly exposes a healthy overlay. It does not publish server settings or option lists.

The versionless service endpoints `/api/healthz`, `/api/healthz.txt`, `/api/swagger`, `/api/graphql`, and `/api/metrics[/{service}]` are aliases. `/healthz` is a separate optional direct alias controlled by `server.healthz.root.enabled`; it defaults to disabled. MCP is also available at `/api/mcp` and `/mcp`, with POST as its only method. The OpenAPI document contains the operations represented in its schema; the service, owner-token, and MCP route details above supplement it. Internal metrics routes are not advertised by discovery.

Every HTTP response includes the public `X-Build` header with the running build stamp (`version-shortcommit`). Browser clients can compare it with the stamp embedded in the rendered page and reload after an update; if a form has unsaved changes, the page offers a reload action instead. For allowed cross-origin requests, `X-Build`, `ETag`, and `Retry-After` are exposed to browser JavaScript. Rendered HTML uses `Cache-Control: no-store` and a build-stamped `ETag`.

The `/api/{api_version}/server/about` and `/api/{api_version}/server/terms` mirrors return `title`, source `content`, and `format`. `/api/{api_version}/server/help` returns those fields plus the ordered registered `pages` list; `/api/{api_version}/server/help/mcp` returns the MCP connection details. `/api/{api_version}/server/security` returns ordered `contacts`, an RFC 3339 `expires` value, the PGP availability object, absolute policy and acknowledgments URLs, and `report_url` only while the report form is enabled. `/api/{api_version}/server/security/thanks` returns `entries` with only `credit`, `date`, and `severity`.

The help mirror includes real API examples and documentation links. It includes Tor and I2P addresses and setup instructions only when the corresponding overlay is enabled, published, and healthy. The privacy mirror's `tracking` object and `third_party.services` reflect the configured analytics provider. Analytics script initialization is served from `/server/tracking.js` only after consent; the route is uncached and returns `404` when disabled, unconsented, or accessed through an overlay host. CSP permits only the configured provider hosts, and tracking never adds inline scripts.

The contact-form descriptor includes a `captcha` object. `simple` returns a single-use question challenge; configured `recaptcha` and `hcaptcha` descriptors include the public `site_key` and a built-in `fallback` challenge for API clients. Submit a provider token as `captcha_response` (or the provider's native response field), or submit the fallback `captcha_id` and `captcha_answer`. Provider verification fails closed. Requests served on Tor or I2P always use the simple challenge and never call a clearnet CAPTCHA provider.

The `/api/{api_version}/server/dpo` mirror returns the configured DPO name, the DPO email or general-contact fallback, and an absolute privacy-policy URL resolved for the current request. It returns an empty email when neither address is configured and never publishes the administrator email. The emailed security-report status link is a token-gated HTML page at `/server/security/report/{tracking_id}`; it has no JSON mirror.

Health responses include the public-safe project, status, version/build, runtime, feature, component-check, and aggregate-stat fields. When GeoIP is enabled, `features.geoip` is `true` and `checks.geoip_country`, `checks.geoip_asn`, and `checks.geoip_city` report `ok`, `degraded`, or `error`. During maintenance, the response includes `maintenance.reason`, operator guidance in `maintenance.message`, the start time, and self-healing attempt state; the same information is rendered at `/server/healthz`. Database guidance covers connectivity and credentials; disk and write guidance covers capacity, permissions, and configured cleanup. Health is `200` for `healthy`, `degraded`, or `restart_required`; critical failure, maintenance, and shutdown states return `503` while retaining the response body. Health `.txt` output uses machine-readable dot notation.

For example, a successful joke request returns `{ "ok": true, "data": { "id": "...", "joke": "..." } }`. Unknown resources return HTTP `404` with `{ "ok": false, "error": "NOT_FOUND", "message": "resource not found", "details": {} }`.

Examples:

```bash
: "${JOKES_SERVER_URL:?Set JOKES_SERVER_URL to the deployed API base URL}"
curl -q -LSsf "$JOKES_SERVER_URL/api/{api_version}/jokes/random"
curl -q -LSsf "$JOKES_SERVER_URL/api/{api_version}/jokes/search?q=programming&limit=10"
curl -q -LSsf "$JOKES_SERVER_URL/api/{api_version}/jokes/categories"
curl -q -LSsf "$JOKES_SERVER_URL/api/{api_version}/server/swagger"
```

Browser references are available at `/server/docs/swagger` and `/server/docs/graphql`. Metrics are internal operational endpoints: `/server/metrics[/prometheus]` returns Prometheus exposition, `/server/metrics/grafana` returns a dashboard, and `/server/metrics/loki` returns bounded sanitized log streams. Each service has a separate Bearer token generated in `secrets.yml`; `/metrics[/{service}]` aliases are enabled by default and can be disabled with `server.metrics.root.enabled: false`. Restrict these paths at the firewall or reverse proxy.

Public joke identifiers are UUID v4 values assigned and persisted by the server. The numeric IDs in the embedded dataset are internal source identifiers and are not accepted by public joke lookup routes. Public joke routes are read-only and do not require authentication. Owner-token management and MCP use Bearer tokens; the API reference does not include operator-only maintenance behavior.

The GraphQL handler validates and executes read-only queries against the published schema, including standard schema introspection. Its browser editor is available at `/server/docs/graphql`; use the editor's default `<code>{ categories }</code>` query as a starting point. Successful GraphQL responses, including execution errors, retain their standard `data` and `errors` fields with HTTP `200`; failures before execution use the canonical API error envelope and their corresponding HTTP status.
