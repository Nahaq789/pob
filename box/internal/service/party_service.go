package service

import (
	"context"
	"log/slog"

	"pob/box/internal/model"
	"pob/box/internal/repository"

	"github.com/google/uuid"
)

type PartyService struct {
	repo *repository.PartyRepository
}

func NewPartyService(repo *repository.PartyRepository) *PartyService {
	return &PartyService{repo: repo}
}

func (s *PartyService) Create(ctx context.Context, userId uuid.UUID, name string) (*model.Party, error) {
	slog.InfoContext(ctx, "create party", slog.String("user_id", userId.String()))

	party := model.NewParty(userId, name)
	if err := s.repo.Create(ctx, party); err != nil {
		slog.ErrorContext(ctx, "failed to create party", slog.String("user_id", userId.String()), slog.Any("error", err))
		return nil, err
	}

	slog.InfoContext(ctx, "party created", slog.String("party_id", party.PartyId.String()))
	return &party, nil
}

func (s *PartyService) GetParties(ctx context.Context, userId uuid.UUID) ([]model.Party, error) {
	parties, err := s.repo.FindByUserId(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get parties", slog.String("user_id", userId.String()), slog.Any("error", err))
		return nil, err
	}
	return parties, nil
}

func (s *PartyService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "failed to delete party", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "party deleted", slog.String("id", id.String()))
	return nil
}

func (s *PartyService) UpdateName(ctx context.Context, id uuid.UUID, name string) error {
    slog.InfoContext(ctx, "update party name", slog.String("id", id.String()))
    if err := s.repo.UpdateName(ctx, id, name); err != nil {
        slog.ErrorContext(ctx, "failed to update party name", slog.String("id", id.String()), slog.Any("error", err))
        return err
    }
    slog.InfoContext(ctx, "party name updated", slog.String("id", id.String()))
    return nil
}