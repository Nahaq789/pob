package player

import "pob/battle/internal/domain/pokemon"

type SwitchRequest struct {
	Outgoing      *pokemon.Pokemon
	Incoming      *pokemon.Pokemon
	IncomingIndex int
}

type MoveRequest struct {
	PlayerId string
	Pokemon  *pokemon.Pokemon
	MoveId   int
}

type ForfeitRequest bool
