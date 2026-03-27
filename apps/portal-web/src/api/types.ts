export interface ApiEnvelope<T> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
    details?: unknown
  }
}

export interface SessionView {
  sessionId: string
  realmId: string
  userId: string
  username: string
  displayName?: string
  realmRoles: string[]
  clientRoles: Record<string, string[]>
  idleTimeoutMinutes: number
  lastActiveAt: string
  expiresAt: string
  absoluteExpiresAt: string
}

export interface RealmProjection {
  realmId: string
  realmName: string
  displayName?: string
  enabled: boolean
  attributes?: Record<string, unknown>
  syncedAt: string
}

export interface UserProjection {
  realmId: string
  userId: string
  username: string
  email?: string
  enabled: boolean
  firstName?: string
  lastName?: string
  attributes?: Record<string, string[]>
  realmRoles: string[]
  clientRoles: Record<string, string[]>
  syncedAt: string
}

export interface PortalSettings {
  id: string
  idleTimeoutMinutes: number
  idleWarnSeconds: number
  updatedAt?: string
}

export interface CurrentUserProfile {
  session: SessionView
  user: UserProjection
  realm: RealmProjection
  settings: PortalSettings
}

export interface PortalAppView {
  clientId: string
  displayName: string
  category?: string
  icon?: string
  launchMode: string
  launchUrl?: string
  canView: boolean
  canLaunch: boolean
  canAdmin: boolean
}

export interface AccessRules {
  anyRealmRoles?: string[]
  anyClientRoles?: string[]
  adminRealmRoles?: string[]
}

export interface PortalClientMeta {
  realmId: string
  clientId: string
  displayName: string
  icon?: string
  category?: string
  sort: number
  launchMode?: string
  launchUrl?: string
  launchConfig?: Record<string, string>
  visible: boolean
  accessRules?: AccessRules
}

export interface PortalLaunchView {
  clientId: string
  displayName: string
  launchMode: string
  launchUrl: string
}

export interface AdminClientRow {
  client: {
    realmId: string
    clientUuid: string
    clientId: string
    name?: string
    enabled: boolean
    baseUrl?: string
    rootUrl?: string
    protocol?: string
    attributes?: Record<string, string>
    syncedAt: string
  }
  meta?: PortalClientMeta
}

export interface SyncStatus {
  realmId: string
  realmSyncedAt: string
  userSyncedAt: string
  clientCount: number
  settingsUpdatedAt: string
}
