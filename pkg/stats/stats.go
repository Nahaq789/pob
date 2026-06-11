package stats

const lv = 50

func CalcStats(base, iv, ev int, natureModifier float64) int {
	f := float64(((((base * 2) + iv + ev) / 4 * lv / 100) + 5))
	return int(f * natureModifier)
}

func CalcHp(base, iv, ev int) int {
	actual := ((((base * 2) + iv + ev) / 4) * lv / 100) + lv + 10
	return actual
}
