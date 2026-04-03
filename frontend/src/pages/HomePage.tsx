import { Bio } from '../components/Bio';
import { NowReadingSection } from '../components/NowReadingSection';
import { BookGrid } from '../components/BookGrid';
import { fetchBooks } from '../api/books';
import { usePageTitle } from '../hooks/usePageTitle';
import './HomePage.css';

export function HomePage() {
  usePageTitle(); // sets "Flo's Library"
  return (
    <main className="home-page">
      <Bio />
      <NowReadingSection />

      {/* Books Read section heading + grid */}
      <section aria-label="Books Read">
        <h2 className="section-heading">Books Read</h2>
        <BookGrid
          queryKey={['books']}
          fetchFn={fetchBooks}
          ariaLabel="Books Read"
        />
      </section>
    </main>
  );
}
