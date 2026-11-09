import { useState, useSyncExternalStore } from 'react'
import { getLanguage, setLanguage, subscribeLanguage, t } from '../lib/i18n'

export function LanguagePicker() {
  const language = useSyncExternalStore(subscribeLanguage, getLanguage, getLanguage)
  const [unsaved, setUnsaved] = useState(false)
  return (
    <div className="language-picker">
      <select
        className="button quiet"
        aria-label={t('language.label')}
        value={language}
        onChange={(event) => setUnsaved(!setLanguage(event.target.value === 'ru' ? 'ru' : 'en'))}
      >
        <option value="en" lang="en">
          English
        </option>
        <option value="ru" lang="ru">
          Русский
        </option>
      </select>
      {unsaved && (
        <span role="status" className="language-warning">
          {t('language.unsaved')}
        </span>
      )}
    </div>
  )
}
