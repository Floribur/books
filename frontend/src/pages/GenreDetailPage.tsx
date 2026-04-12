import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchGenres, fetchBooks } from '../api/books';
import { BookGrid } from '../components/BookGrid';
import { Toast } from '../components/Toast';
import { useState, useEffect } from 'react';
import { usePageTitle } from '../hooks/usePageTitle';
import './GenreDetailPage.css';

export function GenreDetailPage() {
  const { slug } = useParams<{ slug: string }>();
  const [showError, setShowError] = useState(false);

  const { data: allGenres, isError: genresError } = useQuery({
    queryKey: ['genres'],
    queryFn: fetchGenres,
  });
  const { data: allBooks, isPending: booksPending, isError: booksError } = useQuery({
    queryKey: ['books'],
    queryFn: fetchBooks,
  });

  const isError = genresError || booksError;
  useEffect(() => {
    if (isError) setShowError(true);
  }, [isError]);

  const genre = allGenres?.find((g) => g.slug === slug);
  usePageTitle(genre?.name);

  const genreBooks = allBooks?.filter((b) =>
    b.genres.some((g) => g.slug === slug)
  ) ?? [];

  return (
    <main className="genre-detail-page">
      {showError && (
        <Toast
          message="Couldn't load genre. Try refreshing the page."
          onDismiss={() => setShowError(false)}
        />
      )}

      {booksPending ? (
        <div className="genre-detail-heading-skeleton" aria-hidden="true" />
      ) : genre ? (
        <div className="genre-detail-heading">
          <h1 className="genre-detail-name">{genre.name}</h1>
        </div>
      ) : null}

      <BookGrid
        books={genreBooks}
        isPending={booksPending}
        isError={booksError}
        ariaLabel={`Books in ${genre?.name ?? 'this genre'}`}
      />
    </main>
  );
}
