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
