# Verified AI.md implementation backlog

This list records requirements verified as missing from the current code during the repository alignment. Completed work is removed rather than retained as done.

## PART-by-PART whole-repository comparison

| AI.md PART | Current repository comparison |
|---|---|
| 0 AI assistant rules | Operational rules were followed for this work; this is not application behavior. |
| 1 Critical rules | Partial: public joke IDs now use persisted UUIDv4 values and owner/token persistence exists; remaining whole-repository security, dependency, and route conformance needs implementation review. |
| 2 License and attribution | Partial: application license exists; corpus provenance and item-level redistribution rights are unverified. |
| 3 Project structure | Partial: Go sources, docs, CI providers, tests, and instruction loaders exist; required generated metadata and complete directory inventory need comparison. |
| 4 OS paths | Partial: platform-specific path and privilege files exist; full per-OS directory and permission contract remains to be checked. |
| 5 Configuration | Partial: strict primary YAML loading, private secret separation for operator/encryption/metrics/database/backup credentials, selected drop-in merges, trusted proxies, CORS, contact, nested update settings, and a debounced watcher with an atomic runtime snapshot exist; the canonical full schema, complete drop-in merge, per-setting reload hooks, and remaining settings are incomplete. |
| 6 Application modes | Partial: mode support exists; transitions and behavior do not meet the complete contract. |
| 7 Binary requirements | Partial: server and client binaries exist; naming, metadata, hardening, and platform expectations need full comparison. |
| 8 Server CLI | Gap: current flags and command dispatcher cover only a subset of required commands. |
| 9 Error handling and caching | Partial: recovery exists; canonical error envelope, cache semantics, and public error paths are incomplete. |
| 10 Database | Partial: SQLite schema, migrations, persistent generated secrets, joke IDs, owners, tokens, audit and security-report tables exist; libsql, full required tables, backup integration, and several storage flows remain incomplete. |
| 11 Security and logging | Partial: operator/owner tokens, audit events, URL trust helpers, contact validation, rotating security IDs, encrypted security-report storage, and single-use status tokens exist; report notification/acknowledgment delivery, rate limits, sanitization, secret rotation lifecycle, and remaining middleware are incomplete. |
| 12 Server configuration | Partial: strict primary YAML parsing, secrets separation, drop-in handling for selected security/database settings, trusted proxy/CORS/contact config, and watcher startup exist; the complete canonical schema, all drop-in merging, malformed-file per-file fallback, live application of settings, and many required settings are missing. |
| 13 Health and versioning | Partial: health endpoints and checks are wired; the canonical complete response, accurate subsystem checks, and version behavior still need conformance. |
| 14 API structure | Partial: REST/GraphQL/OpenAPI, token routes, contact mirror, MCP, UUID IDs, and discovery exist; full route inventory, schemas, error/content-negotiation behavior, and parity remain incomplete. |
| 15 SSL/TLS | Partial: SSL package exists; required lifecycle and Let's Encrypt behavior need verification and completion. |
| 16 Web frontend | Partial: registered SSR pages, preferences/sync/tokens, help, two-stage contact forms, configured DPO identity, encrypted security-report status, and persisted single-use CAPTCHA challenges exist; layered contact antispam/rate limits, external CAPTCHA providers, report delivery, configurable page content, full preference behavior, language catalogs, accessibility, and PWA conformance remain incomplete. |
| 17 Email and notifications | Gap: notification settings exist but delivery providers, templates, and workflows do not. |
| 18 Scheduler | Partial: timezone-aware cron/interval parsing, catch-up, retry/panic handling, persistent task/history adapters, config task defaults, and always-on startup exist; main currently registers only `healthcheck_self`, leaving required built-in jobs, audit/notification observers, and full shutdown/config integration incomplete. |
| 19 GeoIP | Gap: no GeoIP feature implementation. |
| 20 Metrics | Partial: isolated Prometheus collectors, bounded HTTP labels, database/token/scheduler hooks, gopsutil system/runtime metrics, per-service generated secrets, authenticated route aliases, Grafana and sanitized Loki handlers exist; rate-limit/cache instrumentation, complete dashboard coverage, configuration validation, and broader app hooks remain incomplete. |
| 21 Backup and restore | Partial: archive creation checkpoints local SQLite and writes a verified root-layout archive with config, secrets, optional data/SSL, external certificate files, and optional encryption; retention, full restore/rollback, remote-cache handling, compliance gates, metadata/audit, and TLS path rewriting remain incomplete. |
| 22 Update command | Partial: channel checking, HTTPS SHA-256 verified staging, elevated startup apply, binary replacement, restart helpers, branch persistence, scheduled deferred checks, and exit-75 restart are integrated; audit/notification events, service-aware restart, and failure recovery coverage remain incomplete. |
| 23 Privilege escalation and service | Partial: service command surface exists; escalation/install behavior is incomplete. |
| 24 Service support | Partial: service package exists; supported init systems and lifecycle requirements need comparison. |
| 25 Makefile | Partial: the Makefile now has the six required core targets plus `clean` and `release-check`, a bounded single-container cross-platform build, and coverage enforcement. Release/tag lifecycle and Docker publishing behavior still need a full requirement-by-requirement verification. |
| 26 Docker | Partial: Docker files exist but do not meet the required complete image, rootfs, compose, entrypoint, and runtime contract. |
| 27 CI/CD | Partial: GitHub, GitLab, Gitea, Forgejo, and Jenkins workflow files exist; required release/security/artifact coverage and action pinning still need validation. |
| 28 Testing and development | Partial: unit, Docker, Incus, and browser E2E harnesses exist; prescribed integration/clean-OS coverage and full scripted execution remain to be checked. |
| 29 ReadTheDocs | Partial: documentation pages, RTD config, and root `mkdocs.yml` exist; page accuracy and successful documentation builds remain to be verified. |
| 30 I18N and A11Y | Partial: six embedded locale catalogs, language negotiation, RTL direction, and translated CLI text exist; key parity, complete SSR/browser UI translation, and full WCAG conformance remain incomplete. |
| 31 Tor and I2P | Gap: no overlay network management. |
| 32 Client | Partial: CLI commands, UUID lookup, authentication, and scripts exist; full command/TUI/completion/config and server-operation parity remain incomplete. |
| 33 IDEA.md reference | Partial: IDEA.md was reconciled with the current AI.md product/auth model, but route inventory and project variables still need full verification. |

