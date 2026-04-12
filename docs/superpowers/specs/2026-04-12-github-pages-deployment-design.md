# GitHub Pages Deployment Design

**Date:** 2026-04-12
**Project:** Flo's Library
**Status:** Approved by user

---

## Goal

Deploy Flo's Library at zero cost, forever, with no self-hosted infrastructure. All hosting on GitHub. Daily Goodreads sync continues automatically via GitHub Actions cron.

---

## Architecture

The Go HTTP server is replaced by a Go CLI that generates static files. GitHub Pages serves those files. No server runs in production.

```
CURRENT                          TARGET
───────────────────────────────  ────────────────────────────────────
Go HTTP server (Chi)             Go CLI (cmd/generate)
PostgreSQL (Docker)              SQLite file (committed to repo)
In-process scheduler             GitHub Actions cron (daily)
Cover files on local disk        Covers committed to repo → jsDelivr CDN
React fetches live API           React fetches static JSON files
Single binary deploy             GitHub Pages (no server at all)
```

### Repo Layout After Migration

```
flos-library/
├── cmd/
│   └── generate/
│       └── main.go              ← new entrypoint (replaces cmd/server)
├── internal/
│   ├── sync/                    ← kept: rss.go, enricher.go, covers.go
│   ├── db/                      ← retargeted: sqlc output for SQLite
│   └── generate/                ← new: JSON serializers, output writers
├── data/
│   ├── floslib.db               ← SQLite file (committed)
│   └── covers/                  ← cover images (committed)
├── frontend/
│   ├── src/
│   └── dist/                    ← built output, committed for Pages
├── static/
│   ├── books.json
│   ├── authors.json
│   ├── genres.json
│   └── books/
│       └── {slug}.json          ← per-book detail file
├── .github/
│   └── workflows/
│       ├── sync.yml             ← daily cron: sync → commit → deploy
│       └── deploy.yml           ← on push to main: build React → Pages
└── docs/
    └── superpowers/specs/
        └── this file
```

---

## Go Changes

### Removed

| File/Package | Reason |
|---|---|
| `cmd/server/main.go` | No HTTP server |
| `internal/api/` | No HTTP handlers |
| `internal/scheduler/` | GitHub Actions replaces it |
| `frontend/embed.go` | No binary embed |
| `docker-compose.yml` | No Postgres |

### Kept (unchanged or near-unchanged)

| Package | Notes |
|---|---|
| `internal/sync/rss.go` | No changes |
| `internal/sync/enricher.go` | No changes |
| `internal/sync/covers.go` | Minor: write to `data/covers/` instead of server path |

### Added

**`cmd/generate/main.go`** — new CLI entrypoint:
1. Opens `data/floslib.db` (SQLite)
2. Runs sync pipeline (RSS → enrich → covers)
3. Queries DB and writes static JSON to `static/`
4. Exits non-zero on error (GitHub Actions fails loudly)

**`internal/generate/`** — JSON serializers:
- `books.go` — produces `static/books.json` and `static/books/{slug}.json`
- `authors.go` — produces `static/authors.json`
- `genres.go` — produces `static/genres.json`

Cover URLs use jsDelivr CDN pattern:
```go
const cdnBase = "https://cdn.jsdelivr.net/gh/OWNER/REPO@main/data/covers/"

func coverURL(filename string) string {
    return cdnBase + filename
}
```
`OWNER/REPO` is the GitHub username and repository name (e.g. `florianabel/flos-library`). This value is hardcoded in the Go CLI — it does not change at runtime.

### Database Migration

| Layer | Current | Target |
|---|---|---|
| Driver | `pgx/v5` | `modernc.org/sqlite` (pure Go, no CGO) |
| DSN | `postgres://...` | `file:data/floslib.db` |
| Migrations | golang-migrate (Postgres SQL) | golang-migrate (SQLite SQL) |
| sqlc config | pgx/v5 driver, Postgres dialect | modernc.org/sqlite driver, SQLite dialect |
| Generated code | `internal/db/*.go` | Regenerated from SQLite schema |

**SQLite schema differences from Postgres:**
- `SERIAL`/`BIGSERIAL` → `INTEGER PRIMARY KEY`
- `TEXT[]` arrays → join table or JSON column
- `TIMESTAMPTZ` → `TEXT` (ISO-8601)

### `go.mod` Changes

