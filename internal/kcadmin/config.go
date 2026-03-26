package kcadmin

import (
	"strings"
	"time"

	appconfig "portal/internal/config"
)

// Config stores Keycloak Admin API client configuration.
type Config struct {
	BaseURL           string
	Realm             string
	AdminClientID     string
	AdminClientSecret string
	TokenEndpoint     string
	AdminAPIBaseURL   string
	HTTPTimeout       time.Duration
}

// NewConfig builds a kcadmin config from the app config.
func NewConfig(cfg appconfig.Config) Config {
	baseURL := strings.TrimRight(cfg.Keycloak.BaseURL, "/")
	return Config{
		BaseURL:           baseURL,
		Realm:             cfg.Keycloak.Realm,
		AdminClientID:     cfg.Keycloak.AdminClientID,
		AdminClientSecret: cfg.Keycloak.AdminClientSecret,
		TokenEndpoint:     baseURL + "/realms/" + cfg.Keycloak.Realm + "/protocol/openid-connect/token",
		AdminAPIBaseURL:   baseURL + "/admin/realms",
		HTTPTimeout:       cfg.Keycloak.RequestTimeout,
	}
}
