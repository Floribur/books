# Phase 1: Data Pipeline - Research

**Researched:** 2026-04-01
**Domain:** Go project bootstrap, PostgreSQL schema migrations, Goodreads RSS parsing, Google Books enrichment, cover image download pipeline
**Confidence:** HIGH (core stack verified against pkg.go.dev and official docs; RSS field mapping is MEDIUM pending live feed inspection)

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** Goodreads RSS feed at `https://www.goodreads.com/review/list_rss/79499864?shelf=read` is confirmed live (returns XML). RSS is the primary sync path.
- **D-02:** CSV import (`POST /admin/import-csv`) is the permanent fallback — not a one-time tool.
- **D-03:** User currently has under 200 books on the "read" shelf. Pagination (per_page=200&page=N loop) is implemented as designed but not critical path immediately.
- **D-04:** `currently-reading` shelf fetched first, `read` shelf second; `read` wins on conflict.
- **D-05:** When the 6-hour sync fails (RSS unreachable): Claude's discretion — pick a sensible approach (e.g. log error, keep existing data, optionally retry with backoff).
- **D-06:** When Google Books enrichment fails for a specific book: save the book with partial data from RSS, set `metadata_source='none'`, and retry on the next sync run.
- **D-07:** Enrichment is decoupled from the RSS sync — runs as a separate goroutine (not inline).
- **D-08:** The enrichment job is triggered automatically after each sync run completes (processes all unenriched books). No separate cron needed.
- **D-09:** Google Books lookup: ISBN-13 primary, title+author fallback with confidence gate (returned author must fuzzy-match, title must contain input title); on failure → `metadata_source='none'`.
- **D-10:** OpenLibrary used as cover fallback when Google Books lacks a cover or returns a 1×1 placeholder.
- **D-11:** Cover filename for books with ISBN-13: `{isbn13}.jpg` stored at `data/covers/`.
- **D-12:** For books without ISBN-13: Claude's discretion on filename scheme (e.g. `gr-{goodreads_id}.jpg`).
- **D-13:** If a book's cover changes on re-sync (new URL, different image): overwrite the existing file silently at the same path.
- **D-14:** Cover validation: reject files < 5KB and non-decodable images.
- **Stack (STATE.md):** Chi v5, sqlc + pgx/v5, golang-migrate, Go 1.22+, PostgreSQL 16, Docker Compose for local dev.

### Claude's Discretion

- Sync retry/backoff strategy when RSS is unreachable (keep it simple for a personal site)
- Filename scheme for ISBN-less book covers
- Internal goroutine lifecycle management for the enrichment job

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope.
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| SYNC-01 | App automatically polls Goodreads RSS feed every 6 hours and updates the book database | `time.Ticker` pattern with context cancellation; goroutine lifecycle section |
| SYNC-02 | Sync handles RSS pagination (>200 books) using `?page=N` loop | gofeed pagination loop; Goodreads `per_page=200&page=N` confirmed |
| SYNC-03 | When a book appears on both `currently-reading` and `read` shelves, `read` wins | Delta diff + shelf merge pattern documented |
| SYNC-04 | Books are enriched via Google Books API (ISBN lookup primary, title+author fallback with confidence check) | Google Books API endpoint + confidence gate pattern |
| SYNC-05 | OpenLibrary API used as fallback when Google Books lacks metadata | OpenLibrary Covers API endpoint documented |
| SYNC-06 | Book cover images are downloaded and stored locally (not hotlinked) | Cover download pipeline pattern documented |
| SYNC-07 | Cover images are validated after download (rejects 1×1 pixel placeholders, files < 5KB) | `image.DecodeConfig` + file size check pattern |
| SYNC-08 | Manual sync can be triggered via `POST /admin/sync` | Chi router handler pattern; shared sync function |
| SYNC-09 | Goodreads CSV export can be imported via `POST /admin/import-csv` as permanent fallback | CSV column headers documented; `encoding/csv` + multipart handler pattern |
| DATA-01 | Books have: title, slug, author(s), cover image path, description, genres, page count, publication year, Goodreads ID, read date, ISBN-13, metadata source | Full schema column list documented |
| DATA-02 | Slugs are unique; collision resolved by appending year then author surname | gosimple/slug library + collision strategy documented |
| DATA-03 | Books track read count (for re-reads) | `read_count` column in schema |
| DATA-04 | Authors and genres are normalized (many-to-many join tables) | Schema pattern: `book_authors`, `book_genres` join tables |
| DATA-05 | Schema managed by golang-migrate with numbered SQL files | golang-migrate v4.19.1 install and file naming convention documented |
</phase_requirements>

