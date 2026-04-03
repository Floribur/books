# Phase 2: Go REST API - Context

**Gathered:** 2026-04-03
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver all Go HTTP endpoints consumed by the React frontend: books, authors, genres, years listing, cover file serving, and the production server topology (go:embed React build, SPA catch-all, Open Graph meta tag injection). Admin endpoints from Phase 1 remain unchanged.

</domain>

<decisions>
## Implementation Decisions

### JSON Response Shape

- **D-01:** All paginated list endpoints use an **envelope object** — `{"items": [...], "next_cursor": "...", "has_more": true}`. When no more pages exist, `next_cursor` is `null` and `has_more` is `false`. This is the contract TanStack Query's `getNextPageParam` will read.
- **D-02:** The cursor value is an **opaque base64 token** encoding `read_at + id`. The frontend passes it as `?cursor=<token>`. Go decodes it server-side. This hides DB internals and is safe to use in URLs without additional encoding.
- **D-03:** Non-paginated list endpoints (e.g. `GET /api/years`, `GET /api/authors`) return a plain JSON array `[...]` — no envelope needed since they are not paginated.

### Book Object Shape

- **D-04:** On `GET /api/books` (list), each book item includes: `slug`, `title`, `cover_path`, `read_at`, `publication_year`, plus **inline arrays**: `authors: [{name, slug}]` and `genres: [{name, slug}]`. No extra round trips from the frontend for card display.
- **D-05:** On `GET /api/books/:slug` (detail), return all book fields: `slug`, `title`, `cover_path`, `read_at`, `publication_year`, `description`, `page_count`, `isbn13`, `read_count`, `shelf`, `metadata_source`, plus inline `authors` and `genres` arrays.
- **D-06:** For author/genre list endpoints (`GET /api/authors`, `GET /api/genres`), each item includes `name`, `slug`, and `book_count`.
- **D-07:** For author/genre detail endpoints (`GET /api/authors/:slug`, `GET /api/genres/:slug`), include the entity fields plus a paginated `items` list of books (same shape as D-04).

### Open Graph Meta Tag Injection

- **D-08:** Full OG injection implemented in Phase 2. Go intercepts `/books/:slug` URL paths (non-API, non-asset), queries the DB for the book, and uses `text/template` to inject 5 meta tags into `index.html` before serving:
  - `og:title` — book title
  - `og:description` — book description (truncated to ~200 chars if needed)
  - `og:image` — absolute URL to cover image
  - `og:type` — `"book"`
  - `og:url` — canonical URL for the page
- **D-09:** The `SITE_URL` env var provides the base URL for building absolute OG image URLs. Default: `http://localhost:8081` for dev. Required in production.
- **D-10:** All other non-API, non-asset paths serve `index.html` without modification (standard SPA catch-all).

### Claude's Discretion

- Error response format — use `{"error": "message"}` for all 4xx/5xx responses
- CORS configuration — allow `http://localhost:5173` in dev (per roadmap); no CORS in production (same origin)
- Page size default and maximum — Claude decides reasonable values (e.g. default 24, max 100)
- sqlc query design — JOIN approach for inline authors/genres
- Handler struct organization within `internal/api/` — follow existing `AdminHandlers` pattern

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Requirements
- `.planning/REQUIREMENTS.md` §API — API-01 through API-10: full endpoint list with descriptions

### Roadmap
- `.planning/ROADMAP.md` §Phase 2 — Plan 2.1 and 2.2 task breakdown

### Existing Code Patterns
- `cmd/server/main.go` — Chi router setup, middleware stack, dependency wiring
- `internal/api/admin.go` — Handler struct pattern with `*db.Queries` dependency
- `internal/db/models.go` — sqlc-generated Go types (pgtype.Timestamptz, *string for nullables)
- `internal/db/books.sql.go` — Existing book queries (GetBookBySlug, GetAllGoodreadsIDs)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `chi.NewRouter()` with `middleware.Logger` + `middleware.Recoverer` — already configured in `cmd/server/main.go`; new API routes mount on the same router
- `internal/api.AdminHandlers` struct — established pattern: struct holds `*db.Queries`, methods are `http.HandlerFunc`-compatible; replicate for public API handlers
- `db.New(pool)` — `*db.Queries` is the single DB access object; injected into handler structs at startup

### Established Patterns
- sqlc generated types: nullable columns are `*string`, `*int32`; timestamps are `pgtype.Timestamptz` — API response types will need custom marshaling or explicit mapping to avoid exposing pgtype internals
- `internal/sync/` packages all use `*db.Queries` directly — same pattern for API handlers

### Integration Points
- `cmd/server/main.go`: add `r.Mount("/api", apiRouter)` and `r.Handle("/covers/*", ...)` and catch-all for SPA
- Phase 5 will run `make build-frontend` before `go build` — the `//go:embed frontend/dist` directive added here enables that
- `data/covers/` directory holds cover images served by the `/covers/*` file server

</code_context>

<specifics>
## Specific Ideas

- Cursor shape confirmed during discussion: base64-encode `"<read_at_rfc3339>_<id>"` — simple and debuggable when decoded
- OG injection was clarified: it enables rich link previews when sharing book URLs on WhatsApp, iMessage, Discord etc. — social sharing use case, not primarily SEO

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 02-go-rest-api*
*Context gathered: 2026-04-03*
