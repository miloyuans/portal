package kcadmin

import (
	"net/http"

	appconfig "portal/internal/config"
)

// Client wraps the Keycloak Admin REST API.
type Client struct {
	cfg    Config
	http   *http.Client
	tokens TokenProvider
}

// NewClient creates a Keycloak Admin API client.
func NewClient(cfg appconfig.Config) *Client {
	adminCfg := NewConfig(cfg)
	httpClient := &http.Client{Timeout: adminCfg.HTTPTimeout}
	tokenProvider := NewClientCredentialsTokenProvider(adminCfg, httpClient)
	return &Client{
		cfg:    adminCfg,
		http:   httpClient,
		tokens: tokenProvider,
	}
}
