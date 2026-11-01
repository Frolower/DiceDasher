import ts from 'typescript'
import fs from 'node:fs'
import assert from 'node:assert/strict'

const source = fs.readFileSync(new URL('../src/lib/auth.ts', import.meta.url), 'utf8')
const { outputText } = ts.transpileModule(source, {
  compilerOptions: { module: ts.ModuleKind.ESNext, target: ts.ScriptTarget.ES2023 },
})
const { validateAuth, authenticate, refreshSession, getSession, setSession, sessionFetch } =
  await import(`data:text/javascript;base64,${Buffer.from(outputText).toString('base64')}`)
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

try {
  const calls = []
  globalThis.fetch = async (url, options) => {
    calls.push({ url, ...options })
    return new Response(
      JSON.stringify(
        url === '/api/refresh'
          ? {
              access_token: 'restored-token',
              expires_at: new Date(Date.now() + 600000).toISOString(),
            }
          : { username: 'Alice', email: good.email },
      ),
    )
  }
  const first = refreshSession()
  const second = refreshSession()
  assert.equal(first, second, 'Concurrent refreshes must share token rotation')
  await Promise.all([first, second])
  assert.deepEqual(
    calls.map(({ url }) => url),
    ['/api/refresh', '/api/me'],
  )
  assert.equal(calls[0].method, 'POST')
  assert.equal(calls[0].credentials, 'include')
  assert.equal(calls[1].headers.Authorization, 'Bearer restored-token')
  assert.equal(getSession().username, 'Alice')

  globalThis.fetch = async () => new Response('', { status: 503 })
  await assert.rejects(refreshSession, /restore/)
  assert.equal(getSession().access_token, 'restored-token', 'Outages must preserve login')

  const requests = []
  globalThis.fetch = async (url, options) => {
    requests.push({ url, ...options })
    if (url === '/api/user' && options.headers.Authorization === 'Bearer restored-token')
      return new Response('', { status: 401 })
    return new Response(
      JSON.stringify(
        url === '/api/refresh'
          ? {
              access_token: 'renewed-token',
              expires_at: new Date(Date.now() + 600000).toISOString(),
            }
          : url === '/api/me'
            ? { username: 'Alice' }
            : 'success',
      ),
    )
  }
  const result = await sessionFetch('/api/user', 'old-token', {
    method: 'PATCH',
    body: JSON.stringify({ username: 'AliceNew' }),
  })
  assert.equal(result.status, 200)
  assert.deepEqual(
    requests.map(({ url }) => url),
    ['/api/user', '/api/refresh', '/api/me', '/api/user'],
  )
  assert.equal(requests[3].headers.Authorization, 'Bearer renewed-token')
  assert.equal(requests[3].body, requests[0].body)

  globalThis.fetch = async () => new Response('', { status: 401 })
  assert.equal(await refreshSession(), null)
  assert.equal(getSession(), null, 'Invalid refresh cookies must clear login')
} finally {
  setSession(null)
  globalThis.fetch = originalFetch
}
console.log(
  'Session checks passed: restoration, shared rotation, outage recovery, request retry, and expiry.',
)

const settingsSource = fs.readFileSync(new URL('../src/lib/settings.ts', import.meta.url), 'utf8')
const settingsCode = ts
  .transpileModule(settingsSource, {
    compilerOptions: { module: ts.ModuleKind.ESNext, target: ts.ScriptTarget.ES2023 },
  })
  .outputText.replace(
    "'./auth'",
    JSON.stringify(`data:text/javascript;base64,${Buffer.from(outputText).toString('base64')}`),
  )
const { validateProfile, validatePasswordChange, settingsRequest } = await import(
  `data:text/javascript;base64,${Buffer.from(settingsCode).toString('base64')}`
)
assert.deepEqual(validateProfile({ username: good.username, email: good.email }), {})
assert.ok(validateProfile({ username: 'four', email: 'bad' }).username)
assert.ok(validateProfile({ username: good.username, email: 'bad' }).email)
assert.deepEqual(
  validatePasswordChange({
    oldPassword: 'OldPass1!',
    password: 'NewPass2!',
    repeatPassword: 'NewPass2!',
  }),
  {},
)
assert.ok(
  validatePasswordChange({
    oldPassword: 'OldPass1!',
    password: 'OldPass1!',
    repeatPassword: 'OldPass1!',
  }).password,
)
assert.ok(
  validatePasswordChange({ oldPassword: '', password: 'NewPass2!', repeatPassword: 'different' })
    .oldPassword,
)
assert.ok(
  validatePasswordChange({
    oldPassword: 'OldPass1!',
    password: 'NewPass2!',
    repeatPassword: 'different',
  }).repeatPassword,
)
for (const password of ['short', 'NoSymbol12', 'Has Space1!', 'PassWord1!🎲', 'a'.repeat(73)]) {
  assert.ok(
    validatePasswordChange({ oldPassword: 'OldPass1!', password, repeatPassword: password })
      .password,
  )
}
try {
  let sent
  globalThis.fetch = async (url, options) => {
    sent = { url, ...options }
    return new Response(JSON.stringify('success'))
  }
  await settingsRequest('user', 'token', 'PATCH', { username: 'AliceNew', email: good.email })
  assert.equal(sent.url, '/api/user')
  assert.equal(sent.method, 'PATCH')
  assert.equal(sent.headers.Authorization, 'Bearer token')
  assert.deepEqual(JSON.parse(sent.body), { username: 'AliceNew', email: good.email })
  await settingsRequest('user', 'token', 'DELETE')
  assert.equal(sent.url, '/api/user')
  assert.equal(sent.method, 'DELETE')
  assert.equal(sent.headers.Authorization, 'Bearer token')
  assert.equal(sent.body, undefined)
  await settingsRequest('password', 'token', 'PUT', {
    old_password: 'OldPass1!',
    new_password: 'NewPass2!',
  })
  assert.equal(sent.url, '/api/password')
  assert.equal(sent.method, 'PUT')
  assert.deepEqual(JSON.parse(sent.body), { old_password: 'OldPass1!', new_password: 'NewPass2!' })
  globalThis.fetch = async () => new Response('old password does not match', { status: 400 })
  await assert.rejects(
    () => settingsRequest('password', 'token', 'PUT', {}),
    /old password is incorrect/,
  )
  globalThis.fetch = async () => new Response('', { status: 401 })
  await assert.rejects(() => settingsRequest('me', 'token'), /session has expired/)
} finally {
  globalThis.fetch = originalFetch
}
console.log(
  'Settings checks passed: signup rules, password confirmation and reuse, authenticated API contracts, and errors.',
)
