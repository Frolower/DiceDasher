# Vampire: The Masquerade V5 — frontend

Entry point: authenticated `/characters`. Character creation now replaces the example cards.
Only V5 is offered; other systems have no frontend creator yet.

## State and boundaries

Characters stay in memory for the current tab, partitioned by account username. Navigating within
DiceDasher retains completed sheets; reloading the tab clears them. No character API or database
writes are made. JSON export includes the full current sheet and experience history. There is
currently no JSON import. The existing login/session handling remains in place. The shared collection is implemented in `features/characters`; V5 rules and sheets remain in `features/vtm`.

The system-neutral character collection and system-selection dialog follow the app language (English/Russian). The entire V5 wizard, both sheet pages, catalog, validation, tooltips and XP advancement follow the same English/Russian setting. Original English catalog names remain searchable. The PDF is a visual
reference, not an embedded form or a runtime dependency.

## Creation and editing

Seven stages cover identity, attributes, one of three skill distributions and specialties, clan
disciplines, 18 predator types, advantages/flaws, and 1–3 conviction/touchstone pairs. Step navigation
retains form values. Attribute, skill and discipline quotas block invalid clicks immediately. Attributes begin unassigned (0) in the wizard, and every rating can be cleared. Changing skill distribution retains feasible values and clears values exceeding the new quotas. Completing validates every stage again, including choices changed on earlier
stages. Ancilla budgets are 9/4 and XP 35; neonates start with XP 15.

Predator gifts are separate from the 7/2 (or 9/4) player budget, with locked floors and editable
additional dots. A Siren can upgrade Beautiful to Stunning by paying only the extra two dots.
Thin-blood traits use counts (1–3 merits and the same number of flaws, maximum six combined).
Discipline Affinity and Thin-blood Alchemist supply their starting dots; the latter also provides
one formula slot. Thin-bloods may choose or skip a predator type as requested.

Catalog data is a local snapshot from <https://vtm.paradoxwikis.com/Advantages_and_Flaws>, checked
2026-11-18. It retains names, categories, fixed values, discrete alternatives, and open-ended values.
Descriptions are omitted. Cult sections, ghouls, loresheets and coterie management are excluded.
Unpriced Ingrained Discipline Flaws are separate, require their discipline and never pay the normal
flaw-dot budget. Domain is included only as the Tithe Collector grant/selection.

The source does not state the Roadside Killer's Prey Exclusion (locals) cost. This implementation
uses the catalog's standard one-dot Prey Exclusion. Blood Leech's explicit two-dot variant and Grim
Reaper's one-dot Allies grant are preserved as predator exceptions. They aren't silently rewritten
to the catalog's standard values. Salubri uses Auspex, Dominate, Fortitude.

Rules can be disabled at any point to move to unrestricted sheet editing. Re-enabling checks
**only** `Health = Stamina + 3` and `Willpower = Composure + Resolve`, with exact equality. It
never attempts to reconstruct creation or retroactively validate the remainder of a free-form sheet.
Advancement charges each next level separately, checks available XP and mechanical restrictions,
and applies the purchase atomically. Players award their own XP. Ritual/formula names and discipline
powers are free text; the system doesn't infer narrative eligibility from those names.

Trait compatibility checks cover explicit exclusions, clan/age access, skill limits, discipline
requirements, Mask/Fame/Haven/mobile-haven prerequisites, and the configured thin-blood clan curse.
Narrative assertions are player supplied. The diablerie prerequisite is shown only with a selected diablerie-dependent flaw, not as a general creation question. The technical catalog source/exclusion note is not displayed in the flow.

## Design

The existing React app's semantic colors and controls surround a V5 sheet, with red section rules,
dot ratings and three attribute/skill columns drawn from the supplied two-page reference. The
whole-page view scales the natural page to available height and width. Editing uses readable
controls and reflows to one column on narrow screens. Long text is edited at the readable scale.
The sheet can grow if the character has many traits or disciplines.

Tokens reuse `--panel`, `--ink`, `--muted`, `--line`, and `--accent` across existing themes. The
sheet's red is `#991f36` in light appearances and `#f3a4b5` in dark appearances. Body/control text is
16px; display title 29–42px; labels 13px; controls at least 44px tall except the dot targets (32px
on regular desktop, 40px on narrower wizard layouts). Static paper overview is deliberately scaled;
a clearly labelled action returns to the readable editor. No decorative animation or blur is added.

Applied apple-design references (web foundations rather than native macOS navigation conventions):

- `entering-data.md › Best practices`: “Dynamically validate field values.” Validation appears
  inline; fields and budgets retain their values when corrected.
- `buttons.md › Style`: “Keep the number of prominent buttons to one or two per view.” Each
  stage has one prominent next/completion action.
- `buttons.md › Content`: “Ensure that each button clearly communicates its purpose.” Rule toggle,
  page switching, fit/edit, export, and XP purchases use text labels.

