# Roadmap: Flo's Library

**Created:** 2026-04-01
**Milestone:** v1.0 — Full personal book library, synced from Goodreads

---

## Milestone Goal

A fully working personal book showcase website — synced from Goodreads, displaying all books read with detail pages, author/genre indexes, and a Reading Challenge view. Deployable to a personal VPS as a single Go binary.

---

## Phase Overview

| Phase | Name | Goal | Plans | Research |
|-------|------|------|-------|----------|
| 1 | Data Pipeline | Books in the database, automatically | 3 | YES — verify RSS liveness first |
| 2 | Go REST API | Clean API contract for the frontend | 2 | NO |
| 3 | Frontend Core | Home page with browsable book library | 3 | NO |
| 4 | Frontend Pages | Detail, author, genre, and challenge pages | 2 | NO |
| 5 | Polish & Deploy | Production-ready, deployed on VPS | 2 | NO |

---

## Phase 1: Data Pipeline

**Goal:** Books from Goodreads are in PostgreSQL, enriched with metadata from Google Books, with covers stored locally. The sync runs automatically every 6 hours.

**Why first:** Everything else is blocked on having book data. No data = nothing to display or test.

**Pre-implementation requirement:** Manually verify RSS feed before writing any sync code:
- Open `https://www.goodreads.com/review/list_rss/79499864?shelf=read` in a browser
- If XML returns → RSS is viable (primary path)
- If 403 / redirect → CSV import becomes primary path

**Requirements covered:** SYNC-01–09, DATA-01–05

### Plan 1.1 — Project Foundation & Schema

**Objective:** Go project structure, PostgreSQL schema, and migration toolchain in place.

**Tasks:**
1. Initialize Go module (`flos-library`), project layout (`cmd/`, `internal/`, `migrations/`, `data/covers/`)
2. Docker Compose with PostgreSQL 16 + adminer for local dev
3. golang-migrate installed and configured; initial migration: books, authors, genres, book_authors, book_genres tables
4. Schema columns: `id`, `goodreads_id` (unique), `slug` (unique), `title`, `description`, `cover_path`, `page_count`, `publication_year`, `isbn13`, `metadata_source`, `read_at`, `date_added`, `read_count`, `shelf`, `created_at`, `updated_at`
5. Makefile targets: `make dev` (start docker + run Go), `make migrate`, `make migrate-down`

### Plan 1.2 — Goodreads RSS Sync

**Objective:** RSS polling pipeline that fetches all shelves, paginates past 200, and upserts to DB.

**Pre-task:** Fetch live RSS feed, print all `item.Extensions` keys, document exact key paths for `book_id`, `isbn13`, `author_name`, `user_read_at`, `user_date_added`, `image_url`.

**Tasks:**
1. `gofeed` RSS parser integration; fetch `currently-reading` and `read` shelves
2. Pagination loop: `?per_page=200&page=N` until response < 200 items
3. Field mapping struct from verified `item.Extensions` key paths
4. Shelf merge: process `currently-reading` first, then `read` (read wins conflicts)
5. Delta diff: compare fetched GR IDs against DB; detect new, updated (shelf change), unchanged
6. Slug generation: title → kebab-case; collision resolved by appending year, then author surname
7. CSV import endpoint `POST /admin/import-csv` (shared upsert logic with RSS path)
8. Manual trigger `POST /admin/sync`; 6-hour cron via `time.Ticker`

### Plan 1.3 — Metadata Enrichment & Cover Download

**Objective:** Each book enriched with Google Books metadata; cover images downloaded and validated.

**Pre-task:** Create Google Cloud project, enable Books API, obtain API key (store as `GOOGLE_BOOKS_API_KEY` env var).

**Tasks:**
1. Google Books client: ISBN-based lookup (primary), title+author fallback
2. Confidence gate for title+author fallback: returned author must fuzzy-match, title must contain input title; on failure → `metadata_source = 'none'`, log for review
3. OpenLibrary fallback for cover when Google Books lacks cover or returns 1×1 placeholder
4. Cover download pipeline: fetch URL, validate (file size > 5KB, decodable), store at `data/covers/{isbn13}.jpg`
5. Update book record: description, genres, page_count, publication_year, cover_path, metadata_source
6. Author upsert with slug; genre upsert with slug; join table entries

