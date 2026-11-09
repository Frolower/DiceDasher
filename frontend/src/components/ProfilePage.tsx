import { t } from '../lib/i18n'
import { useEffect, useRef, useState } from 'react'
import { Die } from './Die'

const exampleCampaigns = [
  {
    name: 'text.theEmberCoast',
    kind: 'text.room',
    description: 'text.aVoyageBeyondTheEdgeOfTheMap',
    system: 'Dungeons & Dragons',
    theme: 'coast',
  },
  {
    name: 'text.sundayAdventurers',
    kind: 'text.room',
    description: 'text.aPlaceForThePartyToGather',
    system: 'text.tabletopHangout',
    theme: 'party',
  },
  {
    name: 'text.secretsOfTheHollow',
    kind: 'text.room',
    description: 'text.someMysteriesAreBetterShared',
    system: 'Pathfinder',
    theme: 'hollow',
  },
]

export function ProfilePage({ username }: { username: string }) {
  const [creationNotice, setCreationNotice] = useState(false)
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
          <p className="profile-label">{t('text.yourProfile')}</p>
          <h1 id="profile-title" ref={heading} tabIndex={-1}>
            {username}
          </h1>
          <p>{t('text.yourPartyYourStoriesYourNextAdventure')}</p>
        </div>
        <Die className="profile-die" />
      </section>

      <section
        className="profile-campaigns"
        aria-labelledby="campaigns-title"
        aria-describedby="campaigns-note"
      >
        <div className="campaigns-heading">
          <h2 id="campaigns-title">{t('text.rooms')}</h2>
          <span className="preview-label">{t('text.preview')}</span>
          <button
            className="button primary create-room-button"
            onClick={() => setCreationNotice(true)}
          >
            <span aria-hidden="true">+</span> {t('text.createNewRoom')}
          </button>
        </div>
        <p id="campaigns-note" className="campaigns-note">
          {t('text.yourRoomsWillLiveHereTheseAreExamplesWhileWeGetThingsReady')}
        </p>
        {creationNotice && (
          <p className="room-creation-notice" role="status">
            {t('text.roomCreationIsComingSoonYouLlBeAbleToGatherYourPartyHere')}
          </p>
        )}
        <ul className="campaign-grid">
          {exampleCampaigns.map((campaign) => (
            <li className={`campaign-tile ${campaign.theme}`} key={campaign.name}>
              <div className="campaign-art" aria-hidden="true">
                <Die />
              </div>
              <div className="campaign-content">
                <div className="campaign-meta">
                  <span>{t(campaign.kind)}</span>
                  <span>{t('text.example')}</span>
                </div>
                <h3>{t(campaign.name)}</h3>
                <p>{t(campaign.description)}</p>
                <div className="campaign-system">{t(campaign.system)}</div>
              </div>
            </li>
          ))}
        </ul>
        <p className="campaigns-footer">
          {t('text.moreAdventuresAreOnTheHorizonRoomFeaturesAreComingSoon')}
        </p>
      </section>
    </main>
  )
}
