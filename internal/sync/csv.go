package sync

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"flos-library/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// unquoteISBN strips Excel-style formula quoting from Goodreads CSV ISBNs.
// Input: ="9780385472579"  Output: 9780385472579
func unquoteISBN(raw string) string {
	s := strings.TrimPrefix(raw, "=\"")
	s = strings.TrimSuffix(s, "\"")
	return strings.TrimSpace(s)
}

// ImportCSV parses a Goodreads CSV export and upserts all books.
// r should be the multipart file reader from the HTTP request.
func ImportCSV(ctx context.Context, queries *db.Queries, r io.Reader) (int, error) {
	reader := csv.NewReader(r)
	headers, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("csv: read header: %w", err)
	}
	// Build column index map
	colIdx := make(map[string]int, len(headers))
	for i, h := range headers {
		colIdx[strings.TrimSpace(h)] = i
	}

	// Validate required columns exist
	required := []string{"Book Id", "Title", "Author", "ISBN13", "Exclusive Shelf"}
	for _, col := range required {
		if _, ok := colIdx[col]; !ok {
			return 0, fmt.Errorf("csv: missing required column %q", col)
		}
	}

	// Load existing slugs for collision resolution
	existing, err := queries.GetAllGoodreadsIDs(ctx)
	if err != nil {
		return 0, fmt.Errorf("csv: load existing IDs: %w", err)
	}
	existingSlugs := make(map[string]struct{}, len(existing))
	grIDToSlug := make(map[string]string, len(existing))
	for _, e := range existing {
		existingSlugs[e.Slug] = struct{}{}
		grIDToSlug[e.GoodreadsID] = e.Slug
	}

	count := 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("csv: skip row: %v", err)
			continue
		}
		get := func(col string) string {
			i, ok := colIdx[col]
			if !ok || i >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[i])
		}

		goodreadsID := get("Book Id")
		if goodreadsID == "" {
			continue
		}
		title := get("Title")
		author := get("Author")
		isbn13 := unquoteISBN(get("ISBN13"))
		shelf := get("Exclusive Shelf")
		if shelf == "" {
			shelf = "read"
		}

		// Parse read date
		var readAt pgtype.Timestamptz
		if raw := get("Date Read"); raw != "" {
			if t, err := time.Parse("2006/01/02", raw); err == nil {
				readAt = pgtype.Timestamptz{Time: t, Valid: true}
			}
		}

		// Parse read count
		rc := 1
		if raw := get("Read Count"); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil && n > 0 {
				rc = n
			}
		}

		// Slug
		bookSlug, exists := grIDToSlug[goodreadsID]
		if !exists {
			parts := strings.Fields(author)
			surname := ""
			if len(parts) > 0 {
				surname = parts[len(parts)-1]
			}
			bookSlug = GenerateSlug(title, 0, surname, existingSlugs)
			existingSlugs[bookSlug] = struct{}{}
			grIDToSlug[goodreadsID] = bookSlug
		}

		params := db.UpsertBookParams{
			GoodreadsID:    goodreadsID,
			Slug:           bookSlug,
			Title:          title,
			MetadataSource: "none",
			ReadAt:         readAt,
			ReadCount:      int32(rc),
			Shelf:          shelf,
		}
		if isbn13 != "" {
			params.Isbn13 = &isbn13
		}
		if _, err := queries.UpsertBook(ctx, params); err != nil {
			log.Printf("csv: upsert %s: %v", goodreadsID, err)
			continue
		}
		count++
	}
	log.Printf("csv: imported %d books", count)
	return count, nil
}
