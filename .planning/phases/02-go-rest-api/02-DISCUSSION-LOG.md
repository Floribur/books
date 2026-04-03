# Phase 2: Go REST API — Discussion Log

**Session:** 2026-04-03
**Participant:** Florian

---

## Area 1: JSON Response Shape

**Q: How should paginated list endpoints return data and cursor info?**
Options: Envelope object / Flat array + header
**Selected:** Envelope object — `{"items": [...], "next_cursor": "...", "has_more": true}`
*Rationale: cursor lives in body, easy to consume with TanStack Query's `getNextPageParam`*

**Q: How should the cursor value be encoded?**
Options: Opaque base64 token / Transparent composite param (?after_date + ?after_id)
**Selected:** Opaque base64 token — e.g. `"MjAyNC0wMy0xNVQwMDowMDowMFpfNDI="`
*Rationale: hides internals, safe in URLs, easy to evolve without breaking clients*

---

## Area 2: Book Object Embedding

**Q: On GET /api/books (list), what should each book item include?**
Options: Card fields + inline authors/genres / Scalar fields only
**Selected:** Card fields + inline authors/genres — `authors: [{name, slug}]`, `genres: [{name, slug}]`
*Rationale: one JOIN query, no extra frontend round trips for card display*

**Q: On GET /api/books/:slug (detail), anything extra beyond list fields?**
Options: All book fields + inline authors/genres / Same as list + description only
**Selected:** All book fields + inline authors/genres — adds description, page_count, isbn13, read_count, shelf, metadata_source

---

## Area 3: Open Graph Injection Scope

**Clarification requested:** User asked what Open Graph injection is.
**Explanation given:** OG tags control rich link previews when sharing URLs on social media (WhatsApp, iMessage, Discord). React SPAs serve a static index.html — without Go injecting per-book tags dynamically, all book URLs share the same blank preview. OG tags also help search engine indexing for non-JS crawlers.

**Q: How far to go with Open Graph in Phase 2?**
Options: Full OG injection now / Scaffold only, OG in Phase 5
**Selected:** Full OG injection now
*Rationale: scaffolding cost (go:embed + SPA catch-all) is almost identical to full injection; 20 extra lines for the template; no Phase 5 revisit needed*

---

*Log generated: 2026-04-03*
