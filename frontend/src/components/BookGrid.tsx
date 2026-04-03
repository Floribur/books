import { useCallback, useState, useEffect } from 'react';
import { useInfiniteQuery } from '@tanstack/react-query';
import { BookCard } from './BookCard';
import { SkeletonCard } from './SkeletonCard';
import { Toast } from './Toast';
import { useIntersectionObserver } from '../hooks/useIntersectionObserver';
import { useScrollRestoration } from '../hooks/useScrollRestoration';
import type { PaginatedBooks } from '../api/types';
import './BookGrid.css';

interface BookGridProps {
  queryKey: string[];
  fetchFn: (cursor: string | undefined) => Promise<PaginatedBooks>;
  ariaLabel?: string;
}

export function BookGrid({ queryKey, fetchFn, ariaLabel = 'Books' }: BookGridProps) {
  const [showError, setShowError] = useState(false);
  const { restoreScroll } = useScrollRestoration();

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isPending,   // v5: "no data yet" — use for skeleton display (NOT isLoading)
    isError,
  } = useInfiniteQuery<PaginatedBooks, Error, { pages: PaginatedBooks[] }, string[], string | undefined>({
    queryKey,
    queryFn: ({ pageParam }) => fetchFn(pageParam),
    initialPageParam: undefined,  // REQUIRED in TanStack Query v5
    getNextPageParam: (lastPage) =>
      lastPage.has_more ? lastPage.next_cursor ?? undefined : undefined,
  });

  useEffect(() => {
    if (isError) setShowError(true);
  }, [isError]);

  const books = data?.pages.flatMap((page) => page.items) ?? [];
  useEffect(() => {
    if (books.length > 0) restoreScroll();
  }, [books.length, restoreScroll]);

  const handleIntersect = useCallback(() => {
    if (hasNextPage && !isFetchingNextPage) fetchNextPage();
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  const sentinelRef = useIntersectionObserver(handleIntersect, !!hasNextPage);

  return (
    <section className="book-grid-section" aria-label={ariaLabel}>
      {showError && (
        <Toast
          message="Couldn't load books. Try refreshing the page."
          onDismiss={() => setShowError(false)}
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
          {books.map((book) => (
            <BookCard key={book.slug} book={book} />
          ))}
        </div>
      )}

      <div ref={sentinelRef} style={{ height: 1 }} />

      {hasNextPage && (
        <button
          className="load-more-button"
          onClick={() => fetchNextPage()}
          disabled={isFetchingNextPage}
          aria-label={isFetchingNextPage ? 'Loading more books' : 'Load more books'}
        >
          {isFetchingNextPage ? 'Loading…' : 'Load More Books'}
        </button>
      )}

      {!isPending && !hasNextPage && (data?.pages.length ?? 0) > 1 && (
        <p className="books-end-message">You&apos;ve reached the end.</p>
      )}
    </section>
  );
}
