import { Children, isValidElement, useState, type ReactNode, type ReactElement } from 'react'
import { SelectField, type SelectOption } from '../../components/SelectField'
import { l, copy, formatMessage, type Message } from './localization'
export function Field({
  label,
  value,
  onChange,
  multiline = false,
  readOnly = false,
  placeholder,
}: {
  label: string
  value: string
  onChange: (value: string) => void
  multiline?: boolean
  readOnly?: boolean
  placeholder?: string
}) {
  return (
    <label className="vtm-field">
      <span>{l(label)}</span>
      {multiline ? (
        <textarea
          value={value}
          onChange={(e) => onChange(e.target.value)}
          readOnly={readOnly}
          placeholder={placeholder ? l(placeholder) : undefined}
        />
      ) : (
        <input
          value={value}
          onChange={(e) => onChange(e.target.value)}
          readOnly={readOnly}
          placeholder={placeholder ? l(placeholder) : undefined}
        />
      )}
    </label>
  )
}
export function Select({
  label,
  value,
  onChange,
  children,
  disabled = false,
  readOnly = false,
}: {
  label: string
  value: string
  onChange: (value: string) => void
  children: ReactNode
  disabled?: boolean
  readOnly?: boolean
}) {
  const text = (node: ReactNode): string =>
    Children.toArray(node)
      .map((c) =>
        isValidElement(c) ? text((c.props as { children?: ReactNode }).children) : l(String(c)),
      )
      .join('')
  const options: SelectOption[] = []
  Children.forEach(children, (child) => {
    if (isValidElement(child)) {
      const p = (
        child as ReactElement<{ value?: string | number; children?: ReactNode; disabled?: boolean }>
      ).props
      const label = text(p.children)
      options.push({ value: String(p.value ?? label), label, disabled: p.disabled })
    }
  })
  if (readOnly)
    return (
      <div className="vtm-field vtm-static-field">
        <span>{l(label)}</span>
        <div className="vtm-static-value">
          {options.find((o) => o.value === value)?.label || '—'}
        </div>
      </div>
    )
  return (
    <SelectField
      label={l(label)}
      value={value}
      onChange={onChange}
      options={options}
      disabled={disabled}
    />
  )
}
export function NumberField({
  label,
  value,
  onChange,
  min = 0,
  max,
  fallback = min,
  readOnly = false,
}: {
  label: string
  value: number
  onChange: (value: number) => void
  min?: number
  max?: number
  fallback?: number
  readOnly?: boolean
}) {
  const [text, setText] = useState<string | null>(null)
  const clamp = (n: number) =>
    Math.max(min, Math.min(max ?? Number.MAX_SAFE_INTEGER, Math.trunc(n)))
  return (
    <label className="vtm-field vtm-number">
      <span>{l(label)}</span>
      <input
        aria-label={l(label)}
        type="number"
        inputMode="numeric"
        step="1"
        value={text ?? String(value)}
        onFocus={() => setText(String(value))}
        min={min}
        max={max}
        readOnly={readOnly}
        onChange={(e) => {
          setText(e.target.value)
          const n = e.target.valueAsNumber
          if (Number.isFinite(n)) onChange(clamp(n))
        }}
        onBlur={() => {
          const n = clamp(text === null ? value : text.trim() === '' ? fallback : Number(text))
          setText(null)
          onChange(n)
        }}
      />
    </label>
  )
}
export function Rating({
  label,
  value,
  onChange,
  max = 5,
  min = 0,
  readOnly = false,
  canChange = () => true,
  clearable = false,
}: {
  label: string
  value: number
  onChange: (n: number) => void
  max?: number
  min?: number
  readOnly?: boolean
  canChange?: (n: number) => boolean
  clearable?: boolean
}) {
  return (
    <div className="vtm-rating">
      <span>{l(label)}</span>
      <div className="vtm-dots" role="group" aria-label={`${l(label)}: ${value}`}>
        {Array.from({ length: max }, (_, i) => (
          <button
            type="button"
            key={i}
            className={i > 0 && i % 5 === 0 ? 'vtm-five-gap' : undefined}
            aria-label={`${l(label)}: ${i + 1}`}
            aria-pressed={value >= i + 1}
            disabled={readOnly || !canChange(value === i + 1 ? Math.max(min, i) : i + 1)}
            onClick={() => {
              const n = value === i + 1 ? Math.max(min, i) : i + 1
              if (canChange(n)) onChange(n)
            }}
          >
            <span aria-hidden="true" className={value >= i + 1 ? 'filled' : ''} />
          </button>
        ))}
        {clearable && (
          <button
            type="button"
            className="vtm-clear-rating"
            aria-label={copy('Очистить: {0}', { '0': l(label) })}
            disabled={value === 0 || readOnly}
            onClick={() => onChange(0)}
          >
            ×
          </button>
        )}
      </div>
    </div>
  )
}
export function Errors({ errors }: { errors: Message[] }) {
  return errors.length ? (
    <div role="alert" className="vtm-errors">
      <strong>{l('Проверьте заполнение')}</strong>
      <ul>
        {errors.map((e, i) => (
          <li key={i}>{formatMessage(e)}</li>
        ))}
      </ul>
    </div>
  ) : null
}
