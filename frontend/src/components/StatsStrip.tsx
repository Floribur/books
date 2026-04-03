import './StatsStrip.css';

// totalPages and longestBook are null when page_count data is unavailable
// for the year (D-10). The component hides those stats when null.
// The current GET /api/books?year= endpoint returns Book (list shape) which
// lacks page_count — so these are always null until the API is extended.
// TODO: when API supports page_count on the list endpoint, compute these in ReadingChallengePage.
interface StatsData {
  bookCount: number;
  totalPages: number | null;   // null = hide stat + footnote
  longestBook: { title: string; pageCount: number } | null;  // null = hide stat
}

interface StatsStripProps {
  stats: StatsData;
  year: number;
  isLoading?: boolean;
}


export function StatsStrip({ stats, year, isLoading = false }: StatsStripProps) {
  const showPageStats = stats.totalPages !== null && stats.totalPages > 0;
  const showLongest = stats.longestBook !== null;
  const showFootnote = showPageStats || showLongest;


  if (isLoading) {
    return (
      <div className="stats-strip">
        <div className="stats-strip-inner">
          {[1, 2, 3].map((i) => (
            <div key={i} className="stats-strip-skeleton" aria-hidden="true" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="stats-strip-wrapper">
      <div className="stats-strip">
        {/* Stat 1: Books Read — always visible (CHAL-03) */}
        <div className="stats-strip-stat">
          <span className="stats-strip-value">{stats.bookCount}</span>
          <span className="stats-strip-label">books read in {year}</span>
        </div>

        {/* Stat 2: Total Pages — hidden when totalPages is null or 0 (D-10) */}
        {showPageStats && (
          <div className="stats-strip-stat">
            <span className="stats-strip-value">{stats.totalPages!.toLocaleString()}</span>
            <span className="stats-strip-label">pages read*</span>
          </div>
        )}

        {/* Stat 3: Longest Book — hidden when null (D-10) */}
        {showLongest && (
          <div className="stats-strip-stat">
            <span className="stats-strip-value">{stats.longestBook!.pageCount.toLocaleString()}</span>
            <span className="stats-strip-label">pages, longest book*</span>
          </div>
        )}
      </div>

      {/* Footnote — only when page stats are visible (D-10) */}
      {showFootnote && (
        <p className="stats-strip-footnote">
          * Based on available page data — not all books have page counts.
        </p>
      )}
    </div>
  );
}