---

## Summary

Phase 1 bootstraps the entire Go project from scratch: module, layout, Docker Compose with PostgreSQL 16, migration toolchain, Goodreads RSS sync, Google Books enrichment, and cover image storage. All stack choices are already locked by STATE.md decisions (Chi v5, sqlc + pgx/v5, golang-migrate). Research focused on verifying current versions, confirming API behaviours, and surfacing patterns and pitfalls the planner needs.

The most important gap that CANNOT be resolved by research alone is the exact XML namespace and field paths in the live Goodreads RSS feed. Plan 1.2 correctly requires a pre-task that fetches the live feed, prints all `item.Extensions` keys, and documents the paths before writing the field-mapping struct. The canonical field names from the Node.js ecosystem (e.g., `book_id`, `author_name`, `user_read_at`) are known, but the exact Go `item.Extensions` map key hierarchy (`item.Extensions["namespace"]["field"][0].Value`) must be confirmed against the live feed.

Google Books API is free with a default quota of ~1,000 requests/day. For a personal library of under 200 books, enrichment on initial load plus incremental updates will remain well within quota. Rate limiting (1 req/sec pause) should be implemented defensively. The OpenLibrary Covers API is unauthenticated and rate-limited to 100 requests per 5 minutes per IP for non-OLID key types — cover fallback should use ISBN-based URLs with appropriate throttling.

**Primary recommendation:** Follow the locked stack exactly. The only non-obvious implementation decision is the enrichment goroutine lifecycle: use a `context.Context` passed from main with a `sync.WaitGroup` so the enrichment worker drains cleanly on OS signal before the process exits.

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `github.com/go-chi/chi/v5` | v5.2.5 | HTTP router | Locked in STATE.md; idiomatic net/http, zero dependencies |
| `github.com/jackc/pgx/v5` | v5.9.1 | PostgreSQL driver + toolkit | Locked in STATE.md; lib/pq is maintenance-only; pgx/v5 is current standard |
| `github.com/jackc/pgx/v5/pgxpool` | (same module) | Connection pooling | Bundled with pgx/v5; concurrency-safe pool for sqlc queries |
| `github.com/sqlc-dev/sqlc` (CLI tool) | v1.30.0 | Type-safe query code generation | Locked in STATE.md; compile-time safety, no ORM overhead |
| `github.com/golang-migrate/migrate/v4` | v4.19.1 | Schema migrations CLI + library | Locked in STATE.md; most adopted, simple mental model |
| `github.com/mmcdole/gofeed` | v1.3.0 | RSS/Atom/JSON feed parser | De-facto standard Go feed parser; supports `item.Extensions` for custom namespaces |
| `github.com/gosimple/slug` | latest | URL slug generation with Unicode transliteration | Only well-maintained slug library in the Go ecosystem with multi-language support |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `golang.org/x/text` | (stdlib companion) | Unicode normalization | If gosimple/slug needs supplement for edge-case transliteration |
| `encoding/csv` | stdlib | CSV parsing | Goodreads CSV import — no third-party needed |
| `encoding/json` | stdlib | Google Books API response parsing | Standard HTTP + JSON; no SDK needed |
| `image` / `image/jpeg` | stdlib | Cover image validation | Decode and check dimensions to reject 1×1 placeholders |
| `net/http` | stdlib | HTTP client for API calls and cover downloads | Standard; wrap with timeout via `http.Client{Timeout: 30s}` |
| `time` | stdlib | 6-hour sync ticker | `time.NewTicker(6 * time.Hour)` in goroutine with context cancellation |
| `context` | stdlib | Goroutine lifecycle and cancellation | Pass from `main()` to all long-lived goroutines |
| `sync` | stdlib | WaitGroup for graceful shutdown | Ensure enrichment goroutine drains before exit |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| `gosimple/slug` | hand-rolled kebab | gosimple handles Unicode, accents, edge cases; hand-rolled misses ß→ss, etc. |
| `encoding/csv` (stdlib) | `github.com/gocarina/gocsv` | gocsv adds struct tags but for Goodreads import with known columns, stdlib is sufficient |
| `time.Ticker` (stdlib) | `github.com/robfig/cron/v3` | cron/v3 is only needed for calendar-based schedules; a fixed 6h interval is simpler with Ticker |
| Raw `net/http` client | `google.golang.org/api/books/v1` | Official SDK adds 10MB+ of google-api deps; plain HTTP to `googleapis.com/books/v1/volumes` is simpler and sufficient |

