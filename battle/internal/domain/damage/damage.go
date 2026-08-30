package damage

import "math"

const LEVEL = 50
const PRECISION = 1e6

// DamageInput は第五世代以降のダメージ計算式に必要なパラメータを保持する。
// 各補正倍率は呼び出し側で事前に解決して渡す。補正なしの場合は 1.0 を設定する。
type DamageInput struct {
	Power   int
	Attack  int
	Defense int

	CritMod float64 // 急所補正（通常: 1.0 / 急所: 1.5）
	Random  int     // 乱数補正（85〜100）
	StabMod float64 // タイプ一致補正（不一致: 1.0 / 一致: 1.5 / てきおうりょく: 2.0）
	TypeEff float64 // タイプ相性補正（全タイプ倍率の積）
	BurnMod float64 // やけど補正（物理+やけど: 0.5 / それ以外: 1.0）
	Weather float64 // 天気補正

	// M補正（各補正を math.Round で順次適用）
	Wall        float64 // リフレクター・ひかりのかべ
	Other       float64
	Neuroforce  float64 // ブレインフォース（効果抜群時: 1.25）
	Sniper      float64 // スナイパー（急所時: 1.5）
	TintedLens  float64 // いろめがね（効果いまひとつ時: 2.0）
	Fluffy      float64 // もふもふ（ほのお技被弾時: 2.0）
	Mhalf       float64 // ファントムガード・マルチスケイル等（半減: 0.5）
	Mfilter     float64 // フィルター・ハードロック等（効果抜群時: 0.75）
	Mtwice      float64 // メガランチャー等（対象技: 1.5）
	FriendGuard float64 // フレンドガード（常に 1.0 / ダブルバトル専用のため非実装）
	ExpertBelt  float64 // たつじんのおび（効果抜群時: 1.2）
	Metronome   float64 // メトロノーム（連続使用回数に応じて増加）
	LifeOrb     float64 // いのちのたま（1.3）
	HalfBerry   float64 // 半減の実（0.5）
}

func NewDamageInput(
	power, attack, defense int,
	critMod float64,
	random int,
	stabMod float64,
	typeEff float64,
	burnMod float64,
	weather, wall, other float64,
	neuroforce, sniper, tintedLens, fluffy float64,
	mhalf, mfilter, mtwice float64,
	expertBelt, metronome, lifeOrb, halfBerry float64,

) *DamageInput {
	return &DamageInput{
		Power:       power,
		Attack:      attack,
		Defense:     defense,
		CritMod:     critMod,
		Random:      random,
		StabMod:     stabMod,
		TypeEff:     typeEff,
		BurnMod:     burnMod,
		Weather:     weather,
		Wall:        wall,
		Other:       other,
		Neuroforce:  neuroforce,
		Sniper:      sniper,
		TintedLens:  tintedLens,
		Fluffy:      fluffy,
		Mhalf:       mhalf,
		Mfilter:     mfilter,
		Mtwice:      mtwice,
		FriendGuard: 1.0,
		ExpertBelt:  expertBelt,
		Metronome:   metronome,
		LifeOrb:     lifeOrb,
		HalfBerry:   halfBerry,
	}
}

func (d *DamageInput) CalcDamage() int {
	base := float64(((LEVEL*2/5+2)*d.Power*d.Attack/d.Defense)/50 + 2)

	// 範囲補正
	// ダブルバトルは想定していないので*1で計算
	damage := d.roundHalfDown(base * 1.0)

	// おやこあい補正
	// メガガルーラの専用特性で、メガシンカは実装しない予定なので *1で計算
	damage = d.roundHalfDown(float64(damage) * 1.0)

	// 天気補正
	damage = d.roundHalfDown(float64(damage) * d.Weather)

	// 急所補正
	damage = d.roundHalfDown(float64(damage) * d.CritMod)

	// 乱数補正
	damage = damage * d.Random / 100

	// STAB補正
	damage = d.roundHalfDown(float64(damage) * d.StabMod)

	// 相性補正
	damage = d.roundHalfDown(float64(damage) * d.TypeEff)

	// やけど補正
	damage = d.roundHalfDown(float64(damage) * d.BurnMod)

	// M
	damage = d.roundHalfDown(d.calcM(float64(damage)))

	// Mprotect
	// ダイマックスは実装予定にないので、1.0で計算する
	damage = d.roundHalfDown(float64(damage) * 1.0)

	return damage
}

func (d *DamageInput) calcM(damage float64) float64 {
	m := damage

	// 壁補正
	m = math.Round(m * d.Wall)

	// ブレインフォース補正
	m = math.Round(m * d.Neuroforce)

	// スナイパー補正
	m = math.Round(m * d.Sniper)

	// いろめがね補正
	m = math.Round(m * d.TintedLens)

	// もふもふ(ほのおタイプ)補正
	m = math.Round(m * d.Fluffy)

	// Mhalf
	m = math.Round(m * d.Mhalf)

	// Mfilter
	m = math.Round(m * d.Mfilter)

	// たつじんのおび補正
	m = math.Round(m * d.ExpertBelt)

	// メトロノーム補正
	m = math.Round(m * d.Metronome)

	// いのちのたま補正
	m = math.Round(m * d.LifeOrb)

	// 半減の実補正
	m = math.Round(m * d.HalfBerry)

	// Mtwice
	m = math.Round(m * d.Mtwice)

	return m
}

func (d *DamageInput) roundHalfDown(v float64) int {
	scaled := math.Round(v * PRECISION)
	intPart := math.Floor(scaled / PRECISION)
	frac := scaled - (intPart * PRECISION)
	half := PRECISION / 2

	switch {
	case frac == half:
		return int(intPart)
	case frac > half:
		return int(intPart) + 1
	default:
		return int(intPart)
	}
}
