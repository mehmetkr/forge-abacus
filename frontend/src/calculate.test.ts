import { calculate } from './calculate'

afterEach(() => {
  vi.restoreAllMocks()
})

test('sends correct POST request with JSON body and Content-Type header', async () => {
  const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue(
    new Response(JSON.stringify({ result: 5 }), { status: 200 }),
  )

  await calculate('add', 2, 3)

  expect(fetchSpy).toHaveBeenCalledOnce()
  const [url, options] = fetchSpy.mock.calls[0]
  expect(url).toBe('/api/calculate')
  expect(options?.method).toBe('POST')
  expect(options?.headers).toEqual({ 'Content-Type': 'application/json' })
  expect(JSON.parse(options?.body as string)).toEqual({ operation: 'add', a: 2, b: 3 })
})

test('returns result field on 200 response', async () => {
  vi.spyOn(globalThis, 'fetch').mockResolvedValue(
    new Response(JSON.stringify({ result: 42 }), { status: 200 }),
  )

  const result = await calculate('multiply', 7, 6)
  expect(result).toBe(42)
})

test('throws with backend error message on non-200 JSON response', async () => {
  vi.spyOn(globalThis, 'fetch').mockResolvedValue(
    new Response(JSON.stringify({ error: 'division by zero is undefined' }), { status: 400 }),
  )

  await expect(calculate('divide', 1, 0)).rejects.toThrow('division by zero is undefined')
})

test('throws "Could not reach the server" on network failure', async () => {
  vi.spyOn(globalThis, 'fetch').mockRejectedValue(new TypeError('Failed to fetch'))

  await expect(calculate('add', 1, 2)).rejects.toThrow('Could not reach the server')
})

test('throws "Could not reach the server" on non-JSON response', async () => {
  vi.spyOn(globalThis, 'fetch').mockResolvedValue(
    new Response('<html>502 Bad Gateway</html>', { status: 502 }),
  )

  await expect(calculate('add', 1, 2)).rejects.toThrow('Could not reach the server')
})

test('omits b from body when b is undefined', async () => {
  const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue(
    new Response(JSON.stringify({ result: 3 }), { status: 200 }),
  )

  await calculate('sqrt', 9)

  const body = JSON.parse(fetchSpy.mock.calls[0][1]?.body as string)
  expect(body).toEqual({ operation: 'sqrt', a: 9 })
  expect(body).not.toHaveProperty('b')
})
