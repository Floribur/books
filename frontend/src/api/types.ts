// API response types — matches Go API contract from Phase 2
// D-04: Book list shape (slug, title, cover_path, read_at, publication_year, authors[], genres[])
// D-01: Paginated envelope shape (items, next_cursor, has_more)

export interface Author {
  name: string;
  slug: string;
}

export interface Genre {
  name: string;
  slug: string;
}

export interface Book {
  slug: string;
  title: string;
  cover_path: string;       // e.g. "/covers/9780385490818.jpg" — proxied in dev via Vite
  read_at: string | null;   // ISO 8601 string or null
  publication_year: number | null;
  authors: Author[];
  genres: Genre[];
}

export interface PaginatedBooks {
  items: Book[];
  next_cursor: string | null;  // opaque base64 token; null = no more pages
  has_more: boolean;
}

// BookDetail — full detail shape from GET /api/books/:slug (Phase 2 D-05)
export interface BookDetail extends Book {
  description: string | null;
  page_count: number | null;
  isbn13: string | null;
  read_count: number;
  shelf: string;
  metadata_source: string;
}

// AuthorWithCount — author list item from GET /api/authors (Phase 2 D-06)
export interface AuthorWithCount extends Author {
  book_count: number;
}

// AuthorDetail — author detail from GET /api/authors/:slug (Phase 2 D-07)
// Extends AuthorWithCount; books are in the paginated envelope
export interface AuthorDetail extends AuthorWithCount {
  items: Book[];
  next_cursor: string | null;
  has_more: boolean;
}

// GenreWithCount — genre list item from GET /api/genres (Phase 2 D-06)
export interface GenreWithCount extends Genre {
  book_count: number;
}

// GenreDetail — genre detail from GET /api/genres/:slug (Phase 2 D-07)
export interface GenreDetail extends GenreWithCount {
  items: Book[];
  next_cursor: string | null;
  has_more: boolean;
}

// YearEntry — years list item from GET /api/years
export interface YearEntry {
  year: number;
  book_count: number;
}
