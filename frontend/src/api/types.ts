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
