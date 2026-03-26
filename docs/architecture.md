# Portal Architecture

## Runtime topology

- Keycloak: identity source, role source, client source
- portal-api: Go BFF
- portal-web: Vue SPA
- MongoDB: projection store, settings store, session store

## Login flow

1. Browser opens `/portal`.
2. `portal-web` calls `GET /api/auth/me`.
3. If unauthenticated, `portal-web` redirects browser to `GET /api/auth/login`.
4. `portal-api` redirects to Keycloak OIDC authorization endpoint.
5. Keycloak authenticates the user and redirects to `GET /api/auth/callback?code=...`.
6. `portal-api` exchanges the OIDC code for tokens.
7. `portal-api` uses `portal-sync` client credentials to call Keycloak Admin API.
8. `portal-api` synchronizes:
   - current realm
   - current realm clients
   - current user profile
   - current user effective realm roles
   - current user effective client roles
9. `portal-api` upserts Mongo projections.
10. `portal-api` creates a portal session.
11. Browser returns to `/portal`.

## Session rules

- Portal session idle timeout is controlled by `portal_settings`.
- Each protected portal-api request refreshes `lastActiveAt`.
- Portal logout deletes portal session first, then redirects to Keycloak logout.
- MongoDB is never used as an authentication source.
