package model

import "time"

const (
	// LaunchModeDirect opens the configured launch URL directly.
	LaunchModeDirect = "direct"
	// LaunchModeSPInitiated opens the application URL and lets the service provider start SSO.
	LaunchModeSPInitiated = "sp_initiated"
	// LaunchModeDisabled keeps the app visible but not launchable.
	LaunchModeDisabled = "disabled"
)

// NormalizeLaunchMode returns a safe runtime launch mode.
func NormalizeLaunchMode(mode string) string {
	switch mode {
	case LaunchModeDirect, LaunchModeSPInitiated, LaunchModeDisabled:
		return mode
	default:
		return LaunchModeSPInitiated
	}
}

// RealmProjection stores the projected Keycloak realm snapshot.
type RealmProjection struct {
	RealmID     string                 `json:"realmId" bson:"realmId"`
	RealmName   string                 `json:"realmName" bson:"realmName"`
	DisplayName string                 `json:"displayName,omitempty" bson:"displayName,omitempty"`
	Enabled     bool                   `json:"enabled" bson:"enabled"`
	Attributes  map[string]any         `json:"attributes,omitempty" bson:"attributes,omitempty"`
	SyncedAt    time.Time              `json:"syncedAt" bson:"syncedAt"`
}

// ClientProjection stores the projected Keycloak client snapshot.
type ClientProjection struct {
	RealmID    string                 `json:"realmId" bson:"realmId"`
	ClientUUID string                 `json:"clientUuid" bson:"clientUuid"`
	ClientID   string                 `json:"clientId" bson:"clientId"`
	Name       string                 `json:"name,omitempty" bson:"name,omitempty"`
	Enabled    bool                   `json:"enabled" bson:"enabled"`
	BaseURL    string                 `json:"baseUrl,omitempty" bson:"baseUrl,omitempty"`
	RootURL    string                 `json:"rootUrl,omitempty" bson:"rootUrl,omitempty"`
	Protocol   string                 `json:"protocol,omitempty" bson:"protocol,omitempty"`
	Attributes map[string]string      `json:"attributes,omitempty" bson:"attributes,omitempty"`
	SyncedAt   time.Time              `json:"syncedAt" bson:"syncedAt"`
}

// AccessRules defines how a portal client is exposed.
type AccessRules struct {
	AnyRealmRoles   []string `json:"anyRealmRoles,omitempty" bson:"anyRealmRoles,omitempty"`
	AnyClientRoles  []string `json:"anyClientRoles,omitempty" bson:"anyClientRoles,omitempty"`
	AdminRealmRoles []string `json:"adminRealmRoles,omitempty" bson:"adminRealmRoles,omitempty"`
}

// PortalClientMeta stores portal-only metadata per Keycloak client.
type PortalClientMeta struct {
	RealmID      string            `json:"realmId" bson:"realmId"`
	ClientID     string            `json:"clientId" bson:"clientId"`
	DisplayName  string            `json:"displayName" bson:"displayName"`
	Icon         string            `json:"icon,omitempty" bson:"icon,omitempty"`
	Category     string            `json:"category,omitempty" bson:"category,omitempty"`
	Sort         int               `json:"sort" bson:"sort"`
	LaunchMode   string            `json:"launchMode,omitempty" bson:"launchMode,omitempty"`
	LaunchURL    string            `json:"launchUrl,omitempty" bson:"launchUrl,omitempty"`
	LaunchConfig map[string]string `json:"launchConfig,omitempty" bson:"launchConfig,omitempty"`
	Visible      bool              `json:"visible" bson:"visible"`
	AccessRules  AccessRules       `json:"accessRules,omitempty" bson:"accessRules,omitempty"`
}

// UserProjection stores the current user projection.
type UserProjection struct {
	RealmID     string              `json:"realmId" bson:"realmId"`
	UserID      string              `json:"userId" bson:"userId"`
	Username    string              `json:"username" bson:"username"`
	Email       string              `json:"email,omitempty" bson:"email,omitempty"`
	Enabled     bool                `json:"enabled" bson:"enabled"`
	FirstName   string              `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName    string              `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Attributes  map[string][]string `json:"attributes,omitempty" bson:"attributes,omitempty"`
	RealmRoles  []string            `json:"realmRoles,omitempty" bson:"realmRoles,omitempty"`
	ClientRoles map[string][]string `json:"clientRoles,omitempty" bson:"clientRoles,omitempty"`
	SyncedAt    time.Time           `json:"syncedAt" bson:"syncedAt"`
}

