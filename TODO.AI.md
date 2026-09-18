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
