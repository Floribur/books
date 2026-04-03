# Phase 3: Frontend Core - Research

**Researched:** 2026-04-03
**Domain:** React + TypeScript + Vite frontend, TanStack Query v5, CSS design system, infinite scroll
**Confidence:** HIGH (verified against npm registry; all versions current as of research date)

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Bio content stored in Markdown file `src/content/bio.md` (raw import + `marked`)
- **D-02:** Bio layout: photo left, text right at desktop; agent decides mobile stack order
- **D-03:** Bio section includes Goodreads profile link: `https://www.goodreads.com/user/show/79499864-florian`
- **D-04:** "Now Reading" hidden entirely when API returns empty array — no empty state
- **D-05:** Show maximum 3–4 books in Now Reading (agent caps at 4)
- **D-06:** Now Reading uses fewer columns than Books Read (no horizontal scroll)
- **D-07:** Now Reading cards show larger cover + title + author (more prominent)
- **D-08:** Books Read cards show cover + title + author only (no date, genres, year)
- **D-09:** Section heading "Books Read" — no total count displayed
- **D-10:** Grid columns per breakpoint: agent's discretion (UI-SPEC resolved: 6/5/4/2)
- **D-11:** Dark/light toggle in sidebar near nav links
- **D-12:** Toggle style: sun/moon icon button, no label text
- **D-13:** Detect `prefers-color-scheme`, persist in `localStorage`, toggle via `[data-theme]` on `<html>`
- **D-14:** Loading state: skeleton cards with brand gradient (not spinner)
- **D-15:** API errors: toast/notification at top of page
- **D-16:** Mobile nav: agent's discretion (UI-SPEC resolved: hamburger → left drawer overlay)
- CSS custom properties with `[data-theme]` attribute (no CSS-in-JS)
- TanStack Query v5 `useInfiniteQuery` with cursor-based pagination
- Intersection Observer for infinite scroll sentinel
- React Router v6 (must pin `react-router-dom@6` — v7 is current but not selected)

### Agent's Discretion
- Exact Now Reading cap: **4 books** (resolved in UI-SPEC)
- Grid column counts per breakpoint: **6/5/4/2** (resolved in UI-SPEC)
- Now Reading grid columns: **4/3/2** (resolved in UI-SPEC)
- Mobile nav pattern: **hamburger → left drawer overlay** (resolved in UI-SPEC)
- Toast implementation: **CSS-only + `role="alert"`** (resolved in UI-SPEC)
- Bio mobile stack order: **photo first, then text** (resolved in UI-SPEC)
- Markdown parsing: **raw import + `marked`** (resolved in UI-SPEC)

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| HOME-01 | Bio section with Florian's photo and personal description | Markdown + `marked` for bio content; photo rendered as `<img>` with 2:3 ratio; flex layout |
| HOME-02 | "Now Reading" section with currently-reading books (cover, title, author) | `GET /api/books/currently-reading` → plain array; TanStack Query `useQuery` (not infinite); conditional render |
| HOME-03 | "Books Read" section in descending order | `GET /api/books` paginated; TanStack Query `useInfiniteQuery` v5 with `initialPageParam` |
| HOME-04 | Book cards clickable → navigate to detail page | `<a>` wrapping entire card; React Router v6 `<Link>` or `useNavigate` |
| UI-01 | Infinite scroll — first ~24 books immediately, more on scroll | Intersection Observer sentinel at bottom of grid; `threshold: 0.1` |
| UI-02 | "Load More" button as accessible fallback alongside Intersection Observer | Same `fetchNextPage` call; `disabled` while `isFetchingNextPage`; hidden when `!hasNextPage` |
| UI-03 | Back-button restores scroll position | `sessionStorage` + URL param for book ID; `scrollIntoView({ behavior: 'instant' })` on mount |
| UI-06 | Dark/light mode toggle with OS preference detection and localStorage persistence | `prefers-color-scheme` media query on first load; `localStorage` key `"theme"`; `[data-theme]` on `<html>` |
| UI-07 | Brand colors: `#6d233e` primary, `#c4843a` accent, `#faf8f5` background | CSS custom properties in `:root` and `[data-theme]` selectors |
| UI-08 | Typography: Playfair Display (headings) + Inter (body) | Google Fonts import in `index.html` or CSS `@import` |
| UI-09 | Book cover placeholder: CSS gradient using brand colors (no broken-image icon) | `--color-skeleton` gradient applied to cover container; `onError` hides `<img>` |
| UI-10 | Responsive layout, mobile collapses to 2-column book grid | CSS grid `repeat()` with breakpoint media queries; sidebar drawer at <768px |
</phase_requirements>

