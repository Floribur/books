import { useEffect, useRef } from 'react';

// Attaches an IntersectionObserver to a sentinel element.
// callback is invoked when the sentinel enters the viewport (threshold: 0.1).
// enabled: pass false to disable the observer (e.g. when hasNextPage is false).
export function useIntersectionObserver(
  callback: () => void,
  enabled: boolean
): React.RefObject<HTMLDivElement | null> {
  const sentinelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!enabled) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          callback();
        }
      },
      { threshold: 0.1 }
    );

    if (sentinelRef.current) {
      observer.observe(sentinelRef.current);
    }

    return () => observer.disconnect();
  }, [enabled, callback]);

  return sentinelRef;
}
