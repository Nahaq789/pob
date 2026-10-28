package damage

import "testing"

type modField struct {
	name  string
	set   func(m *DamageMod, v *ModValue)
	get   func(m DamageMod) *ModValue
	input func(d *DamageInput) float64
}

func newModFields() []modField {
	return []modField{
		{"Crit", func(m *DamageMod, v *ModValue) { m.Crit = v }, func(m DamageMod) *ModValue { return m.Crit }, func(d *DamageInput) float64 { return d.critMod }},
		{"Stab", func(m *DamageMod, v *ModValue) { m.Stab = v }, func(m DamageMod) *ModValue { return m.Stab }, func(d *DamageInput) float64 { return d.stabMod }},
		{"TypeEff", func(m *DamageMod, v *ModValue) { m.TypeEff = v }, func(m DamageMod) *ModValue { return m.TypeEff }, func(d *DamageInput) float64 { return d.typeEff }},
		{"Burn", func(m *DamageMod, v *ModValue) { m.Burn = v }, func(m DamageMod) *ModValue { return m.Burn }, func(d *DamageInput) float64 { return d.burnMod }},
		{"Weather", func(m *DamageMod, v *ModValue) { m.Weather = v }, func(m DamageMod) *ModValue { return m.Weather }, func(d *DamageInput) float64 { return d.weather }},
		{"Wall", func(m *DamageMod, v *ModValue) { m.Wall = v }, func(m DamageMod) *ModValue { return m.Wall }, func(d *DamageInput) float64 { return d.wall }},
		{"Neuroforce", func(m *DamageMod, v *ModValue) { m.Neuroforce = v }, func(m DamageMod) *ModValue { return m.Neuroforce }, func(d *DamageInput) float64 { return d.neuroforce }},
		{"Sniper", func(m *DamageMod, v *ModValue) { m.Sniper = v }, func(m DamageMod) *ModValue { return m.Sniper }, func(d *DamageInput) float64 { return d.sniper }},
		{"TintedLens", func(m *DamageMod, v *ModValue) { m.TintedLens = v }, func(m DamageMod) *ModValue { return m.TintedLens }, func(d *DamageInput) float64 { return d.tintedLens }},
		{"Fluffy", func(m *DamageMod, v *ModValue) { m.Fluffy = v }, func(m DamageMod) *ModValue { return m.Fluffy }, func(d *DamageInput) float64 { return d.fluffy }},
		{"MHalf", func(m *DamageMod, v *ModValue) { m.MHalf = v }, func(m DamageMod) *ModValue { return m.MHalf }, func(d *DamageInput) float64 { return d.mhalf }},
		{"MFilter", func(m *DamageMod, v *ModValue) { m.MFilter = v }, func(m DamageMod) *ModValue { return m.MFilter }, func(d *DamageInput) float64 { return d.mfilter }},
		{"MTwice", func(m *DamageMod, v *ModValue) { m.MTwice = v }, func(m DamageMod) *ModValue { return m.MTwice }, func(d *DamageInput) float64 { return d.mtwice }},
		{"ExpertBelt", func(m *DamageMod, v *ModValue) { m.ExpertBelt = v }, func(m DamageMod) *ModValue { return m.ExpertBelt }, func(d *DamageInput) float64 { return d.expertBelt }},
		{"Metronome", func(m *DamageMod, v *ModValue) { m.Metronome = v }, func(m DamageMod) *ModValue { return m.Metronome }, func(d *DamageInput) float64 { return d.metronome }},
		{"LifeOrb", func(m *DamageMod, v *ModValue) { m.LifeOrb = v }, func(m DamageMod) *ModValue { return m.LifeOrb }, func(d *DamageInput) float64 { return d.lifeOrb }},
		{"HalfBerry", func(m *DamageMod, v *ModValue) { m.HalfBerry = v }, func(m DamageMod) *ModValue { return m.HalfBerry }, func(d *DamageInput) float64 { return d.halfBerry }},
	}
}

func newModInput(m DamageMod) *DamageInput {
	return NewDamageInput(NewPower(100), 100, 100, 100, m.ToOptions()...)
}

