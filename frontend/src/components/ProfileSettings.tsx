import { useEffect, useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { DeleteAccount } from './DeleteAccount'
import { passwordChecks } from '../lib/auth'
import { settingsRequest, validatePasswordChange, validateProfile } from '../lib/settings'
import type { PasswordFields, Profile } from '../lib/settings'

export function ProfileSettings({
  token,
  onProfile,
  onBack,
  onDeleted,
}: {
  token: string
  onProfile: (profile: Profile) => void
  onBack: () => void
  onDeleted: () => void
}) {
  const heading = useRef<HTMLHeadingElement>(null)
  const [profile, setProfile] = useState<Profile | null>(null)
  const [loadError, setLoadError] = useState('')
  const [attempt, setAttempt] = useState(0)
  const [passwords, setPasswords] = useState<PasswordFields>({
    oldPassword: '',
    password: '',
    repeatPassword: '',
  })
  const [errors, setErrors] = useState<Record<string, string | undefined>>({})
  const [pending, setPending] = useState('')
  const [message, setMessage] = useState({ section: '', text: '', error: false })
  const [showPassword, setShowPassword] = useState(false)
  useEffect(() => {
    heading.current?.focus()
    document.title = 'Profile settings · DiceDasher'
    return () => {
      document.title = 'DiceDasher'
    }
  }, [])
  useEffect(() => {
    const controller = new AbortController()
    settingsRequest('me', token, 'GET', undefined, controller.signal)
      .then((data: Profile) => setProfile(data))
      .catch((error: Error) => {
        if (!controller.signal.aborted) setLoadError(error.message)
      })
    return () => controller.abort()
  }, [token, attempt])
  async function save(event: FormEvent<HTMLFormElement>, section: 'account' | 'password') {
    event.preventDefault()
    if (pending || !profile) return
    const validation =
      section === 'account' ? validateProfile(profile) : validatePasswordChange(passwords)
    setErrors(validation)
    setMessage({ section: '', text: '', error: false })
    if (Object.values(validation).some(Boolean)) {
      const name = Object.keys(validation).find((key) => validation[key as keyof typeof validation])
      event.currentTarget.querySelector<HTMLInputElement>(`[name="${name}"]`)?.focus()
      return
    }
    setPending(section)
    try {
      if (section === 'account') {
        await settingsRequest('user', token, 'PATCH', profile)
        const updated = {
          username: profile.username.trim(),
          email: profile.email.trim().toLowerCase(),
        }
        setProfile(updated)
        onProfile(updated)
      } else {
        await settingsRequest('password', token, 'PUT', {
          old_password: passwords.oldPassword,
          new_password: passwords.password,
        })
        setPasswords({ oldPassword: '', password: '', repeatPassword: '' })
      }
      setMessage({
        section,
        text: section === 'account' ? 'Profile updated.' : 'Password updated.',
        error: false,
      })
    } catch (error) {
      setMessage({
        section,
        text: error instanceof Error ? error.message : 'Unable to save. Please try again.',
        error: true,
      })
    } finally {
      setPending('')
    }
  }
  function input(
    name: string,
    label: string,
    value: string,
    onChange: (value: string) => void,
    type: string,
    autoComplete: string,
  ) {
    return (
      <div className="form-field">
        <label htmlFor={`settings-${name}`}>{label}</label>
        <input
          id={`settings-${name}`}
          name={name}
          value={value}
          onChange={(event) => {
            onChange(event.target.value)
            setErrors((previous) => ({ ...previous, [name]: undefined }))
            setMessage({ section: '', text: '', error: false })
          }}
          type={type}
          autoComplete={autoComplete}
          autoCapitalize="none"
          spellCheck={false}
          required
          aria-invalid={!!errors[name]}
          aria-describedby={
            errors[name]
              ? `${name}-error`
              : name === 'password'
                ? 'settings-password-help'
                : undefined
          }
        />
        {errors[name] && (
          <small className="field-error" id={`${name}-error`}>
            {errors[name]}
          </small>
        )}
      </div>
    )
  }
  function feedback(section: string) {
    return (
      message.section === section && (
        <p
          className={message.error ? 'form-message' : 'settings-success'}
          role={message.error ? 'alert' : 'status'}
        >
          {message.text}
        </p>
      )
    )
  }
  return (
    <main id="main" className="settings-page shell">
      <button className="button quiet settings-back" onClick={onBack}>
        ← Rooms
      </button>
      <h1 ref={heading} tabIndex={-1}>
        Profile settings
      </h1>
      <p className="settings-intro">Keep your adventuring identity up to date.</p>
      {loadError ? (
        <div role="alert">
          <p>{loadError}</p>
          <button
            className="button quiet"
            onClick={() => {
              setLoadError('')
              setAttempt((value) => value + 1)
            }}
          >
            Try again
          </button>
        </div>
      ) : !profile ? (
        <p role="status">Loading your profile…</p>
      ) : (
        <div className="settings-grid">
          <section className="settings-panel" aria-labelledby="account-title">
            <h2 id="account-title">Account details</h2>
            <p>Your username is how your party knows you.</p>
            <form
              noValidate
              onSubmit={(event) => save(event, 'account')}
              aria-busy={pending === 'account'}
            >
              <fieldset disabled={!!pending}>
                {input(
                  'username',
                  'Username',
                  profile.username,
                  (username) => setProfile({ ...profile, username }),
                  'text',
                  'username',
                )}
                {input(
                  'email',
                  'Email',
                  profile.email,
                  (email) => setProfile({ ...profile, email }),
                  'email',
                  'email',
                )}
                {feedback('account')}
                <button className="button primary" type="submit">
                  {pending === 'account' ? 'Saving…' : 'Save account details'}
                </button>
              </fieldset>
            </form>
          </section>
          <section className="settings-panel" aria-labelledby="password-title">
            <h2 id="password-title">Change password</h2>
            <p>Confirm your old password, then choose a new one.</p>
            <form
              noValidate
              onSubmit={(event) => save(event, 'password')}
              aria-busy={pending === 'password'}
            >
              <fieldset disabled={!!pending}>
                {(['oldPassword', 'password', 'repeatPassword'] as const).map((name, index) => (
                  <div key={name}>
                    {input(
                      name,
                      ['Old password', 'New password', 'Repeat new password'][index],
                      passwords[name],
                      (value) => setPasswords({ ...passwords, [name]: value }),
                      showPassword ? 'text' : 'password',
                      index === 0 ? 'current-password' : 'new-password',
                    )}
                  </div>
                ))}
                <ul className="password-checks" id="settings-password-help">
                  {passwordChecks(passwords.password).map((check) => (
                    <li key={check.label} className={check.valid ? 'met' : ''}>
                      <span aria-hidden="true">{check.valid ? '✓' : '○'}</span>
                      {check.label}
                    </li>
                  ))}
                </ul>
                <label className="show-password">
                  <input
                    type="checkbox"
                    checked={showPassword}
                    onChange={(event) => setShowPassword(event.target.checked)}
                  />
                  Show passwords
                </label>
                {feedback('password')}
                <button className="button primary" type="submit">
                  {pending === 'password' ? 'Updating…' : 'Update password'}
                </button>
              </fieldset>
            </form>
          </section>
          <DeleteAccount token={token} disabled={!!pending} onDeleted={onDeleted} />
        </div>
      )}
    </main>
  )
}
