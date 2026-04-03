// Scroll restoration for back-button navigation.
// Known limitation: if the target book is on page 2+ (not yet loaded), scroll is silently skipped.

export function useScrollRestoration() {
  function saveTarget(slug: string) {
    sessionStorage.setItem('scrollTarget', slug);
  }

  function restoreScroll() {
    const target = sessionStorage.getItem('scrollTarget');
    if (!target) return;

    const el = document.getElementById(`book-${target}`);
    if (!el) return; // Book not yet loaded — skip silently

    el.scrollIntoView({ behavior: 'instant' });
    sessionStorage.removeItem('scrollTarget');
  }

  return { saveTarget, restoreScroll };
}
