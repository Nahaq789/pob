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
