package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"pob/box/internal/model"
	"pob/box/internal/repository"
	"pob/box/internal/service/dto"
	gen "pob/box/proto"

	"github.com/google/uuid"
)

const boxPokemonLimit = 30

type BoxPokemonService struct {
	repo *repository.BoxPokemonRepository
	dex  gen.DexServiceClient
}

func NewBoxPokemonService(repo *repository.BoxPokemonRepository, dex gen.DexServiceClient) *BoxPokemonService {
	return &BoxPokemonService{repo: repo, dex: dex}
}

func (s *BoxPokemonService) Add(ctx context.Context, boxId uuid.UUID, req dto.AddBoxPokemonRequest) (*model.BoxPokemon, error) {
	slog.InfoContext(ctx, "add box_pokemon", slog.String("box_id", boxId.String()), slog.Int("pokemon_id", req.PokemonId))

	if _, err := s.dex.GetPokemon(ctx, &gen.GetPokemonRequest{PokemonId: int32(req.PokemonId)}); err != nil {
		slog.WarnContext(ctx, "pokemon not found", slog.Int("pokemon_id", req.PokemonId), slog.Any("error", err))
		return nil, fmt.Errorf("pokemon %d not found: %w", req.PokemonId, err)
	}

	if _, err := s.dex.GetAbility(ctx, &gen.GetAbilityRequest{AbilityId: int32(req.AbilityId)}); err != nil {
		slog.WarnContext(ctx, "ability not found", slog.Int("ability_id", req.AbilityId), slog.Any("error", err))
		return nil, fmt.Errorf("ability %d not found: %w", req.AbilityId, err)
	}

	for i, moveId := range []*int{req.Move1Id, req.Move2Id, req.Move3Id, req.Move4Id} {
		if moveId == nil {
			continue
		}
		if _, err := s.dex.GetMove(ctx, &gen.GetMoveRequest{MoveId: int32(*moveId)}); err != nil {
			slog.WarnContext(ctx, "move not found", slog.Int("slot", i+1), slog.Int("move_id", *moveId), slog.Any("error", err))
			return nil, fmt.Errorf("move %d not found: %w", *moveId, err)
		}
	}

	existing, err := s.repo.FindByBoxId(ctx, boxId)
	if err != nil {
		return nil, err
	}
	if len(existing) >= boxPokemonLimit {
		return nil, fmt.Errorf("box is full (limit: %d)", boxPokemonLimit)
	}

	now := time.Now()
	bp := model.BoxPokemon{
		BoxPokemonId: uuid.New(),
		BoxId:        boxId,
		PokemonId:    req.PokemonId,
		Nickname:     req.Nickname,
		AbilityId:    req.AbilityId,
		Nature:       req.Nature,
		HeldItemId:   req.HeldItemId,
		IvHp:         req.IvHp,
		IvAttack:     req.IvAttack,
		IvDefense:    req.IvDefense,
		IvSpAttack:   req.IvSpAttack,
		IvSpDefense:  req.IvSpDefense,
		IvSpeed:      req.IvSpeed,
		EvHp:         req.EvHp,
		EvAttack:     req.EvAttack,
		EvDefense:    req.EvDefense,
		EvSpAttack:   req.EvSpAttack,
		EvSpDefense:  req.EvSpDefense,
		EvSpeed:      req.EvSpeed,
		Move1Id:      req.Move1Id,
		Move2Id:      req.Move2Id,
		Move3Id:      req.Move3Id,
		Move4Id:      req.Move4Id,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, bp); err != nil {
		slog.ErrorContext(ctx, "failed to create box_pokemon", slog.Any("error", err))
		return nil, err
	}

	slog.InfoContext(ctx, "box_pokemon added", slog.String("box_pokemon_id", bp.BoxPokemonId.String()))
	return &bp, nil
}

func (s *BoxPokemonService) GetByBoxId(ctx context.Context, boxId uuid.UUID) ([]model.BoxPokemon, error) {
	bps, err := s.repo.FindByBoxId(ctx, boxId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get box_pokemon", slog.String("box_id", boxId.String()), slog.Any("error", err))
		return nil, err
	}
	return bps, nil
}

func (s *BoxPokemonService) Update(ctx context.Context, req dto.UpdateBoxPokemonRequest) (*model.BoxPokemon, error) {
	id, err := uuid.Parse(req.BoxPokemonId)
	if err != nil {
		return nil, fmt.Errorf("invalid box_pokemon_id: %w", err)
	}

	slog.InfoContext(ctx, "update box_pokemon", slog.String("box_pokemon_id", id.String()))

	if _, err := s.dex.GetAbility(ctx, &gen.GetAbilityRequest{AbilityId: int32(req.AbilityId)}); err != nil {
		slog.WarnContext(ctx, "ability not found", slog.Int("ability_id", req.AbilityId), slog.Any("error", err))
		return nil, fmt.Errorf("ability %d not found: %w", req.AbilityId, err)
	}

	for i, moveId := range []*int{req.Move1Id, req.Move2Id, req.Move3Id, req.Move4Id} {
		if moveId == nil {
			continue
		}
		if _, err := s.dex.GetMove(ctx, &gen.GetMoveRequest{MoveId: int32(*moveId)}); err != nil {
			slog.WarnContext(ctx, "move not found", slog.Int("slot", i+1), slog.Int("move_id", *moveId), slog.Any("error", err))
			return nil, fmt.Errorf("move %d not found: %w", *moveId, err)
		}
	}

	existing, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Nickname = req.Nickname
	existing.AbilityId = req.AbilityId
	existing.Nature = req.Nature
	existing.HeldItemId = req.HeldItemId
	existing.IvHp = req.IvHp
	existing.IvAttack = req.IvAttack
	existing.IvDefense = req.IvDefense
	existing.IvSpAttack = req.IvSpAttack
	existing.IvSpDefense = req.IvSpDefense
	existing.IvSpeed = req.IvSpeed
	existing.EvHp = req.EvHp
	existing.EvAttack = req.EvAttack
	existing.EvDefense = req.EvDefense
	existing.EvSpAttack = req.EvSpAttack
	existing.EvSpDefense = req.EvSpDefense
	existing.EvSpeed = req.EvSpeed
	existing.Move1Id = req.Move1Id
	existing.Move2Id = req.Move2Id
	existing.Move3Id = req.Move3Id
	existing.Move4Id = req.Move4Id
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, *existing); err != nil {
		slog.ErrorContext(ctx, "failed to update box_pokemon", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}

	slog.InfoContext(ctx, "box_pokemon updated", slog.String("id", id.String()))
	return existing, nil
}

func (s *BoxPokemonService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "failed to delete box_pokemon", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "box_pokemon deleted", slog.String("id", id.String()))
	return nil
}
