---
name: Flo's Library project context
description: Personal book showcase site for Florian — Go + React TS, syncs from Goodreads RSS, 5-phase roadmap initialized
type: project
---

Flo's Library is Florian's personal book showcase website. Initialized 2026-04-01.

**Why:** Rebuild of a 6-year-old PHP/MySQL/React app. Goodreads API was deprecated; needs new sync mechanism.

**Stack:** Go (Chi v5, sqlc, pgx/v5, golang-migrate) + React TypeScript (Vite, TanStack Query v5). Separate services in dev, Go embeds React build in production. PostgreSQL from day one.

**Key decisions:**
- Goodreads RSS feed (primary sync) + Google Books API enrichment + OpenLibrary fallback + CSV import endpoint
- Self-host book cover images (downloaded, not hotlinked)
- Modernized sidebar animation (CSS keyframes, not Lottie)
- Brand color: #6d233e wine red primary, #c4843a gold accent

**5-phase roadmap:**
1. Data Pipeline (RSS sync, Google Books enrichment, cover download, schema)
2. Go REST API (Chi endpoints, cover serving, Open Graph meta)
3. Frontend Core (home page, book grid, infinite scroll, dark/light mode)
4. Frontend Pages (book detail, author/genre pages, Reading Challenge)
5. Polish & Deploy (sidebar animation, single Go binary, systemd, SSL)

**How to apply:** When resuming this project, check .planning/STATE.md for current phase and .planning/ROADMAP.md for phase goals.

**Pre-Phase 1 blockers:**
- Verify Goodreads RSS: https://www.goodreads.com/review/list_rss/79499864?shelf=read
- Google Books API key needed (Google Cloud project)
