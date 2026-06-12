package stats

import "testing"

func TestCalcHp(t *testing.T) {
	tests := []struct {
		name string
		base int
		iv   int
		ev   int
		want int
	}{
		// HP = (base*2 + iv + ev/4) * 50/100 + 50 + 10  (Level 50)
		{name: "all zero", base: 0, iv: 0, ev: 0, want: 60},
		{name: "base only", base: 50, iv: 0, ev: 0, want: 110},
		{name: "max iv", base: 50, iv: 31, ev: 0, want: 125},
		{name: "max ev", base: 50, iv: 0, ev: 252, want: 141},
		{name: "max iv and ev", base: 50, iv: 31, ev: 252, want: 157},
		// EV floor division: 253 behaves same as 252
		{name: "ev floor division", base: 50, iv: 0, ev: 253, want: 141},
		// Blissey base HP = 255 (highest in game)
		{name: "blissey max", base: 255, iv: 31, ev: 252, want: 362},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcHp(tt.base, tt.iv, tt.ev)
			if got != tt.want {
				t.Errorf("CalcHp(%d, %d, %d) = %d, want %d", tt.base, tt.iv, tt.ev, got, tt.want)
			}
		})
	}
}

func TestCalcStats(t *testing.T) {
	tests := []struct {
		name           string
		base           int
		iv             int
		ev             int
		natureModifier float64
		want           int
	}{
		// Stat = int(float64((base*2 + iv + ev/4)*50/100 + 5) * nature)  (Level 50)
		{name: "all zero neutral", base: 0, iv: 0, ev: 0, natureModifier: 1.0, want: 5},
		{name: "base only neutral", base: 50, iv: 0, ev: 0, natureModifier: 1.0, want: 55},
		{name: "max iv neutral", base: 50, iv: 31, ev: 0, natureModifier: 1.0, want: 70},
		{name: "max ev neutral", base: 50, iv: 0, ev: 252, natureModifier: 1.0, want: 86},
		{name: "max iv ev neutral", base: 50, iv: 31, ev: 252, natureModifier: 1.0, want: 102},
		{name: "nature boost 1.1", base: 50, iv: 31, ev: 252, natureModifier: 1.1, want: 112},
		{name: "nature down 0.9", base: 50, iv: 31, ev: 252, natureModifier: 0.9, want: 91},
		// EV floor division: 255 behaves same as 252 (255/4=63)
		{name: "ev floor division", base: 50, iv: 31, ev: 255, natureModifier: 1.0, want: 102},
		// Mewtwo base Sp.Atk = 154
		{name: "mewtwo spatk max", base: 154, iv: 31, ev: 252, natureModifier: 1.1, want: 226},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcStats(tt.base, tt.iv, tt.ev, tt.natureModifier)
			if got != tt.want {
				t.Errorf("CalcStats(%d, %d, %d, %.1f) = %d, want %d", tt.base, tt.iv, tt.ev, tt.natureModifier, got, tt.want)
			}
		})
	}
}
