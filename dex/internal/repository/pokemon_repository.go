package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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

	cached, err := p.redis.GetClient().Get(ctx, key).Bytes()
	if err != nil {
		slog.WarnContext(ctx, "cache miss", slog.String("key", key), slog.Any("error", err))
	} else {
		var pokemon model.Pokemon
		if err := json.Unmarshal(cached, &pokemon); err != nil {
			slog.WarnContext(ctx, "cache unmarshal error", slog.String("key", key), slog.Any("error", err))
		} else {
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
		if err := p.redis.GetClient().Set(ctx, key, b, pokemonCacheTTL).Err(); err != nil {
			slog.WarnContext(ctx, "cache set error", slog.String("key", key), slog.Any("error", err))
		}
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

func (p *PokemonRepository) FindMovesByPokemonId(ctx context.Context, pokemonId int) ([]model.Move, error) {
	key := p.movesCacheKey(pokemonId)

	cached, err := p.redis.GetClient().Get(ctx, key).Bytes()
	if err != nil {
		slog.WarnContext(ctx, "cache miss", slog.String("key", key), slog.Any("error", err))
	} else {
		var moves []model.Move
		if err := json.Unmarshal(cached, &moves); err != nil {
			slog.WarnContext(ctx, "cache unmarshal error", slog.String("key", key), slog.Any("error", err))
		} else {
			return moves, nil
		}
	}

	var pokemonMoves []entity.PokemonMove
	if err := p.db.GetClient().WithContext(ctx).
		Preload("Move").
		Where("pokemon_id = ?", pokemonId).
		Find(&pokemonMoves).Error; err != nil {
		return nil, err
	}

	moves := make([]model.Move, len(pokemonMoves))
	for i, m := range pokemonMoves {
		var power, accuracy int
		if m.Move.Power != nil {
			power = *m.Move.Power
		}
		if m.Move.Accuracy != nil {
			accuracy = *m.Move.Accuracy
		}
		moves[i] = model.Move{
			MoveId:      m.Move.Id,
			Name:        m.Move.Name,
			TypeId:      m.Move.TypeId,
			DamageClass: m.Move.DamageClass,
			Power:       power,
			Accuracy:    accuracy,
			Pp:          m.Move.Pp,
			Priority:    m.Move.Priority,
		}
	}

	if b, err := json.Marshal(moves); err == nil {
		if err := p.redis.GetClient().Set(ctx, key, b, pokemonCacheTTL).Err(); err != nil {
			slog.WarnContext(ctx, "cache set error", slog.String("key", key), slog.Any("error", err))
		}
	}

	return moves, nil
}

func (p *PokemonRepository) FindAll(ctx context.Context) ([]model.PokemonList, error) {
	const key = "pokemon_list"

	cached, err := p.redis.GetClient().Get(ctx, key).Bytes()
	if err != nil {
		slog.WarnContext(ctx, "cache miss", slog.String("key", key), slog.Any("error", err))
	} else {
		var list []model.PokemonList
		if err := json.Unmarshal(cached, &list); err != nil {
			slog.WarnContext(ctx, "cache unmarshal error", slog.String("key", key), slog.Any("error", err))
		} else {
			return list, nil
		}
	}

	var entities []entity.Pokemon
	if err := p.db.GetClient().WithContext(ctx).Select("id", "name").Find(&entities).Error; err != nil {
		return nil, err
	}

	list := make([]model.PokemonList, len(entities))
	for i, e := range entities {
		list[i] = model.PokemonList{
			PokemonId: e.Id,
			Name:      e.Name,
		}
	}

	if b, err := json.Marshal(list); err == nil {
		if err := p.redis.GetClient().Set(ctx, key, b, pokemonCacheTTL).Err(); err != nil {
			slog.WarnContext(ctx, "cache set error", slog.String("key", key), slog.Any("error", err))
		}
	}

	return list, nil
}
