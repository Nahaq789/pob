package other

import (
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/status"
)

// OtherStatusContext は OtherStatus.Resolve に渡す最小限の情報。
// EntryContext / PreDamageContext と同じ設計思想。
type OtherStatusContext struct {
	ActorId         string
	MoveId          string
	ActionSucceeded bool
	MoveType        ptype.Type       // ClearedOnElectricMove 用
	MainCondition   status.Condition // ClearedOnWakeUp 用（None なら状態異常なし）
}
