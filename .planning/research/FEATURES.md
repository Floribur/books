# Feature Landscape: Flo's Library

**Domain:** Personal book showcase / reading history website
**Researched:** 2026-03-31
**Confidence note:** Web search and WebFetch tools were unavailable during this session. All findings are from training data (cutoff August 2025). Confidence levels reflect this constraint.

---

## Table Stakes

Features users (and the owner) expect. Missing = product feels incomplete.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Book grid with cover images | Visually the core of any library site | Low | Aspect ratio + lazy loading critical |
| Book detail page | Users click covers expecting full info | Medium | Metadata selection matters (see below) |
| "Now Reading" section | Goodreads-style expectation for active readers | Low | Prominent placement on home page |
| Author index | Natural navigation path from a book detail | Low | Alphabetical + book count sort |
| Genre index | Natural browsing path for visitors | Low | Sorted by book count descending |
| Reading Challenge / Year view | Goodreads popularized this; readers expect it | Medium | Year selector, stats summary, book grid per year |
| Responsive layout | Table stakes in 2025 | Low | Mobile grid collapses to 2-col |
| Dark/light mode | Strong user expectation for personal sites | Medium | See dark mode section below |

## Differentiators

What makes Flo's Library feel curated and personal rather than a Goodreads clone.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Sidebar with animated book-reading graphic | Brand continuity from prior version; personality | Medium | CSS animation preferred over Lottie (see below) |
| Self-hosted covers | No broken images over time | Low (infra) | Download pipeline already planned |
| Bio section on home page | Makes it a personal site, not just a database | Low | Photo, passions summary |
| Clean URL structure | /books/the-name-of-the-rose feels intentional | Low | Slug-based routing |
| Elegant typography | Library aesthetic sets tone immediately | Low | Font pairing matters (see below) |

## Anti-Features

Features to explicitly NOT build for this project.

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| Star ratings displayed on site | Source of truth is Goodreads; stale/confusing | Link to Goodreads profile for opinions |
| User reviews / comments | Out of scope; adds auth/moderation complexity | Omit entirely |
| Search box | Small corpus (~hundreds of books) doesn't need it | Browser Ctrl+F sufficient; add later if needed |
| Pagination (numbered pages) | Conflicts with stated requirement of infinite scroll | Infinite scroll with "load more" fallback |
| Social sharing widgets | Add JS weight; personal site doesn't need it | Open Graph meta tags cover share previews |
| Reading progress % bar | Requires real-time sync; Goodreads doesn't expose granular progress | Show "Currently Reading" status only |

---

## 1. Reading Challenge / Year-in-Books Page

**Confidence: MEDIUM** (pattern knowledge from Goodreads, StoryGraph, similar sites)

### What Goodreads Does (Reference Pattern)

Goodreads Year in Books shows:
- Goal vs actual count with a large progress visualization
- Month-by-month breakdown
- Longest / shortest book callouts
- Most popular genres of the year
- A scrollable grid of all books read that year

### Recommended Approach for Flo's Library

Keep it simpler and more elegant than Goodreads:

**Header section** (sticky or prominent):
- Year selector (prev/next arrows, or dropdown) — default to current year
- Single stat: "X books read in [year]" in large type

**Stats strip** (horizontal row of 3–4 cards):
- Books read this year
- Pages read (if page count data is available via Google Books)
- Most-read genre
- Longest book

**Book grid below**: same grid component as the main library, filtered by year. No need for a separate component — reuse with a year filter prop.

**Implementation**: This is a client-side filtered view if all books are fetched at once. If the corpus is large (500+ books), add a `/api/books?year=2024` endpoint. For a personal library, client-side filtering is fine through ~300 books.

**Year navigation**: Store selected year in URL query param (`?year=2023`) so links are shareable.

---

## 2. Book Detail Page

**Confidence: HIGH** (well-established patterns; Google Books API metadata is well-documented)

### Metadata to Display (Priority Order)

| Field | Source | Display Priority | Notes |
|-------|--------|-----------------|-------|
| Cover image (large) | Google Books / self-hosted | Critical | Full-size, not thumbnail |
| Title | Goodreads/Google Books | Critical | H1 |
| Author(s) | Goodreads/Google Books | Critical | Link to author index page |
| Description / synopsis | Google Books | High | Truncate at ~300 chars with "show more" |
| Genres | Google Books categories | High | Tag/pill style, each links to genre page |
| Published year | Google Books | Medium | Not full date — year is sufficient |
| Page count | Google Books | Medium | "432 pages" |
| ISBN | Google Books | Low | Useful for uniqueness but not display-critical |
| Date I read it | Goodreads shelf data | Medium | "Read in March 2024" style |
| External link | Goodreads book page | Low | "View on Goodreads" link |

