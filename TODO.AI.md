# Remaining implementation work

## Product data provenance and completeness

- The embedded dataset contains 5,160 single-part entries across 16 categories. It has no setup/delivery records, so the two-part content requirement cannot be completed without a vetted dataset and provenance/license terms. Do not generate or relabel joke text to fill this gap.
- Confirm and document the dataset source and item-level redistribution terms. The repository MIT license covers the application source but does not establish rights to the embedded joke content.

## GraphQL correctness

- Replace `executeGraphQL`'s substring-based query handling with a standards-compliant GraphQL parser/executor using `src/graphql/schema.go` and `src/graphql/resolver.go`. The current endpoint only recognizes a few field-name substrings, ignores most arguments and selection sets, and can return incorrect data for valid multi-field queries. Preserve the no-external-runtime-service requirement.
- Add integration coverage for query validation, variables, selected fields, category filtering, pagination, and malformed requests after the parser/executor is implemented.

## Client and operator features required by AI.md

- Complete the server CLI standard flags and output behavior (`--debug`, `--color auto|yes|no`, and `NO_COLOR`) required by AI.md.
- Complete the CLI client requirements in AI.md PART 32: user-scoped `cli.yml` configuration, server profiles and selection, interactive TUI, completions/man page, and full parity with supported read-only API features.
- Replace the unimplemented backup, restore, update, maintenance-mode persistence, and configuration-save command responses in `src/main.go`, or remove those commands if they are not part of the project scope under AI.md.
- Either implement actual scheduled work or keep scheduler and notification scheduling disabled and remove the no-op task callbacks. Certificate renewal and notifications are not part of IDEA.md's external-service-free product scope.

## Dataset schema and API docs

- Add source-backed two-part records and update the embedded data model only when content provenance is established.
- Keep `src/swagger/swagger.go` aligned with handler behavior; the current generated spec still needs schema examples, complete query parameters, and response/error definitions.

## Packaging and release integration

- Update and reconcile installer scripts and release workflows with the server and CLI binary names produced by the Makefile (`jokes` and `jokes-cli`) and the actual release archive layout.
- Configure CI image annotations/labels and exercise the production/development Docker builds through the primary agent's validation workflow.
