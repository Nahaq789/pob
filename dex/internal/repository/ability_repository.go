package repository

import (
	"context"

	"pob/dex/internal/model/entity"
	"pob/dex/internal/shared"
)

type AbilityRepository struct {
	db *shared.DBClient
}

func NewAbilityRepository(db *shared.DBClient) *AbilityRepository {
	return &AbilityRepository{db: db}
}

func (a *AbilityRepository) FindById(ctx context.Context, id int) (*entity.Ability, error) {
	var ability entity.Ability
	if err := a.db.GetClient().WithContext(ctx).First(&ability, id).Error; err != nil {
		return nil, err
	}
	return &ability, nil
}
