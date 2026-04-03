import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchGenreBySlug, fetchBooksByGenre } from '../api/books';
import { BookGrid } from '../components/BookGrid';
import { Toast } from '../components/Toast';
import { useState, useEffect } from 'react';
import './GenreDetailPage.css';

export function GenreDetailPage() {
  const { slug } = useParams<{ slug: string }>();
  const [showError, setShowError] = useState(false);

  const { data: genre, isPending, isError } = useQuery({
    queryKey: ['genre', slug],
    queryFn: () => fetchGenreBySlug(slug!),
    enabled: !!slug,
  });

  useEffect(() => {
    if (isError) setShowError(true);
  }, [isError]);

  const fetchFn = (cursor: string | undefined) => fetchBooksByGenre(slug!, cursor);

  return (
    <main className="genre-detail-page">
      {showError && (
        <Toast
          message="Couldn't load genre. Try refreshing the page."
          onDismiss={() => setShowError(false)}
        />
      )}

      {isPending ? (
        <div className="genre-detail-heading-skeleton" aria-hidden="true" />
      ) : genre ? (
        <div className="genre-detail-heading">
          <h1 className="genre-detail-name">{genre.name}</h1>
          <p className="genre-detail-count">
            {genre.book_count === 1 ? '1 book' : `${genre.book_count} books`}
          </p>
        </div>
      ) : null}

      {slug && (
        <BookGrid
          queryKey={['books', 'genre', slug]}
          fetchFn={fetchFn}
          ariaLabel={`Books in ${genre?.name ?? 'this genre'}`}
        />
      )}
    </main>
  );
}
