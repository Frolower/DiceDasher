import assert from 'node:assert/strict'
import { createElement as h } from 'react'
import { renderToStaticMarkup as render } from 'react-dom/server'
import { createServer } from 'vite'
Object.defineProperty(globalThis, 'localStorage', {
  configurable: true,
  value: { getItem: () => null, setItem: () => {} },
})
const server = await createServer({
  server: { middlewareMode: true, hmr: false, ws: false },
  appType: 'custom',
})
try {
  const load = (path) => server.ssrLoadModule('/src/' + path)
  const locale = await load('lib/i18n.ts'),
    labels = await load('features/vtm/localization.ts'),
    rules = await load('features/vtm/rules.ts'),
    data = await load('features/vtm/catalog.ts')
  const { predators } = await load('features/vtm/predators.ts'),
    { CreationWizard } = await load('features/vtm/CreationWizard.tsx'),
    { CharacterSheet } = await load('features/vtm/CharacterSheet.tsx'),
    { TraitPicker, FreeTraits } = await load('features/vtm/TraitPicker.tsx'),
    { Experience } = await load('features/vtm/Experience.tsx'),
    { Errors } = await load('features/vtm/Controls.tsx')
  const d = {
    ...rules.newDraft(),
    name: 'Alice',
    description: 'A restorer seeking lost art.',
    clan: 'toreador',
    age: 'neonate',
    disciplines: { auspex: 2, presence: 1 },
    predator: 'alleycat',
    predatorDiscipline: 'celerity',
    convictions: [{ belief: 'Protect others', name: 'Mary', description: 'My sister' }],
  }
  const entry = (id, value, detail = '') => ({
    key: crypto.randomUUID(),
    id,
    value,
    detail,
    source: 'player',
    floor: 0,
  })
  const c = rules.completeDraft(d, true),
    advanced = rules.buy({ ...c, xp: 200 }, { kind: 'skill', target: 'athletics' }).character,
    invalid = rules.validateStep(rules.newDraft(), 0)
  let count = 0
  for (const language of ['en', 'ru']) {
    locale.setLanguage(language)
    const clean = (element, name) => {
      const markup = render(element)
      assert.doesNotMatch(markup, /vtm\.[\w.]+/, 'Unresolved key: ' + name)
      if (language === 'en')
        assert.doesNotMatch(markup, /[А-Яа-яЁё]/, 'Russian leaked into English: ' + name)
      count++
      return markup
    }
    for (let step = 0; step < 7; step++) {
      const screen = clean(
        h(CreationWizard, { initialDraft: d, initialStep: step }),
        'wizard ' + step,
      )
      assert.doesNotMatch(screen, /type="checkbox"/, 'Wizard must always use rules')
      if (step === 0) assert.ok(screen.includes(d.description))
      if (step === 6) {
        assert.ok(screen.includes('vtm-conviction-cards'))
        assert.doesNotMatch(screen, /(?:Conviction|Touchstone|Убеждение|Опора) 1/)
        if (language === 'en') assert.doesNotMatch(screen, /Persuasion/)
      }
    }
    const sheetEditor = render(
      h(CharacterSheet, { character: { ...c, rules: false }, initialFit: false }),
    )
    assert.doesNotMatch(sheetEditor, /Age category|Возрастная категория/)
    assert.doesNotMatch(sheetEditor, /vtm-description-panel|vtm-sheet-layout/)
    assert.ok(sheetEditor.includes(locale.t('vtm.description')))
    assert.ok(!sheetEditor.includes(d.description))

    assert.ok(sheetEditor.includes('aria-label="Strength: 5"') || language === 'ru')
    assert.equal(rules.completeDraft(d, true).description, d.description)
    clean(
      h(CreationWizard, {
        initialDraft: {
          ...d,
          clan: 'thinBlood',
          predator: '',
          disciplines: {},
          traits: [
            entry('disciplineAffinity', 1),
            entry('thinAlchemist', 1),
            entry('thinCurse', 1),
          ],
        },
        initialStep: 5,
      }),
      'thin-blood configuration',
    )
    for (const page of [1, 2])
      for (const fit of [true, false])
        for (const useRules of [true, false])
          for (const editing of [true, false])
            clean(
              h(CharacterSheet, {
                character: { ...advanced, rules: useRules },
                initialPage: page,
                initialFit: fit,
                initialEditing: editing,
              }),
              'sheet ' + page + '/' + fit + '/' + useRules + '/' + editing,
            )
    const rowEntries = ['merit', 'background', 'flaw'].map((kind) => ({
      ...entry('custom-' + kind, 2),
      kind,
      name: 'Player ' + kind,
    }))
    const rowEditor = clean(
      h(CharacterSheet, {
        character: { ...c, rules: false, traits: rowEntries },
        initialPage: 2,
        initialFit: false,
        initialEditing: true,
      }),
      'inline trait rows',
    )
    assert.equal((rowEditor.match(/class="vtm-trait-edit-row"/g) ?? []).length, 3)
    for (const row of rowEntries) assert.ok(rowEditor.includes(row.name))
    for (const row of rowEditor.split('class="vtm-trait-edit-row"').slice(1))
      assert.ok(row.split('</button>')[0].includes('max="5"'))
    assert.equal((rowEditor.match(/>\+ (?:Add row|Добавить строку)</g) ?? []).length, 3)
    const rowOverview = render(
      h(CharacterSheet, { character: { ...c, traits: rowEntries }, initialPage: 2 }),
    )
    for (const row of rowEntries) assert.ok(rowOverview.includes(row.name))
    assert.ok(rowOverview.includes('aria-label="' + labels.l('Сила крови') + ': 10"'))
    assert.ok(rowOverview.includes('vtm-five-gap'))
    assert.doesNotMatch(rowOverview, /<output/)
    const closeView = render(h(CharacterSheet, { character: c, initialFit: false }))
    assert.match(closeView, /<input readOnly=""/)
    assert.doesNotMatch(closeView, /vtm-trait-edit-row/)
    assert.ok(closeView.includes(locale.t('vtm.editSheet')))
    assert.doesNotMatch(closeView, /role="combobox"/)
    assert.equal(
      (closeView.match(/class="vtm-box-placeholder" aria-hidden="true"/g) ?? []).length,
      13 - c.health + 10 - c.willpower,
    )
    assert.equal((closeView.match(/vtm-tracker-centered/g) ?? []).length, 3)
    const liveBoxes = [...closeView.matchAll(/<button[^>]*aria-label="([^"]*)"[^>]*>/g)].filter(
      ([markup]) =>
        markup.includes(labels.l('Здоровье') + ',') || markup.includes(labels.l('Сила воли') + ','),
    )
    assert.ok(liveBoxes.length >= 2)
    for (const [button] of liveBoxes) assert.ok(!button.includes('disabled'))
    const hunger = closeView.split('class="vtm-tracker vtm-hunger"')[1].split('</div>')[0]
    assert.equal((hunger.match(/<button/g) ?? []).length, 5)
    assert.ok(!hunger.includes('disabled'))
    for (const resonance of ['choleric', 'melancholic', 'phlegmatic', 'sanguine', 'animal']) {
      const screen = clean(
        h(CharacterSheet, {
          character: { ...c, fields: { ...c.fields, resonance } },
          initialFit: false,
        }),
        'resonance ' + resonance,
      )
      assert.ok(screen.includes(locale.t('vtm.resonance.' + resonance)))
      assert.ok(screen.includes('vtm-feeding-row'))
    }
    const longTracks = render(
      h(CharacterSheet, {
        character: { ...c, health: 13, willpower: 10, healthDamage: [1, 2] },
        initialFit: false,
      }),
    )
    const healthTrack = longTracks
      .split(
        'aria-label="' + labels.copy('{0}: максимум {1}', { 0: labels.l('Здоровье'), 1: 13 }) + '"',
      )[1]
      .split('</div>')[0]
    assert.equal((healthTrack.match(/class="vtm-box-group"/g) ?? []).length, 3)
    assert.equal((healthTrack.match(/<button/g) ?? []).length, 13)
    assert.equal((healthTrack.match(/class="vtm-box-mark"/g) ?? []).length, 2)
    const humanityTrack = longTracks
      .split('aria-label="' + labels.l('Человечность') + '"')[1]
      .split('</div>')[0]
    assert.equal((humanityTrack.match(/class="vtm-box-group"/g) ?? []).length, 2)
    assert.equal((humanityTrack.match(/<button/g) ?? []).length, 10)
    assert.ok(!humanityTrack.includes('disabled'))
    assert.doesNotMatch(
      longTracks,
      /<output|aria-label="(?:Maximum Health|Maximum Willpower|Максимум здоровья|Максимум силы воли)"/,
    )
    const descriptionView = render(h(CharacterSheet, { character: c, initialPage: 3 }))
    assert.match(descriptionView, /<textarea[^>]*readOnly=""|<textarea[^>]*readonly=""/)
    clean(h(FreeTraits, { entries: c.traits, onChange: () => {} }), 'free trait editor')
    clean(
      h(TraitPicker, {
        draft: { ...d, traits: [entry('blatantDiablerist', 1)] },
        onChange: () => {},
      }),
      'diablerie prerequisite',
    )
    clean(h(CharacterSheet, { character: c, initialPage: 3 }), 'description page')
    clean(h(Experience, { character: advanced }), 'XP history')
    const errorMarkup = clean(h(Errors, { errors: invalid }), 'validation')
    assert.ok(
      errorMarkup.includes(
        language === 'en' ? 'Enter a character name.' : 'Введите имя персонажа.',
      ),
    )
    for (const trait of data.catalog) {
      const text = labels.l(trait.name) + ' ' + labels.l(trait.group)
      if (language === 'en') assert.doesNotMatch(text, /[А-Яа-яЁё]/, 'Catalog: ' + trait.id)
      for (const message of rules.traitErrors({
        ...c,
        traits: [entry(trait.id, trait.values?.[0] ?? 1)],
      })) {
        if (language === 'en')
          assert.doesNotMatch(
            labels.formatMessage(message),
            /[А-Яа-яЁё]/,
            'Prerequisite: ' + trait.id,
          )
      }
    }
    for (const predator of predators) {
      const text = [
        labels.l(predator.name),
        ...predator.specialties.map(([, text]) => labels.l(text)),
        ...(predator.fixed ?? []).map((g) => labels.l(g.detail ?? '')),
        ...(predator.choices ?? []).flat().map((g) => labels.l(g.detail ?? '')),
        ...(predator.pools ?? []).map((p) => labels.l(p.label)),
      ].join(' ')
      if (language === 'en') assert.doesNotMatch(text, /[А-Яа-яЁё]/, 'Predator: ' + predator.id)
    }
  }
  locale.setLanguage('en')
  const english = labels.formatMessage(invalid[0])
  locale.setLanguage('ru')
  assert.notEqual(labels.formatMessage(invalid[0]), english)
  const before = JSON.stringify(advanced)
  for (const language of ['en', 'ru']) {
    locale.setLanguage(language)
    render(h(CharacterSheet, { character: advanced, initialPage: 2, initialFit: false }))
  }
  assert.equal(JSON.stringify(advanced), before, 'Language must not mutate a character')
  const authored = {
    ...c,
    name: 'Красота',
    specialties: [{ skill: 'craft', text: 'Ресурсы' }],
    fields: { notes: 'Сокрытие' },
    traits: [entry('custom-trait', 1, 'Известность')],
  }
  locale.setLanguage('en')
  const one = render(h(CharacterSheet, { character: authored, initialPage: 1, initialFit: false })),
    two = render(h(CharacterSheet, { character: authored, initialPage: 2, initialFit: false }))
  assert.ok(one.includes('Красота'))
  assert.ok(one.includes('Ресурсы'))
  assert.ok(two.includes('Сокрытие'))
  assert.ok(two.includes('Известность'))
  const free = render(
    h(FreeTraits, { entries: [entry('Красота', 1, 'Сокрытие')], onChange: () => {} }),
  )
  assert.ok(free.includes('Красота'))
  assert.ok(free.includes('Сокрытие'))
  const { en: dictionary } = await server.ssrLoadModule('/src/locales/index.ts')
  for (const [key, value] of Object.entries(dictionary).filter(([key]) => key.startsWith('vtm.')))
    assert.doesNotMatch(value, /[А-Яа-яЁё]/, 'Russian in English dictionary ' + key)
  console.log(
    `V5 localization checks passed: ${count} screens, 7 stages, both pages, rule/free modes, ${data.catalog.length} traits, ${predators.length} predators, prerequisites, XP history, language-independent messages and preserved player text.`,
  )
} finally {
  await server.close()
}
