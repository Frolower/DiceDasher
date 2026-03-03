-- Resolve service database schema
-- Connect to resolve_db first
connect resolve_db

-- Grant schema permissions to service user
GRANT ALL ON SCHEMA public TO resolve_service;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO resolve_service;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO resolve_service;

-- Roll history table: stores results of dice rolls
CREATE TABLE roll_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system_name VARCHAR(50) NOT NULL,
    roll_type VARCHAR(50) NOT NULL,
    request_payload JSONB NOT NULL,
    result_payload JSONB NOT NULL,
    session_id UUID,
    character_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_roll_history_session ON roll_history(session_id) WHERE session_id IS NOT NULL;
CREATE INDEX idx_roll_history_character ON roll_history(character_id) WHERE character_id IS NOT NULL;
CREATE INDEX idx_roll_history_created ON roll_history(created_at);
CREATE INDEX idx_roll_history_system ON roll_history(system_name);
