#!/bin/bash
set -Eeuo pipefail
umask 077
install -d -m 700 /opt/dicedasher/backups
backup=$(mktemp /opt/dicedasher/backups/.backup.XXXXXX)
trap 'rm -f "$backup"' EXIT
docker exec dicedasher-postgres-1 pg_dumpall -U postgres | gzip > "$backup"
gzip -t "$backup"
mv "$backup" "/opt/dicedasher/backups/$(date -u +%Y%m%dT%H%M%SZ).sql.gz"
find /opt/dicedasher/backups -name '*.sql.gz' -mtime +14 -delete
