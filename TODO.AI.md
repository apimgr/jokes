# Verified AI.md implementation backlog

This list records requirements verified as missing from the current code during the repository alignment. Completed work is removed rather than retained as done.

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
