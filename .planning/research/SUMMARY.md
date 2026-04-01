# Project Research Summary

**Project:** Flo's Library
**Domain:** Personal book showcase website with automated Goodreads sync
**Researched:** 2026-04-01
**Confidence:** MEDIUM-HIGH (stack and pitfalls are HIGH; Goodreads RSS liveness is MEDIUM — must be manually verified)

---

## Executive Summary

Flo's Library is a personal reading showcase: a Go REST API + React TypeScript SPA that syncs book shelf data from Goodreads and presents it as a beautiful, always-current library. The core challenge is that Goodreads shut down its public API in 2020. The recommended replacement is RSS feed polling (every 6 hours) with Google Books API for metadata enrichment and OpenLibrary as a fallback — this approach has survived intact through mid-2025 and avoids the fragility of HTML scraping. The entire sync pipeline must treat RSS as a potentially impermanent source and include a CSV import endpoint as a standing fallback from day one.

The recommended architecture is a single deployable artifact: Go embeds the React build at compile time and serves everything from one binary on a personal VPS with systemd. In development, Vite dev server (port 5173) and Go API (port 8080) run separately with CORS enabled only for localhost. This eliminates CORS in production, removes the need for Nginx, and keeps deployment as simple as `scp binary server: && systemctl restart flos-library`. PostgreSQL from day one (Docker in development, system package in production) — there is no benefit to starting with SQLite given the dialect differences.

The primary risks are (1) Goodreads RSS disappearing without notice — mitigated by building CSV import as a first-class path; (2) RSS pagination silently truncating large shelves at 200 items — mitigated by implementing `?page=N` loops from the start; and (3) Google Books title+author searches returning wrong editions — mitigated by confidence-checking matched fields before accepting results. On the frontend, infinite scroll back-button breakage and missing book cover handling are the two most common failure modes, both with well-established patterns documented in PITFALLS.md.

---

## Key Findings

### Recommended Stack

The full stack is well-resolved: Chi v5 router, sqlc + pgx/v5 for the database layer, golang-migrate for migrations, PostgreSQL 16, gofeed for RSS parsing, and robfig/cron for the sync scheduler. On the frontend: Vite + React 18 + TypeScript, TanStack Query v5 for data fetching, and Intersection Observer for infinite scroll with cursor-based pagination. All choices are high-confidence, established 2024–2025 conventions. See STACK.md for full rationale and project layout.

**Core technologies:**
- **Chi v5** — Go HTTP router — idiomatic net/http, composable middleware, no framework magic
- **sqlc + pgx/v5** — database layer — type-safe generated code from raw SQL, no ORM overhead
- **golang-migrate** — migrations — numbered SQL files, embedded into binary, simple mental model
- **PostgreSQL 16 (Docker in dev)** — database — start real from day one; SQLite dialect divergence causes migration pain
- **gofeed** — RSS parsing — standard Go RSS/Atom library, handles GR custom extensions via `item.Extensions`
- **robfig/cron v3** (or `time.Ticker`) — sync scheduling — standard Go cron library
- **Vite + React 18 + TypeScript** — frontend toolchain — standard in 2025; CRA is deprecated
- **TanStack Query v5** — data fetching — handles caching, loading states, and `useInfiniteQuery` for infinite scroll
- **Go embed + http.FileServer** — production serving — single binary serves both API and React SPA

**Explicitly rejected:** GORM (reflection-based, poor complex query behavior), SQLite (dialect divergence), Next.js (SSR overhead not justified for personal site with no Google indexing need), Gin/Echo (non-idiomatic net/http wrappers), Lottie (60kb for a single decorative animation).

### Expected Features

The feature set is well-scoped. FEATURES.md provides a clear MVP build order and explicit anti-features to avoid. The single unresolved tension is SEO: the research recommends at minimum that Go generates per-book Open Graph meta tags in the HTML head, even without full SSR. This is achievable without Next.js.

**Must have (table stakes):**
- Book grid with cover images (aspect-ratio: 2/3, lazy loading, placeholder for missing covers)
- "Now Reading" section prominently on home page
- Book detail page (cover, title, author, description, genres, page count, read date)
- Author index + author detail pages
- Genre index + genre detail pages
- Reading Challenge / year view (year selector, stats strip, year-filtered book grid)
- Infinite scroll on "Books Read" (Intersection Observer + "Load More" button hybrid for accessibility)
- Responsive layout (mobile collapses to 2-column grid)
- Dark/light mode (CSS custom properties on `data-theme`, OS preference detection, localStorage persistence)