### Layout Pattern

Two-column layout on desktop, single column on mobile:

```
[Cover image — left ~35%]  [Title (H1)
                            Author (linked)
                            Year · Pages · Genre pills
                            ─────────────────────────
                            Description (expandable)
                            ─────────────────────────
                            "Read in [Month Year]"
                            [View on Goodreads →]]
```

This is the standard pattern used by Goodreads, StoryGraph, and OpenLibrary. It works because the cover is the primary visual anchor and metadata reads naturally left-to-right.

**"You might also like" / related books**: Skip for MVP. Adds algorithmic complexity with minimal value for a personal showcase.

---

## 3. Infinite Scroll vs Pagination

**Confidence: HIGH** (NNGroup and WCAG guidance is stable; React patterns are established)

### Recommendation: Infinite Scroll with "Load More" Button Fallback

Pure infinite scroll (auto-trigger on scroll) has known accessibility problems:
- Screen readers lose position context
- Users can't bookmark a position in the list
- Footer becomes unreachable

**Best practice (2024–2025 consensus)**: Use a "Load More" button instead of auto-triggering. This gives the user control, is keyboard-navigable, and keeps the footer reachable. It still feels progressive and avoids full pagination UI.

However, the PROJECT.md states infinite scroll as a requirement. The pragmatic middle ground:

**Implement "Load More" with Intersection Observer as the trigger**, but also show the button explicitly. This satisfies the "infinite scroll" UX feel while keeping accessibility viable.

### React Implementation Pattern

```tsx
// Use Intersection Observer on a sentinel element at the bottom of the list
const sentinelRef = useRef<HTMLDivElement>(null);

useEffect(() => {
  const observer = new IntersectionObserver(
    (entries) => {
      if (entries[0].isIntersecting && hasMore && !isLoading) {
        loadMore();
      }
    },
    { threshold: 0.1 }
  );
  if (sentinelRef.current) observer.observe(sentinelRef.current);
  return () => observer.disconnect();
}, [hasMore, isLoading]);

// Render sentinel + explicit button as fallback
<div ref={sentinelRef} aria-hidden="true" />
<button onClick={loadMore} disabled={isLoading}>
  {isLoading ? 'Loading...' : 'Load more books'}
</button>
```

**Library option**: `@tanstack/react-query` with `useInfiniteQuery` handles cursor-based pagination cleanly and is the standard approach for React apps in 2025. It manages loading states, error states, and cache automatically.

**Initial page size**: 24 books (fits cleanly in 4-col, 3-col, and 2-col grids without orphaned rows).

### Accessibility Requirements

- The "Load more" button must be keyboard-focusable
- After loading, move focus to the first newly-added book card (`aria-setsize` if using ARIA list)
- Announce load count to screen readers: `aria-live="polite"` region saying "12 more books loaded"
- Never hide the footer behind auto-triggering scroll

---

## 4. Author Index / Genre Index Pages

**Confidence: HIGH** (standard patterns, no library-specific dependencies)

### Author Index

**Sort options** (implement in this priority order):
1. Alphabetical by last name — default, expected
2. By book count (descending) — useful "most-read authors" view

**Display**: A list (not grid) works better for author index since there are no author cover images. Each entry:
```
[Author Name]          [N books]
```
Clicking navigates to the author detail page (all books by that author).

**Author detail page**: Same grid layout as the main library, filtered. Header: "Books by [Author Name]" with the book count. No author bio needed (would require Wikipedia scraping — skip for MVP).

### Genre Index

**Sort**: By book count descending (most-read genre first) — more useful than alphabetical for genres.

**Display**: Can use a tag cloud or simple sorted list. Tag cloud is visually interesting but can be misleading about relative counts. Recommendation: **sorted list with book count**, optionally with a subtle bar chart visualization using just CSS width (no chart library needed).

```
Fantasy         ████████████  42
Science Fiction ████████      28
Biography       ████          14
```

This is achievable with pure CSS (width as percentage of max count) and is more informative than equal-sized tags.

### URL Structure

- `/authors` — index
- `/authors/ursula-k-le-guin` — author detail (slug from name)
- `/genres` — index
- `/genres/science-fiction` — genre detail

Slugs should be generated consistently: lowercase, spaces to hyphens, remove punctuation.

---

## 5. Book Cover Display

**Confidence: HIGH** (browser APIs stable; aspect ratio CSS is well-established)

### Aspect Ratio

