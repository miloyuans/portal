package session

import (
	"testing"
	"time"

	"portal/internal/config"
	"portal/internal/model"
)

func TestValidateRejectsExpiredSession(t *testing.T) {
	manager := NewManager(nil, nil, config.Config{})
	session := model.PortalSession{
		ExpiresAt:          time.Now().UTC().Add(-time.Minute),
		LastSeenAt:         time.Now().UTC(),
		IdleTimeoutMinutes: 15,
	}

	if err := manager.Validate(session); err != ErrSessionExpired {
		t.Fatalf("expected expired session, got %v", err)
	}
}

func TestValidateRejectsIdleTimeout(t *testing.T) {
	manager := NewManager(nil, nil, config.Config{})
	session := model.PortalSession{
		ExpiresAt:          time.Now().UTC().Add(time.Hour),
		LastSeenAt:         time.Now().UTC().Add(-20 * time.Minute),
		IdleTimeoutMinutes: 15,
	}

	if err := manager.Validate(session); err != ErrSessionExpired {
		t.Fatalf("expected idle timeout expiration, got %v", err)
	}
}
