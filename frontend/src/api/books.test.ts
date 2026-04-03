// TDD RED: Task 1 — API layer extension + BookGrid generic refactor
// These tests verify the new fetch functions and types added in Phase 4 Plan 1

import { describe, expect, it } from 'vitest';
import { http, HttpResponse } from 'msw';
import { server } from '../test/msw-server';
import {
  fetchBookBySlug,
  fetchAuthors,
  fetchAuthorBySlug,
  fetchBooksByAuthor,
  fetchGenres,
  fetchGenreBySlug,
  fetchBooksByGenre,
  fetchYears,
  fetchBooksByYear,
} from './books';

const testBook = {
  slug: 'test-book',
  title: 'Test Book',
  cover_path: '/covers/test.jpg',
  read_at: '2024-01-15T00:00:00Z',
  publication_year: 2024,
  authors: [{ name: 'Test Author', slug: 'test-author' }],
  genres: [{ name: 'Fiction', slug: 'fiction' }],
};

describe('fetchBookBySlug', () => {
  it('calls GET /api/books/:slug and returns BookDetail', async () => {
    const detail = {
      ...testBook,
      description: 'A great book.',
      page_count: 300,
      isbn13: '9780000000000',
      read_count: 1,
      shelf: 'read',
      metadata_source: 'google_books',
    };
    server.use(
      http.get('/api/books/dune', () => HttpResponse.json(detail))
    );
    const result = await fetchBookBySlug('dune');
    expect(result.slug).toBe('test-book');
    expect(result.description).toBe('A great book.');
    expect(result.page_count).toBe(300);
  });
});

describe('fetchAuthors', () => {
  it('calls GET /api/authors and returns AuthorWithCount[]', async () => {
    const authors = [
      { name: 'Frank Herbert', slug: 'frank-herbert', book_count: 5 },
    ];
    server.use(
      http.get('/api/authors', () => HttpResponse.json(authors))
    );
    const result = await fetchAuthors();
    expect(result).toHaveLength(1);
    expect(result[0].book_count).toBe(5);
  });
});

describe('fetchAuthorBySlug', () => {
  it('calls GET /api/authors/:slug and returns AuthorDetail', async () => {
    const detail = {
      name: 'Frank Herbert',
      slug: 'frank-herbert',
      book_count: 5,
      items: [testBook],
      next_cursor: null,
      has_more: false,
    };
    server.use(
      http.get('/api/authors/frank-herbert', () => HttpResponse.json(detail))
    );
    const result = await fetchAuthorBySlug('frank-herbert');
    expect(result.name).toBe('Frank Herbert');
    expect(result.book_count).toBe(5);
    expect(result.items).toHaveLength(1);
  });
});

describe('fetchBooksByAuthor', () => {
  it('calls GET /api/authors/:slug without cursor when cursor is undefined', async () => {
    let capturedUrl = '';
    server.use(
      http.get('/api/authors/tolkien', ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({
          name: 'Tolkien',
          slug: 'tolkien',
          book_count: 3,
          items: [testBook],
          next_cursor: null,
          has_more: false,
        });
      })
    );
    await fetchBooksByAuthor('tolkien', undefined);
    expect(capturedUrl).not.toContain('cursor');
  });

  it('calls GET /api/authors/:slug?cursor=... when cursor provided', async () => {
    let capturedUrl = '';
    server.use(
      http.get('/api/authors/tolkien', ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({
          name: 'Tolkien',
          slug: 'tolkien',
          book_count: 3,
          items: [testBook],
          next_cursor: null,
          has_more: false,
        });
      })
    );
    await fetchBooksByAuthor('tolkien', 'cursor-abc');
    expect(capturedUrl).toContain('cursor=cursor-abc');
  });

  it('returns PaginatedBooks envelope extracted from AuthorDetail', async () => {
    server.use(
      http.get('/api/authors/tolkien', () =>
        HttpResponse.json({
          name: 'Tolkien',
          slug: 'tolkien',
          book_count: 3,
          items: [testBook],
          next_cursor: 'next123',
          has_more: true,
        })
      )
    );
    const result = await fetchBooksByAuthor('tolkien', undefined);
    expect(result.items).toHaveLength(1);
    expect(result.next_cursor).toBe('next123');
    expect(result.has_more).toBe(true);
  });
});

describe('fetchGenres', () => {
  it('calls GET /api/genres and returns GenreWithCount[]', async () => {
    server.use(
      http.get('/api/genres', () =>
        HttpResponse.json([{ name: 'Fantasy', slug: 'fantasy', book_count: 10 }])
      )
    );
    const result = await fetchGenres();
    expect(result[0].book_count).toBe(10);
  });
});

describe('fetchGenreBySlug', () => {
  it('calls GET /api/genres/:slug and returns GenreDetail', async () => {
    server.use(
      http.get('/api/genres/fantasy', () =>
        HttpResponse.json({
          name: 'Fantasy',
          slug: 'fantasy',
          book_count: 10,
          items: [testBook],
          next_cursor: null,
          has_more: false,
        })
      )
    );
    const result = await fetchGenreBySlug('fantasy');
    expect(result.name).toBe('Fantasy');
    expect(result.items).toHaveLength(1);
  });
});

describe('fetchBooksByGenre', () => {
  it('calls GET /api/genres/:slug without cursor when undefined', async () => {
    let capturedUrl = '';
    server.use(
      http.get('/api/genres/fantasy', ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({
          name: 'Fantasy',
          slug: 'fantasy',
          book_count: 10,
          items: [testBook],
          next_cursor: null,
          has_more: false,
        });
      })
    );
    await fetchBooksByGenre('fantasy', undefined);
    expect(capturedUrl).not.toContain('cursor');
  });
});

describe('fetchYears', () => {
  it('calls GET /api/years and returns YearEntry[]', async () => {
    server.use(
      http.get('/api/years', () =>
        HttpResponse.json([{ year: 2024, book_count: 25 }])
      )
    );
    const result = await fetchYears();
    expect(result[0].year).toBe(2024);
    expect(result[0].book_count).toBe(25);
  });
});

describe('fetchBooksByYear', () => {
  it('calls GET /api/books?year=YYYY without cursor when undefined', async () => {
    let capturedUrl = '';
    server.use(
      http.get('/api/books', ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({ items: [testBook], next_cursor: null, has_more: false });
      })
    );
    await fetchBooksByYear(2024, undefined);
    expect(capturedUrl).toContain('year=2024');
    expect(capturedUrl).not.toContain('cursor');
  });

  it('calls GET /api/books?year=YYYY&cursor=... when cursor provided', async () => {
    let capturedUrl = '';
    server.use(
      http.get('/api/books', ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({ items: [testBook], next_cursor: null, has_more: false });
      })
    );
    await fetchBooksByYear(2024, 'cursor-abc');
    expect(capturedUrl).toContain('year=2024');
    expect(capturedUrl).toContain('cursor=cursor-abc');
  });
});
