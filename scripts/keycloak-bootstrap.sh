#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ -f "${ROOT_DIR}/.env" ]]; then
  set -a
  source "${ROOT_DIR}/.env"
  set +a
fi

: "${KEYCLOAK_REALM:=portal}"
: "${KEYCLOAK_BOOTSTRAP_ADMIN_USER:=admin}"
: "${KEYCLOAK_BOOTSTRAP_ADMIN_PASSWORD:=admin123}"
: "${KEYCLOAK_OIDC_CLIENT_ID:=portal-api}"
: "${KEYCLOAK_OIDC_CLIENT_SECRET:=portal-api-secret}"
: "${KEYCLOAK_ADMIN_CLIENT_ID:=portal-sync}"
: "${KEYCLOAK_ADMIN_CLIENT_SECRET:=portal-sync-secret}"
: "${KEYCLOAK_REDIRECT_URL:=http://localhost:8080/api/auth/callback}"

require() {
  command -v "$1" >/dev/null 2>&1 || { echo "missing dependency: $1" >&2; exit 1; }
}

require docker
require jq

kcadm() {
  docker compose exec -T keycloak /opt/keycloak/bin/kcadm.sh "$@"
}

wait_for_keycloak() {
  local attempt=0
  until kcadm config credentials --server http://localhost:8080 --realm master --user "${KEYCLOAK_BOOTSTRAP_ADMIN_USER}" --password "${KEYCLOAK_BOOTSTRAP_ADMIN_PASSWORD}" >/dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [[ ${attempt} -gt 30 ]]; then
      echo "keycloak did not become ready in time" >&2
      exit 1
    fi
    echo "waiting for keycloak..."
    sleep 5
  done
}

client_uuid() {
  local realm="$1"
  local client_id="$2"
  kcadm get clients -r "${realm}" -q clientId="${client_id}" | jq -r '.[0].id // empty'
}

ensure_realm() {
  if ! kcadm get "realms/${KEYCLOAK_REALM}" >/dev/null 2>&1; then
    kcadm create realms -s realm="${KEYCLOAK_REALM}" -s enabled=true -s displayName="Portal Realm" >/dev/null
  fi
}

ensure_realm_role() {
  local role_name="$1"
  if ! kcadm get "roles/${role_name}" -r "${KEYCLOAK_REALM}" >/dev/null 2>&1; then
    kcadm create roles -r "${KEYCLOAK_REALM}" -s name="${role_name}" >/dev/null
  fi
}

ensure_oidc_client() {
  if [[ -z "$(client_uuid "${KEYCLOAK_REALM}" "${KEYCLOAK_OIDC_CLIENT_ID}")" ]]; then
    kcadm create clients -r "${KEYCLOAK_REALM}" \
      -s clientId="${KEYCLOAK_OIDC_CLIENT_ID}" \
      -s name="Portal API" \
      -s protocol="openid-connect" \
      -s publicClient=false \
      -s secret="${KEYCLOAK_OIDC_CLIENT_SECRET}" \
      -s standardFlowEnabled=true \
      -s serviceAccountsEnabled=false \
      -s "redirectUris=[\"${KEYCLOAK_REDIRECT_URL}\"]" >/dev/null
  fi
}

ensure_sync_client() {
  if [[ -z "$(client_uuid "${KEYCLOAK_REALM}" "${KEYCLOAK_ADMIN_CLIENT_ID}")" ]]; then
    kcadm create clients -r "${KEYCLOAK_REALM}" \
      -s clientId="${KEYCLOAK_ADMIN_CLIENT_ID}" \
      -s name="Portal Sync" \
      -s protocol="openid-connect" \
      -s publicClient=false \
      -s secret="${KEYCLOAK_ADMIN_CLIENT_SECRET}" \
      -s standardFlowEnabled=false \
      -s serviceAccountsEnabled=true >/dev/null
  fi
}

