import { useEffect, useRef, useState } from 'react'
import { Die } from './Die'

const characters = [
  {
    name: 'Evelyn Ash',
    type: 'PC',
    system: 'VtM V5',
    description: 'A nocturnal storyteller with a dangerous secret.',
    theme: 'hollow',
  },
  {
    name: 'Thorne Oakshield',
    type: 'PC',
    system: 'DnD',
    description: 'A steadfast guardian of the old roads.',
    theme: 'party',
  },
  {
    name: 'Agent Morgan',
    type: 'NPC',
    system: 'Delta Green',
    description: 'The contact who always knows too much.',
    theme: 'coast',
  },
]
const friends = [
  {
    nickname: 'MoonlitMage',
    initials: 'MM',
    detail: 'VtM V5 · Stories after sundown',
    theme: 'violet',
  },
  {
    nickname: 'CriticalBard',
    initials: 'CB',
    detail: 'DnD · Always a tale to tell',
    theme: 'green',
  },
  { nickname: 'AgentZero', initials: 'AZ', detail: 'Delta Green · Trust the party', theme: 'blue' },
  { nickname: 'DiceGoblin', initials: 'DG', detail: 'Pathfinder · One more roll', theme: 'violet' },
]
export function DemoPage({ kind }: { kind: 'characters' | 'friends' }) {
  const heading = useRef<HTMLHeadingElement>(null)
  const [notice, setNotice] = useState(false)
  const isCharacters = kind === 'characters'
  const title = isCharacters ? 'Characters' : 'Friends'
  useEffect(() => {
    document.title = `${title} · DiceDasher`
    heading.current?.focus()
    return () => {
      document.title = 'DiceDasher'
    }
  }, [title])
  return (
    <main id="main" className="demo-page shell">
      <div className="demo-heading">
        <div>
          <p className="profile-label">Your table</p>
          <h1 ref={heading} tabIndex={-1}>
            {title}
          </h1>
        </div>
        <button className="button primary" onClick={() => setNotice(true)}>
          <span aria-hidden="true">+</span>
          {isCharacters ? 'Create new character' : 'Add friend'}
        </button>
      </div>
      <p className="settings-intro">
        {isCharacters
          ? 'The heroes, strangers, and familiar faces in your stories.'
          : 'Good company makes every adventure better.'}
      </p>
      <p className="demo-disclosure">
        <span className="preview-label">Preview</span>{' '}
        {isCharacters
          ? 'These are example characters. Your characters will live here.'
          : 'These are example friends. Your friends will appear here.'}
      </p>
      {notice && (
        <p role="status" className="room-creation-notice">
          {isCharacters
            ? 'Character creation is coming soon. Start dreaming up your next adventurer.'
            : 'Friend requests are coming soon. You’ll be able to invite your party here.'}
        </p>
      )}
      {isCharacters ? (
        <ul className="campaign-grid">
          {characters.map((character) => (
            <li key={character.name} className={`campaign-tile ${character.theme}`}>
              <div className="campaign-art" aria-hidden="true">
                <Die />
              </div>
              <div className="campaign-content">
                <div className="campaign-meta">
                  <span>
                    <abbr
                      title={character.type === 'PC' ? 'Player character' : 'Non-player character'}
                    >
                      {character.type}
                    </abbr>
                  </span>
                  <span>Example</span>
                </div>
                <h2>{character.name}</h2>
                <p>{character.description}</p>
                <div className="campaign-system">{character.system}</div>
              </div>
            </li>
          ))}
        </ul>
      ) : (
        <ul className="friend-list">
          {friends.map((friend) => (
            <li className="friend-row" key={friend.nickname}>
              <div className={`friend-avatar ${friend.theme}`} aria-hidden="true">
                {friend.initials}
              </div>
              <div className="friend-identity">
                <h2>{friend.nickname}</h2>
                <p>{friend.detail}</p>
              </div>
              <span className="friend-example">Example</span>
            </li>
          ))}
        </ul>
      )}
    </main>
  )
}
