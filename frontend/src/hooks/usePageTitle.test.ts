import { describe, it, expect, afterEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { usePageTitle } from './usePageTitle';

describe('usePageTitle', () => {
  afterEach(() => {
    document.title = '';
  });

  it('sets document.title to "Flo\'s Library" when called with no argument', () => {
    renderHook(() => usePageTitle());
    expect(document.title).toBe("Flo's Library");
  });

  it('sets document.title to "Flo\'s Library — Authors" when called with "Authors"', () => {
    renderHook(() => usePageTitle('Authors'));
    expect(document.title).toBe("Flo's Library \u2014 Authors");
  });

  it('sets document.title to "Flo\'s Library — Dune" when called with "Dune"', () => {
    renderHook(() => usePageTitle('Dune'));
    expect(document.title).toBe("Flo's Library \u2014 Dune");
  });

  it('resets document.title to "Flo\'s Library" on unmount', () => {
    const { unmount } = renderHook(() => usePageTitle('Reading Challenge'));
    expect(document.title).toBe("Flo's Library \u2014 Reading Challenge");
    unmount();
    expect(document.title).toBe("Flo's Library");
  });
});
