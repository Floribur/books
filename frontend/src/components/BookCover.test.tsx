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

  it('uses placehold.co immediately when src is empty', () => {
    render(<BookCover src="" title="Dune" />);
    const img = screen.getByRole('img', { name: 'Dune cover' });
    expect(img.getAttribute('src')).toContain('placehold.co');
  });

  it('switches to placehold.co URL when image errors', () => {
    render(<BookCover src="/covers/broken.jpg" title="Dune" />);
    const img = screen.getByRole('img');
    fireEvent.error(img);
    expect(img).toHaveAttribute('src', expect.stringContaining('placehold.co'));
  });
});
