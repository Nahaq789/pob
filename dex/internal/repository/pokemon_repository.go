package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pob/dex/internal/model"
	"pob/dex/internal/model/entity"
	"pob/dex/internal/shared"
	pkgredis "pob/pkg/redis"
)

const pokemonCacheTTL = 24 * time.Hour

type PokemonRepository struct {
	db    *shared.DBClient
	redis *pkgredis.RedisClient
}

func NewPokemonRepository(db *shared.DBClient, redis *pkgredis.RedisClient) *PokemonRepository {
	return &PokemonRepository{db: db, redis: redis}
}

func (p *PokemonRepository) cacheKey(id int) string {
	return fmt.Sprintf("pokemon:%d", id)
}

func (p *PokemonRepository) movesCacheKey(pokemonId int) string {
	return fmt.Sprintf("learnable_moves:%d", pokemonId)
}

func (p *PokemonRepository) FindById(ctx context.Context, id int) (*model.Pokemon, error) {
	key := p.cacheKey(id)

	if cached, err := p.redis.GetClient().Get(ctx, key).Bytes(); err == nil {
		var pokemon model.Pokemon
		if err := json.Unmarshal(cached, &pokemon); err == nil {
			return &pokemon, nil
		}
	}

	var pokemonEntity entity.Pokemon
	if err := p.db.GetClient().WithContext(ctx).First(&pokemonEntity, id).Error; err != nil {
		return nil, err
	}

	abilities, err := p.FindAbilitiesByPokemonId(ctx, id)
	if err != nil {
		return nil, err
	}

	pokemon := toModel(pokemonEntity, abilities)

	if b, err := json.Marshal(pokemon); err == nil {
		p.redis.GetClient().Set(ctx, key, b, pokemonCacheTTL)
	}

	return pokemon, nil
}

func toModel(e entity.Pokemon, abilities []entity.PokemonAbility) *model.Pokemon {
	var type2Id int
	if e.Type2Id != nil {
		type2Id = *e.Type2Id
	}

	slots := make([]model.AbilitySlot, len(abilities))
	for i, a := range abilities {
		slots[i] = model.AbilitySlot{
			AbilityId:   a.AbilityId,
			AbilityName: a.Ability.Name,
			IsHidden:    a.IsHidden,
		}
	}

	return &model.Pokemon{
		PokemonId:     e.Id,
		Name:          e.Name,
		Type1Id:       e.Type1Id,
		Type2Id:       type2Id,
		BaseHp:        e.BaseHp,
		BaseAttack:    e.BaseAttack,
		BaseDefense:   e.BaseDefense,
		BaseSpAttack:  e.BaseSpAttack,
		BaseSpDefense: e.BaseSpDefense,
		BaseSpeed:     e.BaseSpeed,
		Abilities:     slots,
	}
}

func (p *PokemonRepository) FindAbilitiesByPokemonId(ctx context.Context, pokemonId int) ([]entity.PokemonAbility, error) {
	var abilities []entity.PokemonAbility
	if err := p.db.GetClient().WithContext(ctx).
		Preload("Ability").
		Where("pokemon_id = ?", pokemonId).
		Order("slot ASC").
		Find(&abilities).Error; err != nil {
		return nil, err
	}
	return abilities, nil
}

func (p *PokemonRepository) FindMovesByPokemonId(ctx context.Context, pokemonId int) ([]entity.PokemonMove, error) {
	key := p.movesCacheKey(pokemonId)

	if cached, err := p.redis.GetClient().Get(ctx, key).Bytes(); err == nil {
		var moves []entity.PokemonMove
		if err := json.Unmarshal(cached, &moves); err == nil {
			return moves, nil
		}
	}

	var moves []entity.PokemonMove
	if err := p.db.GetClient().WithContext(ctx).
		Preload("Move").
		Where("pokemon_id = ?", pokemonId).
		Find(&moves).Error; err != nil {
		return nil, err
	}

	if b, err := json.Marshal(moves); err == nil {
		p.redis.GetClient().Set(ctx, key, b, pokemonCacheTTL)
	}

	return moves, nil
}
