package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId       uuid.UUID
	UserName     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(userName, passwordHash string) User {
	return User{
		UserId:       uuid.New(),
		UserName:     userName,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

