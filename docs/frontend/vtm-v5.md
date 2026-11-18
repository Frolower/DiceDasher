# Vampire: The Masquerade V5 — frontend

The authenticated `/characters` page provides a system-neutral collection and system-selection
dialog. V5 is currently the only implemented creator. The collection lives in
`frontend/src/features/characters`; the V5 catalog, rules, wizard and sheet live in
`frontend/src/features/vtm`.

## State and persistence

Completed characters stay in tab memory, partitioned by account username. In-app navigation retains
them; reloading the tab clears them. There are no character API calls or database writes. JSON
export includes the current character and experience history; JSON import is not implemented.

Cards show the system, character name, clan, generation, Chronicle and description. The collection
and creation dialog are system-neutral; V5 terminology appears inside its creator and sheet.

## Guided creation

The setup dialog offers creation with rules or a blank sheet. Guided creation always uses rules
and has seven stages:

1. Identity, clan, age category, generation and optional description.
2. Attributes: one at 4, three at 3, four at 2 and one at 1.
3. Skills: Jack of All Trades, Balanced or Specialist, plus specialties.
4. Starting disciplines: two dots in one discipline and one in another.
5. Predator type and its grants, discipline and specialty choices.
6. Merits, backgrounds and flaws, including prerequisites and exclusions.
7. One to three Conviction/Touchstone pairs, displayed as two cards per pair without numbering.

Attributes begin unassigned at 0. Attribute, skill and discipline quotas prevent invalid clicks;
ratings can be cleared. Changing the skill distribution retains feasible values and clears excess
assignments. Jack of All Trades and Balanced cap starting skills at 3; Specialist has one skill at 4.
There is no extra Specialist/max-rating hint beside the distribution choices.

Returning to an earlier step retains values. Completion validates all seven stages again, including
choices affected by earlier changes. Age category is creation metadata and is not shown on the sheet.

Player trait budgets are 7 advantage dots and 2 flaw dots, or 9/4 for Ancillae. Predator gifts have
separate mandatory floors and do not spend that budget; additional player dots do. The guided trait
picker supports search in either language, category filters, discrete value steppers, collapsed
details and quick removal. Removing an upgraded predator trait resets it to its mandatory floor;
a mandatory grant itself cannot be removed in the wizard. Siren's Beautiful can upgrade to Stunning
for the two additional dots.

Thin-bloods choose one to three merits and the same number of flaws, at most six combined.
Discipline Affinity and Thin-blood Alchemist supply starting dots; the latter supplies one formula
slot. Thin-bloods may choose a predator type or skip it.

## Finished sheet: viewing, editing and play

The sheet opens in view mode. Zoom is separate from editing: whole-page view fits the natural sheet
to available width and height, while close-up view uses readable controls and responsive layout.
**Edit** enables character changes and opens the close-up view; **Done** returns to view mode.
Selected dropdown values display as plain text outside editing. The rule switch is available only
while editing the completed sheet, not inside the guided wizard.

The page switcher has **1 · Character**, **2 · Background** and **Description**. Description has no
page number and occupies its own view; it is not a sidebar or dialog. Its text is editable only in
edit mode and also appears on the collection card.

Attributes, skills and selected disciplines use five dots; Blood Potency uses ten. Increasing them
with rules enabled purchases each intervening level for XP atomically. Free mode changes them
directly. Ratings do not show a duplicate numeric value beside their dots. Merits, backgrounds and
flaws also use individually spaced dots in view mode, matching the skill ratings.

The following play controls remain interactive outside editing, in either zoom mode:

- Health and Willpower damage: click a box to cycle empty → superficial → aggravated → empty.
  Superficial damage is one corner-to-corner diagonal; aggravated damage is a cross.
- Humanity: ten squares representing 0–10.
- Hunger: five squares representing 0–5. Selected Humanity and Hunger squares use crosses.

Squares measure 20 × 20 px at close-up scale and 16 × 16 px on the natural whole-page sheet,
before page scaling. Groups have a gap after every five marks; Blood Potency dots have the same
five-mark grouping. Humanity and Willpower are centered against a ten-slot track, Health against
thirteen slots. Health/Willpower slots beyond the current maximum reserve invisible space without
creating interactive or accessible controls. Hunger stays aligned to the right edge of its row.
Visible values are represented by the marks; accessible labels retain the values and damage states.
Health and Willpower maximum number inputs appear only when editing a free-mode sheet.

In free editing, Merits, Backgrounds and Flaws each support adding and deleting rows with a name and
rating, capped at five dots. The additional catalog editor remains available below the background
page. Rules mode keeps mechanical advancement in the XP workflow. Numeric editors allow clearing
the field while typing; an empty field commits its fallback on blur (10 for generation, otherwise
the field minimum).

Resonance is a compact preset selector on the left: Choleric, Melancholic, Phlegmatic, Sanguine and
Animal, plus an unset state. It follows the same view/edit behavior as other dropdowns. Existing custom resonance text is preserved until replaced.

On the background page, Notes grows as text is entered, pushing rituals downward while reserving
space for rituals, other powers and the footer. Once that space is exhausted, Notes scrolls
internally instead of extending the page. The cap also accounts for a taller left column; adding
rituals reduces the available Notes space.

## Rules and advancement boundaries

Disabling rules on the completed sheet enables free editing. Re-enabling checks only exact equality
for `Health = Stamina + 3` and `Willpower = Composure + Resolve`. It does not reconstruct creation or
retroactively validate the rest of a freely edited sheet.

Neonates start with 15 XP and Humanity 7; Ancillae start with 35 XP and Humanity 6, before predator
modifiers. Players award their own XP. Purchases check the next-level cost, available XP, generation
limits and configured restrictions, then apply atomically. A failed purchase leaves the character
and XP unchanged. Ritual/formula names and discipline powers are player text; narrative eligibility
is not inferred from their names.

Trait compatibility checks include explicit opposing traits, clan/age access, skill and discipline
requirements, Mask/Fame/Haven/mobile-haven prerequisites and the selected thin-blood clan curse.
Narrative assertions are player supplied. The diablerie prerequisite is asked only when relevant to
a selected flaw, not as a general identity question.

## Localization and shared controls

Every screen supports English and Russian through the single application language setting. V5
strings live in `frontend/src/locales/systems/vtm/{en,ru}.json`; the shared character collection has
its own `characters` dictionaries. See the [dictionary guide](../../frontend/src/locales/README.md)
for registration and validation. `labelKeys` maps built-in catalog labels to stable translation keys.

Rule errors and XP history store keys and parameters, including nested label references, and
translate when rendered. IDs, category values and selections are language-independent. Switching
language preserves character data. Predator-generated default details/specialties carry provenance
flags; only those defaults are translated. Player names, descriptions, specialties, notes, powers
and custom trait text are never translated, even when they match built-in terms.

All dropdowns use the shared `SelectField`: labelled combobox/listbox, selected/disabled states,
arrow keys, Home/End, typeahead, Enter, Escape and focus restoration. Popups remain inside a modal
when opened there. The sheet reuses existing semantic light/dark theme tokens and focus styles.

## Verification

Run `npm test`, `npm run build`, `npm run lint` and `npm run format:check` from `frontend/`.
The [frontend testing guide](testing.md) separates automated coverage from manual browser checks.
There is currently no committed automated browser scenario that clicks through the entire wizard.
