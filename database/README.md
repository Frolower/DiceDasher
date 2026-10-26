# Database initialization

Scripts in `init/` run in filename order when PostgreSQL starts with a new volume.
They do not run automatically for an existing volume.

## IDs and PostgreSQL version

Compose uses PostgreSQL 18. UUID primary keys in the roll history, character, TES
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