### Installation

**CLI tools (install once per machine):**
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

**Go module dependencies (add to go.mod):**
```bash
go get github.com/go-chi/chi/v5@latest
go get github.com/jackc/pgx/v5@latest
go get github.com/mmcdole/gofeed@latest
go get github.com/gosimple/slug@latest
```

**Version verification (run before writing Standard Stack table):**
```bash
go list -m github.com/go-chi/chi/v5
go list -m github.com/jackc/pgx/v5
go list -m github.com/mmcdole/gofeed
go list -m github.com/gosimple/slug
```

---

## Architecture Patterns

### Recommended Project Structure

```
flos-library/
├── cmd/
│   └── server/
│       └── main.go           # Entry point: wire deps, start server + scheduler
├── internal/
│   ├── db/                   # sqlc-generated code (queries.sql.go, models.go, db.go)
│   ├── sync/
│   │   ├── rss.go            # Goodreads RSS fetcher + paginator
│   │   ├── csv.go            # Goodreads CSV parser
│   │   ├── enricher.go       # Google Books + OpenLibrary enrichment goroutine
│   │   └── covers.go         # Cover download + validation
│   ├── scheduler/
│   │   └── scheduler.go      # time.Ticker loop, manual trigger channel
│   └── api/
│       └── admin.go          # POST /admin/sync, POST /admin/import-csv handlers
├── migrations/
│   ├── 000001_initial_schema.up.sql
│   └── 000001_initial_schema.down.sql
├── data/
│   └── covers/               # Downloaded cover images (gitignored)
├── sql/
│   ├── schema.sql            # Source of truth for sqlc schema
│   └── queries/
│       └── books.sql         # Annotated SQL queries for sqlc codegen
├── sqlc.yaml
├── docker-compose.yml
├── Makefile
└── go.mod
```

### Pattern 1: sqlc Configuration for pgx/v5

**What:** sqlc reads SQL files and generates type-safe Go code using pgx/v5 types directly (no `database/sql` interface).
**When to use:** All DB queries in this project.

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "sql/schema.sql"
    queries: "sql/queries/"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_pointers_for_null_types: true
```

### Pattern 2: pgxpool Connection Setup

**What:** Create a single connection pool in `main.go`, pass `*pgxpool.Pool` to sqlc `Queries`.
**When to use:** Application startup.

```go
// cmd/server/main.go
pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
if err != nil {
    log.Fatal(err)
}
defer pool.Close()

