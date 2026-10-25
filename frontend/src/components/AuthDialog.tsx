import { useEffect, useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { authenticate, passwordChecks, validateAuth } from '../lib/auth'
import type { AuthFields, AuthMode, FieldErrors } from '../lib/auth'
import { Die } from './Die'

export function AuthDialog({
  mode,
  onClose,
  onSwitch,
  onLogin,
}: {
  mode: AuthMode
  onClose: () => void
  onSwitch: (mode: AuthMode) => void
  onLogin: (username: string, session: { access_token: string; expires_at: string }) => void
}) {
  const dialog = useRef<HTMLDialogElement>(null)
  const request = useRef<AbortController | null>(null)
  const [fields, setFields] = useState<AuthFields>({
    email: '',
    username: '',
    password: '',
    repeatPassword: '',
  })
  const [errors, setErrors] = useState<FieldErrors>({})
  const [touched, setTouched] = useState<Partial<Record<keyof AuthFields, boolean>>>({})
  const [showPassword, setShowPassword] = useState(false)
  const [pending, setPending] = useState(false)
  const [message, setMessage] = useState('')
  const [created, setCreated] = useState(false)
  const signup = mode === 'signup'
  useEffect(() => {
    const previousFocus = document.activeElement as HTMLElement | null
    dialog.current?.showModal()
    dialog.current?.querySelector<HTMLInputElement>('input')?.focus()
    const previousOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => {
      request.current?.abort()
      document.body.style.overflow = previousOverflow
      previousFocus?.focus()
    }
  }, [])
  function update(name: keyof AuthFields, value: string) {
    const next = { ...fields, [name]: value }
    setFields(next)
    const validation = validateAuth(mode, next)
    setErrors(
      Object.fromEntries(
        Object.entries(validation).filter(([key]) => touched[key as keyof AuthFields]),
      ),
    )
    setMessage('')
  }
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (pending) return
    const validation = validateAuth(mode, fields)
    setErrors(validation)
    setTouched({ email: true, username: true, password: true, repeatPassword: true })
    if (Object.keys(validation).length) {
      const firstInvalid = (['email', 'username', 'password', 'repeatPassword'] as const).find(
        (name) => validation[name],
      )
      dialog.current?.querySelector<HTMLInputElement>(`[name="${firstInvalid}"]`)?.focus()
      return
    }
    setPending(true)
    setMessage('')
    const controller = new AbortController()
    request.current = controller
    const timeout = window.setTimeout(() => controller.abort(), 15000)
    try {
      const session = await authenticate(mode, fields, controller.signal)
      if (session) onLogin(fields.username, session)
      else {
        setCreated(true)
        setFields({ email: '', username: '', password: '', repeatPassword: '' })
      }
    } catch (error) {
      setMessage(
        controller.signal.aborted
          ? 'The request timed out. Please try again.'
          : error instanceof Error
            ? error.message
            : 'Unable to connect. Check your connection and try again.',
      )
    } finally {
      window.clearTimeout(timeout)
      setPending(false)
    }
  }
  function field(name: keyof AuthFields, label: string, placeholder: string, type = 'text') {
    return (
      <div className="form-field">
        <label htmlFor={name}>{label}</label>
        <input
          id={name}
          name={name}
          type={type}
          placeholder={placeholder}
          value={fields[name]}
          autoComplete={
            name === 'email'
              ? 'email'
              : name === 'username'
                ? 'username'
                : signup
                  ? 'new-password'
                  : 'current-password'
          }
          autoCapitalize="none"
          spellCheck={false}
          required
          aria-invalid={!!errors[name]}
          aria-describedby={
            errors[name]
              ? `${name}-error`
              : name === 'password' && signup
                ? 'password-help'
                : undefined
          }
          onChange={(event) => update(name, event.target.value)}
          onBlur={() => {
            setTouched((previous) => ({ ...previous, [name]: true }))
            setErrors((previous) => ({ ...previous, [name]: validateAuth(mode, fields)[name] }))
          }}
        />
        {errors[name] && (
          <small className="field-error" id={`${name}-error`}>
            {errors[name]}
          </small>
        )}
      </div>
    )
  }
  return (
    <dialog
      ref={dialog}
      className="auth-dialog"
      aria-labelledby="auth-title"
      onCancel={onClose}
      onClick={(event) => {
        if (event.target === event.currentTarget) {
          const rect = event.currentTarget.getBoundingClientRect()
          if (
            event.clientX < rect.left ||
            event.clientX > rect.right ||
            event.clientY < rect.top ||
            event.clientY > rect.bottom
          )
            onClose()
        }
      }}
    >
      <button className="close-button" aria-label="Close account form" onClick={onClose}>
        ×
      </button>
      {created ? (
        <div className="created-state">
          <span className="success-mark">✓</span>
          <h2 id="auth-title">Your adventure starts here.</h2>
          <p>Your account is ready. Log in to continue to DiceDasher.</p>
          <button className="button primary full" onClick={() => onSwitch('login')}>
            Log in <span aria-hidden="true">↗</span>
          </button>
        </div>
      ) : (
        <>
          <div className="auth-heading">
            <span className="dialog-mark">
              <Die />
            </span>
            <h2 id="auth-title">{signup ? 'Create an account' : 'Log in'}</h2>
          </div>
          <form noValidate onSubmit={submit} aria-busy={pending}>
            <fieldset disabled={pending}>
              {signup && field('email', 'Email', 'you@example.com', 'email')}
              {field('username', 'Username', 'Your adventuring alias')}
              {field(
                'password',
                'Password',
                signup ? 'Make it a strong one' : 'Enter your password',
                showPassword ? 'text' : 'password',
              )}
              {signup && (
                <ul className="password-checks" id="password-help">
                  {passwordChecks(fields.password).map((check) => (
                    <li key={check.label} className={check.valid ? 'met' : ''}>
                      <span aria-hidden="true">{check.valid ? '✓' : '○'}</span> {check.label}
                    </li>
                  ))}
                </ul>
              )}
              {signup &&
                field(
                  'repeatPassword',
                  'Repeat password',
                  'One more time',
                  showPassword ? 'text' : 'password',
                )}
              <label className="show-password">
                <input
                  type="checkbox"
                  checked={showPassword}
                  onChange={(event) => setShowPassword(event.target.checked)}
                />{' '}
                Show {signup ? 'passwords' : 'password'}
              </label>
              {message && (
                <p className="form-message" role="alert">
                  {message}
                </p>
              )}
              <button className="button primary full" type="submit">
                {pending ? 'Please wait…' : signup ? 'Create an account' : 'Log in'}
                <span aria-hidden="true">↗</span>
              </button>
            </fieldset>
          </form>
          <p className="switch-auth">
            {signup ? 'Already have an account?' : 'New around here?'}{' '}
            <button disabled={pending} onClick={() => onSwitch(signup ? 'login' : 'signup')}>
              {signup ? 'Log in' : 'Create an account'}
            </button>
          </p>
        </>
      )}
    </dialog>
  )
}
