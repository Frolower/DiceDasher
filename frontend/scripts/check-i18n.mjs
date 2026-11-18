import { loadLocaleRegistry } from './locale-registry.mjs'
import ts from 'typescript'
import fs from 'node:fs'
import assert from 'node:assert/strict'

const read = (path) => fs.readFileSync(new URL(path, import.meta.url), 'utf8')
const { en, ru, localeModules, mergeDictionaries } = await loadLocaleRegistry()
const localeRoot = new URL('../src/locales/', import.meta.url)
const registeredFiles = []
const parameters = (value) => [...value.matchAll(/\{(\w+)\}/g)].map((match) => match[1]).sort()
for (const [namespace, translations] of Object.entries(localeModules)) {
  assert.deepEqual(
    Object.keys(translations.en).sort(),
    Object.keys(translations.ru).sort(),
    `Language keys differ: ${namespace}`,
  )
  for (const language of ['en', 'ru']) {
    const relative = `${namespace}/${language}.json`
    registeredFiles.push(relative)
    const source = fs.readFileSync(new URL(relative, localeRoot), 'utf8')
    const parsed = JSON.parse(source)
    assert.deepEqual(parsed, translations[language], `Registry differs from ${relative}`)
    const properties = ts.parseJsonText(relative, source).statements[0].expression.properties
    const keys = properties.map((property) => property.name.text)
    assert.equal(new Set(keys).size, keys.length, `Repeated key within ${relative}`)
    for (const [key, value] of Object.entries(parsed)) {
      assert.equal(typeof value, 'string', `Non-string translation: ${relative}: ${key}`)
      assert.ok(value.trim(), `Empty translation: ${relative}: ${key}`)
      assert.deepEqual(
        parameters(value),
        parameters(translations.en[key]),
        `Parameters differ: ${relative}: ${key}`,
      )
    }
  }
}
const actualFiles = fs
  .readdirSync(localeRoot, { recursive: true })
  .filter((file) => file.endsWith('.json'))
assert.deepEqual(
  actualFiles.sort(),
  registeredFiles.sort(),
  'Every dictionary file must be registered',
)
for (const language of ['en', 'ru']) {
  assert.throws(
    () =>
      mergeDictionaries(
        {
          first: { en: { 'test.same': 'One' }, ru: { 'test.same': 'Один' } },
          second: { en: { 'test.same': 'Two' }, ru: { 'test.same': 'Два' } },
        },
        language,
      ),
    /Duplicate translation "test.same".*first and second/,
  )
}
assert.deepEqual(Object.keys(en).sort(), Object.keys(ru).sort())
let instance = 0
const load = (translations = ru) => {
  const source = read('../src/lib/i18n.ts').replace(
    "import { en, ru } from '../locales'",
    `const en = ${JSON.stringify(en)}; const ru = ${JSON.stringify(translations)}`,
  )
  const code = ts.transpileModule(source, {
    compilerOptions: { module: ts.ModuleKind.ESNext, target: ts.ScriptTarget.ES2023 },
  }).outputText
  return import(
    `data:text/javascript;base64,${Buffer.from(code + `\n// ${instance++}`).toString('base64')}`
  )
}
const storage = new Map()
globalThis.localStorage = {
  getItem: (key) => storage.get(key) ?? null,
  setItem: (key, value) => storage.set(key, value),
}
globalThis.document = { documentElement: { lang: '' } }
Object.defineProperty(globalThis, 'navigator', {
  configurable: true,
  value: { language: 'ru-RU', languages: ['ru-RU', 'en'] },
})
const i18n = await load()
assert.equal(i18n.getLanguage(), 'ru')
assert.equal(document.documentElement.lang, 'ru')
assert.equal(i18n.t('text.logIn'), 'Войти')
for (const value of [undefined, '', '   ']) {
  const incomplete = { ...ru, 'text.logIn': value, 'themes.saved': value }
  const fallback = await load(incomplete)
  assert.equal(fallback.getLanguage(), 'ru')
  assert.equal(fallback.t('text.logIn'), 'Log in', 'Missing or blank Russian text uses English')
  assert.equal(fallback.t('themes.saved', { name: 'Cobalt' }), 'Cobalt applied and saved.')
}