Book covers are not uniform, but the vast majority of modern books follow roughly **2:3 portrait ratio** (width:height). Lock all cover containers to `aspect-ratio: 2/3` in CSS. This prevents layout shift when covers load and creates a consistent grid.

```css
.book-cover-container {
  aspect-ratio: 2 / 3;
  overflow: hidden;
  background-color: var(--cover-placeholder-bg);
}

.book-cover-container img {
  width: 100%;
  height: 100%;
  object-fit: cover; /* crops slight mismatches gracefully */
}
```

### Lazy Loading

Use the **native `loading="lazy"` attribute** on `<img>` tags. Browser support is near-universal (97%+ as of 2025). No JavaScript library needed for this.

```tsx
<img
  src={cover.url}
  alt={`Cover of ${book.title}`}
  loading="lazy"
  decoding="async"
  width={200}
  height={300}
/>
```

Always include explicit `width` and `height` to prevent cumulative layout shift (CLS), which is a Core Web Vitals metric.

### Placeholder Pattern for Missing Covers

When no cover image is available (Google Books has no image for ~10-15% of books):

**Recommended: CSS-only generated placeholder** rather than a generic "no image" icon.

```css
.cover-placeholder {
  background: linear-gradient(135deg, var(--primary-dark) 0%, var(--primary) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  text-align: center;
  color: white;
  font-size: 0.75rem;
  font-weight: 600;
}
```

Show the book title text centered on the colored gradient. This matches the brand color and looks intentional rather than broken.

**Alternative**: Generate a SVG placeholder server-side with the title embedded — cleaner for SSR but more complex. Skip for MVP; use CSS approach.

### Cover Quality

Google Books API returns thumbnails (`thumbnail`) and small thumbnails (`smallThumbnail`). For detail pages, `thumbnail` is sufficient (~128px wide). For a high-quality grid display, request the `zoom=1` parameter in Google Books which can return larger images, or strip the `zoom` parameter to get the full-size version. Test per book — quality varies.

---

## 6. Sidebar Navigation with Animated Book-Reading Graphic

**Confidence: MEDIUM** (CSS animation recommendation is well-grounded; Lottie comparison is established)

### Recommendation: CSS Animation, Not Lottie

**Do not use Lottie** for this use case. Reasons:
- Lottie requires the `lottie-web` library (~60kb gzipped) for a single decorative animation
- Lottie files require a designer to create the source animation in After Effects or a compatible tool
- A "reading book" animation is simple enough to implement purely in CSS

**CSS approach**: A small SVG of an open book with CSS keyframe animation is the right tool. The animation can show:
- Pages turning (a simple CSS `rotateY` transform on a page element)
- Or a gentler "breathing" scale effect to suggest reading without being distracting

A page-turn effect with pure CSS:

```css
@keyframes page-turn {
  0%   { transform: rotateY(0deg); }
  45%  { transform: rotateY(-180deg); }
  55%  { transform: rotateY(-180deg); }
  100% { transform: rotateY(0deg); }
}

.book-page-right {
  transform-origin: left center;
  animation: page-turn 3s ease-in-out infinite;
  animation-delay: 0.5s;
}
```

**Alternative**: Use an inline SVG with CSS `animation` on path/group elements. SVG animation (SMIL is deprecated; use CSS keyframes on SVG elements) gives more control for a book graphic.

**React implementation**: The book SVG lives as a React component in the sidebar. `prefers-reduced-motion` must be respected:

```css
@media (prefers-reduced-motion: reduce) {
  .book-animation * {
    animation: none;
  }
}
```

**If the project has a designer producing assets**: Lottie becomes viable at that point. But for a self-built personal site, CSS animation is more maintainable, zero-dependency, and fast.

---

## 7. Dark / Light Mode for a Library Aesthetic

**Confidence: HIGH** (CSS custom properties + prefers-color-scheme is the stable standard)

### Implementation Approach

Use CSS custom properties (variables) as the single source of truth for all colors. Toggle a `data-theme` attribute on `<html>`.

```css
:root {
  /* Light mode (default) */
  --bg-primary: #faf8f5;        /* warm off-white, like aged paper */
  --bg-secondary: #f0ece4;      /* slightly darker warm white */
  --text-primary: #1a1208;      /* very dark warm brown, not pure black */
  --text-secondary: #5c4a3a;    /* medium warm brown */
  --surface: #ffffff;
  --border: #d4c8b8;            /* warm gray */
  --accent: #6d233e;            /* brand primary */
  --accent-hover: #8a2d4f;
}

[data-theme="dark"] {
  --bg-primary: #1a1208;        /* very dark warm brown */
  --bg-secondary: #241a0e;
  --text-primary: #f0ece4;
  --text-secondary: #c4b8a8;
  --surface: #2d2015;
  --border: #3d2e1e;
  --accent: #c45a7a;            /* lighter accent for dark bg contrast */
  --accent-hover: #d46a8a;
}
```

