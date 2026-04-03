# Phase 2: Go REST API - Research

**Researched:** 2026-04-03
**Domain:** Go HTTP (Chi v5), sqlc, cursor-based pagination, go:embed, Open Graph injection
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** All paginated list endpoints use an envelope object `{"items": [...], "next_cursor": "...", "has_more": true}`. When no more pages exist, `next_cursor` is `null` and `has_more` is `false`.
- **D-02:** Cursor value is an opaque base64 token encoding `read_at + id`. Frontend passes it as `?cursor=<token>`. Go decodes server-side.
- **D-03:** Non-paginated list endpoints (e.g. `GET /api/years`, `GET /api/authors`) return a plain JSON array `[...]`.
- **D-04:** `GET /api/books` list items include: `slug`, `title`, `cover_path`, `read_at`, `publication_year`, inline `authors: [{name, slug}]` and `genres: [{name, slug}]`.
- **D-05:** `GET /api/books/:slug` returns all book fields: `slug`, `title`, `cover_path`, `read_at`, `publication_year`, `description`, `page_count`, `isbn13`, `read_count`, `shelf`, `metadata_source`, plus inline `authors` and `genres` arrays.
- **D-06:** Author/genre list endpoints each item includes `name`, `slug`, and `book_count`.
- **D-07:** Author/genre detail endpoints include entity fields plus paginated `items` list of books (same shape as D-04).
- **D-08:** Full OG injection implemented in Phase 2. Go intercepts `/books/:slug` paths, queries DB, uses `text/template` to inject 5 meta tags: `og:title`, `og:description`, `og:image`, `og:type`, `og:url`.
- **D-09:** `SITE_URL` env var provides base URL for absolute OG image URLs. Default: `http://localhost:8081`. Required in production.
- **D-10:** All other non-API, non-asset paths serve `index.html` without modification (standard SPA catch-all).

### Claude's Discretion

- Error response format: use `{"error": "message"}` for all 4xx/5xx responses
- CORS configuration: allow `http://localhost:5173` in dev; no CORS in production (same origin)
- Page size default and maximum: Claude decides reasonable values (e.g. default 24, max 100)
- sqlc query design: JOIN approach for inline authors/genres
- Handler struct organization within `internal/api/`: follow existing `AdminHandlers` pattern

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope.
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| API-01 | `GET /api/books` returns cursor-paginated list (keyed on `read_at` + `id`), descending order | Keyset pagination pattern with composite cursor; new sqlc queries needed |
| API-02 | `GET /api/books/currently-reading` returns all currently-reading books | Simple filter on `shelf = 'currently-reading'`; no pagination needed |
| API-03 | `GET /api/books/:slug` returns full book detail | Extend existing `GetBookBySlug` with author/genre JOINs |
| API-04 | `GET /api/authors` returns all authors with book counts | New sqlc query with COUNT JOIN; plain array response |
| API-05 | `GET /api/authors/:slug` returns author detail with their books (cursor-paginated) | New sqlc queries for author lookup + books by author |
| API-06 | `GET /api/genres` returns all genres with book counts | New sqlc query with COUNT JOIN; plain array response |
| API-07 | `GET /api/genres/:slug` returns genre detail with its books (cursor-paginated) | New sqlc queries for genre lookup + books by genre |
| API-08 | `GET /api/years` returns distinct years with book counts | New sqlc query; plain array response |
| API-09 | `/covers/:filename` serves cover images with `Cache-Control: immutable` headers | Middleware wrapper around `http.FileServer` |
| API-10 | Per-book Open Graph meta tags injected into SPA HTML head by Go | `text/template` renders into `index.html` bytes before write; path detection on `/books/:slug` |
</phase_requirements>

---

## Summary

This phase builds on a fully working Phase 1 codebase. The Chi router, middleware stack, `*db.Queries` injection pattern, and sqlc toolchain are all established. The primary work is (1) writing new sqlc SQL queries for read operations that did not exist in Phase 1, (2) implementing handler structs following the `AdminHandlers` pattern, and (3) wiring cover-file serving, go:embed placeholder, SPA catch-all, and Open Graph injection into `cmd/server/main.go`.

