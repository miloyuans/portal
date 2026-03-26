package kcadmin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type stubTokenProvider struct {
	token           string
	invalidateCalls int32
}

func (s *stubTokenProvider) Token(context.Context) (string, error) { return s.token, nil }
func (s *stubTokenProvider) Invalidate()                           { atomic.AddInt32(&s.invalidateCalls, 1) }

func TestDoJSONRetriesOnceOnUnauthorized(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&calls, 1)
		if current == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"realm":"portal","enabled":true}`))
	}))
	defer server.Close()

	tokens := &stubTokenProvider{token: "token"}
	client := &Client{
		cfg: Config{
			AdminAPIBaseURL: server.URL,
		},
		http:   server.Client(),
		tokens: tokens,
	}

	var out RealmRepresentation
	if err := client.doJSON(context.Background(), http.MethodGet, "/portal", nil, &out); err != nil {
		t.Fatalf("expected retry to succeed, got %v", err)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("expected exactly two calls, got %d", calls)
	}
	if atomic.LoadInt32(&tokens.invalidateCalls) != 1 {
		t.Fatalf("expected one invalidate call, got %d", tokens.invalidateCalls)
	}
}
