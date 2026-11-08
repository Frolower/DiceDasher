import { useEffect, useRef, useState, useSyncExternalStore } from 'react'
import { AuthDialog } from './components/AuthDialog'
import { SideMenu } from './components/SideMenu'
import type { ReactNode } from 'react'
import { Link, Navigate, Route, Routes, useLocation, useNavigate } from 'react-router'
import { DemoPage } from './components/DemoPage'
import { Die } from './components/Die'
import { ProfileSettings } from './components/ProfileSettings'
import { ProfilePage } from './components/ProfilePage'
import { type AuthMode, logout } from './lib/auth'
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
  const location = useLocation()
  const navigate = useNavigate()
  const activeMode =
    mode ??
    (location.pathname === '/login' && !session && !restoring && !sessionError ? 'login' : null)
  const [menuOpen, setMenuOpen] = useState(false)
  const [paused, setPaused] = useState(false)

  const logoutInFlight = useRef(false)
  const [loggingOut, setLoggingOut] = useState(false)
  const [logoutError, setLogoutError] = useState('')

  async function handleLogout(): Promise<void> {
    if (logoutInFlight.current) return
    logoutInFlight.current = true
    setLoggingOut(true)
    setLogoutError('')
    try {
      await logout()
      setSessionError('')
      setMode(null)
      setMenuOpen(false)
      navigate('/', { replace: true })
    } catch (error) {
      setLogoutError(
        error instanceof Error ? error.message : 'Unable to log out. Please try again.',
      )
    } finally {
      logoutInFlight.current = false
      setLoggingOut(false)
    }
  }

  const landing = (
    <main id="main" className="dice-stage">
      <div id="possibilities">
        <div id="systems">
          <button
            className="spinning-die"
            aria-label={paused ? 'Resume dice rotation' : 'Pause dice rotation'}
            aria-pressed={paused}
            data-paused={paused || activeMode !== null}
            onClick={() => setPaused((value) => !value)}
          >
            <Die className="center-die" />
          </button>
        </div>
      </div>
    </main>
  )
  function protect(content: ReactNode) {
    if (session) return content
    if (restoring || sessionError)
      return (
        <main id="main" className="dice-stage">
          <p role="status">{sessionError || 'Restoring your session…'}</p>
        </main>
      )
    return <Navigate to="/login" replace state={{ from: location.pathname + location.search }} />
  }
  return (
    <>
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      <header className="site-header shell">
        {session && (
          <button
            className="button quiet menu-toggle"
            aria-label="Open navigation"
            aria-haspopup="dialog"
            aria-expanded={menuOpen}
            onClick={() => setMenuOpen(true)}
          >
            <svg
              width="22"
              height="22"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.8"
              aria-hidden="true"
            >
              <path d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        )}
        <Link className="wordmark" to={session ? '/profile' : '/'} aria-label="DiceDasher home">
          <Die className="brand-die" />
          DiceDasher
        </Link>
        <div className="account-actions">
          {session ? (
            <Link
              className="button quiet"
              to="/settings"
              aria-current={location.pathname === '/settings' ? 'page' : undefined}
            >
              My profile
            </Link>
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
      <Routes>
        <Route path="/" element={session ? <Navigate to="/profile" replace /> : landing} />
        <Route path="/login" element={session ? <Navigate to="/profile" replace /> : landing} />
        <Route
          path="/profile"
          element={protect(<ProfilePage username={session?.username ?? ''} />)}
        />
        <Route path="/rooms" element={<Navigate to="/profile" replace />} />
        <Route
          path="/settings"
          element={protect(
            <ProfileSettings
              token={session?.access_token ?? ''}
              onBack={() => navigate('/profile')}
              onDeleted={() => {
                setSession(null)
                setMenuOpen(false)
                navigate('/login', { replace: true })
              }}
              onProfile={(profile) =>
                setSession(session ? { ...session, username: profile.username } : null)
              }
            />,
          )}
        />
        <Route
          path="/characters"
          element={protect(<DemoPage key="characters" kind="characters" />)}
        />
        <Route path="/friends" element={protect(<DemoPage key="friends" kind="friends" />)} />
        <Route
          path="*"
          element={
            <main id="main" className="settings-page shell">
              <h1>Page not found</h1>
              <p className="settings-intro">This page doesn’t exist.</p>
              <Link className="button primary" to={session ? '/profile' : '/'}>
                Go home
              </Link>
            </main>
          }
        />
      </Routes>
      {session && menuOpen && (
        <SideMenu
          onLogout={handleLogout}
          loggingOut={loggingOut}
          logoutError={logoutError}
          onClose={() => setMenuOpen(false)}
          onNavigate={() => {
            setMenuOpen(false)
          }}
        />
      )}
      {activeMode && (
        <AuthDialog
          key={activeMode}
          mode={activeMode}
          onClose={() => {
            setMode(null)
            if (location.pathname === '/login') navigate('/')
          }}
          onSwitch={setMode}
          onLogin={(username, data) => {
            setSession({ username, ...data })
            setMode(null)
            const from = location.state?.from
            navigate(
              typeof from === 'string' &&
                /^\/(profile|settings|characters|friends)(\?|$)/.test(from)
                ? from
                : '/profile',
              { replace: true },
            )
          }}
        />
      )}
    </>
  )
}
export default App