queries := db.New(pool)
```

### Pattern 3: golang-migrate File Naming

**What:** Sequential integer prefix (not timestamp) for a single-developer project avoids verbose filenames.
**Convention:** `{NNN}_{description}.up.sql` / `{NNN}_{description}.down.sql`

```
migrations/
├── 000001_initial_schema.up.sql
├── 000001_initial_schema.down.sql
├── 000002_add_search_vector.up.sql   # stub for V2 full-text search
└── 000002_add_search_vector.down.sql
```

**CLI usage:**
```bash
migrate -path ./migrations -database "$DATABASE_URL" up
migrate -path ./migrations -database "$DATABASE_URL" down 1
```

### Pattern 4: gofeed RSS Parsing with Extensions

**What:** Goodreads RSS uses a custom namespace (not Dublin Core, not iTunes). Fields appear in `item.Extensions["namespace"]["fieldname"][0].Value`.
**Critical:** The exact namespace key must be confirmed by the Plan 1.2 pre-task (print all extension keys from live feed).

Known field names (from Node.js ecosystem cross-reference — MEDIUM confidence, verify against live feed):
```
book_id, author_name, isbn13, book_image_url, book_large_image_url,
user_read_at, user_date_added, user_shelves, user_rating, book_description, book_published
```

```go
// After confirming namespace key from pre-task inspection:
func extractGoodreadsField(item *gofeed.Item, ns, key string) string {
    if ext, ok := item.Extensions[ns]; ok {
        if vals, ok := ext[key]; ok && len(vals) > 0 {
            return vals[0].Value
        }
    }
    return ""
}
```

### Pattern 5: Enrichment Goroutine Lifecycle

**What:** Decoupled enrichment worker that processes all unenriched books after each sync completes.
**D-07 and D-08 locked this pattern.**

```go
// internal/sync/enricher.go
func RunEnricher(ctx context.Context, wg *sync.WaitGroup, queries *db.Queries, trigger <-chan struct{}) {
    defer wg.Done()
    for {
        select {
        case <-trigger:
            processUnenrichedBooks(ctx, queries)
        case <-ctx.Done():
            return
        }
    }
}
```

After each RSS sync completes, send to `trigger` channel (non-blocking: use `select { case trigger <- struct{}{}: default: }` to avoid blocking if already queued).

### Pattern 6: 6-Hour Scheduler with Manual Trigger

```go
// internal/scheduler/scheduler.go
func Start(ctx context.Context, wg *sync.WaitGroup, syncFn func(ctx context.Context)) {
    defer wg.Done()
    ticker := time.NewTicker(6 * time.Hour)
    defer ticker.Stop()

    // Run immediately on startup
    syncFn(ctx)

    for {
        select {
        case <-ticker.C:
            syncFn(ctx)
        case <-ctx.Done():
            return
        }
    }
}
```

Manual trigger via `POST /admin/sync` calls the same `syncFn` directly in a goroutine (fire-and-forget with its own timeout context).

### Pattern 7: Cover Validation

```go
// internal/sync/covers.go — D-14: reject < 5KB and non-decodable
func validateCover(data []byte) error {
    if len(data) < 5*1024 {
        return fmt.Errorf("cover too small: %d bytes", len(data))
    }
    _, _, err := image.Decode(bytes.NewReader(data))
    if err != nil {
        return fmt.Errorf("cover not decodable: %w", err)
    }
    return nil
}
```

Import `_ "image/jpeg"` and `_ "image/png"` for format registration (side-effect imports).

### Pattern 8: Graceful Shutdown

```go
// cmd/server/main.go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer cancel()

var wg sync.WaitGroup
wg.Add(2) // scheduler goroutine + enricher goroutine
go scheduler.Start(ctx, &wg, syncFn)
go sync.RunEnricher(ctx, &wg, queries, enrichTrigger)

