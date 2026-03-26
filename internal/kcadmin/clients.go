package kcadmin

import (
	"context"
	"net/url"
)

// ListClients lists clients for a realm.
func (c *Client) ListClients(ctx context.Context, realm string, opt ListClientsOptions) ([]ClientRepresentation, error) {
	values := url.Values{}
	if opt.First > 0 {
		values.Set("first", intToString(opt.First))
	}
	if opt.Max > 0 {
		values.Set("max", intToString(opt.Max))
	}

	var out []ClientRepresentation
	err := c.doJSON(ctx, "GET", withQuery(clientsPath(realm), values), nil, &out)
	return out, err
}

// GetClientByUUID loads a client by internal UUID.
func (c *Client) GetClientByUUID(ctx context.Context, realm, clientUUID string) (ClientRepresentation, error) {
	var out ClientRepresentation
	err := c.doJSON(ctx, "GET", clientPath(realm, clientUUID), nil, &out)
	return out, err
}

// GetClientByClientID searches a client by clientId.
func (c *Client) GetClientByClientID(ctx context.Context, realm, clientID string) (ClientRepresentation, error) {
	values := url.Values{}
	values.Set("clientId", clientID)
	var out []ClientRepresentation
	if err := c.doJSON(ctx, "GET", withQuery(clientsPath(realm), values), nil, &out); err != nil {
		return ClientRepresentation{}, err
	}
	if len(out) == 0 {
		return ClientRepresentation{}, &APIError{StatusCode: 404, Path: clientsPath(realm)}
	}
	return out[0], nil
}

// GetServiceAccountUser returns the service account user for a client.
func (c *Client) GetServiceAccountUser(ctx context.Context, realm, clientUUID string) (UserRepresentation, error) {
	var out UserRepresentation
	err := c.doJSON(ctx, "GET", clientServiceAccountPath(realm, clientUUID), nil, &out)
	return out, err
}
