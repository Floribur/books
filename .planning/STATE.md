---
gsd_state_version: 1.0
milestone: v1.0
milestone_name: milestone
status: verifying
stopped_at: Phase 5 context gathered
last_updated: "2026-04-03T20:36:17.471Z"
progress:
  total_phases: 5
  completed_phases: 4
  total_plans: 10
  completed_plans: 11
  percent: 100
---

# Project State: Flo's Library

## Project Reference

See: .planning/PROJECT.md (updated 2026-04-01)

**Core value:** A beautiful, always-up-to-date view of every book Florian has read — synced from Goodreads without manual effort.
**Current focus:** Phase 05 — production-sync

---

## Current Position

Phase: 04 (frontend-pages) — COMPLETE
Plan: 2 of 2 — COMPLETE (human-verified 2026-04-03)

- **Phase:** 4 of 5 (frontend pages) — COMPLETE
- **Status:** All 6 routes wired and human-verified

- **Next:** Phase 05 (production-sync)
- **Plan:** Not started
- **Status:** Ready to begin

**Progress:** [██████████] 100%

---

## Recent Decisions

| Decision | Rationale |
|----------|-----------|
| sql/schema.sql is sqlc source of truth, manually synced with migrations | sqlc doesn't handle BEGIN/COMMIT transaction wrappers |
| books.search_vector TSVECTOR stub included from day one | Avoids schema migration later when full-text search is added |
| emit_pointers_for_null_types in sqlc.yaml | Ensures nullable columns map to *string in Go for type safety |
| Go backend + React TypeScript frontend (separate services in dev) | User preference |
| Go embeds React build in production (single binary) | Simplicity — no Nginx, no CORS |
| PostgreSQL from day one (not SQLite) | SQL dialect divergence causes migration pain |
| golang-migrate for schema migrations | Most adopted, simple mental model |
| Chi v5 as Go HTTP router | Idiomatic net/http, composable middleware |
| sqlc + pgx/v5 for DB layer | Type-safe generated code, no ORM overhead |
| Goodreads RSS fields in item.Custom map, not item.Extensions | Live feed inspection: gofeed routes non-namespaced custom elements to item.Custom |
| Goodreads RSS exposes 'isbn' field (not 'isbn13') | Stored in isbn13 DB column; Plan 01-03 enrichment can normalize further |
| fetchShelfFromURL / startWithTicker: injectable dependencies for testability | TDD pattern: unexported testable core, public wrapper with defaults |
| Goodreads RSS (primary) + Google Books + OpenLibrary (fallbacks) | Only viable post-API-shutdown approach |
| stdlib sync aliased as stdsync inside package sync | Avoids package name collision; behaviour identical |
| Confidence gate uses substring match (not exact) | "Dune Messiah" passes for inputTitle="Dune" — acceptable false positive tradeoff |
| Title-only confidence gate in fallback (no author join) | GetUnenrichedBooks returns books table only; author join would require query change |
| UpdateBookEnrichmentParams uses *string for nullable TEXT | sqlc with pgx/v5 generates *string not pgtype.Text for nullable columns |
| Self-hosted book cover images | Reliability — hotlinked covers break |
| TanStack Query v5 + Intersection Observer for infinite scroll | Standard 2025 pattern |
| Cursor-based pagination (not offset) | Stable across inserts |
| Modernized sidebar animation (CSS keyframes, no Lottie) | 60KB bundle cost not justified |
| BookGrid accepts queryKey+fetchFn props | All filtered grids (author, genre) reuse the same component |
| fetchBooksByAuthor/fetchBooksByGenre extract PaginatedBooks from detail envelope | Caller doesn't need to know the AuthorDetail/GenreDetail shape |
| DescriptionBlock clamps at 640 chars / 8 lines (CHAR_THRESHOLD constant) | Matches D-04 UI spec |
| StatsStrip accepts totalPages: null for page-based stats | Page_count not on list endpoint; null hides stats; Phase 5 can add without interface change |
| fetchFn for BookGrid wrapped in useMemo per selectedYear | Creates stable closure per year; prevents unnecessary re-fetches when other state changes |

---

## Pending Human Actions

- [x] **CONFIRMED in 01-02:** Goodreads RSS feed is live — `https://www.goodreads.com/review/list_rss/79499864?shelf=read` returns 127 books as RSS XML. Fields delivered via item.Custom map.
- [ ] **REQUIRED before Phase 1 Plan 1.3:** Create Google Cloud project + enable Books API, save API key.

---

## Pending Todos

(None)

---

## Session Continuity

Last session: 2026-04-03T20:36:17.467Z
Stopped at: Phase 5 context gathered
Resume file: .planning/phases/05-polish-production-deployment/05-CONTEXT.md