The most nuanced work is cursor-based pagination. The cursor encodes `read_at` (RFC3339) and `id` as a base64 string. The SQL keyset predicate is `(read_at, id) < ($cursor_time, $cursor_id)` (descending order means `<` not `>`). sqlc cannot auto-generate parameterized keyset queries — these must be hand-written SQL with careful parameter typing for `pgtype.Timestamptz`.

Open Graph injection requires reading `index.html` from the embedded FS at startup, parsing it once with `text/template` (inserting a `{{.OGTags}}` placeholder into a pre-processed copy), then rendering per-request. The alternative — simple `strings.Replace` — is equally valid given the limited template surface and avoids any risk of template injection from book content (description must be HTML-escaped).

**Primary recommendation:** Write all new SQL queries first, run `sqlc generate`, then implement handlers against generated types. This avoids the anti-pattern of writing handlers against types that do not yet exist.

---

## Standard Stack

### Core (already in go.mod)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `github.com/go-chi/chi/v5` | v5.2.5 | HTTP router and middleware | Already in use; idiomatic net/http |
| `github.com/jackc/pgx/v5` | v5.9.1 | PostgreSQL driver | Already in use; best-in-class pgx |
| sqlc | v1.30.0 | Type-safe SQL code generation | Already configured via sqlc.yaml |

### Needs Adding

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `github.com/go-chi/cors` | latest (~v1.2.x) | CORS middleware for Chi | Dev-only: allow `localhost:5173` to call `localhost:8081` |

**Installation:**
```bash
go get github.com/go-chi/cors
```

**Version verification:**
```bash
go list -m github.com/go-chi/cors
```

### Standard Library Packages (no install needed)

| Package | Purpose |
|---------|---------|
| `encoding/base64` | Encode/decode opaque cursor tokens |
| `embed` | `//go:embed frontend/dist` for production SPA serving |
| `text/template` | Open Graph meta tag injection into index.html |
| `net/http` | `http.FileServer`, `http.StripPrefix` for cover serving |
| `io/fs` | `fs.Sub` to scope embedded FS to `frontend/dist` subdirectory |
| `strings` | `strings.Replace` as simpler alternative to template for OG placeholder |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| `github.com/go-chi/cors` | `github.com/rs/cors` | go-chi/cors is a maintained fork of rs/cors; either works, go-chi/cors is the idiomatic pairing |
| `text/template` for OG injection | `strings.Replace` | strings.Replace is simpler and sufficient; text/template adds complexity without benefit for 5 static placeholders |
| Hand-written cursor SQL | kuysor/cursor library | Hand-written is clearer and avoids a dependency; library justified only if pagination spans many entity types |

---

## Architecture Patterns

### Recommended Project Structure

```
internal/
├── api/
│   ├── admin.go          # Existing — AdminHandlers (unchanged)
│   ├── public.go         # New — PublicHandlers struct, all GET endpoints
│   └── response.go       # New — shared response types (BookListItem, BookDetail, etc.)
sql/
└── queries/
    ├── books.sql         # Extend with cursor-paginated queries
    ├── authors.sql       # Extend with list + detail queries
    └── genres.sql        # Extend with list + detail queries
cmd/
└── server/
    └── main.go           # Wire PublicHandlers, covers, embed, OG handler
```

### Pattern 1: Handler Struct (follow AdminHandlers exactly)

```go
// internal/api/public.go
type PublicHandlers struct {
    Queries *db.Queries
    SiteURL string // from SITE_URL env var — needed for og:image absolute URLs
}

func (h *PublicHandlers) GetBooks(w http.ResponseWriter, r *http.Request) { ... }
func (h *PublicHandlers) GetBookBySlug(w http.ResponseWriter, r *http.Request) { ... }
// ... one method per endpoint
```

Wire in `main.go`:
```go
pub := &api.PublicHandlers{
    Queries: queries,
    SiteURL: siteURL,
}
r.Get("/api/books", pub.GetBooks)
r.Get("/api/books/currently-reading", pub.GetCurrentlyReading)
r.Get("/api/books/{slug}", pub.GetBookBySlug)
// etc.
```

### Pattern 2: Cursor Encoding/Decoding

Decision D-02 specifies `base64("<read_at_rfc3339>_<id>")`.

