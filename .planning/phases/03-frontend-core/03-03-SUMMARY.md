---
phase: 03-frontend-core
plan: 03
subsystem: ui
tags: [react, typescript, vite, react-router, tanstack-query, marked, lucide-react, css, msw]

# Dependency graph
requires:
  - phase: 03-01
    provides: BookCover, SkeletonCard, Toast, ThemeToggle, useTheme, CSS design system
  - phase: 03-02
    provides: BookCard, BookGrid, fetchCurrentlyReading, Book type, useIntersectionObserver, useScrollRestoration

provides:
  - Bio section with Markdown-parsed content and Goodreads profile link
  - NowReadingSection with conditional render (hidden when empty, capped at 4 books)
  - NowReadingCard component matching BookCard proportions
  - Sidebar with fixed desktop layout and mobile hamburger drawer (focus trap, Escape, backdrop)
  - HomePage assembling Bio + NowReadingSection + BookGrid
  - App.tsx two-column layout wiring (Sidebar + content area)
  - placehold.co cover placeholder for missing/broken images with word-wrapped title
affects: [04-book-detail, 05-author-genre-pages]

# Tech tracking
tech-stack:
  added: [marked (markdown parsing)]
  patterns: [conditional section rendering via null return, placehold.co for missing covers, always-dark sidebar with CSS vars override]

key-files:
  created:
    - frontend/src/content/bio.md
    - frontend/src/components/Bio.tsx
    - frontend/src/components/Bio.css
    - frontend/src/components/NowReadingCard.tsx
    - frontend/src/components/NowReadingSection.tsx
    - frontend/src/components/Sidebar.tsx
    - frontend/src/components/Sidebar.css
    - frontend/src/pages/HomePage.tsx
    - frontend/src/App.css
  modified:
    - frontend/src/App.tsx
    - frontend/src/components/BookCover.tsx
    - frontend/src/components/BookGrid.css

key-decisions:
  - "marked.parse() returns string synchronously in v17 (parseSync doesn't exist)"
  - "Sidebar always dark (#15100f) regardless of page theme — CSS vars overridden with hardcoded values"
  - "Mobile drawer requires display:none on desktop (not just off-screen transform) to prevent normal-flow layout issues"
  - "App layout uses block (not flex) since sidebar is position:fixed — margin-left:240px on app-content"
  - "Cover placeholder uses placehold.co with word-wrapped multi-line title (max 3 lines, 16 chars each)"
  - "NowReadingCard title uses 13px (same as BookCard) for visual cover size consistency"

patterns-established:
  - "Sidebar pattern: desktop fixed aside + mobile topbar + mobile drawer — both rendered in DOM, drawer hidden via display:none on desktop"
  - "Cover placeholder: switch img src on onError; handle empty src immediately (onError doesn't fire for empty src)"
  - "Conditional section: return null when data empty — section completely absent from DOM"

requirements-completed: [HOME-01, HOME-02, HOME-03, HOME-04, UI-10]

# Metrics
duration: 90min
completed: 2026-04-03
---

# Phase 03-03: Home Page Assembly Summary

**Full home page with Bio (Markdown), Now Reading (conditional grid), responsive sidebar (mobile drawer with focus trap), and placehold.co cover placeholders — user-approved**

## Performance

- **Duration:** ~90 min (including user feedback iterations)
- **Completed:** 2026-04-03
- **Tasks:** 3 (2 auto + 1 human-verify checkpoint)
- **Files modified:** 15+

## Accomplishments
- Bio section renders Florian's photo + marked.parse() HTML + Goodreads link
- NowReadingSection hides completely when API returns empty array (D-04); caps at 4 books
- Sidebar always dark (#15100f background), red #6d233e left-border accent on active items, lucide-react icons
- Mobile hamburger drawer: slides in from left, focus trap, Escape/backdrop/link-click close
- App layout fixed (block not flex — fixed sidebar doesn't participate in flex layout)
- BookCover: placehold.co fallback with word-wrapped title on error or empty src
- User visual approval received after iterative feedback rounds

## Task Commits

1. **Task 1: Bio + NowReadingSection** - `c64c704` (feat)
2. **Task 2: Sidebar + HomePage + App wiring** - `388d673` (feat)
3. **Task 3: Human visual verification** - approved by user
4. **UI feedback fixes** - `abd5013`, `d7fae0e`, `363da8f`, `d451a61`, `71d9241`, `98b7356` (fix)

## Files Created/Modified
- `frontend/src/content/bio.md` - Placeholder bio with Goodreads link
- `frontend/src/components/Bio.tsx` - marked.parse() + photo layout
- `frontend/src/components/NowReadingCard.tsx` - Book cover card matching BookCard proportions
- `frontend/src/components/NowReadingSection.tsx` - Conditional grid, cap at 4, error Toast
- `frontend/src/components/Sidebar.tsx` - Desktop + mobile drawer with lucide icons
- `frontend/src/pages/HomePage.tsx` - Bio + NowReadingSection + BookGrid assembly
- `frontend/src/App.tsx` - Layout wiring with Sidebar
- `frontend/src/App.css` - Block layout, margin-left:240px, overflow-x:hidden
- `frontend/src/components/BookCover.tsx` - placehold.co placeholder, empty src handling

## Decisions Made
- `marked.parseSync` doesn't exist in v17 — used `marked.parse()` which is synchronous
- Sidebar uses hardcoded dark colors (not CSS vars) to stay dark in both themes
- Mobile drawer must be `display:none` on desktop, not just off-screen — otherwise it appears in normal document flow and shifts content down
- App layout changed from `display:flex` to block — `position:fixed` sidebar exits flex flow, causing `app-content` to take full viewport width and overflow

## Deviations from Plan

### Auto-fixed Issues

**1. marked.parseSync API mismatch**
- **Issue:** Plan specified `marked.parseSync()` but marked v17 only has `marked.parse()` (synchronous string return)
- **Fix:** Used `marked.parse() as string`
- **Files modified:** `frontend/src/components/Bio.tsx`

**2. Mobile drawer rendering on desktop**
- **Issue:** `.sidebar--drawer` had no `display:none` rule on desktop — it rendered in normal flow, pushing Bio section down and showing duplicate icons
- **Fix:** Added `display:none` default; `display:flex` only inside `@media (max-width: 767px)`

**3. App layout horizontal overflow**
- **Issue:** `display:flex` on `.app-layout` with `position:fixed` sidebar caused `app-content` to take full viewport width; `margin-left:240px` shifted it without reducing width
- **Fix:** Removed `display:flex`, used block layout

---

**Total deviations:** 3 auto-fixed
**Impact on plan:** All fixes required for correct layout and API compatibility. No scope creep.

## Issues Encountered
- User iterative feedback required 6 fix commits after Task 2 completion (sidebar dark theme, layout overflow, cover placeholders, grid sizing, font loading, mobile drawer visibility)

## Next Phase Readiness
- Full home page complete and user-approved
- All HOME-01 through HOME-04 and UI-10 requirements satisfied
- BookGrid, NowReadingSection, Sidebar all integrated in App.tsx
- Phase 4 (book detail page) can import HomePage as reference for layout patterns

---
*Phase: 03-frontend-core*
*Completed: 2026-04-03*
