import { t as translate } from '../../lib/i18n'
import {
  l,
  msg,
  formatMessage,
  jsxText,
  type Message,
  joinMessages,
  sameMessage,
} from './localization'
import { labelKeys } from './labelKeys'
import { useVtmLanguage } from './useVtmLanguage'
import { useState } from 'react'
import { catalog, traitById, type Trait } from './catalog'
import {
  type Draft,
  type TraitEntry,
  availableTraits,
  creationTraits,
  draftTraitErrors,
  predatorGifts,
  sum,
  traitTotal,
  uid,
} from './rules'
import { Field, NumberField, Select } from './Controls'
export function TraitPicker({ draft, onChange }: { draft: Draft; onChange: (d: Draft) => void }) {
  useVtmLanguage()
  const [search, setSearch] = useState(''),
    [kind, setKind] = useState('merit'),
    [group, setGroup] = useState('all'),
    [error, setError] = useState<Message | null>(null)
  const all = creationTraits(draft),
    gifts = predatorGifts(draft),
    defs = availableTraits(draft.clan)
  const selected = [...new Set(all.filter((t) => t.value > 0).map((t) => t.id))]
  const normal = draft.traits.filter((t) => !traitById[t.id].thin && !traitById[t.id].discipline)
  const merit = sum(normal.filter((t) => traitById[t.id].kind !== 'flaw').map((t) => t.value)),
    flaw = sum(normal.filter((t) => traitById[t.id].kind === 'flaw').map((t) => t.value)),
    bonus = draft.age === 'ancilla' ? 2 : 0
  const thin = draft.traits.filter((t) => traitById[t.id].thin && t.value > 0),
    thinMerits = thin.filter((t) => traitById[t.id].kind !== 'flaw').length
  function change(def: Trait, total: number) {
    const floor =
      traitTotal(gifts, def.id) || (def.id === 'stunning' ? traitTotal(gifts, 'beautiful') : 0)
    if (total < floor) return
    const other = draft.traits.filter(
      (t) => t.id !== def.id && !(def.id === 'stunning' && t.id === 'beautiful'),
    )
    const existing = draft.traits.find((t) => t.id === def.id)
    const next = {
      ...draft,
      traits:
        total > floor
          ? [
              ...other,
              {
                key: existing?.key ?? uid(),
                id: def.id,
                value: total - floor,
                detail: existing?.detail ?? '',
                source: 'player' as const,
                floor: 0,
              },
            ]
          : other,
    }
    if (total > traitTotal(all, def.id)) {
      const before = draftTraitErrors(draft)
      const errors = draftTraitErrors(next).filter(
        (e) =>
          !before.some((p) => sameMessage(p, e)) &&
          e.key !== labelKeys['Выберите клановое проклятие тонкокрового.'] &&
          e.key !== labelKeys['Недостатки диаблери требуют отметки о совершённом диаблери.'],
      )
      const nextThin = next.traits.filter((t) => traitById[t.id].thin && t.value > 0)
      if (def.thin && nextThin.filter((t) => traitById[t.id].kind === def.kind).length > 3)
        errors.push(msg('vtm.chooseNoMoreThan3Thin.3fd73a'))
      const nextNormal = next.traits.filter(
        (t) => !traitById[t.id].thin && !traitById[t.id].discipline,
      )
      if (
        sum(nextNormal.filter((t) => traitById[t.id].kind !== 'flaw').map((t) => t.value)) >
        7 + bonus
      )
        errors.push(msg('vtm.allAdvantageDotsAreAssignedReduce.b09911'))
      if (
        sum(nextNormal.filter((t) => traitById[t.id].kind === 'flaw').map((t) => t.value)) >
        2 + bonus
      )
        errors.push(msg('vtm.allFlawDotsAreAssignedReduce.48116c'))
      if (errors.length) {
        setError(joinMessages(errors))
        return
      }
    }
    setError(null)
    onChange(next)
  }
  function steps(def: Trait) {
    const floor = traitTotal(gifts, def.id)
    if (def.thin || def.discipline) return [0, 1]
    if (!def.values) return null
    return [
      ...new Set([
        0,
        ...def.values,
        ...(def.id === 'beautiful' ? [4] : []),
        ...(floor ? [floor] : []),
      ]),
    ]
      .sort((a, b) => a - b)
      .filter((n) => n >= floor)
  }
  function rating(def: Trait) {
    const value = traitTotal(all, def.id),
      floor = traitTotal(gifts, def.id),
      values = steps(def)
    const previous = values ? [...values].reverse().find((n) => n < value) : value - 1,
      next = values ? values.find((n) => n > value) : value + 1
    return (
      <div
        className="trait-stepper"
        role="group"
        aria-label={translate('vtm.rating.9c6b01', {
          '0': l(def.name),
        })}
      >
        <button
          type="button"
          aria-label={translate('vtm.decrease.76090c', {
            '0': l(def.name),
          })}
          disabled={previous === undefined || previous < floor}
          onClick={() => {
            if (previous !== undefined) change(def, previous)
          }}
        >
          −
        </button>
        <output
          aria-label={translate('vtm.dots.bbff85', {
            '0': l(def.name),
          })}
        >
          {def.thin || def.discipline ? '✓' : value}
        </output>
        <button
          type="button"
          aria-label={translate('vtm.increase.6fbded', {
            '0': l(def.name),
          })}
          disabled={next === undefined}
          onClick={() => {
            if (next !== undefined)
              change(def.id === 'beautiful' && next === 4 ? traitById.stunning : def, next)
          }}
        >
          +
        </button>
      </div>
    )
  }
  const matches = defs.filter(
    (t) =>
      !selected.includes(t.id) &&
      !(t.id === 'beautiful' && selected.includes('stunning')) &&
      (kind === 'thin'
        ? t.thin
        : !t.thin && (kind === 'flaw' ? t.kind === 'flaw' : t.kind !== 'flaw')) &&
      (group === 'all' || t.group === group) &&
      `${t.name} ${t.original}`.toLowerCase().includes(search.toLowerCase()),
  )
  const groups = [
    ...new Set(
      defs
        .filter((t) =>
          kind === 'thin'
            ? t.thin
            : !t.thin && (kind === 'flaw' ? t.kind === 'flaw' : t.kind !== 'flaw'),
        )
        .map((t) => t.group),
    ),
  ]
  return (
    <div className="trait-picker">
      <div className="trait-budgets" aria-live="polite">
        <div>
          <span>{jsxText('vtm.advantages.6ca4d4')}</span>
          <strong>
            {merit}
            <small> / {7 + bonus}</small>
          </strong>
          <span>
            {Math.max(0, 7 + bonus - merit)}
            {jsxText('vtm.remaining.09c1ca')}
          </span>
        </div>
        <div>
          <span>{jsxText('vtm.flaws.685436')}</span>
          <strong>
            {flaw}
            <small> / {2 + bonus}</small>
          </strong>
          <span>
            {Math.max(0, 2 + bonus - flaw)}
            {jsxText('vtm.remaining.09c1ca')}
          </span>
        </div>
        {draft.clan === 'thinBlood' && (
          <div>
            <span>{jsxText('vtm.thinBloods.ba31f5')}</span>
            <strong>
              {thinMerits * 2 - thin.length}
              <small>{jsxText('vtm.balance.77e98f')}</small>
            </strong>
            <span>
              {thinMerits}
              {jsxText('vtm.merits.e4a565')}
              {thin.length - thinMerits}
              {jsxText('vtm.flaws.2fc895')}
              {thin.length} / 6
            </span>
          </div>
        )}
      </div>
      {error && (
        <p role="alert" className="vtm-errors">
          {formatMessage(error)}
        </p>
      )}
      <div className="trait-layout">
        <section className="trait-browser" aria-label={translate('vtm.chooseTraits.4f2d73')}>
          <div
            className="trait-type-tabs"
            role="group"
            aria-label={translate('vtm.traitType.253a40')}
          >
            {[
              ['merit', translate('vtm.advantages.6ca4d4')],
              ['flaw', translate('vtm.flaws.685436')],
              ...(draft.clan === 'thinBlood' ? [['thin', translate('vtm.thinBloods.ba31f5')]] : []),
            ].map(([id, name]) => (
              <button
                type="button"
                key={id}
                aria-pressed={kind === id}
                onClick={() => {
                  setKind(id)
                  setGroup('all')
                }}
              >
                {l(name)}
              </button>
            ))}
          </div>
          <div className="trait-filters">
            <Field
              label={translate('vtm.search.6433b5')}
              value={search}
              onChange={setSearch}
              placeholder={translate('vtm.nameInEitherLanguage.afe3b7')}
            />
            <Select label={translate('vtm.category.19c858')} value={group} onChange={setGroup}>
              <option value="all">{jsxText('vtm.allCategories.b3bf7e')}</option>
              {groups.map((g) => (
                <option key={g} value={g}>
                  {l(g)}
                </option>
              ))}
            </Select>
          </div>
          <div
            className="trait-results"
            role="region"
            aria-label={translate('vtm.traitCatalog.27d087')}
            tabIndex={0}
          >
            {matches.length ? (
              matches.map((def) => {
                const values = steps(def),
                  start = def.thin || def.discipline ? 1 : (values?.find((n) => n > 0) ?? 1)
                return (
                  <div className="trait-result" key={def.id}>
                    <div>
                      <strong>{l(def.name)}</strong>
                      <span>
                        {l(def.group)} ·{' '}
                        {def.thin || def.discipline
                          ? translate('vtm.unpriced.2141a2')
                          : def.values
                            ? def.values.join(' / ')
                            : '1+'}
                      </span>
                      <small>{def.original}</small>
                    </div>
                    <button
                      type="button"
                      className="button quiet"
                      aria-label={translate('vtm.add.1baf01', {
                        '0': l(def.name),
                      })}
                      onClick={() => change(def, start)}
                    >
                      {jsxText('vtm.add.2fc209')}
                    </button>
                  </div>
                )
              })
            ) : (
              <p className="trait-no-results">
                {jsxText('vtm.noResultsTryAnotherCategoryOr.2f5876')}
              </p>
            )}
          </div>
        </section>
        <section className="trait-chosen" aria-label={translate('vtm.selectedTraits.44c1f7')}>
          <h3>
            {jsxText('vtm.selected.dd203a')}
            <span>{selected.length}</span>
          </h3>
          <p className="vtm-hint">{jsxText('vtm.predatorGrantsAreIncludedAndDo.6edea8')}</p>
          {selected.length ? (
            selected.map((id) => {
              const def = traitById[id],
                floor = traitTotal(gifts, id)
              return (
                <article className="trait-chosen-item" key={id}>
                  <div className="trait-chosen-title">
                    <div>
                      <strong>{l(def.name)}</strong>
                      <small>
                        {floor
                          ? translate('vtm.predatorGrantMinimum.3cf5a4', {
                              '0': floor,
                            })
                          : l(def.group)}
                      </small>
                    </div>
                    {rating(def)}
                    <button
                      type="button"
                      className="trait-remove"
                      disabled={traitTotal(all, id) <= floor}
                      aria-label={translate('vtm.removeTrait', { name: l(def.name) })}
                      title={
                        floor && traitTotal(all, id) <= floor
                          ? translate('vtm.lockedGrant')
                          : translate('vtm.removeTrait', { name: l(def.name) })
                      }
                      onClick={() => change(def, floor)}
                    >
                      ×
                    </button>
                  </div>
                  <details>
                    <summary>{jsxText('vtm.details.3f72fa')}</summary>
                    <Field
                      label={translate('vtm.details.17ad5f', {
                        '0': l(def.name),
                      })}
                      value={
                        draft.traits.find((t) => t.id === id)?.detail ??
                        l(gifts.find((t) => t.id === id)?.detail ?? '')
                      }
                      onChange={(detail) => {
                        const existing = draft.traits.find((t) => t.id === id)
                        onChange({
                          ...draft,
                          traits: existing
                            ? draft.traits.map((t) =>
                                t === existing ? { ...t, detail, autoDetail: false } : t,
                              )
                            : [
                                ...draft.traits,
                                { key: uid(), id, value: 0, detail, source: 'player', floor: 0 },
                              ],
                        })
                      }}
                    />
                  </details>
                  {def.group === 'Диаблери' &&
                    !(id === 'inheritedBane' && draft.clan === 'tremere') && (
                      <label className="trait-requirement">
                        <input
                          type="checkbox"
                          checked={draft.diablerist}
                          onChange={(e) => onChange({ ...draft, diablerist: e.target.checked })}
                        />
                        <span>{jsxText('vtm.prerequisiteForThisFlawTheCharacter.7a3372')}</span>
                      </label>
                    )}
                </article>
              )
            })
          ) : (
            <p className="trait-no-results">
              {jsxText('vtm.addTraitsFromTheCatalogAdjust.992fbb')}
            </p>
          )}
        </section>
      </div>
    </div>
  )
}
export function FreeTraits({
  entries,
  onChange,
}: {
  entries: TraitEntry[]
  onChange: (t: TraitEntry[]) => void
}) {
  useVtmLanguage()
  const [id, setId] = useState(catalog[0].id),
    [custom, setCustom] = useState('')
  return (
    <div className="vtm-free-traits">
      {entries.map((t, i) => (
        <div className="vtm-free-trait" key={t.key}>
          <div>
            <strong>{t.name ?? (traitById[t.id] ? l(traitById[t.id].name) : t.id)}</strong>
            <Field
              label={translate('vtm.details.3f72fa')}
              value={t.autoDetail ? l(t.detail) : t.detail}
              onChange={(detail) =>
                onChange(entries.map((v, j) => (j === i ? { ...v, detail, autoDetail: false } : v)))
              }
            />
          </div>
          <NumberField
            label={translate('vtm.rating.9f0b99')}
            value={t.value}
            max={5}
            onChange={(value) => onChange(entries.map((v, j) => (j === i ? { ...v, value } : v)))}
          />
          <button
            type="button"
            className="button quiet"
            aria-label={translate('vtm.remove.f16a3e', {
              '0': t.name ?? (traitById[t.id] ? l(traitById[t.id].name) : t.id),
            })}
            onClick={() => onChange(entries.filter((_, j) => j !== i))}
          >
            ×
          </button>
        </div>
      ))}
      <div className="vtm-grid">
        <Select label={translate('vtm.addFromCatalog.ff67be')} value={id} onChange={setId}>
          {catalog.map((t) => (
            <option key={t.id} value={t.id}>
              {t.name} · {t.group}
            </option>
          ))}
        </Select>
        <button
          type="button"
          className="button quiet"
          onClick={() =>
            onChange([
              ...entries,
              {
                key: uid(),
                id,
                value: traitById[id].values?.[0] ?? 1,
                detail: '',
                source: 'player',
                floor: 0,
              },
            ])
          }
        >
          {jsxText('vtm.add.559a87')}
        </button>
      </div>
      <div className="vtm-grid">
        <Field label={translate('vtm.customTrait.086ebf')} value={custom} onChange={setCustom} />
        <button
          type="button"
          className="button quiet"
          disabled={!custom.trim()}
          onClick={() => {
            onChange([
              ...entries,
              { key: uid(), id: custom.trim(), value: 1, detail: '', source: 'player', floor: 0 },
            ])
            setCustom('')
          }}
        >
          {jsxText('vtm.addCustomTrait.c23b88')}
        </button>
      </div>
    </div>
  )
}
