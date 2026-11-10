package rule

import (
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/phase"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/weather"
)

var _ phase.DamageModHandler = (*SunnyHandler)(nil)

// SunnyHandler ははれ状態の天気補正を行う基本ハンドラー。
// ほのお技に 1.5 倍、みず技に 0.5 倍の補正を返す。特性により天気が打ち消されている場合は補正しない。
type SunnyHandler struct{}

func NewSunnyHandler() *SunnyHandler { return &SunnyHandler{} }

func (s *SunnyHandler) Mod(ctx phase.MoveResolveContext) damage.DamageMod {
	w := ctx.Battle.EffectiveWeather()
	// 晴れか判定
	if w == nil || w.C() != weather.Sunny {
		return damage.DamageMod{}
	}

	// ほのお技判定
	if ctx.Type == ptype.Fire {
		return damage.DamageMod{Weather: damage.Multiply(1.5)}
	}
	if ctx.Type == ptype.Water {
		return damage.DamageMod{Weather: damage.Multiply(0.5)}
	}
	return damage.DamageMod{}
}
