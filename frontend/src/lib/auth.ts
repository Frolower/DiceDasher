export type AuthMode = 'login' | 'signup'
export type AccessTokens = { access_token: string; expires_at: string }
export type Session = AccessTokens & { username: string }
let currentSession: Session | null = null
const listeners = new Set<() => void>()
export const getSession = () => currentSession
export function subscribeSession(listener: () => void) {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}
export function setSession(session: Session | null) {
  currentSession = session
  listeners.forEach((listener) => listener())
}

// Share refreshes across startup, timers and requests. Never abort a rotating
// refresh on React effect cleanup: the server may already have replaced the cookie.
let refreshing: Promise<Session | null> | null = null
let loggingOut: Promise<void> | null = null
export function refreshSession(): Promise<Session | null> {
  if (loggingOut)
    return loggingOut.then(
      () => null,
      () => currentSession,
    )
  if (refreshing) return refreshing
  refreshing = (async () => {
    const response = await fetch('/api/refresh', {
      method: 'POST',
      credentials: 'include',
      cache: 'no-store',
    })
    if (response.status === 401) {
      setSession(null)
      return null
    }
    if (!response.ok) throw new Error('Unable to restore your session. Please try again.')
    const tokens = (await response.json()) as AccessTokens
    if (!tokens.access_token || !Number.isFinite(Date.parse(tokens.expires_at)))
      throw new Error('The account service returned an incomplete session.')
    const profileResponse = await fetch('/api/me', {
      credentials: 'include',
      cache: 'no-store',
      headers: { Authorization: `Bearer ${tokens.access_token}` },
    })
    if (profileResponse.status === 401) {
      setSession(null)
      return null
    }
    if (!profileResponse.ok) throw new Error('Unable to load your account. Please try again.')
    const profile = (await profileResponse.json()) as { username: string }
    const session = { ...tokens, username: profile.username }
    setSession(session)
    return session
  })().finally(() => {
    refreshing = null
  })
  return refreshing
}

export async function sessionFetch(path: string, token: string, options: RequestInit) {
  const send = (access: string) =>
    fetch(path, {
      ...options,
      headers: { ...options.headers, Authorization: `Bearer ${access}` },
    })
  const response = await send(currentSession?.access_token ?? token)
  if (response.status !== 401) return response
  const session = await refreshSession()
  if (!session || options.signal?.aborted) return response
  return send(session.access_token)
}
export type AuthFields = {
  email: string
  username: string
  password: string
  repeatPassword: string
}
export type FieldErrors = Partial<Record<keyof AuthFields, string>>
const bytes = (value: string) => new TextEncoder().encode(value).length
const emailPattern =
  /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$/
export const passwordChecks = (password: string) => [
  { label: '8–72 characters', valid: [...password].length >= 8 && bytes(password) <= 72 },
  {
    label: 'Upper & lowercase (A–Z, a–z)',
    valid: /[a-z]/.test(password) && /[A-Z]/.test(password),
  },
  { label: 'A number (0–9)', valid: /[0-9]/.test(password) },
  {
    label: 'A symbol (! @ # $ …)',
    valid: /[^a-zA-Z0-9]/.test(password) && /^[\x21-\x7e]+$/.test(password),
  },
]
export function validateAuth(mode: AuthMode, fields: AuthFields): FieldErrors {
  const errors: FieldErrors = {}
  if (!fields.username) errors.username = 'Enter your username.'
  else if ([...fields.username].length > 64) errors.username = 'Use no more than 64 characters.'
  if (/[\p{White_Space}\uFEFF]/u.test(fields.username))
    errors.username = 'Username cannot contain spaces.'
  else if (/\p{Cc}/u.test(fields.username))
    errors.username = 'Username contains an invalid character.'
  else if (/[^\x21-\x7e]/.test(fields.username))
    errors.username = 'Use English letters, numbers and symbols only.'
  if (!fields.password) errors.password = 'Enter your password.'
  else if (bytes(fields.password) > 72)
    errors.password = 'This password is too long. Shorten it and try again.'
  if (mode === 'signup') {
    if (fields.username && [...fields.username].length < 5)
      errors.username = 'Use at least 5 characters.'
    if (bytes(fields.email) < 5 || bytes(fields.email) > 254 || !emailPattern.test(fields.email))
      errors.email = 'Enter a valid email address, such as you@example.com.'
    if (
      passwordChecks(fields.password).some((check) => !check.valid) ||
      /[\p{White_Space}\uFEFF]/u.test(fields.password)
    )
      errors.password = 'Meet the password requirements below.'
    if (bytes(fields.password) > 72) errors.password = 'Use no more than 72 characters.'
    if (!fields.repeatPassword) errors.repeatPassword = 'Repeat your password.'
    else if (fields.password !== fields.repeatPassword)
      errors.repeatPassword = 'Your passwords do not match.'
  }
  if (/[^\x21-\x7e]/.test(fields.password))
    errors.password =
      'Use only English letters, numbers and special characters. No spaces or emoji.'
  return errors
}
export async function authenticate(mode: AuthMode, fields: AuthFields, signal: AbortSignal) {
  const response = await fetch(`/api/${mode === 'login' ? 'login' : 'register'}`, {
    method: 'POST',
    credentials: 'include',
    signal,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(
      mode === 'login'
        ? { username: fields.username, password: fields.password }
        : { username: fields.username, email: fields.email, password: fields.password },
    ),
  })
  if (!response.ok) {
    const message =
      response.status === 409
        ? 'That username or email is already in use.'
        : response.status === 401
          ? 'Your username or password is incorrect.'
          : response.status === 403
            ? 'Sign-in is not allowed from this address. Check the backend ALLOWED_ORIGINS setting.'
            : response.status === 400
              ? 'The server could not accept these details. Check the fields and try again.'
              : response.status === 429
                ? 'Too many attempts. Please try again in a moment.'
                : 'We could not reach the account service. Please try again shortly.'
    throw new Error(message)
  }
  if (mode === 'login') {
    const data = (await response.json()) as { access_token: string; expires_at: string }
    if (!data.access_token)
      throw new Error('The account service returned an incomplete response. Please try again.')
    return data
  }
  return null
}
export function logout(): Promise<void> {
  if (loggingOut) return loggingOut
  loggingOut = (async () => {
    // Finish any cookie rotation before revoking the refresh token.
    await refreshing?.catch(() => null)
    const response = await fetch('/api/logout', {
      method: 'POST',
      credentials: 'include',
      cache: 'no-store',
    })
    if (!response.ok) {
      throw new Error(
        response.status === 403
          ? 'Logout is not allowed from this address. Check the backend ALLOWED_ORIGINS setting.'
          : 'Unable to log out. Please try again.',
      )
    }
    setSession(null)
  })().finally(() => {
    loggingOut = null
  })
  return loggingOut
}
