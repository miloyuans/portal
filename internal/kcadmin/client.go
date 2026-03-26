package kcadmin

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"portal/internal/config"
	"portal/internal/model"
)

type Client struct {
	cfg    config.Config
	http   *http.Client
	logger *slog.Logger
}

type realmResponse struct {
	Realm               string `json:"realm"`
	DisplayName         string `json:"displayName"`
	DisplayNameHTML     string `json:"displayNameHtml"`
	LoginTheme          string `json:"loginTheme"`
	AccountTheme        string `json:"accountTheme"`
	AdminTheme          string `json:"adminTheme"`
	RegistrationAllowed bool   `json:"registrationAllowed"`
	SSLRequired         string `json:"sslRequired"`
}

type clientResponse struct {
	ID           string            `json:"id"`
	ClientID     string            `json:"clientId"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	RootURL      string            `json:"rootUrl"`
	BaseURL      string            `json:"baseUrl"`
	RedirectURIs []string          `json:"redirectUris"`
	WebOrigins   []string          `json:"webOrigins"`
	Enabled      bool              `json:"enabled"`
	PublicClient bool              `json:"publicClient"`
	Protocol     string            `json:"protocol"`
	Attributes   map[string]string `json:"attributes"`
}

type userResponse struct {
	ID            string              `json:"id"`
	Username      string              `json:"username"`
	Email         string              `json:"email"`
	FirstName     string              `json:"firstName"`
	LastName      string              `json:"lastName"`
	Enabled       bool                `json:"enabled"`
	EmailVerified bool                `json:"emailVerified"`
	Attributes    map[string][]string `json:"attributes"`
}

type roleResponse struct {
	Name string `json:"name"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewClient(cfg config.Config, logger *slog.Logger) *Client {
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Keycloak.RequestTimeoutSecs) * time.Second,
	}
	if cfg.Keycloak.SkipTLSVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		}
	}
	return &Client{
		cfg:    cfg,
		http:   httpClient,
		logger: logger,
	}
}

func (c *Client) SyncData(ctx context.Context, userID string) (model.RealmProjection, []model.ClientProjection, model.UserProjection, error) {
	token, err := c.adminToken(ctx)
	if err != nil {
		return model.RealmProjection{}, nil, model.UserProjection{}, err
	}

	realm, err := c.getRealm(ctx, token)
	if err != nil {
		return model.RealmProjection{}, nil, model.UserProjection{}, err
	}

	clients, err := c.listClients(ctx, token)
	if err != nil {
		return model.RealmProjection{}, nil, model.UserProjection{}, err
	}

	user, err := c.getUser(ctx, token, userID)
	if err != nil {
		return model.RealmProjection{}, nil, model.UserProjection{}, err
	}

	realmRoles, err := c.getUserRealmRoles(ctx, token, userID)
	if err != nil {
		return model.RealmProjection{}, nil, model.UserProjection{}, err
	}

	clientRoles := make(map[string][]string)
	for _, client := range clients {
		roles, err := c.getUserClientRoles(ctx, token, userID, client.ClientUUID)
		if err != nil {
			return model.RealmProjection{}, nil, model.UserProjection{}, err
		}
		if len(roles) > 0 {
			clientRoles[client.ClientID] = roles
		}
	}

	user.RealmRoles = realmRoles
	user.ClientRoles = clientRoles
	return realm, clients, user, nil
}

func (c *Client) adminToken(ctx context.Context) (string, error) {
	form := url.Values{}
	endpointRealm := c.cfg.Keycloak.Realm

	if c.cfg.Keycloak.AdminUsername != "" && c.cfg.Keycloak.AdminPassword != "" {
		form.Set("grant_type", "password")
		form.Set("client_id", "admin-cli")
		form.Set("username", c.cfg.Keycloak.AdminUsername)
		form.Set("password", c.cfg.Keycloak.AdminPassword)
		endpointRealm = c.cfg.Keycloak.AdminRealm
	} else {
		form.Set("grant_type", "client_credentials")
		form.Set("client_id", c.cfg.Keycloak.AdminClientID)
		form.Set("client_secret", c.cfg.Keycloak.AdminClientSecret)
	}

	endpoint := strings.TrimRight(c.cfg.Keycloak.BaseURL, "/") + "/realms/" + endpointRealm + "/protocol/openid-connect/token"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.http.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("admin token request failed: %s", string(body))
	}

	var payload tokenResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", err
	}
	return payload.AccessToken, nil
}

