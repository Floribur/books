import './YearSelector.css';

interface YearSelectorProps {
  years: number[];          // all available years, sorted descending (newest first)
  selectedYear: number;
  onChange: (year: number) => void;
}

export function YearSelector({ years, selectedYear, onChange }: YearSelectorProps) {
  const currentIndex = years.indexOf(selectedYear);
  const hasPrev = currentIndex < years.length - 1;  // older year = higher index (desc sorted)
  const hasNext = currentIndex > 0;                  // newer year = lower index

  function handlePrev() {
    if (hasPrev) onChange(years[currentIndex + 1]);
  }

  function handleNext() {
    if (hasNext) onChange(years[currentIndex - 1]);
  }

  // Keyboard: ArrowLeft = prev (older), ArrowRight = next (newer) — D-11
  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'ArrowLeft') { e.preventDefault(); handlePrev(); }
    if (e.key === 'ArrowRight') { e.preventDefault(); handleNext(); }
  }

  return (
    <div
      className="year-selector"
      role="group"
      aria-label="Year selector"
      onKeyDown={handleKeyDown}
    >
      {/* Prev arrow — goes to an older year */}
      <button
        className="year-selector-arrow"
        onClick={handlePrev}
        disabled={!hasPrev}
        aria-label="Previous year"
      >
        {/* Left chevron SVG — inline, 20x20 */}
        <svg width="20" height="20" viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="M12.5 15L7.5 10L12.5 5" stroke="currentColor" strokeWidth="1.5"
            strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </button>

      {/* Year display */}
      <span className="year-selector-display" aria-live="polite" aria-atomic="true">
        {selectedYear}
      </span>

      {/* Next arrow — goes to a newer year */}
      <button
        className="year-selector-arrow"
        onClick={handleNext}
        disabled={!hasNext}
        aria-label="Next year"
      >
        {/* Right chevron SVG — inline, 20x20 */}
        <svg width="20" height="20" viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="M7.5 5L12.5 10L7.5 15" stroke="currentColor" strokeWidth="1.5"
            strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </button>
    </div>
  );
}
