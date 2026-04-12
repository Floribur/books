package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// cdnBase is the jsDelivr CDN URL prefix for cover images.
// IMPORTANT: Replace "florianabel/flos-library" with the actual GitHub owner/repo.
const cdnBase = "https://cdn.jsdelivr.net/gh/floribur/books@master/data/covers/"

// Ref is the inline shape for author/genre references embedded in book JSON.
type Ref struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// BookListItem is the shape written to books.json (one entry per book).
type BookListItem struct {
	Slug            string  `json:"slug"`
	Title           string  `json:"title"`
	CoverPath       string  `json:"cover_path"`
	ReadAt          *string `json:"read_at"`
	PublicationYear *int64  `json:"publication_year"`
	PageCount       *int64  `json:"page_count"`
	Shelf           string  `json:"shelf"`
	Authors         []Ref   `json:"authors"`
	Genres          []Ref   `json:"genres"`
}

// BookDetailItem is the shape written to books/{slug}.json.
type BookDetailItem struct {
	BookListItem
	Description    *string `json:"description"`
	Isbn13         *string `json:"isbn13"`
	ReadCount      int64   `json:"read_count"`
	MetadataSource string  `json:"metadata_source"`
}

// AuthorListItem is one entry in authors.json.
type AuthorListItem struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	BookCount int64  `json:"book_count"`
}

// GenreListItem is one entry in genres.json.
type GenreListItem struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	BookCount int64  `json:"book_count"`
}

// coverURL converts a DB cover_path like "data/covers/isbn.jpg" to a jsDelivr CDN URL.
// Returns empty string if path is nil or empty.
func coverURL(dbPath *string) string {
	if dbPath == nil || *dbPath == "" {
		return ""
	}
	return cdnBase + filepath.Base(*dbPath)
}

// parseRefs parses a JSON array of [{name, slug}] objects returned by SQLite json_group_array.
// The sqlc-generated field is interface{} — at runtime it is either a string (JSON) or nil.
// Returns an empty slice on nil input or parse failure.
func parseRefs(v interface{}) []Ref {
	if v == nil {
		return []Ref{}
	}
	s, ok := v.(string)
	if !ok {
		return []Ref{}
	}
	var refs []Ref
	if err := json.Unmarshal([]byte(s), &refs); err != nil {
		return []Ref{}
	}
	return refs
}

// writeJSONFile marshals v to JSON and writes it to path, creating parent dirs as needed.
func writeJSONFile(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
