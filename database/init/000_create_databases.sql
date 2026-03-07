-- Create databases for each microservice
-- This script runs as postgres superuser during container initialization

CREATE DATABASE resolve_db;
CREATE DATABASE character_db;

-- Create dedicated users for each service
CREATE USER resolve_service WITH ENCRYPTED PASSWORD 'resolve_password';
CREATE USER character_service WITH ENCRYPTED PASSWORD 'character_password';