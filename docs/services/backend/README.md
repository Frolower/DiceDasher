# Backend: Authentication and User Management

Backend: `http://localhost:8000`. Registration does not create a session. Signing in
with a username and password creates a separate session and issues an access JWT
and a refresh cookie.

## Code Organization

| File / package | Responsibility |
|---|---|
| `cmd/backend-service/main.go` | Wiring the repository, bcrypt, JWTManager, auth.Service, user.Service, and handler |
| `internal/handler/handler.go` | HTTP Authenticator and User interfaces, dependencies, and cookie settings |
| `internal/handler/register.go` | Registration HTTP contract |
| `internal/handler/me.go` | Authenticated profile HTTP contract; delegates to user.Service |
| `internal/handler/user.go` | Profile update and account deletion HTTP contracts |
| `internal/handler/password.go` | Password change HTTP contract |
| `internal/user/password.go` | Password validation, bcrypt verification, and hashing |
| `internal/repository/password.go` | Password hash lookup and conditional update |
| `internal/handler/auth_middleware.go` | Bearer token and session validation; authenticated claims in context |
| `internal/handler/refresh.go` | Token refresh HTTP contract and cookie rotation |
| `internal/user/profile.go`, `model.go` | Profile lookup and public user model |
| `internal/user/user.go`, `validate.go` | Profile update, normalization, and validation |
| `internal/user/service.go`, `errors.go` | User storage interface, dependencies, and domain errors |
| `internal/handler/login.go` | Login and logout HTTP handling, Origin validation, tokens, and cookie removal |
| `internal/handler/cookie.go` | Refresh cookie |
| `internal/handler/cors.go` | CORS for a specific Origin with credentials |
| `internal/auth/service.go` | Storage interfaces and constructors |
| `internal/auth/register.go` | Registration and bcrypt hash storage |
| `internal/auth/login.go` | User lookup and password verification |
| `internal/auth/logout.go` | Session revocation by refresh token hash |
| `internal/auth/refresh.go`, `authenticate.go` | Refresh rotation and access/session validation |
| `internal/auth/session.go` | Session creation after password verification |
| `internal/auth/jwt.go` | Access JWT issuance and verification |
| `internal/auth/refresh_token.go` | Cryptographically random refresh tokens and SHA-256 hashing |
| `internal/auth/errors.go`, `model.go`, `validate.go` | Errors, models, and validation |
| `internal/repository/user.go`, `session.go` | User and session SQL operations |

The HTTP layer does not execute SQL or bcrypt operations. The service does not set
cookies. The repository does not issue tokens. Context carries request cancellation
and logging data.

`auth` owns registration, login, logout, refresh, and authentication. `user` owns
profile lookup, updates, password changes, and account deletion. Both use the same repository through their own storage
interfaces. Protected routes authenticate through auth before invoking user.

## GET /me

Set `ACCESS_TOKEN` to the `access_token` returned by login or refresh.

```sh
curl -i http://localhost:8000/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

Requires a Bearer access token validated by the auth middleware. The handler passes
only the authenticated user ID to `user.Service.GetProfile`; profile data and
lookup errors belong to the user package. Registration remains in auth.

Success: `200` with `id`, `username`, and `email`. Missing/invalid authentication
or a missing user returns `401`; storage failures return `500`. Responses use
`Cache-Control: no-store` and do not change cookies.

## PATCH /user

Requires the same Bearer authentication as `/me`. The target user ID comes from
verified claims, never from the request body or query string.

```sh
curl -i -X PATCH http://localhost:8000/user \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"username":"AliceNew","email":"alice.new@example.com"}'
```

Send `username`, `email`, or both. Missing fields, empty strings, and JSON `null`
leave the corresponding field unchanged. At least one field must be nonempty
after normalization; `{}` and two blank values return `400`. Unknown fields,
including `id`, are rejected.

- Username is trimmed, then checked for 5–64 ASCII characters (English letters, digits and printable ASCII symbols) and no whitespace
  or control characters.
- Email is trimmed, lowercased, and checked for format and 5–254 bytes.
- PostgreSQL unique indexes on `lower(username)` and `lower(email)` reject
  conflicts with other users. Keeping one's own current values is allowed.
- Both changes are made in one SQL statement. A conflict rejects the entire
  update; `updated_at` changes only when the update succeeds.

Success: `200` with the JSON string `"success"`. Responses use
`Cache-Control: no-store` and do not change cookies. Errors use `text/plain`.

| Status | Reason |
|---|---|
| 400 | Invalid JSON, invalid fields, or no values to update |
| 401 | Missing/invalid access token or inactive session |
| 404 | Authenticated user no longer exists |
| 409 | Username or email is already in use |
| 413 | Body exceeds 16 KiB |
| 500 | Internal failure; storage details are hidden |

## DELETE /user

Requires Bearer authentication. The target ID comes from verified claims. No
request body is required.

```sh
curl -i -X DELETE http://localhost:8000/user \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

