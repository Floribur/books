import { renderHook } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';
import { useIntersectionObserver } from './useIntersectionObserver';

// Note: IntersectionObserver is mocked in src/test/setup.ts
// The mock: vi.fn(() => ({ observe: vi.fn(), disconnect: vi.fn(), unobserve: vi.fn() }))

describe('useIntersectionObserver', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('does not create an observer when enabled is false', () => {
    const callback = vi.fn();
    renderHook(() => useIntersectionObserver(callback, false));
    expect(global.IntersectionObserver).not.toHaveBeenCalled();
  });

  it('creates an IntersectionObserver with threshold 0.1 when enabled is true and ref is attached', () => {
    const callback = vi.fn();
    renderHook(() => useIntersectionObserver(callback, true));
    expect(global.IntersectionObserver).toHaveBeenCalledWith(
      expect.any(Function),
      { threshold: 0.1 }
    );
  });

  it('calls disconnect on cleanup', () => {
    const mockDisconnect = vi.fn();
    (global.IntersectionObserver as ReturnType<typeof vi.fn>).mockImplementation(function() {
      return {
        observe: vi.fn(),
        disconnect: mockDisconnect,
        unobserve: vi.fn(),
      };
    });

    const callback = vi.fn();
    const { unmount } = renderHook(() => useIntersectionObserver(callback, true));
    unmount();
    expect(mockDisconnect).toHaveBeenCalled();
  });
});
