---
phase: 02-go-rest-api
status: complete
plans_executed: [02-01, 02-02]
---

# Phase 02 — Go REST API: Complete

## What Was Built

### Plan 02-01: Public Read API
- 8 GET endpoints: `/api/books`, `/api/books/currently-reading`, `/api/books/:slug`, `/api/authors`, `/api/authors/:slug`, `/api/genres`, `/api/genres/:slug`, `/api/years`
- Cursor-based keyset pagination (stable across inserts)
- `BookStore` interface + 10 unit tests (all green)
- `unmarshalRefs[T]` generic for pgx `interface{}` → typed slice (authors/genres from `json_agg`)

### Plan 02-02: Cover Serving, SPA, OG Tags
- `frontend/embed.go` embeds `dist/` into the binary (workaround: `//go:embed` can't traverse `..`)
- `/covers/:filename` handler with `Cache-Control: immutable` headers
- SPA catch-all at `/*` — injects OG meta tags per-request for `/books/:slug`
- `.env.example` with all 5 env vars documented

## Post-Execution Fixes (Session)

| Fix | Root Cause |
|-----|-----------|
| `filepath.Base()` for cover path (OG tags) | Windows stores `data\covers\...` with backslashes |
| CSV importer: link authors via `book_authors` | Authors were imported but never joined |
| Enricher: link authors + genres from Google Books | `vi.Authors`/`vi.Categories` fetched but not stored |
| Cover path in API: `/covers/abc.jpg` not raw DB path | Frontend needs usable URL, not `data\covers\...` |
| Confidence gate: bidirectional normalized title check | Old check: `returnedTitle.contains(fullGoodreadsTitle)` always failed for subtitles |
| Confidence gate: `normalizeForCompare()` strips hyphens | "GAME-OVER" vs "Game Over" |
| Confidence gate: series name+number check | "NYPD Red 8" matches "The 11:59 Bomber (NYPD Red, #8)" |
| `normalizeAuthor()` collapses whitespace | Goodreads RSS: `"James  Patterson"` (double space) |
| RSS importer: link authors to `book_authors` | Author extracted from feed but never persisted to join table |
| Title search: fixed double URL-encoding | `url.QueryEscape` called twice; queries were garbled |
| Title search: `inauthor:` added when author known | Improves precision for common titles |
| Title search: `maxResults=5`, loop through results | First result wrong → now tries up to 5 |
| `make dev`: auto-loads `.env` via `include` + `export` | Needed manual env var setup before |
| `make dev`: runs migrations before server | Ensures schema always up to date |
| README: full setup + everyday usage guide | No prior documentation |
| Confidence gate log: shows actual failure reason | Always said "title mismatch" even for author failures |
| 18 confidence gate regression tests | Covers all real-world failure cases from logs |

## Test Status
- `internal/api`: 10 tests ✓
- `internal/sync`: all tests ✓ (18 confidence gate cases)
- `internal/scheduler`: all tests ✓

## Key Decisions
- `inauthor:` in search query uses primary author from `book_authors` (populated by CSV/RSS import)
- Genres only come from Google Books enrichment — Goodreads RSS/CSV don't provide genre data
- ISBN lookup (primary) → title+author search with 5-result loop (fallback)
- All API cover paths returned as `/covers/{filename}` (relative URL, frontend-safe)
