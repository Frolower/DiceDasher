# Feature dictionaries

Each module owns an English/Russian pair. The app registers modules in `index.ts` and merges them
by language, so the existing language picker and `t(key, params)` work across every feature.
Current modules:

| Folder         | Contents                                                             |
| -------------- | -------------------------------------------------------------------- |
| `common/`      | Shared controls, navigation, language selection and network messages |
| `auth/`        | Sign-in, sign-up and authentication validation                       |
| `profile/`     | Account settings, profile and room collection                        |
| `themes/`      | Theme settings, palette names and descriptions                       |
| `demo/`        | Sample characters and friends                                        |
| `characters/`  | System-neutral character collection and creation setup               |
| `systems/vtm/` | V5 wizard, sheet, catalog, rule messages and XP                      |

Existing translation keys stay unchanged. A file's folder is its owner, not a required key prefix;
legacy `text.*` keys are retained. Use a descriptive feature prefix for new keys. Shared keys have
one owner in `common/` and can be referenced by any feature.

## Add a system or feature

1. Create `<feature>/en.json` and `<feature>/ru.json`; use `systems/<id>/` for game systems.
2. Add both imports and one entry to `localeModules` in `index.ts`, using the folder as the name:

   ```ts
   import dndEn from './systems/dnd/en.json'
   import dndRu from './systems/dnd/ru.json'

   // Inside localeModules:
   'systems/dnd': { en: dndEn, ru: dndRu },
   ```

3. Add matching, non-empty strings and matching `{parameter}` names to both files.
4. Render with `t('dnd.someKey')`; keep player-authored text out of translations.
5. Run `npm test`, `npm run build` and `npm run lint` from `frontend/`.

No language-picker, screen or route changes are required to register translations. Current modules
load eagerly through explicit imports. This change organizes files; it does not introduce deferred
loading or asynchronous translation states.

`merge.ts` throws if modules define the same key, naming both owners. The test suite also rejects
repeated keys inside JSON, missing language keys, blank/non-string values, different placeholders
and dictionary files omitted from the registry. Existing language persistence, tab synchronization,
English fallback and message-key-based state continue to use `src/lib/i18n.ts`.
