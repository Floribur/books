package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"flos-library/internal/api"
	"flos-library/internal/db"
)

// mockStore is a test double for BookStore.
type mockStore struct {
	books            []db.ListBooksPaginatedRow
	detail           db.GetBookDetailBySlugRow
	currentlyReading []db.GetCurrentlyReadingRow
	authors          []db.ListAuthorsRow
	author           db.GetAuthorBySlugRow
	authorBooks      []db.ListBooksByAuthorRow
	genres           []db.ListGenresRow
	genre            db.GetGenreBySlugRow
	genreBooks       []db.ListBooksByGenreRow
	years            []db.ListYearsRow
	err              error
}

func (m *mockStore) ListBooksPaginated(_ context.Context, _ db.ListBooksPaginatedParams) ([]db.ListBooksPaginatedRow, error) {
	return m.books, m.err
}

func (m *mockStore) GetCurrentlyReading(_ context.Context) ([]db.GetCurrentlyReadingRow, error) {
	return m.currentlyReading, m.err
}

func (m *mockStore) GetBookDetailBySlug(_ context.Context, _ string) (db.GetBookDetailBySlugRow, error) {
	return m.detail, m.err
}

func (m *mockStore) ListAuthors(_ context.Context) ([]db.ListAuthorsRow, error) {
	return m.authors, m.err
}

func (m *mockStore) GetAuthorBySlug(_ context.Context, _ string) (db.GetAuthorBySlugRow, error) {
	return m.author, m.err
}

func (m *mockStore) ListBooksByAuthor(_ context.Context, _ db.ListBooksByAuthorParams) ([]db.ListBooksByAuthorRow, error) {
	return m.authorBooks, m.err
}

func (m *mockStore) ListGenres(_ context.Context) ([]db.ListGenresRow, error) {
	return m.genres, m.err
}

func (m *mockStore) GetGenreBySlug(_ context.Context, _ string) (db.GetGenreBySlugRow, error) {
	return m.genre, m.err
}

func (m *mockStore) ListBooksByGenre(_ context.Context, _ db.ListBooksByGenreParams) ([]db.ListBooksByGenreRow, error) {
	return m.genreBooks, m.err
}

func (m *mockStore) ListYears(_ context.Context) ([]db.ListYearsRow, error) {
	return m.years, m.err
}

func newRouter(h *api.PublicHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/books/currently-reading", h.GetCurrentlyReading)
	r.Get("/api/books", h.GetBooks)
	r.Get("/api/books/{slug}", h.GetBookBySlug)
	r.Get("/api/authors", h.GetAuthors)
	r.Get("/api/authors/{slug}", h.GetAuthorBySlug)
	r.Get("/api/genres", h.GetGenres)
	r.Get("/api/genres/{slug}", h.GetGenreBySlug)
	r.Get("/api/years", h.GetYears)
	return r
}

func TestCursorRoundtrip(t *testing.T) {
	readAt := time.Date(2024, 11, 15, 10, 30, 0, 0, time.UTC)
	encoded := api.EncodeCursor(readAt, 1234)
	gotTime, gotID, err := api.DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("decodeCursor error: %v", err)
	}
	if !gotTime.Equal(readAt) {
		t.Errorf("time mismatch: got %v want %v", gotTime, readAt)
	}
	if gotID != 1234 {
		t.Errorf("id mismatch: got %d want 1234", gotID)
	}
}

func TestGetBooks(t *testing.T) {
	store := &mockStore{}
	h := &api.PublicHandlers{Store: store, SiteURL: "http://localhost:8081"}
	r := newRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("want Content-Type application/json, got %q", ct)
	}
	var resp api.PaginatedResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.HasMore {
		t.Error("want has_more=false for empty store")
	}
	if resp.NextCursor != nil {
		t.Error("want next_cursor=nil for empty store")
	}
}

