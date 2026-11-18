// Names and dot values: https://vtm.paradoxwikis.com/Advantages_and_Flaws
// Snapshot: 2026-11-18. Descriptions intentionally omitted. Cults and ghouls excluded.
export type Trait = {
  id: string
  name: string
  original: string
  group: string
  kind: 'merit' | 'background' | 'flaw'
  values: number[] | null
  thin?: boolean
  discipline?: string
}
export const attributes = [
  ['strength', 'Сила'],
  ['dexterity', 'Ловкость'],
  ['stamina', 'Выносливость'],
  ['charisma', 'Обаяние'],
  ['manipulation', 'Манипуляция'],
  ['composure', 'Самообладание'],
  ['intelligence', 'Интеллект'],
  ['wits', 'Смекалка'],
  ['resolve', 'Упорство'],
] as const
export const skills = [
  ['athletics', 'Атлетика'],
  ['brawl', 'Драка'],
  ['craft', 'Ремесло'],
  ['drive', 'Вождение'],
  ['firearms', 'Стрельба'],
  ['larceny', 'Воровство'],
  ['melee', 'Фехтование'],
  ['stealth', 'Скрытность'],
  ['survival', 'Выживание'],
  ['animalKen', 'Обращение с животными'],
  ['etiquette', 'Этикет'],
  ['insight', 'Проницательность'],
  ['intimidation', 'Запугивание'],
  ['leadership', 'Лидерство'],
  ['performance', 'Исполнение'],
  ['persuasion', 'Убеждение'],
  ['streetwise', 'Уличное чутьё'],
  ['subterfuge', 'Хитрость'],
  ['academics', 'Гуманитарные науки'],
  ['awareness', 'Наблюдательность'],
  ['finance', 'Финансы'],
  ['investigation', 'Расследование'],
  ['medicine', 'Медицина'],
  ['occult', 'Оккультизм'],
  ['politics', 'Политика'],
  ['science', 'Естественные науки'],
  ['technology', 'Техника'],
] as const
export const disciplines: Record<string, string> = {
  animalism: 'Анимализм',
  auspex: 'Ясновидение',
  bloodSorcery: 'Магия Крови',
  celerity: 'Стремительность',
  dominate: 'Доминирование',
  fortitude: 'Стойкость',
  obfuscate: 'Сокрытие',
  oblivion: 'Забвение',
  potence: 'Мощь',
  presence: 'Величие',
  protean: 'Метаморфозы',
  alchemy: 'Алхимия Тонкой Крови',
}
export const clans: Record<string, { name: string; disciplines: string[] }> = {
  banuHaqim: { name: 'Бану Хаким', disciplines: ['bloodSorcery', 'celerity', 'obfuscate'] },
  brujah: { name: 'Бруха', disciplines: ['celerity', 'potence', 'presence'] },
  gangrel: { name: 'Гангрел', disciplines: ['animalism', 'fortitude', 'protean'] },
  hecata: { name: 'Геката', disciplines: ['auspex', 'fortitude', 'oblivion'] },
  lasombra: { name: 'Ласомбра', disciplines: ['dominate', 'oblivion', 'potence'] },
  malkavian: { name: 'Малкавиан', disciplines: ['auspex', 'dominate', 'obfuscate'] },
  ministry: { name: 'Министерство', disciplines: ['obfuscate', 'presence', 'protean'] },
  nosferatu: { name: 'Носферату', disciplines: ['animalism', 'obfuscate', 'potence'] },
  ravnos: { name: 'Равнос', disciplines: ['animalism', 'obfuscate', 'presence'] },
  salubri: { name: 'Салюбри', disciplines: ['auspex', 'dominate', 'fortitude'] },
  toreador: { name: 'Тореадор', disciplines: ['auspex', 'celerity', 'presence'] },
  tremere: { name: 'Тремер', disciplines: ['auspex', 'bloodSorcery', 'dominate'] },
  tzimisce: { name: 'Цимисхи', disciplines: ['animalism', 'dominate', 'protean'] },
  ventrue: { name: 'Вентру', disciplines: ['dominate', 'fortitude', 'presence'] },
  caitiff: { name: 'Каитиф', disciplines: [] },
  thinBlood: { name: 'Тонкокровый', disciplines: [] },
}
export const catalog: Trait[] = []
function group(group: string, kind: Trait['kind'], rows: string, thin = false) {
  for (const row of rows.trim().split('\n')) {
    const [id, name, original, dots] = row.split('|')
    let values: number[] | null = null
    if (dots && dots !== '+') {
      if (dots.includes('-')) {
        const [a, b] = dots.split('-').map(Number)
        values = Array.from({ length: b - a + 1 }, (_, i) => a + i)
      } else values = dots.split(',').map(Number)
    }
    catalog.push({ id, name, original, group, kind, values, thin })
  }
}
group('Языки', 'merit', `linguistics|Языки|Linguistics|+`)
group('Языки', 'flaw', `illiterate|Неграмотность|Illiterate|2`)
group(
  'Внешность',
  'merit',
  `beautiful|Красота|Beautiful|2
stunning|Ослепительная красота|Stunning|4
semblance|Лицо Мафусаила|Semblance of the Methuselah|1-2
famousFace|Знакомое лицо|Famous Face|1
ingenue|Невинный вид|Ingénue|1
remarkableFeature|Запоминающаяся черта|Remarkable Feature|1
upAllNight|Всю ночь напролёт|Up All Night|2,4
sceneKid|Свой в тусовке|Scene Kid|1`,
)
group(
  'Внешность',
  'flaw',
  `ugly|Некрасивый|Ugly|1
repulsive|Отталкивающий|Repulsive|2
stench|Зловоние|Stench|1
transparent|Всё на лице|Transparent|1
unblinking|Неживой облик|Unblinking Visage|2`,
)
group(
  'Зависимости',
  'merit',
  `functioningAddict|Функциональная зависимость|High Functioning Addict|1`,
)
group(
  'Зависимости',
  'flaw',
  `addiction|Зависимость|Addiction|1
hopelessAddiction|Безнадёжная зависимость|Hopeless Addiction|2`,
)
group('Архаичность', 'merit', `custodian|Хранитель истории|Custodian of History|1`)
group(
  'Архаичность',
  'flaw',
  `livingPast|Жизнь в прошлом|Living in the Past|1
archaic|Архаичность|Archaic|2
griefPhobia|Фобия утраты|Grief Phobia|1
oldTricks|Старые приёмы|Old Tricks|1`,
)
group(
  'Узы крови',
  'merit',
  `bondResistance|Сопротивление узам|Bond Resistance|1-3
shortBond|Короткие узы|Short Bond|2
unbondable|Неподвластный узам|Unbondable|5
fealty|Узы верности|Bonds of Fealty|3
enduringBond|Долговечные узы|Enduring Bond|1`,
)
group(
  'Узы крови',
  'flaw',
  `bondJunkie|Зависимость от уз|Bond Junkie|1
longBond|Долгие узы|Long Bond|1
bondslave|Раб уз|Bondslave|2`,
)
group('Сверхъестественное', 'flaw', `twoMasters|Два хозяина|Two Masters|1`)
group(
  'Питание',
  'merit',
  `bloodhound|Ищейка|Bloodhound|1
ironGullet|Железная глотка|Iron Gullet|3
vesselRecognition|Опознание сосудов|Vessel Recognition|1
driveThru|На ходу|Drive-thru|1`,
)
group(
  'Питание',
  'flaw',
  `preyExclusion|Исключённая добыча|Prey Exclusion|1
methuselahThirst|Жажда Мафусаила|Methuselah's Thirst|1
farmer|Фермер|Farmer|2
organovore|Пожиратель органов|Organovore|2
veinTapper|Тайное кормление|Vein Tapper|1
outdatedPreference|Устаревшие вкусы|Outdated Preference|2
resonanceSensitivity|Чувствительность к резонансу|Resonance Sensitivity|1
resonanceMimic|Подражание резонансу|Resonance Mimic|2
sloppyFeeder|Неаккуратное кормление|Sloppy Feeder|2`,
)
group(
  'Мифические',
  'merit',
  `eatFood|Способность есть|Eat Food|2
coldHunger|Мёртвый голод|Cold Dead Hunger|3
packDiablerie|Совместное диаблери|Pack Diablerie|2
devilLuck|Дьявольская удача|Luck of the Devil|4
nuitMode|Новый облик|Nuit Mode|2
objectPower|Предмет силы|Object of Power|1-3
leyLeach|Сила древних путей|Ley Line Leach|1
persistentBlush|Стойкий румянец|Persistent Blush|3`,
)
group(
  'Мифические',
  'flaw',
  `folkloricBane|Фольклорная уязвимость|Folkloric Bane|1
folkloricBlock|Фольклорный запрет|Folkloric Block|1
stigmata|Стигматы|Stigmata|1
stakeBait|Смертельный кол|Stake Bait|2
starvingDecay|Голодное разложение|Starving Decay|2
cursedObject|Проклятый предмет|Cursed Object|1
twiceCursed|Двойное проклятие|Twice Cursed|2
resistantBlush|Слабый румянец|Resistant Blush|1
landLocked|Привязанность к земле|Land Locked|1
corpseFlesh|Мёртвая плоть|Corpse Flesh|2`,
)
group(
  'Психология',
  'merit',
  `unholyWill|Нечестивая воля|Unholy Will|2,4
zealotry|Фанатизм|Zealotry|1-3
penitence|Покаяние|Penitence|1-5
soothedBeast|Усмирённый Зверь|Soothed Beast|1
falseLove|Ложная любовь|False Love|1`,
)
group(
  'Психология',
  'flaw',
  `beaconProfanity|Маяк скверны|Beacon of Profanity|1
crisisFaith|Кризис веры|Crisis of Faith|1
penitenceScars|Шрамы покаяния|Horrible Scars of Penitence|1
grovelingWorm|Пресмыкающийся червь|Groveling Worm|2`,
)
group(
  'Заражение',
  'flaw',
  `diseaseVector|Переносчик болезни|Disease Vector|1
plaguebringer|Разносчик чумы|Plaguebringer|1-2`,
)
group(
  'Родство',
  'merit',
  `bloodSense|Чувство родства|Consanguineous Sense|2
bloodInfluence|Влияние родства|Consanguineous Influence|2
sinsFather|Грехи отцов|Sins of the Father|2,3`,
)
group(
  'Диаблери',
  'flaw',
  `blatantDiablerist|Явный диаблерист|Blatant Diablerist|1
inheritedBane|Унаследованное проклятие|Inherited Bane|2`,
)
group(
  'Прочее',
  'merit',
  `checkTrunk|Запас в багажнике|Check the Trunk|1
sideHustler|Подработка|Side Hustler|2
temperedWill|Закалённая воля|Tempered Will|3
untouchable|Неприкасаемый|Untouchable|5
mysticVoid|Мистик пустоты|Mystic of the Void|1,2`,
)
group(
  'Прочее',
  'flaw',
  `knowledgeHungry|Жажда знаний|Knowledge Hungry|1
prestationDebts|Долги услуг|Prestation Debts|1
riskTaker|Любитель риска|Risk-Taker|1
weakWilled|Слабоволие|Weak-Willed|2`,
)
group(
  'Каитифы',
  'merit',
  `favoredBlood|Благословенная кровь|Favored Blood|4
markCaine|Метка Каина|Mark of Caine|2
mockingbird|Пересмешник|Mockingbird|3
sunScarred|Опалённый солнцем|Sun-Scarred|5
uncleFangs|Дядюшка Клык|Uncle Fangs|3`,
)
group(
  'Каитифы',
  'flaw',
  `befoulingVitae|Оскверняющая витэ|Befouling Vitae|2
caitiffCurse|Клановое проклятие|Clan Curse|2
debtPeon|Долговая кабала|Debt Peon|2
liquidator|Ликвидатор|Liquidator|1
muddledBlood|Мутная кровь|Muddled Blood|1
walkingOmen|Ходячее предзнаменование|Walking Omen|2
wordScarred|Письмена на теле|Word-Scarred|1`,
)
group(
  'Тонкокровые',
  'merit',
  `anarchComrades|Товарищи-анархи|Anarch Comrades
camarillaContact|Контакт в Камарилье|Camarilla Contact
catenatingBlood|Связующая кровь|Catenating Blood
dayDrinker|Дневной гуляка|Day Drinker
disciplineAffinity|Сродство с дисциплиной|Discipline Affinity
lifelike|Как живой|Lifelike
thinAlchemist|Тонкокровый алхимик|Thin-blood Alchemist
vampiricResilience|Вампирская стойкость|Vampiric Resilience
abhorrentBlood|Отвратительная кровь|Abhorrent Blood
faithProof|Иммунитет к вере|Faith-Proof
lowAppetite|Слабый аппетит|Low Appetite
lucidDreamer|Осознанные сновидения|Lucid Dreamer
mortalityMien|Облик смертного|Mortality's Mien
swiftFeeder|Быстрое кормление|Swift Feeder`,
  true,
)
group(
  'Тонкокровые',
  'flaw',
  `anarchShunned|Отвергнутый анархами|Shunned by the Anarchs
camarillaBranded|Клеймо Камарильи|Branded by the Camarilla
bestialTemper|Звериный нрав|Bestial Temper
thinCurse|Клановое проклятие|Clan Curse
vitaeDependency|Зависимость от витэ|Vitae Dependency
deadFlesh|Мёртвое тело|Dead Flesh
babyTeeth|Молочные зубы|Baby Teeth
mortalFrailty|Хрупкость смертного|Mortal Frailty
heliophobia|Гелиофобия|Heliophobia
nightTerrors|Ночные кошмары|Night Terrors
plagueBearer|Носитель заразы|Plague Bearers
sloppyDrinker|Неаккуратное питьё|Sloppy Drinker
sunFaded|Угасание на солнце|Sun-Faded
supernaturalTell|Сверхъестественная примета|Supernatural Tell
twilightPresence|Сумеречное присутствие|Twilight Presence
unendingHunger|Нескончаемый голод|Unending Hunger`,
  true,
)
group('Союзники', 'background', `allies|Союзники|Allies|2-6`)
group('Союзники', 'flaw', `enemy|Враг|Enemy|+`)
group('Контакты', 'background', `contacts|Контакты|Contacts|1-3`)
group('Известность', 'background', `fame|Известность|Fame|1-5`)
group(
  'Известность',
  'merit',
  `influencer|Лидер мнений|Influencer|1
enduringFame|Неувядающая слава|Enduring Fame|1`,
)
group(
  'Известность',
  'flaw',
  `darkSecret|Тёмная тайна|Dark Secret|+
infamy|Дурная слава|Infamy|+
bannedFrom|Изгнание из города|Banned From|1-3`,
)
group('Влияние', 'background', `influence|Влияние|Influence|1-5`)
group(
  'Влияние',
  'flaw',
  `disliked|Неприязнь|Disliked|1
despised|Презрение|Despised|2`,
)
group('Убежище', 'background', `haven|Убежище|Haven|1-3`)
group(
  'Убежище',
  'merit',
  `hiddenArmory|Скрытый арсенал|Hidden Armory|+
cell|Камера|Cell|+
watchmen|Охрана|Watchmen|+
laboratory|Лаборатория|Laboratory|+
library|Библиотека|Library|+
location|Расположение|Location|1
luxury|Роскошь|Luxury|1
postern|Тайный выход|Postern|+
securitySystem|Охранная система|Security System|+
surgery|Операционная|Surgery|1
warding|Обереги|Warding|+
holyGround|Священная земля|Holy Ground|1
shrine|Святилище|Shrine|1-3
business|Предприятие|Business Establishment|2-3
furcus|Фуркус|Furcus|1-3
machineShop|Мастерская|Machine Shop|+
mobile|Передвижное убежище|Mobile|1-3
armored|Бронирование|Armored|1
smugglersStash|Тайник контрабандиста|Smugglers Stash|1-2
sparePlates|Запасные номера|Spare Plates|2`,
)
group(
  'Убежище',
  'flaw',
  `noHaven|Нет убежища|No Haven|1
creepy|Жуткое место|Creepy|1
haunted|Призраки|Haunted|+
compromised|Раскрытое убежище|Compromised|2
shared|Общее убежище|Shared|1,2
onRails|По расписанию|On the Rails|1
temperamental|Ненадёжная техника|Temperamental|1`,
)
group('Стадо', 'background', `herd|Стадо|Herd|1-5`)
group('Стадо', 'flaw', `obviousPredator|Очевидный хищник|Obvious Predator|2`)
group('Маска', 'background', `mask|Маска|Mask|1-2`)
group(
  'Маска',
  'merit',
  `zeroed|Стёртое прошлое|Zeroed|1
cobbler|Изготовитель масок|Cobbler|1`,
)
group(
  'Маска',
  'flaw',
  `knownCorpse|Известный мертвец|Known Corpse|1
knownBlankbody|Известный вампир|Known Blankbody|2`,
)
group('Мавла', 'background', `mawla|Мавла|Mawla|1-5`)
group(
  'Мавла',
  'flaw',
  `adversary|Противник|Adversary|+
secretMaster|Тайный хозяин|Secret Master|1
shamefulChilde|Брошенное дитя|Shameful Childe|1
embracedTouchstone|Обращённая врагами опора|Touchstone Embraced by your Enemies|2`,
)
group('Ресурсы', 'background', `resources|Ресурсы|Resources|1-5`)
group('Ресурсы', 'flaw', `destitute|Нищета|Destitute|1`)
group('Слуги', 'background', `retainers|Слуги|Retainers|1-3`)
group('Слуги', 'flaw', `stalkers|Преследователи|Stalkers|1`)
group('Статус', 'background', `status|Статус|Status|1-5`)
group('Статус', 'merit', `citySecrets|Тайны города|City Secrets|1-3`)
group(
  'Статус',
  'flaw',
  `suspect|Под подозрением|Suspect|1
shunned|Отверженный|Shunned|2
mortalPretender|Подражатель смертным|Mortal Pretender|1`,
)
// Needed by Tithe Collector; domain/coterie management itself is outside this feature.
group('Домен', 'background', `domain|Домен|Domain|+`)
const ingrained: [string, string, string, string][] = [
  ['untamed', 'Неукротимость', 'Untamed', 'animalism'],
  ['daymares', 'Дневные кошмары', 'Daymares', 'auspex'],
  ['sanguinaryAnimism', 'Кровавый анимизм', 'Sanguinary Animism', 'bloodSorcery'],
  ['breakdown', 'Надрыв', 'Breakdown', 'celerity'],
  ['blunt', 'Прямолинейность', 'Blunt', 'dominate'],
  ['scarTissue', 'Рубцовая ткань', 'Scar Tissue', 'fortitude'],
  ['faded', 'Угасание', 'Faded', 'obfuscate'],
  ['monstrous', 'Чудовищность', 'Monstrous', 'oblivion'],
  ['killerInstinct', 'Инстинкт убийцы', 'Killer Instinct', 'potence'],
  ['egomaniac', 'Эгомания', 'Egomaniac', 'presence'],
  ['stasis', 'Стазис', 'Stasis', 'protean'],
]
ingrained.forEach(([id, name, original, discipline]) =>
  catalog.push({
    id,
    name,
    original,
    discipline,
    group: 'Недостатки дисциплин',
    kind: 'flaw',
    values: null,
  }),
)
export const traitById = Object.fromEntries(catalog.map((t) => [t.id, t]))
export const skillName = (id: string) => skills.find((s) => s[0] === id)?.[1] ?? id
