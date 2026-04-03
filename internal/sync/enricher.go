package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	stdsync "sync"
	"time"

	"github.com/gosimple/slug"

	"flos-library/internal/db"
)

var enrichHTTPClient = &http.Client{Timeout: 30 * time.Second}

// googleBooksVolume is the parsed response from Google Books API.
type googleBooksVolume struct {
	Items []struct {
		VolumeInfo struct {
			Title         string   `json:"title"`
			Authors       []string `json:"authors"`
			Description   string   `json:"description"`
			PageCount     int      `json:"pageCount"`
			PublishedDate string   `json:"publishedDate"`
			Categories    []string `json:"categories"`
			ImageLinks    struct {
				SmallThumbnail string `json:"smallThumbnail"`
				Thumbnail      string `json:"thumbnail"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

// normalizeTitle strips parenthetical series info and subtitles for fuzzy matching.
//
//	"Sapiens: A Brief History of Humankind" → "sapiens"
//	"Mockingjay (The Hunger Games, #3)"     → "mockingjay"
//	"Halt die Klappe, Kopf! Ein Selfcare…"  → "halt die klappe, kopf! ein selfcare…"
func normalizeTitle(s string) string {
	// Strip parenthetical series info: " (Series, #3)"
	if idx := strings.Index(s, " ("); idx >= 0 {
		s = s[:idx]
	}
	// Strip subtitle after ": "
	if idx := strings.Index(s, ": "); idx >= 0 {
		s = s[:idx]
	}
	return strings.TrimSpace(strings.ToLower(s))
}

// normalizeForCompare does a more aggressive normalization for comparison only:
// strips hyphens/dashes (handles "GAME-OVER" vs "Game Over").
func normalizeForCompare(s string) string {
	s = normalizeTitle(s)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "–", " ")
	s = strings.ReplaceAll(s, "—", " ")
	// Collapse multiple spaces
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.TrimSpace(s)
}

// confidenceGate checks whether a Google Books result is a confident match.
// Uses bidirectional normalized title containment so that either title being
// a prefix/substring of the other is considered a match.
// inputAuthor must be a case-insensitive substring of at least one returnedAuthor.
func confidenceGate(inputTitle, inputAuthor, returnedTitle string, returnedAuthors []string) bool {
	n1 := normalizeForCompare(inputTitle)
	n2 := normalizeForCompare(returnedTitle)
	if n1 == "" || n2 == "" {
		return false
	}
	titleOK := strings.Contains(n1, n2) || strings.Contains(n2, n1)
	if !titleOK {
		return false
	}
	if inputAuthor == "" {
		return true
	}
	for _, a := range returnedAuthors {
		if strings.Contains(strings.ToLower(a), strings.ToLower(inputAuthor)) {
			return true
		}
	}
	return false
}

// fetchGoogleBooks calls the Google Books API with a query and returns up to maxResults volumes.
func fetchGoogleBooks(query string, maxResults int) (*googleBooksVolume, error) {
	apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
	reqURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s&key=%s&maxResults=%d",
		url.QueryEscape(query), apiKey, maxResults)

	resp, err := enrichHTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("google books request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google books HTTP %d", resp.StatusCode)
	}

	var result googleBooksVolume
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode google books: %w", err)
	}
	return &result, nil
}

// EnrichBook fetches metadata for a single book and updates the database.
// On any failure: sets metadata_source='none' and returns (D-06: retry next run).
func EnrichBook(ctx context.Context, queries *db.Queries, book db.Book) {
	apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
	if apiKey == "" {
		log.Printf("enricher: GOOGLE_BOOKS_API_KEY not set, skipping enrichment")
		return
	}

	var (
		description     string
		pageCount       *int32
		publicationYear *int32
		coverPath       *string
		metadataSource  = "none"
	)

	// Primary: ISBN-13 lookup
	var vol *googleBooksVolume
	isbn13 := ""
	if book.Isbn13 != nil && *book.Isbn13 != "" {
		isbn13 = *book.Isbn13
		query := fmt.Sprintf("isbn:%s", isbn13)
		v, err := fetchGoogleBooks(query, 1)
		if err != nil {
			log.Printf("enricher: google books ISBN lookup for %s: %v", book.GoodreadsID, err)
		} else if len(v.Items) > 0 {
			vol = v
		}
	}

	// Fetch primary author from DB (populated by CSV importer) for richer queries
	primaryAuthor := ""
	if a, err := queries.GetFirstAuthorForBook(ctx, book.ID); err == nil {
		primaryAuthor = a
	}

	// Fallback: title + optional author query, try up to 5 results for confidence gate
	if vol == nil || len(vol.Items) == 0 {
		// Strip series/subtitle from query for better Google Books results.
		// e.g. "Mockingjay (The Hunger Games, #3)" → "Mockingjay"
		queryTitle := normalizeTitle(book.Title)
		if queryTitle == "" {
			queryTitle = book.Title
		}
		query := fmt.Sprintf("intitle:%s", queryTitle)
		if primaryAuthor != "" {
			query += fmt.Sprintf("+inauthor:%s", normalizeTitle(primaryAuthor))
		}
		v, err := fetchGoogleBooks(query, 5)
		if err != nil {
			log.Printf("enricher: google books title lookup for %s: %v", book.GoodreadsID, err)
		} else {
			for _, item := range v.Items {
				vi := item.VolumeInfo
				if confidenceGate(book.Title, primaryAuthor, vi.Title, vi.Authors) {
					// Wrap matched item in a volume struct so downstream code works unchanged
					vol = &googleBooksVolume{Items: []struct {
						VolumeInfo struct {
							Title         string   `json:"title"`
							Authors       []string `json:"authors"`
							Description   string   `json:"description"`
							PageCount     int      `json:"pageCount"`
							PublishedDate string   `json:"publishedDate"`
							Categories    []string `json:"categories"`
							ImageLinks    struct {
								SmallThumbnail string `json:"smallThumbnail"`
								Thumbnail      string `json:"thumbnail"`
							} `json:"imageLinks"`
						} `json:"volumeInfo"`
					}{item}}
					break
				}
			}
			if vol == nil && len(v.Items) > 0 {
				log.Printf("enricher: confidence gate failed for book %s (title mismatch: %q vs %q)",
					book.GoodreadsID, book.Title, v.Items[0].VolumeInfo.Title)
			}
		}
	}

	// Extract metadata from volume
	var gbCoverURL string
	if vol != nil && len(vol.Items) > 0 {
		vi := vol.Items[0].VolumeInfo
		description = vi.Description
		if vi.PageCount > 0 {
			pc := int32(vi.PageCount)
			pageCount = &pc
		}
		if len(vi.PublishedDate) >= 4 {
			year := int32(0)
			if _, err := fmt.Sscanf(vi.PublishedDate[:4], "%d", &year); err == nil && year > 0 {
				publicationYear = &year
			}
		}
		gbCoverURL = vi.ImageLinks.Thumbnail
		if gbCoverURL == "" {
			gbCoverURL = vi.ImageLinks.SmallThumbnail
		}
		metadataSource = "google_books"
	}

	// Cover download: Google Books first, OpenLibrary fallback
	destPath := CoverPath(isbn13, book.GoodreadsID)
	if gbCoverURL != "" {
		if path, err := DownloadCover(gbCoverURL, destPath); err != nil {
			log.Printf("enricher: google books cover failed for %s: %v, trying OpenLibrary", book.GoodreadsID, err)
		} else {
			coverPath = &path
		}
	}

	if coverPath == nil && isbn13 != "" {
		// OpenLibrary fallback (D-10)
		time.Sleep(500 * time.Millisecond) // rate limit: 100 req/5min
		if path := TryOpenLibraryCover(isbn13, destPath); path != "" {
			coverPath = &path
		}
	}

	// Build update params using sqlc-generated types
	params := db.UpdateBookEnrichmentParams{
		ID:             book.ID,
		MetadataSource: metadataSource,
	}
	if description != "" {
		params.Description = &description
	}
	if pageCount != nil {
		params.PageCount = pageCount
	}
	if publicationYear != nil {
		params.PublicationYear = publicationYear
	}
	if coverPath != nil {
		params.CoverPath = coverPath
	}

	if err := queries.UpdateBookEnrichment(ctx, params); err != nil {
		log.Printf("enricher: update book %s: %v", book.GoodreadsID, err)
	}

	// Link authors and genres from Google Books
	if vol != nil && len(vol.Items) > 0 {
		vi := vol.Items[0].VolumeInfo
		for _, authorName := range vi.Authors {
			if authorName == "" {
				continue
			}
			authorRow, err := queries.UpsertAuthor(ctx, db.UpsertAuthorParams{
				Name: authorName,
				Slug: slug.Make(authorName),
			})
			if err != nil {
				log.Printf("enricher: upsert author %q for %s: %v", authorName, book.GoodreadsID, err)
				continue
			}
			if err := queries.LinkBookAuthor(ctx, db.LinkBookAuthorParams{BookID: book.ID, AuthorID: authorRow.ID}); err != nil {
				log.Printf("enricher: link author %q for %s: %v", authorName, book.GoodreadsID, err)
			}
		}
		for _, genreName := range vi.Categories {
			if genreName == "" {
				continue
			}
			genreRow, err := queries.UpsertGenre(ctx, db.UpsertGenreParams{
				Name: genreName,
				Slug: slug.Make(genreName),
			})
			if err != nil {
				log.Printf("enricher: upsert genre %q for %s: %v", genreName, book.GoodreadsID, err)
				continue
			}
			if err := queries.LinkBookGenre(ctx, db.LinkBookGenreParams{BookID: book.ID, GenreID: genreRow.ID}); err != nil {
				log.Printf("enricher: link genre %q for %s: %v", genreName, book.GoodreadsID, err)
			}
		}
	}

	time.Sleep(1 * time.Second) // Google Books rate limiting: ~1 req/sec
}

// RunEnricher is the long-lived goroutine that processes unenriched books.
// Waits on trigger channel; exits cleanly on ctx cancellation.
// trigger must be buffered with capacity 1 (prevents double-trigger, RESEARCH.md Pitfall 7).
func RunEnricher(ctx context.Context, wg *stdsync.WaitGroup, queries *db.Queries, trigger <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-trigger:
			log.Println("enricher: triggered, fetching unenriched books")
			books, err := queries.GetUnenrichedBooks(ctx)
			if err != nil {
				log.Printf("enricher: fetch unenriched: %v", err)
				continue
			}
			log.Printf("enricher: processing %d unenriched books", len(books))
			for _, book := range books {
				select {
				case <-ctx.Done():
					log.Println("enricher: context cancelled mid-run, stopping")
					return
				default:
					EnrichBook(ctx, queries, book)
				}
			}
			log.Println("enricher: batch complete")
		case <-ctx.Done():
			log.Println("enricher: context cancelled, exiting")
			return
		}
	}
}
