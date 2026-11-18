import { t } from '../../lib/i18n'
import { useEffect, useLayoutEffect, useRef, useState, type ReactNode } from 'react'
import { attributes, clans, disciplines, skills, skillName, traitById } from './catalog'
import { predatorById, predators } from './predators'
import { emptyConviction, enableRules, buy, uid, type Character, type Specialty } from './rules'
import { Errors, Field, NumberField, Rating, Select } from './Controls'
import { FreeTraits } from './TraitPicker'
import { Experience } from './Experience'
import { l, copy, type Message } from './localization'
import { useVtmLanguage } from './useVtmLanguage'
const resonanceTypes = ['choleric', 'melancholic', 'phlegmatic', 'sanguine', 'animal'] as const
function Block({
  title,
  children,
  className = '',
}: {
  title: string
  children: ReactNode
  className?: string
}) {
  return (
    <section className={`vtm-sheet-block ${className}`}>
      <h2>{l(title)}</h2>
      {children}
    </section>
  )
}
function BoxMark({ cross = false }: { cross?: boolean }) {
  return (
    <svg
      className="vtm-box-mark"
      aria-hidden="true"
      viewBox="0 0 100 100"
      preserveAspectRatio="none"
    >
      <path d={cross ? 'M0 100L100 0M0 0L100 100' : 'M0 100L100 0'} />
    </svg>
  )
}
function SquareRating({
  label,
  value,
  maximum,
  readOnly = false,
  onChange,
}: {
  label: string
  value: number
  maximum: number
  readOnly?: boolean
  onChange: (value: number) => void
}) {
  return (
    <div
      className={`vtm-tracker ${label === 'Голод' ? 'vtm-hunger' : 'vtm-tracker-centered'}`}
      role="group"
      aria-label={l(label)}
      aria-description={`${l(label)}: ${value} / ${maximum}`}
    >
      {Array.from({ length: Math.ceil(maximum / 5) }, (_, group) => (
        <span className="vtm-box-group" key={group}>
          {Array.from({ length: Math.min(5, maximum - group * 5) }, (_, offset) => {
            const i = group * 5 + offset
            return (
              <button
                type="button"
                key={i}
                aria-label={`${l(label)}: ${i + 1}`}
                aria-pressed={value >= i + 1}
                disabled={readOnly}
                onClick={() => onChange(value === i + 1 ? i : i + 1)}
              >
                {value >= i + 1 && <BoxMark cross />}
              </button>
            )
          })}
        </span>
      ))}
    </div>
  )
}
function Tracker({
  name,
  maximum,
  capacity,
  damage,
  onChange,
  readOnly,
}: {
  name: string
  maximum: number
  capacity: number
  damage: number[]
  onChange: (damage: number[]) => void
  readOnly: boolean
}) {
  const slots = Math.max(capacity, Math.min(Math.max(0, maximum), 50))
  return (
    <div
      className="vtm-tracker vtm-tracker-centered"
      role="group"
      aria-label={copy('{0}: максимум {1}', { '0': l(name), '1': maximum })}
    >
      {Array.from({ length: Math.ceil(slots / 5) }, (_, group) => (
        <span className="vtm-box-group" key={group}>
          {Array.from({ length: Math.min(5, slots - group * 5) }, (_, offset) => {
            const i = group * 5 + offset
            if (i >= maximum)
              return <span key={i} className="vtm-box-placeholder" aria-hidden="true" />
            return (
              <button
                key={i}
                disabled={readOnly}
                aria-label={copy('{0}, ячейка {1}: {2}', {
                  '0': l(name),
                  '1': i + 1,
                  '2': l(
                    damage[i] === 2
                      ? 'тяжёлый урон'
                      : damage[i] === 1
                        ? 'поверхностный урон'
                        : 'без урона',
                  ),
                })}
                onClick={() => {
                  const next = [...damage]
                  next[i] = ((next[i] ?? 0) + 1) % 3
                  onChange(next)
                }}
              >
                {!!damage[i] && <BoxMark cross={damage[i] === 2} />}
              </button>
            )
          })}
        </span>
      ))}
    </div>
  )
}
export function CharacterSheet({
  character: c,
  onChange,
  onBack,
  initialPage = 1,
  initialFit = true,
  initialEditing = false,
}: {
  character: Character
  onChange: (c: Character) => void
  onBack: () => void
  initialPage?: 1 | 2 | 3
  initialFit?: boolean
  initialEditing?: boolean
}) {
  const language = useVtmLanguage()
  const [page, setPage] = useState<number>(initialPage),
    [fit, setFit] = useState(initialFit),
    [editing, setEditing] = useState(initialEditing),
    [scale, setScale] = useState(1),
    [height, setHeight] = useState(1200),
    [errors, setErrors] = useState<Message[]>([]),
    [newDiscipline, setNewDiscipline] = useState('animalism')
  const sheet = useRef<HTMLDivElement>(null),
    viewport = useRef<HTMLDivElement>(null)
  const update = (patch: Partial<Character>) => onChange({ ...c, ...patch })
  const field = (id: string, label: string, multiline = false) => (
    <Field
      label={label}
      value={c.fields[id] ?? ''}
      onChange={(value) => update({ fields: { ...c.fields, [id]: value } })}
      multiline={multiline}
      readOnly={!editing}
    />
  )
  const identity = (key: 'name' | 'concept' | 'sire' | 'chronicle', label: string) => (
    <Field
      label={label}
      value={c[key]}
      onChange={(value) => update({ [key]: value })}
      readOnly={!editing}
    />
  )
  const locked = !editing || c.rules
  useEffect(() => {
    function measure() {
      if (!sheet.current || !viewport.current) return
      const h = sheet.current.scrollHeight
      setHeight(h)
      setScale(
        Math.min(
          1,
          (viewport.current.clientWidth - 4) / 900,
          Math.max(
            120,
            window.innerHeight - Math.max(0, viewport.current.getBoundingClientRect().top) - 24,
          ) / h,
        ),
      )
    }
    const observer = new ResizeObserver(measure)
    if (sheet.current) observer.observe(sheet.current)
    if (viewport.current) observer.observe(viewport.current)
    window.addEventListener('resize', measure)
    measure()
    return () => {
      observer.disconnect()
      window.removeEventListener('resize', measure)
    }
  }, [fit, page, language])
  useLayoutEffect(() => {
    if (page !== 2) return
    const resizeNotes = () => {
      const root = sheet.current
      const notes = root?.querySelector<HTMLTextAreaElement>('.vtm-notes-block textarea')
      const story = root?.querySelector<HTMLElement>('.vtm-background-story')
      const columns = root?.querySelector<HTMLElement>('.vtm-sheet-columns')
      const footer = root?.querySelector<HTMLElement>('.vtm-sheet-footer')
      if (!root || !notes || !story || !columns || !footer) return
      notes.style.height = '0px'
      const unit = root.getBoundingClientRect().width / root.offsetWidth || 1
      const rootTop = root.getBoundingClientRect().top
      const noteTop = (notes.getBoundingClientRect().top - rootTop) / unit
      const last = story.lastElementChild as HTMLElement
      const tail =
        (last.getBoundingClientRect().bottom - notes.getBoundingClientRect().bottom) / unit + 16
      const left = columns.firstElementChild as HTMLElement
      const style = getComputedStyle(root)
      const bottom = footer.offsetHeight + (fit ? 16 : 28) + parseFloat(style.paddingBottom) + 2
      const leftEnd = (left.getBoundingClientRect().bottom - rootTop) / unit
      const pageHeight = Math.max(parseFloat(style.minHeight) || 0, leftEnd + bottom)
      const available = Math.max(90, pageHeight - bottom - noteTop - tail)
      const content = notes.scrollHeight + 2
      notes.style.height = `${Math.min(Math.max(90, content), available)}px`
      notes.style.overflowY = content > available ? 'auto' : 'hidden'
    }
    resizeNotes()
    const observer = new ResizeObserver(resizeNotes)
    if (viewport.current) observer.observe(viewport.current)
    window.addEventListener('resize', resizeNotes)
    return () => {
      observer.disconnect()
      window.removeEventListener('resize', resizeNotes)
    }
  }, [page, fit, editing, language, c.fields, c.rituals, c.traits, c.potency, c.xp, c.spentXp])
  const stat = (id: string, name: string, type: 'attributes' | 'skills' | 'disciplines') => (
    <Rating
      key={id}
      label={name}
      value={c[type][id] ?? 0}
      readOnly={!editing}
      min={type === 'attributes' && c.rules ? 1 : 0}
      canChange={(n) => !c.rules || n > (c[type][id] ?? 0)}
      onChange={(value) => {
        if (!c.rules) {
          update({ [type]: { ...c[type], [id]: value } })
          return
        }
        let next = c
        for (let n = c[type][id] ?? 0; n < value; n++) {
          const result = buy(next, {
            kind: type === 'attributes' ? 'attribute' : type === 'skills' ? 'skill' : 'discipline',
            target: id,
          })
          if (result.error) {
            setErrors([result.error])
            return
          }
          next = result.character
        }
        setErrors([])
        onChange(next)
      }}
    />
  )
  function exportCharacter() {
    const blob = new Blob([JSON.stringify(c, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${c.name || 'vampire-v5'}.json`
    a.click()
    setTimeout(() => URL.revokeObjectURL(url), 1000)
  }
  const narrative = (list: Specialty[], onUpdate: (s: Specialty[]) => void) => (
    <div>
      <h3>{l('Специализации')}</h3>
      {list.map((s, i) => (
        <div className="vtm-specialty" key={i}>
          <Select
            readOnly={!editing}
            label="Навык"
            value={s.skill}
            onChange={(skill) => onUpdate(list.map((x, j) => (i === j ? { ...x, skill } : x)))}
            disabled={locked}
          >
            {skills.map(([id, name]) => (
              <option key={id} value={id}>
                {l(name)}
              </option>
            ))}
          </Select>
          <Field
            label="Специализация"
            value={s.auto ? l(s.text) : s.text}
            readOnly={locked}
            onChange={(text) =>
              onUpdate(list.map((x, j) => (i === j ? { ...x, text, auto: false } : x)))
            }
          />
          {!locked && (
            <button
              className="button quiet"
              aria-label={copy('Удалить специализацию {0}', { '0': i + 1 })}
              onClick={() => onUpdate(list.filter((_, j) => j !== i))}
            >
              ×
            </button>
          )}
        </div>
      ))}
      {!locked && (
        <button
          className="button quiet"
          onClick={() => onUpdate([...list, { skill: 'athletics', text: '' }])}
        >
          {l('+ Специализация')}
        </button>
      )}
    </div>
  )
  return (
    <div className="vtm-sheet-view">
      <div className="vtm-topline">
        <button className="button quiet" onClick={onBack}>
          {l('← Персонажи')}
        </button>
        <label className="vtm-check">
          <input
            type="checkbox"
            checked={c.rules}
            disabled={!editing}
            onChange={(e) => {
              if (e.target.checked) {
                const issues = enableRules(c)
                setErrors(issues)
                if (issues.length) return
              }
              setErrors([])
              update({ rules: e.target.checked })
            }}
          />
          {l('Использовать правила')}
        </label>
        <button className="button quiet" onClick={exportCharacter}>
          {l('Скачать JSON')}
        </button>
      </div>
      <Errors errors={errors} />
      <div className="vtm-sheet-toolbar">
        <button
          className="button primary"
          aria-pressed={editing}
          onClick={() => {
            setEditing(!editing)
            if (!editing) setFit(false)
            setErrors([])
          }}
        >
          {t(editing ? 'vtm.finishEditing' : 'vtm.editSheet')}
        </button>
        <div className="vtm-segment" role="group" aria-label={l('Страница чарника')}>
          {[1, 2, 3].map((p) => (
            <button className="button" key={p} aria-pressed={page === p} onClick={() => setPage(p)}>
              {p === 3 ? t('vtm.description') : l(p === 1 ? '1 · Персонаж' : '2 · Предыстория')}
            </button>
          ))}
        </div>
        {page !== 3 && (
          <button
            className="button quiet"
            aria-pressed={fit}
            onClick={() => {
              setFit(!fit)
              if (!fit)
                viewport.current?.parentElement?.scrollIntoView({
                  block: 'start',
                  behavior: 'instant',
                })
            }}
          >
            {fit ? t('vtm.closeView') : l('Вписать страницу целиком')}
          </button>
        )}
      </div>
      {editing && page !== 3 && (
        <p className="vtm-hint">
          {l(
            c.rules
              ? 'Характеристики развиваются за опыт ниже чарника. Биографию, способности и заметки можно редактировать свободно.'
              : 'Свободное заполнение: бюджеты, требования и стоимость опыта не проверяются.',
          )}
        </p>
      )}

      {page === 3 ? (
        <section className="vtm-description-page" aria-label={t('vtm.description')}>
          <Field
            label={t('vtm.description')}
            multiline
            value={c.description}
            readOnly={!editing}
            onChange={(description) => update({ description })}
          />
        </section>
      ) : (
        <div
          className={`vtm-sheet-viewport ${fit ? 'is-fit' : ''}`}
          ref={viewport}
          style={fit ? { height: height * scale + 8 } : undefined}
        >
          <div
            className={`vtm-sheet ${page === 2 ? 'vtm-sheet-background' : ''}`}
            ref={sheet}
            style={fit ? { transform: `scale(${scale})` } : undefined}
          >
            <header className="vtm-sheet-header">
              <div className="vtm-sheet-brand">
                VAMPIRE<small>THE MASQUERADE · V5</small>
              </div>
              <span>
                {l('Лист персонажа ·')} {page} / 2
              </span>
            </header>
            {page === 1 ? (
              <>
                <div className="vtm-grid vtm-identity">
                  {identity('name', 'Имя')}
                  {identity('concept', 'Концепция')}
                  {identity('sire', 'Сир')}
                  {field('player', 'Игрок')}
                  {identity('chronicle', 'Хроника')}
                  <Select
                    readOnly={!editing}
                    label="Клан"
                    value={c.clan}
                    onChange={(clan) => update({ clan })}
                    disabled={locked}
                  >
                    <option value="">{l('Не выбран')}</option>
                    {Object.entries(clans).map(([id, clan]) => (
                      <option key={id} value={id}>
                        {l(clan.name)}
                      </option>
                    ))}
                  </Select>
                  {field('ambition', 'Амбиция')}
                  <Select
                    readOnly={!editing}
                    label="Хищник"
                    value={c.predator}
                    disabled={locked}
                    onChange={(predator) => update({ predator })}
                  >
                    <option value="">{l('Без типа хищника')}</option>
                    {predators.map((p) => (
                      <option key={p.id} value={p.id}>
                        {l(p.name)}
                      </option>
                    ))}
                  </Select>
                  <NumberField
                    label="Поколение"
                    fallback={10}
                    value={c.generation}
                    readOnly={locked}
                    onChange={(generation) => update({ generation })}
                  />
                </div>
                {editing && c.rules && <p className="vtm-hint">{t('vtm.sheetDotsHelp')}</p>}
                <Block title="Характеристики">
                  <div className="vtm-stat-columns">
                    {['Физические', 'Социальные', 'Ментальные'].map((name, i) => (
                      <section key={i}>
                        <h3>{l(name)}</h3>
                        {attributes
                          .slice(i * 3, i * 3 + 3)
                          .map(([id, label]) => stat(id, label, 'attributes'))}
                      </section>
                    ))}
                  </div>
                </Block>
                <Block title="Навыки">
                  <div className="vtm-stat-columns">
                    {['Физические', 'Социальные', 'Ментальные'].map((name, i) => (
                      <section key={i}>
                        <h3>{l(name)}</h3>
                        {skills
                          .slice(i * 9, i * 9 + 9)
                          .map(([id, label]) => stat(id, label, 'skills'))}
                      </section>
                    ))}
                  </div>
                  {fit ? (
                    <p className="vtm-specialties-preview">
                      {c.specialties.map((s, i) => (
                        <span key={i}>
                          {l(skillName(s.skill))}: {s.auto ? l(s.text) : s.text}
                        </span>
                      ))}
                    </p>
                  ) : (
                    narrative(c.specialties, (specialties) => update({ specialties }))
                  )}
                </Block>
                <div className="vtm-grid vtm-trackers">
                  <Block title="Здоровье">
                    {!locked && (
                      <NumberField
                        label="Максимум здоровья"
                        value={c.health}
                        readOnly={locked}
                        onChange={(health) => update({ health })}
                      />
                    )}
                    <Tracker
                      name="Здоровье"
                      maximum={c.health}
                      capacity={13}
                      damage={c.healthDamage}
                      readOnly={false}
                      onChange={(healthDamage) => update({ healthDamage })}
                    />
                  </Block>
                  <Block title="Сила воли">
                    {!locked && (
                      <NumberField
                        label="Максимум силы воли"
                        value={c.willpower}
                        readOnly={locked}
                        onChange={(willpower) => update({ willpower })}
                      />
                    )}
                    <Tracker
                      name="Сила воли"
                      maximum={c.willpower}
                      capacity={10}
                      damage={c.willpowerDamage}
                      readOnly={false}
                      onChange={(willpowerDamage) => update({ willpowerDamage })}
                    />
                  </Block>
                  <Block title="Человечность">
                    <SquareRating
                      label="Человечность"
                      value={c.humanity}
                      maximum={10}
                      readOnly={false}
                      onChange={(humanity) => update({ humanity })}
                    />
                  </Block>
                </div>
                <div className="vtm-grid vtm-narrative">
                  <Block title="Заповеди хроники">{field('tenets', 'Заповеди', true)}</Block>
                  <Block title="Опоры и убеждения">
                    {c.convictions.map((v, i) => (
                      <div key={i} className="vtm-sheet-conviction">
                        <Field
                          label={t('vtm.conviction')}
                          value={v.belief}
                          multiline
                          readOnly={!editing}
                          onChange={(belief) =>
                            update({
                              convictions: c.convictions.map((x, j) =>
                                j === i ? { ...x, belief } : x,
                              ),
                            })
                          }
                        />
                        <Field
                          label="Имя опоры"
                          value={v.name}
                          readOnly={!editing}
                          onChange={(name) =>
                            update({
                              convictions: c.convictions.map((x, j) =>
                                j === i ? { ...x, name } : x,
                              ),
                            })
                          }
                        />
                        <Field
                          label="Описание опоры"
                          value={v.description}
                          multiline
                          readOnly={!editing}
                          onChange={(description) =>
                            update({
                              convictions: c.convictions.map((x, j) =>
                                j === i ? { ...x, description } : x,
                              ),
                            })
                          }
                        />
                        {!locked && (
                          <button
                            className="button quiet"
                            onClick={() =>
                              update({ convictions: c.convictions.filter((_, j) => j !== i) })
                            }
                          >
                            {l('Удалить пару')}
                          </button>
                        )}
                      </div>
                    ))}
                    {!locked && (
                      <button
                        className="button quiet"
                        onClick={() =>
                          update({ convictions: [...c.convictions, emptyConviction()] })
                        }
                      >
                        {l('+ Пара')}
                      </button>
                    )}
                  </Block>
                  <Block title="Проклятие клана">
                    {field('bane', 'Проклятие', true)}
                    {field('desire', 'Желание')}
                  </Block>
                </div>
                <Block title="Дисциплины">
                  <div className="vtm-discipline-grid">
                    {Object.entries(c.disciplines)
                      .filter(([, n]) => n > 0)
                      .map(([id]) => (
                        <section key={id}>
                          {stat(id, disciplines[id] ? l(disciplines[id]) : id, 'disciplines')}
                          <Field
                            label="Способности"
                            multiline
                            value={c.powers[id] ?? ''}
                            readOnly={!editing}
                            onChange={(text) => update({ powers: { ...c.powers, [id]: text } })}
                          />
                        </section>
                      ))}
                  </div>
                  {!Object.values(c.disciplines).some((n) => n > 0) && (
                    <p className="vtm-hint">
                      {l(
                        c.clan === 'thinBlood'
                          ? 'Дисциплины доступны через мериты тонкокровых.'
                          : 'Дисциплины пока не выбраны.',
                      )}
                    </p>
                  )}
                  {!locked && (
                    <div className="vtm-grid">
                      <Select
                        readOnly={!editing}
                        label="Добавить дисциплину"
                        value={newDiscipline}
                        onChange={setNewDiscipline}
                      >
                        {Object.entries(disciplines).map(([id, name]) => (
                          <option key={id} value={id}>
                            {l(name)}
                          </option>
                        ))}
                      </Select>
                      <button
                        className="button quiet"
                        onClick={() =>
                          update({
                            disciplines: {
                              ...c.disciplines,
                              [newDiscipline]: Math.max(1, c.disciplines[newDiscipline] ?? 0),
                            },
                          })
                        }
                      >
                        {l('Добавить дисциплину')}
                      </button>
                    </div>
                  )}
                </Block>
                <div className="vtm-feeding-row">
                  <div className="vtm-resonance">
                    <Select
                      label="Резонанс"
                      value={c.fields.resonance ?? ''}
                      readOnly={!editing}
                      onChange={(resonance) => update({ fields: { ...c.fields, resonance } })}
                    >
                      <option value="">{t('vtm.noResonance')}</option>
                      {resonanceTypes.map((id) => (
                        <option key={id} value={id}>
                          {t('vtm.resonance.' + id)}
                        </option>
                      ))}
                      {c.fields.resonance &&
                        !resonanceTypes.some((id) => id === c.fields.resonance) && (
                          <option value={c.fields.resonance}>{c.fields.resonance}</option>
                        )}
                    </Select>
                  </div>
                  <div className="vtm-field vtm-hunger-field">
                    <span>{l('Голод')}</span>
                    <SquareRating
                      label="Голод"
                      value={c.hunger}
                      maximum={5}
                      onChange={(hunger) => update({ hunger })}
                    />
                  </div>
                </div>
              </>
            ) : (
              <>
                <div className="vtm-sheet-columns">
                  <div>
                    {(['merit', 'background', 'flaw'] as const).map((kind) => (
                      <Block
                        key={kind}
                        title={
                          kind === 'merit'
                            ? 'Мериты'
                            : kind === 'background'
                              ? 'Предыстории'
                              : 'Недостатки'
                        }
                      >
                        {c.traits
                          .filter(
                            (entry) =>
                              (traitById[entry.id]?.kind ?? entry.kind ?? 'background') === kind,
                          )
                          .map((entry) =>
                            locked ? (
                              <TraitRow key={entry.key} entry={entry} />
                            ) : (
                              <div className="vtm-trait-edit-row" key={entry.key}>
                                <Field
                                  label={t('vtm.name.3de498')}
                                  value={
                                    entry.name ??
                                    (traitById[entry.id] ? l(traitById[entry.id].name) : entry.id)
                                  }
                                  onChange={(name) =>
                                    update({
                                      traits: c.traits.map((v) =>
                                        v.key === entry.key ? { ...v, name } : v,
                                      ),
                                    })
                                  }
                                />
                                <NumberField
                                  label={t('vtm.rating.9f0b99')}
                                  value={entry.value}
                                  max={5}
                                  onChange={(value) =>
                                    update({
                                      traits: c.traits.map((v) =>
                                        v.key === entry.key ? { ...v, value } : v,
                                      ),
                                    })
                                  }
                                />
                                <button
                                  type="button"
                                  className="button quiet"
                                  aria-label={t('vtm.removeTrait', {
                                    name:
                                      entry.name ??
                                      (traitById[entry.id]
                                        ? l(traitById[entry.id].name)
                                        : entry.id),
                                  })}
                                  onClick={() =>
                                    update({ traits: c.traits.filter((v) => v.key !== entry.key) })
                                  }
                                >
                                  ×
                                </button>
                              </div>
                            ),
                          )}
                        {!locked && (
                          <button
                            type="button"
                            className="button quiet"
                            aria-label={t('vtm.addTraitRow', {
                              kind: l(
                                kind === 'merit'
                                  ? 'Мериты'
                                  : kind === 'background'
                                    ? 'Предыстории'
                                    : 'Недостатки',
                              ),
                            })}
                            onClick={() =>
                              update({
                                traits: [
                                  ...c.traits,
                                  {
                                    key: uid(),
                                    id: 'custom-' + uid(),
                                    kind,
                                    name: '',
                                    value: 1,
                                    detail: '',
                                    source: 'player',
                                    floor: 0,
                                  },
                                ],
                              })
                            }
                          >
                            {t('vtm.addRow')}
                          </button>
                        )}
                      </Block>
                    ))}
                    <Block title="Бросок охоты">
                      {field('huntingRoll', 'Пул охоты')}
                      <p className="vtm-hint">
                        {predatorById[c.predator] ? l(predatorById[c.predator].name) : ''}
                      </p>
                    </Block>
                    <Block title="Сила крови">
                      <Rating
                        label="Сила крови"
                        value={c.potency}
                        max={10}
                        readOnly={!editing}
                        canChange={(value) => !c.rules || value > c.potency}
                        onChange={(potency) => {
                          if (!c.rules) {
                            update({ potency })
                            return
                          }
                          let next = c
                          for (let n = c.potency; n < potency; n++) {
                            const result = buy(next, { kind: 'potency', target: '' })
                            if (result.error) {
                              setErrors([result.error])
                              return
                            }
                            next = result.character
                          }
                          setErrors([])
                          onChange(next)
                        }}
                      />
                      <div className="vtm-grid">
                        {field('bloodSurge', 'Прилив крови')}
                        {field('mend', 'Исцеление')}
                        {field('powerBonus', 'Бонус силы')}
                        {field('rouseReroll', 'Переброс Пробуждения')}
                        {field('feedingPenalty', 'Ограничения питания')}
                        {field('baneSeverity', 'Тяжесть проклятия')}
                      </div>
                    </Block>
                    <Block title="Опыт">
                      <div className="vtm-grid">
                        <NumberField
                          label="Всего опыта"
                          value={c.xp}
                          readOnly={locked}
                          onChange={(xp) => update({ xp })}
                        />
                        <NumberField
                          label="Потрачено опыта"
                          value={c.spentXp}
                          readOnly={locked}
                          onChange={(spentXp) => update({ spentXp })}
                        />
                      </div>
                    </Block>
                    <Block title="Оружие">{field('weapons', 'Оружие и урон', true)}</Block>
                  </div>
                  <div className="vtm-background-story">
                    <Block title="Биография">
                      <div className="vtm-grid">
                        {field('trueAge', 'Настоящий возраст')}
                        {field('apparentAge', 'Видимый возраст')}
                        {field('birth', 'Дата рождения')}
                        {field('death', 'Дата смерти')}
                      </div>
                      {field('appearance', 'Внешность', true)}
                      {field('features', 'Отличительные черты', true)}
                      {field('history', 'История', true)}
                    </Block>
                    <Block title="Заметки" className="vtm-notes-block">
                      {field('notes', 'Заметки', true)}
                    </Block>
                    <Block title="Ритуалы и формулы" className="vtm-rituals-block">
                      {c.rituals.map((r, i) => (
                        <div key={i} className="vtm-ritual">
                          <strong>
                            {l(r.kind === 'ritual' ? 'Ритуал' : 'Формула')} {l('· Уровень')}{' '}
                            {r.level}
                          </strong>
                          <Field
                            label="Название ритуала или формулы"
                            value={r.name}
                            readOnly={!editing}
                            onChange={(name) =>
                              update({
                                rituals: c.rituals.map((x, j) => (j === i ? { ...x, name } : x)),
                              })
                            }
                          />
                          {!locked && (
                            <>
                              <NumberField
                                label="Уровень"
                                value={r.level}
                                onChange={(level) =>
                                  update({
                                    rituals: c.rituals.map((x, j) =>
                                      j === i ? { ...x, level } : x,
                                    ),
                                  })
                                }
                              />
                              <button
                                className="button quiet"
                                onClick={() =>
                                  update({ rituals: c.rituals.filter((_, j) => j !== i) })
                                }
                              >
                                {l('Удалить')}
                              </button>
                            </>
                          )}
                        </div>
                      ))}
                      {!locked && (
                        <div className="vtm-topline">
                          {(['ritual', 'formula'] as const).map((kind) => (
                            <button
                              className="button quiet"
                              key={kind}
                              onClick={() =>
                                update({ rituals: [...c.rituals, { kind, name: '', level: 1 }] })
                              }
                            >
                              + {l(kind === 'ritual' ? 'Ритуал' : 'Формула')}
                            </button>
                          ))}
                        </div>
                      )}
                    </Block>
                    {field('otherPowers', 'Другие способности и церемонии', true)}
                  </div>
                </div>
              </>
            )}
            <footer className="vtm-sheet-footer">
              <span>{c.name || l('Новый Сородич')}</span>
              <span>{l('Vampire: The Masquerade · Пятая редакция')}</span>
            </footer>
          </div>
        </div>
      )}
      {editing && page === 2 && !c.rules && (
        <section className="vtm-editor-panel">
          <h2>{l('Свободное редактирование особенностей')}</h2>
          <Select
            readOnly={!editing}
            label="Дисциплина сродства тонкокрового"
            value={c.affinity}
            onChange={(affinity) => update({ affinity })}
          >
            <option value="">{l('Нет сродства')}</option>
            {Object.entries(disciplines)
              .filter(([id]) => id !== 'alchemy')
              .map(([id, name]) => (
                <option key={id} value={id}>
                  {l(name)}
                </option>
              ))}
          </Select>
          <FreeTraits entries={c.traits} onChange={(traits) => update({ traits })} />
        </section>
      )}
      {editing && page !== 3 && c.rules && <Experience character={c} onChange={onChange} />}
    </div>
  )
}
function TraitRow({ entry: t }: { entry: Character['traits'][number] }) {
  const def = traitById[t.id]
  return (
    <div className="vtm-sheet-trait">
      <div>
        <strong>{t.name ?? (def ? l(def.name) : t.id)}</strong>
        {t.detail && <small>{t.autoDetail ? l(t.detail) : t.detail}</small>}
        {t.source === 'predator' && <small>{l('От хищника')}</small>}
      </div>
      <span
        className={def?.thin || def?.discipline ? undefined : 'vtm-trait-dots'}
        role="img"
        aria-label={`${t.name ?? (def ? l(def.name) : t.id)}: ${t.value}`}
      >
        {def?.thin || def?.discipline
          ? '✓'
          : Array.from({ length: Math.max(0, Math.min(t.value, 50)) }, (_, i) => (
              <span className="vtm-trait-dot-slot" key={i} aria-hidden="true" />
            ))}
      </span>
    </div>
  )
}
