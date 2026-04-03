import { Link } from 'react-router-dom';
import { BookCover } from './BookCover';
import { useScrollRestoration } from '../hooks/useScrollRestoration';
import type { Book } from '../api/types';
import './BookCard.css';

interface BookCardProps {
  book: Book;
}

export function BookCard({ book }: BookCardProps) {
  const { saveTarget } = useScrollRestoration();

  function handleClick() {
    saveTarget(book.slug);
  }

  const authorNames = book.authors.map((a) => a.name).join(', ');

  return (
    <article className="book-card" id={`book-${book.slug}`}>
      <Link to={`/books/${book.slug}`} onClick={handleClick} className="book-card-link">
        <BookCover src={book.cover_path} title={book.title} loading="lazy" />
        <div className="book-card-info">
          <p className="book-card-title">{book.title}</p>
          <p className="book-card-author">{authorNames}</p>
        </div>
      </Link>
    </article>
  );
}
