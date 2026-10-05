import { useState, type FormEvent } from 'react'
import { fibonize, METHODS, type FiboResponse } from './api'
import { formatBigNumber, formatDuration } from './format'
import './App.css'

function App() {
  const [input, setInput] = useState('10')
  const [results, setResults] = useState<FiboResponse[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    setLoading(true)
    setError(null)
    setResults([])

    try {
      const n = Number(input)
      const responses: FiboResponse[] = []

      // Sequential on purpose: running them in parallel would make the
      // implementations compete for CPU and skew the measured times
      for (const method of METHODS) {
        responses.push(await fibonize(method, n))
        setResults([...responses])
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
    } finally {
      setLoading(false)
    }
  }

  const slowest = Math.max(1, ...results.map((r) => r.durationNs))
  const fastest = Math.min(...results.map((r) => r.durationNs))

  return (
    <main className="app">
      <header>
        <h1>Fibonizer</h1>
        <p className="subtitle">
          Calculate F(N) and compare how long each implementation takes.
        </p>
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
          {loading ? 'Fibonizing…' : 'Fibonize'}
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
            F({results[0].n}) = <span>{formatBigNumber(results[0].result)}</span>
          </p>
        ) : (
          <p className="placeholder">
            {loading ? 'Calculating…' : 'Enter a number and press Fibonize.'}
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
              const r = results.find((res) => res.method === method)
              const isFastest = r && results.length > 1 && r.durationNs === fastest

              return (
                <tr key={method} className={isFastest ? 'fastest' : undefined}>
                  <td>{method}</td>
                  <td className="mono">
                    {r ? formatBigNumber(r.result) : loading ? '…' : '—'}
                  </td>
                  <td className="mono">{r ? formatDuration(r.durationNs) : '—'}</td>
                  <td className="bar-col">
                    <div className="bar-track">
                      <div
                        className="bar"
                        style={{ width: r ? `${(r.durationNs / slowest) * 100}%` : 0 }}
                      />
                    </div>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </section>
    </main>
  )
}

export default App
