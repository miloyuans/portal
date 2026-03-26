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
MONGO_DB="${MONGO_DB}" mongosh "${MONGO_URI}/${MONGO_DB}" "${ROOT_DIR}/scripts/init-mongo-indexes.js"
