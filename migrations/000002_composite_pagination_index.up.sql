BEGIN;
CREATE INDEX IF NOT EXISTS idx_books_read_at_id
    ON books(read_at DESC NULLS LAST, id DESC);
COMMIT;
