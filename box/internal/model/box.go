package model

import (
	"time"

	"github.com/google/uuid"
)

type Box struct {
	BoxId     uuid.UUID `json:"box_id"`
	UserId    uuid.UUID `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBox(userId uuid.UUID, name string) Box {
	now := time.Now()
	return Box{
		BoxId:     uuid.New(),
		UserId:    userId,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func FromBox(id, userId uuid.UUID, name string, createdAt, updatedAt time.Time) Box {
	return Box{
		BoxId:     id,
		UserId:    userId,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
