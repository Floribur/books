import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchAuthors } from '../api/books';
import { Toast } from '../components/Toast';
import { useState, useEffect } from 'react';
import { usePageTitle } from '../hooks/usePageTitle';
import './AuthorsPage.css';

export function AuthorsPage() {
  usePageTitle('Authors'); // "Flo's Library — Authors"
  const [showError, setShowError] = useState(false);

  const { data: authors, isPending, isError } = useQuery({
    queryKey: ['authors'],
    queryFn: fetchAuthors,
  });

  useEffect(() => {
    if (isError) setShowError(true);
  }, [isError]);

  return (
    <main className="authors-page">
      {showError && (
        <Toast
          message="Couldn't load authors. Try refreshing the page."
          onDismiss={() => setShowError(false)}
        />
      )}

      <h1 className="authors-page-heading">Authors</h1>

      {isPending ? (
        // Loading: 8 skeleton rows
        <ul className="authors-list" aria-busy="true">
          {Array.from({ length: 8 }).map((_, i) => (
            <li key={i} className="authors-list-row authors-list-row--skeleton" aria-hidden="true">
              <div className="authors-skeleton-name" />
            </li>
          ))}
        </ul>
      ) : authors && authors.length > 0 ? (
        <ul className="authors-list">
          {authors.map((author) => (
            <li key={author.slug} className="authors-list-item">
              <Link to={`/authors/${author.slug}`} className="authors-list-row">
                <span className="authors-list-name">{author.name}</span>
                <span className="authors-list-count">
                  {author.book_count === 1 ? '1 book' : `${author.book_count} books`}
                </span>
              </Link>
            </li>
          ))}
        </ul>
      ) : (
        <p className="authors-empty">No authors found.</p>
      )}
    </main>
  );
}
