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

-- name: UpdateBookEnrichment :exec
UPDATE books SET
    description      = $2,
    page_count       = $3,
    publication_year = $4,
    cover_path       = $5,
    metadata_source  = $6,
    updated_at       = NOW()
WHERE id = $1;

-- name: ListBooksPaginated :many
SELECT
    b.id, b.slug, b.title, b.cover_path, b.read_at, b.publication_year,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', a.name, 'slug', a.slug))
        FILTER (WHERE a.id IS NOT NULL), '[]'
    ) AS authors,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', g.name, 'slug', g.slug))
        FILTER (WHERE g.id IS NOT NULL), '[]'
    ) AS genres
FROM books b
LEFT JOIN book_authors ba ON ba.book_id = b.id
LEFT JOIN authors a ON a.id = ba.author_id
LEFT JOIN book_genres bg ON bg.book_id = b.id
LEFT JOIN genres g ON g.id = bg.genre_id
WHERE b.shelf = 'read'
  AND ($1::timestamptz IS NULL OR (b.read_at, b.id) < ($1::timestamptz, $2::bigint))
GROUP BY b.id
ORDER BY b.read_at DESC NULLS LAST, b.id DESC
LIMIT $3;

-- name: ListBooksByYear :many
SELECT
    b.id, b.slug, b.title, b.cover_path, b.read_at, b.publication_year,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', a.name, 'slug', a.slug))
        FILTER (WHERE a.id IS NOT NULL), '[]'
    ) AS authors,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', g.name, 'slug', g.slug))
        FILTER (WHERE g.id IS NOT NULL), '[]'
    ) AS genres
FROM books b
LEFT JOIN book_authors ba ON ba.book_id = b.id
LEFT JOIN authors a ON a.id = ba.author_id
LEFT JOIN book_genres bg ON bg.book_id = b.id
LEFT JOIN genres g ON g.id = bg.genre_id
WHERE b.shelf = 'read'
  AND EXTRACT(YEAR FROM b.read_at) = $1::int
  AND ($2::timestamptz IS NULL OR (b.read_at, b.id) < ($2::timestamptz, $3::bigint))
GROUP BY b.id
ORDER BY b.read_at DESC NULLS LAST, b.id DESC
LIMIT $4;

-- name: GetCurrentlyReading :many
SELECT
    b.id, b.slug, b.title, b.cover_path, b.read_at, b.publication_year,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', a.name, 'slug', a.slug))
        FILTER (WHERE a.id IS NOT NULL), '[]'
    ) AS authors,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', g.name, 'slug', g.slug))
        FILTER (WHERE g.id IS NOT NULL), '[]'
    ) AS genres
FROM books b
LEFT JOIN book_authors ba ON ba.book_id = b.id
LEFT JOIN authors a ON a.id = ba.author_id
LEFT JOIN book_genres bg ON bg.book_id = b.id
LEFT JOIN genres g ON g.id = bg.genre_id
WHERE b.shelf = 'currently-reading'
GROUP BY b.id
ORDER BY b.date_added DESC NULLS LAST;

-- name: GetBookDetailBySlug :one
SELECT
    b.id, b.slug, b.title, b.cover_path, b.read_at, b.publication_year,
    b.description, b.page_count, b.isbn13, b.read_count, b.shelf, b.metadata_source,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', a.name, 'slug', a.slug))
        FILTER (WHERE a.id IS NOT NULL), '[]'
    ) AS authors,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', g.name, 'slug', g.slug))
        FILTER (WHERE g.id IS NOT NULL), '[]'
    ) AS genres
FROM books b
LEFT JOIN book_authors ba ON ba.book_id = b.id
LEFT JOIN authors a ON a.id = ba.author_id
LEFT JOIN book_genres bg ON bg.book_id = b.id
LEFT JOIN genres g ON g.id = bg.genre_id
WHERE b.slug = $1
GROUP BY b.id;
