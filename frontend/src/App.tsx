import { useState, type FormEvent } from "react";
import { fibonize, METHODS, type FiboResponse } from "./api";
import { formatBigNumber, formatDuration } from "./format";
import "./App.css";

function App() {
  const [input, setInput] = useState("10");
  const [results, setResults] = useState<FiboResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setLoading(true);
    setError(null);
    setResults([]);

    try {
      const n = Number(input);

      // All requests start at once and each one updates the UI as soon as it
      // arrives, so a slow method (e.g. recursive/v1) never holds up the rest.
      // Each request handles its own failure, so one timeout doesn't hide the others
      await Promise.all(
        METHODS.map(async (method) => {
          try {
            const fiboValue = await fibonize(method, n);
            setResults((prev) => [...prev, fiboValue]);
          } catch (err) {
            const message = err instanceof Error ? err.message : String(err);
            setError((prev) => (prev ? `${prev}; ${method}: ${message}` : `${method}: ${message}`));
          }
        }),
      );
    } finally {
      setLoading(false);
    }
  }


  const valid = results.filter((r) => r.result !== '-1');

  const fastest = valid.reduce<FiboResponse | null>(
    (best, r) => (!best || r.durationNs < best.durationNs ? r : best),
    null,
  );

  const slowest = valid.reduce<FiboResponse | null>(
    (worst, r) => (!worst || r.durationNs > worst.durationNs ? r : worst),
    null,
  );

  return (
    <main className="app">
      <header>
        <h1>Fibonizer</h1>
        <p className="subtitle">Calculate F(N) and compare how long each implementation takes.</p>
      </header>

      <form className="fibo-form" onSubmit={handleSubmit}>
        <label htmlFor="n">N</label>
        <input
          id="n"
          type="number"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          required
        />
        <button type="submit" disabled={loading}>
          {loading ? "Fibonizing…" : "Fibonize"}
        </button>
      </form>

      {error && (
        <p className="error" role="alert">
          {error}
        </p>
      )}

      <section className="result" aria-live="polite">
        <h2>Result</h2>
        {results.length > 0 ? (
          <p className="result-value">
            F({fastest?.n}) = <span>{formatBigNumber(fastest?.result ?? '—')}</span>
          </p>
        ) : (
          <p className="placeholder">
            {loading ? "Calculating…" : "Enter a number and press Fibonize."}
          </p>
        )}
      </section>

      <section className="comparison">
        <h2>Time comparison</h2>
        <table>
          <thead>
            <tr>
              <th>Method</th>
              <th>Result</th>
              <th>Time</th>
              <th className="bar-col">Relative</th>
            </tr>
          </thead>
          <tbody>
            {METHODS.map((method) => {
              const r = results.find((res) => res.method === method);
              const isValid = valid.find((res) => res.method === method);

              const overloaded = r?.result === "-1";
              const resultText = overloaded
                ? "Result Overloaded"
                : formatBigNumber(r?.result ?? "");

              const isFastest = r && results.length > 1 && r.durationNs === fastest?.durationNs;

              return (
                <tr key={method} className={isFastest ? "fastest" : undefined}>
                  <td>{method}</td>
                  <td className="mono">{r ? resultText : loading ? "…" : "—"}</td>
                  <td className="mono">{r && !overloaded ? formatDuration(r.durationNs) : "—"}</td>
                  <td className="bar-col">
                    <div className="bar-track">
                      <div
                        className="bar"
                        style={{ width: r && isValid ? `${(r.durationNs / (slowest?.durationNs ?? r.durationNs)) * 100}%` : 0 }}
                      />
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </section>
    </main>
  );
}

export default App;