ensure_sample_client() {
  local client_id="$1"
  local client_name="$2"
  local client_url="$3"
  if [[ -z "$(client_uuid "${KEYCLOAK_REALM}" "${client_id}")" ]]; then
    kcadm create clients -r "${KEYCLOAK_REALM}" \
      -s clientId="${client_id}" \
      -s name="${client_name}" \
      -s protocol="openid-connect" \
      -s publicClient=true \
      -s enabled=true \
      -s baseUrl="${client_url}" \
      -s rootUrl="${client_url}" >/dev/null
  fi
  if ! kcadm get "clients/$(client_uuid "${KEYCLOAK_REALM}" "${client_id}")/roles/viewer" -r "${KEYCLOAK_REALM}" >/dev/null 2>&1; then
    kcadm create "clients/$(client_uuid "${KEYCLOAK_REALM}" "${client_id}")/roles" -r "${KEYCLOAK_REALM}" -s name="viewer" >/dev/null
  fi
}

ensure_user() {
  local username="$1"
  local password="$2"
  local first_name="$3"
  local last_name="$4"
  local email="$5"

  local user_id
  user_id="$(kcadm get users -r "${KEYCLOAK_REALM}" -q username="${username}" | jq -r '.[0].id // empty')"
  if [[ -z "${user_id}" ]]; then
    kcadm create users -r "${KEYCLOAK_REALM}" \
      -s username="${username}" \
      -s enabled=true \
      -s firstName="${first_name}" \
      -s lastName="${last_name}" \
      -s email="${email}" \
      -s emailVerified=true >/dev/null
    user_id="$(kcadm get users -r "${KEYCLOAK_REALM}" -q username="${username}" | jq -r '.[0].id // empty')"
  fi
  kcadm set-password -r "${KEYCLOAK_REALM}" --userid "${user_id}" --new-password "${password}" >/dev/null
}

grant_realm_role() {
  local username="$1"
  local role_name="$2"
  kcadm add-roles -r "${KEYCLOAK_REALM}" --uusername "${username}" --rolename "${role_name}" >/dev/null || true
}

grant_client_role() {
  local username="$1"
  local client_id="$2"
  local role_name="$3"
  kcadm add-roles -r "${KEYCLOAK_REALM}" --uusername "${username}" --cclientid "${client_id}" --rolename "${role_name}" >/dev/null || true
}

grant_sync_service_account_roles() {
  local sync_uuid realm_mgmt_uuid svc_user_id
  sync_uuid="$(client_uuid "${KEYCLOAK_REALM}" "${KEYCLOAK_ADMIN_CLIENT_ID}")"
  realm_mgmt_uuid="$(client_uuid "${KEYCLOAK_REALM}" "realm-management")"
  svc_user_id="$(kcadm get "clients/${sync_uuid}/service-account-user" -r "${KEYCLOAK_REALM}" | jq -r '.id')"

  echo "portal-sync client secret: ${KEYCLOAK_ADMIN_CLIENT_SECRET}"
  echo "portal-sync service account user id: ${svc_user_id}"

  for role_name in query-users view-users query-clients view-clients query-realms view-realm; do
    kcadm add-roles -r "${KEYCLOAK_REALM}" --uid "${svc_user_id}" --cclientid realm-management --rolename "${role_name}" >/dev/null || true
  done
  echo "realm-management client uuid: ${realm_mgmt_uuid}"
}

echo "Bootstrapping Keycloak for portal"
docker compose up -d keycloak >/dev/null
wait_for_keycloak
ensure_realm
ensure_realm_role portal_user
ensure_realm_role portal_admin
ensure_oidc_client
ensure_sync_client
ensure_sample_client sales-app "Sales Console" "https://sales.example.local"
ensure_sample_client ops-app "Operations Console" "https://ops.example.local"
ensure_user portal-admin Admin123! Portal Admin portal-admin@example.local
ensure_user alice Alice123! Alice Operator alice@example.local
grant_realm_role portal-admin portal_admin
grant_realm_role portal-admin portal_user
grant_realm_role alice portal_user
grant_client_role portal-admin sales-app viewer
grant_client_role portal-admin ops-app viewer
grant_client_role alice sales-app viewer
grant_sync_service_account_roles
echo "Keycloak bootstrap complete"
