---
phase: 1
slug: data-pipeline
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-04-01
---

# Phase 1 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go stdlib `testing` package (no test framework needed) |
| **Config file** | none — Wave 0 installs test files |
| **Quick run command** | `go test ./internal/... -short -timeout 30s` |
| **Full suite command** | `go test ./... -timeout 60s` |
| **Estimated runtime** | ~30 seconds (short), ~60 seconds (full with integration) |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/... -short -timeout 30s`
- **After every plan wave:** Run `go test ./... -timeout 60s`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 60 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 1-01-01 | 1.1 | 0 | DATA-05 | integration | `go test ./internal/db/... -run TestMigrations` | ❌ W0 | ⬜ pending |
| 1-02-01 | 1.2 | 0 | SYNC-02 | unit | `go test ./internal/sync/... -run TestRSSPagination -short` | ❌ W0 | ⬜ pending |
| 1-02-02 | 1.2 | 1 | SYNC-03 | unit | `go test ./internal/sync/... -run TestShelfMerge -short` | ❌ W0 | ⬜ pending |
| 1-02-03 | 1.2 | 1 | DATA-02 | unit | `go test ./internal/sync/... -run TestSlugCollision -short` | ❌ W0 | ⬜ pending |
| 1-02-04 | 1.2 | 1 | SYNC-09 | unit | `go test ./internal/sync/... -run TestCSVImport -short` | ❌ W0 | ⬜ pending |
| 1-02-05 | 1.2 | 2 | SYNC-01 | unit | `go test ./internal/scheduler/... -run TestScheduler -short` | ❌ W0 | ⬜ pending |
| 1-03-01 | 1.3 | 1 | SYNC-04 | unit | `go test ./internal/sync/... -run TestEnrichmentConfidenceGate -short` | ❌ W0 | ⬜ pending |
| 1-03-02 | 1.3 | 1 | SYNC-07 | unit | `go test ./internal/sync/... -run TestCoverValidation -short` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/sync/rss_test.go` — stubs for SYNC-02, SYNC-03
- [ ] `internal/sync/enricher_test.go` — stubs for SYNC-04, confidence gate logic
- [ ] `internal/sync/covers_test.go` — stubs for SYNC-07 (needs test fixtures: tiny JPEG, 1×1 JPEG, valid JPEG)
- [ ] `internal/sync/csv_test.go` — stubs for SYNC-09, Goodreads ISBN unquoting
- [ ] `internal/sync/slug_test.go` — stubs for DATA-02, collision strategy
- [ ] `internal/scheduler/scheduler_test.go` — stubs for SYNC-01 (mock ticker via dependency injection)
- [ ] `internal/db/migrations_test.go` — stubs for DATA-05 (integration, requires running PostgreSQL)

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| RSS feed returns XML (not 403) | SYNC-01 | External Goodreads service — not mockable pre-implementation | Open `https://www.goodreads.com/review/list_rss/79499864?shelf=read` in browser; confirm XML response |
| RSS extension namespace key paths correct | SYNC-02 | Runtime data from live feed required | Run pre-task in Plan 1.2: fetch feed, dump `item.Extensions` keys |
| Google Books API key valid and quota available | SYNC-04 | External API credential — cannot be automated | Test with `curl "https://www.googleapis.com/books/v1/volumes?q=isbn:9780385472579&key=$GOOGLE_BOOKS_API_KEY"` |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 60s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
