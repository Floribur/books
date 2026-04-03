import { Routes, Route } from 'react-router-dom';
import { Sidebar } from './components/Sidebar';
import { HomePage } from './pages/HomePage';
import { BookDetailPage } from './pages/BookDetailPage';
import { AuthorsPage } from './pages/AuthorsPage';
import { AuthorDetailPage } from './pages/AuthorDetailPage';
import { GenresPage } from './pages/GenresPage';
import { GenreDetailPage } from './pages/GenreDetailPage';
import './App.css';

// ReadingChallengePage stub — will be replaced in Plan 4.2
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
