\connect user_db

\ir functions/set_updated_at.sql

CREATE EXTENSION IF NOT EXISTS pgcrypto;

GRANT CONNECT ON DATABASE user_db TO user_service;
GRANT USAGE ON SCHEMA public TO user_service;

CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    username VARCHAR(64) NOT NULL CHECK (char_length(username) BETWEEN 5 AND 64),
    email VARCHAR(254) NOT NULL,
    hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON public.users
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_at();

CREATE TABLE IF NOT EXISTS public.sessions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    refresh_token_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '30 DAYS'),
    revoked_at TIMESTAMPTZ,
    user_agent TEXT
);


CREATE UNIQUE INDEX IF NOT EXISTS users_username_unique ON public.users (lower(username));
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON public.users (lower(email));

GRANT INSERT ON public.users TO user_service;
GRANT SELECT (id, username, email, hash) ON public.users TO user_service;
GRANT DELETE ON public.users TO user_service;
GRANT UPDATE ON public.users TO user_service;

GRANT INSERT ON public.sessions TO user_service;
GRANT SELECT ON public.sessions TO user_service;
GRANT UPDATE (refresh_token_hash, revoked_at) ON public.sessions TO user_service;

-- The hash uniquely identifies a session; user_id is used to list a user's sessions.
CREATE UNIQUE INDEX IF NOT EXISTS sessions_refresh_token_hash_unique ON public.sessions (refresh_token_hash);
CREATE INDEX IF NOT EXISTS sessions_user_id_idx ON public.sessions (user_id);
