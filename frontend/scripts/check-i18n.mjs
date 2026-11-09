import ts from 'typescript'
import fs from 'node:fs'
import assert from 'node:assert/strict'

const read = (path) => fs.readFileSync(new URL(path, import.meta.url), 'utf8')
const en = JSON.parse(read('../src/locales/en.json'))
const ru = JSON.parse(read('../src/locales/ru.json'))
assert.deepEqual(Object.keys(ru).sort(), Object.keys(en).sort())
for (const key of Object.keys(en)) {
  assert.ok(en[key].trim() && ru[key].trim(), `Empty translation: ${key}`)
  const placeholders = (value) => [...value.matchAll(/\{(\w+)\}/g)].map((match) => match[1]).sort()
  assert.deepEqual(placeholders(en[key]), placeholders(ru[key]), `Parameters differ: ${key}`)
}
let source = read('../src/lib/i18n.ts')
  .replace("import en from '../locales/en.json'", `const en = ${JSON.stringify(en)}`)
  .replace("import ru from '../locales/ru.json'", `const ru = ${JSON.stringify(ru)}`)
const code = ts.transpileModule(source, {
  compilerOptions: { module: ts.ModuleKind.ESNext, target: ts.ScriptTarget.ES2023 },
}).outputText
let instance = 0
const load = () =>
  import(
    `data:text/javascript;base64,${Buffer.from(code + `\n// ${instance++}`).toString('base64')}`
  )
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
      [h(ProfileSettings, { token: 'test' }), 'text.settings'],
      [
        h(SideMenu, { logoutError: 'text.unableToLogOutPleaseTryAgain' }),
        'text.unableToLogOutPleaseTryAgain',
      ],
      [h(LanguagePicker), 'language.label'],
    ]
    for (const [element, key] of screens) {
      const markup = render(element)
      assert.ok(markup.includes(dictionary[key]), `Missing ${language} translation for ${key}`)
      assert.doesNotMatch(
        markup,
        /(?:text|auth|themes|language)\.[A-Za-z]+/,
        'Translation key leaked into markup',
      )
    }
    assert.ok(render(h(ProfilePage, { username: 'Алиса' })).includes('Алиса'))
  }
} finally {
  await server.close()
}
console.log(
  'Screen checks passed: both languages render across auth, profile, settings, themes, navigation, and demo pages.',
)
