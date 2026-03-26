export interface ApiEnvelope<T> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
    details?: unknown
  }
}

export interface PortalSettings {
  realm: string
  idleTimeoutMinutes: number
}

export interface SessionView {
  sessionId: string
  realm: string
  userId: string
  username: string
  email?: string
  displayName?: string
  realmRoles: string[]
  clientRoles: Record<string, string[]>
  idleTimeoutMinutes: number
  lastSeenAt: string
  expiresAt: string
}

export interface CurrentUserProfile {
  realm: string
  user: SessionView
  isAdmin: boolean
  settings: PortalSettings
}

export interface PortalApp {
  clientId: string
  displayName: string
  description?: string
  targetUrl: string
  icon?: string
  category?: string
  tags?: string[]
  sortOrder: number
}

export interface PortalClientMeta {
  realm: string
  clientId: string
  displayName: string
  description?: string
  targetUrl?: string
  icon?: string
  category?: string
  sortOrder: number
  enabled: boolean
  showInPortal: boolean
  requiredRealmRoles?: string[]
  requiredClientRoles?: string[]
  tags?: string[]
  updatedBy?: string
}