Additional reviewed references: accessibility, layout, typography, color, designing-for-macos,
text-fields, toggles and cross-platform. Keyboard focus, native dialog focus containment,
checkbox/radio semantics and number inputs are preserved. All option pickers share the system-neutral `SelectField` component: labelled combobox/listbox, portalled popup, selected/disabled states, arrow keys, Home/End, typeahead, Enter, Escape and focus restoration. The popup stays inside a modal dialog when opened there. Existing reduced-motion
rules cover the new controls. Existing light/dark theme tokens are reused rather than introducing
a separate appearance preference.

## Verification

Run `npm test`, `npm run build`, and `npm run lint` in `frontend`.
`check-vtm.mjs` exercises complete creation, all predator data, budgets, locked gift exceptions,
discrete values, unpriced flaws, compatibility/prerequisites, thin-blood balance and grants,
re-enabling rules, XP costs, health/willpower updates, upgrades and immutable failed purchases.

Browser checks cover the full guided flow, extending predator Resources 1 to 5 within the 7-dot
player budget, an XP purchase with derived Willpower, rejection of Health 99 on re-enabling rules,
two-page switching and narrative editing. Responsive/appearance checks use a temporary local
fixture with the same production CharacterSheet component; no preview route or auth bypass is shipped.

## Terminology and interface revision

Russian attribute/skill names follow the supplied `Charsheet V5 custom rus.pdf`, including
Исполнение, Воровство, Фехтование, Уличное чутьё, Гуманитарные науки, Естественные науки,
Наблюдательность, Техника, Обаяние, Манипуляция, Смекалка and Упорство. Discipline names follow
<https://wta5.ru/vampire/disciplines>, checked 2026-11-18: Сокрытие, Величие, Метаморфозы,
Алхимия Тонкой Крови and the remaining listed disciplines. The age-category label is
«Возрастная категория».

The trait stage has budget summaries, a searchable category-filtered catalog, explicit add buttons
and a separate chosen list with compact discrete steppers. Details are collapsed by default. The
chosen list moves above the catalog at narrow widths. No control forces a label to share a flex row
with a 100%-width select. Budgets and mandatory predator floors are enforced while editing.

Browser regression checks cover English collection/setup, custom pickers inside dialogs, Escape
without dismissing the dialog, arrow/Enter selection, quota-disabled ratings, predator Resources
upgrades, Haven/No Haven rejection, and the trait layout at 390px in dark appearance.

## Complete localization

Both languages are mandatory across the entire frontend. V5 copy and built-in labels live in the
`frontend/src/locales/systems/vtm/{en,ru}.json`; `labelKeys` maps only built-in Russian data labels to stable keys.
Rule errors and XP history store message keys and parameters, with nested label references, rather
than translated strings. They update when the app language changes. IDs, trait categories and
selected option values remain independent of display language.

Predator-generated specialty text and gift details carry provenance flags. Only these built-in
defaults are localized; player-authored names, specialties, background details, powers, notes and
biography remain exactly as entered, even when they happen to match a built-in term.

`check-vtm-i18n.mjs` renders all seven wizard steps and both sheet pages in overview/editor and
rule/free modes in both languages. It checks the full catalog, every predator's defaults,
prerequisites, errors and XP history, rejects Russian service text in English screens and verifies
that changing language does not mutate character data. Character UI code loads on demand.

## Character editing refinements

Guided creation always uses rules; the mode switch is available in setup and on the finished sheet.
Jack of All Trades and Balanced cap initial Skills at 3, while Specialist includes one Skill at 4.
The wizard explains this beside the distribution choices. Description is optional on the identity
step and remains editable on the sheet. Collection cards show system, name, clan, generation and
description. Conviction/Touchstone pairs use two responsive cards without numbered labels.

All sheet Attributes, Skills and selected Disciplines use five clickable dots. Free editing permits
changes directly; rule-mode increases purchase each intervening level for XP atomically. Failed
purchases leave all values and XP unchanged. Overview mode remains read-only. Age category is
retained as creation metadata without displaying it on the sheet.

Numeric editors keep a temporary empty string during editing and commit a field-specific fallback
on blur (10 for generation, the field minimum elsewhere). Trait removal resets predator-derived
traits to their mandatory floor and removes ordinary selections. Mandatory grants cannot be removed.
Browser checks verify all seven English steps, language switching with authored text preserved,
trait removal, numeric clearing and fallback, free dot editing, and the collection card fields.

## Playing from the sheet

View and edit modes are independent of zoom. Outside editing, selected dropdown values render as
plain text. Health and Willpower damage boxes remain interactive, and Hunger has five interactive
squares (0–5), in both page-fit and close-up views. Other character fields stay locked.

On the background page, Notes grows with its text while reserving space for rituals, other powers
and the footer. Its height is capped by the page boundary or the taller left column; further text
scrolls inside Notes. Adding rituals reduces the available Notes space. The collection card also
shows the authored Chronicle name.

Browser checks cover superficial/aggravated damage and Hunger changes while editing is closed,
plain Clan/Predator values, and Notes with 8, 100 and 200 lines. The latter two keep the same field
and page heights with internal scrolling, and the browser reports no console errors.
