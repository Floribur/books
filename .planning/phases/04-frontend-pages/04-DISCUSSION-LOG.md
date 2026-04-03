# Phase 4: Frontend Pages - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-03
**Phase:** 04-frontend-pages
**Areas discussed:** Book detail page, Genre index visualization, Reading Challenge stats, Author index navigation

---

## Book Detail Page

| Option | Description | Selected |
|--------|-------------|----------|
| Stack: cover on top, metadata below | Standard mobile reflow for two-column layouts | |
| Cover stays small + floats left, metadata wraps | Keeps the cover visible alongside metadata, classic inline float layout | ✓ |
| Metadata only on mobile, cover hidden | Maximizes text space, hides cover | |

**User's choice:** Cover stays small + floats left, metadata wraps around it

---

| Option | Description | Selected |
|--------|-------------|----------|
| After 4 lines / ~320 characters | Shorter threshold — more descriptions truncated | |
| After 8 lines / ~640 characters | Longer descriptions only — most descriptions fully visible | ✓ |
| Always fully visible | No truncation at all | |

**User's choice:** After 8 lines / ~640 characters — longer descriptions only

---

| Option | Description | Selected |
|--------|-------------|----------|
| Title → Author links → Genres → Year · Pages · Read date | Most important info at top | ✓ |
| Title → Author links → Read date → Year · Pages → Genres | Date prominent | |
| Title → Author links → Genres → Description → Year · Pages · Read date | Description in the middle | |

**User's choice:** Title → Author links → Genres → Year · Pages · Read date (Recommended)

---

## Genre Index Visualization

| Option | Description | Selected |
|--------|-------------|----------|
| Horizontal bars | Genre name left, bar fills proportionally, count at end | ✓ |
| Tag cloud | Genre names sized by frequency | |
| Numbered list | Rank + name + count, no bars | |

**User's choice:** Horizontal bars

---

| Option | Description | Selected |
|--------|-------------|----------|
| Brand color (#6d233e) full opacity | Bold, consistent with primary palette | |
| Lighter tint of brand color (~30% opacity) | Subtle, understated | ✓ |
| Agent decides | Deferred to implementation | |

**User's choice:** Lighter tint of brand color (~30% opacity)

---

## Reading Challenge Stats

| Option | Description | Selected |
|--------|-------------|----------|
| Show total pages with footnote "*based on available data" | Honest about coverage, still interesting | ✓ |
| Skip total pages | Cleaner, no caveats | |
| Show total pages without caveat | Assumes data is good enough | |

**User's choice:** Show total pages with footnote

---

| Option | Description | Selected |
|--------|-------------|----------|
| Two stats: book count + total pages* | Clean strip | |
| Three stats: book count + total pages* + longest book | More interesting | ✓ |
| Agent decides | Deferred to implementation | |

**User's choice:** Three stats: book count + total pages* + longest book

---

| Option | Description | Selected |
|--------|-------------|----------|
| Most recent year with books | Always shows data, even if current year has no books | ✓ |
| Current calendar year | Matches "Reading Challenge" framing | |
| Agent decides | Deferred to implementation | |

**User's choice:** Most recent year with books (Recommended)

---

## Author Index Navigation

| Option | Description | Selected |
|--------|-------------|----------|
| No A-Z links — simple sorted scrollable list | Simplest, works fine up to ~200 authors | ✓ |
| A-Z letter anchors at the top | Jump links to lettered sections | |
| Agent decides based on author count | Conditional approach | |

**User's choice:** No A-Z links — just a sorted scrollable list

---

## Agent's Discretion

- CSS approach for floating cover on mobile
- Animation/transition for description expand/collapse
- Bar width calculation technique
- Reading Challenge empty state if a year has no books
- BookGrid queryKey shapes for filtered grids

## Deferred Ideas

None — discussion stayed within phase scope.
