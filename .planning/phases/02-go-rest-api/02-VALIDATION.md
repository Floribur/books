---
phase: 2
slug: go-rest-api
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-04-03
---

# Phase 2 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go stdlib `testing` package |
| **Config file** | none — Go test discovery is automatic |
| **Quick run command** | `go test ./internal/api/... -v -run TestPublic` |
| **Full suite command** | `go test ./... -count=1` |
| **Estimated runtime** | ~15 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/api/... -v -run TestPublic`
- **After every plan wave:** Run `go test ./... -count=1`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 15 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 2-01-01 | 01 | 0 | API-01–10 | unit setup | `go test ./internal/api/... -v` | ❌ W0 | ⬜ pending |
| 2-01-02 | 01 | 1 | API-01 | unit (httptest) | `go test ./internal/api/... -run TestGetBooks` | ❌ W0 | ⬜ pending |
| 2-01-03 | 01 | 1 | API-01 | unit | `go test ./internal/api/... -run TestCursor` | ❌ W0 | ⬜ pending |
| 2-01-04 | 01 | 1 | API-02 | unit (httptest) | `go test ./internal/api/... -run TestGetCurrentlyReading` | ❌ W0 | ⬜ pending |
| 2-01-05 | 01 | 1 | API-03 | unit (httptest) | `go test ./internal/api/... -run TestGetBookBySlug` | ❌ W0 | ⬜ pending |
| 2-01-06 | 01 | 1 | API-04 | unit (httptest) | `go test ./internal/api/... -run TestGetAuthors` | ❌ W0 | ⬜ pending |
| 2-01-07 | 01 | 1 | API-05 | unit (httptest) | `go test ./internal/api/... -run TestGetAuthorBySlug` | ❌ W0 | ⬜ pending |
| 2-01-08 | 01 | 1 | API-06 | unit (httptest) | `go test ./internal/api/... -run TestGetGenres` | ❌ W0 | ⬜ pending |
| 2-01-09 | 01 | 1 | API-07 | unit (httptest) | `go test ./internal/api/... -run TestGetGenreBySlug` | ❌ W0 | ⬜ pending |
| 2-01-10 | 01 | 1 | API-08 | unit (httptest) | `go test ./internal/api/... -run TestGetYears` | ❌ W0 | ⬜ pending |
| 2-02-01 | 02 | 1 | API-09 | unit (httptest) | `go test ./internal/api/... -run TestCoversCache` | ❌ W0 | ⬜ pending |
| 2-02-02 | 02 | 1 | API-10 | unit (httptest) | `go test ./internal/api/... -run TestOGInjection` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/api/store.go` — `BookStore` interface wrapping `*db.Queries` methods for mock injection
- [ ] `internal/api/public_test.go` — test stubs for API-01 through API-10 (httptest-based, mock store)

*Wave 0 must complete before Wave 1 tasks begin.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| OG tags render correctly in browser head | API-10 | Requires browser/curl inspection | `curl -s http://localhost:8081/books/<slug> \| grep "og:"` — verify 5 meta tags present |
| Cover file served with correct cache headers | API-09 | Headers inspection | `curl -I http://localhost:8081/covers/<file.jpg>` — verify `Cache-Control: public, max-age=31536000, immutable` |
| SPA catch-all serves index.html for unknown routes | API-10 | Routing behavior | `curl -s http://localhost:8081/some/unknown/path` — verify returns `index.html` content |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 15s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
