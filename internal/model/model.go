package model

import "time"

type RealmProjection struct {
	Realm               string    `json:"realm" bson:"realm"`
	DisplayName         string    `json:"displayName,omitempty" bson:"displayName,omitempty"`
	DisplayNameHTML     string    `json:"displayNameHtml,omitempty" bson:"displayNameHtml,omitempty"`
	LoginTheme          string    `json:"loginTheme,omitempty" bson:"loginTheme,omitempty"`
	AccountTheme        string    `json:"accountTheme,omitempty" bson:"accountTheme,omitempty"`
	AdminTheme          string    `json:"adminTheme,omitempty" bson:"adminTheme,omitempty"`
	RegistrationAllowed bool      `json:"registrationAllowed" bson:"registrationAllowed"`
	SSLRequired         string    `json:"sslRequired,omitempty" bson:"sslRequired,omitempty"`
	SyncedAt            time.Time `json:"syncedAt" bson:"syncedAt"`
	UpdatedAt           time.Time `json:"updatedAt" bson:"updatedAt"`
}

type ClientProjection struct {
	Realm        string            `json:"realm" bson:"realm"`
	ClientUUID   string            `json:"clientUuid" bson:"clientUuid"`
	ClientID     string            `json:"clientId" bson:"clientId"`
	Name         string            `json:"name,omitempty" bson:"name,omitempty"`
	Description  string            `json:"description,omitempty" bson:"description,omitempty"`
	RootURL      string            `json:"rootUrl,omitempty" bson:"rootUrl,omitempty"`
	BaseURL      string            `json:"baseUrl,omitempty" bson:"baseUrl,omitempty"`
	RedirectURIs []string          `json:"redirectUris,omitempty" bson:"redirectUris,omitempty"`
	WebOrigins   []string          `json:"webOrigins,omitempty" bson:"webOrigins,omitempty"`
	Enabled      bool              `json:"enabled" bson:"enabled"`
	PublicClient bool              `json:"publicClient" bson:"publicClient"`
	Protocol     string            `json:"protocol,omitempty" bson:"protocol,omitempty"`
	Attributes   map[string]string `json:"attributes,omitempty" bson:"attributes,omitempty"`
	SyncedAt     time.Time         `json:"syncedAt" bson:"syncedAt"`
	UpdatedAt    time.Time         `json:"updatedAt" bson:"updatedAt"`
}

