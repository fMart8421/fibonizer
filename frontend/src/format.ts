const UNITS: [string, number][] = [
  ['s', 1e9],
  ['ms', 1e6],
  ['µs', 1e3],
]

export function formatDuration(ns: number): string {
  for (const [unit, size] of UNITS) {
    if (ns >= size) {
      return `${(ns / size).toFixed(2)} ${unit}`
    }
  }

  return `${ns} ns`
}

// Groups digits so big Fibonacci numbers stay readable: 832040 -> 832,040
export function formatBigNumber(value: string): string {
  return value.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}