---

## Summary

Phase 3 scaffolds the entire React frontend from a greenfield `frontend/` directory. The key technical areas are: (1) Vite 8 project setup with TypeScript and React 19, (2) CSS design system built entirely from custom properties without any UI library, (3) TanStack Query v5 `useInfiniteQuery` for cursor-based infinite scroll, and (4) several UX patterns that have specific implementation requirements (scroll restoration, dark mode persistence, Intersection Observer).

**Critical version alert:** React Router v7 is the current npm latest (7.14.0), but the locked decision is React Router v6. Must explicitly pin `react-router-dom@6` in package.json. The `marked` library is now at v17.0.5 (ESM-only) and its `parse()` function returns a `Promise<string>` — synchronous use requires `marked.parseSync()`. TanStack Query v5 introduces `initialPageParam` as a required parameter in `useInfiniteQuery` — missing this causes a TypeScript error.

**Primary recommendation:** Scaffold with `npm create vite@latest frontend -- --template react-ts` (produces Vite 8 + React 19 + TypeScript 6), then install all dependencies with explicit versions for the locked packages. Use `marked.parseSync()` for bio rendering. Use `isPending` (not `isLoading`) for skeleton card display in TanStack Query v5.

---

## Standard Stack

### Core
| Library | Version (verified) | Purpose | Why Standard |
|---------|-------------------|---------|--------------|
| react | 19.2.4 | UI component framework | Locked decision |
| react-dom | 19.2.4 | React DOM renderer | Required peer dep |
| typescript | 6.0.2 | Type safety | Locked decision |
| vite | 8.0.3 | Dev server + bundler | Locked decision |
| @vitejs/plugin-react | 6.0.1 | React Fast Refresh in Vite | Standard Vite+React plugin |
| react-router-dom | 6.30.3 | Client-side routing | **Pin to v6** — v7 is current but NOT selected |
| @tanstack/react-query | 5.96.2 | Server state + infinite scroll | Locked decision |
| lucide-react | 1.7.0 | Sun/Moon/Menu/X icons only | Resolved in UI-SPEC (D-12, tree-shakeable) |
| marked | 17.0.5 | Markdown → HTML for bio | Resolved in UI-SPEC; use `marked.parseSync()` |

### Dev / Testing
| Library | Version (verified) | Purpose | When to Use |
|---------|-------------------|---------|-------------|
| vitest | 4.1.2 | Test runner | All unit + component tests |
| @testing-library/react | 16.3.2 | Component rendering in tests | All component tests |
| @testing-library/user-event | latest | Simulate user interactions | Click, scroll, keyboard tests |
| msw | 2.12.14 | Mock Service Worker (v2 API) | API integration tests |
| jsdom | latest | DOM environment for vitest | Component tests |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| CSS custom properties | Tailwind, styled-components | Both excluded by project constraint |
| marked v17 | remark/rehype | Higher complexity; marked is 6KB min+gz and sufficient for bio |
| react-router-dom v6 | react-router-dom v7 | v7 has framework mode complexity; v6 is the locked decision |
| CSS-only toast | react-hot-toast, sonner | Dependencies add bundle size; UI-SPEC resolved to CSS-only |

**Installation:**
```bash
# From within /frontend after vite scaffold
npm install react-router-dom@6 @tanstack/react-query lucide-react marked

# Dev dependencies
npm install -D vitest @testing-library/react @testing-library/user-event msw @vitest/coverage-v8 jsdom
```

**Note on `npm create vite@latest`:** This scaffolds with React 19 and Vite 8. After scaffold, delete boilerplate (`App.css`, `assets/react.svg`, default `App.tsx` content).

---

## Architecture Patterns

