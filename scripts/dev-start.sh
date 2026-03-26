#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "${ROOT_DIR}"

if [[ ! -f ".env" ]]; then
  cp .env.example .env
fi

docker compose up --build -d mongo keycloak
./scripts/keycloak-bootstrap.sh
./scripts/init-mongo-indexes.sh
docker compose up --build -d portal-api portal-web
