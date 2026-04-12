# Flo's Library

Personal book tracking site. Syncs from Goodreads RSS, enriches via Google Books, and publishes as a static site on GitHub Pages.

**Live site:** https://florianabel.github.io/flos-library/

## How it works

- A Go CLI (`cmd/generate`) reads from Goodreads RSS, enriches with Google Books metadata, stores everything in SQLite, and writes static JSON to `frontend/public/static/`
- A GitHub Action runs the CLI daily and commits the updated data
- Another GitHub Action builds the React frontend and deploys to GitHub Pages on every push to `master`
- Cover images are committed to the repo and served via jsDelivr CDN

## Local development

**Prerequisites:** Go 1.22+, Node 18+

**1. Create `.env`**

```bash
echo "GOOGLE_BOOKS_API_KEY=your_key_here" > .env
```

**2. Generate data**

```bash
go run ./cmd/generate
```

This syncs from Goodreads, enriches books, and writes JSON to `frontend/public/static/`.

**3. Start the frontend**

```bash
cd frontend
npm install
npm run dev
```

Opens at **http://localhost:5173**. Book data loads from `/static/books.json` — no backend needed.

## Regenerate DB code

Only needed if you change SQL queries in `sql/queries/`:

```bash
sqlc generate
```

## GitHub Actions

| Workflow | Trigger | What it does |
|----------|---------|--------------|
| `sync.yml` | Daily 6am UTC + manual | Runs `go run ./cmd/generate`, commits updated data |
| `deploy.yml` | Push to `master` | Builds frontend, deploys to GitHub Pages |
