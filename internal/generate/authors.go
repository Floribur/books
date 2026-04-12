package generate

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"flos-library/internal/db"
)

// WriteAuthors queries all authors and writes outDir/authors.json.
func WriteAuthors(ctx context.Context, queries *db.Queries, outDir string) error {
	rows, err := queries.ListAuthors(ctx)
	if err != nil {
		return fmt.Errorf("list authors: %w", err)
	}

	items := make([]AuthorListItem, len(rows))
	for i, row := range rows {
		items[i] = AuthorListItem{
			Name:      row.Name,
			Slug:      row.Slug,
			BookCount: row.BookCount,
		}
	}

	if err := writeJSONFile(filepath.Join(outDir, "authors.json"), items); err != nil {
		return fmt.Errorf("write authors.json: %w", err)
	}
	log.Printf("generate: wrote authors.json (%d authors)", len(rows))
	return nil
}
