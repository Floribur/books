import { useCallback, useState } from 'react';
import { BookCard } from './BookCard';
import { SkeletonCard } from './SkeletonCard';
import { Toast } from './Toast';
import { useIntersectionObserver } from '../hooks/useIntersectionObserver';
import type { Book } from '../api/types';
import './BookGrid.css';

const PAGE_SIZE = 24;

interface BookGridProps {
  books: Book[];
  isPending?: boolean;
  isError?: boolean;
  ariaLabel?: string;
}

export function BookGrid({ books, isPending = false, isError = false, ariaLabel = 'Books' }: BookGridProps) {
  const [visibleCount, setVisibleCount] = useState(PAGE_SIZE);
  const hasMore = visibleCount < books.length;

  const handleIntersect = useCallback(() => {
    if (hasMore) setVisibleCount((c) => Math.min(c + PAGE_SIZE, books.length));
  }, [hasMore, books.length]);

  const sentinelRef = useIntersectionObserver(handleIntersect, hasMore);

  const visibleBooks = books.slice(0, visibleCount);

  return (
    <section className="book-grid-section" aria-label={ariaLabel}>
      {isError && (
        <Toast
          message="Couldn't load books. Try refreshing the page."
          onDismiss={() => {}}
        />
      )}

      {isPending ? (
        <div className="book-grid">
          {Array.from({ length: 12 }).map((_, i) => (
            <SkeletonCard key={i} />
          ))}
        </div>
      ) : (
        <div className="book-grid">
          {visibleBooks.map((book) => (
            <BookCard key={book.slug} book={book} />
          ))}
        </div>
      )}

      <div ref={sentinelRef} style={{ height: 1 }} />

      {hasMore && (
        <button
          className="load-more-button"
          onClick={() => setVisibleCount((c) => Math.min(c + PAGE_SIZE, books.length))}
          aria-label="Load more books"
        >
          Load More Books
        </button>
      )}

      {!isPending && !hasMore && books.length > PAGE_SIZE && (
        <p className="books-end-message">You&apos;ve reached the end.</p>
      )}
    </section>
  );
}
