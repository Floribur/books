# Phase 4: Frontend Pages - Context

**Gathered:** 2026-04-03
**Status:** Ready for planning

<domain>
## Phase Boundary

Build all secondary pages — book detail, author index/detail, genre index/detail, and Reading Challenge. All stub routes are already wired in `App.tsx`; Phase 4 fills them in using the components, hooks, API client, and design system established in Phase 3.

</domain>

<decisions>
## Implementation Decisions

### Book Detail Page Layout

- **D-01:** Desktop: two-column layout — cover left (full height), metadata right.
- **D-02:** Mobile: cover stays **small and floats left**, metadata text wraps around it (classic inline float layout). Cover does NOT collapse into a full-width block on mobile.
- **D-03:** Metadata order (right column / below float): **Title → Author links → Genres (tags) → Publication year · Page count · Read date**. Author names are clickable links to their detail pages; genres are clickable links to genre detail pages.
- **D-04:** Description is **expandable only when it exceeds ~8 lines / ~640 characters**. Shorter descriptions are always fully visible. A "Show more" link expands it in place; once expanded, a "Show less" link collapses it.

### Genre Index Page

- **D-05:** Genre index displays genres sorted by book count descending, with **horizontal bar visualization**: genre name on the left, bar fills proportionally to the max book count on the right, count shown as a number at the end of the bar.
- **D-06:** Bar fill color: **lighter tint of brand color (#6d233e) at ~30% opacity** — subtle, consistent with the design system.
- **D-07:** Each genre row is a clickable link to the genre detail page.

### Reading Challenge Page

- **D-08:** Default year on first load: **most recent year with books** (from `GET /api/years`), not necessarily the current calendar year. This ensures the page always shows data immediately.
- **D-09:** Stats strip shows **three stats**: book count for the year, total pages* (sum of `page_count` for the year), and longest book (title + page count).
- **D-10:** Total pages stat includes a footnote marker `*` with the note "based on available data" (since ~70-80% of books have page_count from Google Books). Longest book also relies on page_count, so if no page data exists for a year, those two stats are hidden.
- **D-11:** Year selector: prev/next arrow buttons flanking the year display. Keyboard accessible (arrow keys). Year is reflected in the URL as `?year=`.

### Author Index Page

- **D-12:** Author index is a **simple sorted scrollable list** — alphabetical by surname, book count per author shown, each row links to the author detail page. No A-Z letter jump anchors.

### Agent's Discretion

- Exact CSS approach for the floating cover on mobile (float or inline flex)
- Animation/transition for the description expand/collapse (e.g., max-height transition)
- Bar width calculation (CSS custom property or inline style with `width: calc(count / max * 100%)`)
- Reading Challenge empty state if a year has no books (though this should not occur for past years)
- Exact BookGrid queryKey shapes for author-filtered and genre-filtered grids

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Requirements
- `.planning/REQUIREMENTS.md` §Book Detail (BOOK-01–04) — detail page content requirements
- `.planning/REQUIREMENTS.md` §Authors (AUTH-01–02) — author index and detail requirements
- `.planning/REQUIREMENTS.md` §Genres (GENR-01–02) — genre index and detail requirements
- `.planning/REQUIREMENTS.md` §Reading Challenge (CHAL-01–04) — challenge page requirements

### Roadmap
- `.planning/ROADMAP.md` §Phase 4 — Plan 4.1 and 4.2 task breakdown

### API Contract (from Phase 2)
- `.planning/phases/02-go-rest-api/02-CONTEXT.md` — D-04/D-05: book list vs detail shapes; D-06/D-07: author/genre shapes; D-01/D-02: paginated envelope + cursor

### Phase 3 Patterns (reuse these)
- `.planning/phases/03-frontend-core/03-CONTEXT.md` — D-08: BookCard shape (cover+title+author only); D-14/D-15: skeleton + toast patterns

### Existing Code
- `frontend/src/components/BookGrid.tsx` — reuse for all author-filtered, genre-filtered, year-filtered grids
- `frontend/src/components/BookCard.tsx` — reuse as-is for all grid views
- `frontend/src/api/types.ts` — Book, Author, Genre, PaginatedBooks interfaces
- `frontend/src/api/client.ts` — apiFetch helper
- `frontend/src/App.tsx` — stub routes already wired; Phase 4 replaces the stub components

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `BookGrid` (accepts `queryKey` + `fetchFn`) — handles infinite scroll, skeleton loading, error toast; use for all filtered grids in Phase 4
- `BookCard` — cover + title + author, links to `/books/:slug`; reuse as-is
- `BookCover` — aspect-ratio 2/3, lazy loading, gradient placeholder; reuse on detail page
- `SkeletonCard` — loading placeholder; used automatically by BookGrid
- `Toast` — error notification; used automatically by BookGrid
- `api/client.ts` `apiFetch` — typed API fetch helper; use for all new endpoints

### Established Patterns
- TanStack Query v5 `useInfiniteQuery` for paginated lists (cursor-based)
- TanStack Query v5 `useQuery` for non-paginated fetches (years list, book detail, author detail, genre detail)
- CSS custom properties from `styles/tokens.css` and `styles/themes.css` — use these for all new styles
- Playfair Display for headings, Inter for body (from typography.css)

### Integration Points
- `App.tsx` stub routes need to be replaced with real page components
- `api/books.ts` needs new fetch functions: `fetchBookBySlug`, `fetchAuthors`, `fetchAuthorBySlug`, `fetchGenres`, `fetchGenreBySlug`, `fetchYears`, `fetchBooksByAuthor`, `fetchBooksByGenre`, `fetchBooksByYear`
- Year-filtered books use `GET /api/books?year={year}` — BookGrid's fetchFn should accept a year param

</code_context>

<specifics>
## Specific Ideas

- The floating cover on mobile (D-02) is intentionally compact — a small cover beside wrapped text, not a full-width hero. This keeps the author bio / metadata visible without much scrolling.
- Genre bars (D-05/D-06) are clickable rows — the entire row (name + bar + count) is a link to the genre detail page.
- Reading Challenge stat "longest book" shows the title + page count, e.g. "The Name of the Wind (662 pages)". If multiple books tie, agent picks one.
- Stats strip for Reading Challenge should be visually compact — a horizontal strip of labeled numbers, not cards or large blocks.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 04-frontend-pages*
*Context gathered: 2026-04-03*