func TestGetCurrentlyReading(t *testing.T) {
	readAt := pgtype.Timestamptz{Time: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	store := &mockStore{
		currentlyReading: []db.GetCurrentlyReadingRow{
			{Slug: "the-great-book", Title: "The Great Book", ReadAt: readAt, Authors: []interface{}{}, Genres: []interface{}{}},
		},
	}
	h := &api.PublicHandlers{Store: store, SiteURL: "http://localhost:8081"}
	r := newRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/books/currently-reading", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var items []api.BookListItem
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("want 1 item, got %d", len(items))
	}
	if items[0].Slug != "the-great-book" {
		t.Errorf("want slug=the-great-book, got %q", items[0].Slug)
	}
}

func TestGetBookBySlug(t *testing.T) {
	readAt := pgtype.Timestamptz{Time: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), Valid: true}
	store := &mockStore{
		detail: db.GetBookDetailBySlugRow{
			Slug: "test-book", Title: "Test Book", ReadAt: readAt,
			Authors: []interface{}{}, Genres: []interface{}{},
		},
	}
	h := &api.PublicHandlers{Store: store, SiteURL: "http://localhost:8081"}
	r := newRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/books/test-book", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var detail api.BookDetail
	if err := json.NewDecoder(w.Body).Decode(&detail); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if detail.Slug != "test-book" {
		t.Errorf("want slug=test-book, got %q", detail.Slug)
	}
}

func TestGetBookBySlugNotFound(t *testing.T) {
	store := &mockStore{err: pgx.ErrNoRows}
	h := &api.PublicHandlers{Store: store, SiteURL: "http://localhost:8081"}
	r := newRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/books/does-not-exist", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestGetAuthors(t *testing.T) {
	store := &mockStore{
		authors: []db.ListAuthorsRow{
			{Name: "Jane Doe", Slug: "jane-doe", BookCount: 3},
		},
	}
	h := &api.PublicHandlers{Store: store, SiteURL: "http://localhost:8081"}
	r := newRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/authors", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var items []api.AuthorListItem
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(items) != 1 || items[0].Slug != "jane-doe" {
		t.Errorf("unexpected authors: %+v", items)
	}
}

func TestGetAuthorBySlug(t *testing.T) {
	store := &mockStore{
		author:      db.GetAuthorBySlugRow{Name: "Jane Doe", Slug: "jane-doe"},
		authorBooks: []db.ListBooksByAuthorRow{},
	}
	h := &api.PublicHandlers{Store: store, SiteURL: "http://localhost:8081"}
	r := newRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/authors/jane-doe", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var detail api.AuthorDetail
	if err := json.NewDecoder(w.Body).Decode(&detail); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if detail.Slug != "jane-doe" {
		t.Errorf("want slug=jane-doe, got %q", detail.Slug)
	}
}

func TestGetGenres(t *testing.T) {
	store := &mockStore{
		genres: []db.ListGenresRow{
			{Name: "Fiction", Slug: "fiction", BookCount: 10},
		},
	}
	h := &api.PublicHandlers{Store: store, SiteURL: "http://localhost:8081"}
	r := newRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/genres", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var items []api.GenreListItem
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(items) != 1 || items[0].Slug != "fiction" {
		t.Errorf("unexpected genres: %+v", items)
	}
}

func TestGetGenreBySlug(t *testing.T) {
	store := &mockStore{
		genre:      db.GetGenreBySlugRow{Name: "Fiction", Slug: "fiction"},
		genreBooks: []db.ListBooksByGenreRow{},
	}
	h := &api.PublicHandlers{Store: store, SiteURL: "http://localhost:8081"}
	r := newRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/genres/fiction", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var detail api.GenreDetail
	if err := json.NewDecoder(w.Body).Decode(&detail); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if detail.Slug != "fiction" {
		t.Errorf("want slug=fiction, got %q", detail.Slug)
	}
}

func TestGetYears(t *testing.T) {
	store := &mockStore{
		years: []db.ListYearsRow{
			{Year: 2024, BookCount: 42},
			{Year: 2023, BookCount: 38},
		},
	}
	h := &api.PublicHandlers{Store: store, SiteURL: "http://localhost:8081"}
	r := newRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/years", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var items []api.YearCount
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("want 2 years, got %d", len(items))
	}
	if items[0].Year != 2024 || items[0].BookCount != 42 {
		t.Errorf("unexpected year data: %+v", items[0])
	}
}
