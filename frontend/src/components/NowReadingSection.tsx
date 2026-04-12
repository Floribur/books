import { useQuery } from '@tanstack/react-query';
import { NowReadingCard } from './NowReadingCard';
import { Toast } from './Toast';
import { fetchBooks } from '../api/books';
import './NowReadingSection.css';

const NOW_READING_CAP = 4;

export function NowReadingSection() {
  const { data: allBooks, isError } = useQuery({
    queryKey: ['books'],
    queryFn: fetchBooks,
  });

  if (isError) {
    return <Toast message="Couldn't load currently reading. Try refreshing." onDismiss={() => {}} />;
  }

  if (!allBooks) return null;

  const books = allBooks
    .filter((b) => b.shelf === 'currently-reading')
    .slice(0, NOW_READING_CAP);

  if (books.length === 0) return null;

  return (
    <section className="now-reading-section" aria-label="Now Reading">
      <h2 className="section-heading">Now Reading</h2>
      <div className="now-reading-grid">
        {books.map((book) => (
          <NowReadingCard key={book.slug} book={book} />
        ))}
      </div>
    </section>
  );
}
