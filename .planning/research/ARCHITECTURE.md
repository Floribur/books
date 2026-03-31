# Architecture Patterns: Goodreads Sync for Flo's Library

**Domain:** Personal book library — Goodreads sync alternatives
**Researched:** 2026-03-31
**Overall confidence:** MEDIUM
**Confidence note:** WebSearch, WebFetch, and Bash tools were all unavailable during this session. All findings are from training data (cutoff August 2025). The single highest-risk claim — "Goodreads RSS feeds are still live" — must be verified manually before implementation begins. Everything else (Google Books API, OpenLibrary) is HIGH confidence from well-documented, stable APIs.

---

## Executive Summary

The Goodreads public API was shut down in December 2020 and will not return. However, Goodreads continued to expose RSS feeds for public user shelves after the API shutdown, and as of mid-2025 (training cutoff) these feeds were still accessible without authentication. This is the recommended primary sync mechanism for a personal single-user app. It is fragile by nature (Goodreads can remove it without notice) but has survived for years and is the simplest, most reliable approach available.

The recommended architecture is:

1. **Primary sync**: Goodreads RSS feed polling (scheduled, ~every 6 hours)
2. **Metadata enrichment**: Google Books API (covers, description, genres, page count)
3. **Fallback metadata**: OpenLibrary API (when Google Books has no match)
4. **Trigger**: Scheduled cron job in Go, plus a manual `/admin/sync` HTTP endpoint

---

## 1. Goodreads RSS Feeds

### Status (as of training cutoff August 2025)

**HIGH confidence: The RSS feeds exist and are publicly documented.**
**MEDIUM confidence: They were still live as of mid-2025.**
**CRITICAL: Must be manually verified before building the sync pipeline.**

Goodreads exposes RSS feeds for public user shelves at:

```
https://www.goodreads.com/review/list_rss/{USER_ID}?shelf={SHELF_NAME}
```

For Florian's profile (user ID `79499864`):

```
# Books currently reading
https://www.goodreads.com/review/list_rss/79499864?shelf=currently-reading

# Books read (all time)
https://www.goodreads.com/review/list_rss/79499864?shelf=read

# All books (all shelves combined)
https://www.goodreads.com/review/list_rss/79499864
```

**Verification step (do this first):** Open each URL in a browser. If you see XML with `<item>` entries containing `<title>` and `<author_name>`, the feed is live.

### Data Fields Available in RSS

Each `<item>` in the feed contains:

| Field | XML element | Example value | Notes |
|-------|------------|---------------|-------|
| Title | `<title>` | `The Name of the Rose` | Full title |
| Author | `<author_name>` | `Umberto Eco` | Single author |
| Book ID (Goodreads) | `<book_id>` | `119494` | GR internal ID |
| Image (small thumbnail) | `<book_image_url>` | `https://...` | Low-res; use Google Books instead |
| Published year | `<book_published>` | `1980` | Year only |
| Average rating | `<average_rating>` | `4.21` | GR community rating |
| ISBN | `<isbn>` | `9780156001311` | ISBN-13; may be blank |
| ISBN-13 | `<isbn13>` | `9780156001311` | Usually more reliable |
| User's read date | `<user_read_at>` | `Mon, 15 Apr 2024 00:00:00 -0700` | Parse this carefully |
| User's date added | `<user_date_added>` | `Tue, 01 Jan 2019 00:00:00 -0800` | When added to shelf |
| User shelf | `<user_shelves>` | `read` or `currently-reading` | Which shelf |
| User rating | `<user_rating>` | `5` | 1-5, or 0 if not rated |
| Goodreads link | `<link>` | `https://www.goodreads.com/review/show/...` | Review URL |
| Book link | (in link) | `https://www.goodreads.com/book/show/...` | Constructed from book_id |

**What is NOT in the RSS feed:**
- Date started reading (only read date, not start date)
- Reading progress percentage
- User's written review text (just rating)
- Genres / categories (must come from Google Books)
- Page count (must come from Google Books)
- High-quality cover images (thumbnails only; use Google Books)

### RSS Feed Limitations

1. **Pagination**: The RSS feed returns a maximum of ~200 items per request (the exact limit is 100 by default but can be extended with `?per_page=200`). For a large shelf, you may need to handle pagination via `?page=2`, etc. **Florian's shelf size is unknown — verify whether pagination is needed.**

