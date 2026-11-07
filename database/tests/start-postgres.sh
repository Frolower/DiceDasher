#!/usr/bin/env bash
# Shared by the database and backend integration CI jobs; run from the repo root.
set -euo pipefail

: "${POSTGRES_CONTAINER:?Set POSTGRES_CONTAINER to an isolated test container name}"

docker build -t dicedasher-postgres:ci database
docker run -d --name "$POSTGRES_CONTAINER" \
  -e POSTGRES_PASSWORD=postgres \
  -p 127.0.0.1:5432:5432 \
  --health-cmd 'pg_isready -h 127.0.0.1 -U postgres' \
  --health-interval 1s --health-timeout 5s --health-retries 30 \
  dicedasher-postgres:ci \
  postgres \
  -c shared_preload_libraries=pg_cron \
  -c cron.database_name=user_db \
  -c cron.timezone=UTC \
  -c cron.use_background_workers=on

for attempt in {1..60}; do
  if [ "$(docker inspect --format '{{.State.Health.Status}}' "$POSTGRES_CONTAINER")" = healthy ]; then
    exit 0
  fi
  sleep 1
done

docker logs "$POSTGRES_CONTAINER"
echo 'PostgreSQL did not become healthy in time' >&2
exit 1