```go
// encode
raw := fmt.Sprintf("%s_%d", readAt.Format(time.RFC3339Nano), id)
cursor := base64.URLEncoding.EncodeToString([]byte(raw))

// decode
decoded, err := base64.URLEncoding.DecodeString(cursorParam)
// decoded = "2024-11-15T10:30:00Z_1234"
parts := strings.SplitN(string(decoded), "_", 2)
readAt, err := time.Parse(time.RFC3339Nano, parts[0])
id, err := strconv.ParseInt(parts[1], 10, 64)
```

Use `base64.URLEncoding` (not `StdEncoding`) — the standard `+/` chars are not URL-safe without additional escaping. URL encoding uses `-_` instead.

### Pattern 3: Keyset Pagination SQL (hand-written, not generated by sqlc)

For descending `read_at, id` order, the keyset predicate is:

```sql
-- name: ListBooksPaginated :many
SELECT b.id, b.slug, b.title, b.cover_path, b.read_at, b.publication_year,
       b.shelf
FROM books b
WHERE b.shelf = 'read'
  AND ($1::timestamptz IS NULL OR (b.read_at, b.id) < ($1::timestamptz, $2::bigint))
ORDER BY b.read_at DESC NULLS LAST, b.id DESC
LIMIT $3;
```

With `$1 = NULL, $2 = 0` for the first page (no cursor). sqlc maps `$1::timestamptz` to `pgtype.Timestamptz` and `$2::bigint` to `int64`.

The `has_more` flag requires fetching `limit + 1` rows and checking if the extra row exists (then dropping it from the response). This is cleaner than a separate COUNT query.

**Important:** The index `idx_books_read_at` already exists on `books(read_at DESC NULLS LAST)`. A composite index `(read_at DESC NULLS LAST, id DESC)` would be more efficient for the keyset predicate — add this in a new migration.

### Pattern 4: Inline Authors/Genres via JSON aggregation

Rather than N+1 queries, use PostgreSQL's `json_agg` to fetch authors and genres inline:

```sql
-- name: GetBookDetailBySlug :one
SELECT
    b.id, b.slug, b.title, b.cover_path, b.read_at, b.publication_year,
    b.description, b.page_count, b.isbn13, b.read_count, b.shelf, b.metadata_source,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', a.name, 'slug', a.slug))
        FILTER (WHERE a.id IS NOT NULL), '[]'
    ) AS authors,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', g.name, 'slug', g.slug))
        FILTER (WHERE g.id IS NOT NULL), '[]'
    ) AS genres
FROM books b
LEFT JOIN book_authors ba ON ba.book_id = b.id
LEFT JOIN authors a ON a.id = ba.author_id
LEFT JOIN book_genres bg ON bg.book_id = b.id
LEFT JOIN genres g ON g.id = bg.genre_id
WHERE b.slug = $1
GROUP BY b.id;
```

**sqlc caveat:** sqlc will generate this as `json.RawMessage` for the `json_agg` columns. The handler must unmarshal these into `[]AuthorRef` / `[]GenreRef` before embedding in the API response. This is the standard pattern when sqlc cannot infer a Go struct for aggregated JSON columns.

### Pattern 5: Cover File Server with Immutable Cache Headers

```go
// In main.go
coversDir := http.Dir("data/covers")
fileServer := http.FileServer(coversDir)
r.Handle("/covers/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
    http.StripPrefix("/covers/", fileServer).ServeHTTP(w, r)
}))
```

Chi's wildcard `*` in `/covers/*` matches everything after `/covers/`. Use `http.StripPrefix` to remove the `/covers/` prefix before the file server resolves the path.

### Pattern 6: go:embed with SPA Catch-All

```go
//go:embed frontend/dist
var frontendFS embed.FS

// Strip the "frontend/dist" prefix so the FS root is the dist directory
distFS, _ := fs.Sub(frontendFS, "frontend/dist")
fileServer := http.FileServer(http.FS(distFS))

// In router wiring:
r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Try to serve the file; fall back to index.html for SPA routing
    f, err := distFS.Open(strings.TrimPrefix(r.URL.Path, "/"))
    if err != nil {
        // File not found → serve index.html (SPA catch-all)
        r.URL.Path = "/"
    }
    if f != nil { f.Close() }
    fileServer.ServeHTTP(w, r)
}))
```

