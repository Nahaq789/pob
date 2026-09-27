package damage

import (
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/ptype"
)

// Spec はダメージ計算に必要な immutable なパラメータ。
type Spec struct {
	ActorId    string
	MoveId     int
	Type       ptype.Type
	Power      int
	Category   move.DamageClass
	MustHit    bool
	CanCrit    bool
	TargetSelf bool
}
