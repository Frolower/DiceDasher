import ts from 'typescript'
import fs from 'node:fs'
import assert from 'node:assert/strict'

const source = fs.readFileSync(new URL('../src/lib/auth.ts', import.meta.url), 'utf8')
const { outputText } = ts.transpileModule(source, {
  compilerOptions: { module: ts.ModuleKind.ESNext, target: ts.ScriptTarget.ES2023 },
})
const { validateAuth, authenticate } = await import(
  `data:text/javascript;base64,${Buffer.from(outputText).toString('base64')}`
)
const good = {
  email: 'alice@example.com',
  username: 'Alice',
  password: 'StrongPass1!',
  repeatPassword: 'StrongPass1!',
}
assert.deepEqual(validateAuth('signup', good), {})
assert.deepEqual(validateAuth('signup', { ...good, username: 'Alice_123' }), {})
for (const password of ['Aa1!' + 'x'.repeat(68), 'Symbols1$', 'Symbols1^', 'Symbols1~']) {
  assert.deepEqual(validateAuth('signup', { ...good, password, repeatPassword: password }), {})
}
for (const username of [
  'four',
  'Alice Smith',
  'Ali\nce',
  'a'.repeat(65),
  'Алиса',
  'Alicе',
  'Alicé',
  'Alice🎲',
]) {
  assert.ok(validateAuth('signup', { ...good, username }).username)
}
for (const email of [
  'bad',
  'a@',
  'a b@example.com',
  'frolowir@gmail',
  'alice@localhost',
  'alice@example.',
  'alice@.com',
]) {
  assert.ok(validateAuth('signup', { ...good, email }).email)
}
for (const password of [
  'Short1!',
  'alllower1!',
  'ALLUPPER1!',
  'NoNumber!!',
  'NoPunct12',
  'StrongPass1!Я',
  'StrongPass1!🎲',
  'StrongPass1!é',
  'StrongPass!١',
  'StrongPass1!\0',
  'Space Pass1!',
  'a'.repeat(64) + 'LongPass1!',
]) {
  assert.ok(validateAuth('signup', { ...good, password, repeatPassword: password }).password)
}
assert.ok(validateAuth('signup', { ...good, repeatPassword: 'different' }).repeatPassword)
assert.deepEqual(validateAuth('login', { ...good, password: 'legacy' }), {})
assert.ok(validateAuth('login', { ...good, password: '' }).password)

const originalFetch = globalThis.fetch
try {
  let sent
  globalThis.fetch = async (url, options) => {
    sent = { url, ...options }
    return new Response(JSON.stringify({ status: 'created' }), { status: 201 })
  }
  await authenticate('signup', good, new AbortController().signal)
  assert.equal(sent.url, '/api/register')
  assert.equal(sent.credentials, 'include')
  assert.deepEqual(JSON.parse(sent.body), {
    email: good.email,
    username: good.username,
    password: good.password,
  })
  globalThis.fetch = async () => new Response('internal details', { status: 409 })
  await assert.rejects(
    () => authenticate('signup', good, new AbortController().signal),
    /already in use/,
  )
  globalThis.fetch = async () => new Response('internal details', { status: 401 })
  await assert.rejects(() => authenticate('login', good, new AbortController().signal), /incorrect/)
  globalThis.fetch = async () =>
    new Response(
      JSON.stringify({ access_token: 'test-token', expires_at: '2026-10-25T18:00:00Z' }),
      { status: 200 },
    )
  assert.equal(
    (await authenticate('login', good, new AbortController().signal)).access_token,
    'test-token',
  )
} finally {
  globalThis.fetch = originalFetch
}
console.log(
  'Auth checks passed: validation boundaries, legacy login, request contract, and API error handling.',
)

for (const value of ['Алиса', 'Alicе', 'Alicé', 'Alice🎲', 'Alice Smith', 'Alice\0']) {
  assert.ok(validateAuth('login', { ...good, username: value }).username)
  assert.ok(validateAuth('login', { ...good, password: value }).password)
}
