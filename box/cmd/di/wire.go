//go:build wireinject

package di

import (
	"pob/box/internal/handler"
	"pob/box/internal/shared"
	gen "pob/box/proto"

	"github.com/google/wire"
)

type Container struct {
	Box   *handler.BoxHandler
	Party *handler.PartyHandler
	Dex   *handler.DexHandler
	Grpc  *handler.GrpcHandler
}

func NewContainer(db *shared.DBClient, dex gen.DexServiceClient) (*Container, error) {
	wire.Build(
		BoxSet,
		BoxPokemonSet,
		PartySet,
		PartyPokemonSet,
		handler.NewDexHandler,
		handler.NewGrpcHandler,
		wire.Struct(new(Container), "*"),
	)
	return nil, nil
}
