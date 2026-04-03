package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"flos-library/internal/db"
)

const defaultPageSize = 24
const maxPageSize = 100

// PublicHandlers holds dependencies for all public read-only API endpoints.
type PublicHandlers struct {
	Store   BookStore
	SiteURL string // from SITE_URL env var; used for og:image absolute URLs in Plan 02-02
}

// writeJSON writes a JSON response with Content-Type: application/json.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes {"error": msg} with the given status code.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// parseCursorParam returns a pgtype.Timestamptz and int64 id from the ?cursor= query param.
// If no cursor is present, returns null Timestamptz and id=0 (first page).
func parseCursorParam(r *http.Request) (pgtype.Timestamptz, int64, error) {
	s := r.URL.Query().Get("cursor")
	if s == "" {
		return pgtype.Timestamptz{Valid: false}, 0, nil
	}
	t, id, err := decodeCursor(s)
	if err != nil {
		return pgtype.Timestamptz{}, 0, err
	}
	ts := pgtype.Timestamptz{Time: t, Valid: true}
	return ts, id, nil
}

// parseLimit reads ?limit= from query params; defaults to defaultPageSize, max maxPageSize.
func parseLimit(r *http.Request) int32 {
	s := r.URL.Query().Get("limit")
	if s == "" {
		return defaultPageSize
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return defaultPageSize
	}
	if n > maxPageSize {
		return maxPageSize
	}
	return int32(n)
}

// unmarshalRefs converts a json_agg interface{} value (from pgx scan) into a typed slice.
// pgx returns JSON aggregates as []interface{} — we re-marshal then unmarshal into the target type.
func unmarshalRefs[T any](raw interface{}) []T {
	if raw == nil {
		return []T{}
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return []T{}
	}
	var out []T
	if err := json.Unmarshal(b, &out); err != nil {
		return []T{}
	}
	return out
}

// toTimePtr converts pgtype.Timestamptz to *time.Time for JSON serialization.
// Returns nil if the timestamp is not valid (NULL in DB).
func toTimePtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time.UTC()
	return &t
}

// --- Handler Methods ---

// GetBooks handles GET /api/books
// Query params: ?cursor=<token>&limit=<n>
// Returns: PaginatedResponse with BookListItem slice. Per API-01, D-01, D-02, D-04.
func (h *PublicHandlers) GetBooks(w http.ResponseWriter, r *http.Request) {
	cursorTS, cursorID, err := parseCursorParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid cursor")
		return
	}
	limit := parseLimit(r)
	rows, err := h.Store.ListBooksPaginated(r.Context(), db.ListBooksPaginatedParams{
		Column1: cursorTS,
		Column2: cursorID,
		Limit:   limit + 1, // fetch one extra to detect has_more
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list books")
		return
	}
	hasMore := len(rows) > int(limit)
	if hasMore {
		rows = rows[:limit]
	}
	items := make([]BookListItem, len(rows))
	for i, row := range rows {
		items[i] = BookListItem{
			Slug:            row.Slug,
			Title:           row.Title,
			CoverPath:       row.CoverPath,
			ReadAt:          toTimePtr(row.ReadAt),
			PublicationYear: row.PublicationYear,
			Authors:         unmarshalRefs[AuthorRef](row.Authors),
			Genres:          unmarshalRefs[GenreRef](row.Genres),
		}
	}
	var nextCursor *string
	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]
		s := encodeCursor(last.ReadAt.Time.UTC(), last.ID)
		nextCursor = &s
	}
	writeJSON(w, http.StatusOK, PaginatedResponse{Items: items, NextCursor: nextCursor, HasMore: hasMore})
}

// GetCurrentlyReading handles GET /api/books/currently-reading
// MUST be registered BEFORE /api/books/{slug} in the router. Per API-02, pitfall 5.
func (h *PublicHandlers) GetCurrentlyReading(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Store.GetCurrentlyReading(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list currently reading")
		return
	}
	items := make([]BookListItem, len(rows))
	for i, row := range rows {
		items[i] = BookListItem{
			Slug:            row.Slug,
			Title:           row.Title,
			CoverPath:       row.CoverPath,
			ReadAt:          toTimePtr(row.ReadAt),
			PublicationYear: row.PublicationYear,
			Authors:         unmarshalRefs[AuthorRef](row.Authors),
			Genres:          unmarshalRefs[GenreRef](row.Genres),
		}
	}
	writeJSON(w, http.StatusOK, items)
}

// GetBookBySlug handles GET /api/books/{slug}. Per API-03, D-05.
func (h *PublicHandlers) GetBookBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	row, err := h.Store.GetBookDetailBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get book")
		return
	}
	detail := BookDetail{
		Slug:            row.Slug,
		Title:           row.Title,
		CoverPath:       row.CoverPath,
		ReadAt:          toTimePtr(row.ReadAt),
		PublicationYear: row.PublicationYear,
		Description:     row.Description,
		PageCount:       row.PageCount,
		Isbn13:          row.Isbn13,
		ReadCount:       row.ReadCount,
		Shelf:           row.Shelf,
		MetadataSource:  row.MetadataSource,
		Authors:         unmarshalRefs[AuthorRef](row.Authors),
		Genres:          unmarshalRefs[GenreRef](row.Genres),
	}
	writeJSON(w, http.StatusOK, detail)
}

