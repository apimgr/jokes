# TODO.AI.md

## Web UI — missing page content templates

`src/web/templates/` only defines `base.html`, `header.html`, `nav.html`,
`footer.html`, and `index.html` (which supplies the `content` block used by
`base.html`). The following handlers in `src/web/handlers.go` render template
names that have no matching `{{define}}` block anywhere in the embedded
template set, so on gin v1.10.0 (which pushes render errors to `c.Errors`
instead of panicking) they currently return `200 OK` with an empty body:

- `ServeBrowse` — renders `browse.html`
- `ServeRandom` — renders `random.html`
- `ServeCategories` — renders `categories.html`
- `ServeAPIDocs` — renders `api-docs.html`

Add real `{{define "content"}}` templates for these four pages (or a
per-page content block wired through `base.html`) so the routes serve
actual page content instead of a blank 200 response. Covered by
`TestServeBrowse`/`TestServeRandom`/`TestServeCategories`/`TestServeAPIDocs`
in `src/web/handlers_test.go`, which currently assert the (broken) status
quo — update those assertions once real templates are added.

## Missing HTML admin dashboard route

`src/routes/routes.go` registers no `/admin/dashboard` (or any other)
protected admin HTML route — only the JSON API admin group under
`/api/v1/admin/*` (bearer-token protected) and the standalone
`ServeAdminLogin`/`HandleAdminLogin` HTML handlers exist, and neither is
wired into `SetupRoutes`. A GET to `/admin/dashboard` currently falls
through to the custom 404 handler. If a protected HTML admin dashboard is
intended, add the route (with session/cookie auth) and register it in
`SetupRoutes`.

## CasjaysDev Go Linting — Pre-existing Issues

19 convention violations found (not introduced by current CI-fix diff):

### Build & Docker (4 issues)
- Makefile lines 41, 50: `go build` missing inline `-buildvcs=false` flag
- Makefile lines 41, 50: direct host builds bypass Docker — must use container via Make
- Makefile line 100: `go test ./...` run on host — must use `make test` (Docker)

### Makefile Configuration (5 issues)
- Makefile line 4: `PROJECT_NAME` hardcoded "jokes" — infer from git remote
- Makefile line 5: `ORG` hardcoded "apimgr" — infer from git remote
- Makefile: missing required `dev` target
- Makefile line 17: `-trimpath` missing from LDFLAGS
- Makefile line 17: LDFLAGS use wrong variable names (`VERSION`→`Version`, `BUILD_TIME`→`BuildEpoch`, `GIT_COMMIT`→`CommitID`)

### CLI Flags (3 issues)
- src/main.go: missing `-h` short flag (only `--help`)
- src/main.go: missing `--debug` flag
- src/main.go: missing `--color` flag with values auto/yes/no (default auto)

### Output & Logging (1 issue)
- src/main.go: emoji output not gated on `NO_COLOR` env var (multiple lines with embedded emojis)

### Directory Layout (4 issues)
- src/handlers/ → must rename to singular `src/handler/`
- src/models/ → must rename to singular `src/model/`
- src/routes/ → must rename to singular `src/route/`
- src/paths/ → must rename to singular `src/path/`
