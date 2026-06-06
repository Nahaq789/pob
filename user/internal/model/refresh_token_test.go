package model

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

// TestNewRefreshToken verifies that NewRefreshToken sets fields correctly
// and calculates expiration 7 days from now.
func TestNewRefreshToken(t *testing.T) {
	uid := uuid.New()
	hash := "hash-abc"
	rt := NewRefreshToken(uid, hash)
	if rt.UserId != uid {
		t.Fatalf("expected UserId %v, got %v", uid, rt.UserId)
	}
	if rt.TokenHash != hash {
		t.Fatalf("expected TokenHash %s, got %s", hash, rt.TokenHash)
	}
	// Expiration should be roughly 7 days ahead; allow small delta.
	if rt.ExpiredAt.Sub(time.Now()) < 6*time.Hour*24 { // 6 days
		t.Fatalf("expected ExpiredAt at least 6 days in future, got %v", rt.ExpiredAt)
	}
}

// TestFromRefreshToken ensures FromRefreshToken creates a struct with
// the supplied values unchanged.
func TestFromRefreshToken(t *testing.T) {
	id := uuid.New()
	uid := uuid.New()
	hash := "hash-xyz"
	expired := time.Now().Add(12 * time.Hour)
	created := time.Now().Add(-1 * time.Hour)
	rt := FromRefreshToken(id, uid, hash, expired, created)
	if rt.RefreshTokenId != id || rt.UserId != uid || rt.TokenHash != hash {
		t.Fatalf("struct fields not set correctly: %+v", rt)
	}
	if !rt.ExpiredAt.Equal(expired) || !rt.CreatedAt.Equal(created) {
		t.Fatalf("timestamps not preserved: %+v", rt)
	}
}
