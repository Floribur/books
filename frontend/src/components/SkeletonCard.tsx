import './SkeletonCard.css';

export function SkeletonCard() {
  return (
    <div className="skeleton-card" aria-hidden="true">
      <div className="skeleton-cover" />
      <div className="skeleton-text skeleton-title" />
      <div className="skeleton-text skeleton-author" />
    </div>
  );
}
