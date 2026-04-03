import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it } from 'vitest';
import { http, HttpResponse } from 'msw';
import { server } from '../test/msw-server';
import { BookGrid } from './BookGrid';

function renderGrid() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <BookGrid />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

const testBook = {
  slug: 'test-book',
  title: 'Test Book',
  cover_path: '/covers/test.jpg',
  read_at: '2024-01-15T00:00:00Z',
  publication_year: 2024,
  authors: [{ name: 'Test Author', slug: 'test-author' }],
  genres: [],
};

describe('BookGrid', () => {
  it('shows 12 skeleton cards while loading (isPending)', async () => {
    server.use(
      http.get('/api/books', () => new Promise(() => {})) // pending forever
    );
    renderGrid();
    const skeletons = document.querySelectorAll('[aria-hidden="true"]');
    expect(skeletons.length).toBe(12);
  });

  it('renders book cards after data loads', async () => {
    server.use(
      http.get('/api/books', () =>
        HttpResponse.json({ items: [testBook], next_cursor: null, has_more: false })
      )
    );
    renderGrid();
    expect(await screen.findByText('Test Book')).toBeInTheDocument();
  });

  it('shows Load More button when hasNextPage is true', async () => {
    server.use(
      http.get('/api/books', () =>
        HttpResponse.json({
          items: [testBook],
          next_cursor: 'abc123',
          has_more: true,
        })
      )
    );
    renderGrid();
    expect(await screen.findByText('Load More Books')).toBeInTheDocument();
  });

  it("shows end message and no Load More button when hasNextPage is false", async () => {
    server.use(
      http.get('/api/books', () =>
        HttpResponse.json({ items: [testBook], next_cursor: null, has_more: false })
      )
    );
    renderGrid();
    expect(await screen.findByText("You've reached the end.")).toBeInTheDocument();
    expect(screen.queryByText('Load More Books')).toBeNull();
  });

  it('shows error toast when API fails', async () => {
    server.use(
      http.get('/api/books', () => new HttpResponse(null, { status: 500 }))
    );
    renderGrid();
    expect(
      await screen.findByText("Couldn't load books. Try refreshing the page.")
    ).toBeInTheDocument();
  });
});