// PortalSession stores a signed portal session persisted in MongoDB.
type PortalSession struct {
	SessionID          string              `json:"sessionId" bson:"sessionId"`
	RealmID            string              `json:"realmId" bson:"realmId"`
	UserID             string              `json:"userId" bson:"userId"`
	Username           string              `json:"username" bson:"username"`
	DisplayName        string              `json:"displayName,omitempty" bson:"displayName,omitempty"`
	RealmRoles         []string            `json:"realmRoles,omitempty" bson:"realmRoles,omitempty"`
	ClientRoles        map[string][]string `json:"clientRoles,omitempty" bson:"clientRoles,omitempty"`
	IdleTimeoutMinutes int                `json:"idleTimeoutMinutes" bson:"idleTimeoutMinutes"`
	LastActiveAt       time.Time           `json:"lastActiveAt" bson:"lastActiveAt"`
	ExpiresAt          time.Time           `json:"expiresAt" bson:"expiresAt"`
	AbsoluteExpiresAt  time.Time           `json:"absoluteExpiresAt" bson:"absoluteExpiresAt"`
	CreatedAt          time.Time           `json:"createdAt" bson:"createdAt"`
	IDToken            string              `json:"-" bson:"idToken"`
}

// PortalSettings stores global portal session settings.
type PortalSettings struct {
	ID                 string    `json:"id" bson:"_id"`
	IdleTimeoutMinutes int       `json:"idleTimeoutMinutes" bson:"idleTimeoutMinutes"`
	IdleWarnSeconds    int       `json:"idleWarnSeconds" bson:"idleWarnSeconds"`
	UpdatedAt          time.Time `json:"updatedAt" bson:"updatedAt"`
}

// SessionView is the frontend-safe session view model.
type SessionView struct {
	SessionID          string              `json:"sessionId"`
	RealmID            string              `json:"realmId"`
	UserID             string              `json:"userId"`
	Username           string              `json:"username"`
	DisplayName        string              `json:"displayName,omitempty"`
	RealmRoles         []string            `json:"realmRoles,omitempty"`
	ClientRoles        map[string][]string `json:"clientRoles,omitempty"`
	IdleTimeoutMinutes int                `json:"idleTimeoutMinutes"`
	LastActiveAt       time.Time           `json:"lastActiveAt"`
	ExpiresAt          time.Time           `json:"expiresAt"`
	AbsoluteExpiresAt  time.Time           `json:"absoluteExpiresAt"`
}

// View converts a PortalSession into a safe response view.
func (s PortalSession) View() SessionView {
	return SessionView{
		SessionID:          s.SessionID,
		RealmID:            s.RealmID,
		UserID:             s.UserID,
		Username:           s.Username,
		DisplayName:        s.DisplayName,
		RealmRoles:         s.RealmRoles,
		ClientRoles:        s.ClientRoles,
		IdleTimeoutMinutes: s.IdleTimeoutMinutes,
		LastActiveAt:       s.LastActiveAt,
		ExpiresAt:          s.ExpiresAt,
		AbsoluteExpiresAt:  s.AbsoluteExpiresAt,
	}
}

// PortalAppView is the portal navigation result for the current user.
type PortalAppView struct {
	ClientID    string `json:"clientId"`
	DisplayName string `json:"displayName"`
	Category    string `json:"category,omitempty"`
	Icon        string `json:"icon,omitempty"`
	LaunchMode  string `json:"launchMode"`
	LaunchURL   string `json:"launchUrl,omitempty"`
	CanView     bool   `json:"canView"`
	CanLaunch   bool   `json:"canLaunch"`
	CanAdmin    bool   `json:"canAdmin"`
}

// PortalLaunchView is the resolved launch result returned by portal-api.
type PortalLaunchView struct {
	ClientID    string `json:"clientId"`
	DisplayName string `json:"displayName"`
	LaunchMode  string `json:"launchMode"`
	LaunchURL   string `json:"launchUrl"`
}

// CurrentUserProfile is the frontend profile payload.
type CurrentUserProfile struct {
	Session  SessionView      `json:"session"`
	User     UserProjection   `json:"user"`
	Realm    RealmProjection  `json:"realm"`
	Settings PortalSettings   `json:"settings"`
}

// IDTokenClaims stores the minimum OIDC claims needed by portal.
type IDTokenClaims struct {
	Subject           string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
}

// SyncStatus summarizes the latest login-triggered sync state.
type SyncStatus struct {
	RealmID           string    `json:"realmId"`
	RealmSyncedAt     time.Time `json:"realmSyncedAt"`
	UserSyncedAt      time.Time `json:"userSyncedAt"`
	ClientCount       int       `json:"clientCount"`
	SettingsUpdatedAt time.Time `json:"settingsUpdatedAt"`
}
