import assert from 'node:assert/strict'
import { createServer } from 'vite'
const server = await createServer({
  server: { middlewareMode: true, hmr: false, ws: false },
  appType: 'custom',
})
try {
  const rules = await server.ssrLoadModule('/src/features/vtm/rules.ts')
  const locale = await server.ssrLoadModule('/src/lib/i18n.ts')
  locale.setLanguage('ru')
  const { formatMessage } = await server.ssrLoadModule('/src/features/vtm/localization.ts')
  const r = {
    ...rules,
    validateStep: (...args) => rules.validateStep(...args).map(formatMessage),
    traitErrors: (...args) => rules.traitErrors(...args).map(formatMessage),
    enableRules: (...args) => rules.enableRules(...args).map(formatMessage),
    buy: (...args) => {
      const result = rules.buy(...args)
      return { ...result, error: result.error ? formatMessage(result.error) : undefined }
    },
  }

  const { attributes, skills, catalog, traitById, clans } = await server.ssrLoadModule(
    '/src/features/vtm/catalog.ts',
  )
  const { predators } = await server.ssrLoadModule('/src/features/vtm/predators.ts')
  const entry = (id, value) => ({
    key: crypto.randomUUID(),
    id,
    value,
    detail: '',
    source: 'player',
    floor: 0,
  })
  const valid = () => {
    const d = r.newDraft()
    Object.assign(d, {
      name: 'Алиса',
      clan: 'toreador',
      age: 'neonate',
      disciplines: { auspex: 2, presence: 1 },
      predator: 'sandman',
      predatorDiscipline: 'obfuscate',
      extraSpecialty: { skill: 'athletics', text: 'Бег' },
      traits: [entry('resources', 4), entry('contacts', 3), entry('enemy', 2)],
      convictions: [{ belief: 'Защищай слабых', name: 'Мария', description: 'Моя сестра' }],
    })
    attributes.forEach(([id], i) => (d.attributes[id] = [4, 3, 3, 3, 2, 2, 2, 2, 1][i]))
    skills.forEach(([id], i) => (d.skills[id] = r.distributions.balanced[i]))
    d.specialties = { craft: 'Живопись', performance: 'Вокал' }
    return d
  }
  let d = valid()
  for (let i = 0; i < 7; i++) assert.deepEqual(r.validateStep(d, i), [], `Valid step ${i}`)
  const c = r.completeDraft(d, true)
  assert.equal(c.health, 6)
  assert.equal(c.willpower, 3)
  assert.equal(c.xp, 15)
  assert.equal(r.traitTotal(c.traits, 'resources'), 5)
  assert.equal(r.traitTotal(r.predatorGifts(d), 'resources'), 1)
  assert.equal(c.disciplines.obfuscate, 1)
  const before = JSON.stringify(c)
  let result = r.buy(c, { kind: 'skill', target: 'athletics' })
  assert.equal(result.character.spentXp, 12)
  assert.equal(c.skills.athletics, 3)
  assert.equal(JSON.stringify(c), before, 'Purchasing must not mutate the current sheet')
  assert.ok(r.buy(result.character, { kind: 'attribute', target: 'strength' }).error)
  result = r.buy(c, { kind: 'attribute', target: 'resolve' })
  assert.equal(result.character.willpower, 4)
  assert.equal(result.character.spentXp, 10)
  assert.equal(r.purchaseCost(c, { kind: 'discipline', target: 'auspex' }), 15)
  assert.equal(r.purchaseCost(c, { kind: 'discipline', target: 'potence' }), 7)
  assert.equal(
    r.purchaseCost({ ...c, clan: 'caitiff' }, { kind: 'discipline', target: 'auspex' }),
    18,
  )
  assert.equal(r.purchaseCost(c, { kind: 'potency', target: '' }), 20)
  assert.equal(
    r.buy({ ...c, xp: 200 }, { kind: 'ritual', target: '', level: 1, text: 'Ритуал' }).error,
    'Уровень превышает соответствующую дисциплину.',
  )
  assert.equal(
    r.buy(
      { ...c, xp: 200, disciplines: { ...c.disciplines, bloodSorcery: 2 } },
      { kind: 'ritual', target: '', level: 2, text: 'Ритуал' },
    ).character.spentXp,
    6,
  )
  assert.deepEqual(traitById.upAllNight.values, [2, 4])
  assert.deepEqual(traitById.unholyWill.values, [2, 4])
  assert.equal(traitById.linguistics.values, null)
  assert.ok(r.enableRules({ ...c, health: c.health + 1 }).length)
  assert.deepEqual(
    r.enableRules({ ...c, skills: { ...c.skills, brawl: 99 } }),
    [],
    'Re-enabling checks only the two formulas',
  )
  const context = { ...c, traits: [entry('haven', 1), entry('noHaven', 1)] }
  assert.ok(r.traitErrors(context).length)
  for (const [a, b] of [
    ['favoredBlood', 'muddledBlood'],
    ['uncleFangs', 'liquidator'],
    ['herd', 'obviousPredator'],
    ['resources', 'destitute'],
    ['anarchComrades', 'anarchShunned'],
    ['lifelike', 'deadFlesh'],
    ['vampiricResilience', 'mortalFrailty'],
  ])
    assert.ok(
      r
        .traitErrors({
          ...c,
          clan: a === 'favoredBlood' || a === 'uncleFangs' ? 'caitiff' : 'thinBlood',
          traits: [entry(a, 1), entry(b, 1)],
        })
        .some((e) => e.includes('Несовместимы')),
      `${a}/${b}`,
    )
  for (const id of [
    'zeroed',
    'cobbler',
    'influencer',
    'enduringFame',
    'cell',
    'laboratory',
    'secretMaster',
    'armored',
    'sparePlates',
  ])
    assert.ok(r.traitErrors({ ...c, traits: [entry(id, 1)] }).length, `Prerequisite: ${id}`)
  d = valid()
  d.clan = 'thinBlood'
  d.age = 'childer'
  d.generation = 15
  d.disciplines = {}
  d.predator = ''
  d.traits = [
    entry('resources', 4),
    entry('contacts', 3),
    entry('enemy', 2),
    entry('disciplineAffinity', 1),
    entry('thinAlchemist', 1),
    entry('babyTeeth', 1),
    entry('nightTerrors', 1),
  ]
  d.affinity = 'auspex'
  assert.deepEqual(r.validateStep(d, 5), [])
  let thin = r.completeDraft(d, true)
  assert.equal(thin.potency, 0)
  assert.equal(thin.disciplines.auspex, 1)
  assert.equal(thin.disciplines.alchemy, 1)
  assert.equal(thin.rituals[0].kind, 'formula')
  assert.ok(r.buy({ ...thin, xp: 100 }, { kind: 'discipline', target: 'potence' }).error)
  d.traits.push(entry('dayDrinker', 1))
  assert.ok(r.validateStep(d, 5).some((e) => e.includes('баланс 0')))
  for (const age of ['neonate', 'ancilla']) {
    d = valid()
    d.age = age
    const char = r.completeDraft(d, true)
    assert.equal(char.xp, age === 'neonate' ? 15 : 35)
    assert.equal(char.humanity, age === 'neonate' ? 7 : 6)
  }
  for (const p of predators) {
    d = valid()
    d.clan = 'hecata'
    d.age = 'ancilla'
    d.predator = p.id
    d.predatorDiscipline = p.disciplines.find((id) => id !== 'bloodSorcery') ?? p.disciplines[0]
    d.pools = {}
    for (const pool of p.pools ?? []) d.pools[pool.ids[0]] = pool.dots
    const gifts = r.predatorGifts(d)
    assert.ok(
      gifts.every((g) => traitById[g.id] && g.value > 0),
      `Predator gifts: ${p.id}`,
    )
    const char = r.completeDraft(d, true)
    assert.equal(char.humanity, 6 + (p.humanity ?? 0))
    assert.equal(char.potency, 2 + (p.potency ?? 0))
  }
  d = valid()
  d.predator = 'bloodLeech'
  d.predatorDiscipline = 'protean'
  assert.equal(r.traitTotal(r.predatorGifts(d), 'preyExclusion'), 2)
  assert.ok(
    !r.validateStep(d, 5).some((e) => e.startsWith('«Исключённая')),
    'Mandatory 2-dot exception is accepted',
  )
  d = valid()
  d.predator = 'siren'
  d.predatorDiscipline = 'presence'
  d.traits = [
    entry('stunning', 2),
    entry('contacts', 3),
    entry('resources', 2),
    entry('darkSecret', 2),
  ]
  assert.deepEqual(r.validateStep(d, 5), [])
  const siren = r.completeDraft(d, true)
  assert.equal(r.traitTotal(siren.traits, 'stunning'), 4)
  assert.equal(r.traitTotal(siren.traits, 'beautiful'), 0)
  assert.equal(siren.traits.find((t) => t.id === 'stunning').floor, 2)
  const beauty = { ...c, traits: [entry('beautiful', 2)], xp: 20 }
  assert.equal(r.purchaseCost(beauty, { kind: 'advantage', target: 'stunning' }), 6)
  const upgraded = r.buy(beauty, { kind: 'advantage', target: 'stunning' })
  assert.equal(upgraded.error, undefined)
  assert.equal(r.traitTotal(upgraded.character.traits, 'stunning'), 4)
  assert.equal(r.traitTotal(upgraded.character.traits, 'beautiful'), 0)
  assert.equal(upgraded.character.spentXp, 6)
  d = valid()
  d.traits.push(entry('daymares', 1))
  assert.deepEqual(
    r.validateStep(d, 5),
    [],
    'Unpriced Discipline flaws do not pay the 2-dot budget',
  )
  d.disciplines = { presence: 2, celerity: 1 }
  assert.ok(r.validateStep(d, 5).some((e) => e.includes('Дневные кошмары')))
  d = valid()
  d.traits = [...d.traits, entry('resources', 0)]
  d.traits[d.traits.length - 1].detail = 'Доход от реставрации'
  const detailed = r.completeDraft(d, true)
  assert.equal(detailed.traits.filter((t) => t.id === 'resources').length, 1)
  assert.equal(detailed.traits.find((t) => t.id === 'resources').detail, 'Доход от реставрации')
  const assigned = { strength: 4, dexterity: 0, stamina: 0 }
  assert.equal(r.canAssignRating(assigned, 'dexterity', 4, r.attributeQuotas), false)
  assert.equal(r.canAssignRating(assigned, 'strength', 0, r.attributeQuotas), true)
  assert.equal(r.canAssignRating(assigned, 'dexterity', 3, r.attributeQuotas), true)
  assert.equal(
    r.canAssignRating({ a: 3, b: 3, c: 3, d: 0 }, 'd', 3, r.skillQuotas('balanced')),
    false,
  )
  assert.equal(r.canAssignRating({ a: 4, b: 0 }, 'b', 4, r.skillQuotas('jack')), false)
  assert.equal(
    r.canAssignRating({ auspex: 2, presence: 1, celerity: 0 }, 'celerity', 1, r.disciplineQuotas),
    false,
  )
  assert.equal(
    r.canAssignRating({ auspex: 2, presence: 0 }, 'presence', 2, r.disciplineQuotas),
    false,
  )
  assert.equal(
    r.canAssignRating({ auspex: 2, presence: 0 }, 'presence', 1, r.disciplineQuotas),
    true,
  )
  const reshaped = r.changeDistribution({ ...valid(), skills: { a: 3, b: 3, c: 3, d: 4 } }, 'jack')
  assert.equal(Object.values(reshaped.skills).filter((n) => n === 3).length, 1)
  assert.equal(reshaped.skills.d, 0)
  for (const distribution of ['balanced', 'jack', 'specialist']) {
    const values = Object.fromEntries(skills.map(([id]) => [id, 0]))
    r.distributions[distribution].forEach((n, i) => {
      const id = skills[i][0]
      assert.equal(
        r.canAssignRating(values, id, n, r.skillQuotas(distribution)),
        true,
        `Legal assignment ${distribution}/${id}`,
      )
      values[id] = n
    })
    assert.equal(r.canAssignRating(values, skills.at(-1)[0], 4, r.skillQuotas(distribution)), false)
  }
  assert.equal(skills.find(([id]) => id === 'performance')[1], 'Исполнение')
  const localizedDisciplines = (await server.ssrLoadModule('/src/features/vtm/catalog.ts'))
    .disciplines
  assert.equal(localizedDisciplines.presence, 'Величие')
  assert.equal(localizedDisciplines.obfuscate, 'Сокрытие')
  assert.equal(localizedDisciplines.protean, 'Метаморфозы')
  assert.equal(clans.salubri.disciplines.join(','), 'auspex,dominate,fortitude')
  assert.equal(new Set(catalog.map((t) => t.id)).size, catalog.length, 'Unique catalog IDs')
  assert.ok(catalog.length > 150)
  assert.ok(!catalog.some((t) => /Cults|Ghouls|Ashfinders/.test(t.group)))
  console.log(
    `V5 checks passed: 7 creation steps, ${predators.length} predators, ${catalog.length} catalog entries, restrictions, thin-blood balance, derived stats, XP prices and atomic purchases.`,
  )
} finally {
  await server.close()
}
