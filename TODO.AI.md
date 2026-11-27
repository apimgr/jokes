# Verified AI.md implementation backlog

This list records requirements verified as missing from the current code during the repository alignment. Completed work is removed rather than retained as done.

## PART-by-PART whole-repository comparison

| AI.md PART | Current repository comparison |
|---|---|
| 0 AI assistant rules | Operational rules were followed for this work; this is not application behavior. |
| 1 Critical rules | Partial: public joke IDs now use persisted UUIDv4 values and owner/token persistence exists; remaining whole-repository security, dependency, and route conformance needs implementation review. |
| 2 License and attribution | Partial: application license exists; corpus provenance and item-level redistribution rights are unverified. |
| 3 Project structure | Partial: Go sources, docs, CI providers, tests, and instruction loaders exist; required generated metadata and complete directory inventory need comparison. |
| 4 OS paths | Partial: system and user path mode is locked at startup, systemd/OpenRC system containers resolve host paths, and Linux runtime directories are assigned to the service account before privilege drop; full per-OS directory and permission contract remains to be checked. |
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
| 16 Web frontend | Partial: registered SSR pages, preferences/sync/tokens, help, two-stage contact forms, configured DPO identity, encrypted security-report status, persisted single-use CAPTCHA challenges, and sanitized configurable page content exist; layered contact antispam/rate limits, external CAPTCHA providers, report delivery, privacy tracking metadata, full preference behavior, language catalogs, accessibility, and PWA conformance remain incomplete. |
| 17 Email and notifications | Partial: SMTP and webhook transports and several event deliveries exist. A reusable embedded/custom template engine now includes all 16 PART 17 default templates, per-send override reload, strict variable validation, safe HTML rendering, and operator email rendering for backup failure and disk-low events. The CLI dry-run/test commands, complete event routing, localized messages, startup/reload status behavior, and all message flows remain incomplete. |
| 18 Scheduler | Partial: timezone-aware cron/interval parsing, catch-up, retry/panic handling, persistent task/history adapters, config task defaults, and always-on startup exist; `geoip_update`, `update_check`, and `healthcheck_self` are registered, while most required built-in jobs, audit/notification observers, and full shutdown/config integration remain. |
| 19 GeoIP | Partial: checksum-verified DB-IP/iptoasn MMDB loading/updating, startup loading, config, after-rate-limit deny-only middleware, allowlist/private-IP bypass, cached health, first-run update, and scheduled update exist; health/API conformance and additional requirements need verification. |
| 20 Metrics | Partial: isolated Prometheus collectors, bounded HTTP labels, database/token/scheduler hooks, gopsutil system/runtime metrics, per-service generated secrets, authenticated route aliases, Grafana and sanitized Loki handlers exist; rate-limit/cache instrumentation, complete dashboard coverage, configuration validation, and broader app hooks remain incomplete. |
| 21 Backup and restore | Partial: archive creation checkpoints local SQLite and writes a verified root-layout archive with config, secrets, optional data/SSL, external certificate files, and optional encryption; retention, incremental backups, restore, remote-cache handling, compliance gates, metadata/audit, and TLS path rewriting are implemented in part. Startup now runs a warning-only backup-directory create/remove probe as the service identity and backup failures carry the same path/owner/mode/user detail. Full PART 21 requirement-by-requirement verification remains. |
| 22 Update command | Partial: channel checking, HTTPS SHA-256 verified staging, elevated startup apply, binary replacement, restart helpers, branch persistence, scheduled deferred checks, and exit-75 restart are integrated; audit/notification events, service-aware restart, and failure recovery coverage remain incomplete. |
| 23 Privilege escalation and service | Partial: Linux root startup now prepares the dedicated account and directories, binds privileged sockets before dropping to the service identity, and runs overlay providers as that identity; interactive escalation and the remaining platform-specific privilege contract still need implementation. |
| 24 Service support | Partial: systemd and OpenRC install/start/restart/stop/disable/uninstall pass lifecycle integration, including systemd root startup, OpenRC parent-directory setup, distinct supervisor/application PID files, and conditional account removal; other managers and full lifecycle conformance remain. |
| 25 Makefile | Partial: the Makefile now has the six required core targets plus `clean` and `release-check`, a bounded single-container cross-platform build, and coverage enforcement. Release/tag lifecycle and Docker publishing behavior still need a full requirement-by-requirement verification. |
| 26 Docker | Partial: Docker files exist but do not meet the required complete image, rootfs, compose, entrypoint, and runtime contract. |
| 27 CI/CD | Partial: GitHub, GitLab, Gitea, Forgejo, and Jenkins workflow files exist; required release/security/artifact coverage and action pinning still need validation. |
| 28 Testing and development | Partial: unit tests with 60.0% aggregate coverage, Docker smoke, Chromium browser E2E, and both systemd/OpenRC Incus service integrations pass; additional prescribed integration and clean-OS coverage still needs comparison. |
| 29 ReadTheDocs | Partial: documentation pages, RTD config, and root `mkdocs.yml` exist; strict local MkDocs build passes, while page accuracy and hosted build still need verification. |
| 30 I18N and A11Y | Partial: six embedded locale catalogs, language negotiation, RTL direction, and translated CLI text exist; key parity, complete SSR/browser UI translation, and full WCAG conformance remain incomplete. |
| 31 Tor and I2P | Partial: dedicated Tor and I2P provider helpers, app-owned backends, config defaults, startup/shutdown, health reporting, discovery host resolution, and scheduled health checks are wired; automatic provider recovery, configured address exposure, bandwidth threshold tracking, vanity generation, and full lifecycle tests remain. |
| 32 Client | Partial: CLI commands, UUID lookup, authentication, and scripts exist; full command/TUI/completion/config and server-operation parity remain incomplete. |
| 33 IDEA.md reference | Partial: IDEA.md was reconciled with the current AI.md product/auth model, but route inventory and project variables still need full verification. |

