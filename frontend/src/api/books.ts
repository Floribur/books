import { apiFetch } from './client';
import type { Book, PaginatedBooks, BookDetail, AuthorWithCount, AuthorDetail, GenreWithCount, GenreDetail, YearEntry } from './types';

// Fetches a page of books from GET /api/books
// cursor: opaque base64 token from previous page's next_cursor; undefined = first page
export async function fetchBooks(cursor: string | undefined): Promise<PaginatedBooks> {
  const url = cursor ? `/api/books?cursor=${encodeURIComponent(cursor)}` : '/api/books';
  return apiFetch<PaginatedBooks>(url);
}

// Fetches currently-reading books from GET /api/books/currently-reading
// Returns a plain Book[] (no envelope — Phase 2 D-03)
export async function fetchCurrentlyReading(): Promise<Book[]> {
  return apiFetch<Book[]>('/api/books/currently-reading');
}

// GET /api/books/:slug — full book detail (BOOK-01 through BOOK-04)
export async function fetchBookBySlug(slug: string): Promise<BookDetail> {
  return apiFetch<BookDetail>(`/api/books/${encodeURIComponent(slug)}`);
}

// GET /api/authors — all authors with book counts, alphabetical by surname (AUTH-01)
export async function fetchAuthors(): Promise<AuthorWithCount[]> {
  return apiFetch<AuthorWithCount[]>('/api/authors');
}

// GET /api/authors/:slug — author metadata + first page of books (AUTH-02)
// Use for heading (name, book_count); returns full AuthorDetail
export async function fetchAuthorBySlug(slug: string): Promise<AuthorDetail> {
  return apiFetch<AuthorDetail>(`/api/authors/${encodeURIComponent(slug)}`);
}

// GET /api/authors/:slug?cursor=... — paginated books for author (AUTH-02)
// Extracts PaginatedBooks envelope from AuthorDetail response (nested under .books)
export async function fetchBooksByAuthor(slug: string, cursor: string | undefined): Promise<PaginatedBooks> {
  const url = cursor
    ? `/api/authors/${encodeURIComponent(slug)}?cursor=${encodeURIComponent(cursor)}`
    : `/api/authors/${encodeURIComponent(slug)}`;
  const data = await apiFetch<AuthorDetail>(url);
  return data.books;
}

// GET /api/genres — all genres with book counts, sorted by book_count desc (GENR-01)
export async function fetchGenres(): Promise<GenreWithCount[]> {
  return apiFetch<GenreWithCount[]>('/api/genres');
}

// GET /api/genres/:slug — genre metadata + first page of books (GENR-02)
export async function fetchGenreBySlug(slug: string): Promise<GenreDetail> {
  return apiFetch<GenreDetail>(`/api/genres/${encodeURIComponent(slug)}`);
}

// GET /api/genres/:slug?cursor=... — paginated books for genre (GENR-02)
export async function fetchBooksByGenre(slug: string, cursor: string | undefined): Promise<PaginatedBooks> {
  const url = cursor
    ? `/api/genres/${encodeURIComponent(slug)}?cursor=${encodeURIComponent(cursor)}`
    : `/api/genres/${encodeURIComponent(slug)}`;
  const data = await apiFetch<GenreDetail>(url);
  return data.books;
}

// GET /api/years — distinct years with book counts (CHAL-02)
export async function fetchYears(): Promise<YearEntry[]> {
  return apiFetch<YearEntry[]>('/api/years');
}

// GET /api/books?year=YYYY&cursor=... — books filtered by year (CHAL-01)
export async function fetchBooksByYear(year: number, cursor: string | undefined): Promise<PaginatedBooks> {
  const base = `/api/books?year=${year}`;
  const url = cursor ? `${base}&cursor=${encodeURIComponent(cursor)}` : base;
  return apiFetch<PaginatedBooks>(url);
}