let notifications = 0
const unsubscribe = i18n.subscribeLanguage(() => notifications++)
assert.equal(i18n.setLanguage('en'), true)
assert.equal(storage.get(i18n.languageStorageKey), 'en')
assert.equal(document.documentElement.lang, 'en')
assert.equal(i18n.t('text.logIn'), 'Log in')
assert.equal(notifications, 1)
assert.equal((await load()).getLanguage(), 'en', 'Saved choice takes priority over browser')
assert.equal(i18n.t('themes.saved', { name: 'Cobalt' }), 'Cobalt applied and saved.')
i18n.setLanguage('ru')
assert.equal(i18n.t('themes.saved', { name: 'Кобальт' }), 'Тема «Кобальт» применена и сохранена.')
assert.equal(i18n.t('unknown.key'), 'unknown.key')
i18n.syncLanguage({ key: 'unrelated', newValue: 'en' })
assert.equal(i18n.getLanguage(), 'ru')
i18n.syncLanguage({ key: i18n.languageStorageKey, newValue: 'en' })
assert.equal(i18n.getLanguage(), 'en')
i18n.syncLanguage({ key: null, newValue: null })
assert.equal(i18n.getLanguage(), 'ru', 'Clearing storage restores browser default')
unsubscribe()
const before = notifications
localStorage.setItem = () => {
  throw new Error('Storage blocked')
}
assert.equal(i18n.setLanguage('en'), false)
assert.equal(i18n.getLanguage(), 'en', 'Language still changes with blocked storage')
assert.equal(notifications, before)
localStorage.getItem = () => {
  throw new Error('Storage blocked')
}
navigator.languages = ['de-DE']
assert.equal(
  (await load()).getLanguage(),
  'en',
  'Unsupported browser language falls back to English',
)
console.log(
  'Localization checks passed: dictionaries, parameters, detection, persistence, switching, blocked storage, and tab sync.',
)

// Render real screens in both languages to catch keys accidentally shown as text.
const { createServer } = await import('vite')
const { createElement: h } = await import('react')
const { renderToStaticMarkup } = await import('react-dom/server')
const { MemoryRouter } = await import('react-router')
const server = await createServer({
  server: { middlewareMode: true, ws: false },
  appType: 'custom',
})
try {
  document.documentElement.dataset = { theme: 'cobalt' }
  const uiLanguage = await server.ssrLoadModule('/src/lib/i18n.ts')
  const { AuthDialog } = await server.ssrLoadModule('/src/components/AuthDialog.tsx')
  const { ThemeSettings } = await server.ssrLoadModule('/src/components/ThemeSettings.tsx')
  const { ProfilePage } = await server.ssrLoadModule('/src/components/ProfilePage.tsx')
  const { DemoPage } = await server.ssrLoadModule('/src/components/DemoPage.tsx')
  const { CharactersPage } = await server.ssrLoadModule(
    '/src/features/characters/CharactersPage.tsx',
  )
  const { ProfileSettings } = await server.ssrLoadModule('/src/components/ProfileSettings.tsx')
  const { SideMenu } = await server.ssrLoadModule('/src/components/SideMenu.tsx')
  const { LanguagePicker } = await server.ssrLoadModule('/src/components/LanguagePicker.tsx')
  const render = (element) => renderToStaticMarkup(h(MemoryRouter, null, element))
  for (const language of ['en', 'ru']) {
    uiLanguage.setLanguage(language)
    const dictionary = language === 'en' ? en : ru
    const screens = [
      [h(AuthDialog, { mode: 'login' }), 'text.logIn'],
      [h(AuthDialog, { mode: 'signup' }), 'auth.passwordLength'],
      [h(ThemeSettings), 'text.cobalt'],
      [h(ProfilePage, { username: 'Алиса' }), 'text.rooms'],
      [h(DemoPage, { kind: 'characters' }), 'text.aNocturnalStorytellerWithADangerousSecret'],
      [h(DemoPage, { kind: 'friends' }), 'text.friends'],
      [h(CharactersPage), 'characters.emptyTitle'],
      [h(ProfileSettings, { token: 'test' }), 'text.settings'],
      [
        h(SideMenu, { logoutError: 'text.unableToLogOutPleaseTryAgain' }),
        'text.unableToLogOutPleaseTryAgain',
      ],
      [h(LanguagePicker), 'language.label'],
    ]
    for (const [element, key] of screens) {
      const markup = render(element)
      assert.ok(
        markup.includes(dictionary[key]?.trim() ? dictionary[key] : en[key]),
        `Missing ${language} translation for ${key}`,
      )
      assert.doesNotMatch(
        markup,
        /(?:text|auth|themes|language)\.[A-Za-z]+/,
        'Translation key leaked into markup',
      )
    }
    assert.ok(render(h(ProfilePage, { username: 'Алиса' })).includes('Алиса'))
    const characterList = render(h(CharactersPage))
    assert.doesNotMatch(
      characterList,
      /Сородич|Vampire|Masquerade|V5|ноч[ьи]/,
      'The empty character collection must stay system neutral',
    )
    if (language === 'en')
      assert.doesNotMatch(
        characterList,
        /[А-Яа-яЁё]/,
        'The English character list must not contain hardcoded Russian',
      )
  }
} finally {
  await server.close()
}
console.log(
  'Screen checks passed: both languages render across auth, profile, settings, themes, navigation, and demo pages.',
)