2. **`per_page` parameter**: `?shelf=read&per_page=200` should work to get up to 200 books in one request. The maximum observed value that Goodreads honors is 200.

3. **Rate limiting**: Goodreads does not publish rate limits for RSS. For a personal app polling every 6 hours, there is no risk of being rate-limited. Do not poll more often than every hour.

4. **No authentication**: Public shelves are accessible without login. Florian's profile must be set to "public" in Goodreads settings (which it is, since the profile URL is public).

5. **Feed may break silently**: Amazon (Goodreads' owner) has shown no commitment to maintaining these feeds. If the feed returns 200 OK but with no items, or a redirect to a login page, the sync pipeline must detect and alert on this.

6. **Date parsing**: The `user_read_at` field is in RFC 2822 format with a timezone offset. Go's `time.Parse` with `time.RFC1123Z` handles this. Some entries may have an empty `user_read_at` even for books on the "read" shelf — handle this gracefully (use `user_date_added` as fallback).

### Go RSS Parsing

Use the `github.com/mmcdole/gofeed` library. It is the standard Go RSS/Atom parser, well-maintained as of 2025.

```go
import "github.com/mmcdole/gofeed"

func fetchShelf(userID, shelf string) ([]*gofeed.Item, error) {
    url := fmt.Sprintf(
        "https://www.goodreads.com/review/list_rss/%s?shelf=%s&per_page=200",
        userID, shelf,
    )
    fp := gofeed.NewParser()
    feed, err := fp.ParseURL(url)
    if err != nil {
        return nil, fmt.Errorf("fetching goodreads shelf %s: %w", shelf, err)
    }
    return feed.Items, nil
}
```

For Goodreads-specific custom fields (`<book_id>`, `<author_name>`, etc.) that are outside the standard RSS spec, `gofeed` exposes them via `item.Extensions` map. These are namespaced under the element name. **This will require implementation testing** — the exact key path in `item.Extensions` must be verified against actual feed output.

---

## 2. Web Scraping Goodreads (Fallback)

If the RSS feeds go down, scraping is the fallback. This section documents the approach but it should NOT be the primary method.

### What Is Accessible Without Login

Goodreads public profile pages are HTML-rendered. For a public profile like `https://www.goodreads.com/user/show/79499864-florian`, the shelf pages are at:

```
https://www.goodreads.com/review/list/79499864?shelf=read&sort=date_read&order=d
https://www.goodreads.com/review/list/79499864?shelf=currently-reading
```

These pages render server-side HTML with book rows that include title, author, cover image, and dates.

### Why Scraping Is the Fallback, Not Primary

1. **Fragility**: Goodreads changes HTML structure periodically. Any class name or DOM structure change breaks the scraper.
2. **Bot detection**: Goodreads (Amazon) increasingly uses bot detection. Headless browser approaches (Playwright/Puppeteer) are needed if they add JavaScript rendering, but even then, CAPTCHAs or IP blocks are possible.
3. **Rate limiting**: Scraping HTML is more likely to trigger rate limits than RSS polling.
4. **Maintenance burden**: Every Goodreads UI update requires scraper maintenance.

### If Scraping Becomes Necessary

Use `golang.org/x/net/html` (stdlib-adjacent) or `github.com/PuerkitoBio/goquery` (jQuery-like HTML traversal in Go). The approach would be:

1. Fetch the HTML page with a realistic `User-Agent` header
2. Parse the HTML and find the book rows (currently under `#booksBody tr.bookalike`)
3. Extract fields from `td` elements by their CSS classes
4. Handle pagination (Goodreads shows ~30 books per page with a "next" link)
5. Respect a delay of at least 2-3 seconds between requests

**Confidence: LOW** on the CSS selectors (HTML structure may have changed since training data). These selectors must be verified against the actual live page before any scraping code is written.

---

## 3. Third-Party APIs and Alternatives

### hardcover.app

Hardcover is a Goodreads alternative that launched around 2022-2023 with an open API. As of mid-2025:

- Hardcover provides a GraphQL API for book data
- They have an import feature from Goodreads (CSV export)
- Their API requires an account

**Assessment: Not viable for this use case.** This would require Florian to migrate from Goodreads to Hardcover, which is a user behavior change outside the project scope. The goal is to sync FROM Goodreads, not to switch platforms.

**Confidence: MEDIUM** (Hardcover's API status as of mid-2025; may have evolved)

### StoryGraph

StoryGraph is a Goodreads alternative with no public API. Not usable for programmatic sync.

### OpenLibrary (for metadata, not sync)

OpenLibrary (archive.org) provides free, open book metadata. **This is not a Goodreads sync source** — it has no concept of a user's reading history. It is relevant as a metadata fallback (see Section 5).

### Goodreads CSV Export (Manual Approach)

Goodreads allows users to export their library as a CSV from `https://www.goodreads.com/review/import`. This CSV contains all books with read dates, shelves, ratings, etc.

**This is the nuclear fallback** if both RSS feeds and scraping fail. The workflow would be:
1. Florian exports the CSV from Goodreads manually
2. Uploads it to the app (or drops it in a watched directory)
3. The Go backend parses and imports it

The CSV format includes: `Book Id, Title, Author, Additional Authors, ISBN, ISBN13, My Rating, Average Rating, Publisher, Binding, Number of Pages, Year Published, Original Publication Year, Date Read, Date Added, Bookshelves, Exclusive Shelf, My Review, Spoiler, Private Notes, Read Count, Owned Copies`

**This is a valid last resort** for a personal app. The app should be designed to accept CSV import regardless — it provides a manual override path.

---

## 4. Google Books API

### Current Status (HIGH confidence)

Google Books API v1 is live, maintained, and extensively documented at `https://developers.google.com/books`. It is one of Google's oldest APIs and shows no signs of deprecation.

### Authentication and Rate Limits

| Tier | Key Required | Daily Limit | Notes |
|------|-------------|-------------|-------|
| Unauthenticated | No | ~1,000 requests/day | IP-based; no guarantee |
| API Key (free) | Yes (Google Cloud project) | 1,000 requests/day | Quota tracked per project |
| Paid tier | Yes | Varies by billing | Not needed for personal use |

**For a personal app**, the free API key tier (1,000 req/day) is more than sufficient. Florian's library is likely a few hundred books total, and metadata is fetched once per book and cached. The entire library could be enriched in a single batch during the initial sync.

**Getting an API key:**
1. Create a project at `console.cloud.google.com`
2. Enable the "Books API"
3. Create credentials → API key
4. Optionally restrict the key to the Books API

No OAuth2 is required — API key is sufficient for reading public book data.

### Search Endpoints

**By ISBN (most reliable):**
```
GET https://www.googleapis.com/books/v1/volumes?q=isbn:{ISBN13}&key={API_KEY}
```

**By title and author (fallback when ISBN is missing):**
```
GET https://www.googleapis.com/books/v1/volumes?q=intitle:{TITLE}+inauthor:{AUTHOR}&key={API_KEY}
```

**By Goodreads-provided data:**
Always try ISBN first. Many books in the Goodreads RSS have a valid `isbn13`. When ISBN is blank or returns no results, fall back to title+author search.

### Data Fields Available

| Field | JSON path | Availability | Notes |
|-------|-----------|--------------|-------|
| Title | `volumeInfo.title` | Nearly always | May differ slightly from GR title |
| Authors | `volumeInfo.authors[]` | Nearly always | Array; usually matches GR |
| Description | `volumeInfo.description` | ~80% of books | HTML content; strip tags |
| Published date | `volumeInfo.publishedDate` | ~90% | May be year only, or YYYY-MM-DD |
| Page count | `volumeInfo.pageCount` | ~70% of books | Missing for some older books |
| Categories/genres | `volumeInfo.categories[]` | ~60% of books | Broad (e.g., "Fiction", "Science Fiction") |
| Thumbnail URL | `volumeInfo.imageLinks.thumbnail` | ~75% of books | Small; ~128px wide |
| Small thumbnail | `volumeInfo.imageLinks.smallThumbnail` | ~75% of books | Even smaller |
| ISBN-10 | `volumeInfo.industryIdentifiers[]` | ~80% | type: "ISBN_10" |
| ISBN-13 | `volumeInfo.industryIdentifiers[]` | ~85% | type: "ISBN_13" |
| Language | `volumeInfo.language` | Nearly always | BCP 47 code |
| Publisher | `volumeInfo.publisher` | ~75% | |
| Print type | `volumeInfo.printType` | Nearly always | "BOOK" or "MAGAZINE" |

**Cover image quality note:** The `thumbnail` URL from Google Books API contains a `zoom=1` parameter. Replacing `zoom=1` with `zoom=0` or removing the zoom parameter entirely sometimes returns a larger image. However, this is undocumented behavior — test case-by-case. For the self-hosting pipeline, download whatever is available and accept that quality will vary.

**Categories note:** Google Books categories are coarse-grained (e.g., "Fiction / Science Fiction / General"). They are useful as genre tags but benefit from normalization — strip the hierarchy separator, deduplicate, and lowercase before storing.

### Rate Limit Strategy

For the initial sync of a library (say 300 books):
- Add a 200ms delay between Google Books API requests
- This gives 5 requests/second — well within limits
- 300 books at 200ms delay = 60 seconds total for initial enrichment
- Subsequent syncs only fetch metadata for **new books** (delta sync)

---

## 5. OpenLibrary API (Fallback for Book Metadata)

### What It Is

OpenLibrary (open.library.org, part of the Internet Archive) provides a free, open, no-key-required API for book metadata.

```
# By ISBN
https://openlibrary.org/api/books?bibkeys=ISBN:{ISBN13}&format=json&jscmd=data

# Search by title/author
https://openlibrary.org/search.json?title={TITLE}&author={AUTHOR}
```

### Coverage and Reliability

| Dimension | Assessment | Confidence |
|-----------|------------|------------|
| Coverage | Wide (20M+ editions) but uneven quality | MEDIUM |
| Cover images | Available via covers.openlibrary.org | HIGH |
| Genres/subjects | Available but crowdsourced (inconsistent) | MEDIUM |
| Page count | Available for many editions | MEDIUM |
| API reliability | Generally up but slower than Google | MEDIUM |
| Rate limits | Not published; be polite (1 req/sec) | MEDIUM |

**Cover image URL format:**
```
https://covers.openlibrary.org/b/isbn/{ISBN13}-L.jpg  # Large
https://covers.openlibrary.org/b/isbn/{ISBN13}-M.jpg  # Medium
https://covers.openlibrary.org/b/isbn/{ISBN13}-S.jpg  # Small
```

No API key needed. `-L` suffix returns the largest available cover.

### When to Use OpenLibrary

OpenLibrary is the **fallback when Google Books returns no results or no cover image**. It is not a replacement — Google Books has better coverage for modern books, better structured metadata, and faster response times. OpenLibrary shines for older or more obscure titles.

**Recommended cascade:**
1. Try Google Books by ISBN → found? Use it.
2. Google Books returned no cover image? Try OpenLibrary covers by ISBN.
3. Google Books returned no results? Try OpenLibrary search by title+author.
4. Both failed? Store the book with Goodreads thumbnail (low-res) and log it for manual review.

---

## 6. Recommended Architecture: Complete Sync Pipeline

### System Components

```
┌─────────────────────────────────────────────────────────┐
│                    Go Backend                           │
│                                                         │
│  ┌──────────────┐    ┌──────────────────────────────┐   │
│  │  Cron Job    │───▶│     Sync Worker              │   │
│  │  (every 6h)  │    │                              │   │
│  └──────────────┘    │  1. Fetch GR RSS feeds       │   │
│                      │     - currently-reading      │   │
│  ┌──────────────┐    │     - read                   │   │
│  │ Manual Sync  │───▶│                              │   │
│  │ POST /admin/ │    │  2. Diff against DB          │   │
│  │ sync         │    │     - new books → enrich     │   │
│  └──────────────┘    │     - moved shelves → update │   │
│                      │     - read date changes      │   │
│                      │       → update               │   │
│                      │                              │   │
│                      │  3. Enrich new books         │   │
│                      │     a. Google Books (ISBN)   │   │
│                      │     b. Google Books (title)  │   │
│                      │     c. OpenLibrary (fallback)│   │
│                      │                              │   │
│                      │  4. Download cover images    │   │
│                      │     → store in /data/covers/ │   │
│                      │                              │   │
│                      │  5. Update DB                │   │
│                      └──────────────────────────────┘   │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  REST API (served to React frontend)             │   │
│  │  GET /api/books                                  │   │
│  │  GET /api/books/:slug                            │   │
│  │  GET /api/authors                                │   │
│  │  GET /api/genres                                 │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
        │                          │
        ▼                          ▼
┌──────────────┐          ┌────────────────┐
│  PostgreSQL  │          │  /data/covers/ │
│  (books,     │          │  (self-hosted  │
│   authors,   │          │   images)      │
│   genres,    │          └────────────────┘
│   shelves)   │
└──────────────┘
```

### Sync Trigger Strategy

**Polling (cron), not webhooks.** Goodreads does not offer webhooks. The options are:

| Trigger | Feasibility | Recommended |
|---------|-------------|-------------|
| Cron (scheduled polling) | Yes — native Go ticker | YES — every 6 hours |
| Webhook from Goodreads | No — Goodreads has no webhooks | No |
| Manual trigger via HTTP | Yes — simple admin endpoint | YES — for development/testing |
| Browser extension push | Possible (complex) | No — out of scope |

**Recommended schedule: every 6 hours.** For a personal reading log, a 6-hour sync lag is entirely acceptable. Florian doesn't add more than a few books per week typically.

**Cron implementation in Go:**

```go
// Use github.com/robfig/cron/v3 — the standard Go cron library
import "github.com/robfig/cron/v3"

c := cron.New()
c.AddFunc("0 */6 * * *", func() {
    if err := syncWorker.Run(ctx); err != nil {
        log.Printf("sync failed: %v", err)
    }
})
c.Start()
defer c.Stop()
```

Alternatively, use Go's built-in `time.Ticker` for simplicity without a cron library dependency:

```go
ticker := time.NewTicker(6 * time.Hour)
go func() {
    for range ticker.C {
        syncWorker.Run(ctx)
    }
}()
```

### Delta Sync (Do Not Re-fetch Everything)

On each sync run:
1. Fetch both RSS feeds (currently-reading + read)
2. Build a map: `goodreads_id → {shelf, read_at}`
3. Compare against DB state
4. Only run Google Books enrichment for **new books** (not in DB yet)
5. Update shelf/date for books already in DB that have changed

This keeps Google Books API usage minimal — typically 0-2 calls per sync run for an active reader.

### Data Model (Core Tables)

```sql
-- Books table
CREATE TABLE books (
    id              SERIAL PRIMARY KEY,
    goodreads_id    TEXT NOT NULL UNIQUE,
    title           TEXT NOT NULL,
    slug            TEXT NOT NULL UNIQUE,   -- URL-safe, e.g. "the-name-of-the-rose"
    author          TEXT NOT NULL,          -- Primary author (denormalized for simplicity)
    isbn            TEXT,                   -- ISBN-13 preferred
    cover_path      TEXT,                   -- Local path: /covers/{isbn}.jpg
    description     TEXT,
    page_count      INTEGER,
    published_year  INTEGER,
    shelf           TEXT NOT NULL,          -- 'read' | 'currently-reading'
    read_at         TIMESTAMPTZ,            -- Null if currently-reading
    added_at        TIMESTAMPTZ NOT NULL,   -- When added to GR shelf
    goodreads_url   TEXT,
    google_books_id TEXT,                   -- For re-fetching metadata if needed
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Authors (normalized for author index page)
CREATE TABLE authors (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE
);

CREATE TABLE book_authors (
    book_id   INTEGER REFERENCES books(id),
    author_id INTEGER REFERENCES authors(id),
    PRIMARY KEY (book_id, author_id)
);

-- Genres (from Google Books categories)
CREATE TABLE genres (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE
);

CREATE TABLE book_genres (
    book_id  INTEGER REFERENCES books(id),
    genre_id INTEGER REFERENCES genres(id),
    PRIMARY KEY (book_id, genre_id)
);
```

### Cover Image Pipeline

```
1. Google Books thumbnail URL available?
   YES → Download to /data/covers/{isbn13}.jpg (or {goodreads_id}.jpg if no ISBN)
   NO  → Try OpenLibrary covers.openlibrary.org/b/isbn/{isbn13}-L.jpg
         NO → Store the low-res Goodreads thumbnail URL (the <book_image_url> from RSS)
              Flag as cover_quality: 'low' in DB for later re-fetch

2. Store cover_path in DB as relative path: /covers/{filename}.jpg
   Serve via Go's http.FileServer or Nginx static files

3. On re-sync: if cover_path is null or cover_quality is 'low', retry download
```

**Do not re-download covers that are already on disk.** Check file existence before downloading.

---

## 7. Resilience and Error Handling

### RSS Feed Failure Detection

```go
// If feed returns 0 items but previously had items, something is wrong
if len(items) == 0 && db.BookCount() > 0 {
    log.Printf("WARNING: GR RSS returned 0 items - possible feed issue")
    // Do NOT delete existing books from DB
    // Alert / log only; skip this sync run
    return
}
```

### Google Books "Not Found" Handling

Some books will not be in Google Books (very new releases, obscure regional titles, some non-English books):
- Store the book with Goodreads data only
- Set `google_books_id = NULL`, `cover_path = NULL` (or use GR thumbnail)
- The frontend must handle books without covers (use the CSS placeholder pattern from FEATURES.md)
- Retry enrichment on next sync (with exponential backoff after 3 failures)

### Feed Goes Down Permanently

If Goodreads RSS is removed (the main risk scenario):
1. The app continues serving existing data from DB (no data loss)
2. The manual CSV import endpoint becomes the new sync mechanism
3. Scraping can be added as a replacement sync module without changing the data model

The data model is intentionally source-agnostic — `goodreads_id` is just a string identifier. If the data ever comes from scraping or CSV instead of RSS, the same ID format works.

---

## 8. Recommended Approach Summary

**For Flo's Library, in 2025, the recommended Goodreads sync approach is:**

**RSS polling every 6 hours, enriched via Google Books API, with OpenLibrary as fallback.**

This is recommended over scraping because:
- RSS is structurally simpler to parse and maintain
- RSS has survived since the 2020 API shutdown, suggesting Goodreads tolerates it
- No bot detection risk (RSS is a documented format, not screen-scraping)
- Sufficient data fields for this use case (title, author, shelf, read date, ISBN)

The approach is robust enough for a personal site and the delta sync model keeps operational costs near zero.

**The single most important implementation step:** Manually verify the RSS URLs are live before writing any sync code. This is a 60-second browser test and eliminates the biggest uncertainty in the entire project.

---

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| Initial RSS parsing | GR custom XML fields not in standard RSS namespace | Use `gofeed` + test `item.Extensions` against live feed |
| Google Books by title | High false-positive match rate | Always cross-check ISBN when available; require high confidence match |
| Cover download | Some Google Books cover URLs redirect | Follow redirects; use `net/http` with `CheckRedirect` set to follow |
| Category normalization | Google Books returns "Fiction / Science Fiction / General" | Split on " / " and normalize each segment; deduplicate |
| Read date parsing | `user_read_at` can be empty, malformed, or in wrong timezone | Parse defensively; fall back to `user_date_added`; store as UTC |
| Pagination | Shelf over 200 books | Implement `?page=N` loop with termination condition |
| RSS disappears | No notice from Amazon/Goodreads | Build CSV import as fallback from day one; don't treat RSS as permanent |

---

## Sources and Confidence Assessment

| Claim | Confidence | Source / Rationale |
|-------|------------|-------------------|
| RSS feed URL format | HIGH | Well-documented in Goodreads developer community; URL format is stable |
| RSS feed still live in 2025 | MEDIUM | Reported as working in community discussions through mid-2025; must verify live |
| RSS field names (`<book_id>`, `<author_name>`, etc.) | MEDIUM | Documented in community; XML structure may have changed |
| Google Books API v1 live and free tier 1,000/day | HIGH | Official Google documentation; stable quota |
| Google Books field availability percentages | MEDIUM | Based on broad community experience; actual rates for this library unknown |
| `gofeed` library recommendation | HIGH | Standard Go RSS library, well-maintained as of 2025 |
| `robfig/cron` library | HIGH | Dominant Go cron library; widely used in production |
| OpenLibrary covers API URL format | HIGH | Stable, documented API; no key required |
| Goodreads CSV export availability | HIGH | Feature documented in Goodreads help; unlikely to be removed |
| Goodreads scraping fragility | HIGH | Universal consensus; DOM changes are frequent |
| hardcover.app viability | MEDIUM | Would require platform migration; not a like-for-like replacement |

**CRITICAL verification steps before implementation:**
1. Open `https://www.goodreads.com/review/list_rss/79499864?shelf=read` in a browser — confirm XML is returned
2. Open `https://www.goodreads.com/review/list_rss/79499864?shelf=currently-reading` — confirm feed works
3. Check `?per_page=200` is honored (or determine actual max page size)
4. Obtain a Google Books API key and test a sample ISBN query