### Recommended Project Structure
```
frontend/
├── src/
│   ├── components/
│   │   ├── BookCard.tsx         # Books Read card (cover + title + author)
│   │   ├── BookCover.tsx        # Shared cover img with gradient placeholder
│   │   ├── BookGrid.tsx         # Infinite-scroll grid (useInfiniteQuery)
│   │   ├── NowReadingCard.tsx   # Now Reading card (larger, more prominent)
│   │   ├── NowReadingSection.tsx# Wraps NowReadingCard grid + hide-when-empty logic
│   │   ├── SkeletonCard.tsx     # Loading placeholder matching BookCard shape
│   │   ├── Sidebar.tsx          # Nav + dark mode toggle; drawer on mobile
│   │   ├── ThemeToggle.tsx      # Sun/Moon icon button
│   │   └── Toast.tsx            # CSS-only error toast with role="alert"
│   ├── pages/
│   │   └── HomePage.tsx         # Assembles Bio + NowReadingSection + BookGrid
│   ├── hooks/
│   │   ├── useTheme.ts          # prefers-color-scheme + localStorage + [data-theme]
│   │   ├── useIntersectionObserver.ts  # Sentinel observer hook
│   │   └── useScrollRestoration.ts    # sessionStorage scroll target logic
│   ├── api/
│   │   ├── client.ts            # fetch wrapper (base URL, error handling)
│   │   ├── books.ts             # fetchBooks (paginated), fetchCurrentlyReading
│   │   └── types.ts             # TypeScript interfaces for API responses
│   ├── content/
│   │   └── bio.md               # Bio markdown file (raw import)
│   ├── styles/
│   │   ├── tokens.css           # All CSS custom properties (colors, spacing, type)
│   │   ├── themes.css           # [data-theme="light"] and [data-theme="dark"]
│   │   ├── typography.css       # Font imports, base type rules
│   │   ├── reset.css            # Minimal CSS reset
│   │   └── global.css           # Imports all above; sets html/body base
│   ├── main.tsx                 # React root, QueryClientProvider, RouterProvider
│   └── App.tsx                  # Route definitions
├── public/
│   └── florian.jpg              # Bio photo (or reference from src/assets/)
├── index.html                   # Google Fonts link tags, [data-theme] initial value
├── vite.config.ts               # Proxy to :8081, build output to dist/
└── tsconfig.json
```

### Pattern 1: TanStack Query v5 useInfiniteQuery (CRITICAL — v5 API differs from v4)

**What:** Fetches paginated books with cursor-based pagination. V5 introduced `initialPageParam` as a **required** field.

**When to use:** `BookGrid` component for "Books Read" section.

```typescript
// Source: TanStack Query v5 docs — initialPageParam is REQUIRED in v5 (breaking from v4)
import { useInfiniteQuery } from '@tanstack/react-query';

interface PaginatedResponse {
  items: Book[];
  next_cursor: string | null;
  has_more: boolean;
}

const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isPending } =
  useInfiniteQuery<PaginatedResponse>({
    queryKey: ['books'],
    queryFn: ({ pageParam }) =>
      fetchBooks(pageParam as string | undefined),
    initialPageParam: undefined,          // ← REQUIRED in v5, was implicit in v4
    getNextPageParam: (lastPage) =>
      lastPage.has_more ? lastPage.next_cursor : undefined,  // undefined = no more pages
  });

// Flatten pages for rendering
const books = data?.pages.flatMap((page) => page.items) ?? [];

// Show skeleton: use isPending, NOT isLoading (v5 rename)
// isPending = no data yet (initial load)
// isLoading in v5 = isPending && isFetching (different from v4!)
if (isPending) return <SkeletonGrid count={12} />;
```

### Pattern 2: Intersection Observer Sentinel

**What:** Triggers `fetchNextPage` when bottom of grid enters viewport.

**When to use:** Inside `BookGrid` component, alongside the "Load More" button (UI-02).

```typescript
// Source: MDN Intersection Observer API; pattern is standard 2025 React
import { useEffect, useRef } from 'react';

function useIntersectionObserver(
  callback: () => void,
  enabled: boolean
): React.RefObject<HTMLDivElement> {
  const sentinelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!enabled || !sentinelRef.current) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) callback();
      },
      { threshold: 0.1 }   // 10% visible triggers fetch
    );

    observer.observe(sentinelRef.current);
    return () => observer.disconnect();  // cleanup is critical to avoid memory leaks
  }, [enabled, callback]);   // re-run when enabled changes (e.g. hasNextPage flips)

  return sentinelRef;
}

// In BookGrid:
const sentinelRef = useIntersectionObserver(
  () => { if (hasNextPage && !isFetchingNextPage) fetchNextPage(); },
  !!hasNextPage    // disable observer when no more pages
);

return (
  <>
    <div className="book-grid">{/* cards */}</div>
    <div ref={sentinelRef} style={{ height: 1 }} aria-hidden="true" />
    <button
      onClick={() => fetchNextPage()}
      disabled={isFetchingNextPage || !hasNextPage}
    >
      {isFetchingNextPage ? 'Loading…' : 'Load More Books'}
    </button>
    {!hasNextPage && <p className="text-muted">You've reached the end.</p>}
  </>
);
```

### Pattern 3: Dark/Light Mode with [data-theme]

**What:** System preference detection → localStorage persistence → `[data-theme]` toggle.

