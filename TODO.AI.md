# Verified AI.md implementation backlog

This list records requirements verified as missing from the current code during the repository alignment. Completed work is removed rather than retained as done.

## PART-by-PART whole-repository comparison

| AI.md PART | Current repository comparison |
|---|---|
| 0 AI assistant rules | Operational rules were followed for this work; this is not application behavior. |
| 1 Critical rules | Partial: project files exist, but code, tests, route inventory, permissions, and dependency rules still need full conformance review. |
| 2 License and attribution | Partial: application license exists; corpus provenance and item-level redistribution rights are unverified. |
| 3 Project structure | Partial: Go sources, docs, and Claude loaders exist; required scripts, generated metadata, and directory/file inventory need comparison. |
| 4 OS paths | Gap: current path package does not implement the complete per-OS directory contract. |
| 5 Configuration | Gap: legacy config schema and defaults do not implement the required full configuration system. |
| 6 Application modes | Partial: mode support exists; transitions and behavior do not meet the complete contract. |
| 7 Binary requirements | Partial: server and client binaries exist; naming, metadata, hardening, and platform expectations need full comparison. |
| 8 Server CLI | Gap: current flags and command dispatcher cover only a small subset of required commands. |
| 9 Error handling and caching | Partial: recovery exists; canonical error envelope, cache semantics, and public error paths are incomplete. |
| 10 Database | Gap: application data is embedded/in-memory; required SQLite schema, migrations, concurrency, and secrets table are absent. |
| 11 Security and logging | Gap: token ownership, audit/security reporting, complete sanitization and security middleware are absent. |
| 12 Server configuration | Gap: no strict merged config, drop-ins, secrets separation, live reload, or complete schema. |
| 13 Health and versioning | Gap: current health response is not the required dynamic/public-safe contract. |
| 14 API structure | Partial: REST/GraphQL/OpenAPI exist; route inventory, canonical schemas, compatibility rules, and page/API parity remain incomplete. |
| 15 SSL/TLS | Partial: SSL package exists; required lifecycle and Let's Encrypt behavior need verification and completion. |
| 16 Web frontend | Partial: SSR pages and guest preference sync exist; contact/security/MCP contracts, all preference controls, translations, and PWA/accessibility requirements remain incomplete. |
| 17 Email and notifications | Gap: notification settings exist but delivery providers, templates, and workflows do not. |
| 18 Scheduler | Partial: scheduler exists; scheduled jobs contain placeholder/no-op callbacks. |
| 19 GeoIP | Gap: no GeoIP feature implementation. |
| 20 Metrics | Partial: Prometheus metrics and a single Bearer token guard exist; token generation, secret storage, exact route/settings, and service-specific authorization are incomplete. |
| 21 Backup and restore | Partial: command stubs exist; required consistent backup/restore and validation are incomplete. |
| 22 Update command | Partial: update command exists; signature/checksum and safe update contract need verification and completion. |
| 23 Privilege escalation and service | Partial: service command surface exists; escalation/install behavior is incomplete. |
| 24 Service support | Partial: service package exists; supported init systems and lifecycle requirements need comparison. |
| 25 Makefile | Partial: Makefile has build/test/check targets; all required targets, names, cross-platform outputs, and gate details need comparison. |
| 26 Docker | Partial: Docker files exist but do not meet the required complete image, rootfs, compose, entrypoint, and runtime contract. |
| 27 CI/CD | Partial: GitHub workflows exist; required release/security/artifact coverage needs validation and gaps fixed. |
| 28 Testing and development | Partial: unit tests exist; required integration, browser, Docker, clean-OS, and scripted test coverage are missing. |
| 29 ReadTheDocs | Partial: required docs pages and RTD config exist; MkDocs configuration, page accuracy, and docs build remain to be checked. |
| 30 I18N and A11Y | Gap: UI is hardcoded English; catalogs, language negotiation, key parity, RTL and full WCAG requirements are absent. |
| 31 Tor and I2P | Gap: no overlay network management. |
| 32 Client | Partial: minimal CLI exists; full commands, TUI, completion, config, auth, and server-operation parity are absent. |
| 33 IDEA.md reference | Partial: IDEA.md was reconciled with the current AI.md product/auth model, but route inventory and project variables still need full verification. |

## Critical

- [ ] Implement the remaining required `/server/*` public pages and API mirrors, including preferences and its sync flow, security reporting, and DPO surfaces where required by AI.md PART 16.
- [ ] Complete server-rendered frontend behavior from AI.md PART 16: cookie-backed theme/language/preferences and consent, no-JavaScript action paths, page metadata, all page variants, PWA recovery behavior, and responsive/accessibility conformance.
- [ ] Reconcile server configuration with AI.md PARTS 5 and 12: canonical `server.yml`, `yaml.d/`, `secrets.yml`, strict validation, complete required settings, reload behavior, and documentation.
- [ ] Complete API alignment with AI.md PARTS 13–14: canonical health response, OpenAPI and GraphQL parity, documented error formats, content negotiation, and configured security boundaries.
- [ ] Implement and document required operational features from AI.md PARTS 15 and 17–24: TLS lifecycle, notifications, scheduler tasks, metrics protection, backup/restore, updates, privilege handling, and native service support.

## High

- [ ] Implement conditional Tor support and the container/default-aware I2P behavior from AI.md PART 31.
- [ ] Complete the required `jokes-cli` capabilities, interactive mode, shell completions, configuration, and server-operation parity from AI.md PART 32.
- [ ] Add the required integration, browser E2E, Docker, and clean-OS test harnesses from AI.md PART 28 while preserving the instruction that browser E2E is on demand.
- [ ] Add provider-compliant CI workflows for the GitHub remote, including security and vulnerability scans, all platform artifacts, and release workflows under AI.md PART 27.
- [ ] Verify and document dataset provenance and item-level licensing before representing the joke corpus as redistributable.

## Documentation

- [ ] Expand configuration, security, CLI, and feature docs after the related behavior is implemented; current docs describe the verified exposed behavior and must not be treated as a complete AI.md feature reference.
- [ ] Confirm the ReadTheDocs project URL and successful hosted build; the repository config uses the AI.md prescribed `apimgr-jokes` default until the project dashboard is available.
