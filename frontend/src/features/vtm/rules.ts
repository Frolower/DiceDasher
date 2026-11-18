import { msg, labelMessage, type Message, joinMessages, sameMessage } from './localization'
import { attributes, catalog, clans, disciplines, skills, traitById } from './catalog'
import { type Gift, predatorById, predatorDisciplines } from './predators'
export type TraitEntry = {
  name?: string
  kind?: 'merit' | 'background' | 'flaw'
  key: string
  id: string
  value: number
  detail: string
  autoDetail?: boolean
  source: 'player' | 'predator'
  floor: number
}
export type Specialty = {
  skill: string
  text: string
  auto?: boolean
}
export type Conviction = {
  belief: string
  name: string
  description: string
}
export type Draft = {
  name: string
  description: string
  concept: string
  chronicle: string
  sire: string
  clan: string
  age: 'childer' | 'neonate' | 'ancilla'
  generation: number
  attributes: Record<string, number>
  skills: Record<string, number>
  distribution: 'balanced' | 'jack' | 'specialist'
  specialties: Record<string, string>
  extraSpecialty: Specialty
  disciplines: Record<string, number>
  predator: string
  predatorDiscipline: string
  predatorSpecialty: number
  predatorDetail: string
  choices: number[]
  pools: Record<string, number>
  traits: TraitEntry[]
  convictions: Conviction[]
  affinity: string
  curseClan: string
  diablerist: boolean
}
export type Character = {
  id: string
  rules: boolean
  name: string
  description: string
  concept: string
  chronicle: string
  sire: string
  clan: string
  age: Draft['age']
  generation: number
  attributes: Record<string, number>
  skills: Record<string, number>
  specialties: Specialty[]
  disciplines: Record<string, number>
  predator: string
  traits: TraitEntry[]
  convictions: Conviction[]
  humanity: number
  potency: number
  health: number
  willpower: number
  hunger: number
  xp: number
  spentXp: number
  fields: Record<string, string>
  powers: Record<string, string>
  rituals: {
    name: string
    level: number
    kind: 'ritual' | 'formula'
  }[]
  healthDamage: number[]
  willpowerDamage: number[]
  log: Message[]
  affinity: string
}
export const distributions = {
  jack: [3, ...Array(8).fill(2), ...Array(10).fill(1), ...Array(8).fill(0)],
  balanced: [3, 3, 3, ...Array(5).fill(2), ...Array(7).fill(1), ...Array(12).fill(0)],
  specialist: [4, 3, 3, 3, 2, 2, 2, 1, 1, 1, ...Array(17).fill(0)],
} as const
export const emptyConviction = (): Conviction => ({ belief: '', name: '', description: '' })
export function newDraft(): Draft {
  return {
    name: '',
    description: '',
    concept: '',
    chronicle: '',
    sire: '',
    clan: '',
    age: 'childer',
    generation: 13,
    attributes: Object.fromEntries(attributes.map(([id]) => [id, 1])),
    skills: Object.fromEntries(skills.map(([id]) => [id, 0])),
    distribution: 'balanced',
    specialties: {},
    extraSpecialty: { skill: '', text: '' },
    disciplines: {},
    predator: '',
    predatorDiscipline: '',
    predatorSpecialty: 0,
    predatorDetail: '',
    choices: [],
    pools: {},
    traits: [],
    convictions: [emptyConviction()],
    affinity: '',
    curseClan: '',
    diablerist: false,
  }
}
export const uid = () => crypto.randomUUID()
export const sum = (values: number[]) => values.reduce((a, b) => a + b, 0)
export const generations = (age: Draft['age'], clan: string) =>
  clan === 'thinBlood'
    ? [14, 15, 16]
    : age === 'ancilla'
      ? [10, 11]
      : age === 'neonate'
        ? [12, 13]
        : [12, 13]
export const basePotency = (d: Draft) =>
  d.clan === 'thinBlood' ? 0 : d.age === 'ancilla' ? 2 : d.generation >= 14 ? 0 : 1
