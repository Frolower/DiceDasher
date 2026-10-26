# In progress

- Backend
  - [x] Registration — POST /register, without a session; [documentation](services/backend/README.md)
  - [x] Login — username + password, access JWT and refresh cookie
  - [x] Session creation — on successful login; stores the refresh token hash
  - [x] Password change — PUT /password with old and new passwords
  - [x] Profile lookup and updates — GET /me and PATCH /user
  - [x] Account deletion — DELETE /user; sessions cascade on deletion
  - [x] UUID v7 — defaults for UUID primary keys and application-generated session IDs
  - [x] Automatic updated_at — database triggers for users and characters
  - [x] Session revocation — POST /logout, refresh cookie removal
  - [x] Token refresh — POST /refresh with refresh cookie rotation
  - [ ] Login rate limiting
  - TODO: timeouts after deployment to the server
- Frontend
  - [x] Landing page, registration, and login forms
  - [ ] Session restoration, refresh, logout, and authenticated app screens
