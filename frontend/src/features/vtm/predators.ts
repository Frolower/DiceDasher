export type Gift = { id: string; value: number; detail?: string }
export type Predator = {
  id: string
  name: string
  original: string
  specialties: [string, string][]
  disciplines: string[]
  humanity?: number
  potency?: number
  fixed?: Gift[]
  choices?: Gift[][]
  pools?: { ids: string[]; dots: number; label: string }[]
}
export const predators: Predator[] = [
  {
    id: 'alleycat',
    name: 'Уличный хищник',
    original: 'Alleycat',
    specialties: [
      ['intimidation', 'Ограбления'],
      ['brawl', 'Захваты'],
    ],
    disciplines: ['celerity', 'potence'],
    humanity: -1,
    fixed: [{ id: 'contacts', value: 3, detail: 'Криминальные контакты' }],
  },
  {
    id: 'bagger',
    name: 'Пакетник',
    original: 'Bagger',
    specialties: [
      ['larceny', 'Взлом замков'],
      ['streetwise', 'Чёрный рынок'],
    ],
    disciplines: ['bloodSorcery', 'oblivion', 'obfuscate'],
    fixed: [
      { id: 'ironGullet', value: 3 },
      { id: 'enemy', value: 2, detail: 'Преследователь, считающий вас должником' },
    ],
  },
  {
    id: 'bloodLeech',
    name: 'Кровосос',
    original: 'Blood Leech',
    specialties: [
      ['brawl', 'Сородичи'],
      ['stealth', 'Против Сородичей'],
    ],
    disciplines: ['celerity', 'protean'],
    humanity: -1,
    potency: 1,
    fixed: [{ id: 'preyExclusion', value: 2, detail: 'Смертные' }],
    choices: [
      [
        { id: 'darkSecret', value: 2, detail: 'Диаблерист' },
        { id: 'shunned', value: 2 },
      ],
    ],
  },
  {
    id: 'cleaver',
    name: 'Семьянин',
    original: 'Cleaver',
    specialties: [
      ['persuasion', 'Газлайтинг'],
      ['subterfuge', 'Сокрытие следов'],
    ],
    disciplines: ['dominate', 'animalism'],
    fixed: [
      { id: 'darkSecret', value: 1, detail: 'Семьянин' },
      { id: 'herd', value: 2 },
    ],
  },
  {
    id: 'consensualist',
    name: 'По согласию',
    original: 'Consensualist',
    specialties: [
      ['medicine', 'Флеботомия'],
      ['persuasion', 'Сосуды'],
    ],
    disciplines: ['auspex', 'fortitude'],
    humanity: 1,
    fixed: [
      { id: 'darkSecret', value: 1, detail: 'Нарушитель Маскарада' },
      { id: 'preyExclusion', value: 1, detail: 'Несогласные' },
    ],
  },
  {
    id: 'farmer',
    name: 'Фермер',
    original: 'Farmer',
    specialties: [
      ['animalKen', 'Вид животного'],
      ['survival', 'Охота'],
    ],
    disciplines: ['animalism', 'protean'],
    humanity: 1,
    fixed: [{ id: 'farmer', value: 2 }],
  },
  {
    id: 'osiris',
    name: 'Осирис',
    original: 'Osiris',
    specialties: [
      ['occult', 'Традиция'],
      ['performance', 'Вид выступления'],
    ],
    disciplines: ['bloodSorcery', 'presence'],
    pools: [
      { ids: ['fame', 'herd'], dots: 3, label: 'Известность и стадо' },
      {
        ids: [
          'enemy',
          'folkloricBane',
          'folkloricBlock',
          'stigmata',
          'stakeBait',
          'starvingDecay',
          'cursedObject',
          'twiceCursed',
          'resistantBlush',
          'landLocked',
          'corpseFlesh',
        ],
        dots: 2,
        label: 'Враги и мифические недостатки',
      },
    ],
  },
  {
    id: 'sandman',
    name: 'Песочный человек',
    original: 'Sandman',
    specialties: [
      ['medicine', 'Анестезия'],
      ['stealth', 'Проникновение'],
    ],
    disciplines: ['auspex', 'obfuscate'],
    fixed: [{ id: 'resources', value: 1 }],
  },
  {
    id: 'sceneQueen',
    name: 'Королева тусовки',
    original: 'Scene Queen',
    specialties: [
      ['etiquette', 'Своя субкультура'],
      ['leadership', 'Своя субкультура'],
      ['streetwise', 'Своя субкультура'],
    ],
    disciplines: ['dominate', 'potence'],
    fixed: [
      { id: 'fame', value: 1 },
      { id: 'contacts', value: 1 },
    ],
    choices: [
      [
        { id: 'disliked', value: 1, detail: 'За пределами субкультуры' },
        { id: 'preyExclusion', value: 1, detail: 'Другая субкультура' },
      ],
    ],
  },
  {
    id: 'siren',
    name: 'Сирена',
    original: 'Siren',
    specialties: [
      ['persuasion', 'Соблазнение'],
      ['subterfuge', 'Соблазнение'],
    ],
    disciplines: ['fortitude', 'presence'],
    fixed: [
      { id: 'beautiful', value: 2 },
      { id: 'enemy', value: 1, detail: 'Отвергнутый любовник или ревнивый партнёр' },
    ],
  },
  {
    id: 'extortionist',
    name: 'Вымогатель',
    original: 'Extortionist',
    specialties: [
      ['intimidation', 'Принуждение'],
      ['larceny', 'Безопасность'],
    ],
    disciplines: ['dominate', 'potence'],
    fixed: [{ id: 'enemy', value: 2, detail: 'Полиция или сбежавшая жертва' }],
    pools: [{ ids: ['contacts', 'resources'], dots: 3, label: 'Контакты и ресурсы' }],
  },
  {
    id: 'graverobber',
    name: 'Расхититель могил',
    original: 'Graverobber',
    specialties: [
      ['occult', 'Погребальные обряды'],
      ['medicine', 'Трупы'],
    ],
    disciplines: ['fortitude', 'oblivion'],
    fixed: [
      { id: 'ironGullet', value: 3 },
      { id: 'haven', value: 1 },
      { id: 'obviousPredator', value: 2 },
    ],
  },
  // The source omits the cost for locals. Use the catalog's standard one-dot Prey Exclusion.
  {
    id: 'roadsideKiller',
    name: 'Придорожный убийца',
    original: 'Roadside Killer',
    specialties: [
      ['survival', 'Дорога'],
      ['investigation', 'Жаргон вампиров'],
    ],
    disciplines: ['fortitude', 'protean'],
    fixed: [
      { id: 'herd', value: 2, detail: 'Кочующее стадо' },
      { id: 'preyExclusion', value: 1, detail: 'Местные жители' },
    ],
  },
  {
    id: 'grimReaper',
    name: 'Жнец',
    original: 'Grim Reaper',
    specialties: [
      ['awareness', 'Смерть'],
      ['larceny', 'Подделка документов'],
    ],
    disciplines: ['auspex', 'oblivion'],
    humanity: 1,
    fixed: [{ id: 'preyExclusion', value: 1, detail: 'Здоровые смертные' }],
    choices: [
      [
        { id: 'allies', value: 1, detail: 'Медицинское сообщество' },
        { id: 'influence', value: 1, detail: 'Медицинское сообщество' },
      ],
    ],
  },
  {
    id: 'montero',
    name: 'Загонщик',
    original: 'Montero',
    specialties: [
      ['leadership', 'Охотничья стая'],
      ['stealth', 'Засада'],
    ],
    disciplines: ['dominate', 'obfuscate'],
    humanity: -1,
    fixed: [{ id: 'retainers', value: 2 }],
  },
  {
    id: 'pursuer',
    name: 'Преследователь',
    original: 'Pursuer',
    specialties: [
      ['investigation', 'Профилирование'],
      ['stealth', 'Слежка'],
    ],
    disciplines: ['animalism', 'auspex'],
    humanity: -1,
    fixed: [
      { id: 'bloodhound', value: 1 },
      { id: 'contacts', value: 1, detail: 'Местные с гибкой моралью' },
    ],
  },
  {
    id: 'trapdoor',
    name: 'Ловчий',
    original: 'Trapdoor',
    specialties: [
      ['persuasion', 'Маркетинг'],
      ['stealth', 'Засады или ловушки'],
    ],
    disciplines: ['protean', 'obfuscate'],
    fixed: [{ id: 'haven', value: 1 }],
    choices: [
      [
        { id: 'retainers', value: 1 },
        { id: 'herd', value: 1 },
        { id: 'haven', value: 1 },
      ],
      [
        { id: 'creepy', value: 1 },
        { id: 'haunted', value: 1 },
      ],
    ],
  },
  {
    id: 'titheCollector',
    name: 'Сборщик дани',
    original: 'Tithe Collector',
    specialties: [
      ['intimidation', 'Сородичи'],
      ['leadership', 'Сородичи'],
    ],
    disciplines: ['dominate', 'presence'],
    fixed: [{ id: 'adversary', value: 2 }],
    choices: [
      [
        { id: 'domain', value: 3 },
        { id: 'status', value: 3 },
      ],
    ],
  },
]
export const predatorById = Object.fromEntries(predators.map((p) => [p.id, p]))
export function predatorDisciplines(id: string, clan: string) {
  return (predatorById[id]?.disciplines ?? []).filter((d) => {
    if (d === 'bloodSorcery') return ['tremere', 'banuHaqim'].includes(clan)
    if (id === 'bagger' && d === 'oblivion') return clan === 'hecata'
    return true
  })
}
