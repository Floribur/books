-- name: UpsertAuthor :one
INSERT INTO authors (name, slug, created_at, updated_at)
VALUES (?, ?, datetime('now'), datetime('now'))
ON CONFLICT (slug) DO UPDATE SET
    name       = excluded.name,
    updated_at = datetime('now')
RETURNING id, name, slug, created_at, updated_at;

-- name: LinkBookAuthor :exec
INSERT INTO book_authors (book_id, author_id) VALUES (?, ?) ON CONFLICT DO NOTHING;

-- name: GetFirstAuthorForBook :one
SELECT a.name FROM authors a
JOIN book_authors ba ON ba.author_id = a.id
WHERE ba.book_id = ? LIMIT 1;

-- name: ListAuthors :many
SELECT a.id, a.name, a.slug, COUNT(ba.book_id) AS book_count
FROM authors a
LEFT JOIN book_authors ba ON ba.author_id = a.id
GROUP BY a.id
ORDER BY a.name ASC;
