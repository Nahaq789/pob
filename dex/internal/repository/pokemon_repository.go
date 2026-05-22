package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

func (p *PokemonRepository) FindById(ctx context.Context, id int) (*entity.Pokemon, error) {
	key := p.cacheKey(id)

	if cached, err := p.redis.GetClient().Get(ctx, key).Bytes(); err == nil {
		var pokemon entity.Pokemon
		if err := json.Unmarshal(cached, &pokemon); err == nil {
			return &pokemon, nil
		}
	}

	var pokemon entity.Pokemon
	if err := p.db.GetClient().WithContext(ctx).First(&pokemon, id).Error; err != nil {
		return nil, err
	}

	if b, err := json.Marshal(pokemon); err == nil {
		p.redis.GetClient().Set(ctx, key, b, pokemonCacheTTL)
	}

	return &pokemon, nil
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
	var moves []entity.PokemonMove
	if err := p.db.GetClient().WithContext(ctx).
		Preload("Move").
		Where("pokemon_id = ?", pokemonId).
		Find(&moves).Error; err != nil {
		return nil, err
	}
	return moves, nil
}
