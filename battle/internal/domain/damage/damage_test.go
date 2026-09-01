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

	// ── 天気補正 ──
	t.Run("天気補正(1.5x): あめ水技/にほんばれ炎技 roundHalfDown(46*1.5)=69", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.weather = 1.5
		if got := d.CalcDamage(); got != 69 {
			t.Errorf("CalcDamage() = %v, want 69", got)
		}
	})

	t.Run("天気補正(0.5x): あめ炎技/にほんばれ水技 roundHalfDown(46*0.5)=23", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.weather = 0.5
		if got := d.CalcDamage(); got != 23 {
			t.Errorf("CalcDamage() = %v, want 23", got)
		}
	})

	// ── タイプ一致補正（てきおうりょく） ──
	t.Run("STABてきおうりょく(2.0x): roundHalfDown(46*2.0)=92", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.stabMod = 2.0
		if got := d.CalcDamage(); got != 92 {
			t.Errorf("CalcDamage() = %v, want 92", got)
		}
	})

	// ── 相性補正（切り捨て） ──
	t.Run("相性補正(4.0x): ダブル弱点 int(46*4.0)=184", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.typeEff = 4.0
		if got := d.CalcDamage(); got != 184 {
			t.Errorf("CalcDamage() = %v, want 184", got)
		}
	})

	t.Run("相性補正(0.25x): 4倍耐性 切り捨て int(7*0.25)=1", func(t *testing.T) {
		// base=7: (22*15*40/50)/50+2 = 264/50+2 = 5+2 = 7
		d := newBaseDamageInput(15, 40, 50)
		d.typeEff = 0.25
		if got := d.CalcDamage(); got != 1 {
			t.Errorf("CalcDamage() = %v, want 1", got)
		}
	})

	// ── M補正（各ステップで math.Round） ──
	t.Run("壁補正(0.5x): リフレクター/ひかりのかべ シングル math.Round(46*0.5)=23", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.wall = 0.5
		if got := d.CalcDamage(); got != 23 {
			t.Errorf("CalcDamage() = %v, want 23", got)
		}
	})

	t.Run("ブレインフォース(1.25x): 効果バツグン時 math.Round(46*1.25)=58", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.neuroforce = 1.25
		if got := d.CalcDamage(); got != 58 {
			t.Errorf("CalcDamage() = %v, want 58", got)
		}
	})

	t.Run("スナイパー急所(×1.5×1.5): math.Round(69*1.5)=104", func(t *testing.T) {
		// 急所: roundHalfDown(46*1.5)=69, スナイパー: math.Round(69*1.5)=104
		d := newBaseDamageInput(100, 100, 100)
		d.critMod = 1.5
		d.sniper = 1.5
		if got := d.CalcDamage(); got != 104 {
			t.Errorf("CalcDamage() = %v, want 104", got)
		}
	})

	t.Run("いろめがね(×2.0): いまひとつ 46→23→46", func(t *testing.T) {
		// typeEff=0.5で切り捨てint(46*0.5)=23, tintedLens=2.0でmath.Round(23*2)=46
		d := newBaseDamageInput(100, 100, 100)
		d.typeEff = 0.5
		d.tintedLens = 2.0
		if got := d.CalcDamage(); got != 46 {
			t.Errorf("CalcDamage() = %v, want 46", got)
		}
	})

	t.Run("もふもふ ほのおタイプ(×2.0): math.Round(46*2.0)=92", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.fluffy = 2.0
		if got := d.CalcDamage(); got != 92 {
			t.Errorf("CalcDamage() = %v, want 92", got)
		}
	})

	t.Run("もふもふ 直接攻撃(×0.5): Mhalf位置 math.Round(46*0.5)=23", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.mhalf = 0.5
		if got := d.CalcDamage(); got != 23 {
			t.Errorf("CalcDamage() = %v, want 23", got)
		}
	})

	t.Run("もふもふ 炎タイプ+直接攻撃(×2.0×0.5): fluffy→math.Round(46*2.0)=92, Mhalf→math.Round(92*0.5)=46", func(t *testing.T) {
		// 炎: fluffy位置 math.Round(46*2.0)=92
		// 直接攻撃: Mhalf位置 math.Round(92*0.5)=46
		d := newBaseDamageInput(100, 100, 100)
		d.fluffy = 2.0
		d.mhalf = 0.5
		if got := d.CalcDamage(); got != 46 {
			t.Errorf("CalcDamage() = %v, want 46", got)
		}
	})

	t.Run("Mhalf(×0.5): こおりのりんぷん/マルチスケイル等 math.Round(46*0.5)=23", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.mhalf = 0.5
		if got := d.CalcDamage(); got != 23 {
			t.Errorf("CalcDamage() = %v, want 23", got)
		}
	})

	t.Run("Mfilter(×0.75): フィルター/ハードロック等 効果バツグン軽減 math.Round(46*0.75)=35", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.mfilter = 0.75
		if got := d.CalcDamage(); got != 35 {
			t.Errorf("CalcDamage() = %v, want 35", got)
		}
	})

	t.Run("たつじんのおび(4915/4096): 効果バツグン時 math.Round(46*4915/4096)=55", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.expertBelt = float64(4915) / 4096
		if got := d.CalcDamage(); got != 55 {
			t.Errorf("CalcDamage() = %v, want 55", got)
		}
	})

	t.Run("メトロノーム2回(4915/4096): math.Round(46*4915/4096)=55", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.metronome = float64(4915) / 4096
		if got := d.CalcDamage(); got != 55 {
			t.Errorf("CalcDamage() = %v, want 55", got)
		}
	})

	t.Run("メトロノーム3回(5734/4096): math.Round(46*5734/4096)=64", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.metronome = float64(5734) / 4096
		if got := d.CalcDamage(); got != 64 {
			t.Errorf("CalcDamage() = %v, want 64", got)
		}
	})

	t.Run("メトロノーム6回以上(8192/4096=2.0): math.Round(46*2.0)=92", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.metronome = float64(8192) / 4096
		if got := d.CalcDamage(); got != 92 {
			t.Errorf("CalcDamage() = %v, want 92", got)
		}
	})

	t.Run("半減の実(2048/4096): 効果バツグン被弾時 math.Round(46*0.5)=23", func(t *testing.T) {
		d := newBaseDamageInput(100, 100, 100)
		d.halfBerry = float64(2048) / 4096
		if got := d.CalcDamage(); got != 23 {
			t.Errorf("CalcDamage() = %v, want 23", got)
		}
	})
}

