---
phase: 02-go-rest-api
plan: 02
subsystem: api
tags: [go, embed, spa, open-graph, cover-serving, chi]

# Dependency graph
requires:
  - phase: 02-go-rest-api
    plan: 01
    provides: PublicHandlers, BookStore interface, GetBookDetailBySlug query, all API routes
provides:
  - frontend/embed.go: embedded dist FS via go:embed
  - Cover file server at /covers/* with immutable cache headers
  - SPA catch-all /* serving embedded index.html
  - OG meta tag injection for /books/<slug> paths (5 tags)
  - .env.example documenting all 5 env vars
affects: [frontend, phase-03, phase-05]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - go:embed in dedicated package (frontend/embed.go) adjacent to embedded directory
    - filepath.Base for OS-agnostic cover path extraction (handles both / and \ separators)
    - SPA catch-all: try real file first, then OG inject for /books/slug, else serve index.html
    - html.EscapeString for XSS-safe OG tag injection
    - Description truncated to ~200 chars at word boundary with ellipsis

key-files:
  created:
    - frontend/embed.go
    - frontend/dist/.gitkeep
    - frontend/dist/index.html
    - .env.example
  modified:
    - cmd/server/main.go (distFS setup, cover handler, SPA catch-all, buildOGTags)

key-decisions:
  - "frontend/embed.go instead of embed in main.go: Go //go:embed paths are relative to source file; cmd/server/ cannot reference ../../frontend/dist"
  - "filepath.Base for cover URL: cover_path may use OS path separators; filepath.Base extracts filename agnostically"
  - "OG injection via strings.Replace on indexHTML: read once at startup, inject per-request — no template overhead"

patterns-established:
  - "Pattern: Put go:embed in a dedicated package adjacent to the embedded directory when main package is in a subdirectory"
  - "Pattern: filepath.Base for path-to-filename extraction instead of strings.TrimPrefix with hardcoded separators"

requirements-completed:
  - API-09
  - API-10

# Metrics
duration: 30min
completed: 2026-04-03
---

# Phase 02-02: Server Topology Summary

**Production server topology: go:embed SPA, immutable cover cache, and per-book OG meta tag injection for WhatsApp/Discord link previews**

## Performance

- **Duration:** 30 min
- **Started:** 2026-04-03T15:50:00Z
- **Completed:** 2026-04-03T16:20:00Z
- **Tasks:** 3 (including human checkpoint)
- **Files modified:** 5

## Accomplishments
- `GET /covers/*` serves files from `data/covers/` with `Cache-Control: public, max-age=31536000, immutable`
- `GET /*` SPA catch-all serves embedded `frontend/dist/index.html` for all unknown paths
- `GET /books/<slug>` injects 5 OG meta tags before `</head>` using `queries.GetBookDetailBySlug`
- `.env.example` documents all 5 env vars with usage comments
- Human smoke test confirmed: covers 200 with immutable cache, SPA serves Flo's Library placeholder, OG tags visible for real book slug

## Task Commits

1. **Task 1: go:embed + cover handler + SPA catch-all + OG injection** - `3cfc3f5` (feat)
   - **Auto-fix: filepath.Base** - `d6dbdda` (fix: OS-agnostic cover path)
2. **Task 2: .env.example** - `7d4064f` (feat)
3. **Task 3: Human smoke test** — ✅ Approved (verified by running server)

## Files Created/Modified
- `frontend/embed.go` - Package `frontend` with `//go:embed dist; var FS embed.FS`
- `frontend/dist/.gitkeep` - Satisfies non-empty directory requirement for go:embed
- `frontend/dist/index.html` - Placeholder SPA shell with `</head>` for OG injection
- `cmd/server/main.go` - distFS setup, cover handler, SPA catch-all, buildOGTags helper
- `.env.example` - All 5 env vars with comments (DATABASE_URL, PORT, SITE_URL, APP_ENV, GOOGLE_BOOKS_API_KEY)

## Decisions Made
- **`frontend/embed.go` pattern**: `//go:embed` paths are relative to the source file. `cmd/server/main.go` cannot embed `../../frontend/dist` (parent directory traversal not allowed). Created `frontend/embed.go` at the same level as `frontend/dist/` and import it in main.go.
- **`filepath.Base` for cover URLs**: `cover_path` may be stored with OS-specific separators (`data\covers\file.jpg` on Windows, `data/covers/file.jpg` on Linux). `filepath.Base` extracts the filename correctly on both platforms.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] OG image URL contained full path instead of filename**
- **Found during:** Task 3 (smoke test via curl)
- **Issue:** `strings.TrimPrefix(*book.CoverPath, "data/covers/")` failed on Windows where path is `data\covers\1649374186.jpg` — produced `http://localhost:8081/covers/data\covers\1649374186.jpg`
- **Fix:** Replaced with `filepath.Base(*book.CoverPath)` for OS-agnostic extraction
- **Files modified:** cmd/server/main.go
- **Verification:** curl showed `og:image` = `http://localhost:8081/covers/1649374186.jpg`
- **Committed in:** d6dbdda

**2. [Structural] go:embed directive in frontend/embed.go instead of cmd/server/main.go**
- **Found during:** Task 1 (build failed with "no matching files found")
- **Issue:** Go's `//go:embed` resolves paths relative to the source file; `cmd/server/` has no `frontend/dist` subdirectory
- **Fix:** Created `frontend/embed.go` (package `frontend`) with the embed directive; imported as `flos-library/frontend` in main.go
- **Files modified:** frontend/embed.go (created), cmd/server/main.go (import added, var removed)
- **Verification:** `go build ./cmd/server/...` exits 0
- **Committed in:** 3cfc3f5

---

**Total deviations:** 2 auto-fixed (1 blocking OS path separator bug, 1 structural embed layout)
**Impact on plan:** Both fixes necessary for correctness. No scope creep. Embed package pattern is idiomatic Go for this project layout.

## Issues Encountered
- Windows path separator (`\`) in `cover_path` database column caused malformed og:image URLs — caught during human smoke test and fixed immediately

## User Setup Required
None - `.env.example` documents required env vars. Copy to `.env` and fill in values.

## Next Phase Readiness
- Backend API complete: all 8 read endpoints + cover serving + SPA routing + OG injection
- Phase 3 (frontend) can build against `/api/*` endpoints
- Phase 5 will replace `frontend/dist/index.html` placeholder with the real React build (same embed pattern, no server changes required)

---
*Phase: 02-go-rest-api*
*Completed: 2026-04-03*
