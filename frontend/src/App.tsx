import { Routes, Route } from 'react-router-dom';
import { Sidebar } from './components/Sidebar';
import { HomePage } from './pages/HomePage';
import './App.css';

// Phase 4 pages — stubs remain
function BookDetailPage() { return <main style={{ padding: '2rem' }}><p>Book detail — Phase 4</p></main>; }
function AuthorsPage() { return <main style={{ padding: '2rem' }}><p>Authors — Phase 4</p></main>; }
function AuthorDetailPage() { return <main style={{ padding: '2rem' }}><p>Author detail — Phase 4</p></main>; }
function GenresPage() { return <main style={{ padding: '2rem' }}><p>Genres — Phase 4</p></main>; }
function GenreDetailPage() { return <main style={{ padding: '2rem' }}><p>Genre detail — Phase 4</p></main>; }
function ReadingChallengePage() { return <main style={{ padding: '2rem' }}><p>Reading Challenge — Phase 4</p></main>; }

export default function App() {
  return (
    <div className="app-layout">
      <Sidebar />
      <div className="app-content">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/books/:slug" element={<BookDetailPage />} />
          <Route path="/authors" element={<AuthorsPage />} />
          <Route path="/authors/:slug" element={<AuthorDetailPage />} />
          <Route path="/genres" element={<GenresPage />} />
          <Route path="/genres/:slug" element={<GenreDetailPage />} />
          <Route path="/reading-challenge" element={<ReadingChallengePage />} />
        </Routes>
      </div>
    </div>
  );
}
