-- Initializing resolver db
\connect resolve_db

-- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Allow the service role to connect to this DB
GRANT CONNECT ON DATABASE resolve_db TO resolve_service;

-- Allow the role to use objects in public schema
GRANT USAGE ON SCHEMA public TO resolve_service;

-- Enums
CREATE TYPE ttrpg_system AS ENUM ('generic', 'tes', 'vtmv5');
CREATE TYPE roll_action AS ENUM ('roll', 'push', 'reroll', 'check');

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
CREATE INDEX idx_roll_history_campaign ON public.roll_history(campaign_id) WHERE campaign_id IS NOT NULL;
CREATE INDEX idx_roll_history_character ON public.roll_history(character_id) WHERE character_id IS NOT NULL;
CREATE INDEX idx_roll_history_created ON public.roll_history(created_at);
CREATE INDEX idx_roll_history_system ON public.roll_history(system_name);

-- Grants
GRANT SELECT, INSERT ON TABLE public.roll_history TO resolve_service;