```
REMOVE: github.com/jackc/pgx/v5
REMOVE: github.com/jackc/pgpassfile
REMOVE: github.com/jackc/pgservicefile
REMOVE: github.com/jackc/puddle/v2
REMOVE: github.com/go-chi/chi/v5
REMOVE: github.com/go-chi/cors
ADD:    modernc.org/sqlite
```

---

## React Changes

### Removed

| Current | Replacement |
|---|---|
| `fetch('/api/...')` calls | `fetch('/static/...')` |
| Server-side pagination params | Client-side slice |
| `VITE_API_BASE_URL` env var | Removed |
| Vite proxy to Go dev server | Removed |

### Data Loading

TanStack Query stays. `queryFn` targets static JSON instead of live endpoints:

```ts
// Before
queryFn: () => fetch('/api/books?page=1&limit=20').then(r => r.json())

// After
queryFn: () => fetch('/static/books.json').then(r => r.json())
```

**Load strategy:**
- `books.json`, `authors.json`, `genres.json` — loaded once, held in TanStack Query cache
- `books/{slug}.json` — loaded on demand per book detail page

### Filtering and Search

All filtering (genre, shelf, author) and search (title) become client-side JS array operations — instant, no loading spinners, no debounce needed.

### Pagination

Keep windowed display (show first N, load more on scroll). No new dependencies. Functionally identical to current UX.

### `vite.config.ts`

Remove the proxy block pointing to the Go dev server. Vite serves `static/` directly in dev mode.

### `frontend/dist/` Committed

React build output is committed to repo so Pages can serve it. `deploy.yml` rebuilds and recommits on every push to `main`. Vite's default `.gitignore` excludes `dist/` — that entry must be removed as part of this migration.

---

## GitHub Actions Workflows

### `sync.yml` — Daily data sync

**Trigger:** `cron '0 6 * * *'` (6am UTC daily) + `workflow_dispatch` (manual)

**Steps:**
1. Checkout repo (full history)
2. Set up Go
3. `go run ./cmd/generate` — runs full sync pipeline
4. `git diff --quiet` — check if anything changed
5. If changed: `git add data/ static/ && git commit && git push`
6. Push to `main` triggers `deploy.yml` automatically

**Secrets required** (set in GitHub repo settings):
- `GOOGLE_BOOKS_API_KEY` — injected as env var, never committed

### `deploy.yml` — Build and deploy React

**Trigger:** push to `main`

**Steps:**
1. Checkout repo
2. Set up Node
3. `npm ci && npm run build` (in `frontend/`)
4. Deploy `frontend/dist/` to GitHub Pages via `actions/deploy-pages`

### Why Two Workflows

Sync runs on a schedule. Deploy runs on code push. Separating them avoids re-running RSS sync on every code change, and keeps sync commits from requiring developer involvement.

---

## Local Dev Workflow

```bash
# One-time setup
cp .env.example .env   # add GOOGLE_BOOKS_API_KEY

# Regenerate data (sync + write JSON)
go run ./cmd/generate

# Start frontend dev server
cd frontend && npm run dev
```

No Docker, no Postgres, no background server.

**For SQLite schema changes:**
```bash
# Edit migration SQL
# Apply migrations + regenerate sqlc
go run ./cmd/generate --migrate-only
sqlc generate
go run ./cmd/generate
```

---

## Error Handling

- Go CLI exits non-zero on any error → GitHub Actions marks the run as failed → no partial commit
- If RSS fetch fails, sync aborts — existing `static/` and `data/` remain unchanged (last good state)
- jsDelivr CDN: covers are committed files, not fetched at runtime — no CDN failure affects data generation
- React: if a static JSON fetch fails, TanStack Query retries and shows an error state (existing behavior)

---

## What This Replaces from Phase 5 Plan

The original Phase 5 planned:
- Lottie sidebar animation — **unaffected, still in scope**
- UI polish pass — **unaffected, still in scope**
- `make build` single binary — **superseded by this design**
- Raspberry Pi deployment runbook — **superseded by GitHub Pages**

The deployment portion of Phase 5 is replaced entirely by this design. Sidebar animation and UI polish proceed as planned.

---

## Out of Scope

- SSL certificate management (GitHub Pages handles HTTPS automatically)
- Custom domain setup (optional, future)
- CI/CD beyond what's described (no staging environment, no automated tests in CI for now)
- WebP cover conversion (v2 requirement)
