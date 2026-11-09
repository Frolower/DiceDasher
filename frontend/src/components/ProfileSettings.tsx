import { t } from '../lib/i18n'
import { useEffect, useRef, useState } from 'react'
import { useSearchParams } from 'react-router'
import { ThemeSettings } from './ThemeSettings'
import type { FormEvent } from 'react'
import { DeleteAccount } from './DeleteAccount'
import { passwordChecks, userErrorMessage } from '../lib/auth'
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
  const [searchParams, setSearchParams] = useSearchParams()
  const section = searchParams.get('section') === 'themes' ? 'themes' : 'profile'
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
  const title = t(section === 'themes' ? 'text.themes' : 'text.profileSettings')
  useEffect(() => {
    heading.current?.focus()
    document.title = `${title} · DiceDasher`
    return () => {
      document.title = 'DiceDasher'
    }
  }, [title])
  useEffect(() => {
    const controller = new AbortController()
    settingsRequest('me', token, 'GET', undefined, controller.signal)
      .then((data: Profile) => setProfile(data))
      .catch((error: unknown) => {
        if (!controller.signal.aborted)
          setLoadError(userErrorMessage(error, 'text.unableToLoadYourProfilePleaseTryAgain'))
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
        text: section === 'account' ? 'text.profileUpdated' : 'text.passwordUpdated',
        error: false,
      })
    } catch (error) {
      setMessage({
        section,
        text: userErrorMessage(
          error,
          section === 'password'
            ? 'text.unableToChangeYourPasswordPleaseTryAgain'
            : 'text.unableToSaveYourProfilePleaseTryAgain',
        ),
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
            {t(errors[name] ?? '')}
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
          {t(message.text)}
        </p>
      )
    )
  }
  return (
    <main id="main" className="settings-page shell">
      <button className="button quiet settings-back" onClick={onBack}>
        {t('navigation.backToRooms')}
      </button>
      <h1 ref={heading} tabIndex={-1}>
        {t('text.settings')}
      </h1>
      <p className="settings-intro">{t('text.makeYourselfAtHomeBeforeYourNextAdventure')}</p>
      <div className="settings-layout">
        <nav className="settings-nav" aria-label={t('text.settingsSections')}>
          {(['profile', 'themes'] as const).map((item) => (
            <button
              key={item}
              type="button"
              className="settings-nav-item"
              aria-current={section === item ? 'page' : undefined}
              onClick={() =>
                setSearchParams((previous) => {
                  const next = new URLSearchParams(previous)
                  if (item === 'profile') next.delete('section')
                  else next.set('section', item)
                  return next
                })
              }
            >
              <span aria-hidden="true">
                {item === 'profile' ? (
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor">
                    <circle cx="12" cy="7" r="4" />
                    <path d="M4 21v-2a8 8 0 0 1 16 0v2H4Z" />
                  </svg>
                ) : (
                  '◐'
                )}
              </span>
              {item === 'profile' ? t('text.profile') : t('text.themes')}
            </button>
          ))}
        </nav>
        <div className="settings-content">
          {section === 'themes' && <ThemeSettings />}
          <div hidden={section !== 'profile'}>
            <h2 className="settings-section-title">{t('text.profile')}</h2>
            <p className="settings-section-intro">
              {t('text.keepYourAdventuringIdentityUpToDate')}
            </p>
            {loadError ? (
              <div role="alert">
                <p>{t(loadError)}</p>
                <button
                  className="button quiet"
                  onClick={() => {
                    setLoadError('')
                    setAttempt((value) => value + 1)
                  }}
                >
                  {t('text.tryAgain')}
                </button>
              </div>
            ) : !profile ? (
              <p role="status">{t('text.loadingYourProfile')}</p>
            ) : (
              <div className="settings-grid">
                <section className="settings-panel" aria-labelledby="account-title">
                  <h2 id="account-title">{t('text.accountDetails')}</h2>
                  <p>{t('text.yourUsernameIsHowYourPartyKnowsYou')}</p>
                  <form
                    noValidate
                    onSubmit={(event) => save(event, 'account')}
                    aria-busy={pending === 'account'}
                  >
                    <fieldset disabled={!!pending}>
                      {input(
                        'username',
                        t('text.username'),
                        profile.username,
                        (username) => setProfile({ ...profile, username }),
                        'text',
                        'username',
                      )}
                      {input(
                        'email',
                        t('text.email'),
                        profile.email,
                        (email) => setProfile({ ...profile, email }),
                        'email',
                        'email',
                      )}
                      {feedback('account')}
                      <button className="button primary" type="submit">
                        {pending === 'account' ? t('text.saving') : t('text.saveAccountDetails')}
                      </button>
                    </fieldset>
                  </form>
                </section>
                <section className="settings-panel" aria-labelledby="password-title">
                  <h2 id="password-title">{t('text.changePassword')}</h2>
                  <p>{t('text.confirmYourOldPasswordThenChooseANewOne')}</p>
                  <form
                    noValidate
                    onSubmit={(event) => save(event, 'password')}
                    aria-busy={pending === 'password'}
                  >
                    <fieldset disabled={!!pending}>
                      {(['oldPassword', 'password', 'repeatPassword'] as const).map(
                        (name, index) => (
                          <div key={name}>
                            {input(
                              name,
                              [
                                t('text.oldPassword'),
                                t('text.newPassword'),
                                t('text.repeatNewPassword'),
                              ][index],
                              passwords[name],
                              (value) => setPasswords({ ...passwords, [name]: value }),
                              showPassword ? 'text' : 'password',
                              index === 0 ? 'current-password' : 'new-password',
                            )}
                          </div>
                        ),
                      )}
                      <ul className="password-checks" id="settings-password-help">
                        {passwordChecks(passwords.password).map((check) => (
                          <li key={check.label} className={check.valid ? 'met' : ''}>
                            <span aria-hidden="true">{check.valid ? '✓' : '○'}</span>
                            {t(check.label)}
                          </li>
                        ))}
                      </ul>
                      <label className="show-password">
                        <input
                          type="checkbox"
                          checked={showPassword}
                          onChange={(event) => setShowPassword(event.target.checked)}
                        />
                        {t('text.showPasswords')}
                      </label>
                      {feedback('password')}
                      <button className="button primary" type="submit">
                        {pending === 'password' ? t('text.updating') : t('text.updatePassword')}
                      </button>
                    </fieldset>
                  </form>
                </section>
                <DeleteAccount token={token} disabled={!!pending} onDeleted={onDeleted} />
              </div>
            )}
          </div>
        </div>
      </div>
    </main>
  )
}
