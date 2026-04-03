# Flo's Library

Personal book tracking app. Go REST API backend, React frontend (coming later).

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (for PostgreSQL)
- [Go 1.25+](https://go.dev/dl/)
- [golang-migrate](https://github.com/golang-migrate/migrate) — for running DB migrations
- [sqlc](https://sqlc.dev/) — only needed if you change SQL queries

## First-time setup

**1. Start the database**

```bash
docker compose up -d
```

This starts PostgreSQL on port 5432 and Adminer (DB browser) on port 8080.

**2. Copy the env file**

```bash
cp .env.example .env
```

Edit `.env` if needed — defaults work out of the box with the Docker setup.

**3. Run database migrations**

```bash
make migrate
```

Or manually:
```bash
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/floslib?sslmode=disable" up
```

**4. Start the server**

```bash
make dev
```

This starts Docker (if not already running), runs any pending migrations, loads `.env`, and starts the server on **http://localhost:8081**.

## Everyday usage (Docker already running)

Just run:
```bash
make dev
```

This picks up all variables from your `.env` file automatically, runs any pending migrations, and starts the server on **http://localhost:8081**.

> **Without Make:** `set -a; source .env; set +a; go run ./cmd/server` (bash) or set env vars manually on Windows.

## Import your Goodreads books

Export your library from Goodreads (Account → Import/Export → Export Library), then:

```bash
curl -X POST http://localhost:8081/admin/import-csv \
  -F "file=@goodreads_library_export.csv"
```

This returns `{"imported": N}` with the number of books upserted.

## Trigger metadata enrichment

After import, enrich books with descriptions, covers, and authors from Google Books:

```bash
curl -X POST http://localhost:8081/admin/sync
```

Requires `GOOGLE_BOOKS_API_KEY` in your `.env`. Runs in the background — check server logs for progress.

## API endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/books` | Paginated book list (`?cursor=`, `?limit=`) |
| GET | `/api/books/currently-reading` | Books on the currently-reading shelf |
| GET | `/api/books/:slug` | Full book detail with authors and genres |
| GET | `/api/authors` | All authors with book count |
| GET | `/api/authors/:slug` | Author detail + paginated books |
| GET | `/api/genres` | All genres with book count |
| GET | `/api/genres/:slug` | Genre detail + paginated books |
| GET | `/api/years` | Read books grouped by year |
| GET | `/covers/:filename` | Cover images (immutable cache headers) |
| POST | `/admin/sync` | Trigger RSS sync + enrichment |
| POST | `/admin/import-csv` | Import Goodreads CSV export |

## Development tools

```bash
make migrate       # Run pending migrations
make migrate-down  # Roll back last migration
make sqlc          # Regenerate DB code after changing SQL queries
make build         # Compile to ./flos-library binary
```

## Database browser

Adminer runs at **http://localhost:8080** when Docker is up.

- System: PostgreSQL
- Server: db
- Username: postgres
- Password: postgres
- Database: floslib
