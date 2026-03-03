-- Create databases for each microservice
-- This script runs as postgres superuser during container initialization

CREATE DATABASE resolve_db;
CREATE DATABASE character_db;

-- Create dedicated users for each service (optional but recommended)
CREATE USER resolver_service WITH ENCRYPTED PASSWORD 'resolver_password';
CREATE USER character_service WITH ENCRYPTED PASSWORD 'character_password';

GRANT ALL PRIVILEGES ON DATABASE resolve_db TO resolve_service;
GRANT ALL PRIVILEGES ON DATABASE character_db TO character_service;