func TestOverride(t *testing.T) {
	v := Override(1.5)
	t.Run("Override(1.5): Value=1.5", func(t *testing.T) {
		if v.Value != 1.5 {
			t.Errorf("expected 1.5, got %v", v.Value)
		}
	})
	t.Run("Override(1.5): Op=ModOpOverride", func(t *testing.T) {
		if v.Op != ModOpOverride {
			t.Errorf("expected ModOpOverride, got %v", v.Op)
		}
	})
}

func TestMultiply(t *testing.T) {
	v := Multiply(1.3)
	t.Run("Multiply(1.3): Value=1.3", func(t *testing.T) {
		if v.Value != 1.3 {
			t.Errorf("expected 1.3, got %v", v.Value)
		}
	})
	t.Run("Multiply(1.3): Op=ModOpMultiply", func(t *testing.T) {
		if v.Op != ModOpMultiply {
			t.Errorf("expected ModOpMultiply, got %v", v.Op)
		}
	})
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name      string
		base      *ModValue
		incoming  *ModValue
		wantNil   bool
		wantValue float64
		wantOp    ModOp
	}{
		{"base=nil, incoming=nil: nil", nil, nil, true, 0, 0},
		{"base=Multiply(1.5), incoming=nil: base を維持", Multiply(1.5), nil, false, 1.5, ModOpMultiply},
		{"base=Override(0), incoming=nil: base を維持", Override(0), nil, false, 0, ModOpOverride},
		{"base=nil, incoming=Multiply(2.0): incoming を使用", nil, Multiply(2.0), false, 2.0, ModOpMultiply},
		{"base=nil, incoming=Override(0.5): incoming を使用", nil, Override(0.5), false, 0.5, ModOpOverride},
		{"base=Multiply(1.5), incoming=Multiply(2.0): 乗算で 3.0", Multiply(1.5), Multiply(2.0), false, 3.0, ModOpMultiply},
		{"base=Multiply(1.5), incoming=Override(1.0): 上書きで 1.0", Multiply(1.5), Override(1.0), false, 1.0, ModOpOverride},
		{"base=Override(2.0), incoming=Override(0.5): 上書きで 0.5", Override(2.0), Override(0.5), false, 0.5, ModOpOverride},
		// base の Op に関わらず incoming が Multiply なら乗算結果は Multiply として返る
		{"base=Override(2.0), incoming=Multiply(1.5): 乗算で 3.0", Override(2.0), Multiply(1.5), false, 3.0, ModOpMultiply},
		{"base=Multiply(1.5), incoming=Multiply(0): 乗算で 0", Multiply(1.5), Multiply(0), false, 0, ModOpMultiply},
	}
	for _, f := range newModFields() {
		for _, tt := range tests {
			t.Run(f.name+"/"+tt.name, func(t *testing.T) {
				var base, incoming DamageMod
				f.set(&base, tt.base)
				f.set(&incoming, tt.incoming)

				got := f.get(Merge(base, incoming))
				if tt.wantNil {
					if got != nil {
						t.Errorf("expected nil, got %+v", *got)
					}
					return
				}
				if got == nil {
					t.Fatal("expected non-nil, got nil")
				}
				if got.Value != tt.wantValue {
					t.Errorf("expected Value=%v, got %v", tt.wantValue, got.Value)
				}
				if got.Op != tt.wantOp {
					t.Errorf("expected Op=%v, got %v", tt.wantOp, got.Op)
				}
			})
		}
	}

	t.Run("フィールド間独立: 他フィールドに影響しない", func(t *testing.T) {
		fields := newModFields()
		for i, target := range fields {
			var incoming DamageMod
			target.set(&incoming, Multiply(2.0))
			merged := Merge(DamageMod{}, incoming)
			for j, other := range fields {
				if i == j {
					continue
				}
				if v := other.get(merged); v != nil {
					t.Errorf("%s を設定したのに %s が非 nil: %+v", target.name, other.name, *v)
				}
			}
		}
	})

	t.Run("3段マージ: Multiply(1.5)→Multiply(2.0)→Multiply(0.5) で 1.5", func(t *testing.T) {
		m := Merge(DamageMod{}, DamageMod{Stab: Multiply(1.5)})
		m = Merge(m, DamageMod{Stab: Multiply(2.0)})
		m = Merge(m, DamageMod{Stab: Multiply(0.5)})
		if m.Stab == nil {
			t.Fatal("expected non-nil Stab, got nil")
		}
		if m.Stab.Value != 1.5 {
			t.Errorf("expected 1.5, got %v", m.Stab.Value)
		}
	})

	t.Run("乗算結果: base の ModValue を書き換えない", func(t *testing.T) {
		b := Multiply(1.5)
		Merge(DamageMod{Crit: b}, DamageMod{Crit: Multiply(2.0)})
		if b.Value != 1.5 {
			t.Errorf("expected base Value=1.5, got %v", b.Value)
		}
	})
}

