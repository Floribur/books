---
phase: 5
slug: polish-production-deployment
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-04-03
---

# Phase 5 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Vitest 4.1.2 |
| **Config file** | `frontend/vitest.config.ts` |
| **Quick run command** | `cd frontend && npm test -- --run --reporter=verbose` |
| **Full suite command** | `cd frontend && npm test -- --run --coverage` |
| **Estimated runtime** | ~15 seconds |

---

## Sampling Rate

- **After every task commit:** Run `cd frontend && npm test -- --run --reporter=verbose`
- **After every plan wave:** Run `cd frontend && npm test -- --run --coverage`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** ~15 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 5-01-01 | 01 | 0 | UI-04 | unit | `cd frontend && npm test -- --run LottieAnimation` | ❌ W0 | ⬜ pending |
| 5-01-02 | 01 | 0 | UI-05 | unit | `cd frontend && npm test -- --run LottieAnimation` | ❌ W0 | ⬜ pending |
| 5-01-03 | 01 | 0 | UI-04 | unit | `cd frontend && npm test -- --run usePageTitle` | ❌ W0 | ⬜ pending |
| 5-01-04 | 01 | 1 | UI-04 | unit | `cd frontend && npm test -- --run Sidebar` | ✅ (modify) | ⬜ pending |
| 5-01-05 | 01 | 1 | UI-05 | unit | `cd frontend && npm test -- --run LottieAnimation` | ❌ W0 | ⬜ pending |
| 5-01-06 | 01 | 2 | UI-04 | unit | `cd frontend && npm test -- --run` | ✅ (modify) | ⬜ pending |
| 5-01-07 | 01 | 2 | UI-04 | manual | see Manual-Only section | N/A | ⬜ pending |
| 5-02-01 | 02 | 1 | DEPL-01 | manual smoke | `./flos-library` + `curl http://localhost:8081/` | N/A | ⬜ pending |
| 5-02-02 | 02 | 1 | DEPL-04 | manual audit | `grep -r GOOGLE_BOOKS_API_KEY frontend/dist/` | N/A | ⬜ pending |
| 5-02-03 | 02 | 2 | DEPL-02 | manual | Read `docs/deployment.md` contains systemd section | N/A | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `frontend/src/components/LottieAnimation.test.tsx` — stubs for UI-04, UI-05 (lottie-react mocked, matchMedia mock extended for `prefers-reduced-motion: reduce`)
- [ ] `frontend/src/hooks/usePageTitle.test.ts` — covers page title hook behavior (home = "Flo's Library", secondary = "Flo's Library — [Name]")

**Mocking strategy (Wave 0):**
```ts
// lottie-react mock (in LottieAnimation.test.tsx)
const mockLottieRef = { stop: vi.fn(), play: vi.fn(), pause: vi.fn() };
vi.mock('lottie-react', () => ({
  default: vi.fn(({ lottieRef }) => {
    lottieRef.current = mockLottieRef;
    return <div data-testid="lottie-animation" />;
  }),
}));
// matchMedia mock (in setup.ts — already mocked; extend for prefers-reduced-motion)
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: query === '(prefers-reduced-motion: reduce)',
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
  })),
});
```

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| `make build` produces binary with embedded frontend | DEPL-01 | Requires full build pipeline + running binary | Run `make build`; start `./flos-library`; visit `http://localhost:8081/` — should serve React app without separate frontend server |
| `GOOGLE_BOOKS_API_KEY` not in built JS bundle | DEPL-04 | Requires built assets | Run `grep -r GOOGLE_BOOKS_API_KEY frontend/dist/` — should return no output |
| Favicon shows brand-red reader in browser | UI-04, D-13 | Requires browser rendering | Open app in browser; check browser tab for favicon |
| Animation pauses with `prefers-reduced-motion` OS setting | UI-05 | Requires OS accessibility setting change | Enable "Reduce motion" in OS; reload app; verify animation is static |
| `docs/deployment.md` covers Pi deployment steps | DEPL-02, DEPL-03 | Documentation review | Read file for: cross-compile command, scp step, systemd unit example, Caddy config |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 15s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
