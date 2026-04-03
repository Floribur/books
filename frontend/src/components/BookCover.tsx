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
      {!hasError && (
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
