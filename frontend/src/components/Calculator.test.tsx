import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { calculate } from '../calculate'
import Calculator from './Calculator'

vi.mock('../calculate')
const mockCalculate = vi.mocked(calculate)

afterEach(() => {
  vi.restoreAllMocks()
})

test('renders operation dropdown with all 7 options, default selected is "add"', () => {
  render(<Calculator />)
  const select = screen.getByRole('combobox', { name: /operation/i })
  expect(select).toHaveValue('add')

  const options = within(select).getAllByRole('option')
  expect(options).toHaveLength(7)
  expect(options.map((o) => o.textContent)).toEqual([
    'Add',
    'Subtract',
    'Multiply',
    'Divide',
    'Power',
    'Square Root',
    'Percentage',
  ])
})

test('renders two number inputs (A and B)', () => {
  render(<Calculator />)
  expect(screen.getByLabelText('A')).toBeInTheDocument()
  expect(screen.getByLabelText('B')).toBeInTheDocument()
})

test('renders a Calculate button', () => {
  render(<Calculator />)
  expect(screen.getByRole('button', { name: /calculate/i })).toBeInTheDocument()
})

test('submits a calculation and displays the result', async () => {
  mockCalculate.mockResolvedValue(5)
  render(<Calculator />)

  await userEvent.type(screen.getByLabelText('A'), '2')
  await userEvent.type(screen.getByLabelText('B'), '3')
  await userEvent.click(screen.getByRole('button', { name: /calculate/i }))

  expect(await screen.findByText('5')).toBeInTheDocument()
  expect(mockCalculate).toHaveBeenCalledWith('add', 2, 3)
})

test('displays backend error messages', async () => {
  mockCalculate.mockRejectedValue(new Error('division by zero is undefined'))
  render(<Calculator />)

  await userEvent.type(screen.getByLabelText('A'), '1')
  await userEvent.type(screen.getByLabelText('B'), '0')
  await userEvent.selectOptions(screen.getByRole('combobox', { name: /operation/i }), 'divide')
  await userEvent.click(screen.getByRole('button', { name: /calculate/i }))

  expect(await screen.findByText('division by zero is undefined')).toBeInTheDocument()
})

test('validates empty input A — shows "Enter a value for A"', async () => {
  render(<Calculator />)

  await userEvent.type(screen.getByLabelText('B'), '3')
  await userEvent.click(screen.getByRole('button', { name: /calculate/i }))

  expect(screen.getByText('Enter a value for A')).toBeInTheDocument()
  expect(mockCalculate).not.toHaveBeenCalled()
})

test('validates empty input B (non-sqrt) — shows "Enter a value for B"', async () => {
  render(<Calculator />)

  await userEvent.type(screen.getByLabelText('A'), '5')
  await userEvent.click(screen.getByRole('button', { name: /calculate/i }))

  expect(screen.getByText('Enter a value for B')).toBeInTheDocument()
  expect(mockCalculate).not.toHaveBeenCalled()
})

test('validates non-numeric input A — shows "A must be a valid number"', async () => {
  render(<Calculator />)

  await userEvent.type(screen.getByLabelText('A'), 'abc')
  await userEvent.type(screen.getByLabelText('B'), '3')
  await userEvent.click(screen.getByRole('button', { name: /calculate/i }))

  expect(screen.getByText('A must be a valid number')).toBeInTheDocument()
  expect(mockCalculate).not.toHaveBeenCalled()
})

test('validates non-numeric input B — shows "B must be a valid number"', async () => {
  render(<Calculator />)

  await userEvent.type(screen.getByLabelText('A'), '5')
  await userEvent.type(screen.getByLabelText('B'), 'xyz')
  await userEvent.click(screen.getByRole('button', { name: /calculate/i }))

  expect(screen.getByText('B must be a valid number')).toBeInTheDocument()
  expect(mockCalculate).not.toHaveBeenCalled()
})

test('hides operand B when sqrt is selected', async () => {
  render(<Calculator />)

  expect(screen.getByLabelText('B')).toBeInTheDocument()
  await userEvent.selectOptions(screen.getByRole('combobox', { name: /operation/i }), 'sqrt')
  expect(screen.queryByLabelText('B')).not.toBeInTheDocument()
})

test('shows loading state during request (button disabled)', async () => {
  let resolve: (value: number) => void
  mockCalculate.mockReturnValue(new Promise((r) => { resolve = r }))
  render(<Calculator />)

  await userEvent.type(screen.getByLabelText('A'), '2')
  await userEvent.type(screen.getByLabelText('B'), '3')
  await userEvent.click(screen.getByRole('button', { name: /calculate/i }))

  expect(screen.getByRole('button', { name: /calculating/i })).toBeDisabled()

  resolve!(5)
  expect(await screen.findByText('5')).toBeInTheDocument()
  expect(screen.getByRole('button', { name: /calculate/i })).toBeEnabled()
})

test('clears previous result/error on new submission', async () => {
  mockCalculate.mockResolvedValueOnce(5)
  render(<Calculator />)

  await userEvent.type(screen.getByLabelText('A'), '2')
  await userEvent.type(screen.getByLabelText('B'), '3')
  await userEvent.click(screen.getByRole('button', { name: /calculate/i }))
  expect(await screen.findByText('5')).toBeInTheDocument()

  mockCalculate.mockRejectedValueOnce(new Error('division by zero is undefined'))
  await userEvent.clear(screen.getByLabelText('A'))
  await userEvent.type(screen.getByLabelText('A'), '1')
  await userEvent.clear(screen.getByLabelText('B'))
  await userEvent.type(screen.getByLabelText('B'), '0')
  await userEvent.selectOptions(screen.getByRole('combobox', { name: /operation/i }), 'divide')
  await userEvent.click(screen.getByRole('button', { name: /calculate/i }))

  expect(await screen.findByText('division by zero is undefined')).toBeInTheDocument()
  expect(screen.queryByText('5')).not.toBeInTheDocument()
})
