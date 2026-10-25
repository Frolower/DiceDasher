import { useState } from 'react'
import { AuthDialog } from './components/AuthDialog'
import { Die } from './components/Die'
import type { AuthMode } from './lib/auth'
import './App.css'

function App() {
  const [mode, setMode] = useState<AuthMode | null>(null)
  const [session, setSession] = useState<{ username: string; access_token: string; expires_at: string } | null>(null)
  const [paused, setPaused] = useState(false)
  return <>
    <a className="skip-link" href="#main">Skip to content</a>
    <header className="site-header shell">
      <a className="wordmark" href="#" aria-label="DiceDasher home"><Die className="brand-die" />DiceDasher</a>
      <div className="account-actions">{session ? <span className="signed-in" role="status">Welcome, {session.username}</span> : <><button className="button quiet" onClick={() => setMode('login')}>Log in</button><button className="button primary" onClick={() => setMode('signup')}>Create an account <span aria-hidden="true">↗</span></button></>}</div>
    </header>
    <main id="main" className="dice-stage">
      <div id="possibilities">
        <div id="systems">
          <button
            className="spinning-die"
            aria-label={paused ? 'Resume dice rotation' : 'Pause dice rotation'}
            aria-pressed={paused}
            data-paused={paused || mode !== null}
            onClick={() => setPaused(value => !value)}
          >
            <Die className="center-die" />
          </button>
        </div>
      </div>
    </main>
    {mode && <AuthDialog key={mode} mode={mode} onClose={() => setMode(null)} onSwitch={setMode} onLogin={(username, data) => { setSession({ username, ...data }); setMode(null) }} />}
  </>
}
export default App