```typescript
// useTheme.ts
export function useTheme() {
  const [theme, setTheme] = useState<'light' | 'dark'>(() => {
    // 1. Check localStorage first
    const stored = localStorage.getItem('theme');
    if (stored === 'light' || stored === 'dark') return stored;
    // 2. Fall back to OS preference
    return window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light';
  });

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  }, [theme]);

  const toggle = () => setTheme((t) => (t === 'light' ? 'dark' : 'light'));
  return { theme, toggle };
}
```

**CSS:**
```css
/* tokens.css — light theme defaults */
:root,
[data-theme="light"] {
  --color-background: #faf8f5;
  --color-surface: #f0ebe3;
  --color-primary: #6d233e;
  --color-accent: #c4843a;
  --color-text: #2c1a1f;
  --color-text-muted: #7a6065;
  --color-border: #e2d9d0;
  --color-destructive: #b91c1c;
  --color-skeleton: linear-gradient(135deg, #6d233e22, #c4843a22);
}

[data-theme="dark"] {
  --color-background: #15100f;
  --color-surface: #221618;
  --color-primary: #e8c4cf;
  --color-accent: #d4943a;
  --color-text: #f0e6e9;
  --color-text-muted: #9e8085;
  --color-border: #3a2428;
  --color-destructive: #ef4444;
  --color-skeleton: linear-gradient(135deg, #6d233e44, #c4843a33);
}
```

> **FOUC prevention:** Set initial `data-theme` in `index.html` via an inline `<script>` before React mounts. This prevents the flash of wrong theme color.

```html
<!-- index.html — before </head> -->
<script>
  (function() {
    var stored = localStorage.getItem('theme');
    var theme = stored || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    document.documentElement.setAttribute('data-theme', theme);
  })();
</script>
```

### Pattern 4: marked v17 — ASYNC by Default

**What:** `marked.parse()` in v17 returns `Promise<string>`. Synchronous rendering requires `marked.parseSync()`.

```typescript
// WRONG in marked v17 — parse() returns a Promise, not a string:
// const html = marked.parse(bioContent);  ← BUG: html = Promise object

// CORRECT for synchronous use in a React component:
import { marked } from 'marked';
import bioRaw from '../content/bio.md?raw';   // Vite ?raw import

const bioHtml = marked.parseSync(bioRaw);     // ← parseSync() for synchronous HTML string

// In component:
<div
  className="bio-text"
  dangerouslySetInnerHTML={{ __html: bioHtml }}
/>
```

**Security:** `bio.md` is a local file bundled with the app — no user input, XSS risk is zero. `dangerouslySetInnerHTML` is appropriate here.

### Pattern 5: Back-Button Scroll Restoration (UI-03)

**What:** Save clicked book's ID to sessionStorage + URL; on mount, scroll to that book.

```typescript
// BookCard: on click, save to sessionStorage
function handleCardClick(slug: string) {
  sessionStorage.setItem('scrollTarget', slug);
}

// BookGrid: on mount, restore scroll
useEffect(() => {
  const target = sessionStorage.getItem('scrollTarget');
  if (!target) return;

  const el = document.getElementById(`book-${target}`);
  if (el) {
    el.scrollIntoView({ behavior: 'instant' });
    sessionStorage.removeItem('scrollTarget');
  }
}, [books]);  // re-run when books data loads (target card may not exist on first render)
```

**Note:** The scroll target element must exist in the DOM. With infinite scroll, the book may not be loaded yet. The approach: store both the cursor position and the target slug. On home page mount, if `scrollTarget` is in sessionStorage, skip fetching from page 1 and start from the saved cursor. This is more complex — see Pitfalls section.

### Pattern 6: Vite Dev Proxy Configuration

```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist',           // matches go:embed frontend/dist
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8081',
      '/covers': 'http://localhost:8081',
    },
  },
});
```

### Pattern 7: BookCover with Gradient Placeholder (UI-09)

```typescript
// BookCover.tsx — CSS gradient placeholder, no broken-image icon
function BookCover({ src, title, loading = 'lazy' }: BookCoverProps) {
  const [error, setError] = useState(false);

  return (
    <div className="book-cover-wrapper">
      {!error && (
        <img
          src={src}
          alt={`${title} cover`}
          loading={loading}
          onError={() => setError(true)}
          className="book-cover-img"
        />
      )}
      {/* Placeholder always rendered; hidden by CSS when image loads successfully */}
    </div>
  );
}
```

```css
.book-cover-wrapper {
  aspect-ratio: 2 / 3;
  background: var(--color-skeleton);   /* brand gradient placeholder */
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.book-cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
```

