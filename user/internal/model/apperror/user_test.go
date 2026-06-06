package apperror

import "testing"

func TestUserErrorMessages(t *testing.T) {
    if ErrUserNotFound == nil {
        t.Fatalf("ErrUserNotFound is nil")
    }
    if got := ErrUserNotFound.Error(); got != "user not found" {
        t.Errorf("ErrUserNotFound.Error() = %q, want %q", got, "user not found")
    }
}
