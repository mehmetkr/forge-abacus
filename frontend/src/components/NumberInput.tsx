import styles from './NumberInput.module.css'

interface NumberInputProps {
  label: string
  value: string
  onChange: (value: string) => void
}

function NumberInput({ label, value, onChange }: NumberInputProps) {
  return (
    <div className={styles.field}>
      <label htmlFor={`input-${label}`}>{label}</label>
      <input
        id={`input-${label}`}
        type="text"
        inputMode="decimal"
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  )
}

export default NumberInput
