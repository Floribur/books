import { Routes, Route } from 'react-router-dom';
import { Sidebar } from './components/Sidebar';
import { HomePage } from './pages/HomePage';
import { BookDetailPage } from './pages/BookDetailPage';
import { AuthorsPage } from './pages/AuthorsPage';
import { AuthorDetailPage } from './pages/AuthorDetailPage';
import { GenresPage } from './pages/GenresPage';
import { GenreDetailPage } from './pages/GenreDetailPage';
import { ReadingChallengePage } from './pages/ReadingChallengePage';
import './App.css';

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
