export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percentage'

export async function calculate(operation: Operation, a: number, b?: number): Promise<number> {
  const body: Record<string, unknown> = { operation, a }
  if (b !== undefined) {
    body.b = b
  }

  let response: Response
  try {
    response = await fetch('/api/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
  } catch {
    throw new Error('Could not reach the server')
  }

  let data: Record<string, unknown>
  try {
    data = await response.json()
  } catch {
    throw new Error('Could not reach the server')
  }

  if (!response.ok) {
    throw new Error(data.error as string)
  }

  return data.result as number
}
