# STACK.md — Flo's Library

**Research type:** Stack decisions
**Confidence:** HIGH (all choices are well-established 2024-2025 conventions)

---

## 1. Go HTTP Router → **Chi v5**

**Use:** `github.com/go-chi/chi/v5`

Chi is the standard choice for idiomatic Go REST APIs in 2025. It's a lightweight router built on `net/http` with composable middleware, no magic, and a large ecosystem. Gin and Echo are both viable but add abstraction over `net/http` that creates friction when reading Go standard library documentation. Chi feels like Go; the others feel like frameworks.

**Do not use:** Gin (non-idiomatic request/response types), Echo (similar issue), Fiber (builds on fasthttp, not net/http — incompatible with standard middleware).

---

## 2. Database Layer → **sqlc**

**Use:** `github.com/sqlc-dev/sqlc` (code generation) + `github.com/jackc/pgx/v5` (PostgreSQL driver)

sqlc generates type-safe Go code from raw SQL queries. You write SQL, sqlc generates the Go functions. This eliminates the ORM abstraction layer while keeping code safe and readable. For a project where the schema is known upfront and queries are not dynamic, sqlc is the right choice.

**Do not use:**
- GORM: Reflection-based ORM with surprising behavior at edge cases, poor performance on complex queries, and generated SQL that's hard to debug.
- sqlx: Requires writing boilerplate scanning code manually. sqlc does this better.

---

## 3. Development Database → **PostgreSQL from day one**

**Use:** Real PostgreSQL instance for development (Docker: `docker run -e POSTGRES_PASSWORD=dev -p 5432:5432 postgres:16`)

Do NOT use SQLite as a dev substitute. The SQL dialects diverge in ways that cause subtle bugs (different date functions, RETURNING clause behavior, JSON support). The PITFALLS.md researcher confirmed: start with PostgreSQL directly. A single `docker compose up` is the correct local dev setup.

**In-memory for tests only:** Use `github.com/jackc/pgx/v5` with a test database, or `testcontainers-go` for integration tests with ephemeral Postgres containers.

---

## 4. Migration Tool → **golang-migrate**

**Use:** `github.com/golang-migrate/migrate/v4`

golang-migrate is the most widely adopted Go migration tool. It uses numbered sequential SQL files (`000001_create_books.up.sql` / `.down.sql`), integrates with the Go binary via embed, and supports both PostgreSQL and SQLite if needed. Simple mental model: up files run forward, down files roll back.

**Do not use:**
- goose: Similar capability but less adoption, slightly different conventions.
- Atlas: Powerful but overkill for a personal project; schema diffing is not needed here.

---

## 5. React Setup → **Vite + React + TypeScript**

**Use:** `npm create vite@latest -- --template react-ts`

Vite is the standard React toolchain in 2025. Fast HMR, minimal config, TypeScript support out of the box.

**Do not use:**
- Create React App: Officially deprecated and unmaintained.
- Next.js: Overkill for this project. The SSR/SEO benefit is real but minimal for a personal book site that Google doesn't need to index. The complexity cost (server components, app router, deployment requirements) is not justified. Plain Vite React + Go serves everything correctly.

**SSR decision:** Skip it. A personal book site does not need Google indexing. If it ever becomes important, migrate to Next.js at that point — the component code is reusable.

---

## 6. Data Fetching → **TanStack Query v5**

**Use:** `@tanstack/react-query` v5

TanStack Query is the standard solution for server state in React. It handles caching, background refetching, loading/error states, and infinite scroll via `useInfiniteQuery`. For a Go REST API consumer, this is the correct choice — it eliminates 80% of boilerplate `useEffect` + `useState` data fetching code.

**Do not use:** Redux Toolkit Query (heavier, more setup for a personal app), SWR (fewer features than TanStack Query v5, no infinite scroll built-in).

---

## 7. Infinite Scroll → **Intersection Observer + useInfiniteQuery**

**Pattern:**
```tsx
const { data, fetchNextPage, hasNextPage } = useInfiniteQuery({
  queryKey: ['books'],
  queryFn: ({ pageParam }) => fetchBooks(pageParam),
  getNextPageParam: (lastPage) => lastPage.nextCursor,
  initialPageParam: null,
})

// Sentinel div at bottom of list
const sentinelRef = useRef(null)
useEffect(() => {
  const observer = new IntersectionObserver(([entry]) => {
    if (entry.isIntersecting && hasNextPage) fetchNextPage()
  })
  observer.observe(sentinelRef.current)
  return () => observer.disconnect()
}, [hasNextPage, fetchNextPage])
```

Use **cursor-based pagination** (not offset), keyed on `last_read_at` timestamp + book ID for stability. This avoids the "item shifts position when new books are added" bug that offset pagination causes.

---

## 8. Book Cover Image Serving → **Chi FileServer + Go embed or fs.FS**

**Pattern:**
```go
// Serve downloaded cover images from disk
r.Handle("/covers/*", http.StripPrefix("/covers", http.FileServer(http.Dir("./data/covers"))))
```

Add cache headers via middleware:
```go
r.Use(func(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if strings.HasPrefix(r.URL.Path, "/covers/") {
      w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
    }
    next.ServeHTTP(w, r)
  })
})
```

Store covers in `./data/covers/{isbn}.jpg`. Filename is the ISBN — stable, deduplicates automatically.

---

## 9. Production Deployment → **Go serves React dist/ (single process)**

**Pattern:**
```go
//go:embed frontend/dist
var staticFiles embed.FS

// Serve React SPA
r.Handle("/*", http.FileServer(http.FS(staticFiles)))
```

Build React (`npm run build` → `frontend/dist/`), embed into Go binary, deploy single binary. No nginx needed for a personal VPS. This eliminates CORS entirely (same origin) and simplifies deployment to `scp binary server:/ && systemctl restart flos-library`.

**Development:** React Vite dev server (port 5173) + Go API (port 8080) with CORS enabled in Go for `localhost:5173` only.

---

## 10. Go Project Layout

```
flos-library/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── api/                  # HTTP handlers
│   ├── db/                   # sqlc generated code + queries
│   ├── sync/                 # Goodreads RSS sync + Google Books enrichment
│   └── models/               # Shared domain types
├── migrations/               # golang-migrate SQL files
├── frontend/                 # Vite React TypeScript app
│   ├── src/
│   └── dist/                 # Built output (embedded into Go binary)
├── data/
│   └── covers/               # Downloaded book cover images
├── docker-compose.yml        # Local Postgres for development
└── Makefile                  # build, migrate, dev targets
```

---

## Summary Decision Table

| Decision | Choice | Version |
|----------|--------|---------|
| Go router | Chi | v5 |
| Go DB layer | sqlc + pgx | sqlc latest, pgx v5 |
| Dev database | PostgreSQL (Docker) | 16 |
| Migration tool | golang-migrate | v4 |
| Frontend toolchain | Vite + React + TypeScript | Vite 5, React 18 |
| Data fetching | TanStack Query | v5 |
| Infinite scroll | Intersection Observer + useInfiniteQuery | — |
| Image serving | Go net/http FileServer | — |
| Production topology | Go embeds React build | — |
| SSR | None (plain SPA) | — |

---

*Last updated: 2026-04-01 — initial stack research*
