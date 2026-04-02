---
phase: 01-data-pipeline
plan: 02
subsystem: database
tags: [go, goodreads, rss, gofeed, gosimple-slug, csv, scheduler, chi, sqlc, pgx]

# Dependency graph
requires:
  - phase: 01-data-pipeline/01-01
    provides: Go module, PostgreSQL schema, sqlc config, migrations, Makefile

provides:
  - Goodreads RSS sync pipeline (FetchShelf, SyncRSS) using item.Custom field access
  - Shelf merge logic (read wins over currently-reading per D-04)
  - Slug generation with collision resolution (year suffix, then author surname)
  - Goodreads CSV import handler (ImportCSV, unquoteISBN)
  - 6-hour scheduler with testable ticker injection
  - Admin HTTP endpoints (POST /admin/sync, POST /admin/import-csv)
  - sqlc-generated queries: UpsertBook, GetAllGoodreadsIDs, GetUnenrichedBooks, UpsertAuthor, LinkBookAuthor, UpsertGenre, LinkBookGenre
  - Wired cmd/server/main.go with chi router, DB pool, scheduler, enrichment trigger channel

affects: [01-data-pipeline/01-03, 02-api, 03-frontend]

# Tech tracking
tech-stack:
  added: [github.com/mmcdole/gofeed v1.3.0, github.com/gosimple/slug v1.15.0]
  patterns:
    - item.Custom map for Goodreads RSS field access (NOT item.Extensions)
    - fetchShelfFromURL base-URL injection for testable pagination
    - startWithTicker injectable ticker channel for testable scheduler
    - Buffered(1) enrichment trigger channel to prevent double-trigger
    - syncMu TryLock for preventing concurrent admin syncs
    - Signal.NotifyContext for graceful shutdown with WaitGroup drain

key-files:
  created:
    - internal/sync/rss.go
    - internal/sync/slug.go
    - internal/sync/csv.go
    - internal/sync/rss_test.go
    - internal/sync/csv_test.go
    - internal/sync/slug_test.go
    - internal/scheduler/scheduler.go
    - internal/scheduler/scheduler_test.go
    - internal/api/admin.go
    - sql/queries/authors.sql
    - sql/queries/genres.sql
    - internal/db/authors.sql.go
    - internal/db/genres.sql.go
  modified:
    - sql/queries/books.sql
    - internal/db/books.sql.go
    - cmd/server/main.go

key-decisions:
  - "Goodreads RSS fields delivered via item.Custom map, not item.Extensions namespace hierarchy — confirmed from live feed inspection"
  - "fetchShelfFromURL accepts base URL parameter for test server injection (TDD pagination test)"
  - "startWithTicker accepts injectable tick channel and stop func for deterministic scheduler testing"
  - "isbn field from RSS (not isbn13) stored in isbn13 column — Goodreads RSS exposes only ISBN-10/13 in 'isbn' key"
  - "Multiple date format fallbacks for user_read_at parsing (RFC1123Z variants)"

patterns-established:
  - "Goodreads RSS access: item.Custom['book_id'], item.Custom['author_name'], etc. (not item.Extensions)"
  - "Testable internals: unexported functions with injected dependencies, public wrappers call with defaults"
  - "Graceful shutdown: signal.NotifyContext + WaitGroup + srv.Shutdown"
  - "Admin sync: TryLock mutex returns 409 Conflict if already running"

requirements-completed: [SYNC-01, SYNC-02, SYNC-03, SYNC-08, SYNC-09]

# Metrics
duration: 35min
completed: 2026-04-02
---

# Phase 01 Plan 02: RSS Sync Pipeline Summary

**Goodreads RSS sync pipeline with pagination, shelf merge, slug collision resolution, CSV fallback import, 6-hour scheduler, and admin endpoints — all wired into a chi HTTP server**

## Performance

- **Duration:** ~35 min
- **Started:** 2026-04-02T11:45:00Z
- **Completed:** 2026-04-02T12:20:00Z
- **Tasks:** 3
- **Files modified:** 17

## Accomplishments

- RSS pipeline fetches both `currently-reading` and `read` shelves, merges (read wins), and upserts books using sqlc-generated queries
- Goodreads CSV import permanently available at `POST /admin/import-csv` with Excel ISBN unquoting (`="..."` → clean string)
- Scheduler fires on startup and every 6 hours; admin trigger available at `POST /admin/sync` (returns 202, non-blocking)
- All unit tests pass: TestRSSPagination (mock httptest server, 205 items across 2 pages), TestShelfMerge, TestSlugCollision, TestGenerateSlug, TestCSVISBNUnquote, TestCSVImport, TestScheduler

## Task Commits

1. **Task 1: Wave 0 test scaffolds** - `dfa8163` (test)
2. **Task 2: RSS sync pipeline** - `2a6f042` (feat)
3. **Task 3: CSV import, scheduler, admin endpoints** - `7d87fb9` (feat)

## Files Created/Modified

