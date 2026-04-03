-- name: GetFirstAuthorForBook :one
SELECT a.name FROM authors a
JOIN book_authors ba ON ba.author_id = a.id
WHERE ba.book_id = $1
LIMIT 1;

-- name: UpsertAuthor :one
INSERT INTO authors (name, slug, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
ON CONFLICT (slug) DO UPDATE SET
    name = EXCLUDED.name,
    updated_at = NOW()
RETURNING *;

-- name: LinkBookAuthor :exec
INSERT INTO book_authors (book_id, author_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: ListAuthors :many
SELECT a.id, a.name, a.slug, COUNT(ba.book_id)::bigint AS book_count
FROM authors a
LEFT JOIN book_authors ba ON ba.author_id = a.id
GROUP BY a.id
ORDER BY a.name ASC;

-- name: GetAuthorBySlug :one
SELECT a.id, a.name, a.slug
FROM authors a
WHERE a.slug = $1;

-- name: ListBooksByAuthor :many
SELECT
    b.id, b.slug, b.title, b.cover_path, b.read_at, b.publication_year,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', a2.name, 'slug', a2.slug))
        FILTER (WHERE a2.id IS NOT NULL), '[]'
    ) AS authors,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', g.name, 'slug', g.slug))
        FILTER (WHERE g.id IS NOT NULL), '[]'
    ) AS genres
FROM books b
JOIN book_authors ba ON ba.book_id = b.id
JOIN authors a ON a.id = ba.author_id
LEFT JOIN book_authors ba2 ON ba2.book_id = b.id
LEFT JOIN authors a2 ON a2.id = ba2.author_id
LEFT JOIN book_genres bg ON bg.book_id = b.id
LEFT JOIN genres g ON g.id = bg.genre_id
WHERE a.slug = $1
  AND ($2::timestamptz IS NULL OR (b.read_at, b.id) < ($2::timestamptz, $3::bigint))
GROUP BY b.id
ORDER BY b.read_at DESC NULLS LAST, b.id DESC
LIMIT $4;
