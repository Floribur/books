# Phase 5: Polish & Production Deployment — Research

**Researched:** 2026-04-03
**Domain:** Lottie animation (React), CSS UI polish, Go embed build pipeline, Windows cross-compilation Makefile, SVG favicon
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Use **Lottie** (`lottie-react`) for the sidebar animation — NOT CSS keyframes. Original "Out of Scope" entry in REQUIREMENTS.md is overridden by CONTEXT.md.
- **D-02:** Source a **free Lottie animation from LottieFiles** — a stylized person sitting and reading a book. Reference URLs:
  - `https://lottiefiles.com/animation/reading-book-animation_11806166`
  - `https://lottiefiles.com/free-animation/boy-reading-a-book-f2fyScGybx`
- **D-03:** Animation style: **monochrome silhouette, single color, brand red (#6d233e)**. One fill color only.
- **D-04:** Animation placement: **replaces the "Flo's Library" text title** in the `sidebar-header` div. Text title removed. Mobile: animation at top of drawer (agent may keep text in topbar `<Link>`).
- **D-05:** `prefers-reduced-motion`: **pause or stop** the animation. Static first frame shown.
- **D-06:** Bundle the `.json` file in `frontend/src/assets/`. Use `lottie-react` (wraps `lottie-web`). No CDN or runtime fetch.
- **D-07:** `make build` = `npm run build` in `frontend/` → `go build -o flos-library ./cmd/server`.
- **D-08:** Existing `//go:embed dist` stub in `frontend/embed.go` is already active; activate by populating `frontend/dist/` before Go build.
- **D-09:** Makefile target must handle cross-directory build sequence correctly on **Windows PowerShell**.
- **D-10:** Env var audit: confirm no `GOOGLE_BOOKS_API_KEY` in built JS bundle.
- **D-11:** Write `docs/deployment.md` runbook for Raspberry Pi (arm64). Documentation only — no live deployment.
- **D-12:** Brand consistency pass: `#6d233e` primary, `#c4843a` accent, WCAG AA contrast, CSS token scale, Playfair Display + Inter.
- **D-13:** Favicon: SVG in brand red (`#6d233e`) — reader silhouette. Include `favicon.svg` + `favicon.ico`. Fall back to open-book icon if silhouette not feasible from Lottie asset.
- **D-14:** Page titles: Home = `Flo's Library`; others = `Flo's Library — [Page Name]`. Agent decides implementation (`document.title` in `useEffect` vs shared hook).

### Agent's Discretion
- Which specific free Lottie JSON to use
- How to recolor to brand red (JSON post-processing vs `colorFilters` prop vs choosing single-fill animation)
- Whether to show animation in mobile drawer header or keep text title on mobile
- Page title implementation (`document.title` in `useEffect` per page vs shared `usePageTitle` hook)
- Exact cross-compile Makefile target syntax (arm64 vs armv7)
- Lottie fill color on dark sidebar — agent decides between `#e8c4cf` (dark-theme primary) or `#f0e6e9` (light text)
- Animation sizing inside sidebar (width ≤ 240px)

### Deferred Ideas (OUT OF SCOPE)
- Live VPS/Raspberry Pi deployment
- SSL certificate automation (certbot/Caddy)
- systemd unit wiring (documented in runbook but NOT deployed)
- Open Graph meta tags (already in Go server; no changes needed)
- WebP conversion for covers (V2-04)
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| UI-04 | Sidebar navigation with animated book-reading SVG | `lottie-react` integration in `Sidebar.tsx`; replaces `.sidebar-title` span |
| UI-05 | `prefers-reduced-motion` respected — animation disabled when set | CSS `animation-play-state: paused` on Lottie wrapper + `matchMedia` hook |
| DEPL-01 | Production Go binary embeds React build — single process, no Nginx | `//go:embed dist` already stubbed; `make build` activates it |
| DEPL-02 | systemd service file for VPS deployment | Documented in `docs/deployment.md` runbook only |
| DEPL-03 | SSL via certbot or Caddy | Documented in `docs/deployment.md` (Caddy preferred for zero-config) |
| DEPL-04 | Google Books API key stored as Go backend env var only (never in React/browser) | Vite ONLY exposes `VITE_*` prefixed vars; key is safe; audit confirms this |
</phase_requirements>

---

## Summary

Phase 5 has two tracks: (1) Lottie animation in the sidebar with UI polish, and (2) wiring the production build pipeline. Both are mechanically straightforward given the existing codebase.

**Track 1 — Lottie animation:** `lottie-react@2.4.1` is compatible with React 19 and bundles with standard Vite JSON import. The animation JSON file goes in `frontend/src/assets/`. The sidebar's `sidebar-header` div already exists; replace the `<span className="sidebar-title">` with a `<Lottie>` component. `prefers-reduced-motion` is handled by reading `matchMedia` and passing the result to Lottie's `lottieRef.stop()` / `lottieRef.play()`. Recoloring to brand red is done by post-processing the downloaded JSON (replace all color arrays with the target hex converted to 0-1 normalized RGB) — simpler and more reliable than runtime `colorFilters`.

**Track 2 — Build pipeline:** The `//go:embed dist` in `frontend/embed.go` already compiles; `frontend/dist/` is already populated (from Phase 4 work). The critical issue is the `make build` Makefile target: it currently only runs `go build`. Adding `cd frontend && npm run build` before the Go build activates the embed with fresh JS. On Windows, GNU Make (via Git Bash or MSYS2) handles `cd frontend && npm run build` correctly. Cross-compilation for Raspberry Pi (arm64) requires `CGO_ENABLED=0 GOOS=linux GOARCH=arm64`. The `docs/deployment.md` runbook is documentation only.

**Primary recommendation:** Use `lottie-react` with post-processed JSON (single brand-color fill), implement page titles via a shared `usePageTitle(page?)` hook, fix the Makefile `make build` target to run `npm run build` first, replace the current purple favicon SVG with a brand-red reader silhouette SVG.

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `lottie-react` | `2.4.1` | React wrapper for Lottie animations | Decided in D-06; React 19 compatible; 62KB (lottie-web); JSON format |
| `lottie-web` | `5.13.0` | Underlying animation engine (peer dep) | Auto-installed with lottie-react |

### Supporting (already installed)
| Library | Version | Purpose |
|---------|---------|---------|
| `vitest` | `^4.1.2` | Unit testing Lottie component, usePageTitle hook |
| `@testing-library/react` | `^16.3.2` | Component rendering in tests |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| `lottie-react` | `@lottiefiles/react-lottie-player@3.6.0` | Supports .lottie format but heavier; D-06 locks lottie-react |
| `lottie-react` | `@lottiefiles/dotlottie-react@0.18.9` | WebAssembly renderer, newest, but overkill for a single decoration |
| Post-process JSON colors | Runtime `colorFilters` prop | `colorFilters` requires keypath matching per layer; post-processing is simpler for single-fill animations |
| `usePageTitle` hook | React Helmet / react-helmet-async | No additional library needed; `document.title` in `useEffect` is sufficient for a non-SEO SPA |

**Installation (new package only):**
```bash
cd frontend && npm install lottie-react
```

**Version verification:** Confirmed `lottie-react@2.4.1` is current as of 2026-04-03.

---

## Architecture Patterns

### Lottie Component in Sidebar

The existing `Sidebar.tsx` `navContent` block renders a `sidebar-header` div with a `<span className="sidebar-title">Flo's Library</span>`. Replace the span with a `<LottieAnimation>` component:

```tsx
// frontend/src/components/LottieAnimation.tsx
import { useEffect, useRef } from 'react';
import Lottie, { type LottieRefCurrentProps } from 'lottie-react';
import animationData from '../assets/reading-animation.json';

export function LottieAnimation() {
  const lottieRef = useRef<LottieRefCurrentProps>(null);

  useEffect(() => {
    const mq = window.matchMedia('(prefers-reduced-motion: reduce)');
    if (mq.matches) {
      lottieRef.current?.stop(); // show first frame
    }
    const onChange = (e: MediaQueryListEvent) => {
      e.matches ? lottieRef.current?.stop() : lottieRef.current?.play();
    };
    mq.addEventListener('change', onChange);
    return () => mq.removeEventListener('change', onChange);
  }, []);

  return (
    <Lottie
      lottieRef={lottieRef}
      animationData={animationData}
      loop={true}
      style={{ width: '100%', maxWidth: 200 }}
      aria-label="Animated reader"
      role="img"
    />
  );
}
```

**Key detail:** `lottie-react` uses a `lottieRef` (not `ref`) prop to expose the Lottie instance. This ref gives access to `.play()`, `.stop()`, `.pause()` methods from `lottie-web`.

### prefers-reduced-motion Pattern

The existing `Sidebar.css` already uses `@media (prefers-reduced-motion: reduce)` for the drawer transition. The Lottie animation needs **both** CSS and JS handling:

```css
/* CSS: fallback — hides animation wrapper when reduced-motion; show static SVG instead */
@media (prefers-reduced-motion: reduce) {
  .sidebar-lottie-wrapper {
    /* lottie-react stops via lottieRef.stop() in JS; CSS is belt-and-suspenders */
  }
}
```

The JS approach (calling `lottieRef.current.stop()`) is the primary mechanism — it shows the first frame statically. CSS `animation-play-state: paused` does NOT work on Lottie canvas-rendered animations; the JS API is the correct method.

### Page Titles — `usePageTitle` Hook

```tsx
// frontend/src/hooks/usePageTitle.ts
import { useEffect } from 'react';

export function usePageTitle(page?: string) {
  useEffect(() => {
    document.title = page ? `Flo's Library — ${page}` : "Flo's Library";
    return () => {
      document.title = "Flo's Library"; // reset on unmount
    };
  }, [page]);
}
```

**Usage in pages:**
```tsx
// BookDetailPage.tsx
usePageTitle(book?.title); // "Flo's Library — Dune"

