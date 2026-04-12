import { useQuery } from '@tanstack/react-query';
import { Bio } from '../components/Bio';
import { NowReadingSection } from '../components/NowReadingSection';
import { BookGrid } from '../components/BookGrid';
import { fetchBooks } from '../api/books';
import { usePageTitle } from '../hooks/usePageTitle';
import './HomePage.css';

export function HomePage() {
  usePageTitle();
  const { data: allBooks, isPending, isError } = useQuery({
    queryKey: ['books'],
    queryFn: fetchBooks,
  });

  const readBooks = allBooks?.filter((b) => b.shelf === 'read') ?? [];

  return (
    <main className="home-page">
      <Bio />
      <NowReadingSection />
      <section aria-label="Books Read">
        <h2 className="section-heading">Books Read</h2>
        <BookGrid
          books={readBooks}
          isPending={isPending}
          isError={isError}
          ariaLabel="Books Read"
        />
      </section>
    </main>
  );
}
