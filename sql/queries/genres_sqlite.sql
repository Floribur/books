-- name: UpsertGenre :one
INSERT INTO genres (name, slug, created_at, updated_at)
VALUES (?, ?, datetime('now'), datetime('now'))
ON CONFLICT (slug) DO UPDATE SET
    name       = excluded.name,
    updated_at = datetime('now')
RETURNING id, name, slug, created_at, updated_at;

-- name: LinkBookGenre :exec
INSERT INTO book_genres (book_id, genre_id) VALUES (?, ?) ON CONFLICT DO NOTHING;

-- name: ListGenres :many
SELECT g.id, g.name, g.slug, COUNT(bg.book_id) AS book_count
FROM genres g
LEFT JOIN book_genres bg ON bg.genre_id = g.id
GROUP BY g.id
ORDER BY book_count DESC, g.name ASC;