// AuthorsPage.tsx
usePageTitle('Authors');

// GenresPage.tsx
usePageTitle('Genres');

// ReadingChallengePage.tsx
usePageTitle('Reading Challenge');

// HomePage.tsx
usePageTitle(); // "Flo's Library"
```

### Makefile — Multi-step Build with Windows Compatibility

GNU Make (used by this project on Windows via Git Bash/MSYS2) handles `cd dir && cmd` correctly:

```makefile
.PHONY: build build-pi

# Production build: React first (populates frontend/dist), then Go binary with embed
build:
	cd frontend && npm run build
	go build -o flos-library ./cmd/server

# Cross-compile for Raspberry Pi 4/5 (arm64 / ARMv8)
build-pi:
	cd frontend && npm run build
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o flos-library-linux-arm64 ./cmd/server
```

**Critical:** `CGO_ENABLED=0` is required for cross-compilation. Without it, Go uses cgo which cannot cross-compile on Windows without a cross-toolchain.

**Windows note:** On Windows PowerShell, `CGO_ENABLED=0 GOOS=linux GOARCH=arm64` as a single-line prefix works in GNU Make (which spawns sh). If using plain PowerShell outside Make, use separate `$env:GOOS = 'linux'` assignments.

### Go Embed — Path Relationship

`frontend/embed.go` is in the `frontend` package at the repo root:

```
books/
├── frontend/
│   ├── embed.go        ← //go:embed dist (embeds frontend/dist/ relative to this file)
│   └── dist/           ← populated by `npm run build`
└── cmd/server/main.go  ← imports "flos-library/frontend"; accesses frontend.FS
```

`main.go` does `fs.Sub(frontend.FS, "dist")` to strip the "dist" prefix. This is already working. The only missing piece is ensuring `npm run build` runs BEFORE `go build` in the Makefile.

**Failure mode:** If `frontend/dist/` is empty or missing when `go build` runs, the embed compiles successfully (the directive doesn't fail on empty dirs) but the served `index.html` will be a previous stale build or missing. Always run `npm run build` first.

### Favicon Replacement

The current `frontend/public/favicon.svg` is a purple lightning bolt (wrong brand). Replace with a reader silhouette in brand red. The SVG must link in `frontend/index.html`:

```html
<!-- Add to <head> in frontend/index.html -->
<link rel="icon" type="image/svg+xml" href="/favicon.svg" />
<link rel="icon" type="image/x-icon" href="/favicon.ico" />
```

Vite automatically copies everything from `frontend/public/` into `frontend/dist/` during `npm run build` — no additional configuration needed.

### Lottie JSON Color Post-Processing

Brand red `#6d233e` as normalized RGB: `r=109/255=0.4275`, `g=35/255=0.1373`, `b=62/255=0.2431`

