import { useEffect, useId, useRef, useState } from 'react'
import { createPortal } from 'react-dom'
import './SelectField.css'
export type SelectOption = { value: string; label: string; disabled?: boolean }
// One listbox implementation for all systems. The popup is portalled out of scrolling cards/dialogs.
export function SelectField({
  label,
  value,
  onChange,
  options,
  disabled = false,
  hideLabel = false,
}: {
  label: string
  value: string
  onChange: (value: string) => void
  options: SelectOption[]
  disabled?: boolean
  hideLabel?: boolean
}) {
  const id = useId(),
    trigger = useRef<HTMLButtonElement>(null),
    popup = useRef<HTMLDivElement>(null)
  const [portal, setPortal] = useState<HTMLElement | null>(null)
  const [open, setOpen] = useState(false),
    [active, setActive] = useState(0),
    [position, setPosition] = useState({ left: 0, top: 0, width: 0, maxHeight: 280 })
  const selected = options.find((o) => o.value === value)
  const typeahead = useRef({ text: '', time: 0 })
  const close = () => {
    setOpen(false)
    trigger.current?.focus()
  }
  function show() {
    if (disabled) return
    setPortal(trigger.current?.closest('dialog') ?? document.body)
    setActive(
      Math.max(
        0,
        options.findIndex((o) => o.value === value),
      ),
    )
    setOpen(true)
  }
  useEffect(() => {
    if (!open) return
    const place = () => {
      const r = trigger.current?.getBoundingClientRect()
      if (!r) return
      const below = window.innerHeight - r.bottom - 12,
        above = r.top - 12
      const up = below < 180 && above > below
      const height = Math.min(320, up ? above : below)
      setPosition({
        left: Math.max(8, Math.min(r.left, window.innerWidth - r.width - 8)),
        top: up ? Math.max(8, r.top - height - 4) : r.bottom + 4,
        width: r.width,
        maxHeight: Math.max(80, height),
      })
    }
    const dismiss = (e: PointerEvent) => {
      if (
        !trigger.current?.contains(e.target as Node) &&
        !popup.current?.contains(e.target as Node)
      )
        setOpen(false)
    }
    place()
    window.addEventListener('resize', place)
    window.addEventListener('scroll', place, true)
    document.addEventListener('pointerdown', dismiss)
    popup.current?.focus()
    return () => {
      window.removeEventListener('resize', place)
      window.removeEventListener('scroll', place, true)
      document.removeEventListener('pointerdown', dismiss)
    }
  }, [open])
  useEffect(() => {
    if (open) document.getElementById(`${id}-${active}`)?.scrollIntoView({ block: 'nearest' })
  }, [active, open, id])
  function move(delta: number) {
    for (let n = 1; n <= options.length; n++) {
      const i = (active + delta * n + options.length) % options.length
      if (!options[i].disabled) {
        setActive(i)
        return
      }
    }
  }
  return (
    <div className="select-field">
      <span id={`${id}-label`} className={hideLabel ? 'sr-only' : 'select-field-label'}>
        {label}
      </span>
      <button
        type="button"
        ref={trigger}
        className="select-field-trigger"
        role="combobox"
        aria-labelledby={`${id}-label`}
        aria-label={label}
        aria-expanded={open}
        aria-controls={open ? `${id}-list` : undefined}
        aria-haspopup="listbox"
        disabled={disabled}
        onClick={() => (open ? close() : show())}
        onKeyDown={(e) => {
          if (['ArrowDown', 'ArrowUp', 'Enter', ' '].includes(e.key)) {
            e.preventDefault()
            show()
          }
        }}
      >
        <span id={`${id}-value`}>{selected?.label ?? ''}</span>
        <svg aria-hidden="true" width="16" height="16" viewBox="0 0 16 16">
          <path d="m4 6 4 4 4-4" fill="none" stroke="currentColor" strokeWidth="1.6" />
        </svg>
      </button>
      {open &&
        portal &&
        createPortal(
          <div
            ref={popup}
            id={`${id}-list`}
            className="select-field-popup"
            role="listbox"
            aria-label={label}
            aria-activedescendant={`${id}-${active}`}
            tabIndex={-1}
            style={{
              left: position.left,
              top: position.top,
              width: position.width,
              maxHeight: position.maxHeight,
            }}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                e.preventDefault()
                e.stopPropagation()
                close()
              } else if (e.key === 'Tab') {
                setOpen(false)
                trigger.current?.focus()
              } else if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
                e.preventDefault()
                move(e.key === 'ArrowDown' ? 1 : -1)
              } else if (e.key === 'Home' || e.key === 'End') {
                e.preventDefault()
                const indexes = options.map((o, i) => (o.disabled ? -1 : i)).filter((i) => i >= 0)
                setActive(e.key === 'Home' ? indexes[0] : (indexes.at(-1) ?? 0))
              } else if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault()
                const option = options[active]
                if (option && !option.disabled) {
                  onChange(option.value)
                  close()
                }
              } else if (e.key.length === 1) {
                const now = Date.now()
                typeahead.current = {
                  text:
                    (now - typeahead.current.time < 700 ? typeahead.current.text : '') +
                    e.key.toLocaleLowerCase(),
                  time: now,
                }
                const i = options.findIndex(
                  (o) =>
                    !o.disabled && o.label.toLocaleLowerCase().startsWith(typeahead.current.text),
                )
                if (i >= 0) setActive(i)
              }
            }}
          >
            {options.map((option, i) => (
              <div
                key={option.value}
                id={`${id}-${i}`}
                role="option"
                aria-selected={option.value === value}
                aria-disabled={option.disabled || undefined}
                data-active={active === i}
                className="select-field-option"
                onPointerMove={() => {
                  if (!option.disabled) setActive(i)
                }}
                onPointerDown={(e) => e.preventDefault()}
                onClick={() => {
                  if (!option.disabled) {
                    onChange(option.value)
                    close()
                  }
                }}
              >
                <span>{option.label}</span>
                <span aria-hidden="true">{option.value === value ? '✓' : ''}</span>
              </div>
            ))}
          </div>,
          portal,
        )}
    </div>
  )
}
