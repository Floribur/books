import { useQuery } from '@tanstack/react-query';
import { NowReadingCard } from './NowReadingCard';
import { Toast } from './Toast';
import { fetchCurrentlyReading } from '../api/books';
import './NowReadingSection.css';

const NOW_READING_CAP = 4; // D-05: cap at 4 books; resolved in UI-SPEC

export function NowReadingSection() {
  const { data: allBooks, isError } = useQuery({
    queryKey: ['currently-reading'],
    queryFn: fetchCurrentlyReading,
  });

  // D-15: Show error toast on fetch failure (checked BEFORE empty-array guard)
  if (isError) {
    return <Toast message="Couldn't load currently reading. Try refreshing." onDismiss={() => {}} />;
  }

  // D-04: Hide section entirely when API returns empty array
  // Also hide while loading (data is undefined) — section is not critical path
  if (!allBooks || allBooks.length === 0) return null;

  // D-05: Cap at 4 books; silently omit extras
  const books = allBooks.slice(0, NOW_READING_CAP);

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