<-ctx.Done()
wg.Wait() // drain before exit
```

### Anti-Patterns to Avoid

- **Inline enrichment:** Running Google Books API calls synchronously during the RSS upsert loop. Enrichment takes 1-3 seconds per book; 200 books = 3-10 minutes blocking the sync. D-07 prevents this.
- **Goroutine leak:** Launching `go syncFn()` on every manual trigger without a semaphore. If multiple admin triggers fire rapidly, they stack up. Use a `sync.Mutex` or channel semaphore to allow only one sync at a time.
- **Orphaned ticker:** Calling `time.NewTicker` without `defer ticker.Stop()`. Tickers are not garbage collected when their goroutine exits without `Stop()`.
- **Unbounded HTTP client:** Using `http.DefaultClient` (no timeout) for cover downloads. A hung cover URL will block the goroutine indefinitely. Always use `&http.Client{Timeout: 30 * time.Second}`.
- **sqlc without schema consistency:** Running `sqlc generate` against a different schema than what golang-migrate has applied. Keep `sql/schema.sql` as the single source of truth, and ensure it matches the latest migration state.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| RSS parsing | Custom XML decoder | `gofeed` | Handles encoding quirks, pagination detection, extension namespaces, malformed feeds |
| URL slug generation | `strings.ToLower + regexp` | `gosimple/slug` | Misses Unicode: "Ö" → "o", "ñ" → "n", "ß" → "ss", ligatures, CJK |
| DB migrations | Schema-per-deploy script | `golang-migrate` | Tracks applied versions in `schema_migrations` table; enables rollback |
| Type-safe queries | Manual `pgx.Rows` scanning | `sqlc` | Eliminates scan-order bugs, nil pointer panics from nullable columns |
| Connection pooling | Single `*pgx.Conn` | `pgxpool.Pool` | Single connection serializes all concurrent DB ops; pool handles goroutine-safe access |
| Image validation | Checking file extension or Content-Type header | `image.Decode` (stdlib) | Content-Type can lie; extension is not content; only decode confirms valid image data |

**Key insight:** The Goodreads RSS extension namespace structure is the trickiest custom parsing in this phase. Do not attempt to parse the raw XML with `encoding/xml` when gofeed already handles the extension map cleanly.

---

## Common Pitfalls

### Pitfall 1: Goodreads RSS Extension Namespace Uncertainty

**What goes wrong:** Code assumes extension key is `"goodreads"` or `"gr"`, but the actual namespace prefix in the live feed may differ or be absent (fields could be in Dublin Core `dc:` namespace, or use a Goodreads-specific prefix).
**Why it happens:** The Goodreads RSS schema is undocumented. Node.js parsers access fields by element name without namespace, hiding this detail.
**How to avoid:** Execute Plan 1.2 pre-task first — fetch the live feed, unmarshal with gofeed, log `spew.Dump(item.Extensions)` or JSON-marshal the extensions map. Document all keys before writing any field mapping code.
**Warning signs:** All extension fields return empty string; gofeed parses the feed without error but no custom data is extracted.

### Pitfall 2: Goodreads CSV ISBN Quoting

**What goes wrong:** Goodreads exports ISBNs in the format `="0060590297"` (Excel formula quoting to prevent numeric truncation). `encoding/csv` returns this literally as `="0060590297"`.
**Why it happens:** Goodreads exports CSV for Excel compatibility.
**How to avoid:** Strip the `="` prefix and `"` suffix when reading ISBN/ISBN13 columns: `strings.Trim(strings.TrimPrefix(val, "="), "\"")`.
**Warning signs:** ISBN13 values in DB contain `="` prefix; Google Books ISBN lookup returns zero results.

### Pitfall 3: Google Books 1x1 Placeholder Cover

**What goes wrong:** Google Books returns a valid HTTP 200 response for a cover URL, but the image is a 1×1 pixel placeholder (`zoom=1` on books without covers). Validation by file size alone may not catch all cases (some placeholders exceed 5KB).
**Why it happens:** Google Books has a default cover thumbnail for books with no cover art.
**How to avoid:** After download and size validation, decode with `image.DecodeConfig` to get dimensions. If `width == 1 && height == 1`, reject and fall back to OpenLibrary. Note: the `image.DecodeConfig` call is cheaper than full `image.Decode` — it reads only the image header.
**Warning signs:** All covers download successfully but display as a tiny dot or broken layout.

### Pitfall 4: OpenLibrary Rate Limit on Cover API

**What goes wrong:** OpenLibrary cover API is rate-limited to 100 requests per 5 minutes per IP for non-OLID lookups (ISBN-based). On initial enrichment of 200 books, hitting all fallbacks will breach this limit, resulting in 429 responses.
**Why it happens:** OpenLibrary protects against bulk scraping.
**How to avoid:** Add a `time.Sleep(500ms)` between cover fallback requests. Since enrichment is decoupled and async (D-07), this does not block any user-facing operation. Log 429s and retry on the next enrichment pass (D-06 already prescribes retry on failure).
**Warning signs:** Cover fallback returns HTTP 429 after ~20-30 rapid requests; log shows many `metadata_source='none'` entries despite OpenLibrary having the covers.

### Pitfall 5: golang-migrate Dirty State

**What goes wrong:** A migration fails partway through and leaves the `schema_migrations` table in a `dirty=true` state. Subsequent `migrate up` commands refuse to run.
**Why it happens:** A SQL syntax error or constraint violation in the `.up.sql` file; the migration ran partially.
**How to avoid:** Test each migration file against a clean PostgreSQL instance before committing. Use transactions in migration SQL where possible (`BEGIN; ... COMMIT;`). Recovery: `migrate force VERSION` to mark the version as clean, then fix and re-run.
**Warning signs:** `migrate up` prints `Dirty database version X. Fix and force version.`

### Pitfall 6: sqlc Generate Against Stale Schema

