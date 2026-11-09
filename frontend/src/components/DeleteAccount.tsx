import { t } from '../lib/i18n'
import { userErrorMessage } from '../lib/auth'
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
      setError(userErrorMessage(cause, 'text.unableToDeleteYourAccountPleaseTryAgain'))
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
      <h2 id="delete-title">{t('text.deleteYourAccount')}</h2>
      <p id="delete-description">
        {t(
          'text.yourAccountAndAllSignInSessionsWillBePermanentlyDeletedThisCannotBeUndoneCharacterDataIsRetainedSeparately',
        )}
      </p>
      {error && (
        <p className="form-message" role="alert">
          {t(error)}
        </p>
      )}
      {pending && <p role="status">{t('text.deletingYourAccount')}</p>}
      <div className="delete-dialog-actions">
        <button className="button quiet" disabled={pending} onClick={onClose}>
          {t('text.cancel')}
        </button>
        <button className="button destructive" disabled={pending} onClick={removeAccount}>
          {pending ? t('text.deleting') : t('text.deleteAccount')}
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
        <h2 id="delete-account-title">{t('text.deleteAccount')}</h2>
        <p>{t('text.permanentlyRemoveYourAccountAndSignOutOnAllDevicesThisCannotBeUndone')}</p>
      </div>
      <button
        className="button destructive-outline"
        disabled={disabled}
        onClick={() => setConfirming(true)}
      >
        {t('text.deleteAccount')}
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
