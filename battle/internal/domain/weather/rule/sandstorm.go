package rule

import (
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/phase"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/weather"
)

var _ phase.StatModHandler = (*SandstormHandler)(nil)

// SandstormHandler はすなあらし状態の能力値補正を行う基本ハンドラー。
// いわタイプのポケモンが特殊技を受けるとき、特防に 1.5 倍の補正を返す。特性により天気が打ち消されている場合は補正しない。
type SandstormHandler struct{}

func NewSandstormHandler() *SandstormHandler { return &SandstormHandler{} }

func (s *SandstormHandler) StatMod(ctx phase.MoveResolveContext) damage.StatMod {
	actor := ctx.Battle.PlayerById(ctx.ActorId)
	defender := ctx.Battle.Opponent(actor).Active()

	w := ctx.Battle.EffectiveWeather()
	if w == nil || w.C() != weather.Sandstorm {
		return damage.StatMod{}
	}
	if !defender.HasType(ptype.Rock) {
		return damage.StatMod{}
	}
	if ctx.Category != move.DamageClassSpecial {
		return damage.StatMod{}
	}
	return damage.StatMod{DefenderStat: damage.Multiply(1.5)}
}
