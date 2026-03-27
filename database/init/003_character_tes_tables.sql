\connect character_db

-- Creating new schema

CREATE SCHEMA IF NOT EXISTS tes;

-- TES specific enums

CREATE TYPE tes.item_type AS ENUM ('weapon', 'neurocaster', 'armor', 'gear', 'drone');

-- TES specific schemas

CREATE TABLE tes.player_cars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    passengers INT NOT NULL,
    fuel VARCHAR(255) NOT NULL,
    description TEXT,
    maneuverability INT NOT NULL,
    speed INT NOT NULL,
    hull INT NOT NULL,
    armor INT NOT NULL,
    cost INT NOT NULL,
    traits TEXT[] NOT NULL DEFAULT '{}',
    gear TEXT[] NOT NULL DEFAULT '{}'
);

CREATE TABLE tes.car_stats (
    id SERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL UNIQUE,
    passengers INT,
    maneuverability INT,
    speed INT NOT NULL,
    hull INT NOT NULL,
    armor INT,
    cost INT
);

CREATE TABLE tes.items (
    id SERIAL PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type tes.item_type NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'
);

-- Indexes

CREATE INDEX idx_car_id ON tes.player_cars(id);
CREATE INDEX idx_car_type ON tes.car_stats(type);
CREATE INDEX idx_item_code ON tes.items(code);

-- Grants

GRANT USAGE ON SCHEMA tes TO character_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE tes.player_cars TO character_service;
GRANT SELECT ON TABLE tes.car_stats TO character_service;
GRANT SELECT ON TABLE tes.items TO character_service;