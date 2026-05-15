package model

import (
	"github.com/google/uuid"
)

type Jwt struct {
	UserId       uuid.UUID
	Token        string
	RefreshToken string
}
