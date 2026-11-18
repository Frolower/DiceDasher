# Frontend verification

Run these commands from `frontend/`:

```sh
npm test
npm run build
npm run lint
npm run format:check
```

`npm test` runs four Node scripts using assertions; it is not a browser E2E suite. Tests use mocked
network/storage where appropriate and Vite's module loader for production TypeScript and components.

## Automated coverage

| Script                       | Coverage                                                                                                                                                                                                                                                                                                                       |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `scripts/check-auth.mjs`     | Auth validation, request contracts, safe errors, session restoration, shared refreshes, retries, logout ordering and account settings                                                                                                                                                                                          |
| `scripts/check-i18n.mjs`     | Dictionary modules, duplicate keys and unregistered files, matching non-empty RU/EN strings and parameters, language detection/persistence, English fallback, blocked storage, tab synchronization and rendered application screens                                                                                            |
| `scripts/check-vtm.mjs`      | Seven creation stages, distributions/quotas, predator grants, budgets, trait prerequisites/exclusions, discrete/open values, thin-blood balance/grants, derived stats, re-enabling rules and XP costs/atomic purchases                                                                                                         |
| `scripts/check-vtm-i18n.mjs` | Currently 70 counted RU/EN render combinations, plus additional assertions: wizard stages, sheet pages, view/edit and rule/free modes, resonance choices, catalog/predator labels, error/history localization, inline trait rows, live tracker availability, five-mark groups, hidden reserved slots and preserved player text |

The V5 UI script renders React components to static HTML. It checks displayed text, control
attributes and markup structure. It does not execute browser click handlers, lay out CSS, measure
pixels, or scroll textareas. The numeric count describes rendered combinations, not separate
browser journeys or a coverage percentage.

`npm run build` checks TypeScript and produces the Vite production bundle. `npm run lint` runs
Oxlint; `npm run format:check` checks Prettier formatting. These are separate from behavioral tests.

## Manual browser verification

During development, browser checks have exercised:

- Guided creation through all seven stages, validation errors and language switching without
  losing entered values.
- Setup and custom pickers, keyboard selection/Escape, and modal focus behavior.
- Quota-blocked ratings, predator trait floors/upgrades, incompatible selections and quick removal.
- Completed-sheet page switching, view/edit and zoom controls, free editing, XP purchases and
  exact derived-stat checks when re-enabling rules.
- Damage and Hunger clicks outside editing; square dimensions and groups of 5 + 5 + 3 for Health.
- Notes with short and long text: initial growth, reserved ritual space and internal scrolling
  after the height cap, without further page growth.
- Numeric clearing/fallback, collection card fields, and selected responsive/dark appearances.

Some checks used temporary local fixtures rendering the same production components. Those fixtures
are removed after use; no preview/auth-bypass route is shipped. These checks document observed
behavior at the time, rather than providing an automatically repeatable regression suite.

## Remaining test boundaries

There is no committed Playwright/Cypress-style E2E suite. A green `npm test` does not prove that a
player can click every control through a complete browser session. CSS alignment, scrolling,
responsive layouts, pointer interactions and focus still need browser verification when changed.
Color contrast is not measured by the automated suite.

Character backend persistence is not tested because it is not implemented; characters currently
live in tab memory. JSON import is also not implemented. See [V5 frontend behavior](vtm-v5.md) for
state and rules boundaries.
