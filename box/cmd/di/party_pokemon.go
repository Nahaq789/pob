package di

import (
	"pob/box/internal/repository"
	"pob/box/internal/service"

	"github.com/google/wire"
)

var PartyPokemonSet = wire.NewSet(
	repository.NewPartyPokemonRepository,
	service.NewPartyPokemonService,
)
