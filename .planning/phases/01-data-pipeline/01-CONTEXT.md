# Phase 1: Data Pipeline - Context

**Gathered:** 2026-04-01
**Status:** Ready for planning

<domain>
## Phase Boundary

Bootstrap the Go project, stand up PostgreSQL schema with golang-migrate, implement the Goodreads RSS sync pipeline, enrich book records via Google Books (with OpenLibrary cover fallback), and download/validate cover images to local storage. The sync runs automatically every 6 hours. This phase delivers: books in the database, automatically.

</domain>

<decisions>
## Implementation Decisions

### RSS Sync
- **D-01:** Goodreads RSS feed at `https://www.goodreads.com/review/list_rss/79499864?shelf=read` is confirmed live (returns XML). RSS is the primary sync path.
- **D-02:** CSV import (`POST /admin/import-csv`) is the permanent fallback — not a one-time tool.
- **D-03:** User currently has under 200 books on the "read" shelf. Pagination (per_page=200&page=N loop) is implemented as designed but not critical path immediately.
- **D-04:** `currently-reading` shelf fetched first, `read` shelf second; `read` wins on conflict.

### Sync Error Handling
- **D-05:** When the 6-hour sync fails (RSS unreachable): Claude's discretion — pick a sensible approach (e.g. log error, keep existing data, optionally retry with backoff).
- **D-06:** When Google Books enrichment fails for a specific book: save the book with partial data from RSS, set `metadata_source='none'`, and retry on the next sync run.

### Metadata Enrichment
- **D-07:** Enrichment is decoupled from the RSS sync — runs as a separate goroutine (not inline).
- **D-08:** The enrichment job is triggered automatically after each sync run completes (processes all unenriched books). No separate cron needed.
- **D-09:** Google Books lookup: ISBN-13 primary, title+author fallback with confidence gate (returned author must fuzzy-match, title must contain input title); on failure → `metadata_source='none'`.
- **D-10:** OpenLibrary used as cover fallback when Google Books lacks a cover or returns a 1×1 placeholder.

### Cover Images
- **D-11:** Cover filename for books with ISBN-13: `{isbn13}.jpg` stored at `data/covers/`.
- **D-12:** For books without ISBN-13: Claude's discretion on filename scheme (e.g. `gr-{goodreads_id}.jpg`).
- **D-13:** If a book's cover changes on re-sync (new URL, different image): overwrite the existing file silently at the same path.
- **D-14:** Cover validation: reject files < 5KB and non-decodable images.

### Claude's Discretion
- Sync retry/backoff strategy when RSS is unreachable (keep it simple for a personal site)
- Filename scheme for ISBN-less book covers
- Internal goroutine lifecycle management for the enrichment job

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Phase Requirements
- `.planning/REQUIREMENTS.md` §Sync Pipeline (SYNC-01–09) — full sync requirements
- `.planning/REQUIREMENTS.md` §Data Model (DATA-01–05) — schema requirements
- `.planning/ROADMAP.md` §Phase 1 — plan-level task breakdown (Plans 1.1, 1.2, 1.3)

### Project Constraints
- `.planning/PROJECT.md` §Constraints — tech stack, DB, cover, and Goodreads sync constraints
- `.planning/STATE.md` §Recent Decisions — locked stack choices (Chi v5, sqlc + pgx/v5, golang-migrate, etc.)

No external ADRs or specs — requirements and decisions fully captured above.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- None — project is greenfield, no existing code

### Established Patterns
- None yet — Phase 1 establishes the foundational patterns

### Integration Points
- Phase 2 (REST API) consumes the PostgreSQL DB written here — schema decisions here lock in Phase 2's query surface
- Cover files stored at `data/covers/` will be served by Phase 2's `/covers/*` file server

</code_context>

<specifics>
## Specific Ideas

- Docker Compose with PostgreSQL 16 + adminer for local dev (from ROADMAP.md Plan 1.1)
- Makefile targets: `make dev`, `make migrate`, `make migrate-down`
- Goodreads RSS `item.Extensions` keys must be documented in Plan 1.2 pre-task before writing field mapping struct
- Google Books API key stored as `GOOGLE_BOOKS_API_KEY` env var (never in frontend)

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 01-data-pipeline*
*Context gathered: 2026-04-01*
