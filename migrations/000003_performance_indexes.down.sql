BEGIN;

DROP INDEX IF EXISTS idx_book_authors_author_id;
DROP INDEX IF EXISTS idx_book_genres_genre_id;
DROP INDEX IF EXISTS idx_books_shelf_read_at_id;
COMMIT;
