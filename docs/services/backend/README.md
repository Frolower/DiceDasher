# Backend: Registration, Login, and Logout

Backend: `http://localhost:8083`. Registration does not create a session. Signing in
with a username and password creates a separate session and issues an access JWT
and a refresh cookie.

## Code Organization

| File / package | Responsibility |
|---|---|
| `cmd/backend-service/main.go` | Wiring the repository, bcrypt, JWTManager, auth.Service, and handler |
| `internal/handler/handler.go` | HTTP Authenticator interface, dependencies, and cookie settings |
| `internal/handler/register.go` | Registration HTTP contract |
| `internal/handler/login.go` | Login and logout HTTP handling, Origin validation, tokens, and cookie removal |
| `internal/handler/cookie.go` | Refresh cookie |
| `internal/handler/cors.go` | CORS for a specific Origin with credentials |
| `internal/auth/service.go` | Storage interfaces and constructors |
| `internal/auth/register.go` | Registration and bcrypt hash storage |
| `internal/auth/login.go` | User lookup and password verification |
| `internal/auth/logout.go` | Session revocation by refresh token hash |
| `internal/auth/session.go` | Session creation after password verification |
| `internal/auth/jwt.go` | Access JWT issuance and verification |
| `internal/auth/refresh_token.go` | Cryptographically random refresh tokens and SHA-256 hashing |
| `internal/auth/errors.go`, `model.go`, `validate.go` | Errors, models, and validation |
| `internal/repository/user.go`, `session.go` | User and session SQL operations |

The HTTP layer does not execute SQL or bcrypt operations. The service does not set
cookies. The repository does not issue tokens. Context carries request cancellation
and logging data.

## POST /register

Input: `username`, `email`, `password`. Success: `201 {"status":"created"}` without cookies.

- Username: 5–64 Unicode characters, with no whitespace or control characters; not trimmed.
- Email: 3–254 bytes; format checked, trimmed, and converted to lowercase.
- Password: at least 8 characters and at most 72 UTF-8 bytes; must include lowercase and uppercase letters and a digit, with no whitespace. Passed unchanged.
- Username and email are case-insensitively unique through PostgreSQL `lower(...)`.
- Errors: 400 — invalid input, 409 — duplicate, 413 — body exceeds 16 KiB, 500 — internal failure.

## POST /login

```sh
curl -i -c /tmp/dicedasher-cookies.txt http://localhost:8083/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"Alice","password":"StrongPass1"}'
```

Success: HTTP `200`, `Cache-Control: no-store`:

```json
{
  "access_token": "<signed JWT>",
  "token_type": "Bearer",
  "expires_at": "<UTC timestamp>"
}
```

The refresh token is sent only through `Set-Cookie`, not in JSON:

- HTTPS: `__Host-refresh_token`, Secure, HttpOnly, SameSite=Strict, Path=/, without Domain.
- Local HTTP with `COOKIE_SECURE=false`: `refresh_token`, HttpOnly, SameSite=Strict, Path=/.
- Expires and Max-Age are limited to the session lifetime.

Browser requests must use `credentials: "include"`, including the first login.
CORS allows only `ALLOWED_ORIGIN`. Login accepts only `application/json` and rejects
an unapproved Origin or cross-site Fetch Metadata without an Origin.
curl/Postman requests without an Origin are allowed. User-Agent is read from the
HTTP header; it is used for display, not as proof of device identity.

| Status | Reason |
|---|---|
| 200 | Password verified, session saved, and tokens issued |
| 400 | Invalid JSON or missing/oversized fields |
| 401 | Invalid username or password (the same response for both) |
| 403 | Origin not allowed |
| 413 | Body exceeds 16 KiB |
| 415 | Content-Type is not application/json |
| 500 | Internal error; no tokens or cookies issued |

Login does not enforce the new registration password complexity rules: existing
passwords must continue to work. Passwords are verified through
bcrypt.CompareHashAndPassword rather than by hashing them again. For a missing
user, a dummy hash is also checked with bcrypt to reduce response timing differences.

## POST /logout

```sh
curl -i -X POST -b /tmp/dicedasher-cookies.txt -c /tmp/dicedasher-cookies.txt \
  http://localhost:8083/logout
```

No request body or Content-Type is required. The refresh token is read from the
cookie: `__Host-refresh_token` when `COOKIE_SECURE=true`, otherwise `refresh_token`.
In browsers, use `credentials: "include"`. Like login, logout rejects an unapproved
Origin and `Sec-Fetch-Site: cross-site` without an Origin.

