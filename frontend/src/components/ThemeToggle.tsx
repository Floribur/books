import { Sun, Moon } from 'lucide-react';
import { useTheme } from '../hooks/useTheme';
import './ThemeToggle.css';

export function ThemeToggle() {
  const { theme, toggle } = useTheme();

  return (
    <button
      className="theme-toggle"
      onClick={toggle}
      aria-label={theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}
    >
      {theme === 'light' ? (
        <Moon size={20} aria-hidden="true" />
      ) : (
        <Sun size={20} aria-hidden="true" />
      )}
    </button>
  );
}
