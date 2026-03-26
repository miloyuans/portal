#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ -f "${ROOT_DIR}/.env" ]]; then
  set -a
  source "${ROOT_DIR}/.env"
  set +a
fi

: "${MONGO_URI:=mongodb://localhost:27017}"
: "${MONGO_DB:=portal}"

echo "Initializing Mongo indexes on ${MONGO_URI}/${MONGO_DB}"

if command -v mongosh >/dev/null 2>&1; then
  MONGO_DB="${MONGO_DB}" mongosh "${MONGO_URI}/${MONGO_DB}" "${ROOT_DIR}/scripts/init-mongo-indexes.js"
  exit 0
fi

echo "Local mongosh not found, falling back to the mongo container"

docker compose exec -e MONGO_DB="${MONGO_DB}" -T mongo sh -lc 'mongosh "mongodb://localhost:27017/${MONGO_DB}"' < "${ROOT_DIR}/scripts/init-mongo-indexes.js"
