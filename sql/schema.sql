CREATE TABLE IF NOT EXISTS authors (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT        NOT NULL,
    slug        TEXT        NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS genres (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT        NOT NULL,
    slug        TEXT        NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS books (
    id                  BIGSERIAL    PRIMARY KEY,
    goodreads_id        TEXT         NOT NULL UNIQUE,
    slug                TEXT         NOT NULL UNIQUE,
    title               TEXT         NOT NULL,
    description         TEXT,
    cover_path          TEXT,
    page_count          INT,
    publication_year    INT,
    isbn13              TEXT,
    metadata_source     TEXT         NOT NULL DEFAULT 'none',
    read_at             TIMESTAMPTZ,
    date_added          TIMESTAMPTZ,
    read_count          INT          NOT NULL DEFAULT 1,
    shelf               TEXT         NOT NULL DEFAULT 'read',
    search_vector       TSVECTOR,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_books_slug         ON books(slug);
CREATE INDEX IF NOT EXISTS idx_books_read_at      ON books(read_at DESC NULLS LAST);
CREATE INDEX IF NOT EXISTS idx_books_goodreads_id ON books(goodreads_id);
CREATE INDEX IF NOT EXISTS idx_books_shelf        ON books(shelf);

CREATE TABLE IF NOT EXISTS book_authors (
    book_id     BIGINT NOT NULL REFERENCES books(id)   ON DELETE CASCADE,
    author_id   BIGINT NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

CREATE TABLE IF NOT EXISTS book_genres (
    book_id     BIGINT NOT NULL REFERENCES books(id)  ON DELETE CASCADE,
    genre_id    BIGINT NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, genre_id)
);
