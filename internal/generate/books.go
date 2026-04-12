package generate

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"flos-library/internal/db"
)

// WriteBooks queries all books and writes:
//   - outDir/books.json — array of BookListItem (all books, all shelves)
//   - outDir/books/{slug}.json — BookDetailItem for each book
func WriteBooks(ctx context.Context, queries *db.Queries, outDir string) error {
	rows, err := queries.ListAllBooks(ctx)
	if err != nil {
		return fmt.Errorf("list all books: %w", err)
	}

	listItems := make([]BookListItem, 0, len(rows))
	for _, row := range rows {
		authors := parseRefs(row.AuthorsJson)
		genres := parseRefs(row.GenresJson)

		item := BookListItem{
			Slug:            row.Slug,
			Title:           row.Title,
			CoverPath:       coverURL(row.CoverPath),
			ReadAt:          row.ReadAt,
			PublicationYear: row.PublicationYear,
			PageCount:       row.PageCount,
			Shelf:           row.Shelf,
			Authors:         authors,
			Genres:          genres,
		}
		listItems = append(listItems, item)

		// Write per-slug detail file
		detail := BookDetailItem{
			BookListItem:   item,
			Description:    row.Description,
			Isbn13:         row.Isbn13,
			ReadCount:      row.ReadCount,
			MetadataSource: row.MetadataSource,
		}
		slugPath := filepath.Join(outDir, "books", row.Slug+".json")
		if err := writeJSONFile(slugPath, detail); err != nil {
			return fmt.Errorf("write book detail %s: %w", row.Slug, err)
		}
	}

	if err := writeJSONFile(filepath.Join(outDir, "books.json"), listItems); err != nil {
		return fmt.Errorf("write books.json: %w", err)
	}
	log.Printf("generate: wrote books.json (%d books) + %d detail files", len(rows), len(rows))
	return nil
}
