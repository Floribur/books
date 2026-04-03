# Requirements: Flo's Library

**Defined:** 2026-04-01
**Core Value:** A beautiful, always-up-to-date view of every book Florian has read — synced from Goodreads without manual effort.

---

## v1 Requirements

### Sync Pipeline

- [x] **SYNC-01**: App automatically polls Goodreads RSS feed every 6 hours and updates the book database
- [x] **SYNC-02**: Sync handles RSS pagination (>200 books) using `?page=N` loop
- [x] **SYNC-03**: When a book appears on both `currently-reading` and `read` shelves, `read` wins
- [x] **SYNC-04**: Books are enriched via Google Books API (ISBN lookup primary, title+author fallback with confidence check)
- [x] **SYNC-05**: OpenLibrary API used as fallback when Google Books lacks metadata
- [x] **SYNC-06**: Book cover images are downloaded and stored locally (not hotlinked)
- [x] **SYNC-07**: Cover images are validated after download (rejects 1×1 pixel placeholders, files < 5KB)
- [x] **SYNC-08**: Manual sync can be triggered via `POST /admin/sync`
- [x] **SYNC-09**: Goodreads CSV export can be imported via `POST /admin/import-csv` as permanent fallback

### Data Model

- [x] **DATA-01**: Books have: title, slug, author(s), cover image path, description, genres, page count, publication year, Goodreads ID, read date, ISBN-13, metadata source
- [x] **DATA-02**: Slugs are unique; collision resolved by appending year then author surname
- [x] **DATA-03**: Books track read count (for re-reads)
- [x] **DATA-04**: Authors and genres are normalized (many-to-many join tables)
- [x] **DATA-05**: Schema managed by golang-migrate with numbered SQL files

### API

- [ ] **API-01**: `GET /api/books` returns cursor-paginated list (keyed on `read_at` + `id`), descending order
- [ ] **API-02**: `GET /api/books/currently-reading` returns all currently-reading books
- [ ] **API-03**: `GET /api/books/:slug` returns full book detail
- [ ] **API-04**: `GET /api/authors` returns all authors with book counts
- [ ] **API-05**: `GET /api/authors/:slug` returns author detail with their books (cursor-paginated)
- [ ] **API-06**: `GET /api/genres` returns all genres with book counts
- [ ] **API-07**: `GET /api/genres/:slug` returns genre detail with its books (cursor-paginated)
- [ ] **API-08**: `GET /api/years` returns distinct years with book counts (for Reading Challenge)
- [ ] **API-09**: `/covers/:filename` serves cover images with `Cache-Control: immutable` headers
- [ ] **API-10**: Per-book Open Graph meta tags injected into the SPA HTML head by Go

### Home Page

- [ ] **HOME-01**: Bio section displays Florian's photo and personal description
- [ ] **HOME-02**: "Now Reading" section displays currently-reading books (cover, title, author)
- [ ] **HOME-03**: "Books Read" section displays all read books in descending order
- [ ] **HOME-04**: Book cards are clickable (cover or title) and navigate to detail page

### Book Detail Page

- [x] **BOOK-01**: Displays book cover, title, author(s), publication year, genres, page count, description
- [x] **BOOK-02**: Author name(s) are links to author pages
- [x] **BOOK-03**: Genre tags are links to genre pages
- [x] **BOOK-04**: "Read on [date]" shown when available

### Author Pages

- [x] **AUTH-01**: Author index lists all authors whose books Florian has read, with book count
- [x] **AUTH-02**: Author detail page lists all books by that author in descending read-date order

### Genre Pages

- [x] **GENR-01**: Genre index lists all genres with book count
- [x] **GENR-02**: Genre detail page lists all books in that genre in descending read-date order

### Reading Challenge

- [x] **CHAL-01**: Reading Challenge page shows all books read in the selected year
- [x] **CHAL-02**: Year selector allows browsing past years
- [x] **CHAL-03**: Stats strip shows book count (and optionally page total) for the selected year
- [x] **CHAL-04**: Current year is selected by default

### UI / UX

- [ ] **UI-01**: Infinite scroll on "Books Read" — first ~24 books load immediately, more load on scroll
- [ ] **UI-02**: Explicit "Load More" button as accessible fallback alongside Intersection Observer
- [ ] **UI-03**: Back-button restores scroll position after navigating to a book and returning
- [ ] **UI-04**: Sidebar navigation with animated reading figure (Lottie via lottie-react) — D-01 overrides original "CSS keyframes, no Lottie" spec
- [ ] **UI-05**: `prefers-reduced-motion` respected — animation disabled when set
- [ ] **UI-06**: Dark/light mode toggle with OS preference detection and localStorage persistence
- [ ] **UI-07**: Brand color `#6d233e` (wine red) as primary, `#c4843a` (antique gold) as accent, `#faf8f5` warm off-white background
- [ ] **UI-08**: Typography: Playfair Display (headings) + Inter (body)
- [ ] **UI-09**: Book cover placeholder is CSS gradient using brand colors (no broken image icon)
- [ ] **UI-10**: Responsive layout (mobile collapses to 2-column book grid)

### Deployment

- [ ] **DEPL-01**: Production Go binary embeds React build — single process, no Nginx needed
- [ ] **DEPL-02**: systemd service file for VPS deployment
- [ ] **DEPL-03**: SSL via certbot or Caddy
- [ ] **DEPL-04**: Google Books API key stored as Go backend env var only (never in React/browser)

---

## v2 Requirements

### Enhancements

- **V2-01**: JSON-LD structured data / Book schema markup for SEO
- **V2-02**: Full-text search (PostgreSQL tsvector; stub column in schema now)
- **V2-03**: Stats visualizations on Reading Challenge page (books per month chart, genre breakdown)
- **V2-04**: WebP conversion for cover images
- **V2-05**: Progressive Web App / offline support
- **V2-06**: Reading stats / annual summary page (total pages, avg books/month)

---

## Out of Scope

| Feature | Reason |
|---------|--------|
| User accounts / auth for visitors | Personal showcase, not a social app |
| Star ratings display | Goodreads is source of truth; surfacing GR ratings as-is is confusing without context |
| User reviews or comments | Auth + moderation complexity; out of scope for v1 |
| Numbered pagination | Conflicts with infinite scroll requirement |
| Social sharing widgets | Open Graph tags are sufficient |
| Reading progress percentage | Goodreads RSS does not expose granular progress data |
| Mobile app | Web-first; mobile later if ever |
| Goodreads API (official) | Deprecated 2020; use RSS + Google Books instead |
| Hotlinked book covers | Reliability — external URLs break; self-hosting is required |
| Multi-user support | This is Florian's library, not a platform |

---

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| SYNC-01 through SYNC-09 | Phase 1 | Pending |
| DATA-01 through DATA-05 | Phase 1 | Pending |
| API-01 through API-10 | Phase 2 | Pending |
| HOME-01 through HOME-04 | Phase 3 | Pending |
| UI-01 through UI-03, UI-06 through UI-10 | Phase 3 | Pending |
| BOOK-01 through BOOK-04 | Phase 4 | Pending |
| AUTH-01 through AUTH-02 | Phase 4 | Pending |
| GENR-01 through GENR-02 | Phase 4 | Pending |
| CHAL-01 through CHAL-04 | Phase 4 | Pending |
| UI-04, UI-05 | Phase 5 | Pending |
| DEPL-01 through DEPL-04 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 47 total
- Mapped to phases: 47
- Unmapped: 0 ✓

---
*Requirements defined: 2026-04-01*
*Last updated: 2026-04-01 after initial definition*
