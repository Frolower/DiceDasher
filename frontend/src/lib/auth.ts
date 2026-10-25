export type AuthMode = 'login' | 'signup'
export type AuthFields = { email: string; username: string; password: string; repeatPassword: string }
export type FieldErrors = Partial<Record<keyof AuthFields, string>>
const bytes = (value: string) => new TextEncoder().encode(value).length
const emailPattern = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$/
export const passwordChecks = (password: string) => [
  { label: '8–72 characters', valid: [...password].length >= 8 && bytes(password) <= 72 },
  { label: 'Upper & lowercase (A–Z, a–z)', valid: /[a-z]/.test(password) && /[A-Z]/.test(password) },
  { label: 'A number (0–9)', valid: /[0-9]/.test(password) },
  { label: 'A symbol (! @ # $ …)', valid: /[^a-zA-Z0-9]/.test(password) && /^[\x21-\x7e]+$/.test(password) },
]
export function validateAuth(mode: AuthMode, fields: AuthFields): FieldErrors {
  const errors: FieldErrors = {}
  if (!fields.username) errors.username = 'Enter your username.'
  else if ([...fields.username].length > 64) errors.username = 'Use no more than 64 characters.'
  if (/[\p{White_Space}\uFEFF]/u.test(fields.username)) errors.username = 'Username cannot contain spaces.'
  else if (/\p{Cc}/u.test(fields.username)) errors.username = 'Username contains an invalid character.'
  else if (/[^\x21-\x7e]/.test(fields.username)) errors.username = 'Use English letters, numbers and symbols only.'
  if (!fields.password) errors.password = 'Enter your password.'
  else if (bytes(fields.password) > 72) errors.password = 'This password is too long. Shorten it and try again.'
  if (mode === 'signup') {
    if (fields.username && [...fields.username].length < 5) errors.username = 'Use at least 5 characters.'
    if (bytes(fields.email) < 5 || bytes(fields.email) > 254 || !emailPattern.test(fields.email)) errors.email = 'Enter a valid email address, such as you@example.com.'
    if (passwordChecks(fields.password).some(check => !check.valid) || /[\p{White_Space}\uFEFF]/u.test(fields.password)) errors.password = 'Meet the password requirements below.'
    if (bytes(fields.password) > 72) errors.password = 'Use no more than 72 characters.'
    if (!fields.repeatPassword) errors.repeatPassword = 'Repeat your password.'
    else if (fields.password !== fields.repeatPassword) errors.repeatPassword = 'Your passwords do not match.'
  }
  if (/[^\x21-\x7e]/.test(fields.password)) errors.password = 'Use only English letters, numbers and special characters. No spaces or emoji.'
  return errors
}
export async function authenticate(mode: AuthMode, fields: AuthFields, signal: AbortSignal) {
  const response = await fetch(`/api/${mode === 'login' ? 'login' : 'register'}`, {
    method: 'POST', credentials: 'include', signal,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(mode === 'login'
      ? { username: fields.username, password: fields.password }
      : { username: fields.username, email: fields.email, password: fields.password }),
  })
  if (!response.ok) {
    const message = response.status === 409 ? 'That username or email is already in use.'
      : response.status === 401 ? 'Your username or password is incorrect.'
      : response.status === 403 ? 'Sign-in is not allowed from this address. Check the backend ALLOWED_ORIGIN setting.'
      : response.status === 400 ? 'The server could not accept these details. Check the fields and try again.'
      : response.status === 429 ? 'Too many attempts. Please try again in a moment.'
      : 'We could not reach the account service. Please try again shortly.'
    throw new Error(message)
  }
  if (mode === 'login') {
    const data = await response.json() as { access_token: string; expires_at: string }
    if (!data.access_token) throw new Error('The account service returned an incomplete response. Please try again.')
    return data
  }
  return null
}
