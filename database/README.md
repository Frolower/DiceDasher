# Database initialization

Scripts in `init/` run in filename order when PostgreSQL starts with a new volume.
They do not run automatically for an existing volume.

## IDs and PostgreSQL version

Compose builds PostgreSQL 18 Alpine with pg_cron using `database/Dockerfile`.
UUID primary keys in the roll history, character, TES
player-car, user, and session tables default to PostgreSQL's built-in `uuidv7()`.
The backend generates session IDs before insertion with Go's `uuid.NewV7()` so
it can sign the access JWT first. IDs remain UUID strings in API responses.
Integer IDs in reference and pregenerated-character tables remain integers.

Changing init scripts does not migrate an existing volume. Existing databases
need PostgreSQL 18 and explicit schema updates for the new defaults and timestamp
triggers; existing UUID values can remain unchanged. Recreating a disposable
volume applies all init scripts, but deletes its data.

Deleting a user from `user_db.public.users` deletes their sessions through
`ON DELETE CASCADE`. Characters live in a separate database; account deletion
does not delete character data.

## Session retention and scheduled cleanup

`init/005_session_cleanup.sql` installs pg_cron in `user_db` and registers
`cleanup-old-sessions`. Every day at **03:00 UTC**, it deletes sessions whose
`LEAST(expires_at, revoked_at)` is more than 30 days old. PostgreSQL ignores NULL
arguments to LEAST, so sessions without revocation use their expiration date;
early logout starts retention at revocation. Active sessions are retained.
Daily scheduling means completed sessions remain for approximately 30–31 days.
Expiry and revocation checks in the backend take effect immediately and do not
depend on cleanup. Account deletion still removes all its sessions immediately.

An expression index supports the cleanup predicate. Jobs are owned by `postgres`;
the application role does not receive DELETE or scheduler permissions.
`cleanup-cron-history` runs at 03:15 UTC and removes completed execution logs
older than 30 days from `cron.job_run_details`.

Compose preloads pg_cron, sets its metadata database to `user_db` and timezone to
UTC, and uses background workers so jobs do not require a separate database
password or relaxed connection authentication. Outside Compose, install pg_cron
for PostgreSQL 18 and set these server options before applying script 005:

```conf
shared_preload_libraries = 'pg_cron'
cron.database_name = 'user_db'
cron.timezone = 'UTC'
cron.use_background_workers = on
```

Preserve any other entries already present in `shared_preload_libraries`.
See the [pg_cron setup documentation](https://github.com/citusdata/pg_cron#setting-up-pg_cron).

### Existing Compose volumes

Rebuild and recreate only the PostgreSQL container, retaining its volume, then
apply the new script explicitly as `postgres` (brief database downtime):

```sh
docker compose up -d --build --wait postgres
docker compose exec -T postgres psql -X -v ON_ERROR_STOP=1 -U postgres -d postgres < database/init/005_session_cleanup.sql
```

Do not recreate the volume or rerun all initialization scripts. Script 005 can
be reapplied as `postgres`: named jobs are updated rather than duplicated.
Creating the index can block writes briefly on an existing table; schedule this
step during a maintenance window if the session table is large.
The next scheduled run also removes existing sessions beyond the retention limit.

Inspect schedules and recent results:

```sh
docker compose exec -T postgres psql -X -U postgres -d user_db -c "SELECT jobid, jobname, schedule, database, username, active FROM cron.job;"
docker compose exec -T postgres psql -X -U postgres -d user_db -c "SELECT jobid, status, return_message, start_time, end_time FROM cron.job_run_details ORDER BY start_time DESC LIMIT 20;"
```

Verify retention of active, recently completed, and boundary sessions. The test
executes the registered cleanup SQL inside a transaction and rolls it all back:

```sh
docker compose exec -T postgres psql -X -v ON_ERROR_STOP=1 -U postgres -d postgres < database/tests/session_cleanup.sql
```

## Automatic update timestamps

`init/004_backend.sql` and `init/002_character.sql` create a
`BEFORE UPDATE FOR EACH ROW` trigger on `user_db.public.users` and
`character_db.public.characters`, respectively. These are the tables
that currently have an `updated_at` column. Each database has its own
`public.set_updated_at()` function.

The function definition lives in `init/functions/set_updated_at.sql`. Each schema
script includes it after connecting to its database using
`\ir functions/set_updated_at.sql`. Trigger declarations remain next to their
tables. Docker runs the top-level init scripts; the `functions/` directory holds
included files rather than standalone initialization steps.

When running a schema script manually, use `psql -f path/to/script.sql` so `\ir`
resolves paths relative to that script. Keep the `functions/` directory beside
the schema scripts when preparing a separate test database.

The triggers set `updated_at = NOW()` for every updated row, including updates
that assign unchanged values. `NOW()` is the transaction start time. Insert
defaults are preserved: users start with NULL, characters start with NOW().
Application queries should omit `updated_at`; the trigger overwrites any explicit
value supplied by an UPDATE. Failed or rolled-back updates do not persist a new
timestamp. New tables with `updated_at` need their own trigger.

Verify timestamps and service-role permissions (all test rows are rolled back):

```sh
docker compose exec -T postgres psql -X -v ON_ERROR_STOP=1 -U postgres -d postgres < database/tests/updated_at.sql
```
