---
phase: 01-data-pipeline
plan: 01
subsystem: database
tags: [go, postgresql, sqlc, pgx, golang-migrate, docker, chi]

# Dependency graph
requires: []
provides:
  - Go module flos-library with Chi, pgx/v5, gofeed, slug dependencies
  - PostgreSQL 16 via Docker Compose on port 5432
  - Database schema: books, authors, genres, book_authors, book_genres tables
  - golang-migrate applied migration 000001 (version 1)
  - sqlc-generated type-safe Go DB layer in internal/db/
  - Makefile targets: migrate, migrate-down, sqlc, build, dev
affects: [01-02, 01-03, 01-04, 01-05, all-plans]

# Tech tracking
tech-stack:
  added:
    - github.com/go-chi/chi/v5 v5.2.5
    - github.com/jackc/pgx/v5 v5.9.1
    - github.com/mmcdole/gofeed v1.3.0
    - github.com/gosimple/slug v1.15.0
    - golang-migrate CLI (latest, postgres tag)
    - sqlc CLI v1.30.0
    - Docker Compose with postgres:16-alpine
  patterns:
    - sqlc generates from sql/schema.sql (not migration files directly)
    - Migration naming: 000001_name.up.sql / 000001_name.down.sql
    - pgxpool for connection pooling in main.go
    - signal.NotifyContext for graceful shutdown

key-files:
  created:
    - go.mod
    - go.sum
    - cmd/server/main.go
    - docker-compose.yml
    - sqlc.yaml
    - Makefile
    - migrations/000001_initial_schema.up.sql
    - migrations/000001_initial_schema.down.sql
    - sql/schema.sql
    - sql/queries/books.sql
    - internal/db/models.go
    - internal/db/db.go
    - internal/db/books.sql.go
    - .env.example
    - .gitignore
  modified: []

key-decisions:
  - "sql/schema.sql is the sqlc source of truth — kept in manual sync with migrations (not auto-derived)"
  - "books table includes search_vector TSVECTOR stub for future full-text search (DATA-05 future-proofing)"
  - "emit_pointers_for_null_types: true in sqlc.yaml for nullable column type safety"

patterns-established:
  - "Pattern: sql/schema.sql mirrors migrations DDL without BEGIN/COMMIT wrappers for sqlc compatibility"
  - "Pattern: Makefile DB_URL uses ?= so it can be overridden by environment"
  - "Pattern: Docker Compose healthcheck on pg_isready before dependent services start"

requirements-completed: [DATA-01, DATA-02, DATA-03, DATA-04, DATA-05]

# Metrics
duration: 25min
completed: 2026-04-01
---

# Phase 1 Plan 01: Bootstrap Summary

**Go module flos-library with PostgreSQL 16 via Docker, golang-migrate schema (5 tables), and sqlc-generated pgx/v5 DB layer**

## Performance

- **Duration:** 25 min
- **Started:** 2026-04-01T21:20:31Z
- **Completed:** 2026-04-01T21:45:00Z
- **Tasks:** 2
- **Files modified:** 15 created, 0 modified

## Accomplishments
- Go module `flos-library` compiles cleanly with all 4 required dependencies (Chi, pgx/v5, gofeed, slug)
- PostgreSQL 16 running via Docker Compose with healthcheck; migration 000001 applied (version 1)
- Complete schema: books (17 columns including goodreads_id, slug UNIQUE, read_count, shelf, search_vector), authors, genres, book_authors, book_genres join tables
- sqlc generated type-safe Go code (`models.go`, `db.go`, `books.sql.go`) in `internal/db/`
- Migration rollback (`migrate-down`) and re-apply both verified working

## Task Commits

Each task was committed atomically:

1. **Task 1: Initialize Go module, project layout, and dependencies** - `dd550a2` (feat)
2. **Task 2: Docker Compose, sqlc config, migration schema, and Makefile** - `e5e6222` (feat)

**Plan metadata:** (docs commit — see below)

## Files Created/Modified
- `go.mod` - Module flos-library, go 1.25.0, all 4 library deps
- `go.sum` - Dependency checksums
- `cmd/server/main.go` - Minimal entry point with pgxpool + signal handling
- `docker-compose.yml` - PostgreSQL 16-alpine + adminer, pg_isready healthcheck
- `sqlc.yaml` - pgx/v5, emit_json_tags, emit_pointers_for_null_types
- `Makefile` - dev, migrate, migrate-down, sqlc, build targets
- `migrations/000001_initial_schema.up.sql` - Full DDL for all 5 tables + indexes
- `migrations/000001_initial_schema.down.sql` - Rollback DDL
- `sql/schema.sql` - sqlc source-of-truth (DDL without transaction wrappers)
- `sql/queries/books.sql` - GetBookBySlug, ListBooks placeholder queries
- `internal/db/models.go` - sqlc-generated Go structs for all tables
- `internal/db/db.go` - sqlc DB/Querier interfaces
- `internal/db/books.sql.go` - sqlc-generated query implementations
- `.env.example` - DATABASE_URL, GOOGLE_BOOKS_API_KEY, PORT
- `.gitignore` - data/covers/*, .env, *.exe, build artifacts

## Decisions Made
- `sql/schema.sql` is kept as sqlc source-of-truth, manually synced with migrations. sqlc reads schema.sql because it doesn't understand golang-migrate's transaction wrappers.
- `books.search_vector` (TSVECTOR) included as a stub column for future full-text search without requiring a schema migration later.
- `emit_pointers_for_null_types: true` in sqlc.yaml ensures nullable columns (e.g., `*string` for description, cover_path) are type-safe in Go.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing Critical] Added *.exe to .gitignore**
- **Found during:** Task 1 (git status revealed pre-existing server.exe binary)
- **Issue:** A compiled Windows binary `server.exe` was present in the project root and would be tracked by git without the pattern
- **Fix:** Added `*.exe` to `.gitignore` to exclude all Windows executables
- **Files modified:** `.gitignore`
- **Verification:** `git status` shows server.exe as ignored after the addition
- **Committed in:** `dd550a2` (Task 1 commit)

---

**Total deviations:** 1 auto-fixed (1 missing critical)
**Impact on plan:** Minor preventive fix. No scope creep.

## Issues Encountered
- Docker Desktop was not running when Task 2 began. Launched programmatically via PowerShell (`Start-Process`), waited for Linux engine to become available (~30s), then proceeded. No user intervention required.
- pgxpool submodule required explicit `go get github.com/jackc/pgx/v5/pgxpool` after the main pgx/v5 get — missing go.sum entry resolved by targeting the subpackage directly.

## User Setup Required
None - no external service configuration required beyond Docker Desktop (already installed).

## Next Phase Readiness
- PostgreSQL 16 running locally, schema version 1 applied — all 5 tables ready
- sqlc-generated DB layer in `internal/db/` ready for use in sync and API plans
- `DATABASE_URL` env var configured in `.env.example` — copy to `.env` before running
- Plan 1.2 (Goodreads RSS sync) can proceed immediately

---
*Phase: 01-data-pipeline*
*Completed: 2026-04-01*

## Self-Check: PASSED

- All key files exist on disk
- Both task commits (`dd550a2`, `e5e6222`) verified in git log
- `go build ./...` exits 0
- Migration version 1 confirmed applied