func (c *Client) getRealm(ctx context.Context, token string) (model.RealmProjection, error) {
	var payload realmResponse
	if err := c.getJSON(ctx, token, "/admin/realms/"+c.cfg.Keycloak.Realm, &payload); err != nil {
		return model.RealmProjection{}, err
	}

	now := time.Now().UTC()
	return model.RealmProjection{
		Realm:               payload.Realm,
		DisplayName:         payload.DisplayName,
		DisplayNameHTML:     payload.DisplayNameHTML,
		LoginTheme:          payload.LoginTheme,
		AccountTheme:        payload.AccountTheme,
		AdminTheme:          payload.AdminTheme,
		RegistrationAllowed: payload.RegistrationAllowed,
		SSLRequired:         payload.SSLRequired,
		SyncedAt:            now,
		UpdatedAt:           now,
	}, nil
}

func (c *Client) listClients(ctx context.Context, token string) ([]model.ClientProjection, error) {
	var payload []clientResponse
	if err := c.getJSON(ctx, token, "/admin/realms/"+c.cfg.Keycloak.Realm+"/clients?max=500", &payload); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	out := make([]model.ClientProjection, 0, len(payload))
	for _, item := range payload {
		out = append(out, model.ClientProjection{
			Realm:        c.cfg.Keycloak.Realm,
			ClientUUID:   item.ID,
			ClientID:     item.ClientID,
			Name:         item.Name,
			Description:  item.Description,
			RootURL:      item.RootURL,
			BaseURL:      item.BaseURL,
			RedirectURIs: item.RedirectURIs,
			WebOrigins:   item.WebOrigins,
			Enabled:      item.Enabled,
			PublicClient: item.PublicClient,
			Protocol:     item.Protocol,
			Attributes:   item.Attributes,
			SyncedAt:     now,
			UpdatedAt:    now,
		})
	}
	return out, nil
}

func (c *Client) getUser(ctx context.Context, token, userID string) (model.UserProjection, error) {
	var payload userResponse
	if err := c.getJSON(ctx, token, "/admin/realms/"+c.cfg.Keycloak.Realm+"/users/"+userID, &payload); err != nil {
		return model.UserProjection{}, err
	}

	now := time.Now().UTC()
	return model.UserProjection{
		Realm:         c.cfg.Keycloak.Realm,
		UserID:        payload.ID,
		Username:      payload.Username,
		Email:         payload.Email,
		FirstName:     payload.FirstName,
		LastName:      payload.LastName,
		Enabled:       payload.Enabled,
		EmailVerified: payload.EmailVerified,
		Attributes:    payload.Attributes,
		SyncedAt:      now,
		UpdatedAt:     now,
	}, nil
}

func (c *Client) getUserRealmRoles(ctx context.Context, token, userID string) ([]string, error) {
	var payload []roleResponse
	if err := c.getJSON(ctx, token, "/admin/realms/"+c.cfg.Keycloak.Realm+"/users/"+userID+"/role-mappings/realm", &payload); err != nil {
		return nil, err
	}

	out := make([]string, 0, len(payload))
	for _, role := range payload {
		out = append(out, role.Name)
	}
	return out, nil
}

func (c *Client) getUserClientRoles(ctx context.Context, token, userID, clientUUID string) ([]string, error) {
	var payload []roleResponse
	if err := c.getJSON(ctx, token, "/admin/realms/"+c.cfg.Keycloak.Realm+"/users/"+userID+"/role-mappings/clients/"+clientUUID, &payload); err != nil {
		return nil, err
	}

	out := make([]string, 0, len(payload))
	for _, role := range payload {
		out = append(out, role.Name)
	}
	return out, nil
}

func (c *Client) getJSON(ctx context.Context, token, path string, target any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(c.cfg.Keycloak.BaseURL, "/")+path, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("keycloak admin request failed (%s): %s", path, string(body))
	}
	return json.NewDecoder(response.Body).Decode(target)
}
