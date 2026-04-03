import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { Bio } from './Bio';

describe('Bio', () => {
  it('renders bio photo with correct alt text', () => {
    render(<Bio />);
    const img = screen.getByRole('img', { name: 'Florian' });
    expect(img).toHaveAttribute('src', '/florian.jpg');
  });

  it('renders parsed markdown content (paragraph)', () => {
    render(<Bio />);
    expect(screen.getByText(/Hoi zäme/i)).toBeInTheDocument();
  });

  it('includes Goodreads profile link in rendered content opening in new tab', () => {
    render(<Bio />);
    const link = screen.getByRole('link', { name: /Goodreads/i });
    expect(link).toHaveAttribute('href', 'https://www.goodreads.com/user/show/79499864-florian');
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveAttribute('rel', 'noopener noreferrer');
  });
});
