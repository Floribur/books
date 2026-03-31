# Flo's Library

## What This Is

A personal book showcase website for Florian — a reader who tracks progress on Goodreads. The site displays currently-reading and read books with cover art, detail pages, author/genre index pages, and a yearly reading challenge view. Book data stays automatically synchronized from Goodreads.

## Core Value

A beautiful, always-up-to-date view of every book Florian has read — synced from Goodreads without manual effort.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Home page with bio section (name, photo, passions)
- [ ] "Now Reading" section showing current books (cover, title, author) with click-through to detail page
- [ ] "Books Read" section in descending read-date order with infinite scroll (first ~20-30 books, load more on scroll)
- [ ] Book detail page: cover, title, author, year published, genres, page count, description
- [ ] Author pages listing all books read by that author (descending order)
- [ ] Genre pages listing all books in that genre (descending order)
- [ ] Reading Challenge page showing books read per year, browseable across past years
- [ ] Automatic sync with Goodreads shelves (currently-reading, read) — mechanism TBD from research
- [ ] Book cover images downloaded and self-hosted (not hotlinked)
- [ ] Book metadata enriched via Google Books API (cover, description, genres, page count, ISBN)
- [ ] Sidebar navigation with modernized animated book-reading graphic
- [ ] Brand color: #6d233e (rgb 109, 35, 62) as primary

### Out of Scope

- User accounts / multi-user support — this is a personal showcase, not a social app
- Writing reviews or ratings — display only, Goodreads stays the source of truth for opinions
- Mobile app — web only
- Goodreads API (deprecated 2020) — need alternative sync mechanism

## Context

**Prior art:** Florian had a working version of this app built ~6 years ago using:
- PHP / Lumen (Laravel micro-framework) backend
- MySQL database
- React frontend
- Goodreads API for shelf sync (now discontinued)
- Google Books API for cover thumbnails (hotlinked, not downloaded)
- Sidebar nav with book-reading animation

**Goodreads sync:** The old Goodreads API is gone. Research phase will determine the best 2025 approach — options include scraping the public profile, using RSS feeds Goodreads still exposes, or a third-party integration. Florian's profile: https://www.goodreads.com/user/show/79499864-florian

**Tech direction:**
- React + TypeScript frontend (separate from backend)
- Go backend (separate service, not embedded single binary)
- Start with in-memory/embedded DB, target PostgreSQL for production
- Self-host book cover images (download via Google Books API, store locally)

## Constraints

- **Tech Stack**: Go backend + React TypeScript frontend — user's preference for the rebuild
- **Data Source**: Goodreads as source of truth for shelf state — no manual entry
- **Database**: Start in-memory (dev/prototype), migrate to PostgreSQL for production
- **Covers**: Must self-host images — no hotlinking to external URLs
- **No Goodreads API**: Must find alternative sync method (RSS, scraping, or export)

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Go backend | Modern, single binary deployment, good for API servers | — Pending |
| Separate frontend/backend | Standard dev experience, easier local iteration | — Pending |
| Self-host book covers | Reliability — hotlinked covers break when source changes | — Pending |
| PostgreSQL as DB target | Solid relational DB for personal VPS hosting | — Pending |
| Modernize sidebar animation | Keep the branded feel, update to current CSS/React patterns | — Pending |
| Goodreads sync mechanism | TBD from research — API gone, alternatives needed | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd:transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-03-31 after initialization*