For dark sidebar background (`#15100f`), the fill should be light. `#f0e6e9` (light text color) as normalized RGB: `r=0.9412`, `g=0.9020`, `b=0.9137`.

Post-processing approach (preferred over `colorFilters` prop):
1. Download animation JSON from LottieFiles
2. Open JSON, search for all `"c": {"k": [r, g, b, 1]}` or `"sc"` color fields
3. Replace all color values with the desired fill
4. Save as `frontend/src/assets/reading-animation.json`

This works for single-fill silhouette animations and is done once at setup time, not at runtime.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Animation playback | Custom requestAnimationFrame loop | `lottie-react` | Lottie handles timing, looping, frame-level control |
| Reduced-motion detection | `setInterval` polling or custom event | `window.matchMedia('(prefers-reduced-motion: reduce)')` with `addEventListener('change')` | Native browser API, event-driven |
| Page titles | Route-aware component | `usePageTitle` hook with `useEffect` | 4 lines; no library needed |
| Cross-compile Go | Docker-based cross toolchain | `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build` | Go's native cross-compile is zero-config when CGO disabled |
| Bundle audit for leaked secrets | Manual binary inspection | `grep` / `strings` on built JS files in `frontend/dist/assets/` | Simple and reliable |

---

## Common Pitfalls

