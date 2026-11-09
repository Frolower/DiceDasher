import en from '../locales/en.json'
import ru from '../locales/ru.json'

export type Language = 'en' | 'ru'
export type TranslationKey = keyof typeof en
const dictionaries: Record<Language, Partial<Record<TranslationKey, string>>> = { en, ru }
export const languageStorageKey = 'dicedasher.language'
const listeners = new Set<() => void>()
const isLanguage = (value: unknown): value is Language => value === 'en' || value === 'ru'
function detectLanguage(): Language {
  if (typeof navigator === 'undefined') return 'en'
  const language = (navigator.languages?.[0] ?? navigator.language).split('-')[0].toLowerCase()
  return isLanguage(language) ? language : 'en'
}
function initialLanguage(): Language {
  try {
    const saved = localStorage.getItem(languageStorageKey)
    if (isLanguage(saved)) return saved
  } catch {
    /* Storage is optional. */
  }
  return detectLanguage()
}
let language = initialLanguage()
export const getLanguage = () => language
export function subscribeLanguage(listener: () => void) {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}
function updateLanguage(next: Language) {
  language = next
  if (typeof document !== 'undefined') document.documentElement.lang = next
  listeners.forEach((listener) => listener())
}
export function setLanguage(next: Language): boolean {
  updateLanguage(next)
  try {
    localStorage.setItem(languageStorageKey, next)
    return true
  } catch {
    return false
  }
}
export function syncLanguage(event: StorageEvent) {
  if (event.key !== languageStorageKey && event.key !== null) return
  updateLanguage(isLanguage(event.newValue) ? event.newValue : detectLanguage())
}
// Stored UI messages use keys. Unknown keys remain visible to aid diagnosis.
export function t(key: string, params: Record<string, string | number> = {}): string {
  const translated = dictionaries[language][key as TranslationKey]
  const template = translated?.trim() ? translated : (en[key as TranslationKey] ?? key)
  return template.replace(/\{(\w+)\}/g, (match, name: string) => String(params[name] ?? match))
}
export function plural(
  key: { one: TranslationKey; few: TranslationKey; many: TranslationKey; other: TranslationKey },
  count: number,
) {
  const form = new Intl.PluralRules(language).select(count)
  return t(key[form === 'one' || form === 'few' || form === 'many' ? form : 'other'], { count })
}
updateLanguage(language)
