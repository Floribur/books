import { useSearchParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchYears, fetchBooksByYear, fetchAllBooksByYear } from '../api/books';
import { BookGrid } from '../components/BookGrid';
import { YearSelector } from '../components/YearSelector';
import { StatsStrip } from '../components/StatsStrip';
import { Toast } from '../components/Toast';
import { useState, useEffect, useMemo } from 'react';
import { usePageTitle } from '../hooks/usePageTitle';
import './ReadingChallengePage.css';

export function ReadingChallengePage() {
  usePageTitle('Reading Challenge'); // "Flo's Library — Reading Challenge"
  const [searchParams, setSearchParams] = useSearchParams();
  const [showError, setShowError] = useState(false);

  // Fetch years list — used to populate YearSelector and determine default year (D-08)
  const { data: yearEntries, isPending: yearsLoading, isError: yearsError } = useQuery({
    queryKey: ['years'],
    queryFn: fetchYears,
  });

  useEffect(() => {
    if (yearsError) setShowError(true);
  }, [yearsError]);

  // Extract sorted year numbers (descending — newest first)
  const years = useMemo(
    () => (yearEntries ?? []).map((e) => e.year).sort((a, b) => b - a),
    [yearEntries]
  );

  // Resolve selected year:
  // 1. If ?year= URL param is set and valid, use it
  // 2. Otherwise default to most recent year with books (D-08)
  const yearParam = searchParams.get('year');
  const parsedParam = yearParam ? parseInt(yearParam, 10) : NaN;
  const selectedYear = useMemo(() => {
    if (!isNaN(parsedParam) && years.includes(parsedParam)) return parsedParam;
    return years[0] ?? new Date().getFullYear(); // fallback to current year if no data yet
  }, [parsedParam, years]);

  function handleYearChange(newYear: number) {
    setSearchParams({ year: String(newYear) }, { replace: false });
  }

  // Book count for selected year (from yearEntries — no extra fetch needed)
  const yearEntry = yearEntries?.find((e) => e.year === selectedYear);
  const bookCount = yearEntry?.book_count ?? 0;

  // Fetch all books for the year (limit=200) to compute page stats
  const { data: allYearBooks } = useQuery({
    queryKey: ['books', 'year', String(selectedYear), 'all'],
    queryFn: () => fetchAllBooksByYear(selectedYear),
    enabled: !yearsLoading && years.length > 0,
  });

  const pageStats = useMemo(() => {
    if (!allYearBooks) return { totalPages: null, longestBook: null };
    const withPages = allYearBooks.filter((b) => b.page_count != null);
    if (withPages.length === 0) return { totalPages: null, longestBook: null };
    const totalPages = withPages.reduce((sum, b) => sum + b.page_count!, 0);
    const longest = withPages.reduce((a, b) => (b.page_count! > a.page_count! ? b : a));
    return {
      totalPages,
      longestBook: { title: longest.title, pageCount: longest.page_count! },
    };
  }, [allYearBooks]);

  // fetchFn for BookGrid — closure over selectedYear (CHAL-01)
  // Must recreate when selectedYear changes — useMemo ensures stable reference per year
  const fetchFn = useMemo(
    () => (cursor: string | undefined) => fetchBooksByYear(selectedYear, cursor),
    [selectedYear]
  );

  return (
    <main className="reading-challenge-page">
      {showError && (
        <Toast
          message="Couldn't load years. Try refreshing the page."
          onDismiss={() => setShowError(false)}
        />
      )}

      {/* Page heading */}
      <h1 className="reading-challenge-heading">Reading Challenge</h1>

      {/* Year selector — D-11 (CHAL-02) */}
      {!yearsLoading && years.length > 0 && (
        <YearSelector
          years={years}
          selectedYear={selectedYear}
          onChange={handleYearChange}
        />
      )}

      {/* Stats strip — D-09, D-10 (CHAL-03) */}
      <StatsStrip
        stats={{
          bookCount,
          totalPages: pageStats.totalPages,
          longestBook: pageStats.longestBook,
        }}
        year={selectedYear}
        isLoading={yearsLoading}
      />

      {/* Empty year message — shown above BookGrid when bookCount is 0 and not loading */}
      {!yearsLoading && bookCount === 0 && (
        <p className="reading-challenge-empty">
          No books recorded for {selectedYear}.
        </p>
      )}

      {/* Book grid filtered to selected year — CHAL-01 */}
      {!yearsLoading && years.length > 0 && (
        <BookGrid
          queryKey={['books', 'year', String(selectedYear)]}
          fetchFn={fetchFn}
          ariaLabel={`Books read in ${selectedYear}`}
        />
      )}
    </main>
  );
}
