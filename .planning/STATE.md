# Project State: Flo's Library

## Project Reference

See: .planning/PROJECT.md (updated 2026-04-01)

**Core value:** A beautiful, always-up-to-date view of every book Florian has read — synced from Goodreads without manual effort.
**Current focus:** Ready to plan Phase 1

---

## Current Position

- **Phase:** 0 of 5 — pre-execution (roadmap complete, no phases started)
- **Plan:** —
- **Status:** Ready to plan Phase 1

**Progress:** `░░░░░░░░░░` 0%

---

## Recent Decisions

| Decision | Rationale |
|----------|-----------|
| Go backend + React TypeScript frontend (separate services in dev) | User preference |
| Go embeds React build in production (single binary) | Simplicity — no Nginx, no CORS |
| PostgreSQL from day one (not SQLite) | SQL dialect divergence causes migration pain |
| golang-migrate for schema migrations | Most adopted, simple mental model |
| Chi v5 as Go HTTP router | Idiomatic net/http, composable middleware |
| sqlc + pgx/v5 for DB layer | Type-safe generated code, no ORM overhead |
| Goodreads RSS (primary) + Google Books + OpenLibrary (fallbacks) | Only viable post-API-shutdown approach |
| Self-hosted book cover images | Reliability — hotlinked covers break |
| TanStack Query v5 + Intersection Observer for infinite scroll | Standard 2025 pattern |
| Cursor-based pagination (not offset) | Stable across inserts |
| Modernized sidebar animation (CSS keyframes, no Lottie) | 60KB bundle cost not justified |

---

## Pending Human Actions

- [ ] **REQUIRED before Phase 1:** Verify Goodreads RSS feed is live — open `https://www.goodreads.com/review/list_rss/79499864?shelf=read` in browser. If XML = RSS works. If 403 = use CSV import path.
- [ ] **REQUIRED before Phase 1 Plan 1.3:** Create Google Cloud project + enable Books API, save API key.

---

## Pending Todos

(None)

---

## Session Continuity

Last session: 2026-04-01
Stopped at: new-project initialization complete — roadmap written
Resume file: (none — handoff superseded by completed initialization)

---

*State initialized: 2026-04-01*
