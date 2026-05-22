package service

import (
	"context"

	"pob/dex/internal/model"
	"pob/dex/internal/repository"
)

type AbilityService struct {
	abilityRepo *repository.AbilityRepository
}

func NewAbilityService(abilityRepo *repository.AbilityRepository) *AbilityService {
	return &AbilityService{abilityRepo: abilityRepo}
}

func (s *AbilityService) GetAbility(ctx context.Context, id int) (*model.Ability, error) {
	return s.abilityRepo.FindById(ctx, id)
}
