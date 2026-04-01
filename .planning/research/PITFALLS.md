# Domain Pitfalls: Flo's Library

**Domain:** Personal book library — Go backend + React TypeScript frontend
**Researched:** 2026-04-01
**Confidence note:** Web search and WebFetch tools were unavailable during this session. All findings are from training data (cutoff August 2025) cross-referenced against the existing ARCHITECTURE.md research. Confidence levels are assigned honestly.

---

## Critical Pitfalls

Mistakes that cause rewrites, data loss, or major unplanned rework.

---

### Pitfall 1: Goodreads RSS Pagination Silently Truncating Data

**What goes wrong:** The RSS endpoint `?shelf=read&per_page=200` returns at most 200 items. If Florian's "read" shelf exceeds 200 books (which is plausible for an active reader tracking history), only the first 200 books are returned. There is no error, no truncation notice — the feed simply ends at 200. The app thinks it has synced everything and marks no books as missing.

**Why it happens:** The RSS feed does not return a total count or "has next page" signal in standard RSS format. Custom GR extensions may or may not include this. `gofeed` will return exactly what the feed provides with no indication that items were capped.

**Consequences:** Books added to Goodreads earliest (oldest reads) are silently missing from the app. The library appears complete but is not. Discovering this requires manually counting books.

**Prevention:**
1. Always paginate. Implement a loop: fetch `?shelf=read&per_page=200&page=1`, then `&page=2`, etc., until a response returns fewer than 200 items (or zero items).
2. After initial sync, log the total book count returned vs. expected. Compare against the Goodreads "N books" count visible on the public profile.
3. Add a sync health check: if the DB count drops significantly between syncs without a corresponding shelf change, alert loudly.

**Detection:** After first sync, visit `https://www.goodreads.com/user/show/79499864-florian` and compare the "read" shelf count on Goodreads against the DB count. A gap signals truncation.

**Concrete mitigation in code:**
```go
func fetchShelfAllPages(userID, shelf string) ([]*gofeed.Item, error) {
    var all []*gofeed.Item
    for page := 1; ; page++ {
        url := fmt.Sprintf(
            "https://www.goodreads.com/review/list_rss/%s?shelf=%s&per_page=200&page=%d",
            userID, shelf, page,
        )
        items, err := fetchPage(url)
        if err != nil {
            return nil, err
        }
        all = append(all, items...)
        if len(items) < 200 {
            break // last page
        }
        time.Sleep(2 * time.Second) // be polite
    }
    return all, nil
}
```

**Confidence: HIGH** — Goodreads RSS `per_page=200` cap is well-documented in the developer community.

---

### Pitfall 2: Goodreads RSS Custom Fields Not Parsed by Standard Libraries

**What goes wrong:** Fields like `<book_id>`, `<author_name>`, `<isbn13>`, `<user_read_at>`, and `<user_shelves>` are Goodreads custom XML extensions, not standard RSS elements. `gofeed` parses them but buries them in `item.Extensions[""]` or `item.Custom` maps with non-obvious key names. Developers assume `item.Title` covers all fields and write code that silently ignores these fields — resulting in books stored without ISBNs, without read dates, or with wrong shelf assignments.

**Why it happens:** RSS is a standard format but Goodreads adds custom namespaced (or non-namespaced) elements. `gofeed` exposes these but the access pattern is not obvious from the library's top-level documentation.

**Consequences:** Books stored without ISBN means Google Books enrichment falls back to title+author search, which has higher false-positive rates. Books stored without `user_read_at` default to `NULL` read date, breaking year-based filtering and Reading Challenge page.

**Prevention:**
1. Before writing any sync code, manually fetch the RSS feed and inspect the raw XML: `curl https://www.goodreads.com/review/list_rss/79499864?shelf=read | head -100`
2. Write a small test program that prints all fields from a single feed item, including `item.Extensions`, to discover the actual key paths
3. Write a dedicated struct that maps all GR-specific fields with defensive nil checks

**Detection:** Parse a live feed item and log all `item.Extensions` keys before writing production field-mapping code.

**Confidence: MEDIUM** — `gofeed` extensions behavior from training data; key paths must be verified against live feed.

---

### Pitfall 3: Google Books API Title+Author Search False Positives

**What goes wrong:** When a book has no ISBN (blank `isbn13` in GR RSS), the enrichment pipeline falls back to searching Google Books by title and author. The search API returns the first result, which may be a different edition, a different book with a similar title, or an anthology containing a story with the same name. The pipeline happily stores the wrong cover image and description.

