\connect user_db

CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Keep sessions for 30 days after expiration or early revocation.
CREATE INDEX IF NOT EXISTS sessions_ended_at_idx
    ON public.sessions (LEAST(expires_at, revoked_at));

-- Named schedules are updated on reapplication by the same owner (postgres).
-- The table owner runs maintenance; user_service keeps its existing permissions.
SELECT cron.schedule(
    'cleanup-old-sessions',
    '0 3 * * *',
    $$DELETE FROM public.sessions
      WHERE LEAST(expires_at, revoked_at) < now() - interval '30 days'$$
);

-- pg_cron does not automatically prune its execution history.
SELECT cron.schedule(
    'cleanup-cron-history',
    '15 3 * * *',
    $$DELETE FROM cron.job_run_details
      WHERE end_time < now() - interval '30 days'$$
);