---

## Phase 2: Go REST API

**Goal:** A clean, documented REST API that the React frontend can consume. Covers served with proper cache headers. Go project structured for embedding React build at production build time.

**Requirements covered:** API-01–10

### Plan 2.1 — Core API Endpoints

**Objective:** All read endpoints implemented with Chi, sqlc queries, and cursor-based pagination.

**Tasks:**
1. Chi router setup: middleware stack (logger, recoverer, CORS for `localhost:5173`)
2. sqlc query generation for all read operations
3. `GET /api/books` — cursor-paginated (`last_read_at` + `id`), supports `?year=` and `?shelf=` filters
4. `GET /api/books/currently-reading` — returns currently-reading shelf
5. `GET /api/books/:slug` — full book detail
6. `GET /api/authors` + `GET /api/authors/:slug` — with cursor-paginated books
7. `GET /api/genres` + `GET /api/genres/:slug` — with cursor-paginated books, sorted by book count
8. `GET /api/years` — distinct years with book count

### Plan 2.2 — Cover Serving & Production Topology

**Objective:** Cover images served with immutable cache headers. Go embed setup for React build. Open Graph meta tag injection.

**Tasks:**
1. `/covers/*` file server with `Cache-Control: public, max-age=31536000, immutable` middleware
2. `//go:embed frontend/dist` embed setup (placeholder; will be populated in Phase 5)
3. SPA catch-all handler serving `index.html` for all non-API routes
4. Open Graph meta tag injection: Go inspects `/books/:slug` path, fetches book from DB, template-renders 5–6 `<meta>` tags into `index.html` head before serving
5. Environment variable documentation (`.env.example` with `DATABASE_URL`, `GOOGLE_BOOKS_API_KEY`, `PORT`)

---

## Phase 3: Frontend Core

**Goal:** Home page with bio, "Now Reading", and infinite-scrolling "Books Read" list. Dark/light mode. Full CSS foundation.

**Requirements covered:** HOME-01–04, UI-01–03, UI-06–10

### Plan 3.1 — Foundation & Design System

**Objective:** Vite + React + TypeScript scaffolded. CSS custom properties system with full brand palette, typography, and dark/light mode.

**Tasks:**
1. `npm create vite@latest frontend -- --template react-ts`; Vite proxy to Go API in dev config
2. TanStack Query v5 (`QueryClientProvider`) setup
3. React Router v6 setup (routes: `/`, `/books/:slug`, `/authors`, `/authors/:slug`, `/genres`, `/genres/:slug`, `/reading-challenge`)
4. CSS custom properties: `--color-primary`, `--color-accent`, `--color-background`, `--color-surface`, `--color-text` for light and dark themes on `[data-theme]` attribute
5. Dark/light mode: detect `prefers-color-scheme`, persist in localStorage, toggle component
6. Typography: import Playfair Display + Inter from Google Fonts; set CSS variables for font families and type scale
7. Book cover component: `aspect-ratio: 2/3`, `object-fit: cover`, lazy loading, CSS gradient placeholder using brand colors

### Plan 3.2 — Book Grid & Infinite Scroll

**Objective:** Reusable book grid component with Intersection Observer infinite scroll and back-button position restoration.

**Tasks:**
1. BookGrid component: accepts `queryKey`, `fetchFn`, renders book cards in responsive CSS grid
2. BookCard component: cover image, title, author name(s)
3. TanStack Query `useInfiniteQuery` with cursor-based pagination
4. Intersection Observer sentinel: trigger `fetchNextPage` when sentinel enters viewport
5. "Load More" button as accessible fallback (same `fetchNextPage` call)
6. Back-button position restoration: persist cursor and scroll target book ID in URL + sessionStorage; on mount, `scrollIntoView` to target if present

### Plan 3.3 — Home Page

**Objective:** Fully assembled home page with bio section, Now Reading, and Books Read.

