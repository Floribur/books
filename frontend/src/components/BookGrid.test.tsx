import { describe, it } from 'vitest';

describe('BookGrid', () => {
  it.todo('renders fetched books as BookCard components');
  it.todo('shows 12 skeleton cards while isPending is true');
  it.todo('shows Load More button when hasNextPage is true');
  it.todo('hides Load More and shows end message when hasNextPage is false');
  it.todo('Load More button is disabled while isFetchingNextPage');
  it.todo('calls fetchNextPage when IntersectionObserver fires');
});
