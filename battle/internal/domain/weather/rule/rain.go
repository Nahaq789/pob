package rule

import (
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/phase"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/weather"
)

var _ phase.DamageModHandler = (*RainHandler)(nil)

// RainHandler はあめ状態の天気補正を行う基本ハンドラー。
// みず技に 1.5 倍、ほのお技に 0.5 倍の補正を返す。特性により天気が打ち消されている場合は補正しない。
type RainHandler struct{}

func NewRainHandler() *RainHandler { return &RainHandler{} }

func (s *RainHandler) Mod(ctx phase.MoveResolveContext) damage.DamageMod {
	w := ctx.Battle.EffectiveWeather()
	// 雨か判定
	if w == nil || w.C() != weather.Rain {
		return damage.DamageMod{}
	}

	// みず技判定
	if ctx.Type == ptype.Water {
		return damage.DamageMod{Weather: damage.Multiply(1.5)}
	}
	if ctx.Type == ptype.Fire {
		return damage.DamageMod{Weather: damage.Multiply(0.5)}
	}
	return damage.DamageMod{}
}
