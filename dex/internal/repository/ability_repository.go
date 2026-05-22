package repository

import (
	"context"

	"pob/dex/internal/model"
	"pob/dex/internal/model/entity"
	"pob/dex/internal/shared"
)

type AbilityRepository struct {
	db *shared.DBClient
}

func NewAbilityRepository(db *shared.DBClient) *AbilityRepository {
	return &AbilityRepository{db: db}
}

func (a *AbilityRepository) FindById(ctx context.Context, id int) (*model.Ability, error) {
	var e entity.Ability
	if err := a.db.GetClient().WithContext(ctx).First(&e, id).Error; err != nil {
		return nil, err
	}
	return &model.Ability{
		AbilityId: e.Id,
		Name:      e.Name,
	}, nil
}
