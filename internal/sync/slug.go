package sync

import (
	"fmt"

	"github.com/gosimple/slug"
)

// GenerateSlug creates a URL-safe slug from a book title.
// Collision resolution: append year, then authorSurname.
// existingSlugs is the set of slugs already in the DB.
func GenerateSlug(title string, year int, authorSurname string, existingSlugs map[string]struct{}) string {
	base := slug.Make(title)
	if _, exists := existingSlugs[base]; !exists {
		return base
	}
	// Append year
	if year > 0 {
		withYear := fmt.Sprintf("%s-%d", base, year)
		if _, exists := existingSlugs[withYear]; !exists {
			return withYear
		}
	}
	// Append author surname
	withAuthor := fmt.Sprintf("%s-%s", base, slug.Make(authorSurname))
	return withAuthor
}
