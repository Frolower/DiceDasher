\set ON_ERROR_STOP on
-- Повторное применение к существующему volume не удаляет данные.
-- Скрипт также можно запустить отдельно от 000 (под администратором).
SELECT 'CREATE ROLE user_service LOGIN PASSWORD ''user_password'''
WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'user_service')
\gexec
SELECT 'CREATE DATABASE user_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'user_db')
\gexec

\connect user_db

CREATE EXTENSION IF NOT EXISTS pgcrypto;

GRANT CONNECT ON DATABASE user_db TO user_service;
GRANT USAGE ON SCHEMA public TO user_service;

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
GRANT SELECT (id, username, hash) ON public.users TO user_service;
GRANT DELETE ON public.users TO user_service;

GRANT INSERT ON public.sessions TO user_service;
GRANT SELECT ON public.sessions TO user_service;
GRANT UPDATE (refresh_token_hash, revoked_at) ON public.sessions TO user_service;

-- Хеш однозначно определяет сессию; user_id нужен для списка сессий пользователя.
CREATE UNIQUE INDEX IF NOT EXISTS sessions_refresh_token_hash_unique ON public.sessions (refresh_token_hash);
CREATE INDEX IF NOT EXISTS sessions_user_id_idx ON public.sessions (user_id);
