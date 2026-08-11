package other

import "pob/battle/internal/domain/battle"

// OtherStatusContext は OtherStatus.Resolve に渡す最小限の情報。
// EntryContext / PreDamageContext と同じ設計思想。
type OtherStatusContext struct {
	ActorId         string
	MoveId          string
	Battle          *battle.Battle
	ActionSucceeded bool // そのポケモンの行動がpost_damageまで完走できたか(pre_damage/damageでhaltした場合false)
}
