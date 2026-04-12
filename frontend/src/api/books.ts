import { apiFetch } from './client';
import type { Book, BookDetail, AuthorWithCount, GenreWithCount } from './types';

// Fetches all books from /static/books.json (all shelves, sorted by read_at desc).
// Held in TanStack Query cache — one fetch per session.
export async function fetchBooks(): Promise<Book[]> {
  return apiFetch<Book[]>('/static/books.json');
}

// Fetches full detail for one book from /static/books/{slug}.json.
// Called on demand when navigating to a book detail page.
export async function fetchBookBySlug(slug: string): Promise<BookDetail> {
  return apiFetch<BookDetail>(`/static/books/${encodeURIComponent(slug)}.json`);
}

// Fetches all authors from /static/authors.json.
export async function fetchAuthors(): Promise<AuthorWithCount[]> {
  return apiFetch<AuthorWithCount[]>('/static/authors.json');
}

// Fetches all genres from /static/genres.json.
export async function fetchGenres(): Promise<GenreWithCount[]> {
  return apiFetch<GenreWithCount[]>('/static/genres.json');
}
