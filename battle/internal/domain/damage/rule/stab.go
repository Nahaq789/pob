package rule

import (
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/phase"
)

var _ phase.DamageModHandler = (*StabHandler)(nil)

// StabHandler はタイプ一致補正を行う基本ハンドラー。
// 技のタイプが攻撃側のタイプのいずれかと一致する場合に 1.5 倍の補正を返す。
// てきおうりょく（2.0 倍）は特性ハンドラー側でタイプ一致を確認したうえで Override する。
type StabHandler struct{}

func NewStabHandler() *StabHandler { return &StabHandler{} }

func (s *StabHandler) Mod(ctx phase.MoveResolveContext) damage.DamageMod {
	attacker := ctx.Battle.PlayerById(ctx.ActorId).Active()
	if !attacker.HasType(ctx.Type) {
		return damage.DamageMod{}
	}

	return damage.DamageMod{Stab: damage.Multiply(1.5)}
}
