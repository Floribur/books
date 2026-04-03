import { useEffect, useRef } from 'react';
import LottieImport, { type LottieRefCurrentProps } from 'lottie-react';
import animationData from '../assets/reading-animation.json';

// lottie-react ships as CJS; Vite's ESM interop may not unwrap .default automatically
const Lottie = (LottieImport as unknown as { default: typeof LottieImport }).default ?? LottieImport;

/**
 * Animated reader figure for the sidebar header.
 * - Plays on loop in default state
 * - Stops (shows frame 0) when prefers-reduced-motion: reduce is active
 * - Listens for mid-session OS setting changes
 *
 * IMPORTANT: Use lottieRef prop (not ref) — lottie-react custom prop.
 * CSS animation-play-state does NOT work on Lottie canvas — JS stop()/play() is required.
 */
export function LottieAnimation() {
  const lottieRef = useRef<LottieRefCurrentProps>(null);

  useEffect(() => {
    const mq = window.matchMedia('(prefers-reduced-motion: reduce)');

    // Apply immediately on mount
    if (mq.matches) {
      lottieRef.current?.stop();
    }

    // Listen for mid-session changes (e.g. user toggles OS accessibility settings)
    const handleChange = (e: MediaQueryListEvent) => {
      if (e.matches) {
        lottieRef.current?.stop();
      } else {
        lottieRef.current?.play();
      }
    };

    mq.addEventListener('change', handleChange);
    return () => mq.removeEventListener('change', handleChange);
  }, []);

  return (
    <div
      className="lottie-animation-wrapper"
      role="img"
      aria-label="Animated reading figure — Flo's Library logo"
    >
      <Lottie
        lottieRef={lottieRef}
        animationData={animationData}
        loop={true}
        style={{ width: '100%', maxWidth: 200, margin: '0 auto', display: 'block' }}
      />
    </div>
  );
}
