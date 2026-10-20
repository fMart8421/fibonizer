const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export type FiboMethod = "loop/v1" | "loop/v2" | "recursive/v1" | "recursive/v2" | "binet";

export const METHODS: FiboMethod[] = [
  "loop/v1",
  "loop/v2",
  "recursive/v1",
  "recursive/v2",
  "binet",
];

export interface FiboResponse {
  method: FiboMethod;
  n: number;
  // String because int64 results above 2^53 lose precision as JS numbers
  result: string;
  durationNs: number;
}

export async function fibonize(method: FiboMethod, n: number): Promise<FiboResponse> {
  const res = await fetch(`${API_URL}/${method}/${n}`, {
    signal: AbortSignal.timeout(5000), // timeout after 5 seconds of request
  });

  if (!res.ok) {
    throw new Error(`${method} request failed with status ${res.status}`);
  }

  return res.json();
}
