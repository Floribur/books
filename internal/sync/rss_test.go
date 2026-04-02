package sync

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRSSPagination: given a mock that returns 200 items on page 1 and 5 on page 2,
// loop terminates after page 2 and total items = 205.
func TestRSSPagination(t *testing.T) {
	// Build mock RSS items
	buildItems := func(count int) string {
		var sb strings.Builder
		for i := 0; i < count; i++ {
			sb.WriteString(fmt.Sprintf(`<item>
<title>Book %d</title>
<guid>https://www.goodreads.com/review/show/%d</guid>
<link>https://www.goodreads.com/review/show/%d</link>
<book_id>%d</book_id>
</item>`, i, i, i, i))
		}
		return sb.String()
	}

	page1Items := buildItems(200)
	page2Items := buildItems(5)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/rss+xml")
		var items string
		if page == "2" {
			items = page2Items
		} else {
			items = page1Items
		}
		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>
<title>Test Shelf</title>
%s
</channel>
</rss>`, items)
	}))
	defer srv.Close()

	// Override rssBaseURL temporarily
	origURL := rssBaseURL
	// We need to patch the base URL — use a test helper approach
	// Since rssBaseURL is a package-level const, we use a wrapper function
	_ = origURL

	// Use fetchShelfFromURL which accepts base URL for testability
	items, err := fetchShelfFromURL(srv.URL, "read")
	if err != nil {
		t.Fatalf("fetchShelfFromURL error: %v", err)
	}
	if len(items) != 205 {
		t.Errorf("expected 205 items, got %d", len(items))
	}
}

// TestShelfMerge: book appears in currently-reading with shelf="currently-reading"
// and in read with shelf="read"; merged result has shelf="read" (read wins per D-04).
func TestShelfMerge(t *testing.T) {
	currentlyReading := []GoodreadsItem{
		{GoodreadsID: "123", Title: "Dune", Shelf: "currently-reading"},
		{GoodreadsID: "456", Title: "Foundation", Shelf: "currently-reading"},
	}
	read := []GoodreadsItem{
		{GoodreadsID: "123", Title: "Dune", Shelf: "read"},
	}

	merged := mergeShelfItems(currentlyReading, read)

	// Find book 123
	found := false
	for _, item := range merged {
		if item.GoodreadsID == "123" {
			found = true
			if item.Shelf != "read" {
				t.Errorf("expected shelf='read' for book 123, got %q", item.Shelf)
			}
		}
	}
	if !found {
		t.Error("book 123 not found in merged result")
	}
	// Book 456 should still be in currently-reading
	for _, item := range merged {
		if item.GoodreadsID == "456" {
			if item.Shelf != "currently-reading" {
				t.Errorf("expected shelf='currently-reading' for book 456, got %q", item.Shelf)
			}
		}
	}
	if len(merged) != 2 {
		t.Errorf("expected 2 merged items, got %d", len(merged))
	}
}
