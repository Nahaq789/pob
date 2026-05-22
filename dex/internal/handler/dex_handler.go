package handler

import (
	"context"
	gen "pob/dex/gen/proto"
	"pob/dex/internal/service"
)

type DexHandler struct {
	gen.UnimplementedDexServiceServer
	pokemon *service.PokemonService
	move    *service.MoveService
	ability *service.AbilityService
}

func NewDexHandler(pokemon *service.PokemonService, move *service.MoveService, ability *service.AbilityService) *DexHandler {
	return &DexHandler{pokemon: pokemon, move: move, ability: ability}
}

func (d *DexHandler) GetPokemon(ctx context.Context, r *gen.GetPokemonRequest) (*gen.PokemonResponse, error) {
	return nil, nil
}

func (d *DexHandler) GetLearnableMoves(ctx context.Context, r *gen.GetLearnableMovesRequest) (*gen.LearnableMovesResponse, error) {
	return nil, nil
}

func (d *DexHandler) GetMove(ctx context.Context, r *gen.GetMoveRequest) (*gen.MoveResponse, error) {
	return nil, nil
}

func (d *DexHandler) GetAbility(ctx context.Context, r *gen.GetAbilityRequest) (*gen.AbilityResponse, error) {
	return nil, nil
}
