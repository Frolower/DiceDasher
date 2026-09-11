-- Скрипт подходит и для первого старта, и для существующего Docker volume.
-- Запускать через psql от администратора; существующие базы не удаляются.
\set ON_ERROR_STOP on

-- Имена и пароль предназначены для локального окружения Docker Compose.
SELECT 'CREATE ROLE backend_service LOGIN PASSWORD ''backend_password'''
WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'backend_service')
\gexec
SELECT 'CREATE DATABASE backend_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'backend_db')
\gexec

\connect backend_db
GRANT CONNECT ON DATABASE backend_db TO backend_service;
GRANT USAGE ON SCHEMA public TO backend_service;

CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(64) NOT NULL CHECK (char_length(username) BETWEEN 3 AND 64),
    email VARCHAR(254) NOT NULL,
    -- В этой колонке хранится только bcrypt-хеш, никогда открытый пароль.
    hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Имя сохраняет исходный регистр для отображения, но Alice и alice занять
-- одновременно нельзя. Индексы защищают и от конкурирующих INSERT.
CREATE UNIQUE INDEX IF NOT EXISTS users_username_unique ON public.users (lower(username));
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON public.users (lower(email));

-- Для текущего сценария нужны INSERT и чтение id через RETURNING.
GRANT INSERT ON public.users TO backend_service;
GRANT SELECT (id) ON public.users TO backend_service;
