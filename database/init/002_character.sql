-- Character service database schema
\connect character_db

-- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Allow the service role to connect to this DB
GRANT CONNECT ON DATABASE character_db TO character_service;

-- Enums
CREATE TYPE ttrpg_system AS ENUM ('generic', 'tes', 'vtmv5');
CREATE TYPE character_type AS ENUM ('pc', 'npc');

-- Characters table
CREATE TABLE public.characters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    system_name ttrpg_system NOT NULL,
    character_type character_type NOT NULL,
    name VARCHAR(255) NOT NULL,
    data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_characters_user ON public.characters(user_id);
CREATE INDEX idx_characters_system ON public.characters(system_name);
CREATE INDEX idx_character_types_system ON public.characters(character_type);

-- Pregenerated characters table
CREATE TABLE public.pregen_characters (
    id SERIAL PRIMARY KEY,
    system_name ttrpg_system NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    template_data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_templates_system ON public.pregen_characters(system_name);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.characters TO character_service;
GRANT SELECT ON TABLE public.pregen_characters TO character_service;