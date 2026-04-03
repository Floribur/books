import { http, HttpResponse } from 'msw';

// Minimal test book fixture — enough to render a BookCard
const testBook = {
  slug: 'test-book-slug',
  title: 'Test Book Title',
  cover_path: '/covers/test.jpg',
  read_at: '2024-01-15T00:00:00Z',
  publication_year: 2024,
  authors: [{ name: 'Test Author', slug: 'test-author' }],
  genres: [{ name: 'Fiction', slug: 'fiction' }],
};

export const handlers = [
  http.get('/api/books', () => {
    return HttpResponse.json({
      items: [testBook],
      next_cursor: null,
      has_more: false,
    });
  }),

  http.get('/api/books/currently-reading', () => {
    return HttpResponse.json([testBook]);
  }),
];
