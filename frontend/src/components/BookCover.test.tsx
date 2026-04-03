import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { BookCover } from './BookCover';

describe('BookCover', () => {
  it('renders img with correct src and alt text', () => {
    render(<BookCover src="/covers/test.jpg" title="Dune" />);
    const img = screen.getByRole('img', { name: 'Dune cover' });
    expect(img).toHaveAttribute('src', '/covers/test.jpg');
  });

  it('uses loading="lazy" by default', () => {
    render(<BookCover src="/covers/test.jpg" title="Dune" />);
    expect(screen.getByRole('img')).toHaveAttribute('loading', 'lazy');
  });

  it('accepts loading="eager" for above-fold covers', () => {
    render(<BookCover src="/covers/test.jpg" title="Dune" loading="eager" />);
    expect(screen.getByRole('img')).toHaveAttribute('loading', 'eager');
  });

  it('hides the img when image errors (shows gradient placeholder only)', () => {
    render(<BookCover src="/covers/broken.jpg" title="Dune" />);
    const img = screen.getByRole('img');
    fireEvent.error(img);
    expect(screen.queryByRole('img')).toBeNull();
  });
});
