import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchAuthorBySlug, fetchBooksByAuthor } from '../api/books';
import { BookGrid } from '../components/BookGrid';
import { Toast } from '../components/Toast';
import { useState, useEffect } from 'react';
import './AuthorDetailPage.css';

export function AuthorDetailPage() {
  const { slug } = useParams<{ slug: string }>();
  const [showError, setShowError] = useState(false);

  const { data: author, isPending, isError } = useQuery({
    queryKey: ['author', slug],
    queryFn: () => fetchAuthorBySlug(slug!),
    enabled: !!slug,
  });

  useEffect(() => {
    if (isError) setShowError(true);
  }, [isError]);

  // fetchFn for BookGrid: closure over slug
  const fetchFn = (cursor: string | undefined) => fetchBooksByAuthor(slug!, cursor);

  return (
    <main className="author-detail-page">
      {showError && (
        <Toast
          message="Couldn't load author. Try refreshing the page."
          onDismiss={() => setShowError(false)}
        />
      )}

      {isPending ? (
        <div className="author-detail-heading-skeleton" aria-hidden="true" />
      ) : author ? (
        <div className="author-detail-heading">
          <h1 className="author-detail-name">{author.name}</h1>
          <p className="author-detail-count">
            {author.book_count === 1 ? '1 book' : `${author.book_count} books`}
          </p>
        </div>
      ) : null}

      {slug && (
        <BookGrid
          queryKey={['books', 'author', slug]}
          fetchFn={fetchFn}
          ariaLabel={`Books by ${author?.name ?? 'this author'}`}
        />
      )}
    </main>
  );
}
