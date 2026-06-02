package di

import (
	"pob/box/internal/handler"
	"pob/box/internal/repository"
	"pob/box/internal/service"

	"github.com/google/wire"
)

var BoxSet = wire.NewSet(
	repository.NewBoxRepository,
	service.NewBoxService,
	handler.NewBoxHandler,
)
