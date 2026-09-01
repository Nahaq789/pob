package damage

import "math"

const LEVEL = 50
const PRECISION = 1e6

// DamageInput は第五世代以降のダメージ計算式に必要なパラメータを保持する。
// 各補正倍率は呼び出し側で事前に解決して渡す。補正なしの場合は 1.0 を設定する。
type DamageInput struct {
	power   Power
	attack  int
	defense int

	critMod float64 // 急所補正（通常: 1.0 / 急所: 1.5）
	random  int     // 乱数補正（85〜100）
	stabMod float64 // タイプ一致補正（不一致: 1.0 / 一致: 1.5 / てきおうりょく: 2.0）
	typeEff float64 // タイプ相性補正（全タイプ倍率の積）
	burnMod float64 // やけど補正（物理+やけど: 0.5 / それ以外: 1.0）
	weather float64 // 天気補正

	// M補正（各補正を math.Round で順次適用）
	wall        float64 // リフレクター・ひかりのかべ（シングル: 0.5 / 急所時は呼び出し側で 1.0 を渡す）
	neuroforce  float64 // ブレインフォース（効果抜群時: 1.25）
	sniper      float64 // スナイパー（急所時: 1.5）
	tintedLens  float64 // いろめがね（効果いまひとつ時: 2.0）
	fluffy      float64 // もふもふ（ほのお技被弾時: 2.0）
	mhalf       float64 // ファントムガード・マルチスケイル等（半減: 0.5）
	mfilter     float64 // フィルター・ハードロック・プリズムアーマー等（効果抜群時: 0.75）
	mtwice      float64 // 状態依存2倍補正（あなをほるへのじしん等）
	friendGuard float64 // フレンドガード（常に 1.0 / ダブルバトル専用のため非実装）
	expertBelt  float64 // たつじんのおび（効果抜群時: 1.2）
	metronome   float64 // メトロノーム（連続使用回数に応じて増加）
	lifeOrb     float64 // いのちのたま（1.3）
	halfBerry   float64 // 半減の実（効果抜群被弾時: 0.5）
}

func NewDamageInput(
	power Power,
	attack, defense int,
	critMod float64,
	random int,
	stabMod float64,
	typeEff float64,
	burnMod float64,
	weather, wall float64,
	neuroforce, sniper, tintedLens, fluffy float64,
	mhalf, mfilter, mtwice float64,
	expertBelt, metronome, lifeOrb, halfBerry float64,
) *DamageInput {
	return &DamageInput{
		power:       power,
		attack:      attack,
		defense:     defense,
		critMod:     critMod,
		random:      random,
		stabMod:     stabMod,
		typeEff:     typeEff,
		burnMod:     burnMod,
		weather:     weather,
		wall:        wall,
		neuroforce:  neuroforce,
		sniper:      sniper,
		tintedLens:  tintedLens,
		fluffy:      fluffy,
		mhalf:       mhalf,
		mfilter:     mfilter,
		mtwice:      mtwice,
		friendGuard: 1.0,
		expertBelt:  expertBelt,
		metronome:   metronome,
		lifeOrb:     lifeOrb,
		halfBerry:   halfBerry,
	}
}

func (d *DamageInput) CalcDamage() int {
	finalPower := d.power.final()
	base := float64(((LEVEL*2/5+2)*finalPower*d.attack/d.defense)/50 + 2)

	// 範囲補正
	// ダブルバトルは想定していないので*1で計算
	damage := roundHalfDown(base * 1.0)

	// おやこあい補正
	// メガガルーラの専用特性で、メガシンカは実装しない予定なので *1で計算
	damage = roundHalfDown(float64(damage) * 1.0)

	// 天気補正
	damage = roundHalfDown(float64(damage) * d.weather)

	// 急所補正
	damage = roundHalfDown(float64(damage) * d.critMod)

	// 乱数補正
	damage = damage * d.random / 100

	// STAB補正
	damage = roundHalfDown(float64(damage) * d.stabMod)

	// 相性補正（切り捨て）
	damage = int(float64(damage) * d.typeEff)

	// やけど補正
	damage = roundHalfDown(float64(damage) * d.burnMod)

	// M
	damage = roundHalfDown(d.calcM(float64(damage)))

	// Mprotect
	// ダイマックスは実装予定にないので、1.0で計算する
	damage = roundHalfDown(float64(damage) * 1.0)

	return damage
}

func (d *DamageInput) calcM(damage float64) float64 {
	m := damage

	// 壁補正
	m = math.Round(m * d.wall)

	// ブレインフォース補正
	m = math.Round(m * d.neuroforce)

	// スナイパー補正
	m = math.Round(m * d.sniper)

	// いろめがね補正
	m = math.Round(m * d.tintedLens)

	// もふもふ(ほのおタイプ)補正
	m = math.Round(m * d.fluffy)

	// Mhalf
	m = math.Round(m * d.mhalf)

	// Mfilter
	m = math.Round(m * d.mfilter)

	// たつじんのおび補正
	m = math.Round(m * d.expertBelt)

	// メトロノーム補正
	m = math.Round(m * d.metronome)

	// いのちのたま補正
	m = math.Round(m * d.lifeOrb)

	// 半減の実補正
	m = math.Round(m * d.halfBerry)

	// Mtwice
	m = math.Round(m * d.mtwice)

	return m
}

func roundHalfDown(v float64) int {
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
