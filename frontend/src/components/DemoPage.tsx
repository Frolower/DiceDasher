import { t } from '../lib/i18n'
import { useEffect, useRef, useState } from 'react'
import { Die } from './Die'

const characters = [
  {
    name: 'text.evelynAsh',
    type: 'PC',
    system: 'VtM V5',
    description: 'text.aNocturnalStorytellerWithADangerousSecret',
    theme: 'hollow',
  },
  {
    name: 'text.thorneOakshield',
    type: 'PC',
    system: 'DnD',
    description: 'text.aSteadfastGuardianOfTheOldRoads',
    theme: 'party',
  },
  {
    name: 'text.agentMorgan',
    type: 'NPC',
    system: 'Delta Green',
    description: 'text.theContactWhoAlwaysKnowsTooMuch',
    theme: 'coast',
  },
]
const friends = [
  {
    nickname: 'MoonlitMage',
    initials: 'MM',
    detail: 'text.vtmV5StoriesAfterSundown',
    theme: 'violet',
  },
  {
    nickname: 'CriticalBard',
    initials: 'CB',
    detail: 'text.dndAlwaysATaleToTell',
    theme: 'green',
  },
  { nickname: 'AgentZero', initials: 'AZ', detail: 'text.deltaGreenTrustTheParty', theme: 'blue' },
  { nickname: 'DiceGoblin', initials: 'DG', detail: 'text.pathfinderOneMoreRoll', theme: 'violet' },
]
export function DemoPage({ kind }: { kind: 'characters' | 'friends' }) {
  const heading = useRef<HTMLHeadingElement>(null)
  const [notice, setNotice] = useState(false)
  const isCharacters = kind === 'characters'
  const title = isCharacters ? t('text.characters') : t('text.friends')
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
          <p className="profile-label">{t('text.yourTable')}</p>
          <h1 ref={heading} tabIndex={-1}>
            {title}
          </h1>
        </div>
        <button className="button primary" onClick={() => setNotice(true)}>
          <span aria-hidden="true">+</span>
          {isCharacters ? t('text.createNewCharacter') : t('text.addFriend')}
        </button>
      </div>
      <p className="settings-intro">
        {isCharacters
          ? t('text.theHeroesStrangersAndFamiliarFacesInYourStories')
          : t('text.goodCompanyMakesEveryAdventureBetter')}
      </p>
      <p className="demo-disclosure">
        <span className="preview-label">{t('text.preview')}</span>{' '}
        {isCharacters
          ? t('text.theseAreExampleCharactersYourCharactersWillLiveHere')
          : t('text.theseAreExampleFriendsYourFriendsWillAppearHere')}
      </p>
      {notice && (
        <p role="status" className="room-creation-notice">
          {isCharacters
            ? t('text.characterCreationIsComingSoonStartDreamingUpYourNextAdventurer')
            : t('text.friendRequestsAreComingSoonYouLlBeAbleToInviteYourPartyHere')}
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
                      title={
                        character.type === 'PC'
                          ? t('text.playerCharacter')
                          : t('text.nonPlayerCharacter')
                      }
                    >
                      {character.type}
                    </abbr>
                  </span>
                  <span>{t('text.example')}</span>
                </div>
                <h2>{t(character.name)}</h2>
                <p>{t(character.description)}</p>
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
                <p>{t(friend.detail)}</p>
              </div>
              <span className="friend-example">{t('text.example')}</span>
            </li>
          ))}
        </ul>
      )}
    </main>
  )
}