### Pattern 8: CSS-Only Error Toast (D-15)

```typescript
// Toast.tsx
interface ToastProps {
  message: string;
  onDismiss: () => void;
}

export function Toast({ message, onDismiss }: ToastProps) {
  useEffect(() => {
    const timer = setTimeout(onDismiss, 5000);
    return () => clearTimeout(timer);
  }, [onDismiss]);

  return (
    <div className="toast" role="alert" aria-live="assertive">
      <span>{message}</span>
      <button
        className="toast-close"
        onClick={onDismiss}
        aria-label="Dismiss"
      >
        ×
      </button>
    </div>
  );
}
```

```css
.toast {
  position: fixed;
  top: 16px;
  right: 16px;
  background: var(--color-destructive);
  color: #fff;
  padding: 12px 16px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  z-index: 1000;
  animation: toast-in 0.2s ease;
}
.toast-close {
  width: 44px;
  height: 44px;
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  font-size: 20px;
}
```

### Anti-Patterns to Avoid

- **Using `isLoading` for skeleton display in TanStack Query v5:** In v5, `isLoading` = `isPending && isFetching`. Use `isPending` instead for "show skeleton while no data exists".
- **Calling `marked.parse()` without await in v17:** `parse()` returns `Promise<string>`. Always use `marked.parseSync()` for synchronous contexts.
- **Forgetting `initialPageParam`:** TanStack Query v5 requires this explicitly. TypeScript will error without it, but verify before testing.
- **Installing `react-router-dom` without version pin:** `npm install react-router-dom` will install v7 (current). Always pin: `react-router-dom@6`.
- **Not disconnecting the Intersection Observer:** Memory leak if observer isn't cleaned up in `useEffect` return.
- **Setting `data-theme` only in React:** React mounts after page parse — causes FOUC. Use inline `<script>` in `index.html`.
- **Using Google Fonts `@import` in CSS:** Blocks rendering. Use `<link rel="preconnect">` + `<link rel="stylesheet">` in `index.html` instead.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Infinite scroll logic | Custom fetch-on-scroll | `useInfiniteQuery` (TanStack Query v5) | Handles loading states, deduplication, cache, error retry |
| API request state management | `useState` + `useEffect` fetches | TanStack Query `useQuery` / `useInfiniteQuery` | Race conditions, stale data, refetch on focus — all handled |
| Intersection Observer hook | Inline useEffect | `useIntersectionObserver` hook (project-local) | Reusable across BookGrid and any future lazy sections |
| Markdown parsing | Custom parser | `marked.parseSync()` | bio.md is simple — no need to hand-roll |
| Toast state | Complex state machine | Simple `useState<string[]>` with array of messages | CSS-only toast is sufficient for 2 error states in Phase 3 |

**Key insight:** The most error-prone area is infinite scroll + scroll restoration combined. TanStack Query handles the data fetching complexity; the scroll restoration is the only custom logic requiring care.

---

## Common Pitfalls

### Pitfall 1: React Router v7 Installed Instead of v6
**What goes wrong:** `npm install react-router-dom` installs v7.14.0. V7 has breaking changes (different `createBrowserRouter` behavior, new data router patterns).
**Why it happens:** npm installs latest by default; v7 was released late 2024.
**How to avoid:** Always specify version: `npm install react-router-dom@6`. Pin in `package.json`: `"react-router-dom": "^6.30.3"`.
**Warning signs:** TypeScript errors on `<BrowserRouter>`, different `useNavigate` type signatures.

### Pitfall 2: TanStack Query v5 — `isPending` vs `isLoading`
**What goes wrong:** Using `isLoading` to show skeleton cards. In v5, `isLoading = isPending && isFetching`, so it's `false` when a query is paused/disabled.
**Why it happens:** v4 used `isLoading` for "no data yet". v5 renamed this to `isPending`.
**How to avoid:** Use `isPending` for "show skeleton when no data exists yet".
**Warning signs:** Skeleton cards never show, or flicker unexpectedly.

### Pitfall 3: `marked.parse()` Returns Promise in v17
**What goes wrong:** `const html = marked.parse(bioContent)` produces `html = "[object Promise]"` rendered into the DOM.
**Why it happens:** marked v5+ made `parse()` async by default. v17 continues this pattern.
**How to avoid:** Use `marked.parseSync(bioContent)` for synchronous string return.
**Warning signs:** Bio section shows "[object Promise]" text.