During Phase 2, `frontend/dist` does not exist yet. Use a placeholder approach:
- Create `frontend/dist/.gitkeep` to satisfy the embed directive
- Add `//go:embed frontend/dist` but wrap SPA serving in a build tag or a runtime check that falls back gracefully when `index.html` is missing

**Simpler alternative for Phase 2:** Skip the SPA embed entirely in Phase 2 (serve a minimal placeholder `index.html` from disk, not embedded). The `//go:embed` directive is implemented as a compile-time feature — the build will fail if the directory is empty or missing. Use `os.DirFS("frontend/dist")` pointing to a placeholder directory for now; Phase 5 switches to `embed.FS`.

### Pattern 7: Open Graph Meta Tag Injection

```go
// At startup: read index.html, prepare template
indexHTML, _ := fs.ReadFile(distFS, "index.html")
// Insert placeholder before </head>
templateSrc := strings.Replace(
    string(indexHTML),
    "</head>",
    "{{.OGTags}}</head>",
    1,
)
ogTemplate, _ := template.New("index").Parse(templateSrc)

// In handler: detect /books/:slug path
r.Get("/books/{slug}", func(w http.ResponseWriter, r *http.Request) {
    slug := chi.URLParam(r, "slug")
    book, err := queries.GetBookBySlug(r.Context(), slug)
    if err != nil {
        // fall back to plain index.html
        w.Write(indexHTML)
        return
    }
    ogTags := renderOGTags(book, siteURL)
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    ogTemplate.Execute(w, map[string]string{"OGTags": ogTags})
})
```

`renderOGTags` returns a pre-escaped HTML string of 5 `<meta>` tags. Because `text/template` auto-escapes `{{.OGTags}}` as plain text (not HTML), use `template.HTML(ogTags)` or switch to `html/template` and mark the value as `template.HTML` to prevent double-escaping.

**Simpler alternative:** Skip `text/template` entirely. Use `strings.Replace(string(indexHTML), "</head>", ogTagsHTML+"</head>", 1)` — five lines of code, zero template risk.

### Anti-Patterns to Avoid

- **N+1 queries for authors/genres:** Never load a list of 24 books then query authors for each. Use `json_agg` in the initial query.
- **Offset pagination:** `OFFSET N` is O(N) in PostgreSQL; breaks on inserts mid-page. Use keyset.
- **Exposing pgtype.Timestamptz in JSON:** The `pgtype.Timestamptz` struct serializes as `{"Time":"...","Valid":true}`. Map to `time.Time` or `string` in response types before JSON encoding.
- **Using `//go:embed` with empty directory:** The build fails at compile time if the embedded path does not exist or contains no files. Always have at least a `.gitkeep` file or use `os.DirFS` during development.
- **Registering `/books/:slug` before `/books/currently-reading` in Chi:** Chi matches routes in registration order for conflicting patterns. Register specific routes first.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| CORS preflight handling | Custom `OPTIONS` handler | `github.com/go-chi/cors` | Handles preflight, wildcard methods, credentials, Vary headers correctly |
| Cursor encoding | Custom binary format | `encoding/base64` (stdlib) | Opaque, URL-safe, debuggable when decoded |
| File content-type detection | Custom MIME mapping | `http.FileServer` | Uses Go's `mime` package with correct sniffing |
| Graceful nil handling in JSON | Custom marshal code | Explicit response struct with `*string` fields | sqlc already uses `*string`; copy the pattern |

**Key insight:** The pagination and OG injection are the most custom work in this phase — everything else has a stdlib or existing-codebase answer.

---

## Common Pitfalls

### Pitfall 1: `pgtype.Timestamptz` leaks into JSON response

**What goes wrong:** Handler returns `db.Book` directly; frontend receives `{"read_at":{"Time":"...","Valid":true}}` instead of `"read_at":"2024-11-15T10:30:00Z"`.

**Why it happens:** `pgtype.Timestamptz` implements `json.Marshaler` to expose its internal struct, not an RFC3339 string.

