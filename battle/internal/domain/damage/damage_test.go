package damage

import "testing"

func newBaseDamageInput(power, attack, defense int) DamageInput {
	return DamageInput{
		power: NewPower(power), attack: attack, defense: defense,
		critMod: 1.0, random: 100, stabMod: 1.0, typeEff: 1.0,
		burnMod: 1.0, weather: 1.0, wall: 1.0,
		neuroforce: 1.0, sniper: 1.0, tintedLens: 1.0, fluffy: 1.0,
		mhalf: 1.0, mfilter: 1.0, mtwice: 1.0, friendGuard: 1.0,
		expertBelt: 1.0, metronome: 1.0, lifeOrb: 1.0, halfBerry: 1.0,
	}
}

func TestDamageInput_roundHalfDown(t *testing.T) {
	tests := []struct {
		name string
		v    float64
		want int
	}{
		{"0.4: 切り捨て", 1.4, 1},
		{"0.5: 五捨（切り捨て）", 1.5, 1},
		{"0.5超: 切り上げ", 1.6, 2},
		{"整数値", 2.0, 2},
		{"0.0", 0.0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roundHalfDown(tt.v)
			if got != tt.want {
				t.Errorf("roundHalfDown(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestDamageInput_CalcDamage(t *testing.T) {
	t.Run("補正なし: base=(22*100*100/100)/50+2=46", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		if got := d.CalcDamage(); got != 46 {
			t.Errorf("CalcDamage() = %v, want 46", got)
		}
	})

	t.Run("STAB補正(1.5x): 46*1.5=69.0→69", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.stabMod = 1.5
		if got := d.CalcDamage(); got != 69 {
			t.Errorf("CalcDamage() = %v, want 69", got)
		}
	})

	t.Run("STAB補正(1.5x): 五捨確認 5*1.5=7.5→7", func(t *testing.T) {
		// base=5: (22*10*35/50)/50+2 = 154/50+2 = 3+2 = 5
		d := newBaseDamageInput(10, 35, 50)
		d.stabMod = 1.5
		if got := d.CalcDamage(); got != 7 {
			t.Errorf("CalcDamage() = %v, want 7", got)
		}
	})

	t.Run("急所補正(1.5x): 46*1.5=69.0→69", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.critMod = 1.5
		if got := d.CalcDamage(); got != 69 {
			t.Errorf("CalcDamage() = %v, want 69", got)
		}
	})

	t.Run("タイプ相性(2.0x): 効果抜群 46*2=92", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.typeEff = 2.0
		if got := d.CalcDamage(); got != 92 {
			t.Errorf("CalcDamage() = %v, want 92", got)
		}
	})

	t.Run("タイプ相性(0.5x): 効果今一つ 46*0.5=23", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.typeEff = 0.5
		if got := d.CalcDamage(); got != 23 {
			t.Errorf("CalcDamage() = %v, want 23", got)
		}
	})

	t.Run("やけど補正(0.5x): 46*0.5=23", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.burnMod = 0.5
		if got := d.CalcDamage(); got != 23 {
			t.Errorf("CalcDamage() = %v, want 23", got)
		}
	})

	t.Run("乱数(85): 最小乱数 46*85/100=39", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.random = 85
		if got := d.CalcDamage(); got != 39 {
			t.Errorf("CalcDamage() = %v, want 39", got)
		}
	})

	t.Run("いのちのたま(1.3x): math.Round(46*1.3)=60", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.lifeOrb = 1.3
		if got := d.CalcDamage(); got != 60 {
			t.Errorf("CalcDamage() = %v, want 60", got)
		}
	})

	t.Run("複合補正: 急所+STAB+効果抜群 46→crit69→stab103→type206", func(t *testing.T) {
		// 急所(1.5): 46*1.5=69.0→69
		// STAB(1.5): 69*1.5=103.5→五捨103
		// タイプ(2.0): 103*2.0=206.0→206
		d := newBaseDamageInput(100, 100, 100)
		d.critMod = 1.5
		d.stabMod = 1.5
		d.typeEff = 2.0
		if got := d.CalcDamage(); got != 206 {
			t.Errorf("CalcDamage() = %v, want 206", got)
		}
	})
}
