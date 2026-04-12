// Static JSON response types — matches generated JSON from cmd/generate CLI
// All data is loaded once from /static/*.json; filtering and pagination are client-side.

export interface Author {
  name: string;
  slug: string;
}

export interface Genre {
  name: string;
  slug: string;
}

// Book — one item from /static/books.json
export interface Book {
  slug: string;
  title: string;
  cover_path: string;       // jsDelivr CDN URL or empty string
  read_at: string | null;   // ISO-8601 string or null
  publication_year: number | null;
  page_count: number | null;
  shelf: string;            // 'read' | 'currently-reading' etc.
  authors: Author[];
  genres: Genre[];
}

// BookDetail — shape of /static/books/{slug}.json
export interface BookDetail extends Book {
  description: string | null;
  isbn13: string | null;
  read_count: number;
  metadata_source: string;
}

// AuthorWithCount — one item from /static/authors.json
export interface AuthorWithCount extends Author {
  book_count: number;
}

// GenreWithCount — one item from /static/genres.json
export interface GenreWithCount extends Genre {
  book_count: number;
}