Success: `200` with the JSON string `"success"`. The user is deleted, all their
sessions are deleted by the database foreign key cascade, and the refresh cookie
is expired using the configured name and security attributes. The client must
clear its access token locally. Subsequent protected requests fail authentication.
Character records in the separate character database are not deleted.

Responses use `Cache-Control: no-store` and `Pragma: no-cache`. Errors use
`text/plain`: `401` for missing/invalid authentication or an inactive session,
`404` if the user disappears after authentication, and `500` for internal failures.
Repeated deletion with the same token normally returns `401` because its session
has already been deleted.

## PUT /password

Requires Bearer authentication. Send both passwords; the user ID comes from
verified claims, and unknown JSON fields are rejected.

```sh
curl -i -X PUT http://localhost:8000/password \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"old_password":"StrongPass1!","new_password":"NewStrongPass2!"}'
```

Both passwords currently undergo the registration complexity checks: 8–72 ASCII
characters, with a lowercase letter, uppercase letter, digit, and printable special
character; no whitespace or control characters. Neither password is trimmed.
This means an older password that fails today's complexity rules cannot be used
here, even if login still accepts it.

The service verifies the old password with bcrypt, hashes the new password, and
updates the stored hash only if it still matches the hash read for verification.
A concurrent password change that replaces that hash results in `404`. The database
trigger updates `updated_at`. Existing sessions, access JWTs, and refresh cookies
remain valid; this operation does not revoke sessions or rotate tokens.

Success: `200` with the JSON string `"success"`, `Cache-Control: no-store`, and
`Pragma: no-cache`. Cookies are unchanged. Errors use `text/plain`.

| Status | Reason |
|---|---|
| 400 | Invalid JSON, unknown/missing fields, password validation failure, or incorrect old password |
| 401 | Missing/invalid access token or inactive session |
| 404 | User not found or stored password hash changed concurrently |
| 413 | Body exceeds 16 KiB |
| 500 | Internal failure; storage details are hidden |

## POST /register

Input: `username`, `email`, `password`. Success: `201 {"status":"created"}` without cookies.

- Username: 5–64 ASCII characters; English letters, digits and printable ASCII symbols only, with no whitespace or control characters. Cyrillic, accented letters and emoji are rejected. Not trimmed.
- Email: 5–254 bytes; format checked (the domain must contain at least one dot, e.g. `example.com`) and converted to lowercase; surrounding whitespace is rejected.
- Password: 8–72 ASCII characters; only English letters (`a-z`, `A-Z`), digits (`0-9`), and printable ASCII special characters are allowed. Must include a lowercase letter, an uppercase letter, a digit, and a special character (including `$`, `^`, or `~`). Whitespace, control characters, Cyrillic, accented letters, and emoji are rejected. Passed unchanged.
- Username and email are case-insensitively unique through PostgreSQL `lower(...)`.
- Errors: 400 — invalid input, 409 — duplicate, 413 — body exceeds 16 KiB, 500 — internal failure.

## POST /login

```sh
curl -i -c /tmp/dicedasher-cookies.txt http://localhost:8000/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"Alice","password":"StrongPass1!"}'
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
CORS allows only `ALLOWED_ORIGINS`. Login accepts only `application/json` and rejects
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

Login accepts only English letters, digits and printable ASCII symbols in usernames and passwords, without spaces. Login does not enforce the new registration password complexity rules: existing
passwords must continue to work. Passwords are verified through
bcrypt.CompareHashAndPassword rather than by hashing them again. For a missing
user, a dummy hash is also checked with bcrypt to reduce response timing differences.

## POST /logout

```sh
curl -i -X POST -b /tmp/dicedasher-cookies.txt -c /tmp/dicedasher-cookies.txt \
  http://localhost:8000/logout
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

## POST /refresh

```sh
curl -i -X POST -b /tmp/dicedasher-cookies.txt -c /tmp/dicedasher-cookies.txt \
  http://localhost:8000/refresh
```

No body or Content-Type is required. Refresh reads the same HttpOnly cookie as
logout and applies the same Origin and cross-site request checks. Browser requests
must use `credentials: "include"`.

Success: `200` with `access_token`, `token_type`, and `expires_at`, as for login,
and a rotated refresh cookie. The old refresh token stops working. Rotation uses
a conditional SQL update so concurrent requests with the same token cannot both
succeed. The session's original expiry is retained; the new access token lasts
up to 10 minutes, capped at that expiry.

Errors: `400` — invalid cookie; `401` — absent/unknown refresh token or an expired,
revoked, or concurrently rotated session; `403` — request origin not allowed;
`500` — internal failure. Responses use `Cache-Control: no-store` and
`Pragma: no-cache`.

