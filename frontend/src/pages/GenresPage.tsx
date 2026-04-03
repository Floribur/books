import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchGenres } from '../api/books';
import { Toast } from '../components/Toast';
import { useState, useEffect } from 'react';
import './GenresPage.css';

export function GenresPage() {
  const [showError, setShowError] = useState(false);

  const { data: genres, isPending, isError } = useQuery({
    queryKey: ['genres'],
    queryFn: fetchGenres,
  });

  useEffect(() => {
    if (isError) setShowError(true);
  }, [isError]);

  // Compute max book_count for proportional bar widths
  const maxCount = genres && genres.length > 0
    ? Math.max(...genres.map((g) => g.book_count))
    : 1;

  return (
    <main className="genres-page">
      {showError && (
        <Toast
          message="Couldn't load genres. Try refreshing the page."
          onDismiss={() => setShowError(false)}
        />
      )}

      <h1 className="genres-page-heading">Genres</h1>

      {isPending ? (
        // Loading: 6 skeleton rows — 3 blocks each
        <div className="genres-list" aria-busy="true">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="genres-row genres-row--skeleton" aria-hidden="true">
              <div className="genres-skeleton-label" />
              <div className="genres-skeleton-bar" />
              <div className="genres-skeleton-count" />
            </div>
          ))}
        </div>
      ) : genres && genres.length > 0 ? (
        <div className="genres-list">
          {genres.map((genre) => (
            <Link
              key={genre.slug}
              to={`/genres/${genre.slug}`}
              className="genres-row"
            >
              <span className="genres-row-name">{genre.name}</span>
              <span className="genres-row-bar-container" aria-hidden="true">
                <span
                  className="genres-row-bar-fill"
                  style={{ width: `${(genre.book_count / maxCount) * 100}%` }}
                />
              </span>
              <span className="genres-row-count">{genre.book_count}</span>
            </Link>
          ))}
        </div>
      ) : (
        <p className="genres-empty">No genres found.</p>
      )}
    </main>
  );
}
