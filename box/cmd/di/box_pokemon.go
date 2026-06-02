package di

import (
	"pob/box/internal/repository"
	"pob/box/internal/service"

	"github.com/google/wire"
)

var BoxPokemonSet = wire.NewSet(
	repository.NewBoxPokemonRepository,
	service.NewBoxPokemonService,
)