func TestDamageMod_ToOptions(t *testing.T) {
	t.Run("全フィールド nil: オプション0件", func(t *testing.T) {
		if got := (DamageMod{}).ToOptions(); len(got) != 0 {
			t.Errorf("expected 0 options, got %d", len(got))
		}
	})

	t.Run("全フィールド nil: 全補正がデフォルト 1.0", func(t *testing.T) {
		d := newModInput(DamageMod{})
		for _, f := range newModFields() {
			if got := f.input(d); got != 1.0 {
				t.Errorf("%s: expected 1.0, got %v", f.name, got)
			}
		}
	})

	for _, f := range newModFields() {
		t.Run(f.name+" のみ設定: オプション1件", func(t *testing.T) {
			var m DamageMod
			f.set(&m, Multiply(1.75))
			if got := m.ToOptions(); len(got) != 1 {
				t.Errorf("expected 1 option, got %d", len(got))
			}
		})

		t.Run(f.name+" のみ設定: 対応する DamageInput フィールドに反映", func(t *testing.T) {
			var m DamageMod
			f.set(&m, Multiply(1.75))
			d := newModInput(m)
			if got := f.input(d); got != 1.75 {
				t.Errorf("expected 1.75, got %v", got)
			}
		})

		t.Run(f.name+" のみ設定: 他フィールドはデフォルト 1.0", func(t *testing.T) {
			var m DamageMod
			f.set(&m, Multiply(1.75))
			d := newModInput(m)
			for _, other := range newModFields() {
				if other.name == f.name {
					continue
				}
				if got := other.input(d); got != 1.0 {
					t.Errorf("%s: expected 1.0, got %v", other.name, got)
				}
			}
		})
	}

	t.Run("Override 指定: Value がそのまま反映", func(t *testing.T) {
		d := newModInput(DamageMod{TypeEff: Override(0)})
		if d.typeEff != 0 {
			t.Errorf("expected 0, got %v", d.typeEff)
		}
	})

	t.Run("全フィールド設定: オプション17件", func(t *testing.T) {
		var m DamageMod
		for _, f := range newModFields() {
			f.set(&m, Multiply(2.0))
		}
		if got := m.ToOptions(); len(got) != 17 {
			t.Errorf("expected 17 options, got %d", len(got))
		}
	})

	t.Run("Merge 結果を適用: STAB1.5 × 急所1.5 で CalcDamage=103", func(t *testing.T) {
		// base=46 → 急所 46*1.5=69 → STAB 69*1.5=103.5→103（五捨）
		m := Merge(DamageMod{Stab: Multiply(1.5)}, DamageMod{Crit: Multiply(1.5)})
		if got := newModInput(m).CalcDamage(); got != 103 {
			t.Errorf("expected 103, got %v", got)
		}
	})

	t.Run("TypeEff を Override(0): CalcDamage=0", func(t *testing.T) {
		m := Merge(DamageMod{TypeEff: Multiply(2.0)}, DamageMod{TypeEff: Override(0)})
		if got := newModInput(m).CalcDamage(); got != 0 {
			t.Errorf("expected 0, got %v", got)
		}
	})
}
