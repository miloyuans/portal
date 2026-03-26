package kcadmin

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TokenProvider provides and caches admin bearer tokens.
type TokenProvider interface {
	Token(ctx context.Context) (string, error)
	Invalidate()
}

// ClientCredentialsTokenProvider gets admin tokens via client_credentials.
type ClientCredentialsTokenProvider struct {
	cfg       Config
	http      *http.Client
	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

// NewClientCredentialsTokenProvider creates a token provider.
func NewClientCredentialsTokenProvider(cfg Config, httpClient *http.Client) *ClientCredentialsTokenProvider {
	return &ClientCredentialsTokenProvider{
		cfg:  cfg,
		http: httpClient,
	}
}

// Token returns a cached token or refreshes it when needed.
func (p *ClientCredentialsTokenProvider) Token(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.token != "" && time.Until(p.expiresAt) > 30*time.Second {
		return p.token, nil
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", p.cfg.AdminClientID)
	form.Set("client_secret", p.cfg.AdminClientSecret)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.TokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := p.http.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return "", &APIError{StatusCode: response.StatusCode, Path: p.cfg.TokenEndpoint}
	}

	var payload TokenResponse
	if err := decodeJSON(response.Body, &payload); err != nil {
		return "", err
	}

	p.token = payload.AccessToken
	p.expiresAt = time.Now().Add(time.Duration(payload.ExpiresIn) * time.Second)
	return p.token, nil
}

// Invalidate forces the next Token call to refresh.
func (p *ClientCredentialsTokenProvider) Invalidate() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.token = ""
	p.expiresAt = time.Time{}
}