## Critical

- [ ] Finish `/server/*` page/API parity: configurable page content, privacy tracking metadata, PGP availability, canonical per-page JSON schemas, and exact error/content-negotiation behavior; wire security-report notification and acknowledgment delivery.
- [ ] Complete remaining PART16 behavior: layered contact antispam/rate limits, configured external CAPTCHA providers, full consent/language coverage, the PWA update watchdog and recovery flow, and final accessibility checks across every page.
- [ ] Reconcile server configuration with AI.md PARTS 5 and 12: canonical `server.yml`, `yaml.d/`, `secrets.yml`, strict validation, complete required settings, reload behavior, and documentation.
- [ ] Complete API alignment with AI.md PARTS 13–14: canonical health response, OpenAPI and GraphQL parity, documented error formats, content negotiation, and configured security boundaries.
- [ ] Implement and document required operational features from AI.md PARTS 15 and 17–24: TLS lifecycle, notifications, scheduler tasks, metrics protection, backup/restore, updates, privilege handling, and native service support.

## High

- [ ] Implement conditional Tor support and the container/default-aware I2P behavior from AI.md PART 31.
- [ ] Complete the required `jokes-cli` capabilities, interactive mode, shell completions, configuration, and server-operation parity from AI.md PART 32.
- [ ] Complete the required integration, browser E2E, Docker, and clean-OS test coverage from AI.md PART 28; harnesses now exist, with browser E2E kept on demand.
- [ ] Add provider-compliant CI workflows for the GitHub remote, including security and vulnerability scans, all platform artifacts, and release workflows under AI.md PART 27.
- [ ] Verify and document dataset provenance and item-level licensing before representing the joke corpus as redistributable.

## Documentation

- [ ] Finish documentation after implementation gaps close; README and current API, CLI, configuration, installation, integrations, translation, and security pages now describe the verified behavior without claiming incomplete features are available.
- [ ] Confirm the ReadTheDocs project URL and successful hosted build; the repository config uses the AI.md prescribed `apimgr-jokes` default until the project dashboard is available.
