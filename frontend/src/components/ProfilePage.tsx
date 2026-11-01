import { useEffect, useRef } from 'react'
import { Die } from './Die'

const exampleCampaigns = [
  {
    name: 'The Ember Coast',
    kind: 'Campaign',
    description: 'A voyage beyond the edge of the map.',
    system: 'Dungeons & Dragons',
    theme: 'coast',
  },
  {
    name: 'Sunday adventurers',
    kind: 'Room',
    description: 'A place for the party to gather.',
    system: 'Tabletop hangout',
    theme: 'party',
  },
  {
    name: 'Secrets of the Hollow',
    kind: 'Campaign',
    description: 'Some mysteries are better shared.',
    system: 'Pathfinder',
    theme: 'hollow',
  },
]

export function ProfilePage({ username }: { username: string }) {
  const heading = useRef<HTMLHeadingElement>(null)

  useEffect(() => {
    document.title = `${username} · DiceDasher`
    heading.current?.focus()
    return () => {
      document.title = 'DiceDasher'
    }
  }, [username])

  return (
    <main id="main" className="profile-page shell">
      <section className="profile-intro" aria-labelledby="profile-title">
        <div className="profile-avatar" aria-hidden="true">
          {username.slice(0, 1).toUpperCase()}
        </div>
        <div className="profile-identity">
          <p className="profile-label">Your profile</p>
          <h1 id="profile-title" ref={heading} tabIndex={-1}>
            {username}
          </h1>
          <p>Your party. Your stories. Your next adventure.</p>
        </div>
        <Die className="profile-die" />
      </section>

      <section
        className="profile-campaigns"
        aria-labelledby="campaigns-title"
        aria-describedby="campaigns-note"
      >
        <div className="campaigns-heading">
          <h2 id="campaigns-title">Rooms & campaigns</h2>
          <span className="preview-label">Preview</span>
        </div>
        <p id="campaigns-note" className="campaigns-note">
          Your rooms and campaigns will live here. These are examples while we get things ready.
        </p>
        <ul className="campaign-grid">
          {exampleCampaigns.map((campaign) => (
            <li className={`campaign-tile ${campaign.theme}`} key={campaign.name}>
              <div className="campaign-art" aria-hidden="true">
                <Die />
              </div>
              <div className="campaign-content">
                <div className="campaign-meta">
                  <span>{campaign.kind}</span>
                  <span>Example</span>
                </div>
                <h3>{campaign.name}</h3>
                <p>{campaign.description}</p>
                <div className="campaign-system">{campaign.system}</div>
              </div>
            </li>
          ))}
        </ul>
        <p className="campaigns-footer">
          More adventures are on the horizon. Room and campaign features are coming soon.
        </p>
      </section>
    </main>
  )
}
