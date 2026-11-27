# Integrations

## REST

Use the versioned API under `/api/v1/`. Joke retrieval is public and read-only. JSON is the default; supported GET endpoints also accept `.txt` or `Accept: text/plain` for plain text. See [API](api.md).

## GraphQL

POST queries to `/api/v1/server/graphql`. The browser editor is at `/server/docs/graphql`; queries are validated against the published schema, and standard schema introspection is available. The versionless alias is `/api/graphql`.

## OpenAPI

The OpenAPI JSON document is available at `/api/v1/server/swagger`; the browser reference is `/server/docs/swagger`.

## Tor and I2P

When a Tor binary is installed, the server starts an isolated Tor process and serves the app through its own hidden service backend. The onion hostname appears in health output after Tor is healthy. `server.tor.expose` controls whether autodiscovery publishes the hostname; Tor has no enable or disable switch. Outbound application requests use Tor only when `server.tor.use_network` is enabled.

I2P is enabled with `server.i2p.enabled`, `I2P_ENABLED=true`, or `--i2p`; it defaults on in containers and off elsewhere. The server prefers an installed `i2pd` binary and otherwise probes a local SAM bridge. If neither provider responds, it creates no I2P state or backend listener. Health output reports the provider and eepsite state when I2P is enabled.

## MCP

Jokes implements stateless Streamable HTTP MCP at `/api/v1/server/mcp`; `/api/mcp` and `/mcp` are direct aliases. Requests require a valid `Authorization: Bearer <token>` using an existing primary owner, linked owner, or operator token. A linked token is recommended for each MCP client so its access can be revoked independently. Browser cookies are not accepted. Requests use MCP revision `2026-07-28` and carry matching `MCP-Protocol-Version`, `_meta.io.modelcontextprotocol/protocolVersion`, and `Mcp-Method` values. `tools/call` also carries `Mcp-Name` matching the tool name.

The tool catalog is derived from the API route registry and uses each OpenAPI `operationId` as its tool name. It includes public API operations and hides cookie-bound preferences, token management, token checks, and internal routes. Public joke lookups are read-only; contact submission is a write operation. The server does not expose MCP resources or prompts. See `/server/help/mcp` for client connection examples. `/.well-known/llms.txt` and the generated OpenAPI document use the same route definitions. `/api/autodiscover` provides the API version, public origins, explicitly exposed overlay addresses, and verified CLI releases.

## Other well-known endpoints

- `/.well-known/security.txt` publishes the security report links, policy URL, and an `Encryption` URL while a publishable PGP key is installed.
- `/.well-known/pgp-key.asc` returns the active ASCII-armored report key only when a key file exists and publication is enabled; otherwise it returns 404.
- `/robots.txt` publishes the crawler policy and sitemap URL.
- `/sitemap.xml` lists eligible public HTML pages and joke details.
- `/feed.xml` and `/rss.xml` serve empty Atom and RSS feeds because embedded jokes do not carry creation timestamps.
- `/favicon.ico` serves a validated configured icon or the embedded default. `server.branding.logo` accepts an HTTPS URL, image data URI, or an image filename in `{config_dir}/assets/`; the default icon is used when no logo is configured.
- `/branding/assets/{filename}` serves validated local branding images from `{config_dir}/assets/`.

## CORS

When `server.cors.allowed_origins` is empty, API CORS uses `Access-Control-Allow-Origin: *`. With an explicit list, it echoes only a matching origin and enables credentialed browser requests for that origin; disallowed preflight requests receive 403. Preflight permits the API methods and the `Authorization`, `Content-Type`, `Accept`, and MCP protocol headers. MCP also validates its `Origin` independently. CORS is not an authentication mechanism.
