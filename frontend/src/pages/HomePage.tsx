import { Bio } from '../components/Bio';
import { NowReadingSection } from '../components/NowReadingSection';
import { BookGrid } from '../components/BookGrid';
import './HomePage.css';

export function HomePage() {
  return (
    <main className="home-page">
      <Bio />
      <NowReadingSection />

      {/* Books Read section heading + grid */}
      <section aria-label="Books Read">
        <h2 className="section-heading">Books Read</h2>
        <BookGrid />
      </section>
    </main>
  );
}
