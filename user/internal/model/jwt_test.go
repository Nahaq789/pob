package model

import (
	"github.com/google/uuid"
	"testing"
)

// TestJwtStruct ensures that a Jwt struct can be created and its fields
// accessed correctly. This test mimics typical usage in the service layer.
func TestJwtStruct(t *testing.T) {
	uid := uuid.New()
	tok := "access-abc"
	rtok := "refresh-xyz"

	j := Jwt{UserId: uid, Token: tok, RefreshToken: rtok}

	if j.UserId != uid {
		t.Fatalf("expected UserId %v, got %v", uid, j.UserId)
	}
	if j.Token != tok {
		t.Fatalf("expected Token %q, got %q", tok, j.Token)
	}
	if j.RefreshToken != rtok {
		t.Fatalf("expected RefreshToken %q, got %q", rtok, j.RefreshToken)
	}
}
