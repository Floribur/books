import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, act } from '@testing-library/react';
import { LottieAnimation } from './LottieAnimation';

// Mock lottie-react — avoids canvas/animation runtime in tests
const mockLottieInstance = { stop: vi.fn(), play: vi.fn(), pause: vi.fn() };

vi.mock('lottie-react', () => ({
  default: vi.fn(({ lottieRef }: { lottieRef?: React.MutableRefObject<typeof mockLottieInstance> }) => {
    if (lottieRef) lottieRef.current = mockLottieInstance;
    return <div data-testid="lottie-animation" />;
  }),
}));

// matchMedia mock — simulate prefers-reduced-motion: reduce
function setReducedMotion(matches: boolean) {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation((query: string) => ({
      matches: matches && query === '(prefers-reduced-motion: reduce)',
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    })),
  });
}

describe('LottieAnimation', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockLottieInstance.stop.mockClear();
    mockLottieInstance.play.mockClear();
  });

  it('renders the lottie animation wrapper', () => {
    setReducedMotion(false);
    const { getByTestId } = render(<LottieAnimation />);
    expect(getByTestId('lottie-animation')).toBeTruthy();
  });

  it('has role=img and aria-label for accessibility', () => {
    setReducedMotion(false);
    const { getByRole } = render(<LottieAnimation />);
    // The wrapper div gets role="img" and aria-label
    expect(getByRole('img')).toBeTruthy();
  });

  it('stops animation on mount when prefers-reduced-motion is active', async () => {
    setReducedMotion(true);
    await act(async () => {
      render(<LottieAnimation />);
    });
    expect(mockLottieInstance.stop).toHaveBeenCalledTimes(1);
  });

  it('does not stop animation when prefers-reduced-motion is inactive', async () => {
    setReducedMotion(false);
    await act(async () => {
      render(<LottieAnimation />);
    });
    expect(mockLottieInstance.stop).not.toHaveBeenCalled();
  });

  it('registers a change listener on the matchMedia object', async () => {
    const addEventListenerMock = vi.fn();
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: vi.fn().mockReturnValue({
        matches: false,
        addEventListener: addEventListenerMock,
        removeEventListener: vi.fn(),
      }),
    });
    await act(async () => {
      render(<LottieAnimation />);
    });
    expect(addEventListenerMock).toHaveBeenCalledWith('change', expect.any(Function));
  });
});
