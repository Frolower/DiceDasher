-- Create databases for each microservice
-- This script runs as postgres superuser during container initialization

CREATE DATABASE resolver_db;
CREATE DATABASE character_db;

-- Create dedicated users for each service
CREATE USER resolver_service WITH ENCRYPTED PASSWORD 'resolver_password';
CREATE USER character_service WITH ENCRYPTED PASSWORD 'character_password';

GRANT ALL PRIVILEGES ON DATABASE character_db TO character_service;
