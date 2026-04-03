import { useState } from 'react';
import './BookCover.css';

interface BookCoverProps {
  src: string;
  title: string;
  loading?: 'lazy' | 'eager';
}

function placeholderUrl(title: string): string {
  // Wrap title into lines of ~14 chars at word boundaries (max 3 lines)
  const words = title.split(' ');
  const lines: string[] = [];
  let current = '';
  for (const word of words) {
    if (current && (current + ' ' + word).length > 14) {
      lines.push(current);
      current = word;
      if (lines.length === 3) break;
    } else {
      current = current ? `${current} ${word}` : word;
    }
  }
  if (current && lines.length < 3) lines.push(current);
  const text = encodeURIComponent(lines.join('\n'));
  return `https://placehold.co/128x192/1a0f12/9e8085?text=${text}`;
}

export function BookCover({ src, title, loading = 'lazy' }: BookCoverProps) {
  // Use placehold.co immediately if src is empty/missing
  const [imgSrc, setImgSrc] = useState(() => src || placeholderUrl(title));

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
