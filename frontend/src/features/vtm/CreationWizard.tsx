import { t } from '../../lib/i18n'
import { useEffect, useRef, useState } from 'react'
import { attributes, clans, disciplines, skills, skillName, traitById } from './catalog'
import { predators, predatorById, predatorDisciplines } from './predators'
import {
  type Character,
  type Draft,
  basePotency,
  completeDraft,
  creationDisciplines,
  emptyConviction,
  generations,
  newDraft,
  predatorGifts,
  sum,
  validateStep,
  canAssignRating,
  attributeQuotas,
  disciplineQuotas,
  skillQuotas,
  changeDistribution,
} from './rules'
import { Errors, Field, NumberField, Rating, Select } from './Controls'
import { TraitPicker } from './TraitPicker'
import { l, copy, type Message } from './localization'
import { useVtmLanguage } from './useVtmLanguage'
const steps = [
  'Личность',
  'Характеристики',
  'Навыки',
  'Дисциплины',
  'Тип хищника',
  'Преимущества',
  'Убеждения',
]
const intros = [
  'Кем вы были — и кем стали после Обращения?',
  'Одна характеристика на 4, три на 3, четыре на 2 и одна на 1. Ноль означает, что значение ещё не выбрано.',
  'Выберите набор и распределите его точки между навыками.',
  'Две точки в одну дисциплину, одна — в другую.',
  'Способ охоты оставляет свой след: выберите бонусы и специализацию.',
  'Добавьте свои преимущества и недостатки поверх полученного от хищника.',
  'Каждое убеждение связывает вас с живым человеком — вашей опорой.',
]
export function CreationWizard({
  onComplete,
  onBack,
  initialDraft,
  initialStep = 0,
}: {
  onComplete: (c: Character) => void
  onBack: () => void
  initialDraft?: Draft
  initialStep?: number
}) {
  useVtmLanguage()
  const [draft, setDraft] = useState<Draft>(
      () =>
        initialDraft ?? {
          ...newDraft(),
          attributes: Object.fromEntries(attributes.map(([id]) => [id, 0])),
        },
    ),
    [step, setStep] = useState(Math.min(6, Math.max(0, initialStep))),
    [errors, setErrors] = useState<Message[]>([])
  const heading = useRef<HTMLHeadingElement>(null)
  useEffect(() => {
    heading.current?.focus()
  }, [step])
  const update = (patch: Partial<Draft>) => {
    setDraft((d) => ({ ...d, ...patch }))
    setErrors([])
  }
  const p = predatorById[draft.predator]
  function next() {
    const problems = validateStep(draft, step)
    setErrors(problems)
    if (problems.length) return
    if (step < 6) setStep(step + 1)
    else {
      const bad = Array.from({ length: 7 }, (_, i) => ({ i, errors: validateStep(draft, i) })).find(
        (s) => s.errors.length,
      )
      if (bad) {
        setStep(bad.i)
        setErrors(bad.errors)
      } else onComplete(completeDraft(draft, true))
    }
  }
  function identity(patch: Partial<Draft>) {
    const d = { ...draft, ...patch }
    if (d.clan === 'thinBlood') d.age = 'childer'
    if (!generations(d.age, d.clan).includes(d.generation))
      d.generation = generations(d.age, d.clan)[0]
    if (patch.clan !== undefined && patch.clan !== draft.clan) {
      d.disciplines = {}
      d.predatorDiscipline = ''
      d.affinity = ''
    }
    setDraft(d)
    setErrors([])
  }
  return (
    <div className="vtm-workspace">
      <div className="vtm-topline">
        <button className="button quiet" onClick={onBack}>
          {l('← Персонажи')}
        </button>
      </div>
      <div className="vtm-wizard-layout">
        <aside className="vtm-steps">
          <p className="vtm-eyebrow">Vampire: The Masquerade V5</p>
          <nav aria-label={l('Шаги создания')}>
            <ol>
              {steps.map((name, i) => (
                <li key={i}>
                  <button
                    aria-current={step === i ? 'step' : undefined}
                    disabled={i > step}
                    onClick={() => {
                      setStep(i)
                      setErrors([])
                    }}
                  >
                    <span aria-hidden="true">{i < step ? '✓' : i + 1}</span>
                    {l(name)}
                  </button>
                </li>
              ))}
            </ol>
          </nav>
          <div className="vtm-character-summary">
            <strong>{draft.name || l('Новый Сородич')}</strong>
            <span>{clans[draft.clan] ? l(clans[draft.clan].name) : l('Клан ещё не выбран')}</span>
            <span>
              {l('Здоровье')} {draft.attributes.stamina + 3} · {l('Сила воли')}{' '}
              {draft.attributes.composure + draft.attributes.resolve}
            </span>
            <span>
              {l('Сила крови')} {basePotency(draft) + (p?.potency ?? 0)}
            </span>
          </div>
        </aside>
        <section className="vtm-wizard-content">
          <p className="vtm-eyebrow">
            {l('Шаг')} {step + 1} {l('из 7')}
          </p>
          <h1 ref={heading} tabIndex={-1}>
            {l(steps[step])}
          </h1>
          <p className="vtm-intro">{l(intros[step])}</p>
          {step === 0 && (
            <>
              <div className="vtm-grid">
                <Field
                  label="Имя персонажа"
                  value={draft.name}
                  onChange={(name) => update({ name })}
                />
                <Field
                  label="Концепция"
                  value={draft.concept}
                  onChange={(concept) => update({ concept })}
                />
                <Field
                  label="Хроника"
                  value={draft.chronicle}
                  onChange={(chronicle) => update({ chronicle })}
                />
                <Field label="Сир" value={draft.sire} onChange={(sire) => update({ sire })} />
                <Select label="Клан" value={draft.clan} onChange={(clan) => identity({ clan })}>
                  <option value="">{l('Выберите клан')}</option>
                  {Object.entries(clans).map(([id, c]) => (
                    <option key={id} value={id}>
                      {l(c.name)}
                    </option>
                  ))}
                </Select>
                <Select
                  label="Возрастная категория"
                  value={draft.age}
                  onChange={(age) => identity({ age: age as Draft['age'] })}
                >
                  <option value="childer">{l('Дитя')}</option>
                  <option value="neonate" disabled={draft.clan === 'thinBlood'}>
                    {l('Неонат')}
                  </option>
                  <option value="ancilla" disabled={draft.clan === 'thinBlood'}>
                    {l('Анцилла')}
                  </option>
                </Select>
                <Select
                  label="Поколение"
                  value={String(draft.generation)}
                  onChange={(n) => update({ generation: Number(n) })}
                >
                  {generations(draft.age, draft.clan).map((n) => (
                    <option key={n} value={n}>
                      {n}
                    </option>
                  ))}
                </Select>
              </div>
              <Field
                label={t('vtm.description')}
                multiline
                value={draft.description}
                onChange={(description) => update({ description })}
              />
              <div className="vtm-note">
                {l(
                  draft.age === 'childer'
                    ? 'Человечность 7 · Без стартового опыта'
                    : draft.age === 'neonate'
                      ? 'Человечность 7 · 15 опыта'
                      : 'Человечность 6 · 35 опыта · +2 точки преимуществ и +2 точки недостатков',
                )}
              </div>
            </>
          )}
          {step === 1 && (
            <>
              <div className="vtm-budget">
                {[4, 3, 2, 1].map((v) => (
                  <span key={v}>
                    {l('На')} {v}{' '}
                    <strong>
                      {Object.values(draft.attributes).filter((n) => n === v).length} /{' '}
                      {attributeQuotas[v]}
                    </strong>
                  </span>
                ))}
              </div>
              <div className="vtm-stat-columns">
                {['Физические', 'Социальные', 'Ментальные'].map((name, i) => (
                  <section key={i}>
                    <h3>{l(name)}</h3>
                    {attributes.slice(i * 3, i * 3 + 3).map(([id, label]) => (
                      <Rating
                        key={id}
                        label={label}
                        value={draft.attributes[id]}
                        min={0}
                        max={4}
                        clearable
                        canChange={(n) => canAssignRating(draft.attributes, id, n, attributeQuotas)}
                        onChange={(n) => {
                          if (canAssignRating(draft.attributes, id, n, attributeQuotas))
                            update({ attributes: { ...draft.attributes, [id]: n } })
                        }}
                      />
                    ))}
                  </section>
                ))}
              </div>
            </>
          )}
          {step === 2 && (
            <>
              <div className="vtm-distribution">
                {(
                  [
                    ['jack', 'Мастер на все руки', '1 × 3 · 8 × 2 · 10 × 1'],
                    ['balanced', 'Сбалансированный', '3 × 3 · 5 × 2 · 7 × 1'],
                    ['specialist', 'Специалист', '1 × 4 · 3 × 3 · 3 × 2 · 3 × 1'],
                  ] as const
                ).map(([id, name, desc]) => (
                  <label key={id} className={draft.distribution === id ? 'chosen' : ''}>
                    <input
                      type="radio"
                      name="distribution"
                      value={id}
                      checked={draft.distribution === id}
                      onChange={() => {
                        setDraft((d) => changeDistribution(d, id))
                        setErrors([])
                      }}
                    />
                    <strong>{l(name)}</strong>
                    <small>{desc}</small>
                  </label>
                ))}
              </div>
              <div className="vtm-budget">
                {[4, 3, 2, 1].map((v) => (
                  <span key={v}>
                    {l('На')} {v}{' '}
                    <strong>
                      {Object.values(draft.skills).filter((n) => n === v).length} /{' '}
                      {skillQuotas(draft.distribution)[v]}
                    </strong>
                  </span>
                ))}
              </div>
              <div className="vtm-stat-columns">
                {['Физические', 'Социальные', 'Ментальные'].map((name, i) => (
                  <section key={i}>
                    <h3>{l(name)}</h3>
                    {skills.slice(i * 9, i * 9 + 9).map(([id, label]) => (
                      <Rating
                        key={id}
                        label={label}
                        value={draft.skills[id]}
                        max={draft.distribution === 'specialist' ? 4 : 3}
                        clearable
                        canChange={(n) =>
                          canAssignRating(draft.skills, id, n, skillQuotas(draft.distribution))
                        }
                        onChange={(n) => {
                          if (canAssignRating(draft.skills, id, n, skillQuotas(draft.distribution)))
                            update({ skills: { ...draft.skills, [id]: n } })
                        }}
                      />
                    ))}
                  </section>
                ))}
              </div>
              <h2>{l('Специализации')}</h2>
              <p className="vtm-hint">
                {l(
                  'Бесплатно по одной для изученных Гуманитарных наук, Ремесла, Исполнения и Естественных наук, плюс одна на выбор.',
                )}
              </p>
              <div className="vtm-grid">
                {['academics', 'craft', 'performance', 'science']
                  .filter((id) => draft.skills[id] > 0)
                  .map((id) => (
                    <Field
                      key={id}
                      label={skillName(id)}
                      value={draft.specialties[id] ?? ''}
                      onChange={(text) =>
                        update({ specialties: { ...draft.specialties, [id]: text } })
                      }
                    />
                  ))}
              </div>
              <div className="vtm-grid">
                <Select
                  label="Навык для свободной специализации"
                  value={draft.extraSpecialty.skill}
                  onChange={(skill) =>
                    update({ extraSpecialty: { ...draft.extraSpecialty, skill } })
                  }
                >
                  <option value="">{l('Выберите навык')}</option>
                  {skills
                    .filter(([id]) => draft.skills[id] > 0)
                    .map(([id, name]) => (
                      <option key={id} value={id}>
                        {l(name)}
                      </option>
                    ))}
                </Select>
                <Field
                  label="Свободная специализация"
                  value={draft.extraSpecialty.text}
                  onChange={(text) => update({ extraSpecialty: { ...draft.extraSpecialty, text } })}
                />
              </div>
            </>
          )}
          {step === 3 && (
            <>
              {draft.clan === 'thinBlood' ? (
                <div className="vtm-note">
                  {l(
                    'Тонкокровые получают дисциплину или алхимию через мериты на шаге «Преимущества».',
                  )}
                </div>
              ) : (
                <div className="vtm-stat-columns">
                  {(draft.clan === 'caitiff'
                    ? Object.keys(disciplines).filter((id) => id !== 'alchemy')
                    : (clans[draft.clan]?.disciplines ?? [])
                  ).map((id) => (
                    <Rating
                      key={id}
                      label={disciplines[id]}
                      max={2}
                      clearable
                      value={draft.disciplines[id] ?? 0}
                      canChange={(n) => canAssignRating(draft.disciplines, id, n, disciplineQuotas)}
                      onChange={(n) => {
                        if (canAssignRating(draft.disciplines, id, n, disciplineQuotas))
                          update({ disciplines: { ...draft.disciplines, [id]: n } })
                      }}
                    />
                  ))}
                </div>
              )}
              <p className="vtm-hint">
                {l('Названия способностей, ритуалов и формул вы сможете вписать в готовый чарник.')}
              </p>
            </>
          )}
          {step === 4 && (
            <>
              <Select
                label="Тип хищника"
                value={draft.predator}
                onChange={(predator) =>
                  update({
                    predator,
                    predatorDiscipline: '',
                    predatorSpecialty: 0,
                    predatorDetail: '',
                    choices: [],
                    pools: {},
                  })
                }
              >
                <option value="">
                  {l(
                    draft.age === 'childer' || draft.clan === 'thinBlood'
                      ? 'Без типа хищника — без бонусов'
                      : 'Выберите тип хищника',
                  )}
                </option>
                {predators.map((p) => (
                  <option
                    key={p.id}
                    value={p.id}
                    disabled={
                      (draft.clan === 'ventrue' && ['bagger', 'farmer'].includes(p.id)) ||
                      (p.id === 'titheCollector' && draft.age !== 'ancilla')
                    }
                  >
                    {l(p.name)}
                  </option>
                ))}
              </Select>
              {p && (
                <>
                  <div className="vtm-grid">
                    <Select
                      label="Бонусная дисциплина: +1"
                      value={draft.predatorDiscipline}
                      onChange={(predatorDiscipline) => update({ predatorDiscipline })}
                    >
                      <option value="">{l('Выберите дисциплину')}</option>
                      {predatorDisciplines(p.id, draft.clan).map((id) => (
                        <option key={id} value={id}>
                          {l(disciplines[id])}
                        </option>
                      ))}
                    </Select>
                    <Select
                      label="Специализация хищника"
                      value={String(draft.predatorSpecialty)}
                      onChange={(v) => update({ predatorSpecialty: Number(v) })}
                    >
                      {p.specialties.map(([skill, text], i) => (
                        <option key={i} value={i}>
                          {l(skillName(skill))}: {l(text)}
                        </option>
                      ))}
                    </Select>
                    <Field
                      label="Уточнение специализации (необязательно)"
                      value={draft.predatorDetail}
                      onChange={(predatorDetail) => update({ predatorDetail })}
                    />
                  </div>
                  {p.humanity !== undefined && (
                    <p className="vtm-note">
                      {l('Человечность:')} {p.humanity > 0 ? '+' : ''}
                      {p.humanity}
                      {p.potency ? l(' · Сила крови: +1') : ''}
                    </p>
                  )}
                  {p.choices?.map((options, i) => (
                    <Select
                      key={`${p.id}-${i}`}
                      label={copy('Бонус на выбор {0}', { '0': i + 1 })}
                      value={String(draft.choices[i] ?? 0)}
                      onChange={(v) => {
                        const choices = [...draft.choices]
                        choices[i] = Number(v)
                        update({ choices })
                      }}
                    >
                      {options.map((g, j) => (
                        <option key={j} value={j}>
                          {l(traitById[g.id].name)} · {g.value}
                        </option>
                      ))}
                    </Select>
                  ))}
                  {p.pools?.map((pool) => (
                    <section key={pool.label}>
                      <h3>
                        {l(pool.label)} · {sum(pool.ids.map((id) => draft.pools[id] ?? 0))} /{' '}
                        {pool.dots}
                      </h3>
                      <div className="vtm-grid">
                        {pool.ids.map((id) => (
                          <NumberField
                            key={id}
                            label={traitById[id].name}
                            value={draft.pools[id] ?? 0}
                            max={pool.dots}
                            onChange={(v) => update({ pools: { ...draft.pools, [id]: v } })}
                          />
                        ))}
                      </div>
                    </section>
                  ))}
                  <div className="vtm-note">
                    <strong>{l('Полученные особенности')}</strong>
                    {predatorGifts(draft).map((g) => (
                      <p key={g.id}>
                        {l(traitById[g.id].name)} · {g.value}
                        {g.detail ? ` · ${l(g.detail)}` : ''}
                      </p>
                    ))}
                  </div>
                </>
              )}
            </>
          )}
          {step === 5 && (
            <>
              <TraitPicker draft={draft} onChange={setDraft} />
              <div className="vtm-grid">
                {draft.traits.some((t) => t.id === 'disciplineAffinity') && (
                  <Select
                    label="Дисциплина сродства"
                    value={draft.affinity}
                    onChange={(affinity) => update({ affinity })}
                  >
                    <option value="">{l('Выберите дисциплину')}</option>
                    {Object.entries(disciplines)
                      .filter(([id]) => id !== 'alchemy')
                      .map(([id, name]) => (
                        <option key={id} value={id}>
                          {l(name)}
                        </option>
                      ))}
                  </Select>
                )}
                {draft.traits.some((t) => t.id === 'thinCurse') && (
                  <Select
                    label="Проклятие клана"
                    value={draft.curseClan}
                    onChange={(curseClan) => update({ curseClan })}
                  >
                    <option value="">{l('Выберите клан')}</option>
                    {Object.entries(clans)
                      .filter(([id]) => !['thinBlood', 'caitiff'].includes(id))
                      .map(([id, c]) => (
                        <option key={id} value={id}>
                          {l(c.name)}
                        </option>
                      ))}
                  </Select>
                )}
              </div>
              {draft.traits.some((t) => t.id === 'thinAlchemist') && (
                <p className="vtm-note">
                  {l(
                    'Получено: Алхимия Тонкой Крови 1 и одна формула 1 уровня. Впишите её название в чарнике.',
                  )}
                </p>
              )}
              <p className="vtm-hint">
                {l('Дисциплины с учётом бонусов:')}{' '}
                {Object.entries(creationDisciplines(draft))
                  .filter(([, n]) => n > 0)
                  .map(([id, n]) => `${l(disciplines[id])} ${n}`)
                  .join(' · ') || l('Нет')}
              </p>
            </>
          )}
          {step === 6 && (
            <>
              {draft.convictions.map((c, i) => (
                <section className="vtm-conviction" key={i}>
                  <div className="vtm-section-heading">
                    <h3>{l('Убеждение и опора')}</h3>
                    {draft.convictions.length > 1 && (
                      <button
                        className="button quiet"
                        onClick={() =>
                          update({ convictions: draft.convictions.filter((_, j) => j !== i) })
                        }
                      >
                        {l('Удалить пару')}
                      </button>
                    )}
                  </div>
                  <div className="vtm-conviction-cards">
                    <div className="vtm-pair-card">
                      <Field
                        label={t('vtm.conviction')}
                        multiline
                        value={c.belief}
                        onChange={(belief) =>
                          update({
                            convictions: draft.convictions.map((x, j) =>
                              j === i ? { ...x, belief } : x,
                            ),
                          })
                        }
                      />
                    </div>
                    <div className="vtm-pair-card">
                      <Field
                        label="Имя опоры"
                        value={c.name}
                        onChange={(name) =>
                          update({
                            convictions: draft.convictions.map((x, j) =>
                              j === i ? { ...x, name } : x,
                            ),
                          })
                        }
                      />
                      <Field
                        label="Описание опоры"
                        multiline
                        value={c.description}
                        onChange={(description) =>
                          update({
                            convictions: draft.convictions.map((x, j) =>
                              j === i ? { ...x, description } : x,
                            ),
                          })
                        }
                      />
                    </div>
                  </div>
                </section>
              ))}
              {draft.convictions.length < 3 && (
                <button
                  className="button quiet"
                  onClick={() => update({ convictions: [...draft.convictions, emptyConviction()] })}
                >
                  {l('+ Добавить убеждение и опору')}
                </button>
              )}
            </>
          )}
          <Errors errors={errors} />
          <div className="vtm-wizard-actions">
            <button
              className="button quiet"
              onClick={() => {
                if (step > 0) {
                  setStep(step - 1)
                  setErrors([])
                } else onBack()
              }}
            >
              {l(step ? 'Назад' : 'Отмена')}
            </button>
            <button className="button primary" onClick={next}>
              {l(step === 6 ? 'Открыть чарник' : 'Далее →')}
            </button>
          </div>
        </section>
      </div>
    </div>
  )
}
