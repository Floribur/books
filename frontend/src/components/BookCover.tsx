import { useState } from 'react';
import './BookCover.css';

interface BookCoverProps {
  src: string;
  title: string;
  loading?: 'lazy' | 'eager';
}

function placeholderUrl(title: string): string {
  const text = encodeURIComponent(title.slice(0, 25));
  return `https://placehold.co/128x192/1a0f12/9e8085?text=${text}`;
}

export function BookCover({ src, title, loading = 'lazy' }: BookCoverProps) {
  const [imgSrc, setImgSrc] = useState(src);

  return (
    <div className="book-cover-wrapper">
      <img
        src={imgSrc}
        alt={`${title} cover`}
        loading={loading}
        onError={() => setImgSrc(placeholderUrl(title))}
        className="book-cover-img"
      />
    </div>
  );
}