The service hashes the token and sets `revoked_at` on the corresponding session.
Success: `204` without a body, with `Cache-Control: no-store` and `Pragma: no-cache`.
The cookie is removed through Set-Cookie with the same name and Path=/, Max-Age=0,
and an Expires value in the past; HttpOnly, Secure, and SameSite match the login settings.
Repeated logout or an unknown token also returns `204` and removes the cookie.
If the cookie is absent, the response is `204` without Set-Cookie.

Errors: `400` — cookie read error, `403` — request origin not allowed,
`500` — internal revocation failure. On an internal failure, the cookie is not
removed so the client can retry. This handler does not remove the access JWT:
the client must clear it locally. JWT signature verification alone does not
account for session revocation in the database.

## Tokens and Sessions

An access JWT (HS256) is valid for 10 minutes. It contains sub, sid, iss, aud, iat,
exp, and jti. Verification is restricted to HS256, the expected issuer/audience,
and mandatory exp/iat claims; user and session UUIDs are validated.
The JWT is signed, but its contents are not encrypted.

A refresh token consists of 32 random bytes encoded as base64url. Only its SHA-256
hash is stored in the database, not the original token. Sessions last 30 days.
Each login creates a new session. The service generates the session ID before
signing the JWT and performing the INSERT.

Tokens are prepared first, then the session is saved, then the response is sent.
A signing failure does not create a database row; a storage failure does not issue
tokens. A network interruption after a successful INSERT can leave a session that
the client never received; login is not idempotent.

Logout is implemented and revokes the session in the database. Refresh currently
returns ErrNotImplemented at the service level, and its route is not wired yet.
The SQL method for rotation is ready, but a history of used refresh tokens and
protection against their reuse have not been implemented. Once an access JWT
expires, another login is currently required. JWT verification is implemented as
a component, but middleware protecting `/me` and other routes is the next step.
A signed JWT does not itself check revoked_at in the database.

## Configuration and Startup

Use the values from `services/backend/.env.example`. Create a local `.env` and set
JWT_SECRET once to the output of `openssl rand -base64 32`. Do not commit the key.
If a local key already exists, retain it across restarts; changing the key
invalidates existing access JWTs. An empty or too-short key prevents startup.

| Variable | Value |
|---|---|
| DATABASE_URL | Required URL for user_db |
| HTTP_ADDR | Defaults to :8083 |
| JWT_SECRET | Required base64 secret containing at least 32 random bytes |
| JWT_ISSUER | dicedasher-backend |
| JWT_AUDIENCE | dicedasher-api |
| COOKIE_SECURE | Defaults to true; false only for local HTTP |
| ALLOWED_ORIGIN | http://localhost:8082, without a trailing slash |
| LOG_MODE | default |

Compose reads `services/backend/.env`, but overrides the database address with the
internal `postgres:5432/user_db` address and the user_service role.
Running Go directly does not automatically read .env.

For an existing PostgreSQL volume:

```sh
docker compose up -d postgres
docker compose exec -T postgres psql -U postgres -d postgres < database/init/004_backend.sql
docker compose up -d --build backend
```

004 can be reapplied: it creates any missing user_db/role, tables, indexes, and
permissions. It does not change an existing role's password or migrate arbitrary
legacy table structures. Data from the former backend_db is not transferred
automatically. A new volume is initialized automatically through 000 and 004.
Do not delete the volume just to apply new GRANT statements or indexes.

In production, use HTTPS, COOKIE_SECURE=true, and your own Origin.
The HS256 key allows both verifying and issuing JWTs; do not share it with clients.
Login rate limiting has not been added yet; a rate limit at the ingress or
application level is required before a public launch.

## Checks

```sh
go test ./services/backend/...
go vet ./services/backend/...
```

For integration tests, apply 004 to a separate test database and set:

```sh
BACKEND_TEST_DATABASE_URL='postgres://user_service:user_password@localhost:55439/user_db?sslmode=disable' \
BACKEND_TEST_ADMIN_URL='postgres://YOUR_ADMIN@localhost:55439/user_db?sslmode=disable' \
go test -race ./services/backend/...
```

Logout unit tests cover token hashing, session revocation, cookie removal,
repeated logout, service errors, and Origin validation. There is no logout
integration test yet.

Tests also cover registration, bcrypt, JWTs, cookies, repeated login, password
errors, concurrent registration, and the role's actual SQL permissions.
Test users have random names and are deleted along with their sessions.
