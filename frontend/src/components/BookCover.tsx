import { useState } from 'react';
import './BookCover.css';

interface BookCoverProps {
  src: string;
  title: string;
  loading?: 'lazy' | 'eager';
}

export function BookCover({ src, title, loading = 'lazy' }: BookCoverProps) {
  const [hasError, setHasError] = useState(false);

  return (
    <div className="book-cover-wrapper">
      {hasError ? (
        <div className="book-cover-placeholder" aria-label={`${title} cover`}>
          <span className="book-cover-placeholder-title">{title}</span>
        </div>
      ) : (
        <img
          src={src}
          alt={`${title} cover`}
          loading={loading}
          onError={() => setHasError(true)}
          className="book-cover-img"
        />
      )}
    </div>
  );
}
