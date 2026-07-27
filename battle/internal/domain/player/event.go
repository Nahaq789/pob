package player

import "pob/battle/internal/domain/pokemon"

type SwitchRequest struct {
	Outgoing *pokemon.Pokemon
	Incoming *pokemon.Pokemon
}
