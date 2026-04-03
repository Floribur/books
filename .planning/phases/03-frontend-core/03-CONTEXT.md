# Phase 3: Frontend Core - Context

**Gathered:** 2026-04-03
**Status:** Ready for planning

<domain>
## Phase Boundary

Scaffold the entire React frontend and deliver a fully working home page: design system (CSS custom properties, dark/light mode, typography), reusable BookCover/BookCard/BookGrid components, infinite-scrolling "Books Read" section, "Now Reading" section, bio section, and responsive sidebar layout. This phase establishes all foundational patterns that Phase 4 pages will reuse.

</domain>

<decisions>
## Implementation Decisions

### Bio Section
- **D-01:** Bio content (text + photo reference) is stored in a **Markdown file** bundled with the app (e.g. `src/content/bio.md`). This allows richer formatting without touching JSX, and easy updates without a redeploy logic change.
- **D-02:** Layout: **photo left, text right** (classic author bio layout) at desktop. Agent decides the mobile stack order.
- **D-03:** Bio section includes a link to the Goodreads profile: `https://www.goodreads.com/user/show/79499864-florian`

### "Now Reading" Section
- **D-04:** When `GET /api/books/currently-reading` returns an empty array, **hide the section entirely** — no empty state message, no placeholder.
- **D-05:** Show a maximum of **3–4 books** (agent decides the exact cap based on what looks best at desktop widths). Books beyond the cap are hidden — not truncated with a "show more" link, just omitted.
- **D-06:** "Now Reading" books display in a **responsive grid with fewer columns than "Books Read"** because the covers are larger/more prominent. No horizontal scroll — books wrap to next row if the grid allows multiple rows within the cap.
- **D-07:** Each "Now Reading" card shows: **larger cover + title + author** (more prominent than Books Read cards — the "currently reading" context warrants emphasis).

### "Books Read" Section
- **D-08:** Each book card shows: **cover + title + author name(s) only** — clean and minimal. No read date, no genres, no year on the card.
- **D-09:** Section heading is simply **"Books Read"** — no total count displayed.
- **D-10:** Grid column count at each breakpoint is **agent's discretion** — choose values that feel balanced (e.g. 6 at wide desktop, 4–5 at medium, 2 at mobile per UI-10).

### Dark/Light Mode Toggle
- **D-11:** Toggle lives **in the sidebar, near the nav links**.
- **D-12:** Toggle style: **sun/moon icon button** — no label text, icon conveys meaning.
- **D-13:** Behavior: detect `prefers-color-scheme` on first load, persist in `localStorage`, toggle on click. Applies via `[data-theme="dark"]` / `[data-theme="light"]` on `<html>` (per ROADMAP Plan 3.1 Task 4).

### Loading & Error States
- **D-14:** While book data is loading: show **skeleton cards** — placeholder rectangles matching the card shape (cover placeholder + text line placeholders). Use brand gradient for the cover skeleton consistent with UI-09.
- **D-15:** When an API fetch fails: show a **toast/notification at the top of the page** with an error message. Agent decides toast library or CSS implementation.

### Mobile Navigation
- **D-16:** Mobile sidebar collapse behavior is **agent's discretion** — pick whatever works best for the layout (tabs, hamburger, etc.). The requirement is that the sidebar nav is accessible on mobile.

### Agent's Discretion
- Exact "Now Reading" display cap (3 or 4 — agent picks based on visual balance)
- Book grid column counts per breakpoint (desktop/tablet/mobile — agent picks balanced values; mobile is 2 per UI-10)
- Mobile nav pattern (agent picks: tabs with icons, hamburger, or other)
- Toast/notification implementation (small CSS component or lightweight lib)
- Bio section mobile stack order (photo then text, or text then photo)
- Markdown parsing approach for bio content (vite-plugin-md, raw import + marked, or similar)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Requirements
- `.planning/REQUIREMENTS.md` §Home Page (HOME-01–04) — bio, Now Reading, Books Read, click-through
- `.planning/REQUIREMENTS.md` §UI/UX (UI-01–03, UI-06–10) — infinite scroll, back-button, dark mode, brand colors, typography, cover placeholder, responsive layout

### Roadmap
- `.planning/ROADMAP.md` §Phase 3 — Plan 3.1, 3.2, 3.3 task breakdown (CSS system, BookGrid, home page assembly)

### API Contract (from Phase 2)
- `.planning/phases/02-go-rest-api/02-CONTEXT.md` — D-01/D-02: paginated envelope shape `{"items":[...], "next_cursor":"...", "has_more":true}`; opaque base64 cursor; `GET /api/books/currently-reading` returns plain array

### Project Constraints
- `.planning/PROJECT.md` §Constraints — tech stack (React + TypeScript + Vite)
- `.planning/STATE.md` §Recent Decisions — TanStack Query v5, cursor-based pagination, CSS keyframes (no Lottie), Intersection Observer for infinite scroll

No external ADRs or specs — requirements and decisions fully captured above.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `frontend/dist/` — placeholder only (`.gitkeep` + empty `index.html`); frontend is greenfield, no existing React code
- `frontend/embed.go` — Go embed directive stub; will be populated when frontend builds

### Established Patterns (backend, for API integration reference)
- Paginated API response: `{"items":[...], "next_cursor":"...", "has_more":true}` — TanStack Query's `getNextPageParam` reads `next_cursor`
- Cover images served at `/covers/{isbn13}.jpg` with immutable cache headers (Plan 2.2)

### Integration Points
- Vite dev proxy: `localhost:8081` (Go API) proxied as `/api` and `/covers` during development
- React Router catch-all handled by Go SPA handler — all client-side routes must be valid SPA routes
- Phase 5 will run `npm run build` and Go embeds `frontend/dist` — the Vite build must output to `frontend/dist/`

</code_context>

<specifics>
## Specific Ideas

- Bio Markdown file path: `src/content/bio.md` (or similar) — Vite can import raw text or use a plugin
- Goodreads profile link in bio: `https://www.goodreads.com/user/show/79499864-florian`
- "Now Reading" cards are visually larger/more prominent than Books Read cards — this distinction should be clear in the component design
- Skeleton cards use brand gradient (consistent with CSS gradient placeholder from UI-09)

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 03-frontend-core*
*Context gathered: 2026-04-03*