**Auto-detect OS preference** on first load:
```tsx
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
document.documentElement.setAttribute('data-theme', prefersDark ? 'dark' : 'light');
```

**Persist user choice** in `localStorage`. Don't force the OS preference if the user has explicitly chosen.

### Library Aesthetic Notes

- **Light mode**: Warm paper tones (`#faf8f5`, not pure white) evoke physical book pages. This is the detail that separates a "library" feel from a generic app.
- **Dark mode**: Very dark warm brown (`#1a1208`) rather than pure black (`#000000`) maintains the library warmth in dark mode.
- **Avoid cool grays** in both modes — they conflict with the warm burgundy brand color.
- **Cover images**: In dark mode, add a very subtle CSS shadow around book covers to prevent them from floating on dark backgrounds: `box-shadow: 0 2px 8px rgba(0,0,0,0.4)`.

---

## 8. SEO Considerations

**Confidence: HIGH** (Next.js SSR patterns are well-established; personal site SEO is stable knowledge)

### Should Book Pages Be Server-Side Rendered?

**Yes, book detail pages should be SSR or SSG** (static site generation). Reasons:

1. Search engines index individual book pages (someone searching "Florian reading The Name of the Rose" should find it)
2. Open Graph / Twitter Card meta tags (for social sharing previews) must be in the initial HTML — client-rendered meta tags are often missed by social scrapers
3. Core Web Vitals (LCP) is better with SSR — the cover image URL is in initial HTML, not injected after JS runs

**However**, this site uses React + Go (not Next.js). This creates a tension:

- **Option A**: Use Next.js for the frontend (React meta-framework with built-in SSR/SSG) — recommended if SEO is important
- **Option B**: Keep Vite/CRA React but add server-side rendering via Go template rendering for the initial HTML shell + meta tags only
- **Option C**: Accept client-side only rendering; add a prerender middleware (like `rendertron` or similar) — complex

**Recommendation**: If the tech stack allows flexibility, adopt Next.js instead of plain React. The SSR/SSG benefit for book pages (Open Graph previews, Google indexing) is significant. If the React + Go split is firm, implement at minimum a Go endpoint that returns page-specific Open Graph meta tags in the HTML head for the `/books/:slug` route.

### Minimum SEO Requirements

For each book detail page, the server must return in initial HTML:

```html
<title>The Name of the Rose — Flo's Library</title>
<meta name="description" content="[First 150 chars of description]" />
<meta property="og:title" content="The Name of the Rose" />
<meta property="og:description" content="[Description]" />
<meta property="og:image" content="https://[domain]/covers/the-name-of-the-rose.jpg" />
<meta property="og:type" content="book" />
<link rel="canonical" href="https://[domain]/books/the-name-of-the-rose" />
```

For the home page, author pages, and genre pages: standard title/description meta is sufficient.

**Structured data (JSON-LD)**: Optional but valuable for book pages. Google supports `Book` schema type:

```json
{
  "@context": "https://schema.org",
  "@type": "Book",
  "name": "The Name of the Rose",
  "author": { "@type": "Person", "name": "Umberto Eco" },
  "numberOfPages": 502,
  "isbn": "9780156001311"
}
```

Add this as a `<script type="application/ld+json">` in the page head. Google uses it for rich snippets.

---

## 9. Design Inspiration for Personal Book Library Sites

**Confidence: MEDIUM** (site-specific examples from training knowledge; may have evolved)

### Reference Sites Worth Studying

**StoryGraph (app.thestorygraph.com)**
- Strong library grid with clear typography hierarchy
- Good genre/mood filtering approach
- Their "stats" page is a reference for the Reading Challenge view

**Literal.club**
- Beautiful, minimal book display
- Strong typographic focus
- Clean author/genre navigation

**OpenLibrary (openlibrary.org)**
- Dense but functional book catalog patterns
- Good cover grid reference

**Personal sites to look for**: Search "bookshelf website personal" or "reading log site" on GitHub Pages / Netlify showcase — many developers have built these and some have excellent taste.

### What Sets Beautiful Book Sites Apart (Pattern Summary)

1. **Typography does more work than color** — a good serif for titles, clean sans-serif for metadata
2. **Generous whitespace** — books feel precious when given room
3. **Consistent cover grid** — uniform aspect ratios create visual calm
4. **Minimal navigation** — book sites do not need complex nav; sidebar or simple top nav suffices
5. **Cover image quality** — a blurry or stretched cover ruins the aesthetic