**What goes wrong:** `sqlc generate` generates code based on `sql/schema.sql`, but the running database has a different schema (a migration was applied but `schema.sql` wasn't updated, or vice versa). Queries work at compile time but fail at runtime with column-not-found errors.
**Why it happens:** `schema.sql` drifts from the migration files.
**How to avoid:** Keep `sql/schema.sql` as the authoritative cumulative schema that matches the latest migration. Add a Makefile target `make sqlc` that runs `sqlc generate`. Run it after every migration file change.
**Warning signs:** sqlc-generated code compiles but runtime queries fail with `column "X" does not exist`.

### Pitfall 7: Enrichment Goroutine Double-Trigger

**What goes wrong:** The RSS sync completes while a previous enrichment run is still in progress. A second enrichment goroutine launches and both try to update the same book records concurrently, causing either duplicate work or DB constraint violations.
**Why it happens:** The trigger channel is unbuffered, or the enrichment function is called directly instead of going through the single-goroutine worker.
**How to avoid:** Use the single enricher goroutine pattern (Pattern 5). The `trigger` channel should be buffered with capacity 1. A non-blocking send (`select { case trigger <- struct{}{}: default: }`) means: "if not already triggered, queue one more run." The enricher processes the trigger serially — only one run at a time.
**Warning signs:** DB logs show concurrent upserts on the same `goodreads_id`; enrichment takes longer than expected.

---

## Code Examples

Verified patterns from official sources and confirmed library APIs:

### sqlc.yaml for pgx/v5 (verified from docs.sqlc.dev)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "sql/schema.sql"
    queries: "sql/queries/"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_pointers_for_null_types: true
```

### pgxpool Connection (verified from pkg.go.dev/github.com/jackc/pgx/v5/pgxpool)

```go
pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
if err != nil {
    log.Fatalf("unable to connect to database: %v", err)
}
defer pool.Close()
```

### gofeed RSS Fetch (verified from pkg.go.dev/github.com/mmcdole/gofeed)

```go
fp := gofeed.NewParser()
feed, err := fp.ParseURL("https://www.goodreads.com/review/list_rss/79499864?shelf=read&per_page=200&page=1")
if err != nil {
    return fmt.Errorf("fetch RSS: %w", err)
}
for _, item := range feed.Items {
    // item.Extensions["NAMESPACE"]["FIELD"][0].Value
    // NAMESPACE = confirmed by pre-task inspection of live feed
}
```

### Google Books ISBN Lookup (verified from developers.google.com/books/docs/v1/using)

```go
// Primary path: ISBN-13
url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=isbn:%s&key=%s", isbn13, apiKey)

// Fallback: title+author
url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=intitle:%s+inauthor:%s&key=%s",
    url.QueryEscape(title), url.QueryEscape(author), apiKey)
```

### OpenLibrary Cover URL (verified from openlibrary.org/dev/docs/api/covers)

```go
// Large cover by ISBN-13 (rate-limited: 100 req / 5 min per IP)
coverURL := fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-L.jpg", isbn13)
```

### Goodreads CSV ISBN Unquoting

```go
func unquoteISBN(raw string) string {
    // Goodreads exports: ="9780385472579" -> strip =", trailing "
    s := strings.TrimPrefix(raw, "=\"")
    s = strings.TrimSuffix(s, "\"")
    return strings.TrimSpace(s)
}
```

### Goodreads CSV Column Headers (confirmed from live export sample)

Key columns for import (exact header names from Goodreads export):
```
Book Id, Title, Author, ISBN, ISBN13, Number of Pages,
Year Published, Original Publication Year, Date Read, Date Added,
Bookshelves, Exclusive Shelf, Read Count
```

### golang-migrate Makefile Targets

```makefile
DB_URL ?= postgres://postgres:postgres@localhost:5432/floslib?sslmode=disable

migrate:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `lib/pq` PostgreSQL driver | `pgx/v5` | ~2022 (lib/pq maintenance mode) | pgx has better type support, `pgxpool` built-in, `pgconn` for copy protocol |
| `database/sql` interface | `pgx/v5` native types directly via sqlc | ~2023 (sqlc pgx/v5 support v1.18) | No interface overhead; pgx types (e.g., `pgtype.Timestamptz`) instead of nullable wrappers |
| Offset-based pagination | Cursor-based pagination (locked: `read_at` + `id`) | Ongoing | Stable across concurrent inserts; Phase 2 concern but schema must support it |
| ORM (GORM, ent) | sqlc type-safe code generation | 2021-present | No reflection at runtime; SQL stays readable; schema changes surface as compile errors |
| Manual migration scripts | golang-migrate | ~2019-present | Version tracking in `schema_migrations` table; rollback support |