## Critical

- [ ] Finish `/server/*` page/API parity: configurable page content, privacy tracking metadata, PGP availability, canonical per-page JSON schemas, and exact error/content-negotiation behavior; security-report notifications, researcher acknowledgments, and status updates are wired but still need full contract verification.
- [ ] Complete remaining PART16 behavior: layered contact antispam/rate limits, configured external CAPTCHA providers, full consent/language coverage, the PWA update watchdog and recovery flow, and final accessibility checks across every page.
- [ ] Reconcile server configuration with AI.md PARTS 5 and 12: canonical `server.yml`, `yaml.d/`, `secrets.yml`, strict validation, complete required settings, reload behavior, and documentation.
- [ ] Complete API alignment with AI.md PARTS 13–14: canonical health response, OpenAPI and GraphQL parity, documented error formats, content negotiation, and configured security boundaries.
- [ ] Finish AI.md PART 19 verification: confirm health/API status parity and all failure, stale-data, and configuration behavior; attribution and scheduled loading/update are implemented.
- [ ] Implement and document required operational features from AI.md PARTS 15 and 17–24: TLS lifecycle, customizable email templates and complete notification workflows, scheduler tasks, metrics protection, backup/restore, updates, privilege handling, and native service support.

## High

- [ ] Complete AI.md PART 31: provider recovery, configured autodiscovery exposure, bandwidth thresholds, vanity address support, and provider-backed lifecycle tests.
- [ ] Complete the required `jokes-cli` capabilities, interactive mode, shell completions, configuration, and server-operation parity from AI.md PART 32.
- [ ] Compare remaining integration and clean-OS coverage with AI.md PART 28; unit, Docker smoke, Chromium browser E2E, and systemd/OpenRC Incus integration gates have now passed.
- [ ] Add provider-compliant CI workflows for the GitHub remote, including security and vulnerability scans, all platform artifacts, and release workflows under AI.md PART 27.
- [ ] Verify and document dataset provenance and item-level licensing before representing the joke corpus as redistributable.

## Documentation

- [ ] Finish documentation after implementation gaps close; README and current API, CLI, configuration, installation, integrations, translation, and security pages now describe the verified behavior without claiming incomplete features are available. The strict local MkDocs build passes.
- [ ] Confirm the ReadTheDocs project URL and successful hosted build; the repository config uses the AI.md prescribed `apimgr-jokes` default until the project dashboard is available.