**How to avoid:** Define explicit API response structs (`BookListItem`, `BookDetail`) that map `pgtype.Timestamptz` → `time.Time` (which marshals as RFC3339 by default) or `string`. Never return `db.Book` directly from handlers.

**Warning signs:** `curl /api/books` shows nested objects for timestamp fields.

### Pitfall 2: Cursor predicate direction

**What goes wrong:** Using `>` instead of `<` in the keyset predicate for descending order. Returns no rows or skips to the wrong page.

**Why it happens:** Keyset direction depends on sort order. For `ORDER BY read_at DESC`, the "next page" is items *before* the cursor in time (i.e., `read_at < cursor_time`).

**How to avoid:** Always verify: `ORDER BY col DESC` → keyset uses `<`. `ORDER BY col ASC` → keyset uses `>`.

**Warning signs:** Second page returns empty or same books as first page.

### Pitfall 3: `//go:embed` fails at compile time if directory is empty

**What goes wrong:** `//go:embed frontend/dist` causes `go build` to fail with `pattern frontend/dist: no matching files found` when the React build has not been run.

**Why it happens:** `embed` validates patterns at compile time.

**How to avoid:** Either create `frontend/dist/.gitkeep` (embedded, but the SPA handler gracefully 404s on missing `index.html`), or use `os.DirFS("frontend/dist")` for Phase 2 and switch to `embed.FS` in Phase 5.

**Warning signs:** Build error `pattern frontend/dist: no matching files found`.

### Pitfall 4: sqlc and `json_agg` return type

**What goes wrong:** sqlc generates `interface{}` or `json.RawMessage` for `json_agg(...)` columns. Handler tries to range over it directly and panics.

**Why it happens:** sqlc cannot infer a Go struct from a PostgreSQL aggregate expression.

**How to avoid:** Accept `json.RawMessage` from sqlc, then `json.Unmarshal` into `[]AuthorRef{Name, Slug}` in the handler before building the response.

**Warning signs:** Compile error "cannot range over interface{}" or runtime panic on type assertion.

### Pitfall 5: Chi route ordering for `/api/books/currently-reading` vs `/api/books/{slug}`

**What goes wrong:** `GET /api/books/currently-reading` matches `{slug}` pattern and returns a 404 from the DB (no book with slug "currently-reading").

**Why it happens:** Chi resolves routes in registration order when patterns overlap.

**How to avoid:** Register `/api/books/currently-reading` before `/api/books/{slug}` in the router.

**Warning signs:** `/api/books/currently-reading` returns a 404 or a DB "not found" error.

### Pitfall 6: CORS only in dev, not prod

**What goes wrong:** CORS middleware applied globally in production. In production, Go serves the React build from the same origin — CORS headers are unnecessary and add overhead.

**Why it happens:** Dev and prod differ in topology: dev has two processes (`localhost:5173` + `localhost:8081`); prod has one (`go:embed`).

**How to avoid:** Gate CORS middleware on an env var (e.g., `APP_ENV=development`) or `GO_ENV`. Only apply `cors.Handler(...)` when in dev mode.

---

## Code Examples

Verified patterns from the existing codebase and Go standard library.

### Error Response Helper

```go
// internal/api/public.go
func writeError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
```

### Paginated Envelope Response

```go
type PaginatedResponse struct {
    Items      any     `json:"items"`
    NextCursor *string `json:"next_cursor"` // nil when no more pages
    HasMore    bool    `json:"has_more"`
}

func writePaginated(w http.ResponseWriter, items any, nextCursor *string) {
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(PaginatedResponse{
        Items:      items,
        NextCursor: nextCursor,
        HasMore:    nextCursor != nil,
    })
}
```

### Cursor Encode/Decode

```go
import (
    "encoding/base64"
    "fmt"
    "strconv"
    "strings"
    "time"
    "github.com/jackc/pgx/v5/pgtype"
)

func encodeCursor(readAt pgtype.Timestamptz, id int64) string {
    t := readAt.Time.UTC().Format(time.RFC3339Nano)
    raw := fmt.Sprintf("%s_%d", t, id)
    return base64.URLEncoding.EncodeToString([]byte(raw))
}

func decodeCursor(s string) (time.Time, int64, error) {
    b, err := base64.URLEncoding.DecodeString(s)
    if err != nil {
        return time.Time{}, 0, err
    }
    parts := strings.SplitN(string(b), "_", 2)
    if len(parts) != 2 {
        return time.Time{}, 0, fmt.Errorf("invalid cursor")
    }
    t, err := time.Parse(time.RFC3339Nano, parts[0])
    if err != nil {
        return time.Time{}, 0, err
    }
    id, err := strconv.ParseInt(parts[1], 10, 64)
    return t, id, err
}
```

