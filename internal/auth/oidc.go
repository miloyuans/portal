package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"portal/internal/config"
	"portal/internal/model"
)

// OIDCClient wraps the Keycloak OIDC browser login flow.
type OIDCClient struct {
	cfg      config.Config
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2   *oauth2.Config
	http     *http.Client
}

// TokenBundle stores the exchanged OIDC tokens and verified claims.
type TokenBundle struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	Claims       model.IDTokenClaims
}

// NewOIDCClient creates a new OIDC browser-flow client.
func NewOIDCClient(ctx context.Context, cfg config.Config) (*OIDCClient, error) {
	httpClient := &http.Client{Timeout: cfg.Keycloak.RequestTimeout}
	issuer := strings.TrimRight(cfg.Keycloak.BaseURL, "/") + "/realms/" + cfg.Keycloak.Realm

	oidcContext := oidc.ClientContext(ctx, httpClient)
	provider, err := oidc.NewProvider(oidcContext, issuer)
	if err != nil {
		return nil, err
	}

	return &OIDCClient{
		cfg:      cfg,
		provider: provider,
		verifier: provider.Verifier(&oidc.Config{ClientID: cfg.Keycloak.OIDCClientID}),
		oauth2: &oauth2.Config{
			ClientID:     cfg.Keycloak.OIDCClientID,
			ClientSecret: cfg.Keycloak.OIDCClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  oidcAuthorizationURL(cfg),
				TokenURL: provider.Endpoint().TokenURL,
			},
			RedirectURL:  cfg.Keycloak.RedirectURL,
			Scopes:       cfg.Keycloak.OIDCScopes,
		},
		http: httpClient,
	}, nil
}

// AuthCodeURL returns the Keycloak authorization URL.
func (c *OIDCClient) AuthCodeURL(state, nonce string) string {
	return c.oauth2.AuthCodeURL(state, oauth2.AccessTypeOffline, oidc.Nonce(nonce))
}

// Exchange exchanges the authorization code and verifies the returned ID token.
func (c *OIDCClient) Exchange(ctx context.Context, code, expectedNonce string) (TokenBundle, error) {
	oidcContext := oidc.ClientContext(ctx, c.http)
	token, err := c.oauth2.Exchange(oidcContext, code)
	if err != nil {
		return TokenBundle{}, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return TokenBundle{}, fmt.Errorf("missing id_token in oauth response")
	}

	idToken, err := c.verifier.Verify(oidcContext, rawIDToken)
	if err != nil {
		return TokenBundle{}, err
	}
	if idToken.Nonce != expectedNonce {
		return TokenBundle{}, fmt.Errorf("oidc nonce mismatch")
	}

	var claims model.IDTokenClaims
	if err := idToken.Claims(&claims); err != nil {
		return TokenBundle{}, err
	}

	return TokenBundle{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      rawIDToken,
		Claims:       claims,
	}, nil
}

// LogoutURL returns the Keycloak RP-initiated logout URL.
func (c *OIDCClient) LogoutURL(idTokenHint string, redirectURI string) string {
	if redirectURI == "" {
		redirectURI = c.cfg.Keycloak.PostLogoutRedirectURL
	}

	values := url.Values{}
	if idTokenHint != "" {
		values.Set("id_token_hint", idTokenHint)
	}
	values.Set("post_logout_redirect_uri", redirectURI)
	values.Set("client_id", c.cfg.Keycloak.OIDCClientID)

	base := strings.TrimRight(c.cfg.Keycloak.PublicURL, "/") + "/realms/" + c.cfg.Keycloak.Realm + "/protocol/openid-connect/logout"
	return base + "?" + values.Encode()
}

// Ready reports whether the OIDC provider metadata is initialized.
func (c *OIDCClient) Ready() bool {
	return c.provider != nil
}

func oidcAuthorizationURL(cfg config.Config) string {
	return strings.TrimRight(cfg.Keycloak.PublicURL, "/") + "/realms/" + cfg.Keycloak.Realm + "/protocol/openid-connect/auth"
}
