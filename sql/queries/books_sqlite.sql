-- name: UpsertBook :one
INSERT INTO books (
    goodreads_id, slug, title, description, cover_path, page_count,
    publication_year, isbn13, metadata_source, read_at, date_added,
    read_count, shelf, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now')
)
ON CONFLICT (goodreads_id) DO UPDATE SET
    slug       = excluded.slug,
    title      = excluded.title,
    read_at    = excluded.read_at,
    date_added = excluded.date_added,
    read_count = excluded.read_count,
    shelf      = excluded.shelf,
    updated_at = datetime('now')
RETURNING id, goodreads_id, slug, title, description, cover_path, page_count,
          publication_year, isbn13, metadata_source, read_at, date_added,
          read_count, shelf, created_at, updated_at;

-- name: GetAllGoodreadsIDs :many
SELECT goodreads_id, slug, shelf FROM books;

-- name: GetUnenrichedBooks :many
SELECT id, goodreads_id, slug, title, description, cover_path, page_count,
       publication_year, isbn13, metadata_source, read_at, date_added,
       read_count, shelf, created_at, updated_at
FROM books WHERE metadata_source = 'none' ORDER BY created_at ASC;

-- name: UpdateBookEnrichment :exec
UPDATE books SET
    description      = ?,
    page_count       = ?,
    publication_year = ?,
    cover_path       = ?,
    metadata_source  = ?,
    updated_at       = datetime('now')
WHERE id = ?;

-- name: ListAllBooks :many
SELECT
  b.id, b.slug, b.title, b.cover_path, b.read_at, b.publication_year,
  b.page_count, b.description, b.isbn13, b.read_count, b.shelf, b.metadata_source,
  (SELECT json_group_array(json_object('name', a.name, 'slug', a.slug))
   FROM book_authors ba JOIN authors a ON a.id = ba.author_id
   WHERE ba.book_id = b.id) AS authors_json,
  (SELECT json_group_array(json_object('name', g.name, 'slug', g.slug))
   FROM book_genres bg JOIN genres g ON g.id = bg.genre_id
   WHERE bg.book_id = b.id) AS genres_json
FROM books b
ORDER BY b.read_at DESC, b.id DESC;
