package service

import (
	"context"
	"fmt"
	"log/slog"

	"pob/box/internal/model"
	"pob/box/internal/repository"
	"pob/box/internal/service/dto"

	"github.com/google/uuid"
)

const partyPokemonLimit = 6

type PartyPokemonService struct {
	repo    *repository.PartyPokemonRepository
	boxRepo *repository.BoxPokemonRepository
}

func NewPartyPokemonService(repo *repository.PartyPokemonRepository, boxRepo *repository.BoxPokemonRepository) *PartyPokemonService {
	return &PartyPokemonService{repo: repo, boxRepo: boxRepo}
}

func (s *PartyPokemonService) AddPokemon(ctx context.Context, partyId uuid.UUID, boxPokemonId uuid.UUID, slot int) error {
	slog.InfoContext(ctx, "add party_pokemon", slog.String("party_id", partyId.String()), slog.String("box_pokemon_id", boxPokemonId.String()))

	if _, err := s.boxRepo.FindById(ctx, boxPokemonId); err != nil {
		slog.WarnContext(ctx, "box_pokemon not found", slog.String("box_pokemon_id", boxPokemonId.String()), slog.Any("error", err))
		return fmt.Errorf("box_pokemon %s not found: %w", boxPokemonId, err)
	}

	existing, err := s.repo.FindByPartyId(ctx, partyId)
	if err != nil {
		return err
	}
	if len(existing) >= partyPokemonLimit {
		return fmt.Errorf("party is full (limit: %d)", partyPokemonLimit)
	}

	pp := model.NewPartyPokemon(partyId, boxPokemonId, slot)
	if err := s.repo.AddPokemon(ctx, pp); err != nil {
		slog.ErrorContext(ctx, "failed to add party_pokemon", slog.Any("error", err))
		return err
	}

	slog.InfoContext(ctx, "party_pokemon added", slog.String("party_pokemon_id", pp.PartyPokemonId.String()))
	return nil
}

func (s *PartyPokemonService) GetByPartyId(ctx context.Context, partyId uuid.UUID) ([]model.PartyPokemon, error) {
	pps, err := s.repo.FindByPartyId(ctx, partyId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get party_pokemon", slog.String("party_id", partyId.String()), slog.Any("error", err))
		return nil, err
	}
	return pps, nil
}

func (s *PartyPokemonService) RemovePokemon(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.RemovePokemon(ctx, id); err != nil {
		slog.ErrorContext(ctx, "failed to remove party_pokemon", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "party_pokemon removed", slog.String("id", id.String()))
	return nil
}

func (s *PartyPokemonService) UpdateSlot(ctx context.Context, id uuid.UUID, slot int) error {
	if err := s.repo.UpdateSlot(ctx, id, slot); err != nil {
		slog.ErrorContext(ctx, "failed to update party_pokemon slot", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "party_pokemon slot updated", slog.String("id", id.String()), slog.Int("slot", slot))
	return nil
}

func (s *PartyPokemonService) SetPokemon(ctx context.Context, partyId uuid.UUID, entries []dto.PartyPokemonEntry) error {
	existing, err := s.repo.FindByPartyId(ctx, partyId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get party_pokemon", slog.String("party_id", partyId.String()), slog.Any("error", err))
		return err
	}

	for _, pp := range existing {
		if err := s.repo.RemovePokemon(ctx, pp.PartyPokemonId); err != nil {
			slog.ErrorContext(ctx, "failed to remove party_pokemon", slog.String("id", pp.PartyPokemonId.String()), slog.Any("error", err))
			return err
		}
	}

	for _, e := range entries {
		boxPokemonId, err := uuid.Parse(e.BoxPokemonId)
		if err != nil {
			return fmt.Errorf("invalid box_pokemon_id %q: %w", e.BoxPokemonId, err)
		}
		if _, err := s.boxRepo.FindById(ctx, boxPokemonId); err != nil {
			slog.WarnContext(ctx, "box_pokemon not found", slog.String("box_pokemon_id", e.BoxPokemonId), slog.Any("error", err))
			return fmt.Errorf("box_pokemon %s not found: %w", e.BoxPokemonId, err)
		}
		pp := model.NewPartyPokemon(partyId, boxPokemonId, e.Slot)
		if err := s.repo.AddPokemon(ctx, pp); err != nil {
			slog.ErrorContext(ctx, "failed to add party_pokemon", slog.String("box_pokemon_id", e.BoxPokemonId), slog.Any("error", err))
			return err
		}
	}

	slog.InfoContext(ctx, "party_pokemon set", slog.String("party_id", partyId.String()), slog.Int("count", len(entries)))
	return nil
}
