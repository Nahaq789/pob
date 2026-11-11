package damage

type ModOp int

const (
	ModOpOverride ModOp = iota // 上書き
	ModOpMultiply              // 乗算
)

type ModValue struct {
	Value float64
	Op    ModOp
}

// Override は上書き補正値を生成するヘルパー。
func Override(v float64) *ModValue { return &ModValue{Value: v, Op: ModOpOverride} }

// Multiply は乗算補正値を生成するヘルパー。
func Multiply(v float64) *ModValue { return &ModValue{Value: v, Op: ModOpMultiply} }

// mergeValue は incoming を base に適用した補正値を返す。
// incoming が nil なら base、base が nil または incoming が Override なら incoming、それ以外は乗算した値を返す。
func mergeValue(base, incoming *ModValue) *ModValue {
	if incoming == nil {
		return base
	}
	if base == nil || incoming.Op == ModOpOverride {
		return incoming
	}
	return Multiply(base.Value * incoming.Value)
}

// DamageMod はダメージ計算ハンドラーが返す補正値の集合。
// nil フィールドはそのハンドラーが関与しないことを示す。
type DamageMod struct {
	Crit       *ModValue
	Stab       *ModValue
	TypeEff    *ModValue
	Burn       *ModValue
	Weather    *ModValue
	Wall       *ModValue
	Neuroforce *ModValue
	Sniper     *ModValue
	TintedLens *ModValue
	Fluffy     *ModValue
	MHalf      *ModValue
	MFilter    *ModValue
	MTwice     *ModValue
	ExpertBelt *ModValue
	Metronome  *ModValue
	LifeOrb    *ModValue
	HalfBerry  *ModValue
}

// Merge は incoming の非 nil フィールドを base に適用して返す。
// Override は上書き、Multiply は乗算（base が nil なら incoming をそのまま使用）。
func Merge(base, incoming DamageMod) DamageMod {
	return DamageMod{
		Crit:       mergeValue(base.Crit, incoming.Crit),
		Stab:       mergeValue(base.Stab, incoming.Stab),
		TypeEff:    mergeValue(base.TypeEff, incoming.TypeEff),
		Burn:       mergeValue(base.Burn, incoming.Burn),
		Weather:    mergeValue(base.Weather, incoming.Weather),
		Wall:       mergeValue(base.Wall, incoming.Wall),
		Neuroforce: mergeValue(base.Neuroforce, incoming.Neuroforce),
		Sniper:     mergeValue(base.Sniper, incoming.Sniper),
		TintedLens: mergeValue(base.TintedLens, incoming.TintedLens),
		Fluffy:     mergeValue(base.Fluffy, incoming.Fluffy),
		MHalf:      mergeValue(base.MHalf, incoming.MHalf),
		MFilter:    mergeValue(base.MFilter, incoming.MFilter),
		MTwice:     mergeValue(base.MTwice, incoming.MTwice),
		ExpertBelt: mergeValue(base.ExpertBelt, incoming.ExpertBelt),
		Metronome:  mergeValue(base.Metronome, incoming.Metronome),
		LifeOrb:    mergeValue(base.LifeOrb, incoming.LifeOrb),
		HalfBerry:  mergeValue(base.HalfBerry, incoming.HalfBerry),
	}
}

// ToOptions は DamageMod を DamageInput のオプション列に変換する。
// nil フィールドはスキップされ、DamageInput のデフォルト値（1.0）が使われる。
func (m DamageMod) ToOptions() []DamageInputOption {
	var opts []DamageInputOption
	if m.Crit != nil {
		opts = append(opts, WithCrit(m.Crit.Value))
	}
	if m.Stab != nil {
		opts = append(opts, WithStab(m.Stab.Value))
	}
	if m.TypeEff != nil {
		opts = append(opts, WithTypeEff(m.TypeEff.Value))
	}
	if m.Burn != nil {
		opts = append(opts, WithBurn(m.Burn.Value))
	}
	if m.Weather != nil {
		opts = append(opts, WithWeather(m.Weather.Value))
	}
	if m.Wall != nil {
		opts = append(opts, WithWall(m.Wall.Value))
	}
	if m.Neuroforce != nil {
		opts = append(opts, WithNeuroforce(m.Neuroforce.Value))
	}
	if m.Sniper != nil {
		opts = append(opts, WithSniper(m.Sniper.Value))
	}
	if m.TintedLens != nil {
		opts = append(opts, WithTintedLens(m.TintedLens.Value))
	}
	if m.Fluffy != nil {
		opts = append(opts, WithFluffy(m.Fluffy.Value))
	}
	if m.MHalf != nil {
		opts = append(opts, WithMHalf(m.MHalf.Value))
	}
	if m.MFilter != nil {
		opts = append(opts, WithMFilter(m.MFilter.Value))
	}
	if m.MTwice != nil {
		opts = append(opts, WithMTwice(m.MTwice.Value))
	}
	if m.ExpertBelt != nil {
		opts = append(opts, WithExpertBelt(m.ExpertBelt.Value))
	}
	if m.Metronome != nil {
		opts = append(opts, WithMetronome(m.Metronome.Value))
	}
	if m.LifeOrb != nil {
		opts = append(opts, WithLifeOrb(m.LifeOrb.Value))
	}
	if m.HalfBerry != nil {
		opts = append(opts, WithHalfBerry(m.HalfBerry.Value))
	}
	return opts
}