### CORS Middleware (dev only)

```go
import "github.com/go-chi/cors"

if os.Getenv("APP_ENV") == "development" {
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"http://localhost:5173"},
        AllowedMethods: []string{"GET", "OPTIONS"},
        AllowedHeaders: []string{"Accept", "Content-Type"},
        MaxAge:         300,
    }))
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Offset pagination | Keyset/cursor pagination | 2020+ | Stable across inserts; O(1) vs O(N) |
| `net/http.ServeMux` | Chi v5 router | Ongoing | Named URL params, middleware composability |
| `database/sql` + hand-written mapping | sqlc + pgx/v5 | 2022+ | Type-safe generated code; no ORM overhead |
| Separate nginx for static files | `//go:embed` + `http.FileServer` | Go 1.16 (2021) | Single binary deployment |

**Deprecated/outdated:**
- `github.com/rs/cors`: Still maintained but `github.com/go-chi/cors` is the idiomatic choice for Chi projects (same API, chi-ecosystem maintained).
- `offset` pagination: Explicitly deferred in REQUIREMENTS.md as out of scope; keyset is the required approach.

---

## Open Questions

1. **Composite index for keyset pagination**
   - What we know: `idx_books_read_at` exists on `(read_at DESC NULLS LAST)`. The keyset predicate uses `(read_at, id)` together.
   - What's unclear: Whether PostgreSQL can efficiently use the single-column index for the composite predicate.
   - Recommendation: Add migration `000002_composite_pagination_index.up.sql` with `CREATE INDEX idx_books_read_at_id ON books(read_at DESC NULLS LAST, id DESC)`. Low risk, measurable improvement.

2. **OG description truncation**
   - What we know: D-08 specifies truncation to ~200 chars if needed.
   - What's unclear: Whether to truncate at word boundary or hard-cut.
   - Recommendation: Word-boundary truncation (find last space before 200 chars, append `…`). Simple stdlib implementation, no library needed.

3. **`frontend/dist` placeholder strategy for Phase 2**
   - What we know: `//go:embed` fails if directory is empty. Phase 5 populates it.
   - What's unclear: Whether to use `os.DirFS` (Phase 2) and switch to `embed.FS` (Phase 5) or use `frontend/dist/.gitkeep` now.
   - Recommendation: Create `frontend/dist/index.html` with a minimal HTML placeholder and `frontend/dist/.gitkeep`. Use `embed.FS` from the start. The SPA handler will serve the placeholder `index.html` for all non-API routes in Phase 2 — this is fine and exercises the full code path early.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go toolchain | All compilation | Yes | 1.25.0 | — |
| PostgreSQL | All DB queries | Yes (Docker) | via docker-compose | — |
| sqlc CLI | Code generation | Yes (per Makefile) | v1.30.0 | — |
| `go-chi/cors` | CORS middleware | No (not in go.mod) | — | `go get github.com/go-chi/cors` |
| `data/covers/` directory | Cover file serving | Yes | — | — |
| `frontend/dist/` directory | go:embed, SPA | No (not yet built) | — | Create placeholder directory |

**Missing dependencies with no fallback:**
- `github.com/go-chi/cors` must be added with `go get`; no existing alternative in go.mod

