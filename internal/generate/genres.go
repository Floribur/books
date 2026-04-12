package generate

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"flos-library/internal/db"
)

// WriteGenres queries all genres and writes outDir/genres.json.
func WriteGenres(ctx context.Context, queries *db.Queries, outDir string) error {
	rows, err := queries.ListGenres(ctx)
	if err != nil {
		return fmt.Errorf("list genres: %w", err)
	}

	items := make([]GenreListItem, len(rows))
	for i, row := range rows {
		items[i] = GenreListItem{
			Name:      row.Name,
			Slug:      row.Slug,
			BookCount: row.BookCount,
		}
	}

	if err := writeJSONFile(filepath.Join(outDir, "genres.json"), items); err != nil {
		return fmt.Errorf("write genres.json: %w", err)
	}
	log.Printf("generate: wrote genres.json (%d genres)", len(rows))
	return nil
}
