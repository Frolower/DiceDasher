# DiceDasher frontend

React + TypeScript + Vite. Minimal landing page with a rotating d20 and account forms.

## Docker

From the repository root:

```sh
docker compose up -d --build frontend
```

Open http://localhost:5174. Compose also starts PostgreSQL and the backend.
The image builds the app with Node and serves only the resulting static files with Nginx.
Nginx forwards `/api/*` to the backend on the internal Docker network, stripping `/api`.
It preserves the browser Origin and cookies. `/health` checks the frontend itself.

After frontend changes, rebuild with the same command. To follow its logs:

```sh
docker compose logs -f frontend
```

## Development with live reload

```sh
cd frontend
npm install
npm run dev
```

Open http://localhost:5174. The port is fixed so the backend Origin check stays predictable.
Stop the frontend container first (`docker compose stop frontend` from the repository root)
when using Vite locally, because both use port 5174.

The Vite development proxy forwards `/api/*` to `http://localhost:8000/*`.
From the repository root, start or recreate the backend with the updated Compose configuration:

```sh
docker compose up -d --build backend
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
- On login, the landing page shows a welcome message. Access tokens are kept in React memory only; persistent sessions, refresh, logout, and authenticated app screens are a later step. Reloading the page clears the displayed login state.
- The central d20 rotates slowly; click it or press Enter/Space to pause or resume. Rotation pauses while an account dialog is open and respects the system reduced-motion preference.

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
