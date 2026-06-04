package handler

import (
	"context"
	"pob/box/internal/service"
	"pob/box/proto"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	proto.UnimplementedBoxServiceServer
	partyPokemonSvc *service.PartyPokemonService
	boxPokemonSvc   *service.BoxPokemonService
}

func NewGrpcHandler(partyPokemonSvc *service.PartyPokemonService, boxPokemonSvc *service.BoxPokemonService) *GrpcHandler {
	return &GrpcHandler{partyPokemonSvc: partyPokemonSvc, boxPokemonSvc: boxPokemonSvc}
}

func (h *GrpcHandler) GetParty(ctx context.Context, req *proto.GetPartyRequest) (*proto.PartyResponse, error) {
	partyId, err := uuid.Parse(req.PartyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid party_id")
	}

	partyPokemons, err := h.partyPokemonSvc.GetByPartyId(ctx, partyId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*proto.PartyPokemon, len(partyPokemons))
	eg, egCtx := errgroup.WithContext(ctx)
	for i, pp := range partyPokemons {
		eg.Go(func() error {
			bp, err := h.boxPokemonSvc.GetById(egCtx, pp.BoxPokemonId)
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}
			result[i] = &proto.PartyPokemon{
				BoxPokemonId: bp.BoxPokemonId.String(),
				Slot:         int32(pp.Slot),
				PokemonId:    int32(bp.PokemonId),
				Nickname:     strOrEmpty(bp.Nickname),
				Nature:       bp.Nature,
				Gender:       int32(bp.Gender),
				AbilityId:    int32(bp.AbilityId),
				Move1Id:      intOrZero(bp.Move1Id),
				Move2Id:      intOrZero(bp.Move2Id),
				Move3Id:      intOrZero(bp.Move3Id),
				Move4Id:      intOrZero(bp.Move4Id),
				IvHp:         int32(bp.IvHp),
				IvAttack:     int32(bp.IvAttack),
				IvDefense:    int32(bp.IvDefense),
				IvSpAttack:   int32(bp.IvSpAttack),
				IvSpDefense:  int32(bp.IvSpDefense),
				IvSpeed:      int32(bp.IvSpeed),
				EvHp:         int32(bp.EvHp),
				EvAttack:     int32(bp.EvAttack),
				EvDefense:    int32(bp.EvDefense),
				EvSpAttack:   int32(bp.EvSpAttack),
				EvSpDefense:  int32(bp.EvSpDefense),
				EvSpeed:      int32(bp.EvSpeed),
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.PartyResponse{Pokemon: result}, nil
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func intOrZero(v *int) int32 {
	if v == nil {
		return 0
	}
	return int32(*v)
}
