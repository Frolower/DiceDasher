import en from '../locales/en.json'
import type { TranslationKey } from './i18n'
// Only explicitly written messages are safe to display in the interface.
export class UserFacingError extends Error {
  readonly key: string
  constructor(key: string) {
    super(en[key as TranslationKey] ?? key)
    this.key = key
  }
}
export function userErrorMessage(error: unknown, fallback: string): string {
  return error instanceof UserFacingError ? error.key : fallback
}

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
    if (!response.ok) throw new UserFacingError('text.unableToRestoreYourSessionPleaseTryAgain')
    const tokens = (await response.json()) as AccessTokens
    if (!tokens.access_token || !Number.isFinite(Date.parse(tokens.expires_at)))
      throw new UserFacingError('text.unableToRestoreYourSessionPleaseTryAgain')
    const profileResponse = await fetch('/api/me', {
      credentials: 'include',
      cache: 'no-store',
      headers: { Authorization: `Bearer ${tokens.access_token}` },
    })
    if (profileResponse.status === 401) {
      setSession(null)
      return null
    }
    if (!profileResponse.ok) throw new UserFacingError('text.unableToLoadYourAccountPleaseTryAgain')
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
  { label: 'auth.passwordLength', valid: [...password].length >= 8 && bytes(password) <= 72 },
  {
    label: 'text.upperLowercaseAZAZ',
    valid: /[a-z]/.test(password) && /[A-Z]/.test(password),
  },
  { label: 'text.aNumber09', valid: /[0-9]/.test(password) },
  {
    label: 'auth.passwordSymbol',
    valid: /[^a-zA-Z0-9]/.test(password) && /^[\x21-\x7e]+$/.test(password),
  },
]
export function validateAuth(mode: AuthMode, fields: AuthFields): FieldErrors {
  const errors: FieldErrors = {}
  if (!fields.username) errors.username = 'text.enterYourUsername'
  else if ([...fields.username].length > 64) errors.username = 'text.useNoMoreThan64Characters'
  if (/[\p{White_Space}\uFEFF]/u.test(fields.username))
    errors.username = 'text.usernameCannotContainSpaces'
  else if (/\p{Cc}/u.test(fields.username))
    errors.username = 'text.usernameContainsAnInvalidCharacter'
  else if (/[^\x21-\x7e]/.test(fields.username))
    errors.username = 'text.useEnglishLettersNumbersAndSymbolsOnly'
  if (!fields.password) errors.password = 'errors.passwordRequired'
  else if (bytes(fields.password) > 72)
    errors.password = 'text.thisPasswordIsTooLongShortenItAndTryAgain'
  if (mode === 'signup') {
    if (fields.username && [...fields.username].length < 5)
      errors.username = 'text.useAtLeast5Characters'
    if (bytes(fields.email) < 5 || bytes(fields.email) > 254 || !emailPattern.test(fields.email))
      errors.email = 'text.enterAValidEmailAddressSuchAsYouExampleCom'
    if (
      passwordChecks(fields.password).some((check) => !check.valid) ||
      /[\p{White_Space}\uFEFF]/u.test(fields.password)
    )
      errors.password = 'text.meetThePasswordRequirementsBelow'
    if (bytes(fields.password) > 72) errors.password = 'text.useNoMoreThan72Characters'
    if (!fields.repeatPassword) errors.repeatPassword = 'text.repeatYourPassword'
    else if (fields.password !== fields.repeatPassword)
      errors.repeatPassword = 'text.yourPasswordsDoNotMatch'
  }
  if (/[^\x21-\x7e]/.test(fields.password))
    errors.password = 'text.useOnlyEnglishLettersNumbersAndSpecialCharactersNoSpacesOrEmoji'
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
        ? 'text.thatUsernameOrEmailIsAlreadyInUse'
        : response.status === 401
          ? 'text.yourUsernameOrPasswordIsIncorrect'
          : response.status === 403
            ? 'text.unableToSignInRightNowPleaseTryAgainLater'
            : response.status === 400
              ? 'text.checkYourDetailsAndTryAgain'
              : response.status === 429
                ? 'text.tooManyAttemptsPleaseTryAgainInAMoment'
                : mode === 'login'
                  ? 'text.unableToSignInRightNowPleaseTryAgainLater'
                  : 'auth.unableToCreate'
    throw new UserFacingError(message)
  }
  if (mode === 'login') {
    const data = (await response.json()) as { access_token: string; expires_at: string }
    if (!data.access_token)
      throw new UserFacingError('text.unableToSignInRightNowPleaseTryAgainLater')
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
      throw new UserFacingError('text.unableToLogOutPleaseTryAgain')
    }
    setSession(null)
  })().finally(() => {
    loggingOut = null
  })
  return loggingOut
}