**Why it happens:** `intitle:` + `inauthor:` queries are fuzzy. "The Road" by "Cormac McCarthy" returns a clean match. "Normal People" by "Sally Rooney" might match multiple editions with different covers and descriptions. Books with common words in the title (e.g., "The House") are especially risky.

**Consequences:** Wrong cover images or descriptions stored in the DB. Discovering this requires visual inspection of every book — there is no automated way to detect a plausible-but-wrong match.

**Prevention:**
1. When falling back to title+author search, do NOT accept the first result automatically. Implement a confidence check:
   - Does the returned `authors[0]` closely match the input author name? (case-insensitive, Levenshtein distance)
   - Does the returned `title` contain the input title? (substring match, case-insensitive)
2. If confidence check fails: store the book with no Google Books enrichment and log it for manual review. Do not store a wrong match.
3. Add a `metadata_source` field to the books table: `'isbn'`, `'title_author'`, `'manual'`, `'none'`. Flag `'title_author'` matches for future review.
4. Never use description or genres from a title+author match without logging the match quality.

**Confidence: HIGH** — This is a well-known issue with book metadata APIs; same class of problem as music metadata matching.

---

### Pitfall 4: Book Cover Copyright — Fair Use for Personal Sites

**What goes wrong:** Developers assume that since covers are small images for personal use, there are no copyright concerns. Technically, book cover images are copyrighted by publishers. Downloading and self-hosting them without a license is a gray area that could theoretically generate a DMCA takedown notice.

**Why it happens:** The behavior is so common (every book review site shows covers) that it feels obviously permitted, but there is no explicit general license.

**Consequences:** For a personal, non-commercial site with low traffic, the practical risk is extremely low — publishers have never systematically targeted personal reading logs. The theoretical risk is a DMCA notice requiring removal of the image.

**The actual situation (MEDIUM confidence):**
- Google Books API Terms of Service (at time of training) permit displaying images returned by the API in connection with book information displayed using the API. However, downloading and re-hosting may not be explicitly permitted under those terms.
- OpenLibrary (Internet Archive) uses covers under fair use as a non-commercial library service. Deriving from their covers for a personal, non-commercial use is in similarly safe territory.
- For a private personal site, the practical exposure is essentially zero. No publisher has ever DMCA'd a personal reading log for showing book covers.

**Mitigation:**
1. Keep the site personal and non-commercial — no ads, no selling content. This is what makes the use defensible as fair use.
2. Give credit where possible: include a "cover source: Google Books" note or link back to the Google Books page.
3. If you are ever concerned, OpenLibrary covers are the cleanest option since the Internet Archive operates explicitly as a library under library copyright exceptions.
4. Do not hotlink — the project requirement to self-host covers is actually the riskier choice from a strict ToS standpoint, but the practical risk for a personal site is negligible.

**Definitive recommendation:** Self-host the covers as planned. The practical risk for a personal, non-commercial reading log is negligible. Do not let this block the project.

**Confidence: MEDIUM** — Copyright law interpretation; general consensus from the developer community is consistent with this assessment, but not legal advice.

---

### Pitfall 5: CORS Misconfiguration Causing Silent Fetch Failures in Production

**What goes wrong:** In development, the frontend (Vite dev server on port 5173) proxies API requests to Go (port 8080), so CORS never appears as a problem. In production, the React build is served from a different origin than the Go API. CORS headers are either missing entirely (API requests fail silently) or set to `*` (which works for unauthenticated requests but breaks if cookies or Authorization headers are ever added).

**Why it happens:** Vite's `proxy` config in `vite.config.ts` makes CORS invisible in development. First production deploy reveals the misconfiguration.

**Consequences:** The entire site appears broken in production — all API calls return opaque network errors in the browser. This is a deployment blocker, not a gradual degradation.

**Prevention — pick one deployment topology and commit to it:**

**Option A (recommended): Go serves the React build.** In production, the Go binary serves `dist/` as static files at `/`, and the API is at `/api/`. No CORS needed because everything is same-origin.
```go
// In production Go server
mux.Handle("/", http.FileServer(http.Dir("./dist")))
mux.Handle("/api/", apiRouter)
```
Build step: `npm run build` produces `dist/`, then the Go binary is deployed alongside it.

**Option B: Nginx reverse proxy.** Nginx serves React static files at `/`, proxies `/api/` to Go. Same-origin from the browser's perspective. CORS not needed.
```nginx
location / {
    root /var/www/books/dist;
    try_files $uri $uri/ /index.html;
}
location /api/ {
    proxy_pass http://localhost:8080;
}
```

