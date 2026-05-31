package model

import "github.com/google/uuid"

type PartyPokemon struct {
	PartyPokemonId uuid.UUID
	PartyId        uuid.UUID
	BoxPokemonId   uuid.UUID
	Slot           int
}

func NewPartyPokemon(partyId, boxPokemonId uuid.UUID, slot int) PartyPokemon {
	return PartyPokemon{
		PartyPokemonId: uuid.New(),
		PartyId:        partyId,
		BoxPokemonId:   boxPokemonId,
		Slot:           slot,
	}
}

func FromPartyPokemon(id, partyId, boxPokemonId uuid.UUID, slot int) PartyPokemon {
	return PartyPokemon{
		PartyPokemonId: id,
		PartyId:        partyId,
		BoxPokemonId:   boxPokemonId,
		Slot:           slot,
	}
}
