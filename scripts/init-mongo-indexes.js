db = db.getSiblingDB(process.env.MONGO_DB || 'portal');

function ensureCollection(name) {
  const exists = db.getCollectionNames().includes(name);
  if (!exists) {
    db.createCollection(name);
  }
}

[
  'kc_realms',
  'kc_clients',
  'portal_client_meta',
  'kc_users',
  'portal_sessions',
  'portal_settings',
].forEach(ensureCollection);

db.kc_realms.createIndex({ realmId: 1 }, { unique: true, name: 'ux_realm_id' });
db.kc_clients.createIndex({ realmId: 1, clientId: 1 }, { unique: true, name: 'ux_realm_client_id' });
db.kc_clients.createIndex({ realmId: 1, clientUuid: 1 }, { unique: true, name: 'ux_realm_client_uuid' });
db.portal_client_meta.createIndex({ realmId: 1, clientId: 1 }, { unique: true, name: 'ux_realm_client_meta' });
db.kc_users.createIndex({ realmId: 1, userId: 1 }, { unique: true, name: 'ux_realm_user_id' });
db.kc_users.createIndex({ realmId: 1, username: 1 }, { unique: true, name: 'ux_realm_username' });
db.portal_sessions.createIndex({ sessionId: 1 }, { unique: true, name: 'ux_session_id' });
db.portal_sessions.createIndex({ expiresAt: 1 }, { expireAfterSeconds: 0, name: 'ttl_expires_at' });
db.portal_settings.createIndex({ _id: 1 }, { unique: true, name: 'ux_settings_id' });

db.portal_settings.updateOne(
  { _id: 'global' },
  {
    $setOnInsert: {
      _id: 'global',
      idleTimeoutMinutes: 15,
      idleWarnSeconds: 60,
      updatedAt: new Date(),
    },
  },
  { upsert: true },
);
