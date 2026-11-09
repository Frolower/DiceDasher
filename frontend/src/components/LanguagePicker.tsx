import { useEffect, useId, useRef, useState, useSyncExternalStore } from 'react'
import { getLanguage, setLanguage, subscribeLanguage, t } from '../lib/i18n'

export function LanguagePicker() {
  const language = useSyncExternalStore(subscribeLanguage, getLanguage, getLanguage)
  const [unsaved, setUnsaved] = useState(false)
  const [open, setOpen] = useState(false)
  const root = useRef<HTMLDivElement>(null)
  const trigger = useRef<HTMLButtonElement>(null)
  const optionsId = useId()
  const warningId = useId()
  useEffect(() => {
    if (!open) return
    function outside(event: PointerEvent) {
      if (!root.current?.contains(event.target as Node)) setOpen(false)
    }
    document.addEventListener('pointerdown', outside)
    return () => document.removeEventListener('pointerdown', outside)
  }, [open])
  return (
    <div
      className="language-picker"
      ref={root}
      onBlur={(event) => {
        if (!event.currentTarget.contains(event.relatedTarget)) setOpen(false)
      }}
      onKeyDown={(event) => {
        if (event.key === 'Escape' && open) {
          event.preventDefault()
          event.stopPropagation()
          setOpen(false)
          trigger.current?.focus()
        }
      }}
    >
      <button
        ref={trigger}
        type="button"
        className="language-trigger"
        aria-label={`${t('language.label')}: ${language === 'en' ? 'English' : 'Русский'}`}
        aria-expanded={open}
        aria-controls={optionsId}
        aria-describedby={unsaved ? warningId : undefined}
        onClick={() => setOpen((value) => !value)}
      >
        <span lang={language}>{language === 'en' ? 'English' : 'Русский'}</span>
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
          <path
            d="m3 4.5 3 3 3-3"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
      </button>
      {open && (
        <div
          id={optionsId}
          className="language-options"
          role="group"
          aria-label={t('language.label')}
        >
          {(['en', 'ru'] as const).map((option) => (
            <button
              key={option}
              type="button"
              className="language-option"
              lang={option}
              aria-pressed={language === option}
              onClick={() => {
                setUnsaved(!setLanguage(option))
                setOpen(false)
                trigger.current?.focus()
              }}
            >
              {option === 'en' ? 'English' : 'Русский'}
              <span aria-hidden="true">{language === option ? '✓' : ''}</span>
            </button>
          ))}
        </div>
      )}
      {unsaved && (
        <span id={warningId} role="status" className="language-warning">
          {t('language.unsaved')}
        </span>
      )}
    </div>
  )
}
