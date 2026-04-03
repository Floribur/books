import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it, beforeEach } from 'vitest';
import { BookCard } from './BookCard';
import type { Book } from '../api/types';

const testBook: Book = {
  slug: 'dune-1965',
  title: 'Dune',
  cover_path: '/covers/test.jpg',
  read_at: '2024-01-15T00:00:00Z',
  publication_year: 1965,
  authors: [{ name: 'Frank Herbert', slug: 'frank-herbert' }],
  genres: [{ name: 'Science Fiction', slug: 'science-fiction' }],
};

function renderCard(book = testBook) {
  return render(
    <MemoryRouter>
      <BookCard book={book} />
    </MemoryRouter>
  );
}

describe('BookCard', () => {
  beforeEach(() => {
    sessionStorage.clear();
  });

  it('renders a link to /books/:slug', () => {
    renderCard();
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/books/dune-1965');
  });

  it('renders the book title', () => {
    renderCard();
    expect(screen.getByText('Dune')).toBeInTheDocument();
  });

  it('renders the author name', () => {
    renderCard();
    expect(screen.getByText('Frank Herbert')).toBeInTheDocument();
  });

  it('renders multiple authors joined by comma', () => {
    const multiAuthorBook: Book = {
      ...testBook,
      authors: [
        { name: 'Author One', slug: 'author-one' },
        { name: 'Author Two', slug: 'author-two' },
      ],
    };
    renderCard(multiAuthorBook);
    expect(screen.getByText('Author One, Author Two')).toBeInTheDocument();
  });

  it('has id="book-{slug}" for scroll restoration', () => {
    renderCard();
    expect(document.getElementById('book-dune-1965')).toBeInTheDocument();
  });

  it('saves slug to sessionStorage on click', () => {
    renderCard();
    fireEvent.click(screen.getByRole('link'));
    expect(sessionStorage.getItem('scrollTarget')).toBe('dune-1965');
  });
});