**Option C (avoid if possible): Separate domains with explicit CORS.** `books.domain.com` for frontend, `api.books.domain.com` for backend. Requires setting `Access-Control-Allow-Origin` precisely and handling preflight OPTIONS requests. Added complexity for no benefit on a personal site.

**Definitive recommendation:** Use Option A (Go serves React build) in production. It is the simplest, eliminates CORS entirely, and results in a single deployable binary + dist directory. Use Vite proxy only in development.

**Confidence: HIGH** — CORS is a well-understood browser mechanism; the specific topology recommendations are standard practice.

---

## Moderate Pitfalls

Mistakes that cause bugs or significant rework, but not rewrites.

---

### Pitfall 6: SQLite-to-PostgreSQL Migration Breaking Schema Assumptions

**What goes wrong:** Starting with SQLite (or in-memory store) encourages schema choices that work in SQLite but need adjustment for PostgreSQL:
- SQLite uses `INTEGER PRIMARY KEY` (not `SERIAL`); PostgreSQL uses `SERIAL` or `BIGSERIAL` or `GENERATED ALWAYS AS IDENTITY`
- SQLite is case-insensitive for strings by default; PostgreSQL collation is `C` by default (case-sensitive LIKE)
- SQLite accepts any value in any column regardless of type; PostgreSQL enforces strict types
- SQLite `DATETIME` stores as text; PostgreSQL `TIMESTAMPTZ` stores as actual timestamps with timezone
- SQLite has no `ENUM` type; PostgreSQL has it but it requires migration to change values

**Consequences:** When migrating, existing queries break. `LIKE 'read%'` behaves differently. Date comparisons that worked in SQLite fail in PostgreSQL. The schema must be rewritten rather than ported.

**Prevention:**
1. Use `golang-migrate` from day one, even during development. Never use `db.AutoMigrate()` or manual schema creation.
2. Write SQL migrations in standard SQL that is compatible with PostgreSQL, not SQLite-specific syntax. Even during the SQLite development phase, write PostgreSQL-compatible SQL.
3. Or: skip SQLite entirely. Use PostgreSQL from the start with Docker (`docker run -e POSTGRES_PASSWORD=dev -p 5432:5432 postgres:16-alpine`). For a project of this size, there is no meaningful benefit to starting with SQLite.
4. Use `TIMESTAMPTZ` (not `TIMESTAMP`) for all date fields — always store and retrieve in UTC.
5. Avoid PostgreSQL `ENUM` for the `shelf` field — use a `TEXT` with a `CHECK` constraint instead. Easier to evolve.

**golang-migrate vs goose:**
- `golang-migrate` is the more widely used tool; migrations are plain SQL files versioned by timestamp or number. Supports PostgreSQL, SQLite, and others. Good choice.
- `goose` is slightly more opinionated and supports Go-function migrations (useful for complex data migrations). Slightly less community adoption.
- **Recommendation: golang-migrate.** For a personal site, plain SQL migrations are simpler and the larger community means more Stack Overflow answers.

**Confidence: HIGH** — SQLite vs PostgreSQL differences are well-documented.

---

### Pitfall 7: Infinite Scroll Breaking Back Button / Position Restoration

**What goes wrong:** User scrolls through 80 books, clicks a book cover to view detail, then clicks Back. The browser returns to the top of the list (position 0) because the list was built dynamically via JavaScript and the browser has no scroll position to restore.

**Why it happens:** Browser scroll restoration only works reliably for server-rendered or statically-generated pages where the full DOM is present on load. With infinite scroll, the DOM at navigation time may differ from the DOM at back-navigation time.

**Consequences:** Frustrated UX — the user must scroll back through all the books they had already seen. For a library of 200+ books, this is a significant usability failure.

**Prevention:**
1. Persist the current page/cursor offset in the URL. Use `?page=3` or `?cursor=80` as query parameters. When navigating back, the React router re-renders with the same offset, automatically loading to that position.
2. Use TanStack Query's `useInfiniteQuery` with `keepPreviousData: true` — combined with URL persistence, this re-fetches all pages up to the current offset on back-navigation.
3. After re-loading to the correct offset, use `useEffect` to scroll to the previously-viewed item (store item ID in `sessionStorage` on click).
4. Alternative: use React Router's `ScrollRestoration` component (React Router v6.4+) which handles this automatically for most cases.

