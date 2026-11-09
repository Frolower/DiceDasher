import { useEffect, useState } from 'react'
import type { CSSProperties } from 'react'
import { applyTheme, currentTheme, themes, themeStorageKey } from '../lib/themes'
import type { ThemeId } from '../lib/themes'

export function ThemeSettings() {
  const [selected, setSelected] = useState(currentTheme)
  const [message, setMessage] = useState('')
  useEffect(() => {
    function sync(event: StorageEvent) {
      if (event.key !== themeStorageKey && event.key !== null) return
      setSelected(currentTheme())
      setMessage('')
    }
    window.addEventListener('storage', sync)
    return () => window.removeEventListener('storage', sync)
  }, [])
  function choose(theme: ThemeId) {
    const saved = applyTheme(theme)
    setSelected(theme)
    const name = themes.find((item) => item.id === theme)!.name
    setMessage(
      saved
        ? `${name} applied and saved.`
        : `${name} applied. Your browser couldn’t save it for next time.`,
    )
  }
  return (
    <section aria-labelledby="themes-title">
      <h2 id="themes-title" className="settings-section-title">
        Themes
      </h2>
      <p className="settings-section-intro">Choose the atmosphere for your next chapter.</p>
      <fieldset className="theme-picker" aria-describedby="theme-storage-note">
        <legend className="sr-only">Interface theme</legend>
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
              >
                <span className="theme-swatch-die">◇</span>
              </span>
              <span className="theme-option-heading">
                <span>{theme.name}</span>
                <span className="theme-check" aria-hidden="true">
                  {selected === theme.id ? '✓' : ''}
                </span>
              </span>
              <span className="theme-description">{theme.description}</span>
            </label>
          ))}
        </div>
      </fieldset>
      <p className="theme-storage-note" id="theme-storage-note">
        Changes apply immediately and are saved in this browser.
      </p>
      <p className="theme-feedback" role="status">
        {message}
      </p>
    </section>
  )
}
