import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import NumberInput from './NumberInput'

test('renders a labeled input with the correct label text', () => {
  render(<NumberInput label="A" value="" onChange={() => {}} />)
  expect(screen.getByLabelText('A')).toBeInTheDocument()
})

test('calls onChange handler when user types', async () => {
  const handleChange = vi.fn()
  render(<NumberInput label="A" value="" onChange={handleChange} />)

  await userEvent.type(screen.getByLabelText('A'), '5')
  expect(handleChange).toHaveBeenCalledWith('5')
})

test('displays the current value', () => {
  render(<NumberInput label="A" value="42" onChange={() => {}} />)
  expect(screen.getByLabelText('A')).toHaveValue('42')
})

test('uses inputMode="decimal" for mobile keyboard', () => {
  render(<NumberInput label="A" value="" onChange={() => {}} />)
  expect(screen.getByLabelText('A')).toHaveAttribute('inputMode', 'decimal')
})
