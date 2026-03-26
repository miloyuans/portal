package kcadmin

import "context"

// GetRealm loads a realm representation.
func (c *Client) GetRealm(ctx context.Context, realm string) (RealmRepresentation, error) {
	var out RealmRepresentation
	err := c.doJSON(ctx, "GET", realmPath(realm), nil, &out)
	return out, err
}
