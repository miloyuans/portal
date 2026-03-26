package kcadmin

import (
	"context"
	"net/url"
)

// ListUsers lists users in a realm.
func (c *Client) ListUsers(ctx context.Context, realm string, opt ListUsersOptions) ([]UserRepresentation, error) {
	values := url.Values{}
	if opt.First > 0 {
		values.Set("first", intToString(opt.First))
	}
	if opt.Max > 0 {
		values.Set("max", intToString(opt.Max))
	}
	if opt.Search != "" {
		values.Set("search", opt.Search)
	}

	var out []UserRepresentation
	err := c.doJSON(ctx, "GET", withQuery(usersPath(realm), values), nil, &out)
	return out, err
}

// GetUserByID loads a user by ID.
func (c *Client) GetUserByID(ctx context.Context, realm, userID string) (UserRepresentation, error) {
	var out UserRepresentation
	err := c.doJSON(ctx, "GET", userPath(realm, userID), nil, &out)
	return out, err
}
