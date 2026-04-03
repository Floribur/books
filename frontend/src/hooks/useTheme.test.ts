import { renderHook, act } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { useTheme } from './useTheme';

describe('useTheme', () => {
  beforeEach(() => {
    localStorage.clear();
    document.documentElement.removeAttribute('data-theme');
    // matchMedia mock already set in setup.ts — default returns matches: false (light preference)
  });

  it('defaults to light when no localStorage value and OS prefers light', () => {
    const { result } = renderHook(() => useTheme());
    expect(result.current.theme).toBe('light');
  });

  it('reads stored theme from localStorage on first load', () => {
    localStorage.setItem('theme', 'dark');
    const { result } = renderHook(() => useTheme());
    expect(result.current.theme).toBe('dark');
  });

  it('sets [data-theme] attribute on document.documentElement after mount', () => {
    localStorage.setItem('theme', 'dark');
    renderHook(() => useTheme());
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark');
  });

  it('persists new theme to localStorage on toggle', () => {
    const { result } = renderHook(() => useTheme());
    act(() => { result.current.toggle(); });
    expect(localStorage.getItem('theme')).toBe('dark');
  });

  it('toggles between light and dark', () => {
    const { result } = renderHook(() => useTheme());
    expect(result.current.theme).toBe('light');
    act(() => { result.current.toggle(); });
    expect(result.current.theme).toBe('dark');
    act(() => { result.current.toggle(); });
    expect(result.current.theme).toBe('light');
  });
});