type PortalClientMeta struct {
	Realm               string    `json:"realm" bson:"realm"`
	ClientID            string    `json:"clientId" bson:"clientId"`
	DisplayName         string    `json:"displayName" bson:"displayName"`
	Description         string    `json:"description,omitempty" bson:"description,omitempty"`
	TargetURL           string    `json:"targetUrl,omitempty" bson:"targetUrl,omitempty"`
	Icon                string    `json:"icon,omitempty" bson:"icon,omitempty"`
	Category            string    `json:"category,omitempty" bson:"category,omitempty"`
	SortOrder           int       `json:"sortOrder" bson:"sortOrder"`
	Enabled             bool      `json:"enabled" bson:"enabled"`
	ShowInPortal        bool      `json:"showInPortal" bson:"showInPortal"`
	RequiredRealmRoles  []string  `json:"requiredRealmRoles,omitempty" bson:"requiredRealmRoles,omitempty"`
	RequiredClientRoles []string  `json:"requiredClientRoles,omitempty" bson:"requiredClientRoles,omitempty"`
	Tags                []string  `json:"tags,omitempty" bson:"tags,omitempty"`
	UpdatedBy           string    `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	CreatedAt           time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UserProjection struct {
	Realm         string              `json:"realm" bson:"realm"`
	UserID        string              `json:"userId" bson:"userId"`
	Username      string              `json:"username" bson:"username"`
	Email         string              `json:"email,omitempty" bson:"email,omitempty"`
	FirstName     string              `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName      string              `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Enabled       bool                `json:"enabled" bson:"enabled"`
	EmailVerified bool                `json:"emailVerified" bson:"emailVerified"`
	Attributes    map[string][]string `json:"attributes,omitempty" bson:"attributes,omitempty"`
	RealmRoles    []string            `json:"realmRoles,omitempty" bson:"realmRoles,omitempty"`
	ClientRoles   map[string][]string `json:"clientRoles,omitempty" bson:"clientRoles,omitempty"`
	SyncedAt      time.Time           `json:"syncedAt" bson:"syncedAt"`
	UpdatedAt     time.Time           `json:"updatedAt" bson:"updatedAt"`
}

type PortalSession struct {
	SessionID          string              `json:"sessionId" bson:"sessionId"`
	Realm              string              `json:"realm" bson:"realm"`
	UserID             string              `json:"userId" bson:"userId"`
	Username           string              `json:"username" bson:"username"`
	Email              string              `json:"email,omitempty" bson:"email,omitempty"`
	DisplayName        string              `json:"displayName,omitempty" bson:"displayName,omitempty"`
	RealmRoles         []string            `json:"realmRoles,omitempty" bson:"realmRoles,omitempty"`
	ClientRoles        map[string][]string `json:"clientRoles,omitempty" bson:"clientRoles,omitempty"`
	AccessToken        string              `json:"-" bson:"accessToken"`
	RefreshToken       string              `json:"-" bson:"refreshToken"`
	IDToken            string              `json:"-" bson:"idToken"`
	IdleTimeoutMinutes int                 `json:"idleTimeoutMinutes" bson:"idleTimeoutMinutes"`
	LastSeenAt         time.Time           `json:"lastSeenAt" bson:"lastSeenAt"`
	ExpiresAt          time.Time           `json:"expiresAt" bson:"expiresAt"`
	CreatedAt          time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time           `json:"updatedAt" bson:"updatedAt"`
}

type PortalSettings struct {
	Realm              string    `json:"realm" bson:"realm"`
	IdleTimeoutMinutes int       `json:"idleTimeoutMinutes" bson:"idleTimeoutMinutes"`
	CreatedAt          time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" bson:"updatedAt"`
}

type SessionView struct {
	SessionID          string              `json:"sessionId"`
	Realm              string              `json:"realm"`
	UserID             string              `json:"userId"`
	Username           string              `json:"username"`
	Email              string              `json:"email,omitempty"`
	DisplayName        string              `json:"displayName,omitempty"`
	RealmRoles         []string            `json:"realmRoles,omitempty"`
	ClientRoles        map[string][]string `json:"clientRoles,omitempty"`
	IdleTimeoutMinutes int                 `json:"idleTimeoutMinutes"`
	LastSeenAt         time.Time           `json:"lastSeenAt"`
	ExpiresAt          time.Time           `json:"expiresAt"`
}

func (s PortalSession) View() SessionView {
	return SessionView{
		SessionID:          s.SessionID,
		Realm:              s.Realm,
		UserID:             s.UserID,
		Username:           s.Username,
		Email:              s.Email,
		DisplayName:        s.DisplayName,
		RealmRoles:         s.RealmRoles,
		ClientRoles:        s.ClientRoles,
		IdleTimeoutMinutes: s.IdleTimeoutMinutes,
		LastSeenAt:         s.LastSeenAt,
		ExpiresAt:          s.ExpiresAt,
	}
}

type PortalApp struct {
	ClientID    string   `json:"clientId"`
	DisplayName string   `json:"displayName"`
	Description string   `json:"description,omitempty"`
	TargetURL   string   `json:"targetUrl"`
	Icon        string   `json:"icon,omitempty"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	SortOrder   int      `json:"sortOrder"`
}

type CurrentUserProfile struct {
	Realm    string         `json:"realm"`
	User     SessionView    `json:"user"`
	IsAdmin  bool           `json:"isAdmin"`
	Settings PortalSettings `json:"settings"`
}

type IDTokenClaims struct {
	Subject           string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
}