### Pitfall 4: Scroll Restoration with Infinite Scroll — Book Not in DOM
**What goes wrong:** User clicks book #50 (loaded via infinite scroll page 3), goes to detail, returns. On home page mount, only 24 books are in DOM (page 1). `document.getElementById('book-slug')` returns null.
**Why it happens:** The target book may be on page 2 or 3, not yet fetched.
**How to avoid:** Along with `scrollTarget` in sessionStorage, also store the `cursor` that was active when the book was clicked. On mount, if `scrollTarget` exists, restore the cursor to pre-fetch enough pages before scrolling. Simpler fallback: just store the cursor and show the books from where the user left off.
**Warning signs:** Scroll jumps to top instead of the clicked book on navigation back.

### Pitfall 5: FOUC (Flash of Unstyled/Wrong Theme) on Dark Mode
**What goes wrong:** Page loads with light theme, then React mounts and switches to dark — visible flash.
**Why it happens:** React reads localStorage only after JS executes, which is after initial paint.
**How to avoid:** Add an inline `<script>` in `index.html` (before `</head>`) that reads localStorage and sets `data-theme` on `<html>` synchronously — before React loads.
**Warning signs:** Dark mode users see a light flash on every page load.

### Pitfall 6: Google Fonts Render Blocking
**What goes wrong:** `@import url('...')` in CSS blocks rendering until fonts load. Text invisible until fonts arrive.
**Why it happens:** CSS `@import` is synchronous and render-blocking.
**How to avoid:** Use `<link rel="preconnect" href="https://fonts.googleapis.com">` and `<link rel="stylesheet" href="...">` in `index.html` `<head>`.
**Warning signs:** Slow font rendering, layout shift score issues.

### Pitfall 7: Vite 8 + `npm create vite@latest` — Node Version
**What goes wrong:** Vite 8 requires Node.js ≥20. Older Node versions fail.
**Why it happens:** Vite 8 uses modern Node APIs.
**How to avoid:** Current environment has Node 20.14.0 ✓ — no action needed.

### Pitfall 8: `?raw` Vite Import for Markdown
**What goes wrong:** `import bioContent from '../content/bio.md'` without `?raw` tries to process it as JavaScript and fails.
**Why it happens:** Vite doesn't know to treat `.md` as raw text without the `?raw` suffix.
**How to avoid:** Always use `import bioRaw from '../content/bio.md?raw'`. Add TypeScript declaration: `declare module '*.md?raw' { const content: string; export default content; }`.

---

## Runtime State Inventory

> Step 2.5: SKIPPED — this is a greenfield frontend phase, not a rename/refactor/migration. No runtime state to inventory.

None — verified: `frontend/` directory contains only `dist/` (empty with `.gitkeep`) and `embed.go` (stub). No existing state to migrate.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|-------------|-----------|---------|----------|
| Node.js ≥20 | Vite 8 | ✓ | 20.14.0 | — |
| npm | Package management | ✓ | 10.1.0 | — |
| Go backend (port 8081) | Vite dev proxy | Conditional | N/A | Mock with MSW in tests |
| PostgreSQL | Book data (via API) | Conditional | N/A | Mock with MSW in tests |
| Google Fonts CDN | Playfair Display + Inter | ✓ (internet) | N/A | System font fallbacks in CSS |

**Missing dependencies with no fallback:** None that block the build. Go backend needed for end-to-end dev testing but not for component development.

**Missing dependencies with fallback:** Go API not running → use MSW in Vitest tests; use manual test data for visual development.

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `react-router-dom` v5 createBrowserHistory | v6 `<BrowserRouter>` / `createBrowserRouter` | React Router v6 (2021) | Must use `<Routes>`, `<Route>` components |
| TanStack Query v4 `isLoading` | v5 `isPending` | v5 (2023) | Skeleton display logic changes |
| TanStack Query v4 implicit `initialPageParam` | v5 explicit `initialPageParam: undefined` | v5 (2023) | Required field — missing causes TS error |
| `marked(text)` sync call (v1-v4) | `marked.parseSync(text)` in v17 | marked v5+ (2023) | Breaking — sync parse API changed |
| MSW v1 `rest.get()` handlers | v2 `http.get()` + `HttpResponse.json()` | MSW v2 (2023) | Different test setup/handler syntax |
| React 17/18 `ReactDOM.render()` | React 19 `createRoot()` | React 18 (2022) | Must use `createRoot` in main.tsx |

**Deprecated/outdated:**
- `marked.parse(text)` in synchronous context: now returns Promise in v17
- `react-router-dom` v5 patterns: `<Switch>`, `useHistory` removed in v6
- TanStack Query v4 `onSuccess`/`onError` callbacks on `useQuery`: removed in v5; use `useEffect` watching `data`/`error` instead

