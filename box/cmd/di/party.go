package di

import (
	"pob/box/internal/handler"
	"pob/box/internal/repository"
	"pob/box/internal/service"

	"github.com/google/wire"
)

var PartySet = wire.NewSet(
	repository.NewPartyRepository,
	service.NewPartyService,
	handler.NewPartyHandler,
)
