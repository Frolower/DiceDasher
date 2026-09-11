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


CREATE UNIQUE INDEX IF NOT EXISTS users_username_unique ON public.users (lower(username));
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON public.users (lower(email));

GRANT INSERT ON public.users TO user_service;
GRANT SELECT (id) ON public.users TO user_service;
