-- name: GetBookBySlug :one
SELECT * FROM books WHERE slug = $1 LIMIT 1;

-- name: ListBooks :many
SELECT * FROM books ORDER BY read_at DESC NULLS LAST, id DESC;

-- name: UpsertBook :one
INSERT INTO books (
    goodreads_id, slug, title, description, cover_path, page_count,
    publication_year, isbn13, metadata_source, read_at, date_added,
    read_count, shelf, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW()
)
ON CONFLICT (goodreads_id) DO UPDATE SET
    slug             = EXCLUDED.slug,
    title            = EXCLUDED.title,
    read_at          = EXCLUDED.read_at,
    date_added       = EXCLUDED.date_added,
    read_count       = EXCLUDED.read_count,
    shelf            = EXCLUDED.shelf,
    updated_at       = NOW()
RETURNING *;

-- name: GetAllGoodreadsIDs :many
SELECT goodreads_id, slug, shelf FROM books;

-- name: GetUnenrichedBooks :many
SELECT * FROM books WHERE metadata_source = 'none' ORDER BY created_at ASC;
