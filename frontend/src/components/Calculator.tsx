import { useState } from 'react'
import { calculate, type Operation } from '../calculate'
import NumberInput from './NumberInput'
import styles from './Calculator.module.css'

const OPERATIONS: { value: Operation; label: string }[] = [
  { value: 'add', label: 'Add' },
  { value: 'subtract', label: 'Subtract' },
  { value: 'multiply', label: 'Multiply' },
  { value: 'divide', label: 'Divide' },
  { value: 'power', label: 'Power' },
  { value: 'sqrt', label: 'Square Root' },
  { value: 'percentage', label: 'Percentage' },
]

function Calculator() {
  const [operation, setOperation] = useState<Operation>('add')
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const needsB = operation !== 'sqrt'

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setResult(null)
    setError(null)

    if (a.trim() === '') {
      setError('Enter a value for A')
      return
    }
    if (needsB && b.trim() === '') {
      setError('Enter a value for B')
      return
    }

    const numA = Number(a)
    if (!Number.isFinite(numA)) {
      setError('A must be a valid number')
      return
    }

    let numB: number | undefined
    if (needsB) {
      numB = Number(b)
      if (!Number.isFinite(numB)) {
        setError('B must be a valid number')
        return
      }
    }

    setLoading(true)
    try {
      const value = await calculate(operation, numA, numB)
      setResult(value)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An unexpected error occurred')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className={styles.calculator} onSubmit={handleSubmit}>
      <div className={styles.field}>
        <label htmlFor="operation">Operation</label>
        <select
          id="operation"
          value={operation}
          onChange={(e) => setOperation(e.target.value as Operation)}
        >
          {OPERATIONS.map((op) => (
            <option key={op.value} value={op.value}>
              {op.label}
            </option>
          ))}
        </select>
      </div>

      <NumberInput label="A" value={a} onChange={setA} />
      {needsB && <NumberInput label="B" value={b} onChange={setB} />}

      <button type="submit" disabled={loading}>
        {loading ? 'Calculating…' : 'Calculate'}
      </button>

      <div className={styles.output} aria-live="polite">
        {result !== null && <p className={styles.result}>{result}</p>}
        {error !== null && <p className={styles.error}>{error}</p>}
      </div>
    </form>
  )
}

export default Calculator
