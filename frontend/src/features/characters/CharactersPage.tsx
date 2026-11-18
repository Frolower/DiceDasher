import { useEffect, useRef, useState, useSyncExternalStore } from 'react'
import { getLanguage, subscribeLanguage, t } from '../../lib/i18n'
import { SelectField } from '../../components/SelectField'
import { clans } from '../vtm/catalog'
import { l } from '../vtm/localization'
import { completeDraft, newDraft, type Character } from '../vtm/rules'
import { CharacterSheet } from '../vtm/CharacterSheet'
import { CreationWizard } from '../vtm/CreationWizard'
import '../vtm/vtm.css'
import './characters.css'
// Memory-only: no character request is made to the backend.
const characters = new Map<string, Character[]>()
const empty: Character[] = []
const listeners = new Set<() => void>()
const subscribe = (listener: () => void) => {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}
const snapshot = (owner: string) => characters.get(owner) ?? empty
function put(owner: string, character: Character) {
  const items = snapshot(owner)
  const next = items.some((c) => c.id === character.id)
    ? items.map((c) => (c.id === character.id ? character : c))
    : [...items, character]
  characters.set(owner, next)
  listeners.forEach((fn) => fn())
}
export function CharactersPage({ owner = 'preview' }: { owner?: string }) {
  const language = useSyncExternalStore(subscribeLanguage, getLanguage, getLanguage)
  const items = useSyncExternalStore(
    subscribe,
    () => snapshot(owner),
    () => snapshot(owner),
  )
  const [view, setView] = useState<'list' | 'wizard' | 'sheet'>('list'),
    [selected, setSelected] = useState(''),
    [setup, setSetup] = useState(false)
  const heading = useRef<HTMLHeadingElement>(null)
  useEffect(() => {
    document.title = `${t('characters.title')} · DiceDasher`
    heading.current?.focus()
    return () => {
      document.title = 'DiceDasher'
    }
  }, [view, language])
  const character = items.find((c) => c.id === selected)
  const complete = (c: Character) => {
    put(owner, c)
    setSelected(c.id)
    setView('sheet')
  }
  return (
    <main
      id="main"
      lang={language}
      className={view === 'list' ? 'characters-page shell' : 'vtm-page'}
    >
      {view === 'list' ? (
        <>
          <div className="demo-heading">
            <div>
              <p className="profile-label">{t('characters.collection')}</p>
              <h1 ref={heading} tabIndex={-1}>
                {t('characters.title')}
              </h1>
            </div>
            <button className="button primary" onClick={() => setSetup(true)}>
              {t('characters.create')}
            </button>
          </div>
          <p className="vtm-intro">{t('characters.intro')}</p>
          <p className="vtm-local-notice">{t('characters.temporary')}</p>
          {items.length ? (
            <ul className="character-list">
              {items.map((c) => (
                <li key={c.id}>
                  <button
                    onClick={() => {
                      setSelected(c.id)
                      setView('sheet')
                    }}
                  >
                    <span className="character-card-system">Vampire: The Masquerade V5</span>
                    <strong>{c.name || t('characters.unnamed')}</strong>
                    <span className="character-card-meta">
                      {clans[c.clan] ? l(clans[c.clan].name) : t('characters.noClan')} ·{' '}
                      {t('characters.generation', { value: c.generation })}
                    </span>
                    {c.chronicle && (
                      <span className="character-card-chronicle">
                        {t('characters.chronicle')}: {c.chronicle}
                      </span>
                    )}
                    <span className="character-card-description">{c.description || c.concept}</span>
                    <span className="character-card-open" aria-hidden="true">
                      ↗
                    </span>
                  </button>
                </li>
              ))}
            </ul>
          ) : (
            <section className="characters-empty">
              <CharacterIcon />
              <h2>{t('characters.emptyTitle')}</h2>
              <p>{t('characters.emptyBody')}</p>
              <button className="button primary" onClick={() => setSetup(true)}>
                {t('characters.first')}
              </button>
            </section>
          )}
        </>
      ) : view === 'wizard' ? (
        <CreationWizard onComplete={complete} onBack={() => setView('list')} />
      ) : character ? (
        <CharacterSheet
          character={character}
          onChange={(c) => put(owner, c)}
          onBack={() => setView('list')}
        />
      ) : null}
      {setup && (
        <Setup
          onClose={() => setSetup(false)}
          onCreate={(rules) => {
            setSetup(false)
            if (rules) setView('wizard')
            else complete(completeDraft(newDraft(), false))
          }}
        />
      )}
    </main>
  )
}
function Setup({ onClose, onCreate }: { onClose: () => void; onCreate: (rules: boolean) => void }) {
  const dialog = useRef<HTMLDialogElement>(null),
    [rules, setRules] = useState(true)
  useEffect(() => {
    const previous = document.activeElement as HTMLElement | null
    dialog.current?.showModal()
    return () => previous?.focus()
  }, [])
  return (
    <dialog
      className="auth-dialog vtm-setup"
      ref={dialog}
      aria-labelledby="vtm-setup-title"
      onCancel={onClose}
    >
      <button className="close-button" onClick={onClose} aria-label={t('characters.close')}>
        ×
      </button>
      <p className="vtm-eyebrow">{t('characters.newStory')}</p>
      <h2 id="vtm-setup-title">{t('characters.create')}</h2>
      <SelectField
        label={t('characters.system')}
        value="vtm-v5"
        onChange={() => {}}
        options={[{ value: 'vtm-v5', label: 'Vampire: The Masquerade V5' }]}
      />
      <label className="vtm-check vtm-rules-option">
        <input type="checkbox" checked={rules} onChange={(e) => setRules(e.target.checked)} />
        <span>
          <strong>{t('characters.rules')}</strong>
          <small>{t('characters.rulesHelp')}</small>
        </span>
      </label>
      <p className="vtm-hint">{rules ? t('characters.guidedHelp') : t('characters.freeHelp')}</p>
      <div className="vtm-wizard-actions">
        <button className="button quiet" onClick={onClose}>
          {t('characters.cancel')}
        </button>
        <button className="button primary" onClick={() => onCreate(rules)}>
          {rules ? t('characters.start') : t('characters.blank')}
        </button>
      </div>
    </dialog>
  )
}

function CharacterIcon() {
  return (
    <svg
      className="character-icon"
      aria-hidden="true"
      width="48"
      height="48"
      viewBox="0 0 32 32"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.4"
    >
      <path d="M19 4H8a2 2 0 0 0-2 2v20a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V11Z" />
      <path d="M19 4v7h7M11 22h10" />
      <circle cx="14" cy="15" r="2.5" />
    </svg>
  )
}
