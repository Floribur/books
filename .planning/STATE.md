---
gsd_state_version: 1.0
milestone: v1.0
milestone_name: milestone
status: executing
stopped_at: Completed 01-02-PLAN.md
last_updated: "2026-04-02T12:20:00.000Z"
progress:
  total_phases: 5
  completed_phases: 0
  total_plans: 3
  completed_plans: 2
  percent: 67
---

# Project State: Flo's Library

## Project Reference

See: .planning/PROJECT.md (updated 2026-04-01)

**Core value:** A beautiful, always-up-to-date view of every book Florian has read — synced from Goodreads without manual effort.
**Current focus:** Phase 01 — data-pipeline

---

## Current Position

Phase: 01 (data-pipeline) — EXECUTING
Plan: 3 of 3

- **Phase:** 1 of 5 — Phase 01 in progress
- **Plan:** 01-02 complete, proceeding to 01-03
- **Status:** Executing Phase 01

**Progress:** [██████░░░░] 67%

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
| Self-hosted book cover images | Reliability — hotlinked covers break |
| TanStack Query v5 + Intersection Observer for infinite scroll | Standard 2025 pattern |
| Cursor-based pagination (not offset) | Stable across inserts |
| Modernized sidebar animation (CSS keyframes, no Lottie) | 60KB bundle cost not justified |

---

## Pending Human Actions

- [x] **CONFIRMED in 01-02:** Goodreads RSS feed is live — `https://www.goodreads.com/review/list_rss/79499864?shelf=read` returns 127 books as RSS XML. Fields delivered via item.Custom map.
- [ ] **REQUIRED before Phase 1 Plan 1.3:** Create Google Cloud project + enable Books API, save API key.

---

## Pending Todos

(None)

---

## Session Continuity

Last session: 2026-04-02T12:20:00.000Z
Stopped at: Completed 01-02-PLAN.md
Resume file: None

---

*State initialized: 2026-04-01*
