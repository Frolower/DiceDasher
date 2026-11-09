export const themes = [
  {
    id: 'cobalt',
    name: 'Cobalt',
    description: 'A clear sky for your next adventure.',
    background: '#f5f7fc',
    accent: '#1555d5',
  },
  {
    id: 'midnight',
    name: 'Midnight',
    description: 'Quiet blues for late-night sessions.',
    background: '#141923',
    accent: '#99b8ff',
  },
  {
    id: 'parchment',
    name: 'Parchment',
    description: 'Warm pages from your adventure journal.',
    background: '#f4efe3',
    accent: '#805027',
  },
  {
    id: 'forest',
    name: 'Forest',
    description: 'Deep greens for the path less travelled.',
    background: '#121d1a',
    accent: '#91c7ac',
  },
] as const
export type ThemeId = (typeof themes)[number]['id']
export const themeStorageKey = 'dicedasher.theme'
export function isTheme(value: unknown): value is ThemeId {
  return themes.some((theme) => theme.id === value)
}
export function currentTheme(): ThemeId {
  const value = document.documentElement.dataset.theme
  return isTheme(value) ? value : 'cobalt'
}
export function applyTheme(theme: ThemeId): boolean {
  document.documentElement.dataset.theme = theme
  try {
    localStorage.setItem(themeStorageKey, theme)
    return true
  } catch {
    return false
  }
}

// Keep the whole app in sync when another tab changes or clears the preference.
export function syncTheme(event: StorageEvent) {
  if (event.key !== themeStorageKey && event.key !== null) return
  const theme = isTheme(event.newValue)
    ? event.newValue
    : window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'midnight'
      : 'cobalt'
  document.documentElement.dataset.theme = theme
}