export function predatorGifts(d: Draft): TraitEntry[] {
  const p = predatorById[d.predator]
  if (!p) return []
  const gifts: Gift[] = [
    ...(p.fixed ?? []),
    ...(p.choices ?? []).flatMap((c, i) => (c[d.choices[i] ?? 0] ? [c[d.choices[i] ?? 0]] : [])),
    ...(p.pools ?? []).flatMap((pool) =>
      pool.ids.filter((id) => (d.pools[id] ?? 0) > 0).map((id) => ({ id, value: d.pools[id] })),
    ),
  ]
  const merged: Record<string, TraitEntry> = {}
  gifts.forEach((g) => {
    if (merged[g.id]) {
      merged[g.id].value += g.value
      merged[g.id].floor += g.value
    } else
      merged[g.id] = {
        key: `predator-${g.id}`,
        id:
          g.id === 'beautiful' && d.traits.some((t) => t.id === 'stunning' && t.value > 0)
            ? 'stunning'
            : g.id,
        value: g.value,
        floor: g.value,
        detail: g.detail ?? '',
        autoDetail: true,
        source: 'predator',
      }
  })
  return Object.values(merged)
}
export function creationTraits(d: Draft) {
  const gifts = predatorGifts(d)
  // Additional dots in a predator gift are represented separately until completion.
  return [...gifts, ...d.traits]
}
export function creationDisciplines(d: Draft) {
  const result = { ...d.disciplines }
  if (d.predator && d.predatorDiscipline)
    result[d.predatorDiscipline] = (result[d.predatorDiscipline] ?? 0) + 1
  if (d.clan === 'thinBlood') {
    if (d.traits.some((t) => t.id === 'disciplineAffinity') && d.affinity)
      result[d.affinity] = (result[d.affinity] ?? 0) + 1
    if (d.traits.some((t) => t.id === 'thinAlchemist')) result.alchemy = (result.alchemy ?? 0) + 1
  }
  return result
}
export function traitTotal(entries: TraitEntry[], id: string) {
  return sum(entries.filter((e) => e.id === id).map((e) => e.value))
}
const conflicts = [
  ['beautiful', 'stunning', 'ugly', 'repulsive'],
  ['addiction', 'hopelessAddiction'],
  ['shortBond', 'longBond'],
  ['unbondable', 'bondslave', 'bondJunkie', 'twoMasters'],
  ['herd', 'obviousPredator'],
  ['resources', 'destitute'],
  ['favoredBlood', 'muddledBlood'],
  ['uncleFangs', 'liquidator'],
  ['anarchComrades', 'anarchShunned'],
  ['lifelike', 'deadFlesh'],
  ['vampiricResilience', 'mortalFrailty'],
  ['penitence', 'grovelingWorm'],
  ['persistentBlush', 'resistantBlush', 'corpseFlesh'],
  ['upAllNight', 'corpseFlesh'],
] // Groups are pairwise exclusive except unbondable and the blush group; handled below.
export function traitErrors(ctx: {
  clan: string
  age: Draft['age']
  potency: number
  skills: Record<string, number>
  disciplines: Record<string, number>
  traits: TraitEntry[]
  curseClan?: string
  diablerist?: boolean
}) {
  const { traits, clan } = ctx
  const has = (id: string) => traitTotal(traits, id) > 0
  const errors: Message[] = []
  const issue = (text: Message) => {
    if (!errors.some((e) => sameMessage(e, text))) errors.push(text)
  }
  for (const group of conflicts)
    for (let i = 0; i < group.length; i++)
      for (let j = i + 1; j < group.length; j++) {
        const a = group[i],
          b = group[j]
        if (group[0] === 'unbondable' && a !== 'unbondable') continue
        if (group[0] === 'persistentBlush' && a === 'resistantBlush' && b === 'corpseFlesh')
          continue
        if (has(a) && has(b))
          issue(
            msg('vtm.incompatibleAnd.2088b8', {
              '0': labelMessage(traitById[a].name),
              '1': labelMessage(traitById[b].name),
            }),
          )
      }
  if (
    has('noHaven') &&
    traits.some(
      (t) => t.value > 0 && traitById[t.id]?.group === 'Убежище' && traitById[t.id].kind !== 'flaw',
    )
  )
    issue(msg('vtm.noHavenIsIncompatibleWithHaven.7dc342'))
  for (const t of traits.filter((t) => t.value > 0)) {
    const def = traitById[t.id]
    if (!def) continue
    if (def.discipline && !(ctx.disciplines[def.discipline] > 0))
      issue(
        msg('vtm.requiresTheDiscipline.a7c8df', {
          '0': labelMessage(def.name),
          '1': labelMessage(disciplines[def.discipline]),
        }),
      )
    if (def.thin && clan !== 'thinBlood') issue(msg('vtm.thinBloodMeritsAndFlawsAre.0de32b'))
    if (def.group === 'Каитифы' && clan !== 'caitiff')
      issue(msg('vtm.caitiffTraitsAreAvailableOnlyTo.6899ab'))
    if (def.group === 'Архаичность' && ctx.age !== 'ancilla')
      issue(msg('vtm.archaicTraitsAreAvailableOnlyTo.619a49'))
    if (def.group === 'Родство' && ['caitiff', 'thinBlood'].includes(clan))
      issue(msg('vtm.bloodTiesMeritsRequireAClan.c5fe13'))
    if (t.id === 'farmer' && (clan === 'ventrue' || ctx.potency >= 3))
      issue(msg('vtm.farmerIsUnavailableToVentrueAnd.bcb033'))
    if (t.id === 'nuitMode' && ctx.potency > 1)
      issue(msg('vtm.nuitModeRequiresBloodPotency1.2a5afa'))
    if (t.id === 'illiterate' && ((ctx.skills.academics ?? 0) > 1 || (ctx.skills.science ?? 0) > 1))
      issue(msg('vtm.illiterateLimitsAcademicsAndScienceTo.f748d4'))
    if (t.id === 'transparent' && (ctx.skills.subterfuge ?? 0) > 0)
      issue(msg('vtm.transparentRequiresSubterfuge0.d96217'))
    if (t.id === 'archaic' && (ctx.skills.technology ?? 0) > 0)
      issue(msg('vtm.archaicRequiresTechnology0.4cbc0d'))
    if (
      t.id === 'temperedWill' &&
      ((ctx.disciplines.dominate ?? 0) > 0 || (ctx.disciplines.presence ?? 0) > 0)
    )
      issue(msg('vtm.temperedWillRequiresNoDotsIn.0ff83f'))
    if (t.id === 'fealty' && !(ctx.disciplines.dominate > 0))
      issue(msg('vtm.bondsOfFealtyRequiresDominate.5c122d'))
    if (t.id === 'mysticVoid' && t.value === 2 && !['hecata', 'lasombra'].includes(clan))
      issue(msg('vtm.mysticOfTheVoidAt2.58bc35'))
    const requirements: Record<string, [string, number]> = {
      influencer: ['fame', 2],
      enduringFame: ['fame', 3],
      cell: ['haven', 2],
      laboratory: ['haven', 2],
      zeroed: ['mask', 2],
      cobbler: ['mask', 2],
      secretMaster: ['mawla', 1],
      armored: ['mobile', 1],
      smugglersStash: ['mobile', 1],
      sparePlates: ['mobile', 1],
      onRails: ['mobile', 1],
      temperamental: ['mobile', 1],
    }
    const req = requirements[t.id]
    if (req && traitTotal(traits, req[0]) < req[1])
      issue(
        msg('vtm.requires.6ccd25', {
          '0': labelMessage(def.name),
          '1': labelMessage(traitById[req[0]].name),
          '2': labelMessage(req[1]),
        }),
      )
    if (
      t.id === 'library' &&
      traitTotal(traits, 'haven') === 1 &&
      traitTotal(traits, 'library') > 1
    )
      issue(msg('vtm.a1DotHavenCanHave.8c1c3f'))
    if (
      def.group === 'Диаблери' &&
      !ctx.diablerist &&
      !(t.id === 'inheritedBane' && clan === 'tremere')
    )
      issue(msg('vtm.diablerieFlawsRequireConfirmationThatThe.e740a2'))
  }
  if (has('thinCurse')) {
    if (!ctx.curseClan) issue(msg('vtm.chooseTheThinBloodSClan.dc8a37'))
    if (['banuHaqim', 'gangrel', 'brujah'].includes(ctx.curseClan ?? '') && !has('bestialTemper'))
      issue(msg('vtm.theSelectedCurseRequiresBestialTemper.718ef3'))
    if (ctx.curseClan === 'tremere' && !has('catenatingBlood'))
      issue(msg('vtm.theTremereCurseRequiresCatenatingBlood.d6561a'))
  }
  return errors
}
export function draftTraitErrors(d: Draft) {
  return traitErrors({
    ...d,
    potency: basePotency(d) + (predatorById[d.predator]?.potency ?? 0),
    traits: creationTraits(d),
    disciplines: creationDisciplines(d),
  })
}
export function validateStep(d: Draft, step: number): Message[] {
  const errors: Message[] = []
  if (step === 0) {
    if (!d.name.trim()) errors.push(msg('vtm.enterACharacterName.37fec0'))
    if (!clans[d.clan]) errors.push(msg('vtm.chooseAClan.ee7f43'))
    if (!generations(d.age, d.clan).includes(d.generation))
      errors.push(msg('vtm.chooseAGenerationThatMatchesThe.633726'))
    if (d.clan === 'thinBlood' && d.age !== 'childer')
      errors.push(msg('vtm.chooseChildeAsTheAgeCategory.2c2e59'))
  }
  if (
    step === 1 &&
    Object.values(d.attributes).sort().join() !== [1, 2, 2, 2, 2, 3, 3, 3, 4].join()
  )
    errors.push(msg('vtm.assignAttributesOneAt4Three.243617'))
  if (step === 2) {
    if (
      Object.values(d.skills)
        .sort((a, b) => b - a)
        .join() !== [...distributions[d.distribution]].sort((a, b) => b - a).join()
    )
      errors.push(msg('vtm.yourSkillsDoNotMatchThe.7e876e'))
    for (const id of ['academics', 'craft', 'performance', 'science'])
      if (d.skills[id] > 0 && !d.specialties[id]?.trim())
        errors.push(
          msg('vtm.enterTheFreeSpecialtyFor.cbce33', {
            '0': labelMessage(skills.find((s) => s[0] === id)?.[1]),
          }),
        )
    if (
      !d.extraSpecialty.skill ||
      !d.extraSpecialty.text.trim() ||
      !d.skills[d.extraSpecialty.skill]
    )
      errors.push(msg('vtm.addOneFreeSpecialtyToA.b74b36'))
  }
  if (step === 3) {
    const vals = Object.values(d.disciplines)
      .filter((n) => n > 0)
      .sort()
    if (d.clan === 'thinBlood') {
      if (vals.length) errors.push(msg('vtm.thinBloodsHaveNoStartingDisciplines.f8d803'))
    } else if (vals.join() !== '1,2') errors.push(msg('vtm.assign2DotsToOneDiscipline.c1e860'))
    const allowed =
      d.clan === 'caitiff'
        ? Object.keys(disciplines).filter((id) => id !== 'alchemy')
        : (clans[d.clan]?.disciplines ?? [])
    if (Object.entries(d.disciplines).some(([id, n]) => n > 0 && !allowed.includes(id)))
      errors.push(msg('vtm.chooseDisciplinesAvailableToYourClan.0e57b3'))
  }
  if (step === 4) {
    const p = predatorById[d.predator]
    if (!p && !(d.age === 'childer' || d.clan === 'thinBlood'))
      errors.push(msg('vtm.chooseAPredatorType.60a56f'))
    if (p) {
      if (!predatorDisciplines(p.id, d.clan).includes(d.predatorDiscipline))
        errors.push(msg('vtm.chooseYourPredatorTypeSBonus.d80dfe'))
      if (!p.specialties[d.predatorSpecialty])
        errors.push(msg('vtm.chooseYourPredatorSpecialty.20c4d4'))
      if (p.pools?.some((pool) => sum(pool.ids.map((id) => d.pools[id] ?? 0)) !== pool.dots))
        errors.push(msg('vtm.assignAllOfYourPredatorType.1a0f3d'))
      if (d.clan === 'ventrue' && ['bagger', 'farmer'].includes(p.id))
        errors.push(msg('vtm.ventrueCannotChooseThisPredatorType.354469'))
      if (p.id === 'titheCollector' && d.age !== 'ancilla')
        errors.push(msg('vtm.titheCollectorIsAvailableToAncillae.852431'))
      for (const pool of p.pools ?? [])
        for (const id of pool.ids) {
          const n = d.pools[id] ?? 0
          if (n > 0 && traitById[id].values && !traitById[id].values!.includes(n))
            errors.push(
              msg('vtm.invalidDotRating.880606', {
                '0': labelMessage(traitById[id].name),
              }),
            )
        }
    }
  }
  if (step === 5) {
    const normal = d.traits.filter((t) => !traitById[t.id]?.thin && !traitById[t.id]?.discipline)
    const advantage = sum(
      normal.filter((t) => traitById[t.id]?.kind !== 'flaw').map((t) => t.value),
    )
    const flaws = sum(normal.filter((t) => traitById[t.id]?.kind === 'flaw').map((t) => t.value))
    const bonus = d.age === 'ancilla' ? 2 : 0
    if (advantage !== 7 + bonus || flaws !== 2 + bonus)
      errors.push(
        msg('vtm.assignExactlyAdvantageDotsAndFlaw.434a78', {
          '0': labelMessage(7 + bonus),
          '1': labelMessage(2 + bonus),
        }),
      )
    const thin = d.traits.filter((t) => traitById[t.id]?.thin)
    const merits = thin.filter((t) => traitById[t.id].kind !== 'flaw').length
    if (d.clan === 'thinBlood' && (merits < 1 || merits > 3 || thin.length !== merits * 2))
      errors.push(msg('vtm.choose13ThinBloodMerits.50d29b'))
    if (d.traits.some((t) => t.id === 'disciplineAffinity') && !d.affinity)
      errors.push(msg('vtm.chooseADisciplineForYourAffinity.b1602b'))
    // Validate combined ratings, including mandatory predator dots and their exceptional ratings.
    for (const id of new Set(creationTraits(d).map((t) => t.id))) {
      const def = traitById[id]
      if (!def || def.thin) continue
      const total = traitTotal(creationTraits(d), id),
        floor = traitTotal(predatorGifts(d), id)
      if (def.values && !def.values.includes(total) && total !== floor)
        errors.push(
          msg('vtm.allowedRatingsAre.8781b1', {
            '0': labelMessage(def.name),
            '1': labelMessage(def.values.join(', ')),
          }),
        )
    }
    errors.push(...draftTraitErrors(d))
  }
  if (step === 6) {
    if (d.convictions.length < 1 || d.convictions.length > 3)
      errors.push(msg('vtm.choose13ConvictionAndTouchstone.db586d'))
    if (d.convictions.some((c) => !c.belief.trim() || !c.name.trim() || !c.description.trim()))
      errors.push(msg('vtm.enterEachConvictionAndItsTouchstone.5d16d4'))
  }
  return errors.filter((e, i) => errors.findIndex((p) => sameMessage(p, e)) === i)
}
export function completeDraft(d: Draft, rules: boolean): Character {
  const p = predatorById[d.predator]
  const traits = creationTraits(d).reduce<TraitEntry[]>((out, t) => {
    const existing = out.find((x) => x.id === t.id)
    if (existing) {
      existing.value += t.value
      existing.floor += t.floor
      if (t.detail) {
        existing.detail = t.detail
        existing.autoDetail = t.autoDetail ?? false
      }
    } else out.push({ ...t })
    return out
  }, [])
  const specialties: Specialty[] = [
    ...Object.entries(d.specialties)
      .filter(([id, text]) => d.skills[id] > 0 && text.trim())
      .map(([skill, text]) => ({ skill, text })),
    ...(d.extraSpecialty.text.trim() ? [{ ...d.extraSpecialty }] : []),
  ]
  if (p) {
    const [skill, text] = p.specialties[d.predatorSpecialty] ?? p.specialties[0]
    specialties.push({
      skill,
      text: d.predatorDetail.trim() || text,
      auto: !d.predatorDetail.trim(),
    })
  }
  return {
    id: uid(),
    rules,
    name: d.name,
    description: d.description,
    concept: d.concept,
    chronicle: d.chronicle,
    sire: d.sire,
    clan: d.clan,
    age: d.age,
    generation: d.generation,
    attributes: { ...d.attributes },
    skills: { ...d.skills },
    specialties,
    disciplines: creationDisciplines(d),
    predator: d.predator,
    traits,
    convictions: d.convictions.map((c) => ({ ...c })),
    humanity: 7 - (d.age === 'ancilla' ? 1 : 0) + (p?.humanity ?? 0),
    potency: basePotency(d) + (p?.potency ?? 0),
    health: d.attributes.stamina + 3,
    willpower: d.attributes.composure + d.attributes.resolve,
    hunger: 1,
    xp: d.age === 'neonate' ? 15 : d.age === 'ancilla' ? 35 : 0,
    spentXp: 0,
    fields: { curseClan: d.curseClan, diablerist: String(d.diablerist) },
    powers: {},
    rituals: d.traits.some((t) => t.id === 'thinAlchemist')
      ? [{ name: '', level: 1, kind: 'formula' }]
      : [],
    healthDamage: [],
    willpowerDamage: [],
    log: [],
    affinity: d.affinity,
  }
}
export function enableRules(c: Character): Message[] {
  const errors: Message[] = []
  if (c.health !== c.attributes.stamina + 3)
    errors.push(
      msg('vtm.maximumHealthMustBeStamina3.b9d785', {
        '0': labelMessage(c.attributes.stamina + 3),
      }),
    )
  if (c.willpower !== c.attributes.composure + c.attributes.resolve)
    errors.push(
      msg('vtm.maximumWillpowerMustBeComposureResolve.15a84e', {
        '0': labelMessage(c.attributes.composure + c.attributes.resolve),
      }),
    )
  return errors
}
export type Purchase = {
  kind:
    | 'attribute'
    | 'skill'
    | 'discipline'
    | 'specialty'
    | 'advantage'
    | 'potency'
    | 'ritual'
    | 'formula'
  target: string
  text?: string
  level?: number
}
export function purchaseCost(c: Character, p: Purchase) {
  if (p.kind === 'attribute') return ((c.attributes[p.target] ?? 0) + 1) * 5
  if (p.kind === 'skill') return ((c.skills[p.target] ?? 0) + 1) * 3
  if (p.kind === 'specialty') return 3
  if (p.kind === 'potency') return (c.potency + 1) * 10
  if (p.kind === 'ritual' || p.kind === 'formula') return (p.level ?? 1) * 3
  if (p.kind === 'discipline')
    return (
      ((c.disciplines[p.target] ?? 0) + 1) *
      (c.clan === 'caitiff' ? 6 : clans[c.clan]?.disciplines.includes(p.target) ? 5 : 7)
    )
  if (p.kind === 'advantage') {
    const def = traitById[p.target]
    const current =
      traitTotal(c.traits, p.target) ||
      (p.target === 'stunning' ? traitTotal(c.traits, 'beautiful') : 0)
    const next = def?.values?.find((n) => n > current) ?? current + 1
    return (next - current) * 3
  }
  return 0
}
export function buy(
  c: Character,
  p: Purchase,
): {
  character: Character
  error?: Message
} {
  const fail = (error: Message) => ({ character: c, error })
  if (!c.rules) return fail(msg('vtm.enableRulesToSpendXP.088bb6'))
  const cost = purchaseCost(c, p)
  if (!Number.isFinite(cost) || cost < 1 || c.xp - c.spentXp < cost)
    return fail(msg('vtm.notEnoughUnspentXP.cfc8b0'))
  const n: Character = {
    ...c,
    attributes: { ...c.attributes },
    skills: { ...c.skills },
    disciplines: { ...c.disciplines },
    traits: c.traits.map((t) => ({ ...t })),
    specialties: c.specialties.map((s) => ({ ...s })),
    rituals: c.rituals.map((r) => ({ ...r })),
  }
  let label: Message | string = ''
  if (p.kind === 'attribute' || p.kind === 'skill') {
    const list = p.kind === 'attribute' ? attributes : skills
    const item = list.find(([id]) => id === p.target)
    if (!item) return fail(msg('vtm.chooseAnAttributeOrSkill.d0b6ef'))
    const values = p.kind === 'attribute' ? n.attributes : n.skills
    if (values[p.target] >= 5) return fail(msg('vtm.theMaximumRatingIs5.dedb5e'))
    values[p.target]++
    label = msg('vtm.ratingChanged', {
      '0': labelMessage(item[1]),
      '1': labelMessage(values[p.target]),
    })
    n.health = n.attributes.stamina + 3
    n.willpower = n.attributes.composure + n.attributes.resolve
  }
  if (p.kind === 'discipline') {
    if (!disciplines[p.target]) return fail(msg('vtm.chooseADiscipline.c59def'))
    if ((n.disciplines[p.target] ?? 0) >= 5)
      return fail(msg('vtm.theMaximumDisciplineRatingIs5.8f5eec'))
    if (n.clan === 'thinBlood' && p.target !== 'alchemy' && p.target !== n.affinity)
      return fail(msg('vtm.thinBloodsCanAdvanceOnlyTheir.d44afd'))
    if (n.clan !== 'thinBlood' && p.target === 'alchemy')
      return fail(msg('vtm.alchemyIsAvailableOnlyToThin.1a65cc'))
    n.disciplines[p.target] = (n.disciplines[p.target] ?? 0) + 1
    label = msg('vtm.ratingChanged', {
      '0': labelMessage(disciplines[p.target]),
      '1': labelMessage(n.disciplines[p.target]),
    })
  }
  if (p.kind === 'specialty') {
    if (!n.skills[p.target] || !p.text?.trim())
      return fail(msg('vtm.chooseALearnedSkillAndEnter.e16db8'))
    if (
      n.specialties.some(
        (s) => s.skill === p.target && s.text.toLowerCase() === p.text!.trim().toLowerCase(),
      )
    )
      return fail(msg('vtm.thatSpecialtyAlreadyExists.075f95'))
    n.specialties.push({ skill: p.target, text: p.text.trim() })
    label = msg('vtm.specialty.1fa0b7', {
      '0': p.text.trim(),
    })
  }
  if (p.kind === 'advantage') {
    const def = traitById[p.target]
    if (!def || def.kind === 'flaw' || def.thin || def.discipline)
      return fail(msg('vtm.chooseAnAdvantageWithADot.16e99b'))
    const previousBeauty =
      p.target === 'stunning' ? n.traits.find((t) => t.id === 'beautiful') : undefined
    const current = traitTotal(n.traits, p.target) || previousBeauty?.value || 0
    if (previousBeauty) n.traits = n.traits.filter((t) => t.id !== 'beautiful')
    const next = def.values?.find((v) => v > current) ?? (def.values ? undefined : current + 1)
    if (next === undefined) return fail(msg('vtm.thisAdvantageIsAlreadyAtIts.d23652'))
    const entry = n.traits.find((t) => t.id === p.target)
    if (entry) entry.value += next - current
    else
      n.traits.push({
        key: uid(),
        id: p.target,
        value: next,
        detail: p.text ?? '',
        source: previousBeauty?.source ?? 'player',
        floor: previousBeauty?.floor ?? 0,
      })
    label = msg('vtm.ratingChanged', {
      '0': labelMessage(def.name),
      '1': labelMessage(next),
    })
  }
  if (p.kind === 'potency') {
    if (n.clan === 'thinBlood') return fail(msg('vtm.thinBloodBloodPotencyIs0.7b97a8'))
    const max: Record<number, number> = { 10: 4, 11: 4, 12: 3, 13: 3, 14: 0, 15: 0, 16: 0 }
    if (n.potency >= (max[n.generation] ?? 10))
      return fail(msg('vtm.bloodPotencyHasReachedTheGeneration.9c0233'))
    n.potency++
    label = msg('vtm.bloodPotency.684d69', {
      '0': labelMessage(n.potency),
    })
  }
  if (p.kind === 'ritual' || p.kind === 'formula') {
    const level = p.level ?? 1
    if (!p.text?.trim() || !Number.isInteger(level) || level < 1 || level > 5)
      return fail(msg('vtm.enterANameAndALevel.0baad1'))
    const discipline = p.kind === 'ritual' ? 'bloodSorcery' : 'alchemy'
    if ((n.disciplines[discipline] ?? 0) < level)
      return fail(msg('vtm.theLevelExceedsTheRelevantDiscipline.b0def1'))
    if (
      n.rituals.some(
        (r) => r.kind === p.kind && r.name.toLowerCase() === p.text!.trim().toLowerCase(),
      )
    )
      return fail(msg('vtm.thatEntryAlreadyExists.2e3689'))
    n.rituals.push({ name: p.text.trim(), level, kind: p.kind })
    label = msg('vtm.namedEntry', {
      '0': labelMessage(p.kind === 'ritual' ? 'Ритуал' : 'Формула'),
      '1': p.text.trim(),
    })
  }
  // Existing free-form values aren't retroactively validated. Only prevent newly introduced rule violations.
  const context = (x: Character) =>
    traitErrors({ ...x, curseClan: x.fields.curseClan, diablerist: x.fields.diablerist === 'true' })
  const prior = context(c)
  const introduced = context(n).filter((e) => !prior.some((p) => sameMessage(p, e)))
  if (introduced.length) return fail(joinMessages(introduced))
  n.spentXp += cost
  n.log = [
    msg('vtm.xP.014a39', {
      '0': labelMessage(label),
      '1': labelMessage(cost),
    }),
    ...c.log,
  ]
  return { character: n }
}
export const availableTraits = (clan: string) =>
  catalog.filter(
    (t) => (!t.thin || clan === 'thinBlood') && (t.group !== 'Каитифы' || clan === 'caitiff'),
  )
