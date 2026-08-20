package damage

import "math"

const LEVEL = 50
const PRECISION = 1e6

type DamageInput struct {
	Power    int
	Attack   int
	Defense  int
	IsCrit   bool
	Random   int
	IsStab   bool
	TypeEff  float64
	IsBurned bool
	Weather  float64
	Wall     float64
	Other    float64
}

func NewDamageInput(
	power, attack, defense int,
	isCrit bool,
	random int,
	isStab bool,
	typeEff float64,
	isBurned bool,
	weather, wall, other float64,
) *DamageInput {
	return &DamageInput{
		Power:    power,
		Attack:   attack,
		Defense:  defense,
		IsCrit:   isCrit,
		Random:   random,
		IsStab:   isStab,
		TypeEff:  typeEff,
		IsBurned: isBurned,
		Weather:  weather,
		Wall:     wall,
		Other:    other,
	}
}

func (d *DamageInput) CalcDamage() int {
	base := ((LEVEL*2/5+2)*d.Power*d.Attack)/50 + 2
	damage := base * int(d.Weather)
	return damage
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