**Tasks:**
1. Bio section: Florian's photo + description text (hardcoded or from a config file)
2. "Now Reading" section: horizontal shelf of currently-reading books using `GET /api/books/currently-reading`
3. "Books Read" section: BookGrid with infinite scroll using `GET /api/books`
4. Sidebar navigation: links to Home, Authors, Genres, Reading Challenge; placeholder for animation (Phase 5)
5. Layout: sidebar + main content area; mobile responsive (sidebar collapses to top nav)

---

## Phase 4: Frontend Pages

**Goal:** All secondary pages — book detail, author index/detail, genre index/detail, and Reading Challenge.

**Requirements covered:** BOOK-01–04, AUTH-01–02, GENR-01–02, CHAL-01–04

### Plan 4.1 — Book Detail & Author/Genre Pages

**Objective:** Book detail page and all index/detail pages for authors and genres.

**Tasks:**
1. Book detail page: two-column layout (cover left, metadata right), title, author links, genre tag links, description (expandable if long), publication year, page count, read date
2. Author index page: alphabetical list with book count per author; links to author detail
3. Author detail page: author name + BookGrid filtered to that author
4. Genre index page: genres sorted by book count, CSS bar visualization
5. Genre detail page: genre name + BookGrid filtered to that genre

### Plan 4.2 — Reading Challenge Page

**Objective:** Year-in-books view with year selector, stats, and filtered book grid.

**Tasks:**
1. `GET /api/years` powers year selector (current year default, URL param `?year=`)
2. Stats strip: book count for selected year; optionally total pages (from `page_count` sum, noting ~70-80% coverage)
3. BookGrid filtered by selected year
4. Year selector component: prev/next arrows + year display; keyboard accessible
5. Verify slug collision handling from Phase 1 is correct (books with identical titles)

---

## Phase 5: Polish & Production Deployment

**Goal:** Sidebar animation complete, Go binary embeds React build, deployed on VPS with SSL and systemd.

**Requirements covered:** UI-04, UI-05, DEPL-01–04

### Plan 5.1 — Sidebar Animation & UI Polish

**Objective:** Animated book SVG in sidebar. Final palette and typography pass.

**Tasks:**
1. SVG book graphic: open book with pages; CSS keyframe animation (page-turn or breathing/glow effect)
2. `@media (prefers-reduced-motion: reduce)` — disable animation, show static SVG
3. Final color/typography pass: verify all brand colors applied consistently, check contrast ratios
4. Favicon and page title configuration

### Plan 5.2 — Production Build & Deployment

**Objective:** Single Go binary serving everything, deployed on VPS via systemd.

**Tasks:**
1. `npm run build` produces `frontend/dist/`; verify `//go:embed frontend/dist` compiles correctly
2. systemd unit file (`flos-library.service`): restart on failure, environment file for secrets
3. Makefile `make build` target: builds React, compiles Go binary with embed
4. SSL: certbot with systemd timer (or Caddy for zero-config)
5. Environment variable audit: inspect built JS bundle to confirm no `GOOGLE_BOOKS_API_KEY` leakage
6. Deployment runbook: `make build` → `scp binary server:~/ && systemctl restart flos-library`

---

## Critical Path

```
Phase 1 (data) → Phase 2 (API) → Phase 3 (frontend core) → Phase 4 (pages) → Phase 5 (deploy)
```

No phase can start before the previous completes. The critical path is strictly sequential because:
- Phase 2 depends on Phase 1 data to test against
- Phase 3 depends on Phase 2 API contract to build against
- Phase 4 reuses Phase 3 components
- Phase 5 wraps Phase 3+4 output

---

## Pre-Implementation Checklist

Before starting Phase 1:

- [ ] Verify Goodreads RSS feed: open `https://www.goodreads.com/review/list_rss/79499864?shelf=read` in browser
- [ ] Note total book count on Goodreads "read" shelf (if >200, pagination is critical path in Plan 1.2)
- [ ] Create Google Cloud project and enable Books API; save API key for Plan 1.3
- [ ] Decide: Open Graph meta tags from Go (yes/no) — recommended YES, implemented in Plan 2.2

---

*Roadmap created: 2026-04-01*
*Next: `/gsd:plan-phase 1`*
