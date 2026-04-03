import { Routes, Route } from 'react-router-dom';

// Page stubs — to be implemented in Plans 02 and 03
function HomePage() { return <main><p>Home page — coming soon</p></main>; }
function BookDetailPage() { return <main><p>Book detail — coming soon</p></main>; }
function AuthorsPage() { return <main><p>Authors — coming soon</p></main>; }
function AuthorDetailPage() { return <main><p>Author detail — coming soon</p></main>; }
function GenresPage() { return <main><p>Genres — coming soon</p></main>; }
function GenreDetailPage() { return <main><p>Genre detail — coming soon</p></main>; }
function ReadingChallengePage() { return <main><p>Reading Challenge — coming soon</p></main>; }

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/books/:slug" element={<BookDetailPage />} />
      <Route path="/authors" element={<AuthorsPage />} />
      <Route path="/authors/:slug" element={<AuthorDetailPage />} />
      <Route path="/genres" element={<GenresPage />} />
      <Route path="/genres/:slug" element={<GenreDetailPage />} />
      <Route path="/reading-challenge" element={<ReadingChallengePage />} />
    </Routes>
  );
}
