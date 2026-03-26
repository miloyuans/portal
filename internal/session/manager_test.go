package session

import (
	"testing"
	"time"

	"portal/internal/config"
	"portal/internal/model"
)

func TestValidateRejectsExpiredSession(t *testing.T) {
	manager := NewManager(nil, config.Config{})
	session := model.PortalSession{
		ExpiresAt:         time.Now().UTC().Add(-time.Minute),
		AbsoluteExpiresAt: time.Now().UTC().Add(time.Hour),
	}

	if err := manager.Validate(session); err != ErrSessionExpired {
		t.Fatalf("expected expired session, got %v", err)
	}
}

func TestValidateRejectsAbsoluteExpiration(t *testing.T) {
	manager := NewManager(nil, config.Config{})
	session := model.PortalSession{
		ExpiresAt:         time.Now().UTC().Add(time.Hour),
		AbsoluteExpiresAt: time.Now().UTC().Add(-time.Minute),
	}

	if err := manager.Validate(session); err != ErrSessionExpired {
		t.Fatalf("expected absolute expiration, got %v", err)
	}
}
