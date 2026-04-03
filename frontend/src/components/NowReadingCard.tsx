import { Link } from 'react-router-dom';
import { BookCover } from './BookCover';
import type { Book } from '../api/types';
import './NowReadingCard.css';

interface NowReadingCardProps {
  book: Book;
}

export function NowReadingCard({ book }: NowReadingCardProps) {
  const authorNames = book.authors.map((a) => a.name).join(', ');

  return (
    <article className="now-reading-card">
      <Link to={`/books/${book.slug}`} className="now-reading-card-link">
        {/* loading="eager" — Now Reading is above the fold */}
        <BookCover src={book.cover_path} title={book.title} loading="eager" />
        <div className="now-reading-card-info">
          <p className="now-reading-card-title">{book.title}</p>
          <p className="now-reading-card-author">{authorNames}</p>
        </div>
      </Link>
    </article>
  );
}
