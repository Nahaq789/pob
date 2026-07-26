package phase

import (
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/pokemon"
)

type EntryContext struct {
	ActorId   string
	AbilityId int
	ItemId    int
	Battle    *battle.Battle
}

func NewEntryContext(actorId string, abilityId, itemId int, battle *battle.Battle) EntryContext {
	return EntryContext{
		ActorId:   actorId,
		AbilityId: abilityId,
		ItemId:    itemId,
		Battle:    battle,
	}
}

type ExitContext struct {
	ActorId   string
	AbilityId int
	ItemId    int
	Incoming  *pokemon.Pokemon
	Battle    *battle.Battle
}

func NewExitContext(actorId string, abilityId, itemId int, incoming *pokemon.Pokemon, battle *battle.Battle) ExitContext {
	return ExitContext{
		ActorId:   actorId,
		AbilityId: abilityId,
		ItemId:    itemId,
		Incoming:  incoming,
		Battle:    battle,
	}
}
