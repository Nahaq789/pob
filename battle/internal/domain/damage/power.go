package damage

const modBase = 4096 // 威力補正の基準値（等倍 = 4096）

type Power struct {
	base int   // 基礎威力
	mods []int // 威力補正リスト（4096 が等倍 / テクニシャン: 6144 / ちからのハチマキ: 4505 等）
}

// NewPower は Power を生成する。mods は可変長で複数の補正を渡せる。
func NewPower(base int, mods ...int) Power {
	return Power{base: base, mods: mods}
}

// final は最終威力を返す。
// 補正は 4096 を初期値として順次乗算し、各ステップで五捨五超入する。
// 最終威力 = roundHalfDown(base × 累積mod / 4096)
func (p Power) final() int {
	accumulated := modBase
	for _, mod := range p.mods {
		accumulated = roundHalfDown(float64(accumulated*mod) / modBase)
	}
	return roundHalfDown(float64(p.base*accumulated) / modBase)
}