---

## 10. Color Palette for #6d233e (Dark Burgundy/Wine Red)

**Confidence: HIGH** (color theory is stable; specific hex values are calculated)

### Primary Color Analysis

`#6d233e` is a **dark desaturated wine red** (HSL approximately 339°, 51%, 28%). It is deep, warm, and associated with leather-bound books, wine, and academic spaces — perfect for a library.

### Recommended Complementary Palette

```
Primary:        #6d233e   Wine red (brand)
Primary light:  #8a2d4f   Lighter wine (hover states)
Primary dark:   #4a1728   Deeper wine (pressed states, borders)

Accent warm:    #c4843a   Antique gold/amber
                          (HSL ~30°, 55%, 50% — split-complementary to #6d233e)
                          Use for: call-to-action highlights, year badges, stats

Neutral warm:   #5c4a3a   Dark warm brown (body text)
Neutral mid:    #a08878   Medium warm taupe (secondary text, borders)
Neutral light:  #f0ece4   Warm off-white (page background)
Neutral paper:  #faf8f5   Near-white with warmth (card backgrounds)
```

**The gold accent (#c4843a)** is the key complementary color. It appears in library contexts naturally (gold lettering on book spines, candlelight, wood shelves) and creates strong contrast against the burgundy without clashing.

**Avoid**: Cool blue/green complements (they fight the warmth), pure orange (too aggressive), or gray (lifeless against this palette).

### Dark Mode Palette Extension

```
Dark bg:        #1a1208   Very dark warm brown (not pure black)
Dark surface:   #2d2015   Slightly lighter for cards
Dark accent:    #c45a7a   Lighter wine (primary accent on dark bg, contrast-safe)
Dark gold:      #d4944a   Slightly lighter gold for dark backgrounds
```

### Typography Pairing Recommendation

For a library aesthetic:
- **Headings**: A serif — `Playfair Display`, `Lora`, or `Libre Baskerville` (all available on Google Fonts, free)
- **Body/metadata**: `Inter` or `Source Sans 3` — clean, highly legible sans-serif
- **Accent/labels**: The same sans-serif at small caps or uppercase with letter-spacing

`Playfair Display` (headings) + `Inter` (body) is a widely-used combination that feels literary without being precious.

---

## Feature Dependencies

```
Book detail page → Cover self-hosting pipeline
Author index → Author detail page → Book detail page
Genre index → Genre detail page → Book detail page
Reading Challenge page → Year-filtered book grid (reuse book grid component)
Infinite scroll → Book grid component
Dark/light mode → CSS custom properties (foundation; all other components depend on it)
```

## MVP Recommendation

**Build in this order:**

1. CSS foundation (custom properties, dark/light mode, typography) — everything depends on this
2. Book grid component with lazy-loaded covers and placeholder pattern
3. Home page: bio section + "Now Reading" + "Books Read" grid with infinite scroll
4. Book detail page (the most content-rich single page)
5. Author index + Author detail pages (reuses book grid)
6. Genre index + Genre detail pages (reuses book grid)
7. Reading Challenge page (reuses book grid with year filter)
8. Sidebar with animated book graphic (can be added without affecting data flow)

**Defer for post-MVP:**
- JSON-LD structured data (low-effort later addition)
- Stats visualizations beyond basic counts on Reading Challenge page
- Progressive Web App / offline support

---

## Sources

**Note:** Web search and WebFetch were unavailable during this research session. All findings are from training data (cutoff August 2025). Confidence levels reflect this:

- NNGroup infinite scroll UX guidance — HIGH confidence (stable, well-documented research; last known publication on the topic ~2022 still current)
- CSS `loading="lazy"` — HIGH confidence (browser API, widely documented)
- CSS custom properties for theming — HIGH confidence (stable CSS standard, universal support)
- Intersection Observer API — HIGH confidence (stable Web API since ~2019, universal support)
- `@tanstack/react-query` `useInfiniteQuery` — MEDIUM confidence (API shape as of v5; verify current docs)
- Color palette recommendations — HIGH confidence (color theory; specific hex values are calculated)
- StoryGraph / Literal.club site descriptions — MEDIUM confidence (sites may have redesigned since training)
- Lottie vs CSS animation recommendation — HIGH confidence (bundle size and tooling tradeoffs are factual)
- Google Books API cover URL parameters — MEDIUM confidence (API behavior may have changed; verify with actual API calls during implementation)
- Next.js SSR recommendation — HIGH confidence (well-established pattern; the tension with the Go+React stack is a real architectural decision that needs resolution)
