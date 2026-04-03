import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchBookBySlug } from '../api/books';
import { BookCover } from '../components/BookCover';
import { DescriptionBlock } from '../components/DescriptionBlock';
import { Toast } from '../components/Toast';
import { useState, useEffect } from 'react';
import './BookDetailPage.css';

// Format read_at ISO string as "Read Month D, YYYY" — D-03 metadata order
function formatReadDate(readAt: string | null): string | null {
  if (!readAt) return null;
  const date = new Date(readAt);
  return `Read ${date.toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })}`;
}

export function BookDetailPage() {
  const { slug } = useParams<{ slug: string }>();
  const [showError, setShowError] = useState(false);

  const { data: book, isPending, isError } = useQuery({
    queryKey: ['book', slug],
    queryFn: () => fetchBookBySlug(slug!),
    enabled: !!slug,
  });

  useEffect(() => {
    if (isError) setShowError(true);
  }, [isError]);

  // Build metadata line: "Published {year} · {N} pages · Read {date}" — omit null fields
  // Separator: · (U+00B7 middle dot) per copywriting contract
  const metaParts: string[] = [];
  if (book?.publication_year) metaParts.push(`Published ${book.publication_year}`);
  if (book?.page_count) metaParts.push(`${book.page_count} pages`);
  const readDate = formatReadDate(book?.read_at ?? null);
  if (readDate) metaParts.push(readDate);
  const metaLine = metaParts.join(' \u00B7 ');

  return (
    <main className="book-detail-page">
      {showError && (
        <Toast
          message="Couldn't load book details. Try refreshing the page."
          onDismiss={() => setShowError(false)}
        />
      )}

      {isPending ? (
        // Loading skeleton — two-column structure
        <div className="book-detail-layout">
          <div className="book-detail-cover-col book-detail-cover--skeleton" aria-hidden="true" />
          <div className="book-detail-meta-col">
            <div className="book-detail-skeleton-line" style={{ width: '70%', height: '28px' }} />
            <div className="book-detail-skeleton-line" style={{ width: '40%', height: '14px' }} />
            <div className="book-detail-skeleton-line" style={{ width: '60%', height: '14px' }} />
            <div className="book-detail-skeleton-line" style={{ width: '90%', height: '14px' }} />
          </div>
        </div>
      ) : book ? (
        <>
          {/* Desktop: CSS grid two-column. Mobile: cover floats left 96px (D-01, D-02) */}
          <div className="book-detail-layout">
            <div className="book-detail-cover-col">
              <BookCover
                src={book.cover_path}
                title={book.title}
                loading="eager"
              />
            </div>

            <div className="book-detail-meta-col">
              {/* D-03 metadata order: Title → Author links → Genres → meta line → Description */}

              {/* 1. Title */}
              <h1 className="book-detail-title">{book.title}</h1>

              {/* 2. Author links — comma-separated, link to /authors/:slug (BOOK-02) */}
              {book.authors.length > 0 && (
                <p className="book-detail-authors">
                  {book.authors.map((author, i) => (
                    <span key={author.slug}>
                      <Link to={`/authors/${author.slug}`} className="book-detail-author-link">
                        {author.name}
                      </Link>
                      {i < book.authors.length - 1 ? ', ' : ''}
                    </span>
                  ))}
                </p>
              )}

              {/* 3. Genre tags — pill links to /genres/:slug (BOOK-03) */}
              {book.genres.length > 0 && (
                <div className="book-detail-genres">
                  {book.genres.map((genre) => (
                    <Link
                      key={genre.slug}
                      to={`/genres/${genre.slug}`}
                      className="genre-tag"
                    >
                      {genre.name}
                    </Link>
                  ))}
                </div>
              )}

              {/* 4. Metadata line — BOOK-04, BOOK-01 */}
              {metaLine && (
                <p className="book-detail-meta">{metaLine}</p>
              )}

              {/* 5. Description with expand/collapse — D-04 */}
              <DescriptionBlock description={book.description} />
            </div>

            {/* Clearfix for mobile float layout — D-02 */}
            <div style={{ clear: 'both' }} />
          </div>
        </>
      ) : null}
    </main>
  );
}
