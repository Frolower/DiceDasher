import { useEffect, useState, useSyncExternalStore } from 'react'
import { AuthDialog } from './components/AuthDialog'
import { Die } from './components/Die'
import { ProfileSettings } from './components/ProfileSettings'
import { ProfilePage } from './components/ProfilePage'
import type { AuthMode } from './lib/auth'
import { getSession, refreshSession, setSession, subscribeSession } from './lib/auth'
import './App.css'

function App() {
  const [mode, setMode] = useState<AuthMode | null>(null)
  const session = useSyncExternalStore(subscribeSession, getSession)
  const [restoring, setRestoring] = useState(true)
  const [sessionError, setSessionError] = useState('')
  const [attempt, setAttempt] = useState(0)
  useEffect(() => {
    let active = true
    let timer: ReturnType<typeof setTimeout>
    async function restore() {
      try {
        await refreshSession()
        if (active) setSessionError('')
      } catch {
        if (active) {
          setSessionError('Unable to reconnect to your account. Retrying…')
          timer = setTimeout(restore, 5000)
        }
      } finally {
        if (active) setRestoring(false)
      }
    }
    void restore()
    return () => {
      active = false
      clearTimeout(timer)
    }
  }, [attempt])
  useEffect(() => {
    if (!session) return
    let timer: ReturnType<typeof setTimeout>
    let active = true
    async function renew() {
      try {
        await refreshSession()
        if (active) setSessionError('')
      } catch {
        if (active) timer = setTimeout(renew, 5000)
      }
    }
    timer = setTimeout(renew, Math.max(0, Date.parse(session.expires_at) - Date.now() - 60000))
    const onVisible = () => {
      if (
        document.visibilityState === 'visible' &&
        Date.parse(session.expires_at) - Date.now() < 60000
      )
        void renew()
    }
    document.addEventListener('visibilitychange', onVisible)
    return () => {
      active = false
      clearTimeout(timer)
      document.removeEventListener('visibilitychange', onVisible)
    }
  }, [session])
  const [settings, setSettings] = useState(false)
  const [paused, setPaused] = useState(false)
  return (
    <>
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      <header className="site-header shell">
        <a
          className="wordmark"
          href="#"
          onClick={() => setSettings(false)}
          aria-label="DiceDasher home"
        >
          <Die className="brand-die" />
          DiceDasher
        </a>
        <div className="account-actions">
          {session ? (
            <button
              className="button quiet"
              onClick={() => setSettings(true)}
              aria-current={settings ? 'page' : undefined}
            >
              My profile
            </button>
          ) : restoring ? (
            <p role="status">Restoring your session…</p>
          ) : sessionError ? (
            <button className="button quiet" onClick={() => setAttempt((value) => value + 1)}>
              Reconnect
            </button>
          ) : (
            <>
              <button className="button quiet" onClick={() => setMode('login')}>
                Log in
              </button>
              <button className="button primary" onClick={() => setMode('signup')}>
                Create an account <span aria-hidden="true">↗</span>
              </button>
            </>
          )}
        </div>
      </header>
      {sessionError && (
        <p className="shell" role="status">
          {sessionError}
        </p>
      )}
      {session ? (
        <>
          {settings ? (
            <ProfileSettings
              token={session.access_token}
              onBack={() => setSettings(false)}
              onProfile={(profile) =>
                setSession(session ? { ...session, username: profile.username } : null)
              }
            />
          ) : (
            <ProfilePage username={session.username} />
          )}
        </>
      ) : (
        <main id="main" className="dice-stage">
          <div id="possibilities">
            <div id="systems">
              <button
                className="spinning-die"
                aria-label={paused ? 'Resume dice rotation' : 'Pause dice rotation'}
                aria-pressed={paused}
                data-paused={paused || mode !== null}
                onClick={() => setPaused((value) => !value)}
              >
                <Die className="center-die" />
              </button>
            </div>
          </div>
        </main>
      )}
      {mode && (
        <AuthDialog
          key={mode}
          mode={mode}
          onClose={() => setMode(null)}
          onSwitch={setMode}
          onLogin={(username, data) => {
            setSession({ username, ...data })
            setMode(null)
          }}
        />
      )}
    </>
  )
}
export default App
