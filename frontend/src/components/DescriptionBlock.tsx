import { useState } from 'react';
import './DescriptionBlock.css';

interface DescriptionBlockProps {
  description: string | null;
}

const CHAR_THRESHOLD = 640; // D-04: expand trigger threshold

export function DescriptionBlock({ description }: DescriptionBlockProps) {
  const [isExpanded, setIsExpanded] = useState(false);

  if (!description) return null;

  const isLong = description.length > CHAR_THRESHOLD;

  return (
    <div className="description-block">
      <div
        className={`description-block-text${isLong && !isExpanded ? ' description-block-text--clamped' : ''}`}
      >
        {description}
      </div>
      {isLong && (
        <button
          className="description-block-toggle"
          onClick={() => setIsExpanded((prev) => !prev)}
          aria-expanded={isExpanded}
        >
          {isExpanded ? 'Show less' : 'Show more'}
        </button>
      )}
    </div>
  );
}
