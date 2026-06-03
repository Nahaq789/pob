package model

import (
	"time"

	"github.com/google/uuid"
)

type Party struct {
	PartyId   uuid.UUID `json:"party_id"`
	UserId    uuid.UUID `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewParty(userId uuid.UUID, name string) Party {
	now := time.Now()
	return Party{
		PartyId:   uuid.New(),
		UserId:    userId,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func FromParty(id, userId uuid.UUID, name string, createdAt, updatedAt time.Time) Party {
	return Party{
		PartyId:   id,
		UserId:    userId,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
