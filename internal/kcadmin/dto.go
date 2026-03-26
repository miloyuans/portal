package kcadmin

// RealmRepresentation is the minimal realm DTO used by portal.
type RealmRepresentation struct {
	ID          string         `json:"id,omitempty"`
	Realm       string         `json:"realm"`
	DisplayName string         `json:"displayName,omitempty"`
	Enabled     bool           `json:"enabled"`
	Attributes  map[string]any `json:"attributes,omitempty"`
}

// ClientRepresentation is the minimal client DTO used by portal.
type ClientRepresentation struct {
	ID          string            `json:"id"`
	ClientID    string            `json:"clientId"`
	Name        string            `json:"name,omitempty"`
	Enabled     bool              `json:"enabled"`
	BaseURL     string            `json:"baseUrl,omitempty"`
	RootURL     string            `json:"rootUrl,omitempty"`
	Protocol    string            `json:"protocol,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

// UserRepresentation is the minimal user DTO used by portal.
type UserRepresentation struct {
	ID         string              `json:"id"`
	Username   string              `json:"username"`
	Email      string              `json:"email,omitempty"`
	Enabled    bool                `json:"enabled"`
	FirstName  string              `json:"firstName,omitempty"`
	LastName   string              `json:"lastName,omitempty"`
	Attributes map[string][]string `json:"attributes,omitempty"`
}

// RoleRepresentation is the minimal role DTO used by portal.
type RoleRepresentation struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ClientRole  bool   `json:"clientRole,omitempty"`
}

// MappingsRepresentation is the minimal role mapping DTO used by portal.
type MappingsRepresentation struct {
	RealmMappings  []RoleRepresentation                    `json:"realmMappings,omitempty"`
	ClientMappings map[string]ClientMappingsRepresentation `json:"clientMappings,omitempty"`
}

// ClientMappingsRepresentation is the per-client role mapping DTO.
type ClientMappingsRepresentation struct {
	ID       string               `json:"id,omitempty"`
	Client   string               `json:"client,omitempty"`
	Mappings []RoleRepresentation `json:"mappings,omitempty"`
}

// TokenResponse is the client_credentials token response DTO.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// ListClientsOptions controls client listing.
type ListClientsOptions struct {
	First int
	Max   int
}

// ListUsersOptions controls user listing.
type ListUsersOptions struct {
	First  int
	Max    int
	Search string
}
