package sync

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mmcdole/gofeed"

	"flos-library/internal/db"
)

const (
	rssBaseURL  = "https://www.goodreads.com/review/list_rss/79499864"
	rssPerPage  = 200
)

// GoodreadsItem holds parsed fields from one RSS item.
type GoodreadsItem struct {
	GoodreadsID string
	Title       string
	AuthorName  string
	ISBN        string
	ImageURL    string
	ReadAt      time.Time
	DateAdded   time.Time
	Shelf       string
	ReadCount   int
}

// extractCustom reads a field from item.Custom (the actual Goodreads RSS format).
// Goodreads delivers book fields via item.Custom, not item.Extensions.
func extractCustom(item *gofeed.Item, key string) string {
	if item.Custom == nil {
		return ""
	}
	return strings.TrimSpace(item.Custom[key])
}

// FetchShelf fetches all pages for a given shelf name (e.g., "read", "currently-reading").
func FetchShelf(shelfName string) ([]GoodreadsItem, error) {
	return fetchShelfFromURL(rssBaseURL, shelfName)
}

// fetchShelfFromURL is the testable implementation of FetchShelf.
// It accepts an arbitrary base URL so tests can inject a mock server.
func fetchShelfFromURL(baseURL, shelfName string) ([]GoodreadsItem, error) {
	fp := gofeed.NewParser()
	var all []GoodreadsItem
	page := 1
	for {
		url := fmt.Sprintf("%s?shelf=%s&per_page=%d&page=%d", baseURL, shelfName, rssPerPage, page)
		feed, err := fp.ParseURL(url)
		if err != nil {
			return nil, fmt.Errorf("fetch shelf %s page %d: %w", shelfName, page, err)
		}
		for _, item := range feed.Items {
			gi := parseItem(item)
			gi.Shelf = shelfName
			all = append(all, gi)
		}
		if len(feed.Items) < rssPerPage {
			break // last page
		}
		page++
		time.Sleep(500 * time.Millisecond) // be polite to Goodreads
	}
	return all, nil
}

// parseItem extracts GoodreadsItem data from a gofeed.Item using item.Custom.
func parseItem(item *gofeed.Item) GoodreadsItem {
	gi := GoodreadsItem{
		GoodreadsID: extractCustom(item, "book_id"),
		Title:       item.Title,
		AuthorName:  strings.Join(strings.Fields(extractCustom(item, "author_name")), " "), // normalize whitespace
		ISBN:        extractCustom(item, "isbn"),
		ImageURL:    extractCustom(item, "book_large_image_url"),
	}
	if gi.ImageURL == "" {
		gi.ImageURL = extractCustom(item, "book_image_url")
	}
	// Parse dates — Goodreads uses RFC1123Z variants; try multiple formats
	readAtFormats := []string{
		time.RFC1123Z,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 +0000",
	}
	if raw := extractCustom(item, "user_read_at"); raw != "" {
		for _, fmt := range readAtFormats {
			if t, err := time.Parse(fmt, raw); err == nil {
				gi.ReadAt = t
				break
			}
		}
	}
	addedFormats := []string{
		time.RFC1123Z,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 +0000",
	}
	if raw := extractCustom(item, "user_date_added"); raw != "" {
		for _, fmt := range addedFormats {
			if t, err := time.Parse(fmt, raw); err == nil {
				gi.DateAdded = t
				break
			}
		}
	}
	gi.ReadCount = 1
	return gi
}

// mergeShelfItems merges currently-reading and read items; read wins on conflict (D-04).
func mergeShelfItems(currentlyReading, read []GoodreadsItem) []GoodreadsItem {
	merged := make(map[string]GoodreadsItem)
	// Process currently-reading first
	for _, item := range currentlyReading {
		merged[item.GoodreadsID] = item
	}
	// read overrides (D-04: read wins)
	for _, item := range read {
		merged[item.GoodreadsID] = item
	}
	result := make([]GoodreadsItem, 0, len(merged))
	for _, item := range merged {
		result = append(result, item)
	}
	return result
}

// SyncRSS fetches both shelves, merges, and upserts all books to the database.
func SyncRSS(ctx context.Context, queries *db.Queries) error {
	log.Println("sync: starting RSS sync")

	currentlyReading, err := FetchShelf("currently-reading")
	if err != nil {
		log.Printf("sync: error fetching currently-reading shelf: %v (continuing)", err)
		currentlyReading = nil
	}

	read, err := FetchShelf("read")
	if err != nil {
		return fmt.Errorf("sync: error fetching read shelf: %w", err)
	}

	items := mergeShelfItems(currentlyReading, read)
	log.Printf("sync: %d items after merge", len(items))

	// Load existing slugs for collision resolution
	existing, err := queries.GetAllGoodreadsIDs(ctx)
	if err != nil {
		return fmt.Errorf("sync: load existing IDs: %w", err)
	}
	existingSlugs := make(map[string]struct{}, len(existing))
	grIDToSlug := make(map[string]string, len(existing))
	for _, e := range existing {
		existingSlugs[e.Slug] = struct{}{}
		grIDToSlug[e.GoodreadsID] = e.Slug
	}

	for _, item := range items {
		if item.GoodreadsID == "" {
			continue
		}
		// Reuse existing slug if book already in DB; otherwise generate new
		bookSlug, exists := grIDToSlug[item.GoodreadsID]
		if !exists {
			// Extract author surname for collision fallback
			authorParts := strings.Fields(item.AuthorName)
			authorSurname := ""
			if len(authorParts) > 0 {
				authorSurname = authorParts[len(authorParts)-1]
			}
			bookSlug = GenerateSlug(item.Title, 0, authorSurname, existingSlugs)
			existingSlugs[bookSlug] = struct{}{}
			grIDToSlug[item.GoodreadsID] = bookSlug
		}

		var readAt pgtype.Timestamptz
		if !item.ReadAt.IsZero() {
			readAt = pgtype.Timestamptz{Time: item.ReadAt, Valid: true}
		}
		var dateAdded pgtype.Timestamptz
		if !item.DateAdded.IsZero() {
			dateAdded = pgtype.Timestamptz{Time: item.DateAdded, Valid: true}
		}

		params := db.UpsertBookParams{
			GoodreadsID:    item.GoodreadsID,
			Slug:           bookSlug,
			Title:          item.Title,
			MetadataSource: "none",
			ReadAt:         readAt,
			DateAdded:      dateAdded,
			ReadCount:      int32(item.ReadCount),
			Shelf:          item.Shelf,
		}
		if item.ISBN != "" {
			params.Isbn13 = &item.ISBN
		}
		book, err := queries.UpsertBook(ctx, params)
		if err != nil {
			log.Printf("sync: upsert book %s: %v", item.GoodreadsID, err)
			continue
		}

		// Link author from RSS feed (ON CONFLICT DO NOTHING — idempotent)
		if item.AuthorName != "" {
			authorRow, err := queries.UpsertAuthor(ctx, db.UpsertAuthorParams{
				Name: item.AuthorName,
				Slug: slug.Make(item.AuthorName),
			})
			if err != nil {
				log.Printf("sync: upsert author %q for %s: %v", item.AuthorName, item.GoodreadsID, err)
			} else if err := queries.LinkBookAuthor(ctx, db.LinkBookAuthorParams{
				BookID: book.ID, AuthorID: authorRow.ID,
			}); err != nil {
				log.Printf("sync: link author %q to %s: %v", item.AuthorName, item.GoodreadsID, err)
			}
		}
	}

	log.Printf("sync: completed, upserted %d books", len(items))
	return nil
}
