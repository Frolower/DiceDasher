# Frontend localization

The site always has complete English and Russian versions, everywhere. Every new or changed
user-facing screen, dialog, label, option, error, status, tooltip and accessibility label must
support both languages through the existing locale dictionaries. System-specific flows and
character sheets are included; never hardcode a language on them.

Use stable IDs for data and selection values. Store generated errors/history as translation keys
and parameters so they follow language changes. Switching language must preserve in-progress forms
and character data. Never translate or overwrite player-authored text.

Add coverage for both languages when changing a flow. Check the actual rendered UI, including
later wizard stages and error states, rather than checking only dictionary parity or the first page.

Translations are split by feature under `src/locales/<feature>/{en,ru}.json`; game systems live in
`src/locales/systems/<system>/{en,ru}.json`. Register both language files together in
`src/locales/index.ts`. Keep common strings in `common`; do not duplicate shared keys in modules.
Preserve existing keys when moving strings. Tests require matching non-empty keys and parameters
in both languages, and reject duplicate keys and unregistered dictionary files.
