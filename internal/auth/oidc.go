package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"portal/internal/config"
	"portal/internal/model"
)

type OIDCClient struct {
	cfg      config.Config
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2   *oauth2.Config
	http     *http.Client
}

type TokenBundle struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	Expiry       time.Time
	Claims       model.IDTokenClaims
}

func NewOIDCClient(ctx context.Context, cfg config.Config) (*OIDCClient, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Keycloak.RequestTimeoutSecs) * time.Second,
	}
	if cfg.Keycloak.SkipTLSVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		}
	}

	issuer := strings.TrimRight(cfg.Keycloak.BaseURL, "/") + "/realms/" + cfg.Keycloak.Realm
	oidcContext := oidc.ClientContext(ctx, httpClient)
	provider, err := oidc.NewProvider(oidcContext, issuer)
	if err != nil {
		return nil, err
	}

	return &OIDCClient{
		cfg:      cfg,
		provider: provider,
		verifier: provider.Verifier(&oidc.Config{ClientID: cfg.Keycloak.ClientID}),
		oauth2: &oauth2.Config{
			ClientID:     cfg.Keycloak.ClientID,
			ClientSecret: cfg.Keycloak.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  cfg.Keycloak.RedirectURL,
			Scopes:       cfg.Keycloak.Scopes,
		},
		http: httpClient,
	}, nil
}

func (c *OIDCClient) AuthCodeURL(state, nonce string) string {
	return c.oauth2.AuthCodeURL(state, oauth2.AccessTypeOffline, oidc.Nonce(nonce))
}

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
		return TokenBundle{}, fmt.Errorf("nonce mismatch")
	}

	var claims model.IDTokenClaims
	if err := idToken.Claims(&claims); err != nil {
		return TokenBundle{}, err
	}

	return TokenBundle{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      rawIDToken,
		Expiry:       token.Expiry,
		Claims:       claims,
	}, nil
}

func (c *OIDCClient) LogoutURL(idTokenHint string, redirectURI string) string {
	base := strings.TrimRight(c.cfg.Keycloak.BaseURL, "/") + "/realms/" + c.cfg.Keycloak.Realm + "/protocol/openid-connect/logout"
	values := url.Values{}
	if idTokenHint != "" {
		values.Set("id_token_hint", idTokenHint)
	}
	if redirectURI == "" {
		redirectURI = c.cfg.Keycloak.LogoutRedirectURL
	}
	values.Set("post_logout_redirect_uri", redirectURI)
	values.Set("client_id", c.cfg.Keycloak.ClientID)
	return base + "?" + values.Encode()
}

func (c *OIDCClient) Ready() bool {
	return c.provider != nil
}
