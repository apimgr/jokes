# Jokes documentation

Jokes is a Go service with a server-rendered catalog, public read-only joke APIs, a GraphQL endpoint, a companion `jokes-cli` client, and an authenticated Streamable HTTP MCP endpoint. Joke retrieval does not require an account. The service also exposes owner-token management and protected system endpoints; read the [security guide](security.md) before exposing an operator token.

## Start here

- [Install and run the server](installation.md)
- [Configure a deployment](configuration.md)
- [Use the REST and GraphQL APIs](api.md)
- [Use the CLI](cli.md)
- [Security and privacy](security.md)
- [Integrations and current protocol support](integrations.md)
- [Contribute](development.md)
- [Translation guidance](translating.md)

The browser provides home, browse/search, random, category, joke detail, API reference, GraphQL query, about, help, health, privacy, terms, DPO, contact, the security overview/policy/acknowledgments, token-gated security-report status, guest preferences, owner-token management, and preference export/import under `/server/preferences/sync`. Unknown pages and joke IDs render themed 404 responses. Public discovery files include `robots.txt`, `sitemap.xml`, empty Atom and RSS feeds, `/favicon.ico`, `/.well-known/security.txt`, the feature-gated `/.well-known/pgp-key.asc`, and `/.well-known/llms.txt`; the manifest, service worker, and locale catalogs support the installable, offline-enhanced interface.

## License

The application code is licensed under MIT; see the [license and third-party attributions](https://github.com/apimgr/jokes/blob/main/LICENSE.md). The project owner states that the embedded joke corpus is their own work and is in the public domain. This public-domain status is separate from the application's MIT license.