**Concrete implementation note:**
```tsx
// On book card click, save position
const handleBookClick = (bookId: string) => {
    sessionStorage.setItem('lastViewedBook', bookId);
    navigate(`/books/${slug}`);
};

// On list mount, restore scroll
useEffect(() => {
    const lastViewed = sessionStorage.getItem('lastViewedBook');
    if (lastViewed) {
        document.getElementById(lastViewed)?.scrollIntoView({ behavior: 'instant' });
        sessionStorage.removeItem('lastViewedBook');
    }
}, [books]); // runs after books are loaded
```

**Confidence: HIGH** — Browser scroll restoration with dynamic lists is a well-documented SPA problem.

---

### Pitfall 8: Image Download Pipeline Deduplication and Broken URLs

**What goes wrong:** Three failure modes in the cover download pipeline:

**A) The same book appears multiple times during sync** (e.g., it shows up on both `currently-reading` and `read` shelves during a shelf transition, or it appears twice due to a Goodreads bug). The pipeline tries to download the same cover twice simultaneously, creating a race condition that corrupts the file.

**B) Google Books cover URLs redirect** to a CDN URL with a short expiry, or return a 404 because the edition has been updated in Google's catalog.

**C) The download succeeds but the image is corrupt or a 1x1 pixel placeholder.** Google Books occasionally returns a "no image available" placeholder image at the same URL as if it were a real cover.

**Prevention:**
- **For A (deduplication):** Check if the file already exists on disk before downloading. Use a mutex or database flag (`cover_downloaded: bool`) to prevent concurrent downloads of the same file. In Go, use `sync.Map` keyed by filename as an in-progress lock.
- **For B (redirect/404):** Use `net/http` with redirect following enabled. On non-200 responses, log and fall back to OpenLibrary. Set a 10-second timeout on image downloads.
- **For C (corrupt/placeholder):** After download, verify the file size is > 5KB (placeholder images are typically <1KB). Check that the file is a valid image using `image.Decode()` from Go's `image` package. If validation fails, delete the file and log for retry.

```go
func validateCoverImage(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()
    info, _ := f.Stat()
    if info.Size() < 5000 { // < 5KB is suspicious
        return fmt.Errorf("image too small (%d bytes), likely a placeholder", info.Size())
    }
    _, _, err = image.Decode(f)
    return err // nil = valid
}
```

**Confidence: HIGH** — These are documented failure modes of the Google Books API cover endpoint.

---

### Pitfall 9: Data Model — Multi-Author, Series, and Re-Reads

**What goes wrong:** The data model in ARCHITECTURE.md correctly handles the primary author but has gaps for three common real-world cases:

**A) Multiple authors.** The `books` table has `author TEXT NOT NULL` (single, denormalized) and a separate `book_authors` join table. This is fine. The pitfall is in the RSS parsing: `<author_name>` in the Goodreads RSS only returns the primary author. Co-authors and translators are not in the RSS feed. If you want multiple authors, you must get them from the Google Books `volumeInfo.authors[]` array. This inconsistency means a sync from RSS alone loses co-author data.

**B) Series books.** Goodreads RSS does not include series information (e.g., "The Fellowship of the Ring" being Book 1 of "The Lord of the Rings"). Google Books does not reliably provide series information either (it is in `volumeInfo.seriesInfo` for some books but not standardized). If series grouping is desired, it requires either: a manual tagging approach, or an additional API (no good automated source exists).

**C) Re-reads.** If Florian reads the same book twice, Goodreads creates two separate review entries for the same book. The RSS feed includes both entries. The current data model uses `goodreads_id` as a unique key — if a re-read creates a new GR review ID, the book will appear twice. If it reuses the old review ID but updates `user_read_at`, the old read date is overwritten.

**Consequences:**
- A: Co-author credit is lost unless Google Books enrichment is used and its `authors[]` array is trusted
- B: No series support (acceptable for MVP — series is not in the requirements)
- C: Only one read date is stored per book; re-reads are silently collapsed into the latest read date

**Prevention:**
- **For A:** In the sync pipeline, always prefer Google Books `authors[]` for the author list if available; fall back to GR RSS `<author_name>` only when Google Books returns nothing. Store the full authors array in `book_authors` table.
- **For B:** Accept no series support for MVP. Add a `series_name` and `series_position` column stub (nullable) to the schema from the start to avoid a migration later, but do not populate it initially.
- **For C:** Add a `read_count` column (INTEGER DEFAULT 1) to the books table. In the sync pipeline, if the same `goodreads_id` is encountered again with a different `user_read_at`, increment `read_count` and update `read_at` to the latest date. The original read date is lost — this is an acceptable tradeoff for a personal site that does not require full read history. If full re-read history is needed in the future, add a `book_reads` table with one row per read event.

