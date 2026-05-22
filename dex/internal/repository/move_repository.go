package repository

import (
	"context"

	"pob/dex/internal/model/entity"
	"pob/dex/internal/shared"
)

type MoveRepository struct {
	db *shared.DBClient
}

func NewMoveRepository(db *shared.DBClient) *MoveRepository {
	return &MoveRepository{db: db}
}

func (m *MoveRepository) FindById(ctx context.Context, id int) (*entity.Move, error) {
	var move entity.Move
	if err := m.db.GetClient().WithContext(ctx).First(&move, id).Error; err != nil {
		return nil, err
	}
	return &move, nil
}
