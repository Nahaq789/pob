// Package rule はダメージ計算の基本ルールと handler を提供する。
package rule

import (
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/phase"
)

var _ phase.DamageModHandler = (*CritHandler)(nil)

// CritHandler は急所判定を行う基本ハンドラー。
// 攻撃側の急所ランクに技の急所補正を加算した確率で判定し、急所時は 1.5 倍の補正を返す。
type CritHandler struct{}

func NewCritHandler() *CritHandler { return &CritHandler{} }

func (h *CritHandler) Mod(ctx phase.MoveResolveContext) damage.DamageMod {
	if !ctx.CanCrit {
		return damage.DamageMod{}
	}

	attacker := ctx.Battle.PlayerById(ctx.ActorId).Active()
	crit, err := attacker.RollCritical(ctx.CritBonus)
	if err != nil || !crit {
		return damage.DamageMod{}
	}

	return damage.DamageMod{Crit: damage.Override(1.5)}
}
