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
	apply := func(b, n *ModValue) *ModValue {
		if n == nil {
			return b
		}
		if b == nil || n.Op == ModOpOverride {
			return n
		}
		return Multiply(b.Value * n.Value)
	}
	return DamageMod{
		Crit:       apply(base.Crit, incoming.Crit),
		Stab:       apply(base.Stab, incoming.Stab),
		TypeEff:    apply(base.TypeEff, incoming.TypeEff),
		Burn:       apply(base.Burn, incoming.Burn),
		Weather:    apply(base.Weather, incoming.Weather),
		Wall:       apply(base.Wall, incoming.Wall),
		Neuroforce: apply(base.Neuroforce, incoming.Neuroforce),
		Sniper:     apply(base.Sniper, incoming.Sniper),
		TintedLens: apply(base.TintedLens, incoming.TintedLens),
		Fluffy:     apply(base.Fluffy, incoming.Fluffy),
		MHalf:      apply(base.MHalf, incoming.MHalf),
		MFilter:    apply(base.MFilter, incoming.MFilter),
		MTwice:     apply(base.MTwice, incoming.MTwice),
		ExpertBelt: apply(base.ExpertBelt, incoming.ExpertBelt),
		Metronome:  apply(base.Metronome, incoming.Metronome),
		LifeOrb:    apply(base.LifeOrb, incoming.LifeOrb),
		HalfBerry:  apply(base.HalfBerry, incoming.HalfBerry),
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
