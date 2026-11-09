import { t } from '../lib/i18n'
import { useEffect, useState } from 'react'
import type { CSSProperties } from 'react'
import { applyTheme, currentTheme, themes, themeStorageKey } from '../lib/themes'
import type { ThemeId } from '../lib/themes'

export function ThemeSettings() {
  const [selected, setSelected] = useState(currentTheme)
  const [message, setMessage] = useState<{ key: string; name: string } | null>(null)
  useEffect(() => {
    function sync(event: StorageEvent) {
      if (event.key !== themeStorageKey && event.key !== null) return
      setSelected(currentTheme())
      setMessage(null)
    }
    window.addEventListener('storage', sync)
    return () => window.removeEventListener('storage', sync)
  }, [])
  function choose(theme: ThemeId) {
    const saved = applyTheme(theme)
    setSelected(theme)
    const name = themes.find((item) => item.id === theme)!.name
    setMessage({ key: saved ? 'themes.saved' : 'themes.unsaved', name })
  }
  return (
    <section aria-labelledby="themes-title">
      <h2 id="themes-title" className="settings-section-title">
        {t('text.themes')}
      </h2>
      <p className="settings-section-intro">{t('text.chooseTheAtmosphereForYourNextChapter')}</p>
      <fieldset className="theme-picker">
        <legend className="sr-only">{t('text.interfaceTheme')}</legend>
        <div className="theme-options">
          {themes.map((theme) => (
            <label className="theme-option" key={theme.id}>
              <input
                type="radio"
                name="interface-theme"
                value={theme.id}
                checked={selected === theme.id}
                onChange={() => choose(theme.id)}
              />
              <span
                className="theme-swatch"
                aria-hidden="true"
                style={
                  {
                    '--swatch-background': theme.background,
                    '--swatch-accent': theme.accent,
                  } as CSSProperties
                }
              />
              <span className="theme-option-heading">
                <span>{t(theme.name)}</span>
                <span className="theme-check" aria-hidden="true">
                  {selected === theme.id ? '✓' : ''}
                </span>
              </span>
              <span className="theme-description">{t(theme.description)}</span>
            </label>
          ))}
        </div>
      </fieldset>
      <p className="theme-feedback" role="status">
        {message && t(message.key, { name: t(message.name) })}
      </p>
    </section>
  )
}
