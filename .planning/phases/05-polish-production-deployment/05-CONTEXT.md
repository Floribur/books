# Phase 5: Polish & Production Deployment - Context

**Gathered:** 2026-04-03
**Status:** Ready for planning

<domain>
## Phase Boundary

Complete the sidebar animation (Lottie, sourced from LottieFiles), do a final UI polish pass (brand consistency, favicon, page titles), and wire up the local production build pipeline (`make build` = React build + Go binary with embedded frontend). No live VPS deployment this phase — instead, produce a working single-binary local build and a written Raspberry Pi deployment runbook for when deployment is needed.

</domain>

<decisions>
## Implementation Decisions

### Sidebar Animation

- **D-01:** Use **Lottie** for the sidebar animation — this overrides the Phase 3 "CSS keyframes only" decision. The original concern was 60KB bundle cost for a simple decorative element; the updated vision (an animated reading character) is not achievable with raw CSS keyframes.
- **D-02:** Source a **free Lottie animation from LottieFiles** — a stylized person sitting and reading a book, turning pages periodically. Close references to the desired look:
  - `https://lottiefiles.com/animation/reading-book-animation_11806166`
  - `https://lottiefiles.com/free-animation/boy-reading-a-book-f2fyScGybx`
  Find the closest match: simple, upper-body or full sitting figure, clean vector style. The character does not need to be a photorealistic likeness.
