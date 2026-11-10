export const themes = [
  {
    id: 'cobalt',
    name: 'text.cobalt',
    description: 'text.aClearSkyForYourNextAdventure',
    background: '#f5f7fc',
    accent: '#1555d5',
  },
  {
    id: 'classic-dark',
    name: 'themes.classic_dark.name',
    description: 'themes.classic_dark.description',
    background: '#101010',
    accent: '#e6e6e6',
  },
  {
    id: 'midnight',
    name: 'text.midnight',
    description: 'text.quietBluesForLateNightSessions',
    background: '#141923',
    accent: '#99b8ff',
  },
  {
    id: 'parchment',
    name: 'text.parchment',
    description: 'text.warmPagesFromYourAdventureJournal',
    background: '#f4efe3',
    accent: '#805027',
  },
  {
    id: 'forest',
    name: 'text.forest',
    description: 'text.deepGreensForThePathLessTravelled',
    background: '#121d1a',
    accent: '#91c7ac',
  },
  {
    id: 'blossom',
    name: 'themes.blossom.name',
    description: 'themes.blossom.description',
    background: '#fff2f6',
    accent: '#a52d62',
  },
  {
    id: 'cyberpunk-light',
    name: 'themes.cyberpunk_light.name',
    description: 'themes.cyberpunk_light.description',
    background: '#f7f7f7',
    accent: '#cc202b',
  },
  {
    id: 'cyberpunk-dark',
    name: 'themes.cyberpunk_dark.name',
    description: 'themes.cyberpunk_dark.description',
    background: '#171410',
    accent: '#e18067',
  },
  {
    id: 'vampire',
    name: 'themes.vampire.name',
    description: 'themes.vampire.description',
    background: '#19131e',
    accent: '#c69be0',
  },
  {
    id: 'chthonic-horror',
    name: 'themes.chthonic_horror.name',
    description: 'themes.chthonic_horror.description',
    background: '#edf0e7',
    accent: '#365d3c',
  },
  {
    id: 'the-loop',
    name: 'themes.the_loop.name',
    description: 'themes.the_loop.description',
    background: '#faf8f3',
    accent: '#a94708',
  },
  {
    id: 'the-flood',
    name: 'themes.the_flood.name',
    description: 'themes.the_flood.description',
    background: '#f6f6f8',
    accent: '#803450',
  },
  {
    id: 'nuclear',
    name: 'themes.nuclear.name',
    description: 'themes.nuclear.description',
    background: '#f7f8f3',
    accent: '#20549a',
  },
  {
    id: 'cosmic-horror',
    name: 'themes.cosmic_horror.name',
    description: 'themes.cosmic_horror.description',
    background: '#101713',
    accent: '#b0d478',
  },
  {
    id: 'field-manual',
    name: 'themes.field_manual.name',
    description: 'themes.field_manual.description',
    background: '#e9ede6',
    accent: '#4d5d43',
  },
  {
    id: 'noir',
    name: 'themes.noir.name',
    description: 'themes.noir.description',
    background: '#f3f3f3',
    accent: '#333333',
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
