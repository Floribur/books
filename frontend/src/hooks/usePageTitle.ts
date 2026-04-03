import { useEffect } from 'react';

/**
 * Sets document.title for the current page.
 * - usePageTitle()          → "Flo's Library"
 * - usePageTitle('Authors') → "Flo's Library — Authors"
 * Resets to base title on unmount.
 */
export function usePageTitle(page?: string): void {
  useEffect(() => {
    document.title = page ? `Flo's Library \u2014 ${page}` : "Flo's Library";
    return () => {
      document.title = "Flo's Library";
    };
  }, [page]);
}
