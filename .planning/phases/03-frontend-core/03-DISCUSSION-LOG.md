# Phase 3: Frontend Core - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-03
**Phase:** 03-frontend-core
**Areas discussed:** Bio section, Now Reading, Book card density, Dark mode toggle, Loading & error states, Mobile nav

---

## Bio Section

| Option | Description | Selected |
|--------|-------------|----------|
| Hardcoded in JSX | Simplest — personal site, rarely changes | |
| JSON/TS config file | Easy to update without touching JSX | |
| Markdown file bundled | Allows richer formatting | ✓ |

**User's choice:** Markdown file bundled (e.g. `src/content/bio.md`)

---

| Option | Description | Selected |
|--------|-------------|----------|
| Photo left, text right | Classic author bio layout | ✓ |
| Photo centered above, text below | Portrait / vertical stack | |
| Large hero with photo background | Text overlay | |

**User's choice:** Photo left, text right (classic author bio layout)

---

| Option | Description | Selected |
|--------|-------------|----------|
| Yes — link to Goodreads profile | `https://www.goodreads.com/user/show/79499864-florian` | ✓ |
| No — keep bio self-contained | No outbound links | |

**User's choice:** Yes — include Goodreads profile link

---

## Now Reading

| Option | Description | Selected |
|--------|-------------|----------|
| Hide section entirely | No empty state shown | ✓ |
| Friendly message | "Nothing right now — check back soon!" | |
| Subtle placeholder | Empty covers with brand gradient | |

**User's choice:** Hide the section entirely when shelf is empty

---

| Option | Description | Selected |
|--------|-------------|----------|
| Show all in a horizontal row | No limit | |
| Maximum 3–4, truncate the rest | Clean limit | ✓ |
| Horizontal scroll if > N | Shelf-like | |

**User's choice:** Max 3–4 books, truncate the rest (no horizontal scroll)

**Notes (from free-text clarification):** Books wrap vertically in a responsive grid (not horizontal scroll). Fewer columns than "Books Read" since covers are larger. Agent decides exact cap and column counts based on screen size.

---

| Option | Description | Selected |
|--------|-------------|----------|
| Same as Books Read cards | Cover + title + author | |
| Larger covers with title/author | More prominent "currently reading" feel | ✓ |
| Cover only | Minimal | |

**User's choice:** Larger covers with title/author (more prominent)

---

## Book Card Density

| Option | Description | Selected |
|--------|-------------|----------|
| Cover + title + author only | Clean, minimal | ✓ |
| + read date | "Read January 2024" | |
| + year published | Publication year | |

**User's choice:** Cover + title + author only

---

| Option | Description | Selected |
|--------|-------------|----------|
| No count | Just heading "Books Read" | ✓ |
| Yes — show total count | "Books Read · 342" | |

**User's choice:** No count — just "Books Read"

---

## Dark Mode Toggle

| Option | Description | Selected |
|--------|-------------|----------|
| In the sidebar, near nav links | Integrated with navigation | ✓ |
| Top-right of main content area | Floats in content zone | |
| Top-right of nav bar | Global header placement | |

**User's choice:** In the sidebar, near nav links

---

| Option | Description | Selected |
|--------|-------------|----------|
| Sun/moon icon button | Widely understood | ✓ |
| Toggle switch | Explicit on/off feel | |
| Simple text link | "Dark" / "Light" | |

**User's choice:** Sun/moon icon button

---

## Loading & Error States

| Option | Description | Selected |
|--------|-------------|----------|
| Skeleton cards | Placeholder rectangles matching card shape | ✓ |
| Spinner | Simple loading indicator | |
| Nothing | Content pops in when ready | |

**User's choice:** Skeleton cards

---

| Option | Description | Selected |
|--------|-------------|----------|
| Inline error message | In the section itself | |
| Toast/notification | At top of page | ✓ |
| Agent's discretion | Keep it simple | |

**User's choice:** Toast/notification at top of page

---

## Mobile Navigation

| Option | Description | Selected |
|--------|-------------|----------|
| Tabs with icons | Always visible | |
| Hamburger menu | Reveals full sidebar nav | |
| Agent's discretion | Pick what works best | ✓ |

**User's choice:** Agent's discretion

---

## Agent's Discretion

- Exact "Now Reading" display cap (3 or 4)
- Book grid column counts per breakpoint
- Mobile nav pattern
- Toast/notification implementation
- Bio section mobile stack order
- Markdown parsing approach for bio content

## Deferred Ideas

None — discussion stayed within phase scope.
