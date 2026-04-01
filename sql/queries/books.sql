-- name: GetBookBySlug :one
SELECT * FROM books WHERE slug = $1 LIMIT 1;

-- name: ListBooks :many
SELECT * FROM books ORDER BY read_at DESC NULLS LAST, id DESC;
