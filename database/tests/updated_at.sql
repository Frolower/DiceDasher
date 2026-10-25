-- Run as the local database administrator after database initialization.
-- Fixtures are isolated in transactions and never committed.
\set ON_ERROR_STOP on

\connect user_db
BEGIN;

DO $$
DECLARE
    test_id uuid := uuidv7();
    test_name text := 'Timestamp' || replace(test_id::text, '-', '');
    stamp timestamptz;
    saved_hash text;
BEGIN
    INSERT INTO public.users (id, username, email, hash)
    VALUES (test_id, test_name, test_name || '@example.com', 'initial-hash')
    RETURNING updated_at INTO stamp;
    IF stamp IS NOT NULL THEN
        RAISE EXCEPTION 'users: insert default must remain NULL';
    END IF;

    SET LOCAL ROLE user_service;
    UPDATE public.users SET username = test_name || 'New' WHERE id = test_id;
    RESET ROLE;
    SELECT updated_at INTO stamp FROM public.users WHERE id = test_id;
    IF stamp IS DISTINCT FROM NOW() THEN
        RAISE EXCEPTION 'users: profile update did not set updated_at';
    END IF;

    SET LOCAL ROLE user_service;
    UPDATE public.users SET hash = 'new-hash' WHERE id = test_id AND hash = 'initial-hash';
    RESET ROLE;
    SELECT hash, updated_at INTO saved_hash, stamp FROM public.users WHERE id = test_id;
    IF saved_hash <> 'new-hash' OR stamp IS DISTINCT FROM NOW() THEN
        RAISE EXCEPTION 'users: password update did not set hash and updated_at';
    END IF;

    -- Explicit timestamps cannot bypass the trigger, even for a no-op update.
    SET LOCAL ROLE user_service;
    UPDATE public.users SET hash = hash, updated_at = '2000-01-01' WHERE id = test_id;
    RESET ROLE;
    SELECT updated_at INTO stamp FROM public.users WHERE id = test_id;
    IF stamp IS DISTINCT FROM NOW() THEN
        RAISE EXCEPTION 'users: trigger did not override an explicit timestamp';
    END IF;
END;
$$;

ROLLBACK;

\connect character_db
BEGIN;

DO $$
DECLARE
    test_id uuid := uuidv7();
    stamp timestamptz;
BEGIN
    INSERT INTO public.characters (id, user_id, system_name, character_type, name)
    VALUES (test_id, uuidv7(), 'generic', 'pc', 'Timestamp test')
    RETURNING updated_at INTO stamp;
    IF stamp IS DISTINCT FROM NOW() THEN
        RAISE EXCEPTION 'characters: insert default must remain NOW()';
    END IF;

    SET LOCAL ROLE character_service;
    UPDATE public.characters SET name = 'Updated character' WHERE id = test_id;
    RESET ROLE;
    SELECT updated_at INTO stamp FROM public.characters WHERE id = test_id;
    IF stamp IS DISTINCT FROM NOW() THEN
        RAISE EXCEPTION 'characters: update did not set updated_at';
    END IF;

    SET LOCAL ROLE character_service;
    UPDATE public.characters SET name = name, updated_at = '2000-01-01' WHERE id = test_id;
    RESET ROLE;
    SELECT updated_at INTO stamp FROM public.characters WHERE id = test_id;
    IF stamp IS DISTINCT FROM NOW() THEN
        RAISE EXCEPTION 'characters: trigger did not override an explicit timestamp';
    END IF;
END;
$$;

ROLLBACK;
