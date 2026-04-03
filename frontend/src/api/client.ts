// Base fetch wrapper — throws Error with message on non-OK responses.
// All API callers (books.ts etc.) use this instead of raw fetch.

export class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, init);

  if (!response.ok) {
    let message = `API error ${response.status}`;
    try {
      const body = await response.json() as { error?: string };
      if (body.error) message = body.error;
    } catch {
      // body not JSON — use default message
    }
    throw new ApiError(message, response.status);
  }

  return response.json() as Promise<T>;
}
