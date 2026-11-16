#!/bin/bash
set -Eeuo pipefail
sha=${1:-}
[[ "$sha" =~ ^[0-9a-f]{40}$ ]] || exit 2
exec 9>/var/lock/dicedasher-deploy.lock
flock -x 9
cd /opt/dicedasher/source
git fetch -q origin master
[[ $(git rev-parse origin/master) == "$sha" ]] || { echo 'Commit is not the current master'; exit 1; }
git checkout -q --detach "$sha"
export IMAGE_ROOT=ghcr.io/frolower/dicedasher IMAGE_TAG="$sha"
export DOCKER_CONFIG
DOCKER_CONFIG=$(mktemp -d)
trap 'rm -rf "$DOCKER_CONFIG"' EXIT
IFS= read -r actor
IFS= read -r token
printf '%s' "$token" | docker login ghcr.io -u "$actor" --password-stdin >/dev/null
unset token
compose=(docker compose --project-name dicedasher --env-file /opt/dicedasher/secrets.env -f deploy/docker-compose.yml)
"${compose[@]}" pull
"${compose[@]}" up -d --remove-orphans --wait --wait-timeout 180
for path in /health /api/health; do
  curl --fail --silent --show-error --retry 20 --retry-all-errors --retry-delay 3 --resolve dicedasher.dmitriy.space:443:127.0.0.1 "https://dicedasher.dmitriy.space$path" >/dev/null
done
printf '%s\n' "$sha" > /opt/dicedasher/deployed-sha
