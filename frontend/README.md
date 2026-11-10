# DiceDasher frontend

React + TypeScript + Vite. Landing page with a rotating d20, account management,
English/Russian UI, and selectable color themes.

## Docker

From the repository root:

```sh
make up SERVICE=frontend
```

Open http://localhost:5174. Compose also starts PostgreSQL and the backend.
The image builds the app with Node and serves only the resulting static files with Nginx.
Nginx forwards `/api/*` to the backend on the internal Docker network, stripping `/api`.
It preserves the browser Origin and cookies. `/health` checks the frontend itself.

After frontend changes, rebuild with the same command. To follow its logs:

```sh
make logs SERVICE=frontend
```

## Development with live reload

```sh
cd frontend
npm install
npm run dev
```

Open http://localhost:5174. The port is fixed so the backend Origin check stays predictable.
Stop the frontend container first (`docker compose -f docker-compose.yml -f docker-compose.dev.yml stop frontend` from the repository root)
when using Vite locally, because both use port 5174.

The Vite development proxy forwards `/api/*` to `http://localhost:8000/*`.
From the repository root, start or recreate the backend with the updated Compose configuration:

```sh
make up SERVICE=backend
```

Compose reads `ALLOWED_ORIGINS` from `services/backend/.env` without overriding it.
The default backend allowlist contains Swagger at `http://localhost:8888` and the frontend at `http://localhost:5174`.
For a backend running outside Compose, set `ALLOWED_ORIGINS=http://localhost:8888,http://localhost:5174` and `COOKIE_SECURE=false` for local HTTP.
The Docker image includes the `/api` reverse proxy; Vite's development proxy is not included in `dist`.
For deployment on a domain, set `FRONTEND_ORIGIN` to the public HTTPS Origin,
terminate HTTPS at your ingress, and set `COOKIE_SECURE=true` in the backend environment.

## Account behavior

- Registration sends email, username, and password to `POST /register`; repeat password is validated locally and never sent.
- Login sends username and password to `POST /login` with credentials included for the HttpOnly refresh cookie.
- Validation follows the current Go registration rules, including 8–72 ASCII characters, English uppercase/lowercase letters, a digit, and a printable special character. Spaces, Cyrillic, accented letters, and emoji are rejected. Login applies the same ASCII character restrictions, but does not apply the registration password complexity checks.
- Access tokens are kept in memory only. On page load, the app restores the session using the HttpOnly refresh cookie (`POST /api/refresh`, then `GET /api/me`). It also refreshes tokens before expiry and retries authenticated requests once after a `401`.
- Signed-in users can log out using the red **Log out** button at the bottom of the left navigation menu. It sends `POST /api/logout` with credentials included. On success, the app clears its in-memory session, closes the menu, and returns to `/`.
- During logout, the button is disabled and shows **Logging out…**. Failed requests leave the session in place and display an error beside the button so the user can retry. HTTP errors use messages defined in `src/lib/auth.ts`; network errors currently display the browser's error message.
- Logout waits for any active refresh to finish before revoking the session. Concurrent logout calls share one request, and new refresh calls wait for logout to settle.
- The central d20 rotates slowly; click it or press Enter/Space to pause or resume. Rotation pauses while an account dialog is open and respects the system reduced-motion preference.

### Authentication code

- `src/lib/auth.ts`: session store, login and refresh requests, and `logout()` (server request and session clearing).
- `src/App.tsx`: `handleLogout()` controls pending/error state, closes the menu, and navigates after success.
- `src/components/SideMenu.tsx`: renders the logout button, pending label, and error message using props from `App`.
- `src/App.css`: logout button colors and hover, pressed, keyboard focus, and disabled states.
- `scripts/check-auth.mjs`: checks request contracts, failure recovery, duplicate logout calls, and refresh/logout ordering.

## Checks

```sh
npm run build
npm test
npm run lint
```

## Formatting

```sh
npm run format        # Format frontend source, styles and configuration
npm run format:check  # Check formatting without changing files
```

