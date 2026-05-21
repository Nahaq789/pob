package repository

import "pob/dex/internal/shared"

type PokemonRepository struct {
	db *shared.DBClient
}

func NewPokemonRepository(db *shared.DBClient) *PokemonRepository {
	return &PokemonRepository{db: db}
}