### Pitfall 1: Go Build Before npm Build (Stale Embed)
**What goes wrong:** `go build` runs before `npm run build` — binary embeds old frontend or missing `index.html`
**Why it happens:** Makefile `build` target currently only runs `go build`; developer habit
**How to avoid:** Makefile target MUST run `cd frontend && npm run build` as the FIRST step before `go build`
**Warning signs:** Binary serves wrong JS bundle hash; Go compiles fine but app shows previous version

### Pitfall 2: lottieRef vs ref Prop
**What goes wrong:** Passing `ref={lottieRef}` to `<Lottie>` instead of `lottieRef={lottieRef}`
**Why it happens:** React component convention is `ref`; lottie-react uses a custom `lottieRef` prop
**How to avoid:** Always use `lottieRef={lottieRef}` with `useRef<LottieRefCurrentProps>(null)`
**Warning signs:** `lottieRef.current` is null; `.play()` / `.stop()` throw errors

### Pitfall 3: CSS `animation-play-state` on Lottie Canvas
**What goes wrong:** Adding `animation-play-state: paused` in CSS and expecting it to pause Lottie
**Why it happens:** Lottie renders to canvas (or SVG), not CSS animations
**How to avoid:** Use `lottieRef.current.stop()` and `.play()` — the JS API is the only working method
**Warning signs:** `prefers-reduced-motion` test passes but animation still plays

### Pitfall 4: CGO_ENABLED Missing in Cross-Compile
**What goes wrong:** `GOOS=linux GOARCH=arm64 go build` fails with linker errors or silently produces wrong binary
**Why it happens:** cgo requires a native C cross-compiler for the target; not present on Windows by default
**How to avoid:** Always set `CGO_ENABLED=0` before cross-compiling: `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build`
**Warning signs:** `exec: "arm-linux-gnueabihf-gcc": executable file not found in $PATH`

### Pitfall 5: GOOGLE_BOOKS_API_KEY Leakage Check (False Confidence)
**What goes wrong:** Assuming the key is safe without actually checking the bundle
**Why it happens:** Vite's `VITE_*` prefix convention is well-known but devs sometimes add `VITE_` by accident
**How to avoid:** After `npm run build`, run `grep -r "GOOGLE_BOOKS_API_KEY" frontend/dist/` — should return nothing. Also check the `.env` file has no `VITE_GOOGLE_BOOKS_API_KEY` accidentally
**Warning signs:** `grep` returns results in `frontend/dist/assets/*.js`

### Pitfall 6: Lottie JSON Import TypeScript Error
**What goes wrong:** `import animationData from '../assets/reading-animation.json'` throws TS error
**Why it happens:** TypeScript `resolveJsonModule` may not be enabled
**How to avoid:** Verify `tsconfig.app.json` has `"resolveJsonModule": true` (Vite scaffold includes this by default). If missing, add it.
**Warning signs:** TS2732 or TS2307 errors on JSON import

