import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { Bio } from './Bio';

describe('Bio', () => {
  it('renders bio photo with correct alt text', () => {
    render(<Bio />);
    const img = screen.getByRole('img', { name: 'Florian' });
    expect(img).toHaveAttribute('src', '/florian.jpg');
  });

  it('renders parsed markdown content (heading or paragraph)', () => {
    render(<Bio />);
    // bio.md has a heading "# Florian" — marked parses it to <h1>
    expect(screen.getByRole('heading', { name: 'Florian' })).toBeInTheDocument();
  });

  it('includes Goodreads profile link in rendered content', () => {
    render(<Bio />);
    const link = screen.getByRole('link', { name: /Goodreads/i });
    expect(link).toHaveAttribute('href', 'https://www.goodreads.com/user/show/79499864-florian');
  });
});
