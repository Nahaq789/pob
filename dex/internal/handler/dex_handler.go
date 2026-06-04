package handler

import (
	"context"

	gen "pob/dex/proto"
	"pob/dex/internal/service"
)

type DexHandler struct {
	gen.UnimplementedDexServiceServer
	pokemon *service.PokemonService
	move    *service.MoveService
	ability *service.AbilityService
	item    *service.ItemService
}

func NewDexHandler(pokemon *service.PokemonService, move *service.MoveService, ability *service.AbilityService, item *service.ItemService) *DexHandler {
	return &DexHandler{pokemon: pokemon, move: move, ability: ability, item: item}
}

func (d *DexHandler) GetPokemon(ctx context.Context, r *gen.GetPokemonRequest) (*gen.PokemonResponse, error) {
	p, err := d.pokemon.GetPokemon(ctx, int(r.PokemonId))
	if err != nil {
		return nil, err
	}

	abilities := make([]*gen.AbilityInfo, len(p.Abilities))
	for i, a := range p.Abilities {
		abilities[i] = &gen.AbilityInfo{
			AbilityId:   int32(a.AbilityId),
			AbilityName: a.AbilityName,
			IsHidden:    a.IsHidden,
		}
	}

	return &gen.PokemonResponse{
		PokemonId:     int32(p.PokemonId),
		Name:          p.Name,
		Type1Id:       int32(p.Type1Id),
		Type2Id:       int32(p.Type2Id),
		BaseHp:        int32(p.BaseHp),
		BaseAttack:    int32(p.BaseAttack),
		BaseDefense:   int32(p.BaseDefense),
		BaseSpAttack:  int32(p.BaseSpAttack),
		BaseSpDefense: int32(p.BaseSpDefense),
		BaseSpeed:     int32(p.BaseSpeed),
		Abilities:     abilities,
	}, nil
}

func (d *DexHandler) GetLearnableMoves(ctx context.Context, r *gen.GetLearnableMovesRequest) (*gen.LearnableMovesResponse, error) {
	moves, err := d.pokemon.GetLearnableMoves(ctx, int(r.PokemonId))
	if err != nil {
		return nil, err
	}

	res := make([]*gen.MoveResponse, len(moves))
	for i, m := range moves {
		res[i] = &gen.MoveResponse{
			MoveId:      int32(m.MoveId),
			Name:        m.Name,
			TypeId:      int32(m.TypeId),
			DamageClass: m.DamageClass,
			Power:       int32(m.Power),
			Accuracy:    int32(m.Accuracy),
			Pp:          int32(m.Pp),
			Priority:    int32(m.Priority),
		}
	}

	return &gen.LearnableMovesResponse{Moves: res}, nil
}

func (d *DexHandler) GetMove(ctx context.Context, r *gen.GetMoveRequest) (*gen.MoveResponse, error) {
	m, err := d.move.GetMove(ctx, int(r.MoveId))
	if err != nil {
		return nil, err
	}

	return &gen.MoveResponse{
		MoveId:      int32(m.MoveId),
		Name:        m.Name,
		TypeId:      int32(m.TypeId),
		DamageClass: m.DamageClass,
		Power:       int32(m.Power),
		Accuracy:    int32(m.Accuracy),
		Pp:          int32(m.Pp),
		Priority:    int32(m.Priority),
	}, nil
}

func (d *DexHandler) GetAbility(ctx context.Context, r *gen.GetAbilityRequest) (*gen.AbilityResponse, error) {
	a, err := d.ability.GetAbility(ctx, int(r.AbilityId))
	if err != nil {
		return nil, err
	}

	return &gen.AbilityResponse{
		AbilityId:   int32(a.AbilityId),
		AbilityName: a.Name,
	}, nil
}

func (d *DexHandler) GetPokemonList(ctx context.Context, r *gen.GetPokemonListRequest) (*gen.PokemonListResponse, error) {
	list, err := d.pokemon.GetPokemonList(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*gen.PokemonList, len(list))
	for i, p := range list {
		items[i] = &gen.PokemonList{
			PokemonId: int32(p.PokemonId),
			Name:      p.Name,
		}
	}

	return &gen.PokemonListResponse{Pokemon: items}, nil
}

func (d *DexHandler) GetItem(ctx context.Context, r *gen.GetItemRequest) (*gen.ItemResponse, error) {
	item, err := d.item.GetItem(ctx, int(r.ItemId))
	if err != nil {
		return nil, err
	}

	return &gen.ItemResponse{
		ItemId:     int32(item.Id),
		Name:       item.Name,
		Category:   item.Category,
		FlavorText: item.FlavorText,
	}, nil
}
