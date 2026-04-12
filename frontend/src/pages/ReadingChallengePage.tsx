import { useSearchParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchBooks } from '../api/books';
import { BookGrid } from '../components/BookGrid';
import { YearSelector } from '../components/YearSelector';
import { StatsStrip } from '../components/StatsStrip';
import { Toast } from '../components/Toast';
import { useState, useEffect, useMemo } from 'react';
import { usePageTitle } from '../hooks/usePageTitle';
import './ReadingChallengePage.css';

export function ReadingChallengePage() {
  usePageTitle('Reading Challenge');
  const [searchParams, setSearchParams] = useSearchParams();
  const [showError, setShowError] = useState(false);

  const { data: allBooks, isPending, isError } = useQuery({
    queryKey: ['books'],
    queryFn: fetchBooks,
  });

  useEffect(() => {
    if (isError) setShowError(true);
  }, [isError]);

  // Derive years list from books.json: extract unique years with counts for 'read' shelf
  const yearEntries = useMemo(() => {
    if (!allBooks) return [];
    const counts = new Map<number, number>();
    for (const b of allBooks) {
      if (b.shelf === 'read' && b.read_at) {
        const year = new Date(b.read_at).getFullYear();
        if (!isNaN(year)) counts.set(year, (counts.get(year) ?? 0) + 1);
      }
    }
    return Array.from(counts.entries())
      .map(([year, book_count]) => ({ year, book_count }))
      .sort((a, b) => b.year - a.year);
  }, [allBooks]);

  const years = yearEntries.map((e) => e.year);

  const yearParam = searchParams.get('year');
  const parsedParam = yearParam ? parseInt(yearParam, 10) : NaN;
  const selectedYear = useMemo(() => {
    if (!isNaN(parsedParam) && years.includes(parsedParam)) return parsedParam;
    return years[0] ?? new Date().getFullYear();
  }, [parsedParam, years]);

  function handleYearChange(newYear: number) {
    setSearchParams({ year: String(newYear) }, { replace: false });
  }

  const yearEntry = yearEntries.find((e) => e.year === selectedYear);
  const bookCount = yearEntry?.book_count ?? 0;

  // Filter books for selected year (used for both grid and stats)
  const yearBooks = useMemo(
    () =>
      (allBooks ?? []).filter(
        (b) =>
          b.shelf === 'read' &&
          b.read_at &&
          new Date(b.read_at).getFullYear() === selectedYear
      ),
    [allBooks, selectedYear]
  );

  const pageStats = useMemo(() => {
    const withPages = yearBooks.filter((b) => b.page_count != null);
    if (withPages.length === 0) return { totalPages: null, longestBook: null };
    const totalPages = withPages.reduce((sum, b) => sum + b.page_count!, 0);
    const longest = withPages.reduce((a, b) => (b.page_count! > a.page_count! ? b : a));
    return {
      totalPages,
      longestBook: { title: longest.title, pageCount: longest.page_count! },
    };
  }, [yearBooks]);

  return (
    <main className="reading-challenge-page">
      {showError && (
        <Toast
          message="Couldn't load books. Try refreshing the page."
          onDismiss={() => setShowError(false)}
        />
      )}

      <h1 className="reading-challenge-heading">Reading Challenge</h1>

      {!isPending && years.length > 0 && (
        <YearSelector
          years={years}
          selectedYear={selectedYear}
          onChange={handleYearChange}
        />
      )}

      <StatsStrip
        stats={{
          bookCount,
          totalPages: pageStats.totalPages,
          longestBook: pageStats.longestBook,
        }}
        year={selectedYear}
        isLoading={isPending}
      />

      {!isPending && bookCount === 0 && (
        <p className="reading-challenge-empty">
          No books recorded for {selectedYear}.
        </p>
      )}

      {!isPending && years.length > 0 && (
        <BookGrid
          books={yearBooks}
          isPending={isPending}
          isError={isError}
          ariaLabel={`Books read in ${selectedYear}`}
        />
      )}
    </main>
  );
}
