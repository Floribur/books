import { describe, expect, it, beforeEach, vi } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useScrollRestoration } from './useScrollRestoration';

describe('useScrollRestoration', () => {
  beforeEach(() => {
    sessionStorage.clear();
  });

  it('saveTarget stores slug in sessionStorage', () => {
    const { result } = renderHook(() => useScrollRestoration());
    result.current.saveTarget('dune-1965');
    expect(sessionStorage.getItem('scrollTarget')).toBe('dune-1965');
  });

  it('restoreScroll scrolls to element with id book-{slug} if found', () => {
    sessionStorage.setItem('scrollTarget', 'dune-1965');

    const el = document.createElement('div');
    el.id = 'book-dune-1965';
    el.scrollIntoView = vi.fn();
    document.body.appendChild(el);

    const { result } = renderHook(() => useScrollRestoration());
    result.current.restoreScroll();

    expect(el.scrollIntoView).toHaveBeenCalledWith({ behavior: 'instant' });
    expect(sessionStorage.getItem('scrollTarget')).toBeNull();

    document.body.removeChild(el);
  });

  it('removes sessionStorage entry after successful scroll restoration', () => {
    sessionStorage.setItem('scrollTarget', 'dune-1965');

    const el = document.createElement('div');
    el.id = 'book-dune-1965';
    el.scrollIntoView = vi.fn();
    document.body.appendChild(el);

    const { result } = renderHook(() => useScrollRestoration());
    result.current.restoreScroll();

    expect(sessionStorage.getItem('scrollTarget')).toBeNull();
    document.body.removeChild(el);
  });

  it('does nothing when element is not found (book not yet loaded)', () => {
    sessionStorage.setItem('scrollTarget', 'book-not-in-dom');

    const { result } = renderHook(() => useScrollRestoration());
    expect(() => result.current.restoreScroll()).not.toThrow();
    expect(sessionStorage.getItem('scrollTarget')).toBe('book-not-in-dom');
  });
});