// Creation quotas are enforced before state changes, not just when leaving a stage.
export function canAssignRating(
  values: Record<string, number>,
  id: string,
  next: number,
  quotas: Record<number, number>,
) {
  if (!Number.isInteger(next) || next < 0) return false
  if (next === 0 || next === values[id]) return true
  const used = Object.entries(values).filter(([key, n]) => key !== id && n === next).length
  return used < (quotas[next] ?? 0)
}
export const attributeQuotas: Record<number, number> = { 1: 1, 2: 4, 3: 3, 4: 1 }
export const disciplineQuotas: Record<number, number> = { 1: 1, 2: 1 }
export function skillQuotas(distribution: Draft['distribution']): Record<number, number> {
  return Object.fromEntries(
    [1, 2, 3, 4].map((n) => [n, distributions[distribution].filter((v) => v === n).length]),
  )
}
export function changeDistribution(d: Draft, distribution: Draft['distribution']): Draft {
  const quotas = skillQuotas(distribution),
    used: Record<number, number> = {}
  const skills = Object.fromEntries(
    Object.entries(d.skills).map(([id, n]) => {
      used[n] = (used[n] ?? 0) + 1
      return [id, n > 0 && used[n] > (quotas[n] ?? 0) ? 0 : n]
    }),
  )
  return { ...d, distribution, skills }
}
