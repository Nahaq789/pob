package model

import (
	"github.com/google/uuid"
	"testing"
)

// TestNewUser verifies that NewUser sets fields appropriately and
// generates a new UUID.
func TestNewUser(t *testing.T) {
	name := "alice"
	hash := "passhash"
	u := NewUser(name, hash)
	if u.UserName != name {
		t.Fatalf("expected UserName %s, got %s", name, u.UserName)
	}
	if u.PasswordHash != hash {
		t.Fatalf("expected PasswordHash %s, got %s", hash, u.PasswordHash)
	}
	if u.UserId == uuid.Nil {
		t.Fatalf("expected non-nil UserId, got %v", u.UserId)
	}
}
