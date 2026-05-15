package apperror

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrAlreadyLoggedOut = errors.New("already logged out")
