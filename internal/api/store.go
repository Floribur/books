package api

import (
	"context"

	"flos-library/internal/db"
)

// BookStore is the subset of *db.Queries used by PublicHandlers.
// A mock implementation is injected during tests.
type BookStore interface {
	ListBooksPaginated(ctx context.Context, arg db.ListBooksPaginatedParams) ([]db.ListBooksPaginatedRow, error)
	GetCurrentlyReading(ctx context.Context) ([]db.GetCurrentlyReadingRow, error)
	GetBookDetailBySlug(ctx context.Context, slug string) (db.GetBookDetailBySlugRow, error)
	ListAuthors(ctx context.Context) ([]db.ListAuthorsRow, error)
	GetAuthorBySlug(ctx context.Context, slug string) (db.GetAuthorBySlugRow, error)
	ListBooksByAuthor(ctx context.Context, arg db.ListBooksByAuthorParams) ([]db.ListBooksByAuthorRow, error)
	ListGenres(ctx context.Context) ([]db.ListGenresRow, error)
	GetGenreBySlug(ctx context.Context, slug string) (db.GetGenreBySlugRow, error)
	ListBooksByGenre(ctx context.Context, arg db.ListBooksByGenreParams) ([]db.ListBooksByGenreRow, error)
	ListYears(ctx context.Context) ([]db.ListYearsRow, error)
}

// Verify *db.Queries satisfies BookStore at compile time.
var _ BookStore = (*db.Queries)(nil)
