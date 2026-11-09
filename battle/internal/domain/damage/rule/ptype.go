package rule

import (
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/phase"
	"pob/battle/internal/domain/ptype"
)

var _ phase.DamageModHandler = (*PtypeHandler)(nil)

// PtypeHandler はタイプ相性補正を行う基本ハンドラー。
// 技のタイプと防御側のタイプから相性倍率を算出して返す。
type PtypeHandler struct{}

func NewPtypeHandler() *PtypeHandler { return &PtypeHandler{} }

func (p *PtypeHandler) Mod(ctx phase.MoveResolveContext) damage.DamageMod {
	attackerP := ctx.Battle.PlayerById(ctx.ActorId)
	defence := ctx.Battle.Opponent(attackerP).Active()

	v := ptype.Effectiveness(ctx.Spec.Type, defence.Types())
	return damage.DamageMod{TypeEff: damage.Multiply(v)}
}
