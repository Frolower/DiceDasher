import { useEffect, useRef, useState } from 'react'
import { settingsRequest } from '../lib/settings'

function DeleteConfirmation({
  token,
  onClose,
  onDeleted,
}: {
  token: string
  onClose: () => void
  onDeleted: () => void
}) {
  const dialog = useRef<HTMLDialogElement>(null)
  const inFlight = useRef(false)
  const [pending, setPending] = useState(false)
  const [error, setError] = useState('')
  useEffect(() => {
    const previous = document.activeElement as HTMLElement | null
    const overflow = document.body.style.overflow
    dialog.current?.showModal()
    document.body.style.overflow = 'hidden'
    return () => {
      document.body.style.overflow = overflow
      previous?.focus()
    }
  }, [])
  async function removeAccount() {
    if (inFlight.current) return
    inFlight.current = true
    setPending(true)
    setError('')
    try {
      await settingsRequest('user', token, 'DELETE')
      onDeleted()
    } catch (cause) {
      setError(
        cause instanceof Error ? cause.message : 'Unable to delete your account. Please try again.',
      )
      inFlight.current = false
      setPending(false)
    }
  }
  return (
    <dialog
      ref={dialog}
      className="auth-dialog delete-dialog"
      aria-labelledby="delete-title"
      aria-describedby="delete-description"
      onCancel={(event) => {
        if (pending) event.preventDefault()
        else onClose()
      }}
    >
      <h2 id="delete-title">Delete your account?</h2>
      <p id="delete-description">
        Your account and all sign-in sessions will be permanently deleted. This cannot be undone.
        Character data is retained separately.
      </p>
      {error && (
        <p className="form-message" role="alert">
          {error}
        </p>
      )}
      {pending && <p role="status">Deleting your account…</p>}
      <div className="delete-dialog-actions">
        <button className="button quiet" disabled={pending} onClick={onClose}>
          Cancel
        </button>
        <button className="button destructive" disabled={pending} onClick={removeAccount}>
          {pending ? 'Deleting…' : 'Delete account'}
        </button>
      </div>
    </dialog>
  )
}

export function DeleteAccount({
  token,
  disabled,
  onDeleted,
}: {
  token: string
  disabled: boolean
  onDeleted: () => void
}) {
  const [confirming, setConfirming] = useState(false)
  return (
    <section className="settings-panel delete-account" aria-labelledby="delete-account-title">
      <div>
        <h2 id="delete-account-title">Delete account</h2>
        <p>Permanently remove your account and sign out on all devices. This cannot be undone.</p>
      </div>
      <button
        className="button destructive-outline"
        disabled={disabled}
        onClick={() => setConfirming(true)}
      >
        Delete account
      </button>
      {confirming && (
        <DeleteConfirmation
          token={token}
          onClose={() => setConfirming(false)}
          onDeleted={onDeleted}
        />
      )}
    </section>
  )
}