**Missing dependencies with fallback:**
- `frontend/dist/` directory: create with a placeholder `index.html` — Phase 5 replaces it with the React build

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` package |
| Config file | none (Go test discovery is automatic) |
| Quick run command | `go test ./internal/api/... -v -run TestPublic` |
| Full suite command | `go test ./... -count=1` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| API-01 | `GET /api/books` returns paginated envelope, second page uses cursor | unit (httptest) | `go test ./internal/api/... -run TestGetBooks` | No — Wave 0 |
| API-01 | Cursor encode/decode roundtrip | unit | `go test ./internal/api/... -run TestCursor` | No — Wave 0 |
| API-02 | `GET /api/books/currently-reading` returns currently-reading shelf | unit (httptest) | `go test ./internal/api/... -run TestGetCurrentlyReading` | No — Wave 0 |
| API-03 | `GET /api/books/:slug` returns full detail with inline authors/genres | unit (httptest) | `go test ./internal/api/... -run TestGetBookBySlug` | No — Wave 0 |
| API-04 | `GET /api/authors` returns array with book_count | unit (httptest) | `go test ./internal/api/... -run TestGetAuthors` | No — Wave 0 |
| API-05 | `GET /api/authors/:slug` returns author + paginated books | unit (httptest) | `go test ./internal/api/... -run TestGetAuthorBySlug` | No — Wave 0 |
| API-06 | `GET /api/genres` returns array with book_count | unit (httptest) | `go test ./internal/api/... -run TestGetGenres` | No — Wave 0 |
| API-07 | `GET /api/genres/:slug` returns genre + paginated books | unit (httptest) | `go test ./internal/api/... -run TestGetGenreBySlug` | No — Wave 0 |
| API-08 | `GET /api/years` returns array of year+count | unit (httptest) | `go test ./internal/api/... -run TestGetYears` | No — Wave 0 |
| API-09 | `/covers/x.jpg` sets `Cache-Control: public, max-age=31536000, immutable` | unit (httptest) | `go test ./internal/api/... -run TestCoversCache` | No — Wave 0 |
| API-10 | `/books/:slug` response contains og:title meta tag | unit (httptest) | `go test ./internal/api/... -run TestOGInjection` | No — Wave 0 |

**Handler testing strategy:** The handlers call `*db.Queries` which requires a real Postgres connection. Two options:
1. **Interface abstraction:** Define a `BookStore` interface wrapping the `*db.Queries` methods used; inject a mock in tests. This is the clean approach but requires adding an interface layer.
2. **httptest + real DB:** Wire up a test DB (docker-compose already available). Integration tests that run against real Postgres — slower but test the full stack including SQL.

Recommendation: Use a mock interface for unit tests of handler logic (cursor parsing, response shape) and rely on manual smoke tests or separate integration test suite for SQL correctness. The existing test files in `internal/sync/` all use httptest + mocked HTTP servers (no DB) — follow the same approach.

### Wave 0 Gaps

- [ ] `internal/api/public_test.go` — covers API-01 through API-10 handler unit tests
- [ ] Interface `internal/api/store.go` — `BookStore` interface for mock injection in tests

---

## Sources

### Primary (HIGH confidence)

- Existing codebase (`cmd/server/main.go`, `internal/api/admin.go`, `internal/db/models.go`, `sql/queries/*.sql`) — direct inspection; all patterns verified from working Phase 1 code
- Go standard library: `embed`, `encoding/base64`, `text/template`, `html/template`, `net/http` — HIGH (stdlib, stable)
- go-chi/chi/v5 v5.2.5 — in go.mod, inspected router pattern directly

### Secondary (MEDIUM confidence)

- [github.com/go-chi/cors](https://github.com/go-chi/cors) — CORS package for Chi; pkg.go.dev verified; not yet in go.mod
- [Keyset Cursors for Postgres Pagination](https://blog.sequinstream.com/keyset-cursors-not-offsets-for-postgres-pagination/) — keyset predicate direction verified against PostgreSQL behavior
- [Chi fileserver example](https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go) — covers + SPA serving pattern

### Tertiary (LOW confidence)

- [Cursor pagination guide (uptrace/bun)](https://bun.uptrace.dev/guide/cursor-pagination.html) — general pattern reference; verified against known Go/PostgreSQL behavior

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all packages already in use or verified on pkg.go.dev
- Architecture: HIGH — patterns derived directly from existing codebase and stdlib docs
- sqlc query patterns: MEDIUM — `json_agg` + `json.RawMessage` pattern is standard; exact generated types verified by inspection of existing generated code
- Pitfalls: HIGH — derived from direct code inspection and known Go/PostgreSQL behavior

**Research date:** 2026-04-03
**Valid until:** 2026-07-03 (stable stack; 90-day window is conservative)
