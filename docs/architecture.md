# Portal Architecture

## Core rules

- `portal-web` and `portal-api` are independently deployed.
- Keycloak is the only identity and permission source of truth.
- MongoDB is used only as the portal projection store.
- Frontend only talks to `portal-api`.

## Login flow

1. User enters `portal-web`.
2. `portal-web` redirects browser to `portal-api /api/v1/auth/login`.
3. `portal-api` redirects to Keycloak OIDC authorization endpoint.
4. Keycloak sends the browser back to `portal-api /api/v1/auth/callback`.
5. `portal-api` exchanges `code` for tokens.
6. `portal-api` calls Keycloak Admin REST API and synchronizes:
   - current realm base info
   - current realm clients
   - current user profile
   - current user realm roles
   - current user client roles
7. Projection data is upserted into MongoDB.
8. Only after sync completes does `portal-api` create `portal_sessions`.
9. Browser is redirected back to `portal-web`.

## Authorization and session

- Session idle timeout defaults to 15 minutes and is controlled by portal middleware plus the web idle timer.
- Logout always deletes the portal session first, then redirects to Keycloak logout.
- Admin page access uses synced realm roles from the current portal session.
