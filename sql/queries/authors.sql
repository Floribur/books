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
