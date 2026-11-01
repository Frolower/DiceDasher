import { useEffect, useRef } from 'react'

export type Page = 'rooms' | 'settings' | 'characters' | 'friends'
export function SideMenu({
  page,
  onNavigate,
  onClose,
}: {
  page: Page
  onNavigate: (page: Page) => void
  onClose: () => void
}) {
  const dialog = useRef<HTMLDialogElement>(null)
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
  return (
    <dialog
      ref={dialog}
      className="side-menu"
      aria-labelledby="menu-title"
      onCancel={onClose}
      onClick={(event) => {
        if (event.target === event.currentTarget) {
          const rect = event.currentTarget.getBoundingClientRect()
          if (event.clientX > rect.right || event.clientX < rect.left) onClose()
        }
      }}
    >
      <div className="side-menu-heading">
        <h2 id="menu-title">Your table</h2>
        <button className="button quiet" onClick={onClose} aria-label="Close navigation">
          ×
        </button>
      </div>
      <nav aria-label="Main navigation">
        {(['characters', 'friends'] as const).map((item) => (
          <button
            key={item}
            className="side-menu-item"
            aria-current={page === item ? 'page' : undefined}
            onClick={() => onNavigate(item)}
          >
            <svg
              className="side-menu-icon"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.7"
              strokeLinecap="round"
              strokeLinejoin="round"
              aria-hidden="true"
              focusable="false"
            >
              {item === 'characters' ? (
                <>
                  <path d="M14 3H6a1 1 0 0 0-1 1v16a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V8Z" />
                  <path d="M14 3v5h5M13 12h3M8 17h8" />
                  <circle cx="9" cy="10.5" r="1.5" />
                  <path d="M6.8 14a2.3 2.3 0 0 1 4.4 0" />
                </>
              ) : (
                <>
                  <circle cx="9" cy="8" r="3" />
                  <path d="M3 20v-2a6 6 0 0 1 12 0v2ZM16 5a3 3 0 0 1 0 6M18 14a5 5 0 0 1 3 4v2h-3" />
                </>
              )}
            </svg>
            {item === 'characters' ? 'Characters' : 'Friends'}
            <span aria-hidden="true">→</span>
          </button>
        ))}
      </nav>
      <p className="side-menu-note">Every adventure starts with good company.</p>
    </dialog>
  )
}
