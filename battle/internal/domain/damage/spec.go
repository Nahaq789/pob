package damage

import (
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/ptype"
)

// Spec はダメージ計算に必要な immutable なパラメータ。
// Battle を持たないため damage パッケージ内で完結し、循環参照を避けられる。
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
