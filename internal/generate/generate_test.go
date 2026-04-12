package generate_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"flos-library/internal/db"
	"flos-library/internal/generate"
	_ "modernc.org/sqlite"
)

// openTestDB opens an in-memory SQLite DB and applies the schema.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })

	schema := `
CREATE TABLE IF NOT EXISTS authors (
    id INTEGER PRIMARY KEY, name TEXT NOT NULL, slug TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE TABLE IF NOT EXISTS genres (
    id INTEGER PRIMARY KEY, name TEXT NOT NULL, slug TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE TABLE IF NOT EXISTS books (
    id INTEGER PRIMARY KEY, goodreads_id TEXT NOT NULL UNIQUE, slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL, description TEXT, cover_path TEXT, page_count INTEGER,
    publication_year INTEGER, isbn13 TEXT, metadata_source TEXT NOT NULL DEFAULT 'none',
    read_at TEXT, date_added TEXT, read_count INTEGER NOT NULL DEFAULT 1,
    shelf TEXT NOT NULL DEFAULT 'read',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE TABLE IF NOT EXISTS book_authors (
    book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);
CREATE TABLE IF NOT EXISTS book_genres (
    book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    genre_id INTEGER NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, genre_id)
);`
	if _, err := sqlDB.Exec(schema); err != nil {
		t.Fatalf("apply schema: %v", err)
	}
	return sqlDB
}

func seedTestData(t *testing.T, queries *db.Queries) {
	t.Helper()
	ctx := context.Background()

	// Insert one book
	readAt := "2024-03-15T00:00:00Z"
	coverPath := "data/covers/9780441013593.jpg"
	desc := "A hero's journey on a desert planet."
	pc := int64(896)
	py := int64(1965)
	isbn := "9780441013593"
	book, err := queries.UpsertBook(ctx, db.UpsertBookParams{
		GoodreadsID:     "7767",
		Slug:            "dune",
		Title:           "Dune",
		ReadAt:          &readAt,
		MetadataSource:  "google_books",
		ReadCount:       1,
		Shelf:           "read",
		CoverPath:       &coverPath,
		Description:     &desc,
		PageCount:       &pc,
		PublicationYear: &py,
		Isbn13:          &isbn,
	})
	if err != nil {
		t.Fatalf("upsert book: %v", err)
	}

	// Insert one author
	author, err := queries.UpsertAuthor(ctx, db.UpsertAuthorParams{
		Name: "Frank Herbert",
		Slug: "frank-herbert",
	})
	if err != nil {
		t.Fatalf("upsert author: %v", err)
	}
	if err := queries.LinkBookAuthor(ctx, db.LinkBookAuthorParams{BookID: book.ID, AuthorID: author.ID}); err != nil {
		t.Fatalf("link author: %v", err)
	}

	// Insert one genre
	genre, err := queries.UpsertGenre(ctx, db.UpsertGenreParams{
		Name: "Science Fiction",
		Slug: "science-fiction",
	})
	if err != nil {
		t.Fatalf("upsert genre: %v", err)
	}
	if err := queries.LinkBookGenre(ctx, db.LinkBookGenreParams{BookID: book.ID, GenreID: genre.ID}); err != nil {
		t.Fatalf("link genre: %v", err)
	}
}

func TestWriteBooks_Shape(t *testing.T) {
	sqlDB := openTestDB(t)
	queries := db.New(sqlDB)
	ctx := context.Background()
	seedTestData(t, queries)

	outDir := t.TempDir()
	if err := generate.WriteBooks(ctx, queries, outDir); err != nil {
		t.Fatalf("WriteBooks: %v", err)
	}

	// Verify books.json exists and has correct shape
	booksJSON, err := os.ReadFile(filepath.Join(outDir, "books.json"))
	if err != nil {
		t.Fatalf("read books.json: %v", err)
	}
	var books []generate.BookListItem
	if err := json.Unmarshal(booksJSON, &books); err != nil {
		t.Fatalf("unmarshal books.json: %v", err)
	}
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
	b := books[0]
	if b.Slug != "dune" {
		t.Errorf("slug: got %q, want %q", b.Slug, "dune")
	}
	if b.Title != "Dune" {
		t.Errorf("title: got %q, want %q", b.Title, "Dune")
	}
	if b.Shelf != "read" {
		t.Errorf("shelf: got %q, want %q", b.Shelf, "read")
	}
	// CoverPath should be jsDelivr CDN URL
	wantCover := "https://cdn.jsdelivr.net/gh/florianabel/flos-library@main/data/covers/9780441013593.jpg"
	if b.CoverPath != wantCover {
		t.Errorf("cover_path: got %q, want %q", b.CoverPath, wantCover)
	}
	if len(b.Authors) != 1 || b.Authors[0].Name != "Frank Herbert" {
		t.Errorf("authors: got %v, want [{Frank Herbert frank-herbert}]", b.Authors)
	}
	if len(b.Genres) != 1 || b.Genres[0].Slug != "science-fiction" {
		t.Errorf("genres: got %v, want [{Science Fiction science-fiction}]", b.Genres)
	}

	// Verify per-slug detail file
	detailJSON, err := os.ReadFile(filepath.Join(outDir, "books", "dune.json"))
	if err != nil {
		t.Fatalf("read books/dune.json: %v", err)
	}
	var detail generate.BookDetailItem
	if err := json.Unmarshal(detailJSON, &detail); err != nil {
		t.Fatalf("unmarshal dune.json: %v", err)
	}
	if detail.Description == nil || *detail.Description != "A hero's journey on a desert planet." {
		t.Errorf("description: got %v", detail.Description)
	}
	if detail.Isbn13 == nil || *detail.Isbn13 != "9780441013593" {
		t.Errorf("isbn13: got %v", detail.Isbn13)
	}
}

func TestWriteAuthors_Shape(t *testing.T) {
	sqlDB := openTestDB(t)
	queries := db.New(sqlDB)
	ctx := context.Background()
	seedTestData(t, queries)

	outDir := t.TempDir()
	if err := generate.WriteAuthors(ctx, queries, outDir); err != nil {
		t.Fatalf("WriteAuthors: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "authors.json"))
	if err != nil {
		t.Fatalf("read authors.json: %v", err)
	}
	var authors []generate.AuthorListItem
	if err := json.Unmarshal(data, &authors); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(authors) != 1 || authors[0].Name != "Frank Herbert" {
		t.Errorf("authors: got %v", authors)
	}
	if authors[0].BookCount != 1 {
		t.Errorf("book_count: got %d, want 1", authors[0].BookCount)
	}
}

func TestWriteGenres_Shape(t *testing.T) {
	sqlDB := openTestDB(t)
	queries := db.New(sqlDB)
	ctx := context.Background()
	seedTestData(t, queries)

	outDir := t.TempDir()
	if err := generate.WriteGenres(ctx, queries, outDir); err != nil {
		t.Fatalf("WriteGenres: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "genres.json"))
	if err != nil {
		t.Fatalf("read genres.json: %v", err)
	}
	var genres []generate.GenreListItem
	if err := json.Unmarshal(data, &genres); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(genres) != 1 || genres[0].Name != "Science Fiction" {
		t.Errorf("genres: got %v", genres)
	}
}
