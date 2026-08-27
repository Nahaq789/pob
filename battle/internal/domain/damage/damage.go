package damage

import "math"

const LEVEL = 50
const PRECISION = 1e6

type DamageInput struct {
	Power   int
	Attack  int
	Defense int
	CritMod float64
	Random  int
	StabMod float64
	TypeEff float64
	BurnMod float64
	Weather float64
	Wall    float64
	Other   float64
}

func NewDamageInput(
	power, attack, defense int,
	critMod float64,
	random int,
	stabMod float64,
	typeEff float64,
	burnMod float64,
	weather, wall, other float64,
) *DamageInput {
	return &DamageInput{
		Power:   power,
		Attack:  attack,
		Defense: defense,
		CritMod: critMod,
		Random:  random,
		StabMod: stabMod,
		TypeEff: typeEff,
		BurnMod: burnMod,
		Weather: weather,
		Wall:    wall,
		Other:   other,
	}
}

func (d *DamageInput) CalcDamage() int {
	base := float64(((LEVEL*2/5+2)*d.Power*d.Attack)/50 + 2)

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

	return damage
}

// func (d *DamageInput) m()

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
