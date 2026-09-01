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

type DamageInputOption func(*DamageInput)

// WithCrit は急所補正（通常: 1.0 / 急所: 1.5）を設定する。
func WithCrit(mod float64) DamageInputOption { return func(d *DamageInput) { d.critMod = mod } }

// WithStab はタイプ一致補正（不一致: 1.0 / 一致: 1.5 / てきおうりょく: 2.0）を設定する。
func WithStab(mod float64) DamageInputOption { return func(d *DamageInput) { d.stabMod = mod } }

// WithTypeEff はタイプ相性補正（全タイプ倍率の積）を設定する。
func WithTypeEff(mod float64) DamageInputOption { return func(d *DamageInput) { d.typeEff = mod } }

// WithBurn はやけど補正（物理+やけど: 0.5 / それ以外: 1.0）を設定する。
func WithBurn(mod float64) DamageInputOption { return func(d *DamageInput) { d.burnMod = mod } }

// WithWeather は天気補正を設定する。
func WithWeather(mod float64) DamageInputOption { return func(d *DamageInput) { d.weather = mod } }

// WithWall はリフレクター・ひかりのかべ（シングル: 0.5 / 急所時は呼び出し側で 1.0 を渡す）を設定する。
func WithWall(mod float64) DamageInputOption { return func(d *DamageInput) { d.wall = mod } }

// WithNeuroforce はブレインフォース（効果抜群時: 1.25）を設定する。
func WithNeuroforce(mod float64) DamageInputOption {
	return func(d *DamageInput) { d.neuroforce = mod }
}

// WithSniper はスナイパー（急所時: 1.5）を設定する。
func WithSniper(mod float64) DamageInputOption { return func(d *DamageInput) { d.sniper = mod } }

// WithTintedLens はいろめがね（効果いまひとつ時: 2.0）を設定する。
func WithTintedLens(mod float64) DamageInputOption {
	return func(d *DamageInput) { d.tintedLens = mod }
}

// WithFluffy はもふもふ（ほのお技被弾時: 2.0）を設定する。
func WithFluffy(mod float64) DamageInputOption { return func(d *DamageInput) { d.fluffy = mod } }

// WithMHalf はファントムガード・マルチスケイル等（半減: 0.5）を設定する。
func WithMHalf(mod float64) DamageInputOption { return func(d *DamageInput) { d.mhalf = mod } }

// WithMFilter はフィルター・ハードロック・プリズムアーマー等（効果抜群時: 0.75）を設定する。
func WithMFilter(mod float64) DamageInputOption { return func(d *DamageInput) { d.mfilter = mod } }

// WithMTwice は状態依存2倍補正（あなをほるへのじしん等）を設定する。
func WithMTwice(mod float64) DamageInputOption { return func(d *DamageInput) { d.mtwice = mod } }

// WithExpertBelt はたつじんのおび（効果抜群時: 1.2）を設定する。
func WithExpertBelt(mod float64) DamageInputOption {
	return func(d *DamageInput) { d.expertBelt = mod }
}

// WithMetronome はメトロノーム（連続使用回数に応じて増加）を設定する。
func WithMetronome(mod float64) DamageInputOption {
	return func(d *DamageInput) { d.metronome = mod }
}

// WithLifeOrb はいのちのたま（1.3）を設定する。
func WithLifeOrb(mod float64) DamageInputOption { return func(d *DamageInput) { d.lifeOrb = mod } }

// WithHalfBerry は半減の実（効果抜群被弾時: 0.5）を設定する。
func WithHalfBerry(mod float64) DamageInputOption {
	return func(d *DamageInput) { d.halfBerry = mod }
}

func NewDamageInput(power Power, attack, defense, random int, opts ...DamageInputOption) *DamageInput {
	d := &DamageInput{
		power:       power,
		attack:      attack,
		defense:     defense,
		random:      random,
		critMod:     1.0,
		stabMod:     1.0,
		typeEff:     1.0,
		burnMod:     1.0,
		weather:     1.0,
		wall:        1.0,
		neuroforce:  1.0,
		sniper:      1.0,
		tintedLens:  1.0,
		fluffy:      1.0,
		mhalf:       1.0,
		mfilter:     1.0,
		mtwice:      1.0,
		friendGuard: 1.0,
		expertBelt:  1.0,
		metronome:   1.0,
		lifeOrb:     1.0,
		halfBerry:   1.0,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
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