- `internal/sync/rss.go` - FetchShelf, fetchShelfFromURL, parseItem, mergeShelfItems, SyncRSS
- `internal/sync/slug.go` - GenerateSlug with year-then-surname collision resolution
- `internal/sync/csv.go` - ImportCSV, unquoteISBN
- `internal/scheduler/scheduler.go` - Start, startWithTicker (testable core)
- `internal/api/admin.go` - AdminHandlers.PostSync, AdminHandlers.PostImportCSV
- `cmd/server/main.go` - Full wiring: DB pool, queries, enrichTrig channel, chi router, scheduler
- `sql/queries/books.sql` - Added UpsertBook, GetAllGoodreadsIDs, GetUnenrichedBooks
- `sql/queries/authors.sql` - UpsertAuthor, LinkBookAuthor (new file)
- `sql/queries/genres.sql` - UpsertGenre, LinkBookGenre (new file)
- `internal/db/books.sql.go` - Regenerated with new queries
- `internal/db/authors.sql.go` - New generated file
- `internal/db/genres.sql.go` - New generated file
- `internal/sync/rss_test.go` - TestRSSPagination (httptest server), TestShelfMerge
- `internal/sync/csv_test.go` - TestCSVImport, TestCSVISBNUnquote
- `internal/sync/slug_test.go` - TestSlugCollision, TestGenerateSlug
- `internal/scheduler/scheduler_test.go` - TestScheduler (channel-based mock ticker)

## Decisions Made

- `item.Custom` map is the actual Goodreads RSS field access pattern, not `item.Extensions` — confirmed from live feed inspection of `https://www.goodreads.com/review/list_rss/79499864?shelf=read`
- ISBN field in Goodreads RSS is `isbn` (not `isbn13`); stored in the `isbn13` DB column
- Scheduler uses injectable ticker channel (`startWithTicker`) to make TestScheduler deterministic without real time delays

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] item.Custom replaces item.Extensions for Goodreads RSS field access**
- **Found during:** Task 2 (MANDATORY PRE-TASK: live RSS feed inspection)
- **Issue:** Plan specified `item.Extensions["CONFIRMED_NAMESPACE"]["field"][0].Value` access pattern. Live feed inspection showed `item.Extensions` is `null`. All Goodreads-specific fields arrive in `item.Custom` as a flat `map[string]string` (e.g., `item.Custom["book_id"]`).
- **Fix:** Replaced `extractField(item, key)` using Extensions with `extractCustom(item, key)` using Custom. No namespace key needed.
- **Files modified:** internal/sync/rss.go
- **Verification:** Live RSS returns 127 items correctly parsed; TestRSSPagination passes with mock server
- **Committed in:** 2a6f042 (Task 2 commit)

**2. [Rule 1 - Bug] Multiple RFC1123Z date format fallbacks for user_read_at**
- **Found during:** Task 2 (live RSS inspection)
- **Issue:** Live feed shows `user_read_at: "Sat, 7 Feb 2026 00:00:00 +0000"` (single-digit day, no timezone abbreviation) which doesn't match standard `time.RFC1123Z`. Plan only specified `time.RFC1123Z`.
- **Fix:** Added fallback format `"Mon, 2 Jan 2006 15:04:05 -0700"` and `"Mon, 2 Jan 2006 15:04:05 +0000"` to handle Goodreads' non-standard formatting.
- **Files modified:** internal/sync/rss.go
- **Verification:** Date parsing works correctly for the live feed format
- **Committed in:** 2a6f042 (Task 2 commit)

**3. [Rule 1 - Bug] Scheduler TestScheduler race condition fixed with atomic counter and channel sync**
- **Found during:** Task 3 (TestScheduler failing)
- **Issue:** First test iteration used `time.Sleep(50ms)` to wait for initial sync, which was a race condition (goroutine hadn't called syncFn yet when checked).
- **Fix:** Replaced sleep-based timing with `chan struct{}` signaled after each syncFn call, plus `atomic.Int32` for the call counter.
- **Files modified:** internal/scheduler/scheduler_test.go
- **Verification:** TestScheduler passes consistently
- **Committed in:** 7d87fb9 (Task 3 commit)

---

**Total deviations:** 3 auto-fixed (2 Rule 1 bugs from live feed discovery, 1 Rule 1 test race)
**Impact on plan:** All auto-fixes necessary for correctness. The Extensions→Custom fix was the critical pre-task finding. No scope creep.

## Issues Encountered

- Live RSS feed returned `null` for `item.Extensions` — all fields in `item.Custom` flat map. This is a gofeed behavior difference: Goodreads uses non-standard custom elements without XML namespaces, so gofeed routes them to `item.Custom` rather than `item.Extensions`.

## User Setup Required

None - no external service configuration required for this plan. Google Books API key needed for Plan 01-03 (enrichment).

## Next Phase Readiness

- RSS sync pipeline complete and ready — Plan 01-03 (Google Books enrichment) can now call `queries.GetUnenrichedBooks()` to get books needing enrichment
- Admin endpoints wired and available at port 8081
- `POST /admin/sync` + `POST /admin/import-csv` ready for manual use
- Enrichment trigger channel `enrichTrig` available in main.go for Plan 01-03 to wire up an enrichment goroutine

---
*Phase: 01-data-pipeline*
*Completed: 2026-04-02*

## Self-Check: PASSED

- internal/sync/rss.go — FOUND
- internal/sync/slug.go — FOUND
- internal/sync/csv.go — FOUND
- internal/scheduler/scheduler.go — FOUND
- internal/api/admin.go — FOUND
- cmd/server/main.go — FOUND
- Commit dfa8163 — FOUND (test stubs)
- Commit 2a6f042 — FOUND (RSS pipeline)
- Commit 7d87fb9 — FOUND (CSV/scheduler/admin)
