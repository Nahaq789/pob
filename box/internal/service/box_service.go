package service

import (
	"context"
	"log/slog"

	"pob/box/internal/model"
	"pob/box/internal/repository"

	"github.com/google/uuid"
)

type BoxService struct {
	repo *repository.BoxRepository
}

func NewBoxService(repo *repository.BoxRepository) *BoxService {
	return &BoxService{repo: repo}
}

func (s *BoxService) Create(ctx context.Context, userId uuid.UUID, name string) (*model.Box, error) {
	slog.InfoContext(ctx, "create box", slog.String("user_id", userId.String()))

	box := model.NewBox(userId, name)
	if err := s.repo.Create(ctx, box); err != nil {
		slog.ErrorContext(ctx, "failed to create box", slog.String("user_id", userId.String()), slog.Any("error", err))
		return nil, err
	}

	slog.InfoContext(ctx, "box created", slog.String("box_id", box.BoxId.String()))
	return &box, nil
}

func (s *BoxService) GetBoxes(ctx context.Context, userId uuid.UUID) ([]model.Box, error) {
	boxes, err := s.repo.FindByUserId(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get boxes", slog.String("user_id", userId.String()), slog.Any("error", err))
		return nil, err
	}
	return boxes, nil
}

func (s *BoxService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "failed to delete box", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "box deleted", slog.String("id", id.String()))
	return nil
}
