db = db.getSiblingDB('portal');

db.createCollection('kc_realms');
db.createCollection('kc_clients');
db.createCollection('portal_client_meta');
db.createCollection('kc_users');
db.createCollection('portal_sessions');
db.createCollection('portal_settings');

db.kc_realms.createIndex({ realm: 1 }, { unique: true, name: 'ux_realm' });
db.kc_clients.createIndex({ realm: 1, clientId: 1 }, { unique: true, name: 'ux_realm_client_id' });
db.kc_clients.createIndex({ realm: 1, clientUuid: 1 }, { unique: true, name: 'ux_realm_client_uuid' });
db.portal_client_meta.createIndex({ realm: 1, clientId: 1 }, { unique: true, name: 'ux_realm_client_meta' });
db.kc_users.createIndex({ realm: 1, userId: 1 }, { unique: true, name: 'ux_realm_user' });
db.portal_sessions.createIndex({ sessionId: 1 }, { unique: true, name: 'ux_session_id' });
db.portal_sessions.createIndex({ expiresAt: 1 }, { expireAfterSeconds: 0, name: 'ttl_expires_at' });
db.portal_settings.createIndex({ realm: 1 }, { unique: true, name: 'ux_settings_realm' });