### Pitfall 7: Favicon SVG Not Showing in Built App
**What goes wrong:** Updated `favicon.svg` in `frontend/public/` doesn't appear in built app
**Why it happens:** Missing `<link rel="icon" href="/favicon.svg">` in `frontend/index.html` (currently absent)
**How to avoid:** Add both icon link tags to `index.html` `<head>` — Vite passes them through to `dist/index.html`
**Warning signs:** Browser still shows old favicon (check DevTools > Application > Manifest)

---

## Code Examples

### Lottie Component — Full Pattern (Verified from lottie-react docs)
```tsx
// frontend/src/components/LottieAnimation.tsx
import { useEffect, useRef } from 'react';
import Lottie, { type LottieRefCurrentProps } from 'lottie-react';
import animationData from '../assets/reading-animation.json';

export function LottieAnimation() {
  const lottieRef = useRef<LottieRefCurrentProps>(null);

  useEffect(() => {
    const mq = window.matchMedia('(prefers-reduced-motion: reduce)');
    // Apply on mount
    if (mq.matches) {
      lottieRef.current?.stop();
    }
    // Listen for changes (e.g., user changes OS accessibility settings mid-session)
    const handleChange = (e: MediaQueryListEvent) => {
      e.matches ? lottieRef.current?.stop() : lottieRef.current?.play();
    };
    mq.addEventListener('change', handleChange);
    return () => mq.removeEventListener('change', handleChange);
  }, []);

  return (
    <Lottie
      lottieRef={lottieRef}
      animationData={animationData}
      loop={true}
      style={{ width: '100%', maxWidth: 200, margin: '0 auto' }}
      aria-label="Animated reading character"
      role="img"
    />
  );
}
```

### Sidebar Header Replacement
```tsx
// In Sidebar.tsx navContent, replace:
// <span className="sidebar-title">Flo's Library</span>
// with:
<div className="sidebar-header">
  <LottieAnimation />
</div>
```

### usePageTitle Hook
```tsx
// frontend/src/hooks/usePageTitle.ts
import { useEffect } from 'react';

export function usePageTitle(page?: string) {
  useEffect(() => {
    document.title = page ? `Flo's Library \u2014 ${page}` : "Flo's Library";
    return () => { document.title = "Flo's Library"; };
  }, [page]);
}

// Usage: usePageTitle('Authors') → "Flo's Library — Authors"
// Usage: usePageTitle()         → "Flo's Library"
```

### Updated Makefile
```makefile
.PHONY: dev migrate migrate-down sqlc build build-pi

# ... existing targets ...

# Production build: frontend first (populates embed), then Go binary
build:
	cd frontend && npm run build
	go build -o flos-library ./cmd/server

# Raspberry Pi 4/5 cross-compile (ARMv8 / arm64)
build-pi:
	cd frontend && npm run build
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o flos-library-linux-arm64 ./cmd/server
```

### Bundle Audit Command
```bash
# After npm run build — should return no output (key not present)
grep -r "GOOGLE_BOOKS_API_KEY" frontend/dist/

# Also verify no VITE_ prefixed version was accidentally added
grep -r "VITE_GOOGLE" frontend/dist/
```

### favicon.svg Link Tags (index.html)
```html
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Flo's Library</title>
  <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
  <link rel="icon" type="image/x-icon" href="/favicon.ico" />
  <!-- ... rest of head ... -->
</head>
```

---

## Existing Code State (Critical Context)

| File | Current State | Phase 5 Change |
|------|---------------|----------------|
| `frontend/embed.go` | `//go:embed dist` stub — already compiles | No change; just needs `dist/` populated |
| `frontend/dist/` | Populated (from Phase 4 build, 2026-04-03) | Refresh via `npm run build` in `make build` |
| `Makefile` `build` target | `go build -o flos-library ./cmd/server` only | Add `cd frontend && npm run build` as first step |
| `Sidebar.tsx` | `<span className="sidebar-title">Flo's Library</span>` in `sidebar-header` div | Replace with `<LottieAnimation />` |
| `Sidebar.css` | Has `@media (prefers-reduced-motion: reduce)` for drawer transition | Add `.sidebar-lottie-wrapper` reduced-motion rule |
| `frontend/public/favicon.svg` | Purple lightning bolt (wrong brand) | Replace with brand-red reader silhouette |
| `frontend/index.html` | No `<link rel="icon">` tag | Add favicon link tags |
| Page components | No `document.title` management | Add `usePageTitle()` hook to each page |
| `.env.example` | Has `GOOGLE_BOOKS_API_KEY=...` (not VITE_ prefixed) | No change; already safe |

