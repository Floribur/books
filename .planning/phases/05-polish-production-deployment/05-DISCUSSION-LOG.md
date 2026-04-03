# Phase 5: Polish & Production Deployment - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-03
**Phase:** 05-polish-production-deployment
**Areas discussed:** Animation style & feel, SSL & deployment, UI polish focus

---

## Animation Style & Feel

| Option | Description | Selected |
|--------|-------------|----------|
| Open book page-turn | SVG open book, pages flip periodically (CSS keyframes) | |
| Soft pulse / glow | Static open book with breathing/glow CSS animation | |
| Reading motion dot | Cursor drifting across text lines (CSS keyframes) | |
| Abstract stacked books | Shimmer/ripple effect on book stack | |
| Lottie: free LottieFiles animation | Person sitting and reading, turning pages | ✓ |

**User's choice:** Lottie animation of a person reading — specifically wanted a character that resembles themselves (tall, sporty, glasses, short hair that stands up slightly, shirt). Initially considered CSS keyframes variants but changed mind for a more expressive character.

**Notes:** User referenced LottieFiles examples:
- `https://lottiefiles.com/animation/reading-book-animation_11806166`
- `https://lottiefiles.com/free-animation/boy-reading-a-book-f2fyScGybx`

This overrides the Phase 3 decision "CSS keyframes only, no Lottie" — original concern was 60KB bundle cost for a simple decorative element. New vision warrants Lottie.

### Animation Color

| Option | Description | Selected |
|--------|-------------|----------|
| Full brand palette recolor | Match #6d233e primary, #c4843a gold accent, dark sidebar bg | |
| Monochrome silhouette | Single color, brand red or similar, dead simple | ✓ |

**User's choice:** Monochrome silhouette — one fill color, keep it simple.

### Animation Placement

| Option | Description | Selected |
|--------|-------------|----------|
| Above nav links | Between title text and nav links | |
| Below nav links | Above the theme toggle | |
| Replaces title | Animation IS the header — text title removed | ✓ |

**User's choice:** Replace the "Flo's Library" text title with the Lottie player.

---

## SSL & Deployment

| Option | Description | Selected |
|--------|-------------|----------|
| Caddy reverse proxy | Auto-HTTPS, proxy to Go on :8081 | |
| Go serves TLS directly | autocert built into binary | |
| certbot + HTTP binary | Manual renewal, most complex | |
| No deployment this phase | Local build only + deployment runbook | ✓ |

**User's choice:** No live deployment this phase. Focus on local production build workflow (make build → single binary). Write a Raspberry Pi deployment runbook (docs/deployment.md) with Caddy + systemd unit documented but not automated.

**Notes:** User mentioned Raspberry Pi as target (not a VPS). Deployment will be done manually when ready.

---

## UI Polish Focus

| Option | Description | Selected |
|--------|-------------|----------|
| General consistency pass | Brand colors, contrast, spacing, typography across all pages | ✓ |
| Specific rough edges | Fix particular issues noticed during testing | |
| Both | Specific first, then general | |

**User's choice:** General consistency pass — no specific bugs reported, just want everything to feel polished and consistent.

### Favicon

| Option | Description | Selected |
|--------|-------------|----------|
| Book icon | Simple open book silhouette in brand red | |
| 'F' lettermark | Stylized F in Playfair Display, brand red | |
| Reader silhouette | Matches Lottie character, small favicon version | ✓ |

**User's choice:** Reader silhouette favicon matching the Lottie character style.

### Page Titles

| Option | Description | Selected |
|--------|-------------|----------|
| 'Flo's Library' everywhere | Same title on all pages | |
| 'Flo's Library — [Page Name]' | Scoped title on secondary pages | ✓ |

**User's choice:** `Flo's Library — [Page Name]` format for secondary pages (book, author, genre, reading challenge). Home page: `Flo's Library`.

---

## Agent's Discretion

- Which specific Lottie animation file to use
- How to recolor to monochrome (Lottie colorFilters, JSON post-process, or find single-layer animation)
- Lottie wrapper library choice (`lottie-react` vs `@lottiefiles/react-lottie-player`)
- Animation sizing in sidebar
- Mobile drawer: whether to show animation or text title
- Page title implementation approach (hook vs per-component `document.title`)
- Pi cross-compile target (arm64 vs armv7)

## Deferred Ideas

- Live Raspberry Pi deployment (automation, SSH scripts, CI/CD)
- SSL automation
- systemd unit wiring (documented only)
