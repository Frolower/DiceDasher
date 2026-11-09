import { sessionFetch, validateAuth, UserFacingError } from './auth'

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
  if (!fields.repeatPassword) errors.repeatPassword = 'text.repeatYourNewPassword'
  else if (fields.password !== fields.repeatPassword)
    errors.repeatPassword = 'text.yourPasswordsDoNotMatch'
  if (fields.oldPassword === fields.password && fields.password)
    errors.password = 'text.chooseAPasswordDifferentFromYourOldPassword'
  return Object.fromEntries(Object.entries(errors).filter(([, value]) => value)) as typeof errors
}
export async function settingsRequest(
  path: string,
  token: string,
  method = 'GET',
  body?: unknown,
  signal?: AbortSignal,
) {
  const response = await sessionFetch(`/api/${path}`, token, {
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
    throw new UserFacingError(
      response.status === 401
        ? 'text.yourSessionHasExpiredReloadAndLogInAgain'
        : response.status === 409
          ? 'text.thatUsernameOrEmailIsAlreadyInUse'
          : response.status === 400 &&
              path === 'password' &&
              detail.includes('old password does not match')
            ? 'text.yourOldPasswordIsIncorrect'
            : response.status === 400
              ? 'text.checkYourDetailsAndTryAgain'
              : method === 'DELETE'
                ? 'text.unableToDeleteYourAccountPleaseTryAgain'
                : method === 'GET'
                  ? 'text.unableToLoadYourProfilePleaseTryAgain'
                  : path === 'password'
                    ? 'text.unableToChangeYourPasswordPleaseTryAgain'
                    : 'text.unableToSaveYourProfilePleaseTryAgain',
    )
  }
  return response.json()
}