---

## Runtime State Inventory

> This is NOT a rename/refactor phase. Runtime state inventory is not applicable. Skipped.

---

## Environment Availability

| Dependency | Required By | Available | Version | Notes |
|------------|------------|-----------|---------|-------|
| Node.js | `npm run build` | ✓ | v24.14.1 | Confirmed |
| npm | frontend build | ✓ | 10.1.0 | Confirmed |
| Go | `go build` | ✓ | 1.26.1 | Confirmed |
| GNU Make | `make build` | ✓ | (existing Makefile works) | Confirmed (project uses it) |
| `lottie-react` | Lottie animation | ✗ | — | Install: `npm install lottie-react` |

**Missing dependencies with no fallback:**
- `lottie-react` — must be installed; Wave 0 task

**Missing dependencies with fallback:**
- None

**Cross-compile note:** `CGO_ENABLED=0 GOOS=linux GOARCH=arm64` requires no additional toolchain on Windows — Go's built-in cross-compiler handles this when CGO is disabled.

---

## Validation Architecture

> `nyquist_validation: true` in `.planning/config.json` — section required.

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Vitest 4.1.2 |
| Config file | `frontend/vitest.config.ts` |
| Quick run command | `cd frontend && npm test -- --run --reporter=verbose` |
| Full suite command | `cd frontend && npm test -- --run --coverage` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| UI-04 | `LottieAnimation` component renders without error | Unit | `cd frontend && npm test -- --run LottieAnimation` | ❌ Wave 0 |
| UI-04 | Lottie player appears in Sidebar (replaces text title) | Unit | `cd frontend && npm test -- --run Sidebar` | ✅ (Sidebar.test.tsx needed) |
| UI-05 | `prefers-reduced-motion` pauses Lottie animation | Unit | `cd frontend && npm test -- --run LottieAnimation` | ❌ Wave 0 |
| UI-05 | `matchMedia` change event triggers stop/play | Unit | `cd frontend && npm test -- --run LottieAnimation` | ❌ Wave 0 |
| DEPL-01 | `make build` produces binary that serves embedded frontend | Manual smoke | `./flos-library` + `curl http://localhost:8081/` returns HTML | N/A |
| DEPL-04 | `GOOGLE_BOOKS_API_KEY` not in built JS bundle | Manual audit | `grep -r GOOGLE_BOOKS_API_KEY frontend/dist/` returns empty | N/A |
| UI-04 | `usePageTitle` sets `document.title` correctly | Unit | `cd frontend && npm test -- --run usePageTitle` | ❌ Wave 0 |
| UI-04 | Home page title = "Flo's Library" | Unit | `cd frontend && npm test -- --run HomePage` | ❌ (add to existing) |
| UI-04 | Secondary page title = "Flo's Library — [Name]" | Unit | `cd frontend && npm test -- --run AuthorsPage` | ❌ (add to existing) |

### Sampling Rate
- **Per task commit:** `cd frontend && npm test -- --run --reporter=verbose`
- **Per wave merge:** `cd frontend && npm test -- --run --coverage`
- **Phase gate:** Full suite green before `/gsd-verify-work`

### Wave 0 Gaps
- [ ] `frontend/src/components/LottieAnimation.test.tsx` — covers UI-04, UI-05
- [ ] `frontend/src/hooks/usePageTitle.test.ts` — covers page title behavior

