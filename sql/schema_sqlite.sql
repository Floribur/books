CREATE TABLE IF NOT EXISTS authors (
    id          INTEGER PRIMARY KEY,
    name        TEXT    NOT NULL,
    slug        TEXT    NOT NULL UNIQUE,
    created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS genres (
    id          INTEGER PRIMARY KEY,
    name        TEXT    NOT NULL,
    slug        TEXT    NOT NULL UNIQUE,
    created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS books (
    id                  INTEGER PRIMARY KEY,
    goodreads_id        TEXT    NOT NULL UNIQUE,
    slug                TEXT    NOT NULL UNIQUE,
    title               TEXT    NOT NULL,
    description         TEXT,
    cover_path          TEXT,
    page_count          INTEGER,
    publication_year    INTEGER,
    isbn13              TEXT,
    metadata_source     TEXT    NOT NULL DEFAULT 'none',
    read_at             TEXT,
    date_added          TEXT,
    read_count          INTEGER NOT NULL DEFAULT 1,
    shelf               TEXT    NOT NULL DEFAULT 'read',
    created_at          TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at          TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_books_slug         ON books(slug);
CREATE INDEX IF NOT EXISTS idx_books_read_at      ON books(read_at DESC);
CREATE INDEX IF NOT EXISTS idx_books_goodreads_id ON books(goodreads_id);
CREATE INDEX IF NOT EXISTS idx_books_shelf        ON books(shelf);

CREATE TABLE IF NOT EXISTS book_authors (
    book_id     INTEGER NOT NULL REFERENCES books(id)   ON DELETE CASCADE,
    author_id   INTEGER NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

CREATE TABLE IF NOT EXISTS book_genres (
    book_id     INTEGER NOT NULL REFERENCES books(id)  ON DELETE CASCADE,
    genre_id    INTEGER NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, genre_id)
);
