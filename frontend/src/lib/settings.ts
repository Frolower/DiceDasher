import { validateAuth } from './auth'

export type Profile = { username: string; email: string }
export type PasswordFields = { oldPassword: string; password: string; repeatPassword: string }
export function validateProfile(profile: Profile) {
  const errors = validateAuth('signup', {
    ...profile,
    password: 'ValidPass1!',
    repeatPassword: 'ValidPass1!',
  })
  return errors
}
export function validatePasswordChange(fields: PasswordFields) {
  const errors: Partial<Record<keyof PasswordFields, string>> = {}
  const check = (password: string) =>
    validateAuth('signup', {
      username: 'ValidUser',
      email: 'valid@example.com',
      password,
      repeatPassword: password,
    }).password
  errors.oldPassword = check(fields.oldPassword)
  errors.password = check(fields.password)
  if (!fields.repeatPassword) errors.repeatPassword = 'Repeat your new password.'
  else if (fields.password !== fields.repeatPassword)
    errors.repeatPassword = 'Your passwords do not match.'
  if (fields.oldPassword === fields.password && fields.password)
    errors.password = 'Choose a password different from your old password.'
  return Object.fromEntries(Object.entries(errors).filter(([, value]) => value)) as typeof errors
}
export async function settingsRequest(
  path: string,
  token: string,
  method = 'GET',
  body?: unknown,
  signal?: AbortSignal,
) {
  const response = await fetch(`/api/${path}`, {
    method,
    signal,
    credentials: 'include',
    cache: 'no-store',
    headers: {
      Authorization: `Bearer ${token}`,
      ...(body ? { 'Content-Type': 'application/json' } : {}),
    },
    ...(body ? { body: JSON.stringify(body) } : {}),
  })
  if (!response.ok) {
    const detail = await response.text()
    throw new Error(
      response.status === 401
        ? 'Your session has expired. Reload and log in again.'
        : response.status === 409
          ? 'That username or email is already in use.'
          : response.status === 400 &&
              path === 'password' &&
              detail.includes('old password does not match')
            ? 'Your old password is incorrect.'
            : response.status === 400
              ? 'Check your details and try again.'
              : 'Unable to save or load your profile. Please try again.',
    )
  }
  return response.json()
}
