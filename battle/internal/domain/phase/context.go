package phase

import "pob/battle/internal/domain/battle"

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
