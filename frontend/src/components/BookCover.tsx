import { useState } from 'react';
import './BookCover.css';

interface BookCoverProps {
  src: string;
  title: string;
  loading?: 'lazy' | 'eager';
}

function placeholderUrl(title: string): string {
  // Wrap title into word-boundary lines, show up to 3, add … if truncated
  const maxChars = 16;
  const maxLines = 3;
  const words = title.split(' ');
  const lines: string[] = [];
  let current = '';

  for (const word of words) {
    const candidate = current ? `${current} ${word}` : word;
    if (current && candidate.length > maxChars) {
      lines.push(current);
      current = word;
    } else {
      current = candidate;
    }
  }
  if (current) lines.push(current);

  const display = lines.slice(0, maxLines);
  if (lines.length > maxLines) display[maxLines - 1] += '…';

  const text = encodeURIComponent(display.join('\n'));
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
