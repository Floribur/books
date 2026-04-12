package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"path/filepath"

	"flos-library/internal/db"
	"flos-library/internal/generate"
	syncp "flos-library/internal/sync"

	_ "modernc.org/sqlite"
)

// sqliteSchema is applied on every startup (CREATE TABLE IF NOT EXISTS is idempotent).
const sqliteSchema = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS authors (
    id         INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS genres (
    id         INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS books (
    id               INTEGER PRIMARY KEY,
    goodreads_id     TEXT    NOT NULL UNIQUE,
    slug             TEXT    NOT NULL UNIQUE,
    title            TEXT    NOT NULL,
    description      TEXT,
    cover_path       TEXT,
    page_count       INTEGER,
    publication_year INTEGER,
    isbn13           TEXT,
    metadata_source  TEXT NOT NULL DEFAULT 'none',
    read_at          TEXT,
    date_added       TEXT,
    read_count       INTEGER NOT NULL DEFAULT 1,
    shelf            TEXT NOT NULL DEFAULT 'read',
    created_at       TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at       TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_books_slug         ON books(slug);
CREATE INDEX IF NOT EXISTS idx_books_read_at      ON books(read_at DESC);
CREATE INDEX IF NOT EXISTS idx_books_goodreads_id ON books(goodreads_id);
CREATE INDEX IF NOT EXISTS idx_books_shelf        ON books(shelf);

CREATE TABLE IF NOT EXISTS book_authors (
    book_id   INTEGER NOT NULL REFERENCES books(id)   ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

CREATE TABLE IF NOT EXISTS book_genres (
    book_id  INTEGER NOT NULL REFERENCES books(id)  ON DELETE CASCADE,
    genre_id INTEGER NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, genre_id)
);
`

func main() {
	migrateOnly := flag.Bool("migrate-only", false, "apply schema and exit without syncing")
	flag.Parse()

	// Ensure data/ directory exists
	if err := os.MkdirAll("data/covers", 0755); err != nil {
		log.Fatalf("create data/covers: %v", err)
	}

	sqlDB, err := sql.Open("sqlite", "data/floslib.db")
	if err != nil {
		log.Fatalf("open sqlite: %v", err)
	}
	defer sqlDB.Close()

	// Apply schema (idempotent)
	if _, err := sqlDB.Exec(sqliteSchema); err != nil {
		log.Fatalf("apply schema: %v", err)
	}
	log.Println("generate: schema applied")

	if *migrateOnly {
		log.Println("generate: --migrate-only set, exiting")
		return
	}

	queries := db.New(sqlDB)
	ctx := context.Background()

	// Step 1: RSS sync
	log.Println("generate: running RSS sync")
	if err := syncp.SyncRSS(ctx, queries); err != nil {
		log.Fatalf("sync RSS: %v", err)
	}

	// Step 2: Enrich unenriched books
	log.Println("generate: enriching unenriched books")
	books, err := queries.GetUnenrichedBooks(ctx)
	if err != nil {
		log.Fatalf("get unenriched books: %v", err)
	}
	log.Printf("generate: %d books to enrich", len(books))
	for _, book := range books {
		select {
		case <-ctx.Done():
			log.Fatalf("context cancelled")
		default:
			syncp.EnrichBook(ctx, queries, book)
		}
	}

	// Step 3: Write static JSON files
	outDir := filepath.Join("frontend", "public", "static")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("create output dir: %v", err)
	}

	if err := generate.WriteBooks(ctx, queries, outDir); err != nil {
		log.Fatalf("write books: %v", err)
	}
	if err := generate.WriteAuthors(ctx, queries, outDir); err != nil {
		log.Fatalf("write authors: %v", err)
	}
	if err := generate.WriteGenres(ctx, queries, outDir); err != nil {
		log.Fatalf("write genres: %v", err)
	}

	log.Println("generate: all done")
}
