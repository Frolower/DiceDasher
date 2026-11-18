type Dictionary = Record<string, string>
type Modules = Record<string, { en: Dictionary; ru: Dictionary }>
type Intersection<T> = (T extends unknown ? (value: T) => void : never) extends (
  value: infer Result,
) => void
  ? Result
  : never

// Keep the exact key union for TranslationKey while rejecting silent overwrites.
export function mergeDictionaries<M extends Modules, L extends 'en' | 'ru'>(
  modules: M,
  language: L,
): Intersection<M[keyof M][L]> {
  const dictionary: Dictionary = Object.create(null)
  const owners = new Map<string, string>()
  for (const [name, translations] of Object.entries(modules)) {
    for (const [key, value] of Object.entries(translations[language])) {
      if (owners.has(key))
        throw new Error(
          `Duplicate translation "${key}" (${language}): ${owners.get(key)} and ${name}`,
        )
      owners.set(key, name)
      dictionary[key] = value
    }
  }
  return dictionary as Intersection<M[keyof M][L]>
}
