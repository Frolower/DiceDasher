#!/bin/sh
set -eu
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres \
  -v user_password="$USER_PASSWORD" -v character_password="$CHARACTER_PASSWORD" -v resolve_password="$RESOLVE_PASSWORD" <<'SQL'
CREATE DATABASE resolve_db;
CREATE DATABASE character_db;
CREATE DATABASE user_db;
CREATE USER resolve_service WITH ENCRYPTED PASSWORD :'resolve_password';
CREATE USER character_service WITH ENCRYPTED PASSWORD :'character_password';
CREATE USER user_service WITH ENCRYPTED PASSWORD :'user_password';
SQL