- **D-03:** Animation style: **monochrome silhouette, single color, brand red (#6d233e)**. Keep it dead simple — one fill color only, no multicolor.
- **D-04:** Animation placement: **replaces the "Flo's Library" text title** — the Lottie player IS the sidebar header. The text title is removed. On mobile (drawer), the animation appears at the top of the drawer as well.
- **D-05:** `prefers-reduced-motion`: **pause or stop** the animation when the media query matches. Static first frame shown instead.
- **D-06:** Bundle approach: use `lottie-react` (wraps `lottie-web`). The `.lottie` or `.json` file is bundled with the app (stored in `frontend/src/assets/` or similar). No CDN or runtime fetch needed.

### Production Build Pipeline

- **D-07:** `make build` should: (1) run `npm run build` in `frontend/` to produce `frontend/dist/`, (2) then run `go build -o flos-library ./cmd/server` which picks up the `//go:embed frontend/dist` directive. Output: a single `flos-library` binary that serves the full app.
- **D-08:** The existing `//go:embed frontend/dist` directive in `frontend/embed.go` is already stubbed — Phase 5 activates it by ensuring the dist directory is populated before the Go build.
- **D-09:** Add a `make build-local` (or update `make build`) that handles the cross-directory build sequence correctly on Windows PowerShell (since the project runs on Windows dev).
- **D-10:** Environment variable audit: verify no `GOOGLE_BOOKS_API_KEY` leaks into the React/JS bundle (it should live only in Go backend `.env`).

### Deployment Runbook (non-automated)

- **D-11:** Write a `docs/deployment.md` runbook describing Raspberry Pi deployment:
  - Makefile target to cross-compile the Go binary for Linux ARM (`GOOS=linux GOARCH=arm64`)
  - Copy binary + `.env` to Pi via `scp`
  - systemd unit file example (included in the doc, not wired up)
  - Caddy configuration for reverse proxy + automatic HTTPS
  - `systemctl restart flos-library`
  This is documentation only — no CI/CD, no automation in this phase.

### UI Polish

- **D-12:** General consistency pass across all pages: verify brand colors (#6d233e primary, #c4843a accent), contrast ratios (WCAG AA minimum), spacing (CSS token scale), and typography (Playfair Display headings, Inter body) are applied consistently everywhere.
- **D-13:** **Favicon**: SVG favicon in brand red (#6d233e) — a small reader silhouette (matching the Lottie character style). If a clean SVG favicon is not feasible from the Lottie file, fall back to a simple open book icon. Include both `favicon.svg` and `favicon.ico` (for IE/legacy).
- **D-14:** **Page titles** (via React `<title>` or a `useEffect`/`document.title` pattern):
  - Home page: `Flo's Library`
  - Secondary pages: `Flo's Library — [Page Name]` (e.g. `Flo's Library — Dune`, `Flo's Library — Authors`, `Flo's Library — Reading Challenge`)
  Agent decides the implementation (React Helmet, or simple `document.title` in each page component).

### Agent's Discretion

- Which specific free Lottie animation to use (pick the cleanest/simplest reading figure from LottieFiles that can be recolored to a single fill)
- How to recolor the Lottie to brand red (Lottie's `colorFilters` prop, or post-process the JSON, or choose an animation that has a single-layer fill)
- Whether to use `lottie-react` or `@lottiefiles/react-lottie-player` (both are viable wrappers)
- Animation sizing in sidebar (width constrained to sidebar width ~240px, height proportional)
- Mobile: whether to show the animation in the mobile drawer header or just show the text title on mobile
- Page title implementation (`document.title` in `useEffect` per page vs a shared `usePageTitle` hook)
- Exact cross-compile Makefile target syntax for Pi (arm64 vs armv7)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Requirements
- `.planning/REQUIREMENTS.md` §UI/UX (UI-04, UI-05) — animated sidebar, prefers-reduced-motion
- `.planning/REQUIREMENTS.md` §Deployment (DEPL-01–DEPL-04) — single binary embed, env var safety

### Roadmap
- `.planning/ROADMAP.md` §Phase 5 — Plan 5.1 and 5.2 task breakdown

### Existing Sidebar Code
- `frontend/src/components/Sidebar.tsx` — current sidebar structure; animation replaces the `sidebar-header` div's text
- `frontend/src/components/Sidebar.css` — existing styles; already has `@media (prefers-reduced-motion)` pattern to follow

### Design System
- `frontend/src/styles/tokens.css` — CSS custom properties (spacing, typography scale)
- `frontend/src/styles/themes.css` — dark theme: `--color-background: #15100f`, `--color-primary: #e8c4cf` (dark mode lightened primary)
- `frontend/src/styles/tokens.css` `:root` — light theme: `--color-primary: #6d233e`, `--color-accent: #c4843a`

### Phase 3 Patterns (carried forward)
- `.planning/phases/03-frontend-core/03-CONTEXT.md` — established dark/light mode, sidebar placement (D-11), CSS token usage

### Go Embed Stub
- `frontend/embed.go` — already has `//go:embed frontend/dist` stub; needs `frontend/dist/` populated before Go build

### Lottie Animation References (user-provided)
- `https://lottiefiles.com/animation/reading-book-animation_11806166` — reference style
- `https://lottiefiles.com/free-animation/boy-reading-a-book-f2fyScGybx` — reference style

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `Sidebar.tsx` — fully built desktop + mobile drawer; animation drops into the `sidebar-header` div (replace `<span className="sidebar-title">Flo's Library</span>` with Lottie player)
- `Sidebar.css` — `@media (prefers-reduced-motion: reduce)` already used for drawer transition — same pattern for animation pause
- `frontend/src/styles/` — complete design system; all colors/tokens available as CSS custom properties

### Established Patterns
- CSS custom properties for all colors, spacing, and typography — no hardcoded values in new CSS
- `@media (prefers-reduced-motion: reduce)` guard (in Sidebar.css) — replicate for animation

### Integration Points
- `frontend/embed.go` — Go embed stub waiting for `frontend/dist/` to be populated
- `Makefile` — `make build` currently only runs `go build`; needs `npm run build` step prepended
- `frontend/package.json` — `"build": "tsc -b && vite build"` outputs to `frontend/dist/`

</code_context>

<specifics>
## Specific Ideas

- The animation replacing the title means the sidebar loses its text wordmark. This is intentional — the animated character IS the brand identity for the sidebar. On mobile, agent may choose to keep the text in the topbar `<Link>` or also use a static frame of the animation.
- Lottie color: target `#6d233e` (light theme primary) as the fill. In dark mode the sidebar is always dark (#15100f bg), so the monochrome fill should be light enough to read — consider `#e8c4cf` (dark theme primary) or `#f0e6e9` (light text) as the single fill color instead of the dark red. **Agent decides which reads better on the dark sidebar background.**
- The deployment runbook should note that the Pi build needs `data/covers/` directory to exist alongside the binary (covers are not embedded, only the React frontend is embedded).

</specifics>

<deferred>
## Deferred Ideas

- **Live VPS/Raspberry Pi deployment** — automation (SSH deploy script, GitHub Actions CI/CD) deferred to a post-v1 phase when ready to go live
- **SSL certificate automation** — certbot timer or Caddy auto-HTTPS deferred with deployment
- **systemd unit wiring** — documented in runbook but not deployed this phase
- **Open Graph meta tags** (from Phase 2 Plan 2.2) — already in the Go server; no changes needed in Phase 5
- **WebP conversion for covers** — v2 requirement (V2-04)

</deferred>

---

*Phase: 05-polish-production-deployment*
*Context gathered: 2026-04-03*
