-- Initializing resolver db

\connect resolver_db

-- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Allow the service role to connect to this DB
GRANT CONNECT, TEMP ON DATABASE resolver_db TO resolver_service;

-- Allow the role to use objects in public schema
GRANT USAGE ON SCHEMA public TO resolver_service;

-- Enums
DO $$ BEGIN
    CREATE TYPE ttrpg_system AS ENUM ('generic', 'tes', 'vtmv5');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE roll_action AS ENUM ('roll', 'push', 'reroll', 'check');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Roll history table
CREATE TABLE IF NOT EXISTS public.roll_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL,
    system_name ttrpg_system NOT NULL,
    action_type roll_action NOT NULL,
    request_payload JSONB,
    response_payload JSONB NOT NULL,
    campaign_id UUID,
    character_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
    );

-- Indexes
CREATE INDEX IF NOT EXISTS idx_roll_history_campaign
    ON public.roll_history(campaign_id) WHERE campaign_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_roll_history_character
    ON public.roll_history(character_id) WHERE character_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_roll_history_created
    ON public.roll_history(created_at);

CREATE INDEX IF NOT EXISTS idx_roll_history_system
    ON public.roll_history(system_name);

-- Table privileges
GRANT SELECT, INSERT ON TABLE public.roll_history TO resolver_service;