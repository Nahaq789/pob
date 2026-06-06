package model

import "testing"

// TestNewAuth verifies that NewAuth correctly constructs an Auth
// struct with the provided username and password.
func TestNewAuth(t *testing.T) {
	const (
		username = "testuser"
		password = "secret123"
	)

	au := NewAuth(username, password)
	if au.UserName != username {
		t.Fatalf("expected UserName %q, got %q", username, au.UserName)
	}
	if au.PasswordPlane != password {
		t.Fatalf("expected PasswordPlane %q, got %q", password, au.PasswordPlane)
	}
}
