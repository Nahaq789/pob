package stats

const lv = 50

func CalcStats(base, iv, ev int, natureModifier float64) int {
	f := float64(calcBase(base, iv, ev) + 5)
	return int(f * natureModifier)
}

func CalcHp(base, iv, ev int) int {
	actual := calcBase(base, iv, ev) + lv + 10
	return actual
}

func calcBase(base, iv, ev int) int {
	return (base*2 + iv + ev/4) * lv / 100
}