**Deprecated/outdated:**
- `lib/pq`: Still works but maintenance-only since 2021. Do not use for new projects.
- Goodreads official API: Deprecated 2020. Use RSS + Google Books as locked in STATE.md.
- `github.com/golang-migrate/migrate/v3`: Use v4 (v3 is abandoned).
- sqlc v1 YAML config (`version: "1"`): Use `version: "2"` which supports multiple SQL blocks and richer config.

---

## Open Questions

1. **Goodreads RSS Extension Namespace Key**
   - What we know: Fields like `book_id`, `author_name`, `user_read_at` exist in the feed per Node.js implementations
   - What's unclear: The exact string key for the namespace in gofeed's `item.Extensions` map (could be `"goodreads_book"`, `"gr"`, or the namespace URI as key)
   - Recommendation: Plan 1.2 pre-task MUST run first. Print `item.Extensions` keys from the live feed before writing any field mapping struct. This is a blocking dependency for Plan 1.2 Task 3.

2. **isbn13 Availability in RSS vs. Google Books**
   - What we know: Goodreads RSS has an `isbn13` field; some books (especially older or obscure ones) may have no ISBN-13
   - What's unclear: What fraction of Florian's ~200 books lack ISBN-13 — affects how many will need title+author fallback for enrichment
   - Recommendation: Accept this as runtime data; the confidence gate (D-09) handles it correctly. No pre-research needed.