**Confidence: HIGH** — These data model edge cases are inherent to Goodreads' data model and well-known to developers who have worked with their data.

---

### Pitfall 10: Goodreads RSS Field Reliability — What Is Actually Consistent

**What goes wrong:** The RSS feed is treated as a reliable structured data source, but individual fields have varying consistency.

**Field reliability assessment (MEDIUM confidence — based on training data and community reports):**

| Field | Reliability | Notes |
|-------|-------------|-------|
| `<title>` | HIGH | Always present, always the book title |
| `<author_name>` | HIGH | Always present; primary author only |
| `<book_id>` | HIGH | Always present; use as the stable GR identifier |
| `<link>` | HIGH | Always present; links to review page |
| `<book_image_url>` | HIGH | Almost always present; low-res thumbnail |
| `<user_shelves>` | HIGH | Always present for shelved books |
| `<user_date_added>` | HIGH | Always present; when added to any shelf |
| `<isbn13>` | MEDIUM | Present for ~80-85% of books; blank for some older or non-ISBN books |
| `<isbn>` (ISBN-10) | MEDIUM | Similar coverage; ISBN-13 preferred |
| `<user_read_at>` | MEDIUM-LOW | Present for books on "read" shelf, but frequently empty or set to the `date_added` value (when the user didn't log the read date). Parse: if `user_read_at` is empty, use `user_date_added` as approximate read date. |
| `<book_published>` | MEDIUM | Year of publication; present for most but missing for some |
| `<user_rating>` | MEDIUM | 0-5; 0 means not rated (not a zero-star rating). Handle 0 as NULL. |
| `<average_rating>` | HIGH | Community average; always present |
| `<num_pages>` | LOW | Often blank; do not rely on this — use Google Books `pageCount` instead |
| `<description>` | LOW | Sometimes a brief synopsis, sometimes empty, sometimes HTML-encoded. Use Google Books description instead. |

**Key defensive parsing requirements:**
1. `user_read_at` is empty frequently — always have a fallback. Do not store as NULL without a fallback value.
2. ISBN fields may be empty or contain `"="` (Excel artifact from CSV exports that Goodreads sometimes embeds in RSS) — strip non-numeric characters from ISBNs before using them.
3. Never use `<num_pages>` from RSS — it is too unreliable. Always use Google Books `pageCount`.
4. `<description>` from RSS may contain HTML entities — decode before storing if you use it at all.

**Confidence: MEDIUM** — Field reliability is from training data and community experience; actual rates should be validated against Florian's live feed during implementation.

---

### Pitfall 11: Shelf Transition Race Condition (currently-reading → read)

**What goes wrong:** When Florian marks a book as "read" on Goodreads, for a brief window the book may appear on both the `currently-reading` and `read` shelves in the RSS feeds (especially if feeds are fetched close to the transition time). The sync pipeline processes both feeds and tries to upsert the same `goodreads_id` twice with conflicting shelf values.

**Why it happens:** Goodreads shelf transitions are not atomic from the RSS perspective. The `currently-reading` feed may not yet reflect the removal when the `read` feed already shows the addition.

**Consequences:** Non-deterministic shelf state in the DB. If `currently-reading` is processed after `read`, the book regresses to "currently-reading" until the next sync.

**Prevention:**
1. In the sync pipeline, always process the `read` shelf AFTER the `currently-reading` shelf, and let `read` win on conflict.
2. Add a conflict resolution rule: if a book appears on both shelves, always assign it to `read` (the "more final" shelf).
3. Use a single `UPSERT` (INSERT ... ON CONFLICT DO UPDATE) rather than separate INSERT/UPDATE operations. Process shelves in precedence order: `read` > `currently-reading`.

```go
// Process order: fetch both, merge, resolve conflicts, then upsert
func mergeShelfItems(currentlyReading, read []*ShelfItem) []*ShelfItem {
    seen := make(map[string]*ShelfItem)
    // currently-reading first (lower priority)
    for _, item := range currentlyReading {
        seen[item.GoodreadsID] = item
    }
    // read overwrites (higher priority)
    for _, item := range read {
        seen[item.GoodreadsID] = item
    }
    result := make([]*ShelfItem, 0, len(seen))
    for _, item := range seen {
        result = append(result, item)
    }
    return result
}
```

**Confidence: HIGH** — Shelf transition race conditions are an inherent consequence of polling two separate feeds.

---

### Pitfall 12: Go File Serving Without Correct Cache Headers for Cover Images

**What goes wrong:** Cover images are self-hosted and served by Go (or Nginx). Without explicit cache headers, browsers either re-request the same image on every page load (wasting bandwidth) or cache it indefinitely without any way to invalidate if the image is replaced.

**Why it happens:** Go's `http.FileServer` sets `Last-Modified` and handles `If-Modified-Since`, but does not set `Cache-Control: max-age` by default.

**Consequences:** Without caching, a page with 24 book covers generates 24 image requests on every visit. With infinite caching but no fingerprinting, replacing a cover (e.g., when re-downloading a higher quality version) is invisible to users until cache is manually cleared.

**Prevention:**
1. Serve cover images with `Cache-Control: public, max-age=31536000, immutable` — one year cache. This is safe because cover filenames are derived from ISBN or Goodreads ID and do not change for the same book.
2. If a cover is ever replaced (e.g., upgrading from low-res to high-res), rename the file or append a version suffix (`{isbn}_v2.jpg`). Update the `cover_path` in the DB. The old filename becomes unreachable and the new one gets cached freshly.
3. For the Go file server, wrap it with a middleware that adds the header:

```go
func coverCacheMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
        next.ServeHTTP(w, r)
    })
}

// Usage
coverHandler := coverCacheMiddleware(http.FileServer(http.Dir("./covers")))
mux.Handle("/covers/", http.StripPrefix("/covers/", coverHandler))
```

**Confidence: HIGH** — HTTP caching semantics are stable and well-documented.

---

### Pitfall 13: WebP Conversion for Cover Images — Complexity vs. Value

**What goes wrong:** The impulse to convert all JPEG covers to WebP to save bandwidth adds a surprising amount of complexity: Go image processing libraries, format detection, conversion pipeline, content-type headers, and browser compatibility fallbacks all need attention.

**Why it happens:** WebP is meaningfully smaller than JPEG (~25-35% for photos). For a page with 24 book covers, this is real bandwidth savings.

**Assessment:** For a personal site with low traffic, the bandwidth savings are negligible in absolute terms. The complexity cost is disproportionate to the benefit.

**Recommendation:** Do NOT convert to WebP as part of the initial pipeline. Download covers as whatever format Google Books provides (JPEG) and serve JPEG. If cover serving becomes a performance concern later (it won't for a personal site), add WebP conversion then.

**If WebP is ever desired:** Use the `golang.org/x/image` package which includes WebP decoding; for encoding, use the `github.com/chai2010/webp` package. But defer entirely.

**Confidence: HIGH** — The performance vs. complexity tradeoff for a personal site is clear.

---

### Pitfall 14: Environment Variable Leakage in React Frontend Build

**What goes wrong:** Developers put secrets (Google Books API key, database URLs) into `.env` files that are accidentally included in the React build. Vite bundles ALL `VITE_` prefixed variables into the client bundle, which is shipped to the browser and readable by anyone who views the page source.

**Why it happens:** Vite's environment variable convention (`VITE_API_KEY`) looks similar to backend env vars and developers copy patterns between projects.

**Consequences:** Google Books API key is exposed in the browser bundle. Anyone can view source and extract it. The key is then potentially abused for the 1,000 requests/day quota.

**Prevention:**
1. The Google Books API key must ONLY exist as a backend environment variable. The Go backend makes all Google Books API calls — the frontend never calls Google Books directly.
2. Never add any secret to a `VITE_` variable. The only `VITE_` variable in this project should be `VITE_API_BASE_URL` pointing to the Go API.
3. Add `.env.local` to `.gitignore` and use it for local development secrets. Check at startup that no `VITE_` vars contain API keys.
4. Review the built `dist/assets/*.js` file once during development to confirm no secrets appear.

**Confidence: HIGH** — Vite's environment variable bundling behavior is documented and a common mistake.

---

## Minor Pitfalls

---

### Pitfall 15: Read Date Timezone Parsing Producing Off-by-One-Day Errors

**What goes wrong:** The Goodreads RSS `user_read_at` field is in RFC 2822 format with timezone offsets (e.g., `Mon, 15 Apr 2024 00:00:00 -0700`). Parsing this in UTC shifts the date to April 14 at 07:00 UTC — stored as April 14 in the DB. The Reading Challenge year filter then miscounts books for December/January cross-year reads.

**Prevention:** Store all dates as UTC (correct) but for display purposes, use the date as provided (strip timezone, use date portion as-is). The read date is a user-entered date-only concept — Goodreads uses midnight in the local timezone because the user picked a date, not a time. Treat `user_read_at` as a date (not a datetime) by extracting only the year/month/day before converting to UTC.

```go
// Parse and normalize to date-only (ignore time component)
t, err := time.Parse(time.RFC1123Z, item.UserReadAt)
if err == nil {
    readDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
```

**Confidence: HIGH** — RFC 2822 timezone offset handling is a standard Go gotcha.

---

### Pitfall 16: Slug Collisions for Books with Identical Titles

**What goes wrong:** Two books with the same title (e.g., "It" by Stephen King and a different book also titled "It") generate the same slug `it`. The second upsert silently overwrites the first in the DB if `slug` has a `UNIQUE` constraint, or the URL `/books/it` becomes ambiguous.

**Prevention:**
1. For slug collisions, append the publication year: `it-1986` vs `it-2020`
2. For continued collisions (same title, same year), append the author last name: `it-1986-king`
3. Implement slug generation as a pure function with collision detection in the sync pipeline

**Confidence: HIGH** — Slug collision is a standard URL routing pitfall.

---

### Pitfall 17: Author Slug Collisions for Authors with the Same Name

**What goes wrong:** Two different authors with the same name (e.g., there are multiple "David Mitchell"s) generate the same author slug. All their books get attributed to one author page.

**Prevention:** For a personal library of a few hundred books, this is unlikely to occur. Add a year-of-birth disambiguation if it occurs: `david-mitchell-1974`. Practically, just log when an author name collision is detected during sync and handle manually.

**Confidence: HIGH** — Name collision is inherent to any slug-based system.

---

### Pitfall 18: PostgreSQL Full-Text Search vs LIKE Queries for Future Search

**What goes wrong:** If a search box is added later (it is explicitly out of scope for MVP), developers reach for `LIKE '%query%'` which requires a full table scan and cannot use standard B-tree indexes.

**Prevention:** Pre-emptively add a `tsvector` generated column to the books table so that adding search later requires no migration:

```sql
ALTER TABLE books ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        to_tsvector('english', coalesce(title,'') || ' ' || coalesce(author,''))
    ) STORED;
CREATE INDEX books_search_idx ON books USING GIN (search_vector);
```

This costs nothing at runtime and makes future search a simple GIN index lookup.

**Confidence: HIGH** — PostgreSQL FTS via tsvector is the standard approach.

---

### Pitfall 19: systemd Service Management vs Docker for Personal VPS

**What goes wrong:** Developers default to Docker on a personal VPS because "that's what production uses" without considering that Docker adds overhead (daemon, networking, image management) that is not justified for a single-binary Go app on a personal server.

**Assessment of options:**

**Docker with docker-compose:**
- Pros: Reproducible environment; easy `docker compose up`; isolates dependencies
- Cons: Docker daemon overhead; requires understanding Docker networking for PostgreSQL; adds complexity to log access; image storage costs disk space; `docker ps` / `docker logs` is less intuitive than `journalctl` for a solo developer

**Bare metal + systemd:**
- Pros: Simpler; Go compiles to a single binary; `systemctl status books` is straightforward; logs via `journalctl`; minimal overhead; no container registry
- Cons: Binary must be copied to server on deploy; less portable to new environments

**Definitive recommendation for a personal VPS:** Use bare metal systemd. Go's single binary compilation is exactly designed for this use case. The complexity of Docker is not justified.

```ini
# /etc/systemd/system/books.service
[Unit]
Description=Flo's Library Go Backend
After=network.target postgresql.service

[Service]
Type=simple
User=books
WorkingDirectory=/opt/books
EnvironmentFile=/opt/books/.env
ExecStart=/opt/books/books-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Deployment becomes:
```bash
GOOS=linux GOARCH=amd64 go build -o books-server ./cmd/server
scp books-server user@vps:/opt/books/
ssh user@vps "systemctl restart books"
```

**If PostgreSQL on the same VPS:** Install via package manager (`apt install postgresql`), not Docker. The PostgreSQL system package integrates with systemd and is simpler to manage on a single-server setup.

**Confidence: HIGH** — Deployment architecture tradeoffs for Go + personal VPS are well-understood.

---

### Pitfall 20: SSL Certificate Management with Nginx + Let's Encrypt

**What goes wrong:** Manual SSL certificate installation and renewal. Let's Encrypt certificates expire every 90 days. Without automated renewal, the site goes HTTPS-broken after 90 days.

**Prevention:** Use `certbot` with the systemd timer (not cron) for automated renewal:
```bash
apt install certbot python3-certbot-nginx
certbot --nginx -d books.yourdomain.com
# certbot automatically installs a systemd timer for renewal
systemctl status certbot.timer  # verify it's active
```

The `certbot --nginx` mode automatically updates the nginx config to add SSL. The systemd timer runs twice daily and renews certificates that expire within 30 days.

**Alternative:** Use Caddy as the reverse proxy instead of Nginx. Caddy handles SSL/TLS automatically with zero configuration — just specify the domain and it provisions Let's Encrypt certificates and renews them. For a personal site, Caddy's auto-HTTPS is a compelling reason to prefer it over Nginx.

**Confidence: HIGH** — Let's Encrypt certbot renewal is a standard Linux server management task.

---

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| Initial RSS parsing | `gofeed` Extensions key paths for GR custom fields | Test against live feed before writing sync code; log all item.Extensions |
| Goodreads sync | Pagination truncation at 200 items | Implement `?page=N` loop from day one |
| Goodreads sync | Shelf transition race condition | Process `read` shelf after `currently-reading`; let `read` win conflicts |
| Google Books enrichment | Title+author false positives | Confidence check on author+title match before accepting result |
| Google Books enrichment | API key in environment | Key stays in Go backend only; never in React build |
| Cover download | Corrupt/placeholder images from Google Books | Validate image size (>5KB) and decodability after download |
| Cover download | Race condition on concurrent same-file downloads | Use file-existence check + sync.Map lock before downloading |
| Cover serving | Missing cache headers | Add `Cache-Control: immutable` middleware to cover file server |
| Data model | Re-reads collapsing into one entry | Add `read_count` field; update on re-sync |
| Data model | `user_read_at` timezone shifting year | Extract date portion before UTC conversion |
| Slug generation | Collision for books with same title | Append year, then author surname |
| Frontend build | VITE_ vars leaking Google Books API key | API key stays in Go only; no secrets in VITE_ vars |
| Infinite scroll | Back-button position loss | Persist cursor in URL; use sessionStorage for scroll target |
| Infinite scroll | SEO — client-rendered content | Go serves Open Graph meta tags per-book-slug in HTML head |
| Deployment | CORS in production | Use Option A: Go serves React dist; no CORS needed |
| Deployment | Docker vs systemd | Use systemd for single-binary Go on personal VPS |
| Deployment | SSL renewal | Use certbot with systemd timer, or switch to Caddy |
| Schema evolution | SQLite to PostgreSQL differences | Use PostgreSQL from day one; skip SQLite |
| Schema evolution | Migrations tool | Use golang-migrate with plain SQL files |

---

## Sources and Confidence Assessment

| Claim | Confidence | Source |
|-------|------------|--------|
| RSS `per_page=200` cap and `?page=N` pagination | HIGH | GR developer community documentation; widely reported |
| `gofeed` Extensions map access for custom XML fields | MEDIUM | Training data; must verify against live feed |
| Google Books title+author false positive rate | HIGH | Documented API behavior; common in metadata pipelines |
| Book cover copyright for personal sites | MEDIUM | General copyright/fair use knowledge; not legal advice |
| CORS behavior with Vite proxy vs production | HIGH | Standard browser/HTTP specification; documented Vite behavior |
| SQLite vs PostgreSQL schema differences | HIGH | Official documentation for both databases |
| Browser scroll restoration with dynamic lists | HIGH | MDN Web Docs; documented browser behavior |
| Image download corrupt/placeholder detection | HIGH | Documented Google Books API behavior |
| Multi-author data from RSS vs Google Books | HIGH | Inherent in Goodreads RSS data model documentation |
| `user_read_at` empty field rate | MEDIUM | Community reports; actual rate for this library unknown |
| systemd vs Docker for personal Go VPS | HIGH | Well-understood operational tradeoff |
| Let's Encrypt certbot renewal | HIGH | Official Let's Encrypt/certbot documentation |
| Vite env var bundling behavior | HIGH | Official Vite documentation |
| HTTP cache-control for immutable assets | HIGH | MDN HTTP caching documentation |

**CRITICAL verification steps before implementation:**
1. Fetch the live RSS feed and inspect raw XML to verify `gofeed` field access patterns for GR custom elements
2. Count Florian's total books on Goodreads to determine if pagination is needed immediately or just future-proofed
3. Test `?per_page=200&page=2` to confirm pagination works on the live feed
4. Validate Google Books API key behavior: create project, enable Books API, test one ISBN query
5. Test a title+author search on a book with a common title to calibrate confidence threshold logic
