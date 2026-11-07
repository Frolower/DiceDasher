-- Run as postgres; all fixtures and deletions are rolled back.
\set ON_ERROR_STOP on
\connect user_db
BEGIN;

DO $$
DECLARE
    test_user uuid := uuidv7();
    test_name text := 'Cleanup' || replace(test_user::text, '-', '');
    cleanup_sql text;
    remaining text[];
BEGIN
    SELECT command INTO STRICT cleanup_sql
    FROM cron.job
    WHERE jobname = 'cleanup-old-sessions' AND database = 'user_db'
        AND username = 'postgres' AND schedule = '0 3 * * *' AND active;

    INSERT INTO public.users (id, username, email, hash)
    VALUES (test_user, test_name, test_name || '@example.com', 'test');

    INSERT INTO public.sessions
        (user_id, refresh_token_hash, created_at, expires_at, revoked_at, user_agent)
    SELECT test_user, test_name || label, now() - interval '90 days', expiry, revoked, label
    FROM (VALUES
        ('active', now() + interval '1 day', NULL::timestamptz),
        ('expired-recent', now() - interval '29 days', NULL::timestamptz),
        ('expired-boundary', now() - interval '30 days', NULL::timestamptz),
        ('expired-old', now() - interval '31 days', NULL::timestamptz),
        ('revoked-recent', now() + interval '1 day', now() - interval '29 days'),
        ('revoked-old', now() + interval '1 day', now() - interval '31 days'),
        ('expired-before-revoked', now() - interval '31 days', now() - interval '1 day')
    ) AS cases(label, expiry, revoked);

    EXECUTE cleanup_sql;
    SELECT array_agg(user_agent ORDER BY user_agent) INTO remaining
    FROM public.sessions WHERE user_id = test_user;
    IF remaining IS DISTINCT FROM ARRAY[
        'active', 'expired-boundary', 'expired-recent', 'revoked-recent'
    ]::text[] THEN
        RAISE EXCEPTION 'Unexpected retained sessions: %', remaining;
    END IF;
END;
$$;

ROLLBACK;
