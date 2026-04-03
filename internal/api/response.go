package api

import "time"

// AuthorRef is the inline author shape in book list and detail responses.
type AuthorRef struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// GenreRef is the inline genre shape in book list and detail responses.
type GenreRef struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// BookListItem is the shape returned by GET /api/books,
// GET /api/books/currently-reading, and paginated author/genre detail endpoints.
// Per decisions D-04 and D-07.
type BookListItem struct {
	Slug            string      `json:"slug"`
	Title           string      `json:"title"`
	CoverPath       *string     `json:"cover_path"`
	ReadAt          *time.Time  `json:"read_at"`
	PublicationYear *int32      `json:"publication_year"`
	Authors         []AuthorRef `json:"authors"`
	Genres          []GenreRef  `json:"genres"`
}

// BookDetail is the shape returned by GET /api/books/:slug. Per decision D-05.
type BookDetail struct {
	Slug            string      `json:"slug"`
	Title           string      `json:"title"`
	CoverPath       *string     `json:"cover_path"`
	ReadAt          *time.Time  `json:"read_at"`
	PublicationYear *int32      `json:"publication_year"`
	Description     *string     `json:"description"`
	PageCount       *int32      `json:"page_count"`
	Isbn13          *string     `json:"isbn13"`
	ReadCount       int32       `json:"read_count"`
	Shelf           string      `json:"shelf"`
	MetadataSource  string      `json:"metadata_source"`
	Authors         []AuthorRef `json:"authors"`
	Genres          []GenreRef  `json:"genres"`
}

// AuthorListItem is the shape returned by GET /api/authors. Per decision D-06.
type AuthorListItem struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	BookCount int64  `json:"book_count"`
}

// AuthorDetail is the shape returned by GET /api/authors/:slug. Per decision D-07.
type AuthorDetail struct {
	Name  string            `json:"name"`
	Slug  string            `json:"slug"`
	Books PaginatedResponse `json:"books"`
}

// GenreListItem is the shape returned by GET /api/genres. Per decision D-06.
type GenreListItem struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	BookCount int64  `json:"book_count"`
}

// GenreDetail is the shape returned by GET /api/genres/:slug. Per decision D-07.
type GenreDetail struct {
	Name  string            `json:"name"`
	Slug  string            `json:"slug"`
	Books PaginatedResponse `json:"books"`
}

// YearCount is one item in the GET /api/years array response.
type YearCount struct {
	Year      int   `json:"year"`
	BookCount int64 `json:"book_count"`
}

// PaginatedResponse is the envelope for all paginated list endpoints. Per decision D-01.
type PaginatedResponse struct {
	Items      any     `json:"items"`
	NextCursor *string `json:"next_cursor"` // nil when no more pages (has_more=false)
	HasMore    bool    `json:"has_more"`
}
