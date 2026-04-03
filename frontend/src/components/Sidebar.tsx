import { useState, useEffect, useRef } from 'react';
import { NavLink } from 'react-router-dom';
import { Menu, X, Library, Users, Tag, Trophy } from 'lucide-react';
import { ThemeToggle } from './ThemeToggle';
import './Sidebar.css';

const navLinks = [
  { to: '/', label: 'Home', icon: Library },
  { to: '/authors', label: 'Authors', icon: Users },
  { to: '/genres', label: 'Genres', icon: Tag },
  { to: '/reading-challenge', label: 'Reading Challenge', icon: Trophy },
];

export function Sidebar() {
  const [isOpen, setIsOpen] = useState(false);
  const drawerRef = useRef<HTMLElement>(null);

  // Close drawer on nav link click
  function handleNavClick() {
    setIsOpen(false);
  }

  // Close drawer on backdrop click
  function handleBackdropClick() {
    setIsOpen(false);
  }

  // Prevent body scroll when drawer is open
  useEffect(() => {
    document.body.style.overflow = isOpen ? 'hidden' : '';
    return () => { document.body.style.overflow = ''; };
  }, [isOpen]);

  // Focus trap: Tab cycles through focusable elements inside the drawer
  useEffect(() => {
    if (!isOpen || !drawerRef.current) return;

    const focusable = drawerRef.current.querySelectorAll<HTMLElement>(
      'a[href], button:not([disabled])'
    );
    const first = focusable[0];
    const last = focusable[focusable.length - 1];

    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') {
        setIsOpen(false);
        return;
      }
      if (e.key !== 'Tab') return;
      if (e.shiftKey) {
        if (document.activeElement === first) {
          e.preventDefault();
          last.focus();
        }
      } else {
        if (document.activeElement === last) {
          e.preventDefault();
          first.focus();
        }
      }
    }

    document.addEventListener('keydown', handleKeyDown);
    first?.focus();
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen]);

  const navContent = (
    <>
      <div className="sidebar-header">
        <span className="sidebar-title">Flo's Library</span>
      </div>
      <nav className="sidebar-nav" aria-label="Main navigation">
        <ul className="sidebar-nav-list">
          {navLinks.map(({ to, label, icon: Icon }) => (
            <li key={to}>
              <NavLink
                to={to}
                end={to === '/'}
                className={({ isActive }) =>
                  ['sidebar-nav-link', isActive ? 'sidebar-nav-link--active' : ''].join(' ')
                }
                onClick={handleNavClick}
              >
                <Icon size={16} className="sidebar-nav-icon" aria-hidden="true" />
                {label}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
      <div className="sidebar-toggle-wrapper">
        <ThemeToggle />
      </div>
    </>
  );

  return (
    <>
      {/* Desktop sidebar — visible at ≥768px */}
      <aside className="sidebar sidebar--desktop" aria-label="Site navigation">
        {navContent}
      </aside>

      {/* Mobile top bar — visible at <768px */}
      <div className="mobile-topbar">
        <span className="sidebar-title">Flo's Library</span>
        <button
          className="hamburger-button"
          onClick={() => setIsOpen(true)}
          aria-label="Open navigation"
          aria-expanded={isOpen}
          aria-controls="mobile-drawer"
        >
          <Menu size={24} aria-hidden="true" />
        </button>
      </div>

      {/* Mobile drawer backdrop */}
      {isOpen && (
        <div
          className="drawer-backdrop"
          onClick={handleBackdropClick}
          aria-hidden="true"
        />
      )}

      {/* Mobile drawer */}
      <aside
        id="mobile-drawer"
        className={['sidebar', 'sidebar--drawer', isOpen ? 'sidebar--drawer-open' : ''].join(' ')}
        ref={drawerRef}
        aria-label="Site navigation"
        aria-hidden={!isOpen}
      >
        <button
          className="drawer-close-button"
          onClick={() => setIsOpen(false)}
          aria-label="Close navigation"
        >
          <X size={24} aria-hidden="true" />
        </button>
        {navContent}
      </aside>
    </>
  );
}
