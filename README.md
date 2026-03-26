# portal

`portal` is an independently deployed web portal that sits in front of Keycloak.

Core boundaries:

- Keycloak is the only identity source, role source, and client source.
- `portal-api` handles OIDC callback, portal session, Keycloak Admin API access, login-triggered sync, and permission resolution.
- `portal-web` only calls `portal-api`.
- MongoDB is a projection store, settings store, and session store. It is not an authentication source.

## Stack

- Go 1.23+
- Gin
- MongoDB official Go driver
- Vue 3 + TypeScript + Vite + Pinia + Vue Router + Element Plus
- Docker Compose
- Kubernetes manifests

## Layout

```text
portal/
  apps/
    portal-api/
    portal-web/
  internal/
    auth/
    config/
    handler/
    kcadmin/
    middleware/
    model/
    permission/
    repository/
    service/
    session/
    sync/
  deployments/
    docker-compose/
    k8s/
  docs/
  scripts/
  tests/
```

## Prerequisites

- Linux shell
- Docker and Docker Compose plugin
- `make`
- `jq`
- `mongosh`
- Go 1.23+ for local backend development
- Node 22+ and npm for local frontend development

## Environment

```bash
cp .env.example .env
```

Edit `.env` if you need different hostnames, ports, or secrets.

## Start the full development stack

```bash
make dev-start
```

This does the following:

1. Starts MongoDB and Keycloak
2. Bootstraps the Keycloak realm, clients, roles, and sample users
3. Ensures Mongo collections and indexes
4. Starts `portal-api` and `portal-web`

## Default URLs

- portal-web: `http://localhost:5173`
- portal-api: `http://localhost:8080`
- Keycloak: `http://localhost:8081`
- OpenAPI: `http://localhost:8080/openapi.yaml`

## Default Keycloak bootstrap data

- realm: `portal`
- OIDC client: `portal-api`
- admin service-account client: `portal-sync`
- realm roles: `portal_user`, `portal_admin`
- admin user: `portal-admin / Admin123!`
- normal user: `alice / Alice123!`

## Scripts

- `scripts/dev-start.sh`: start the local Docker-based development environment
- `scripts/keycloak-bootstrap.sh`: idempotently create realm, clients, roles, service-account bindings, and sample users
- `scripts/init-mongo-indexes.sh`: create required collections and indexes in MongoDB

## API overview

Authentication:

- `GET /api/auth/login`
- `GET /api/auth/callback`
- `POST /api/auth/logout`
- `GET /api/auth/me`

Portal:

- `GET /api/portal/apps`
- `GET /api/portal/realms`
- `GET /api/portal/profile`

Admin:

- `GET /api/admin/realms`
- `GET /api/admin/clients`
- `PUT /api/admin/clients/:clientId/meta`
- `GET /api/admin/users/:userId`
- `GET /api/admin/settings/session`
- `PUT /api/admin/settings/session`
- `GET /api/admin/sync-status`

## Local backend run

```bash
go run ./apps/portal-api
```

## Local frontend run

```bash
cd apps/portal-web
npm install
npm run dev
```

## Tests

Unit tests:

```bash
go test ./...
```

Integration test skeletons:

```bash
go test -tags integration ./tests/integration/...
```

## Notes

- The portal idle timeout defaults to 15 minutes and is enforced by `portal-api`, not by changing Keycloak global session timeout.
- On successful login, `portal-api` synchronizes the current realm, realm clients, current user profile, effective realm roles, and effective client roles before creating the portal session.
- On logout, the portal deletes its own session first and then redirects the browser to Keycloak logout.
