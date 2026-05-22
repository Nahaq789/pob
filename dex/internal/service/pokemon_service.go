package service

import (
	"context"
	"pob/dex/internal/repository"
)

type PokemonService struct {
	pokemonRepo *repository.PokemonRepository
}

func NewPokemonService(pokemonRepo *repository.PokemonRepository) *PokemonService {
	return &PokemonService{pokemonRepo: pokemonRepo}
}

func (p *PokemonService) GetPokemon(ctx context.Context, id int) error {
	return nil
}