func TestDamageInput_random(t *testing.T) {
	// base=46: (22*100*100/100)/50+2=46
	// 乱数補正は切り捨て: damage * random / 100
	tests := []struct {
		name   string
		random int
		want   int
	}{
		{"乱数85(最小): 46*85/100=39", 85, 39},
		{"乱数86: 46*86/100=39", 86, 39},
		{"乱数87: 46*87/100=40", 87, 40},
		{"乱数88: 46*88/100=40", 88, 40},
		{"乱数89: 46*89/100=40", 89, 40},
		{"乱数90: 46*90/100=41", 90, 41},
		{"乱数91: 46*91/100=41", 91, 41},
		{"乱数92: 46*92/100=42", 92, 42},
		{"乱数93: 46*93/100=42", 93, 42},
		{"乱数94: 46*94/100=43", 94, 43},
		{"乱数95: 46*95/100=43", 95, 43},
		{"乱数96: 46*96/100=44", 96, 44},
		{"乱数97: 46*97/100=44", 97, 44},
		{"乱数98: 46*98/100=45", 98, 45},
		{"乱数99: 46*99/100=45", 99, 45},
		{"乱数100(最大): 46*100/100=46", 100, 46},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newBaseDamageInput(100, 100, 100)
			d.random = tt.random
			if got := d.CalcDamage(); got != tt.want {
				t.Errorf("CalcDamage() = %v, want %v", got, tt.want)
			}
		})
	}
}
