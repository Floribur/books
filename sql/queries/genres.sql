-- name: UpsertGenre :one
INSERT INTO genres (name, slug, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
ON CONFLICT (slug) DO UPDATE SET
    name = EXCLUDED.name,
    updated_at = NOW()
RETURNING *;

-- name: LinkBookGenre :exec
INSERT INTO book_genres (book_id, genre_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: ListGenres :many
SELECT g.id, g.name, g.slug, COUNT(bg.book_id)::bigint AS book_count
FROM genres g
LEFT JOIN book_genres bg ON bg.genre_id = g.id
GROUP BY g.id
ORDER BY COUNT(bg.book_id) DESC, g.name ASC;

-- name: GetGenreBySlug :one
SELECT g.id, g.name, g.slug
FROM genres g
WHERE g.slug = $1;

-- name: ListBooksByGenre :many
SELECT
    b.id, b.slug, b.title, b.cover_path, b.read_at, b.publication_year,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', a.name, 'slug', a.slug))
        FILTER (WHERE a.id IS NOT NULL), '[]'
    ) AS authors,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object('name', g2.name, 'slug', g2.slug))
        FILTER (WHERE g2.id IS NOT NULL), '[]'
    ) AS genres
FROM books b
JOIN book_genres bg ON bg.book_id = b.id
JOIN genres g ON g.id = bg.genre_id
LEFT JOIN book_authors ba ON ba.book_id = b.id
LEFT JOIN authors a ON a.id = ba.author_id
LEFT JOIN book_genres bg2 ON bg2.book_id = b.id
LEFT JOIN genres g2 ON g2.id = bg2.genre_id
WHERE g.slug = $1
  AND ($2::timestamptz IS NULL OR (b.read_at, b.id) < ($2::timestamptz, $3::bigint))
GROUP BY b.id
ORDER BY b.read_at DESC NULLS LAST, b.id DESC
LIMIT $4;

-- name: ListYears :many
SELECT
    EXTRACT(YEAR FROM read_at)::int AS year,
    COUNT(*)::bigint AS book_count
FROM books
WHERE read_at IS NOT NULL AND shelf = 'read'
GROUP BY EXTRACT(YEAR FROM read_at)
ORDER BY year DESC;