**Should have (differentiators):**
- Sidebar with animated book-reading SVG (CSS keyframes, no Lottie, `prefers-reduced-motion` respected)
- Brand palette fully applied: `#6d233e` wine red primary, antique gold `#c4843a` accent, warm paper tones for backgrounds
- Typography: Playfair Display (headings) + Inter (body) — serif/sans pairing with library feel
- Bio section on home page (photo, passions)
- Clean slug-based URLs (`/books/the-name-of-the-rose`)
- Per-book Open Graph meta tags served from Go (for social sharing previews)

**Defer to post-MVP:**
- JSON-LD structured data / Book schema markup
- Stats visualizations beyond basic counts on Reading Challenge page
- Search (small corpus; browser Ctrl+F sufficient; PostgreSQL tsvector column can be stubbed in schema now)
- Progressive Web App / offline support
- WebP conversion for cover images (complexity-to-benefit ratio too high for personal site traffic)

**Explicit anti-features (do not build):**
- Star ratings display (Goodreads is source of truth; showing GR data as-is is confusing)
- User reviews / comments (auth + moderation complexity, out of scope)
- Numbered pagination (conflicts with infinite scroll requirement)
- Social sharing widgets (Open Graph tags are sufficient)
- Reading progress percentage (Goodreads doesn't expose granular progress via RSS)

### Architecture Approach

The system has four functional layers: (1) the Goodreads sync pipeline (RSS polling → diff → enrichment → cover download), (2) the PostgreSQL data store with a source-agnostic schema, (3) the Go REST API, and (4) the React SPA. The sync pipeline is intentionally decoupled from the data model — `goodreads_id` is just a string identifier, so switching from RSS to CSV import or scraping requires no schema changes. Cover images are stored at `./data/covers/{isbn13}.jpg` and served with `Cache-Control: immutable` headers.

**Major components:**
1. **Sync worker** — fetches Goodreads RSS (all pages), diffs against DB, calls Google Books by ISBN (then title+author with confidence check), falls back to OpenLibrary, downloads and validates covers, upserts to DB
2. **Cron scheduler** — runs sync worker every 6 hours via `time.Ticker` or robfig/cron; also exposed as `POST /admin/sync` for manual trigger
3. **PostgreSQL schema** — books, authors, book_authors, genres, book_genres tables; `goodreads_id` unique key; `slug` for URL routing; `cover_path`, `metadata_source`, `read_count` columns
4. **REST API** — Chi router; `GET /api/books` (cursor-paginated), `GET /api/books/:slug`, `GET /api/authors`, `GET /api/genres`, `GET /api/years`; covers served from `/covers/*` with immutable cache headers
5. **React SPA** — TanStack Query for server state; Intersection Observer infinite scroll; CSS custom properties theming; book grid with placeholder pattern for missing covers
6. **CSV import endpoint** — `POST /admin/import-csv` for Goodreads CSV export — standing fallback if RSS goes down

**Data flow for sync:**
```
GR RSS (currently-reading + read, all pages)
  → merge shelves (read wins conflicts)
  → diff against DB
  → new books: Google Books (ISBN) → Google Books (title+author) → OpenLibrary → GR thumbnail
  → download + validate cover image
  → upsert books, authors, genres
```

### Critical Pitfalls

Five pitfalls that require explicit design decisions before implementation begins:

1. **RSS pagination truncation at 200 items** — Use `?per_page=200&page=N` loop from the first sync implementation; break when response returns fewer than 200 items. No error is returned when truncated — the feed silently ends. See PITFALLS.md #1 for the Go implementation pattern.

2. **Google Books title+author false positives** — When ISBN is missing, title+author search can match the wrong book (different edition, anthology, similarly titled work). Implement a confidence check: returned author must fuzzy-match input author, returned title must contain input title. On failure: store book with no enrichment and log for manual review; never store a wrong match. Add `metadata_source` column: `'isbn' | 'title_author' | 'manual' | 'none'`.

3. **CORS in production** — Vite's proxy config makes CORS invisible in development; first production deploy reveals the issue. Prevention: use Go-serves-React topology (Option A) in production — everything is same-origin, CORS headers never needed. See STACK.md #9 and PITFALLS.md #5.

4. **Goodreads RSS custom fields not in standard RSS** — `<book_id>`, `<author_name>`, `<isbn13>`, `<user_read_at>` are GR extensions accessed via `item.Extensions` in gofeed, not standard fields. Key paths must be verified against the live feed before writing sync code. Pre-implementation step: fetch live feed, print all `item.Extensions` keys, build field mapping struct. See PITFALLS.md #2 and #10 for field reliability table.

5. **Infinite scroll back-button position loss** — User scrolls 80 books deep, clicks a book, hits Back, returns to position 0. Prevention: persist cursor offset in URL as query param; use `sessionStorage` to store last-viewed book ID on click; on list mount, `scrollIntoView` to that element. See PITFALLS.md #7 for the React implementation.

**Additional pitfalls requiring attention:**
- `user_read_at` timezone shifts year assignment (extract date portion before UTC conversion; see PITFALLS.md #15)
- Shelf transition race condition when `currently-reading` and `read` both contain the same book (process `read` last, let it win; see PITFALLS.md #11)
- Cover image validation — Google Books can return 1×1 pixel placeholders at valid URLs; validate file size (>5KB) and decodability after download (see PITFALLS.md #8)
- `VITE_` prefixed variables are bundled into the browser — Google Books API key must only exist as a Go backend env var, never in React (see PITFALLS.md #14)
- Slug collisions for books with identical titles — append publication year, then author surname (see PITFALLS.md #16)
- `user_read_at` is frequently empty even for "read" shelf books — fall back to `user_date_added` (see PITFALLS.md #10)

---

## Open Questions — Decisions Required Before Implementation

These are unresolved questions that have direct implementation impact and must be answered before or during Phase 1:

| # | Question | Decision Needed | Impact |
|---|----------|-----------------|--------|
| 1 | **Is the Goodreads RSS feed live?** | Manually open `https://www.goodreads.com/review/list_rss/79499864?shelf=read` before writing any sync code | If dead: CSV import becomes primary; architecture changes |
| 2 | **How many books are on Florian's "read" shelf?** | Check Goodreads profile for total count | If >200: pagination must be implemented in Phase 1, not deferred |
| 3 | **Does `?per_page=200&page=2` actually work on the live feed?** | Test manually during RSS verification | Affects sync pagination implementation |
| 4 | **Open Graph meta tags from Go?** | Decide whether Go template-renders per-book meta tags into the SPA HTML shell, or skip entirely | If yes: adds a Go template rendering step; small but real complexity |
| 5 | **Sidebar animation: existing asset or new creation?** | Does the prior PHP/React version have an SVG that can be ported? | If yes: simplifies sidebar phase; if no: SVG + CSS animation must be created from scratch |
| 6 | **Cover quality preference** | Accept Google Books thumbnail quality (~128px), or invest in fetching largest available size? | Affects cover download pipeline complexity |
| 7 | **Re-read handling** | If Florian re-reads a book, should both read dates be stored, or just the latest? | Affects schema: single `read_at` vs `book_reads` join table; `read_count` column is the minimum viable approach |
| 8 | **Production hosting target** | VPS with systemd (recommended) or Docker? | Affects deployment docs and Makefile targets |
| 9 | **Google Books API key** | Create Google Cloud project + enable Books API before enrichment phase | Blocks Phase 2 if not done; 1,000 req/day free tier is sufficient |

---

## Implications for Roadmap

### Phase 1: Foundation + Data Pipeline

**Rationale:** Everything downstream depends on having books in the database. Without a working sync pipeline, there is nothing to display. The data model must be correct from the start — schema changes with existing data are painful.

**Delivers:**
- PostgreSQL schema with golang-migrate (books, authors, genres, join tables)
- Goodreads RSS sync pipeline (pagination, field parsing, shelf merge, delta diff)
- Google Books enrichment (ISBN lookup, title+author fallback with confidence check)
- OpenLibrary fallback for covers and missing metadata
- Cover image download pipeline (deduplication, validation, filesystem storage)
- CSV import endpoint (`POST /admin/import-csv`) as standing fallback
- Manual sync endpoint (`POST /admin/sync`)
- Cron scheduler (6-hour interval)
- Go project structure (cmd/, internal/, migrations/, data/covers/)
- Docker Compose for local PostgreSQL

**Must avoid:**
- Pitfall #1 (RSS truncation) — implement `?page=N` loop immediately
- Pitfall #2 (custom field parsing) — verify `item.Extensions` key paths against live feed before writing struct
- Pitfall #3 (Google Books false positives) — implement confidence check from the start
- Pitfall #10 (field reliability) — treat `user_read_at` as optional; use `user_date_added` as fallback
- Pitfall #11 (shelf race condition) — process `read` shelf after `currently-reading`
- Pitfall #15 (timezone shift) — extract date portion before UTC conversion

**Research flag:** YES — verify Goodreads RSS liveness (60 seconds in browser) before writing a single line of sync code. If dead, Phase 1 becomes CSV-first.

### Phase 2: Go REST API

**Rationale:** Once books are in the database, expose them through a clean API. This phase is well-patterned (Chi + sqlc) and does not require research. The API contract here defines the React data shapes.

**Delivers:**
- Chi router setup with middleware (logging, recovery, CORS for dev)
- sqlc-generated queries for all read operations
- `GET /api/books` with cursor-based pagination (`last_read_at` + `id`)
- `GET /api/books/:slug`
- `GET /api/authors`, `GET /api/authors/:slug`
- `GET /api/genres`, `GET /api/genres/:slug`
- `GET /api/years` (distinct years for Reading Challenge selector)
- `/covers/*` static file server with `Cache-Control: immutable` middleware
- CORS middleware restricted to `localhost:5173` in development
- Go embed setup for React dist (production topology established even if not used yet)

**Must avoid:**
- Pitfall #12 (missing cache headers on covers) — add middleware in this phase, not later
- Pitfall #14 (API key in frontend) — API key stays in Go env; never exposed via API responses

**Research flag:** NO — Chi + sqlc patterns are well-documented in STACK.md.

### Phase 3: React Frontend — Core Library

**Rationale:** Build the primary browsing experience. The book grid is the central UI component that every other page reuses — build it correctly once. Foundation (CSS variables, theming, typography) must come before components.

**Delivers:**
- Vite + React 18 + TypeScript project scaffolded inside `frontend/`
- CSS custom properties foundation (light/dark mode, full color palette, typography scale)
- Dark/light mode toggle with OS preference detection and localStorage persistence
- Book cover component (aspect-ratio: 2/3, lazy loading, CSS gradient placeholder for missing covers)
- Book grid component (reusable, accepts filter props for year/author/genre)
- Infinite scroll with Intersection Observer + explicit "Load More" button fallback
- Back-button position restoration (URL cursor param + sessionStorage scroll target)
- Home page: bio section + "Now Reading" shelf + "Books Read" grid
- TanStack Query v5 setup with cursor-based infinite query

**Must avoid:**
- Pitfall #7 (back-button loss) — URL cursor persistence and sessionStorage scroll target from first implementation
- Pitfall #5 (CORS) — Vite proxy to Go in dev; Go-serves-React in production (establish this topology in this phase)

**Research flag:** NO — TanStack Query infinite scroll and Intersection Observer patterns are in STACK.md and FEATURES.md.

### Phase 4: React Frontend — Detail and Index Pages

**Rationale:** Secondary pages that reuse the book grid component built in Phase 3. Low risk, high completeness value. Establishes full slug-based routing.

**Delivers:**
- React Router setup (slug-based routing for all entity types)
- Book detail page (two-column layout, full metadata display, expandable description)
- Author index page (alphabetical list with book count)
- Author detail page (book grid filtered by author)
- Genre index page (sorted by book count, CSS bar visualization)
- Genre detail page (book grid filtered by genre)
- Reading Challenge page (year selector stored in URL, stats strip, year-filtered book grid)
- Open Graph meta tags from Go (per-book HTML head injection for social sharing previews)

**Must avoid:** Slug collision handling (Pitfall #16) — append year and author surname on conflict — should already be in the sync pipeline from Phase 1, but must be verified here.

**Research flag:** NO — standard React Router + component composition patterns.

### Phase 5: Polish and Production Deployment

**Rationale:** Final refinements and production hardening. Kept separate so it does not block the working application from being usable.

**Delivers:**
- Sidebar animated book SVG (CSS keyframe page-turn or breathing effect, `prefers-reduced-motion` respected)
- Production Go build embedding React dist (verify `//go:embed frontend/dist` setup)
- systemd service file for VPS deployment
- SSL via certbot with systemd timer (or Caddy for zero-config auto-HTTPS)
- Makefile targets: `make dev`, `make build`, `make migrate`, `make sync`
- Environment variable audit (no secrets in `VITE_` vars; confirm by inspecting built JS bundle)
- Final palette and typography pass

**Research flag:** NO — deployment patterns are in PITFALLS.md #19 and #20; CSS animation in FEATURES.md #6.

### Phase Ordering Rationale

- Phase 1 must come first: no data = nothing to display or test
- Phase 2 before Phase 3: the API contract defines React's data shapes; building UI before API means building against assumptions
- Phase 3 before Phase 4: the book grid component is shared by all index and detail pages; building it first means Phase 4 is assembly, not invention
- Phase 5 last: decoration and deployment can be done iteratively; the app should be usable (if not polished) after Phase 3

### Research Flags Summary

| Phase | Research Needed? | Reason |
|-------|-----------------|--------|
| Phase 1 | YES — pre-implementation | Verify Goodreads RSS feed is live; inspect raw XML for `item.Extensions` key paths |
| Phase 2 | NO | Chi + sqlc patterns fully documented in STACK.md |
| Phase 3 | NO | Vite + TanStack Query + Intersection Observer patterns fully documented |
| Phase 4 | NO | Standard React Router + component composition |
| Phase 5 | NO | Deployment patterns (systemd, certbot/Caddy) documented in PITFALLS.md |

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | All choices are well-established 2024–2025 conventions; no experimental dependencies |
| Features | HIGH | Clear MVP scope; differentiators and anti-features are explicit |
| Architecture | MEDIUM | Sync pipeline architecture is solid; Goodreads RSS liveness is the one live unknown |
| Pitfalls | HIGH | Pitfalls are specific, well-sourced, and all include concrete mitigations |

**Overall confidence:** MEDIUM-HIGH

The only genuinely uncertain element is whether Goodreads RSS feeds are still live at implementation time. Everything else is HIGH confidence based on stable APIs, established libraries, and well-documented patterns. The architecture is designed to be resilient to RSS going down (CSV import fallback, source-agnostic data model), so even the uncertain element has a mitigation path.

### Gaps to Address

- **Goodreads RSS liveness:** Must be manually verified before Phase 1 begins. This takes 60 seconds. Do it first.
- **`gofeed` extension key paths for GR custom fields:** The exact `item.Extensions` map keys for `<book_id>`, `<author_name>`, `<isbn13>`, `<user_read_at>` must be confirmed against the live feed. Write a test harness that fetches one item and prints the full extensions map before building the field-mapping struct.
- **Florian's "read" shelf size:** Check total book count on Goodreads profile before Phase 1. If >200, pagination is critical path in Phase 1. If ≤200 today, still implement pagination as a first-class concern (the shelf will grow).
- **Open Graph meta tag approach:** Decide in Phase 2 whether Go template-renders per-book meta tags into the SPA HTML shell. This is a small feature with a disproportionate impact on social sharing UX. Recommend yes — it requires Go to inspect the URL and inject 5–6 `<meta>` tags into the index.html head.
- **Google Books API key:** Obtain before Phase 1 enrichment work begins. Free tier (1,000 req/day) is sufficient; getting the key requires creating a Google Cloud project and enabling the Books API.

---

## Sources

### Primary (HIGH confidence)
- STACK.md (this project) — Chi v5, sqlc, golang-migrate, gofeed, TanStack Query v5, Go embed patterns
- PITFALLS.md (this project) — 20 documented pitfalls with concrete mitigations
- FEATURES.md (this project) — complete feature inventory with priority and anti-feature rationale
- ARCHITECTURE.md (this project) — Goodreads sync pipeline, data model, API routes, resilience patterns
- Google Books API v1 official documentation — field availability, authentication, rate limits
- OpenLibrary API documentation — cover URL format, search endpoints

### Secondary (MEDIUM confidence)
- Goodreads developer community documentation — RSS feed URL format and custom field names
- Community reports (via training data) — RSS feed liveness through mid-2025; `per_page=200` cap
- Google Books cover URL `zoom` parameter behavior — undocumented; test per-book

### Tertiary (LOW confidence)
- Goodreads HTML scraping selectors — confirmed fragile; use only if RSS is confirmed dead
- hardcover.app API status — not applicable for this use case (would require Florian to migrate platforms)

---

*Research completed: 2026-04-01*
*Ready for roadmap: yes — pending manual RSS verification*
