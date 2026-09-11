\connect user_db

CREATE EXTENSION IF NOT EXISTS pgcrypto;

GRANT CONNECT ON DATABASE user_db TO user_service;

CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(64) NOT NULL CHECK (char_length(username) BETWEEN 5 AND 64),
    email VARCHAR(254) NOT NULL,
    hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
GRANT SELECT (id) ON public.users TO user_service;
GRANT DELETE ON public.users TO user_service;

GRANT INSERT ON public.sessions TO user_service;
GRANT SELECT (user_id) ON public.sessions TO user_service;
