BEGIN;

-- Join table lookups by author_id (used by ListBooksByAuthor and ListAuthors COUNT)
CREATE INDEX IF NOT EXISTS idx_book_authors_author_id
    ON book_authors(author_id);

-- Join table lookups by genre_id (used by ListBooksByGenre and ListGenres COUNT)
CREATE INDEX IF NOT EXISTS idx_book_genres_genre_id
    ON book_genres(genre_id);

-- Main book list: shelf filter + pagination sort (replaces single-column read_at index)
CREATE INDEX IF NOT EXISTS idx_books_shelf_read_at_id
    ON books(shelf, read_at DESC NULLS LAST, id DESC);

COMMIT;