3. **Cover Filename for ISBN-less Books (Claude's Discretion)**
   - What we know: D-12 delegates this to implementation
   - Recommendation: Use `gr-{goodreads_id}.jpg`. Goodreads IDs are stable numeric identifiers that appear in the RSS feed. Simple, deterministic, human-readable. No collision risk with ISBN-based filenames.

4. **Sync Retry/Backoff for RSS Failure (Claude's Discretion)**
   - What we know: D-05 delegates this to implementation; "keep it simple"
   - Recommendation: Log the error, keep existing DB data (no deletion), retry at next scheduled tick (6h). No exponential backoff needed for a personal site — the 6h cadence is already the retry.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | All Go code | ✓ | go1.22.4 windows/amd64 | — |
| Docker | PostgreSQL 16 + adminer dev containers | ✓ | 26.1.4 | — |
| Docker Compose | `make dev` | ✓ | v2.27.1 | — |
| golang-migrate CLI | `make migrate` | ✗ | — | `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest` |
| sqlc CLI | `make sqlc` (codegen) | ✗ | — | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` |
| psql client | DB inspection / debugging | ✗ | — | Use adminer (bundled in Docker Compose) |
| make | Makefile targets | ✗ (no output) | — | Run commands directly; or install via `winget install GnuWin32.Make` |
| Node.js | Not needed for Phase 1 | ✓ | v20.14.0 | N/A |

**Missing dependencies with no fallback:**
- None — all missing tools have install paths or viable alternatives.

**Missing dependencies with fallback:**
- `golang-migrate` CLI: Install via `go install` with postgres tag (Wave 0 task)
- `sqlc` CLI: Install via `go install` (Wave 0 task)
- `psql`: adminer available via Docker Compose on `localhost:8080`
- `make`: either install or run Makefile commands directly in shell

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` package (no test framework needed) |
| Config file | none (go test discovers `*_test.go` files automatically) |
| Quick run command | `go test ./internal/sync/... -run TestUnit -short` |
| Full suite command | `go test ./... -timeout 60s` |

### Phase Requirements -> Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| SYNC-02 | Pagination loop stops when response < 200 items | unit | `go test ./internal/sync/... -run TestRSSPagination -short` | ❌ Wave 0 |
| SYNC-03 | `read` shelf wins when same book on both shelves | unit | `go test ./internal/sync/... -run TestShelfMerge -short` | ❌ Wave 0 |
| SYNC-04 | ISBN lookup primary, title+author fallback with confidence gate | unit | `go test ./internal/sync/... -run TestEnrichmentConfidenceGate -short` | ❌ Wave 0 |
| SYNC-07 | Cover validation rejects files < 5KB and non-decodable | unit | `go test ./internal/sync/... -run TestCoverValidation -short` | ❌ Wave 0 |
| SYNC-09 | CSV import unquotes ISBNs and maps column headers correctly | unit | `go test ./internal/sync/... -run TestCSVImport -short` | ❌ Wave 0 |
| DATA-02 | Slug collision resolved by appending year then author surname | unit | `go test ./internal/sync/... -run TestSlugCollision -short` | ❌ Wave 0 |
| DATA-05 | Migration applies cleanly against PostgreSQL 16 | integration | `go test ./internal/db/... -run TestMigrations` (requires Docker) | ❌ Wave 0 |
| SYNC-01 | 6-hour ticker fires sync function | unit | `go test ./internal/scheduler/... -run TestScheduler -short` | ❌ Wave 0 |

### Sampling Rate

- **Per task commit:** `go test ./internal/... -short -timeout 30s`
- **Per wave merge:** `go test ./... -timeout 60s`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps

- [ ] `internal/sync/rss_test.go` — covers SYNC-02, SYNC-03
- [ ] `internal/sync/enricher_test.go` — covers SYNC-04, confidence gate logic
- [ ] `internal/sync/covers_test.go` — covers SYNC-07 (needs test fixture: tiny JPEG, 1x1 JPEG, valid JPEG)
- [ ] `internal/sync/csv_test.go` — covers SYNC-09, Goodreads ISBN unquoting
- [ ] `internal/sync/slug_test.go` — covers DATA-02, collision strategy
- [ ] `internal/scheduler/scheduler_test.go` — covers SYNC-01 (mock ticker via dependency injection)
- [ ] `internal/db/migrations_test.go` — covers DATA-05 (integration, requires running PostgreSQL)

---

## Sources

### Primary (HIGH confidence)

- `pkg.go.dev/github.com/mmcdole/gofeed` — version v1.3.0 confirmed, Extensions structure verified
- `pkg.go.dev/github.com/jackc/pgx/v5` — version v5.9.1 confirmed, pgxpool.New API verified
- `pkg.go.dev/github.com/sqlc-dev/sqlc` — version v1.30.0 confirmed
- `pkg.go.dev/github.com/go-chi/chi/v5` — version v5.2.5 confirmed
- `github.com/golang-migrate/migrate/releases` — v4.19.1 confirmed; CLI install command verified
- `developers.google.com/books/docs/v1/using` — ISBN query URL format `?q=isbn:{isbn}` verified
- `openlibrary.org/dev/docs/api/covers` — Cover URL format and rate limit (100 req/5min) verified
- `gist.github.com/tmcw/f077b2f174a0194f62b94bec4e88f4d0` — Goodreads CSV column headers confirmed (31 columns)
- `pkg.go.dev/github.com/gosimple/slug` — library confirmed as standard choice

### Secondary (MEDIUM confidence)

- `leohuynh.dev/blog/crawling-goodreads-books-data` — Goodreads RSS field names (book_id, author_name, etc.) verified via Node.js cross-reference; namespace key for Go's gofeed must be confirmed against live feed
- `docs.sqlc.dev/en/stable/guides/using-go-and-pgx.html` — sqlc.yaml pgx/v5 configuration pattern (403 on direct fetch; confirmed from search result snippet + multiple version docs)
- `betterstack.com/community/guides/scaling-go/golang-migrate` — golang-migrate usage patterns

### Tertiary (LOW confidence)

- Google Books API daily quota (~1,000 req/day free tier) — from community reports; not officially documented in current Google Cloud docs page. Sufficient for personal use regardless.

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all versions verified against pkg.go.dev and official releases
- Architecture: HIGH — patterns follow official docs and locked decisions from CONTEXT.md
- Pitfalls: HIGH (Goodreads CSV quoting, goroutine leaks) / MEDIUM (RSS namespace — must be confirmed by pre-task)
- Validation architecture: HIGH — Go stdlib testing, no framework needed

**Research date:** 2026-04-01
**Valid until:** 2026-07-01 (stable ecosystem; gofeed, chi, pgx are slow-moving)
