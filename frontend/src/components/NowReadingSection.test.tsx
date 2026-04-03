import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it } from 'vitest';
import { http, HttpResponse } from 'msw';
import { server } from '../test/msw-server';
import { NowReadingSection } from './NowReadingSection';

const testBook = {
  slug: 'test-book',
  title: 'Now Reading Test Book',
  cover_path: '/covers/test.jpg',
  read_at: null,
  publication_year: 2024,
  authors: [{ name: 'Test Author', slug: 'test-author' }],
  genres: [],
};

function renderSection() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <NowReadingSection />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('NowReadingSection', () => {
  it('renders nothing (null) when API returns empty array', async () => {
    server.use(
      http.get('/api/books/currently-reading', () => HttpResponse.json([]))
    );
    const { container } = renderSection();
    // Wait for query to resolve
    await new Promise((r) => setTimeout(r, 50));
    expect(container).toBeEmptyDOMElement();
  });

  it('renders cards for currently reading books', async () => {
    server.use(
      http.get('/api/books/currently-reading', () =>
        HttpResponse.json([testBook])
      )
    );
    renderSection();
    expect(await screen.findByText('Now Reading Test Book')).toBeInTheDocument();
  });

  it('shows at most 4 books regardless of API response length', async () => {
    const fiveBooks = Array.from({ length: 5 }, (_, i) => ({
      ...testBook,
      slug: `book-${i}`,
      title: `Book ${i}`,
    }));
    server.use(
      http.get('/api/books/currently-reading', () =>
        HttpResponse.json(fiveBooks)
      )
    );
    renderSection();
    // Wait for all books to render
    await screen.findByText('Book 0');
    // Should have exactly 4 NowReadingCard article elements
    const cards = document.querySelectorAll('article.now-reading-card');
    expect(cards.length).toBe(4);
  });

  it('shows error toast when API returns 500', async () => {
    server.use(
      http.get('/api/books/currently-reading', () =>
        HttpResponse.error()
      )
    );
    renderSection();
    expect(
      await screen.findByText("Couldn't load currently reading. Try refreshing.")
    ).toBeInTheDocument();
  });
});
