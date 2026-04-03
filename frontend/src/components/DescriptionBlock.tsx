import { useState, useRef, useEffect } from 'react';
import './DescriptionBlock.css';

interface DescriptionBlockProps {
  description: string | null;
}

export function DescriptionBlock({ description }: DescriptionBlockProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const [isClamped, setIsClamped] = useState(false);
  const textRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const el = textRef.current;
    if (!el) return;

    const check = () => setIsClamped(el.scrollHeight > el.clientHeight);
    check();

    const observer = new ResizeObserver(check);
    observer.observe(el);
    return () => observer.disconnect();
  }, [description]);

  if (!description) return null;

  return (
    <div className="description-block">
      <div
        ref={textRef}
        className={`description-block-text${!isExpanded ? ' description-block-text--clamped' : ''}`}
      >
        {description}
      </div>
      {(isClamped || isExpanded) && (
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
