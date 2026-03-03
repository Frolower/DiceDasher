-- Character service database schema
-- Connect to character_db first
\connect character_db

-- Grant schema permissions to service user
GRANT ALL ON SCHEMA public TO character_service;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO character_service;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO character_service;

-- Characters table
CREATE TABLE characters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    system_name VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_characters_user ON characters(user_id);
CREATE INDEX idx_characters_system ON characters(system_name);

-- Character templates (predefined character archetypes)
CREATE TABLE character_templates (
    id SERIAL PRIMARY KEY,
    system_name VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    template_data JSONB NOT NULL,
    is_official BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_templates_system ON character_templates(system_name);
