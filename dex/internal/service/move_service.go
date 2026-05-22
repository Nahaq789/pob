package service

import (
	"context"

	"pob/dex/internal/model"
	"pob/dex/internal/repository"
)

type MoveService struct {
	moveRepo *repository.MoveRepository
}

func NewMoveService(moveRepo *repository.MoveRepository) *MoveService {
	return &MoveService{moveRepo: moveRepo}
}

func (s *MoveService) GetMove(ctx context.Context, id int) (*model.Move, error) {
	e, err := s.moveRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	var power, accuracy int
	if e.Power != nil {
		power = *e.Power
	}
	if e.Accuracy != nil {
		accuracy = *e.Accuracy
	}

	return &model.Move{
		MoveId:      e.Id,
		Name:        e.Name,
		TypeId:      e.TypeId,
		DamageClass: e.DamageClass,
		Power:       power,
		Accuracy:    accuracy,
		Pp:          e.Pp,
	}, nil
}
