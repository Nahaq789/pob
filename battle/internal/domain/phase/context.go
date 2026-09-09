package phase

import (
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/pokemon"
	"pob/battle/internal/domain/ptype"
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

type PreDamageContext struct {
	ActorId string
	MoveId  int
	Battle  *battle.Battle
}

func NewPreDamageContext(actorId string, moveId int, battle *battle.Battle) PreDamageContext {
	return PreDamageContext{
		ActorId: actorId,
		MoveId:  moveId,
		Battle:  battle,
	}
}

type DamageContext struct {
	Battle     *battle.Battle
	ActorId    string
	MoveId     int        // 通常時のみ有効。自傷時は無視
	Type       ptype.Type // 自傷時はptype.None
	Power      int
	Category   move.DamageClass // Physical / Special
	MustHit    bool             // true なら命中判定スキップ
	CanCrit    bool             // false なら急所判定スキップ
	TargetSelf bool             // true ならActor自身がTarget
	IsCrit     bool             // CritHandler が解決後に設定。スナイパー等の後続ハンドラーが参照する
}
