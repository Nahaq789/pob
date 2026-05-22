package service

import (
	"context"

	"pob/dex/internal/model"
	"pob/dex/internal/repository"
)

type PokemonService struct {
	pokemonRepo *repository.PokemonRepository
}

func NewPokemonService(pokemonRepo *repository.PokemonRepository) *PokemonService {
	return &PokemonService{pokemonRepo: pokemonRepo}
}

func (p *PokemonService) GetPokemon(ctx context.Context, id int) (*model.Pokemon, error) {
	return p.pokemonRepo.FindById(ctx, id)
}

func (p *PokemonService) GetLearnableMoves(ctx context.Context, pokemonId int) ([]model.Move, error) {
	return p.pokemonRepo.FindMovesByPokemonId(ctx, pokemonId)
}