## Tokens and Sessions

An access JWT (HS256) is valid for 10 minutes. It contains sub, sid, iss, aud, iat,
exp, and jti. Verification is restricted to HS256, the expected issuer/audience,
and mandatory exp/iat claims; user and session UUIDs are validated.
The JWT is signed, but its contents are not encrypted.

A refresh token consists of 32 random bytes encoded as base64url. Only its SHA-256
hash is stored in the database, not the original token. Sessions last 30 days.
Each login creates a new session. User IDs default to PostgreSQL 18 `uuidv7()`.
The service generates a UUID v7 session ID before
signing the JWT and performing the INSERT.

Tokens are prepared first, then the session is saved, then the response is sent.
A signing failure does not create a database row; a storage failure does not issue
tokens. A network interruption after a successful INSERT can leave a session that
the client never received; login is not idempotent.

Protected requests verify the access JWT and load its session by `sid`. Revoked
or expired sessions and mismatched user IDs are rejected. Logout therefore blocks
subsequent protected requests even when their access JWT has not expired.

Refresh rotation is implemented. Used refresh tokens are not retained in a
history, so an old token is rejected without revoking the entire session family.

## Configuration and Startup

Use the values from `services/backend/.env.example`. Create a local `.env` and set
JWT_SECRET once to the output of `openssl rand -base64 32`. Do not commit the key.
If a local key already exists, retain it across restarts; changing the key
invalidates existing access JWTs. An empty or too-short key prevents startup.

| Variable | Value |
|---|---|
| DATABASE_URL | Required URL for user_db |
| HTTP_ADDR | Defaults to :8000 |
| JWT_SECRET | Required base64 secret containing at least 32 random bytes |
| JWT_ISSUER | dicedasher-backend |
| JWT_AUDIENCE | dicedasher-api |
| COOKIE_SECURE | Defaults to true; false only for local HTTP |
| ALLOWED_ORIGINS | `http://localhost:8888,http://localhost:5174`, without trailing slashes |
| LOG_MODE | default |

Compose reads `services/backend/.env`, but overrides the database address with the
internal `postgres:5432/user_db` address and the user_service role. It does not
override `ALLOWED_ORIGINS`. Set it to a comma-separated list of exact browser
origins without paths or trailing slashes. The default is `http://localhost:8888,http://localhost:5174`
(Swagger and frontend). Resolve and character default to `http://localhost:8888,http://localhost:8000`.
Use `COOKIE_SECURE=false` for local HTTP.
Running Go directly does not automatically read .env.

Start PostgreSQL and the backend:

```sh
docker compose up -d postgres
docker compose up -d --build backend
```

On a new volume, PostgreSQL runs the scripts in database/init in filename order.
Script 000 creates the databases and roles; 004 creates the backend tables,
indexes, permissions, and the users timestamp trigger. Script 002 creates the
characters timestamp trigger. The database sets updated_at on every matching
UPDATE, including assignments of unchanged values. NOW() uses the transaction
start time. Users start with NULL in updated_at until the first update.

In production, use HTTPS, COOKIE_SECURE=true, and your own Origin.
The HS256 key allows both verifying and issuing JWTs; do not share it with clients.
Login rate limiting has not been added yet; a rate limit at the ingress or
application level is required before a public launch.

## Checks

```sh
go test ./services/backend/...
go vet ./services/backend/...
```

For integration tests, create a separate test database and service role, apply 004
(adjust its database name if needed), and set:

```sh
BACKEND_TEST_DATABASE_URL='postgres://user_service:user_password@localhost:55439/user_db?sslmode=disable' \
BACKEND_TEST_ADMIN_URL='postgres://YOUR_ADMIN@localhost:55439/user_db?sslmode=disable' \
go test -race ./services/backend/...
```

Logout unit tests cover token hashing, session revocation, cookie removal,
repeated logout, service errors, and Origin validation. There is no logout
integration test yet.

Tests also cover registration, bcrypt, JWTs, refresh rotation, session revocation,
profile lookup, update validation/normalization, and protected-route HTTP contracts.
Update tests check authenticated identity, malformed/oversized JSON, and
400/404/409/500 responses without internal error details.

Password tests cover validation, old-password verification, hashing, conditional
storage updates, authenticated HTTP handling, and errors. The password PostgreSQL
integration test checks rejected changes, persisted bcrypt hashes, old-password
rejection after a change, and rejection of a stale hash update.

PostgreSQL integration tests cover registration, login, password changes, and profile updates under
the service role. Update tests check one-field and two-field changes, uniqueness,
atomicity on conflicts, updated_at, and unchanged profiles of other users.
They are skipped unless both database URLs are set. Test users have random names
and are deleted along with their sessions.