**Mocking requirement for Wave 0 tests:**
- `lottie-react` must be mocked in tests (Lottie doesn't render in jsdom environment):
  ```ts
  vi.mock('lottie-react', () => ({
    default: vi.fn(() => <div data-testid="lottie-animation" />),
  }));
  ```
- The `LottieRefCurrentProps` mock (for `lottieRef.current.stop()` / `.play()`):
  ```ts
  const mockLottieRef = { stop: vi.fn(), play: vi.fn(), pause: vi.fn() };
  vi.mock('lottie-react', () => ({
    default: vi.fn(({ lottieRef }) => {
      lottieRef.current = mockLottieRef;
      return <div data-testid="lottie-animation" />;
    }),
  }));
  ```
- `matchMedia` is already mocked in `frontend/src/test/setup.ts` — extend mock to support `prefers-reduced-motion: reduce` in specific tests

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| CSS keyframe SVG animation (Phase 3 plan) | Lottie JSON animation via `lottie-react` | Phase 5 CONTEXT.md D-01 | ~62KB bundle addition; richer character animation |
| Manual Go cross-compile flags | `CGO_ENABLED=0 GOOS=linux GOARCH=arm64` | Go 1.5+ (stable) | Zero-config cross-compile; no external toolchain |
| React Helmet for page titles | `document.title` in `useEffect` | Current standard for lightweight SPAs | No additional dependency; works in React 19 |

**Deprecated / not applicable:**
- `.lottie` (DotLottie binary format) — NOT used. Use `.json` format for simpler integration with `lottie-react`.
- `@lottiefiles/dotlottie-react` — WebAssembly-based, overkill for a single decorative animation.

---

## Open Questions

1. **Which LottieFiles animation to use**
   - What we know: Two reference URLs provided in CONTEXT.md D-02
   - What's unclear: These URLs require LottieFiles account download; animation must be single-fill silhouette OR post-processable; agent must evaluate at download time
   - Recommendation: Agent downloads both, picks the cleanest one with fewest layers and a single `sc` (stroke color) or `fc` (fill color) field per shape

2. **Lottie animation dark-mode fill color**
   - What we know: Sidebar background is always `#15100f` (dark); brand red `#6d233e` may be too dark to see
   - What's unclear: Agent needs to judge visually — `#e8c4cf` (dark theme primary) or `#f0e6e9` (near-white) reads better
   - Recommendation: Use `#e8c4cf` as it's already the dark-theme primary color and matches the design system

3. **`favicon.ico` generation**
   - What we know: `favicon.ico` is requested alongside `favicon.svg` for legacy browser support
   - What's unclear: No ICO generation tool is currently available; may require manual conversion
   - Recommendation: Use an online converter (e.g., favicon.io) to convert the `favicon.svg` → `favicon.ico` at 32×32 and 16×16 sizes. Add to `frontend/public/` as a static asset. Document this as a manual one-time step.

---

## Sources

### Primary (HIGH confidence)
- Direct code inspection: `frontend/embed.go`, `frontend/src/components/Sidebar.tsx`, `frontend/src/components/Sidebar.css`, `Makefile`, `frontend/package.json`, `frontend/index.html` — verified 2026-04-03
- npm registry: `lottie-react@2.4.1` peer deps (`react: ^16.8 || ^17 || ^18 || ^19`) — verified 2026-04-03
- npm registry: `lottie-web@5.13.0`, `@lottiefiles/react-lottie-player@3.6.0` — verified 2026-04-03

### Secondary (MEDIUM confidence)
- lottie-react documentation (https://lottiereact.com) — `lottieRef` prop, `LottieRefCurrentProps` type, `stop()`/`play()` API
- Go embed documentation — `//go:embed` directive requires directory to exist at compile time; `fs.Sub` usage
- Go cross-compilation documentation — `CGO_ENABLED=0` required for cross-compile without C toolchain

### Tertiary (LOW confidence)
- LottieFiles animation availability — sourcing the specific animation requires agent to download and evaluate at implementation time

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — `lottie-react@2.4.1` verified on npm with React 19 peer dep support
- Architecture: HIGH — based on direct code inspection of all integration points
- Pitfalls: HIGH — based on direct code inspection (missing favicon link, wrong favicon SVG, missing Makefile step)
- Lottie animation sourcing: LOW — requires human action (download + evaluate) at implementation time

**Research date:** 2026-04-03
**Valid until:** 2026-05-03 (stable libraries; Lottie-react API is mature)
