import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it } from 'vitest';
import { BookGrid } from './BookGrid';
import type { Book } from '../api/types';

function makeBook(overrides: Partial<Book> = {}): Book {
  return {
    slug: 'test-book',
    title: 'Test Book',
    cover_path: '/covers/test.jpg',
    read_at: '2024-01-15T00:00:00Z',
    publication_year: 2024,
    page_count: 300,
    shelf: 'read',
    authors: [{ name: 'Test Author', slug: 'test-author' }],
    genres: [],
    ...overrides,
  };
}

function renderGrid(props: Partial<Parameters<typeof BookGrid>[0]> = {}) {
  return render(
    <MemoryRouter>
      <BookGrid books={[]} ariaLabel="Books Read" {...props} />
    </MemoryRouter>
  );
}

describe('BookGrid', () => {
  it('shows 12 skeleton cards while loading (isPending)', () => {
    renderGrid({ isPending: true });
    const skeletons = document.querySelectorAll('[aria-hidden="true"]');
    expect(skeletons.length).toBe(12);
  });

  it('renders book cards after data loads', () => {
    renderGrid({ books: [makeBook()] });
    expect(screen.getByText('Test Book')).toBeInTheDocument();
  });

  it('shows Load More button when there are more than 24 books', () => {
    const books = Array.from({ length: 25 }, (_, i) =>
      makeBook({ slug: `book-${i}`, title: `Book ${i}` })
    );
    renderGrid({ books });
    expect(screen.getByText('Load More Books')).toBeInTheDocument();
  });

  it('does not show Load More button when all books fit on one page', () => {
    renderGrid({ books: [makeBook()] });
    expect(screen.queryByText('Load More Books')).toBeNull();
  });

  it('shows end message after loading all books beyond one page', () => {
    const books = Array.from({ length: 25 }, (_, i) =>
      makeBook({ slug: `book-${i}`, title: `Book ${i}` })
    );
    renderGrid({ books });
    fireEvent.click(screen.getByText('Load More Books'));
    expect(screen.getByText("You've reached the end.")).toBeInTheDocument();
  });

  it('shows error toast when isError is true', () => {
    renderGrid({ isError: true });
    expect(
      screen.getByText("Couldn't load books. Try refreshing the page.")
    ).toBeInTheDocument();
  });
});