// GetAuthors handles GET /api/authors. Returns plain array. Per API-04, D-03, D-06.
func (h *PublicHandlers) GetAuthors(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Store.ListAuthors(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list authors")
		return
	}
	items := make([]AuthorListItem, len(rows))
	for i, row := range rows {
		items[i] = AuthorListItem{Name: row.Name, Slug: row.Slug, BookCount: row.BookCount}
	}
	writeJSON(w, http.StatusOK, items)
}

// GetAuthorBySlug handles GET /api/authors/{slug}. Per API-05, D-07.
func (h *PublicHandlers) GetAuthorBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	author, err := h.Store.GetAuthorBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get author")
		return
	}
	cursorTS, cursorID, err := parseCursorParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid cursor")
		return
	}
	limit := parseLimit(r)
	bookRows, err := h.Store.ListBooksByAuthor(r.Context(), db.ListBooksByAuthorParams{
		Slug:    slug,
		Column2: cursorTS,
		Column3: cursorID,
		Limit:   limit + 1,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list author books")
		return
	}
	hasMore := len(bookRows) > int(limit)
	if hasMore {
		bookRows = bookRows[:limit]
	}
	books := make([]BookListItem, len(bookRows))
	for i, row := range bookRows {
		books[i] = BookListItem{
			Slug: row.Slug, Title: row.Title, CoverPath: row.CoverPath,
			ReadAt: toTimePtr(row.ReadAt), PublicationYear: row.PublicationYear,
			Authors: unmarshalRefs[AuthorRef](row.Authors),
			Genres:  unmarshalRefs[GenreRef](row.Genres),
		}
	}
	var nextCursor *string
	if hasMore && len(bookRows) > 0 {
		last := bookRows[len(bookRows)-1]
		s := encodeCursor(last.ReadAt.Time.UTC(), last.ID)
		nextCursor = &s
	}
	writeJSON(w, http.StatusOK, AuthorDetail{
		Name: author.Name, Slug: author.Slug,
		Books: PaginatedResponse{Items: books, NextCursor: nextCursor, HasMore: hasMore},
	})
}

// GetGenres handles GET /api/genres. Returns plain array. Per API-06, D-03, D-06.
func (h *PublicHandlers) GetGenres(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Store.ListGenres(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list genres")
		return
	}
	items := make([]GenreListItem, len(rows))
	for i, row := range rows {
		items[i] = GenreListItem{Name: row.Name, Slug: row.Slug, BookCount: row.BookCount}
	}
	writeJSON(w, http.StatusOK, items)
}

// GetGenreBySlug handles GET /api/genres/{slug}. Per API-07, D-07.
func (h *PublicHandlers) GetGenreBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	genre, err := h.Store.GetGenreBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get genre")
		return
	}
	cursorTS, cursorID, err := parseCursorParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid cursor")
		return
	}
	limit := parseLimit(r)
	bookRows, err := h.Store.ListBooksByGenre(r.Context(), db.ListBooksByGenreParams{
		Slug:    slug,
		Column2: cursorTS,
		Column3: cursorID,
		Limit:   limit + 1,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list genre books")
		return
	}
	hasMore := len(bookRows) > int(limit)
	if hasMore {
		bookRows = bookRows[:limit]
	}
	books := make([]BookListItem, len(bookRows))
	for i, row := range bookRows {
		books[i] = BookListItem{
			Slug: row.Slug, Title: row.Title, CoverPath: row.CoverPath,
			ReadAt: toTimePtr(row.ReadAt), PublicationYear: row.PublicationYear,
			Authors: unmarshalRefs[AuthorRef](row.Authors),
			Genres:  unmarshalRefs[GenreRef](row.Genres),
		}
	}
	var nextCursor *string
	if hasMore && len(bookRows) > 0 {
		last := bookRows[len(bookRows)-1]
		s := encodeCursor(last.ReadAt.Time.UTC(), last.ID)
		nextCursor = &s
	}
	writeJSON(w, http.StatusOK, GenreDetail{
		Name: genre.Name, Slug: genre.Slug,
		Books: PaginatedResponse{Items: books, NextCursor: nextCursor, HasMore: hasMore},
	})
}

// GetYears handles GET /api/years. Returns plain array. Per API-08, D-03.
func (h *PublicHandlers) GetYears(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Store.ListYears(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list years")
		return
	}
	items := make([]YearCount, len(rows))
	for i, row := range rows {
		items[i] = YearCount{Year: int(row.Year), BookCount: row.BookCount}
	}
	writeJSON(w, http.StatusOK, items)
}