Prettier uses `.prettierrc.json`: 2-space indentation, single quotes, no semicolons,
and a 100-character line width. Generated output, dependencies, and lockfiles are excluded.
For formatting on save, select the Prettier extension in your editor and enable format on save.

### Client-side routes

React Router handles navigation without reloading the document:

- `/`: public landing page; signed-in users continue to `/profile`.
- `/profile`: profile and Rooms (`/rooms` redirects here).
- `/settings`: account settings.
- `/characters` and `/friends`: demo pages.
- `/login`: login dialog; successful login returns to the requested private page.

Private routes wait for session restoration before redirecting to login. Unknown
paths show a page-not-found view. Docker's Nginx configuration serves `index.html`
for frontend routes so direct links and refreshes work; `/api/` is proxied separately.

## UI languages

The UI supports English and Russian. The text selector in the header is available before sign-in. After sign-in,
interface language is available in Settings instead of the header.
The saved `dicedasher.language` preference takes priority over the browser's primary
language; unsupported languages fall back to English. Changes also sync between tabs.
Language selection continues to work when browser storage is unavailable.

Translations live in `src/locales/en.json` and `src/locales/ru.json`. Add the same
key to the English dictionary and render it with `t(key)` from `src/lib/i18n.ts`.
Add the Russian translation under the same key when ready. Missing or blank
Russian translations fall back to English, including placeholder substitution.
Use `{name}` placeholders and `t(key, { name })` for dynamic values rather than
joining translated fragments. `plural()` selects one/few/many/other keys using
`Intl.PluralRules` for the current UI language.

Keep translation keys in state, validation errors and demo data; translate them
when rendering, so existing messages update when the language changes. API errors
use `UserFacingError` keys, while unexpected errors use safe fallback keys.
Usernames and other user-provided content are not translated. UI language remains
independent of any future generated-content language and does not change API requests.

`npm test` checks dictionary coverage, interpolation, language detection, saved
preferences, storage failures and cross-tab synchronization, alongside auth checks.

## Color themes

Signed-in users can choose a theme in Settings. The 16 palettes are:

- Light: Cobalt, Parchment, Blossom (Цветение), Cyberpunk Light, Chthonic Horror, The Loop, The Flood,
  Nuclear (Атомная эра), Field Manual, and Noir.
- Dark: Midnight, Forest, Cyberpunk Dark, Vampire, Cosmic Horror, and Classic Dark.

The choice is stored in browser local storage as `dicedasher.theme`, independently
of UI language and account data. It applies across the app and syncs between tabs.
Without a valid saved choice, the app initially selects Midnight or Cobalt based
on `prefers-color-scheme`. Clearing the preference in another tab restores that
system-based default. If storage is unavailable, selection still applies to the
current page, and Settings reports that it could not be saved.

### Adding or changing a palette

1. Define its ID, translation keys, and preview colors in `src/lib/themes.ts`.
2. Add a `:root[data-theme='…']` block in `src/themes.css`. Set `color-scheme` and
   every semantic token, including text, surfaces, borders, selection, status
   colors, hover states, and shadows. Components use these tokens through CSS
   custom properties; avoid adding palette-specific colors to component styles.
3. Keep preview colors in `themes.ts` consistent with `--page` and `--accent`.
4. Add the ID to the startup allowlist in `index.html`, so a saved palette applies
   before the first paint and survives direct navigation or reloads.
5. Add its name and description to both locale dictionaries. Translate theme
   metadata when rendering, just like other UI text.

Check text contrast against every surface where it appears, including button
labels, hover states, and status badges; target at least 4.5:1 for normal text.
Check focus indicators and control boundaries separately. Color must not be the
only way to distinguish errors or results, especially in red-accented palettes
and the monochrome Noir theme.

Run the frontend checks listed above and inspect the palette on real screens at
compact and wide widths. The current automated suite renders the theme selector
in both languages, but does not measure palette contrast or verify visual layout.
