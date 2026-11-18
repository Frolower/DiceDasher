import { useState } from 'react'
import { attributes, catalog, disciplines, skills } from './catalog'
import { buy, purchaseCost, type Character, type Purchase } from './rules'
import { Field, NumberField, Select } from './Controls'
import { l, message, formatMessage, type Message } from './localization'
import { useVtmLanguage } from './useVtmLanguage'
export function Experience({
  character: c,
  onChange,
}: {
  character: Character
  onChange: (c: Character) => void
}) {
  useVtmLanguage()
  const [kind, setKind] = useState<Purchase['kind']>('attribute'),
    [target, setTarget] = useState('strength'),
    [text, setText] = useState(''),
    [level, setLevel] = useState(1),
    [award, setAward] = useState(1),
    [feedback, setFeedback] = useState<Message | null>(null)
  const options =
    kind === 'attribute'
      ? attributes
      : kind === 'skill' || kind === 'specialty'
        ? skills
        : kind === 'discipline'
          ? Object.entries(disciplines)
          : kind === 'advantage'
            ? catalog
                .filter((t) => t.kind !== 'flaw' && !t.thin)
                .map((t) => [t.id, t.name] as const)
            : []
  const purchase: Purchase = { kind, target, text, level }
  return (
    <section className="vtm-experience">
      <div className="vtm-section-heading">
        <div>
          <h2>{l('Развитие персонажа')}</h2>
          <p>
            {l('Доступно')} <strong>{c.xp - c.spentXp}</strong> {l('опыта · Всего')} {c.xp}{' '}
            {l('· Потрачено')} {c.spentXp}
          </p>
        </div>
      </div>
      <div className="vtm-grid">
        <NumberField label="Начислить опыт" value={award} min={1} onChange={setAward} />
        <button
          className="button quiet"
          onClick={() => {
            onChange({
              ...c,
              xp: c.xp + award,
              log: [message('Начислено {0} опыта', { '0': award }), ...c.log],
            })
            setFeedback(message('Начислено {0} опыта.', { '0': award }))
          }}
        >
          {l('Начислить')}
        </button>
      </div>
      <div className="vtm-grid">
        <Select
          label="Что улучшить"
          value={kind}
          onChange={(value) => {
            const k = value as Purchase['kind']
            setKind(k)
            setTarget(
              k === 'attribute'
                ? 'strength'
                : k === 'skill' || k === 'specialty'
                  ? 'athletics'
                  : k === 'discipline'
                    ? 'animalism'
                    : k === 'advantage'
                      ? 'linguistics'
                      : '',
            )
            setFeedback(null)
          }}
        >
          <option value="attribute">{l('Характеристика · новый уровень × 5')}</option>
          <option value="skill">{l('Навык · новый уровень × 3')}</option>
          <option value="specialty">{l('Специализация · 3')}</option>
          <option value="discipline">{l('Дисциплина · новый уровень × 5 / 6 / 7')}</option>
          <option value="advantage">{l('Преимущество · 3 за точку')}</option>
          <option value="potency">{l('Сила крови · новый уровень × 10')}</option>
          <option value="ritual">{l('Ритуал магии крови · уровень × 3')}</option>
          <option value="formula">{l('Формула тонкокровых · уровень × 3')}</option>
        </Select>
        {options.length > 0 && (
          <Select label="Выбор улучшения" value={target} onChange={setTarget}>
            {options.map(([id, name]) => (
              <option key={id} value={id}>
                {l(name)}
              </option>
            ))}
          </Select>
        )}
        {['specialty', 'ritual', 'formula'].includes(kind) && (
          <Field label="Название" value={text} onChange={setText} />
        )}{' '}
        {['ritual', 'formula'].includes(kind) && (
          <NumberField label="Уровень" min={1} max={5} value={level} onChange={setLevel} />
        )}
      </div>
      <div className="vtm-section-heading">
        <p>
          {l('Стоимость:')}{' '}
          <strong>
            {purchaseCost(c, purchase)} {l('опыта')}
          </strong>
        </p>
        <button
          className="button primary"
          onClick={() => {
            const result = buy(c, purchase)
            if (result.error) setFeedback(result.error)
            else {
              onChange(result.character)
              setFeedback(message('Приобретено за {0} опыта.', { '0': purchaseCost(c, purchase) }))
              setText('')
            }
          }}
        >
          {l('Потратить опыт')}
        </button>
      </div>
      <p role="status" className="vtm-hint">
        {formatMessage(feedback)}
      </p>
      {c.log.length > 0 && (
        <details>
          <summary>{l('История опыта')}</summary>
          <ul>
            {c.log.map((item, i) => (
              <li key={i}>{formatMessage(item)}</li>
            ))}
          </ul>
        </details>
      )}
    </section>
  )
}
