package apperror

import "testing"

func TestAuthErrorMessages(t *testing.T) {
    if ErrInvalidCredentials == nil {
        t.Fatalf("ErrInvalidCredentials is nil")
    }
    if got := ErrInvalidCredentials.Error(); got != "invalid credentials" {
        t.Errorf("ErrInvalidCredentials.Error() = %q, want %q", got, "invalid credentials")
    }

    if ErrAlreadyLoggedOut == nil {
        t.Fatalf("ErrAlreadyLoggedOut is nil")
    }
    if got := ErrAlreadyLoggedOut.Error(); got != "already logged out" {
        t.Errorf("ErrAlreadyLoggedOut.Error() = %q, want %q", got, "already logged out")
    }
}
