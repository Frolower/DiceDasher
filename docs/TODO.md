# In progress

- Backend
  - [x] Registration — POST /register, without a session; [documentation](services/backend/README.md)
  - [x] Login — username + password, access JWT and refresh cookie
  - [x] Session creation — on successful login; stores the refresh token hash
  - [ ] Password change
  - [x] Session revocation — POST /logout, refresh cookie removal
  - [ ] Token refresh
  - TODO: timeouts after deployment to the server
