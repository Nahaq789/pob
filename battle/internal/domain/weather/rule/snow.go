package rule

import (
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/phase"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/weather"
)

var _ phase.StatModHandler = (*SnowHandler)(nil)

// SnowHandler はゆき状態の能力値補正を行う基本ハンドラー。
// こおりタイプのポケモンが物理技を受けるとき、防御に 1.5 倍の補正を返す。特性により天気が打ち消されている場合は補正しない。
type SnowHandler struct{}

func NewSnowHandler() *SnowHandler { return &SnowHandler{} }

func (s *SnowHandler) StatMod(ctx phase.MoveResolveContext) damage.StatMod {
	actor := ctx.Battle.PlayerById(ctx.ActorId)
	defender := ctx.Battle.Opponent(actor).Active()

	w := ctx.Battle.EffectiveWeather()
	if w == nil || w.C() != weather.Snow {
		return damage.StatMod{}
	}
	if !defender.HasType(ptype.Ice) {
		return damage.StatMod{}
	}
	if ctx.Category != move.DamageClassPhysical {
		return damage.StatMod{}
	}
	return damage.StatMod{DefenderStat: damage.Multiply(1.5)}
}
