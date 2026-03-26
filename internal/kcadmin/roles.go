package kcadmin

import (
	"context"
	"net/url"
)

// GetUserRoleMappings returns the raw role mappings of a user.
func (c *Client) GetUserRoleMappings(ctx context.Context, realm, userID string) (MappingsRepresentation, error) {
	var out MappingsRepresentation
	err := c.doJSON(ctx, "GET", userRoleMappingsPath(realm, userID), nil, &out)
	return out, err
}

// GetUserEffectiveRealmRoles returns effective realm roles for a user.
func (c *Client) GetUserEffectiveRealmRoles(ctx context.Context, realm, userID string, brief bool) ([]RoleRepresentation, error) {
	values := url.Values{}
	if brief {
		values.Set("briefRepresentation", "true")
	}

	var out []RoleRepresentation
	err := c.doJSON(ctx, "GET", withQuery(userEffectiveRealmRolesPath(realm, userID), values), nil, &out)
	return out, err
}

// GetUserEffectiveClientRoles returns effective client roles for a user on one client.
func (c *Client) GetUserEffectiveClientRoles(ctx context.Context, realm, userID, clientUUID string) ([]RoleRepresentation, error) {
	var out []RoleRepresentation
	err := c.doJSON(ctx, "GET", userEffectiveClientRolesPath(realm, userID, clientUUID), nil, &out)
	return out, err
}
