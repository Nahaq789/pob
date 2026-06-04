package model

import (
	"time"

	"github.com/google/uuid"
)

type BoxPokemon struct {
	BoxPokemonId uuid.UUID
	BoxId        uuid.UUID
	PokemonId    int
	Nickname     *string
	AbilityId    int
	Nature       string
	Gender       Gender
	HeldItemId   *int
	IvHp         int
	IvAttack     int
	IvDefense    int
	IvSpAttack   int
	IvSpDefense  int
	IvSpeed      int
	EvHp         int
	EvAttack     int
	EvDefense    int
	EvSpAttack   int
	EvSpDefense  int
	EvSpeed      int
	Move1Id      *int
	Move2Id      *int
	Move3Id      *int
	Move4Id      *int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewBoxPokemon(
	boxId uuid.UUID,
	pokemonId int,
	abilityId int,
	nature string,
	gender Gender,
) BoxPokemon {
	now := time.Now()
	return BoxPokemon{
		BoxPokemonId: uuid.New(),
		BoxId:        boxId,
		PokemonId:    pokemonId,
		AbilityId:    abilityId,
		Nature:       nature,
		Gender:       gender,
		IvHp:         31,
		IvAttack:     31,
		IvDefense:    31,
		IvSpAttack:   31,
		IvSpDefense:  31,
		IvSpeed:      31,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func FromBoxPokemon(
	id, boxId uuid.UUID,
	pokemonId int,
	nickname *string,
	abilityId int,
	nature string,
	gender Gender,
	heldItemId *int,
	ivHp, ivAttack, ivDefense, ivSpAttack, ivSpDefense, ivSpeed int,
	evHp, evAttack, evDefense, evSpAttack, evSpDefense, evSpeed int,
	move1Id, move2Id, move3Id, move4Id *int,
	createdAt, updatedAt time.Time,
) BoxPokemon {
	return BoxPokemon{
		BoxPokemonId: id,
		BoxId:        boxId,
		PokemonId:    pokemonId,
		Nickname:     nickname,
		AbilityId:    abilityId,
		Nature:       nature,
		Gender:       gender,
		HeldItemId:   heldItemId,
		IvHp:         ivHp,
		IvAttack:     ivAttack,
		IvDefense:    ivDefense,
		IvSpAttack:   ivSpAttack,
		IvSpDefense:  ivSpDefense,
		IvSpeed:      ivSpeed,
		EvHp:         evHp,
		EvAttack:     evAttack,
		EvDefense:    evDefense,
		EvSpAttack:   evSpAttack,
		EvSpDefense:  evSpDefense,
		EvSpeed:      evSpeed,
		Move1Id:      move1Id,
		Move2Id:      move2Id,
		Move3Id:      move3Id,
		Move4Id:      move4Id,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}