---

## Code Examples

### main.tsx Setup (React 19 + TanStack Query v5 + React Router v6)
```typescript
// main.tsx
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';       // v6
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import App from './App';
import './styles/global.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,   // 5 min — books don't change often
      retry: 2,
    },
  },
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>
);
```

### API Types (from Phase 2 contract)
```typescript
// api/types.ts — matches Phase 2 D-01, D-04 contract
export interface Author { name: string; slug: string; }
export interface Genre  { name: string; slug: string; }

export interface Book {
  slug: string;
  title: string;
  cover_path: string;     // e.g. "/covers/9780385490818.jpg"
  read_at: string | null;
  publication_year: number | null;
  authors: Author[];
  genres: Genre[];
}

export interface PaginatedBooks {
  items: Book[];
  next_cursor: string | null;
  has_more: boolean;
}
```

### MSW v2 Handlers (for Vitest)
```typescript
// src/mocks/handlers.ts — MSW v2 syntax (http.get, HttpResponse)
import { http, HttpResponse } from 'msw';
import { PaginatedBooks } from '../api/types';

export const handlers = [
  http.get('/api/books', () => {
    const response: PaginatedBooks = {
      items: [/* test data */],
      next_cursor: null,
      has_more: false,
    };
    return HttpResponse.json(response);
  }),

  http.get('/api/books/currently-reading', () => {
    return HttpResponse.json([/* test currently-reading books */]);
  }),
];
```

### Vitest Setup
```typescript
// vitest.config.ts
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
  },
});

// src/test/setup.ts
import { beforeAll, afterAll, afterEach } from 'vitest';
import { cleanup } from '@testing-library/react';
import { server } from './msw-server';

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));
afterEach(() => { cleanup(); server.resetHandlers(); });
afterAll(() => server.close());
```

---

## Validation Architecture

> `nyquist_validation: true` in `.planning/config.json` — validation section is required.

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Vitest 4.1.2 |
| Config file | `frontend/vitest.config.ts` — Wave 0 creates this |
| Quick run command | `cd frontend && npm test -- --run` |
| Full suite command | `cd frontend && npm test -- --run --coverage` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| HOME-01 | Bio section renders photo and text from bio.md | unit | `npm test -- --run src/components/Bio.test.tsx` | ❌ Wave 0 |
| HOME-02 | NowReadingSection renders books from API; hides when empty | unit | `npm test -- --run src/components/NowReadingSection.test.tsx` | ❌ Wave 0 |
| HOME-03 | BookGrid renders fetched books in grid | integration | `npm test -- --run src/components/BookGrid.test.tsx` | ❌ Wave 0 |
| HOME-04 | BookCard is a link navigating to /books/:slug | unit | `npm test -- --run src/components/BookCard.test.tsx` | ❌ Wave 0 |
| UI-01 | Intersection Observer triggers fetchNextPage | integration | `npm test -- --run src/hooks/useIntersectionObserver.test.ts` | ❌ Wave 0 |
| UI-02 | Load More button calls fetchNextPage; disabled while fetching | unit | `npm test -- --run src/components/BookGrid.test.tsx` | ❌ Wave 0 |
| UI-03 | sessionStorage scroll target set on card click; restored on mount | unit | `npm test -- --run src/hooks/useScrollRestoration.test.ts` | ❌ Wave 0 |
| UI-06 | Theme toggle reads prefers-color-scheme on first load; persists to localStorage; toggles data-theme | unit | `npm test -- --run src/hooks/useTheme.test.ts` | ❌ Wave 0 |
| UI-07 | CSS tokens correctly defined (smoke — check computed styles) | manual | visual inspection | N/A |
| UI-08 | Font loading: Playfair Display and Inter applied | manual | visual inspection | N/A |
| UI-09 | BookCover shows gradient placeholder before image loads / on error | unit | `npm test -- --run src/components/BookCover.test.tsx` | ❌ Wave 0 |
| UI-10 | Grid renders 2 columns at mobile viewport | integration | `npm test -- --run src/components/BookGrid.test.tsx` | ❌ Wave 0 |

### Testing Notes

- **Intersection Observer in jsdom:** jsdom does not implement Intersection Observer. Must mock it in test setup: `global.IntersectionObserver = vi.fn(() => ({ observe: vi.fn(), disconnect: vi.fn(), unobserve: vi.fn() }))`.
- **MSW v2 is required** for any test that exercises `useInfiniteQuery` or `useQuery` — mock the Go API endpoints.
- **No snapshot tests:** CSS-heavy components with custom properties are brittle for snapshots. Prefer behavior assertions.
- **Theme toggle test:** Mock `window.matchMedia` in test setup for `prefers-color-scheme` queries.
- **TanStack Query in tests:** Wrap test components in a fresh `QueryClientProvider` with `QueryClient({ defaultOptions: { queries: { retry: false } } })` to avoid retry delays.

