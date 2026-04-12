import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchAuthors, fetchBooks } from '../api/books';
import { BookGrid } from '../components/BookGrid';
import { Toast } from '../components/Toast';
import { useState, useEffect } from 'react';
import { usePageTitle } from '../hooks/usePageTitle';
import './AuthorDetailPage.css';

export function AuthorDetailPage() {
  const { slug } = useParams<{ slug: string }>();
  const [showError, setShowError] = useState(false);

  const { data: allAuthors, isError: authorsError } = useQuery({
    queryKey: ['authors'],
    queryFn: fetchAuthors,
  });
  const { data: allBooks, isPending: booksPending, isError: booksError } = useQuery({
    queryKey: ['books'],
    queryFn: fetchBooks,
  });

  const isError = authorsError || booksError;
  useEffect(() => {
    if (isError) setShowError(true);
  }, [isError]);

  const author = allAuthors?.find((a) => a.slug === slug);
  usePageTitle(author?.name);

  const authorBooks = allBooks?.filter((b) =>
    b.authors.some((a) => a.slug === slug)
  ) ?? [];

  return (
    <main className="author-detail-page">
      {showError && (
        <Toast
          message="Couldn't load author. Try refreshing the page."
          onDismiss={() => setShowError(false)}
        />
      )}

      {booksPending ? (
        <div className="author-detail-heading-skeleton" aria-hidden="true" />
      ) : author ? (
        <div className="author-detail-heading">
          <h1 className="author-detail-name">{author.name}</h1>
        </div>
      ) : null}

      <BookGrid
        books={authorBooks}
        isPending={booksPending}
        isError={booksError}
        ariaLabel={`Books by ${author?.name ?? 'this author'}`}
      />
    </main>
  );
}
