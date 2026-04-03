import { apiFetch } from './client';
import type { Book, PaginatedBooks } from './types';

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