### Sampling Rate
- **Per task commit:** `cd frontend && npm test -- --run` (all unit tests, ~5-10s)
- **Per wave merge:** `cd frontend && npm test -- --run --coverage`
- **Phase gate:** Full suite green before `/gsd-verify-work`

### Wave 0 Gaps
- [ ] `frontend/vitest.config.ts` — test environment configuration
- [ ] `frontend/src/test/setup.ts` — MSW server setup, IntersectionObserver mock, matchMedia mock
- [ ] `frontend/src/test/msw-server.ts` — MSW server instance
- [ ] `frontend/src/mocks/handlers.ts` — API handler mocks for `/api/books` and `/api/books/currently-reading`
- [ ] Framework install: `npm install -D vitest @testing-library/react @testing-library/user-event msw @vitest/coverage-v8 jsdom`
- [ ] Add `"test": "vitest"` script to `frontend/package.json`

---

## Open Questions

1. **marked v17 — `parseSync` XSS sanitization**
   - What we know: `bio.md` is a bundled local file — not user input. XSS risk is zero.
   - What's unclear: Whether marked v17 has any built-in sanitization that could strip valid HTML tags used in bio.
   - Recommendation: Use `marked.parseSync()` without sanitization for bio.md. Document the assumption that bio.md is trusted content.

2. **Vite `?raw` imports — TypeScript declaration**
   - What we know: Vite supports `?raw` suffix for raw text import. TypeScript needs a `.d.ts` declaration.
   - What's unclear: Whether `@vitejs/plugin-react` already provides these types or a manual declaration is needed.
   - Recommendation: Add `/// <reference types="vite/client" />` to `src/vite-env.d.ts` — this includes `?raw` typings in Vite 4+.

3. **Scroll restoration — deep pagination scenario**
   - What we know: If user scrolled to book #100 and clicks through, returning user will only see books 1-24 initially.
   - What's unclear: Whether we restore the scroll cursor (fetch pages 1-N) or accept that the position will be lost for deeply-paginated items.
   - Recommendation: For Phase 3, implement basic scroll restoration (save slug to sessionStorage). If slug not found in initial 24 books, skip restoration. Document as a known limitation. The complex multi-page restoration can be a v2 enhancement.

---

## Sources

### Primary (HIGH confidence)
- npm registry — `@tanstack/react-query@5.96.2`, `react-router-dom@6.30.3`/`7.14.0`, `vite@8.0.3`, `marked@17.0.5`, `lucide-react@1.7.0`, `vitest@4.1.2`, `msw@2.12.14` — all verified via `npm view` 2026-04-03
- `.planning/phases/03-frontend-core/03-CONTEXT.md` — locked implementation decisions
- `.planning/phases/03-frontend-core/03-UI-SPEC.md` — all component contracts, resolved agent discretion items
- `.planning/phases/02-go-rest-api/02-CONTEXT.md` — API response shapes (D-01/D-02/D-04)
- MDN Intersection Observer API — `IntersectionObserver({ threshold: 0.1 })`
- MDN `prefers-color-scheme` media query

### Secondary (MEDIUM confidence)
- TanStack Query v5 migration guide (training knowledge, verified via version number): `initialPageParam` required, `isPending` replaces `isLoading`
- marked changelog (training knowledge, v5+ made async default): `parseSync()` for sync use in v17
- MSW v2 handler syntax (training knowledge, verified via version number): `http.get()` + `HttpResponse.json()`

### Tertiary (LOW confidence — flag for validation)
- React Router v7 breaking changes: training knowledge — verify when pinning v6 that `<BrowserRouter>` and `<Routes>` APIs are unchanged in v6.30.x

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all versions verified against npm registry 2026-04-03
- Architecture: HIGH — patterns derived from UI-SPEC.md (authoritative) + verified API shapes
- TanStack Query v5 API: HIGH — version confirmed, key breaking changes well-documented
- marked v17 async API: MEDIUM — version confirmed, sync/async behavior inferred from known v5 migration pattern
- Pitfalls: HIGH — derived from verified version differences and known ecosystem changes

**Research date:** 2026-04-03
**Valid until:** 2026-05-03 (30 days — stable libraries, but marked/react-router could patch)
