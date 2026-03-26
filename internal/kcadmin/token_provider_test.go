package kcadmin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientCredentialsTokenProviderCachesToken(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"token-1","expires_in":300,"token_type":"Bearer"}`))
	}))
	defer server.Close()

	provider := NewClientCredentialsTokenProvider(Config{
		TokenEndpoint:     server.URL,
		AdminClientID:     "portal-sync",
		AdminClientSecret: "secret",
		HTTPTimeout:       2 * time.Second,
	}, server.Client())

	first, err := provider.Token(context.Background())
	if err != nil {
		t.Fatalf("first token call failed: %v", err)
	}
	second, err := provider.Token(context.Background())
	if err != nil {
		t.Fatalf("second token call failed: %v", err)
	}

	if first != second {
		t.Fatalf("expected cached token reuse")
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("expected one token endpoint call, got %d", calls)
	}
}
