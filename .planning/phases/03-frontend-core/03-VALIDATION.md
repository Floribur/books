---
phase: 3
slug: frontend-core
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-04-03
---

# Phase 3 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Vitest 4.1.2 |
| **Config file** | `frontend/vitest.config.ts` — Wave 0 creates this |
| **Quick run command** | `cd frontend && npm test -- --run` |
| **Full suite command** | `cd frontend && npm test -- --run --coverage` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `cd frontend && npm test -- --run`
- **After every plan wave:** Run `cd frontend && npm test -- --run --coverage`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 3-01-05 | 01 | 1 | UI-06 | unit | `cd frontend && npm test -- --run src/hooks/useTheme.test.ts` | ❌ W0 | ⬜ pending |
| 3-01-07 | 01 | 1 | UI-09 | unit | `cd frontend && npm test -- --run src/components/BookCover.test.tsx` | ❌ W0 | ⬜ pending |
| 3-02-01 | 02 | 2 | HOME-03, UI-10 | integration | `cd frontend && npm test -- --run src/components/BookGrid.test.tsx` | ❌ W0 | ⬜ pending |
| 3-02-02 | 02 | 2 | HOME-04 | unit | `cd frontend && npm test -- --run src/components/BookCard.test.tsx` | ❌ W0 | ⬜ pending |
| 3-02-03 | 02 | 2 | UI-01, UI-02 | integration | `cd frontend && npm test -- --run src/components/BookGrid.test.tsx` | ❌ W0 | ⬜ pending |
| 3-02-04 | 02 | 2 | UI-01 | integration | `cd frontend && npm test -- --run src/hooks/useIntersectionObserver.test.ts` | ❌ W0 | ⬜ pending |
| 3-02-05 | 02 | 2 | UI-02 | unit | `cd frontend && npm test -- --run src/components/BookGrid.test.tsx` | ❌ W0 | ⬜ pending |
| 3-02-06 | 02 | 2 | UI-03 | unit | `cd frontend && npm test -- --run src/hooks/useScrollRestoration.test.ts` | ❌ W0 | ⬜ pending |
| 3-03-01 | 03 | 3 | HOME-01 | unit | `cd frontend && npm test -- --run src/components/Bio.test.tsx` | ❌ W0 | ⬜ pending |
| 3-03-02 | 03 | 3 | HOME-02 | unit | `cd frontend && npm test -- --run src/components/NowReadingSection.test.tsx` | ❌ W0 | ⬜ pending |
| 3-03-03 | 03 | 3 | HOME-03 | integration | `cd frontend && npm test -- --run src/components/BookGrid.test.tsx` | ❌ W0 | ⬜ pending |
| UI-07 | — | — | UI-07 | manual | visual inspection | N/A | ⬜ pending |
| UI-08 | — | — | UI-08 | manual | visual inspection | N/A | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `frontend/vitest.config.ts` — vitest config with jsdom environment
- [ ] `frontend/src/test/setup.ts` — MSW server setup, IntersectionObserver mock (`global.IntersectionObserver = vi.fn(...)`), matchMedia mock
- [ ] `frontend/src/test/msw-server.ts` — MSW server instance
- [ ] `frontend/src/mocks/handlers.ts` — API handler mocks for `GET /api/books` and `GET /api/books/currently-reading`
- [ ] Install test deps: `npm install -D vitest @testing-library/react @testing-library/user-event msw @vitest/coverage-v8 jsdom`
- [ ] Add `"test": "vitest"` script to `frontend/package.json`
- [ ] Stub test files for all Wave 0 gaps above (❌ W0 entries)

*Wave 0 installs testing infrastructure as part of Plan 3.1.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| CSS custom properties render correct brand colors in light/dark mode | UI-07 | Computed styles not reliably testable in jsdom | Open app in browser, inspect `--color-primary` etc. in DevTools; toggle theme and verify token values switch |
| Playfair Display and Inter fonts applied to headings/body | UI-08 | Font loading is a network/browser concern, not unit-testable | Open app in browser, check heading and body fonts in DevTools Elements panel |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
