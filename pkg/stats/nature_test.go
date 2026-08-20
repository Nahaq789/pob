package stats

import "testing"

func TestGetNatureModifiers(t *testing.T) {
	t.Run("未知の性格: エラーを返す", func(t *testing.T) {
		_, err := GetNatureModifiers("unknown")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("補正なし性格: すべて 1.0", func(t *testing.T) {
		for _, name := range []string{"がんばりや", "すなお", "てれや", "きまぐれ", "まじめ"} {
			m, err := GetNatureModifiers(name)
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", name, err)
			}
			for _, v := range []float64{m.H, m.A, m.B, m.C, m.D, m.S} {
				if v != 1.0 {
					t.Errorf("%s: expected all 1.0, got %+v", name, m)
					break
				}
			}
		}
	})

	t.Run("A上昇性格: A=1.1", func(t *testing.T) {
		for _, name := range []string{"さみしがり", "いじっぱり", "やんちゃ", "ゆうかん"} {
			m, err := GetNatureModifiers(name)
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", name, err)
			}
			if m.A != 1.1 {
				t.Errorf("%s: A = %v, want 1.1", name, m.A)
			}
		}
	})

	t.Run("B上昇性格: B=1.1", func(t *testing.T) {
		for _, name := range []string{"ずぶとい", "わんぱく", "のうてんき", "のんき"} {
			m, err := GetNatureModifiers(name)
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", name, err)
			}
			if m.B != 1.1 {
				t.Errorf("%s: B = %v, want 1.1", name, m.B)
			}
		}
	})

	t.Run("C上昇性格: C=1.1", func(t *testing.T) {
		for _, name := range []string{"ひかえめ", "おっとり", "うっかりや", "れいせい"} {
			m, err := GetNatureModifiers(name)
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", name, err)
			}
			if m.C != 1.1 {
				t.Errorf("%s: C = %v, want 1.1", name, m.C)
			}
		}
	})

	t.Run("D上昇性格: D=1.1", func(t *testing.T) {
		for _, name := range []string{"おだやか", "おとなしい", "しんちょう", "なまいき"} {
			m, err := GetNatureModifiers(name)
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", name, err)
			}
			if m.D != 1.1 {
				t.Errorf("%s: D = %v, want 1.1", name, m.D)
			}
		}
	})

	t.Run("S上昇性格: S=1.1", func(t *testing.T) {
		for _, name := range []string{"おくびょう", "せっかち", "ようき", "むじゃき"} {
			m, err := GetNatureModifiers(name)
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", name, err)
			}
			if m.S != 1.1 {
				t.Errorf("%s: S = %v, want 1.1", name, m.S)
			}
		}
	})

	t.Run("各性格の上昇と下降が対になっている", func(t *testing.T) {
		tests := []struct {
			name   string
			upStat func(NatureModifier) float64
			dnStat func(NatureModifier) float64
		}{
			{"さみしがり", func(m NatureModifier) float64 { return m.A }, func(m NatureModifier) float64 { return m.B }},
			{"いじっぱり", func(m NatureModifier) float64 { return m.A }, func(m NatureModifier) float64 { return m.C }},
			{"やんちゃ", func(m NatureModifier) float64 { return m.A }, func(m NatureModifier) float64 { return m.D }},
			{"ゆうかん", func(m NatureModifier) float64 { return m.A }, func(m NatureModifier) float64 { return m.S }},
			{"ひかえめ", func(m NatureModifier) float64 { return m.C }, func(m NatureModifier) float64 { return m.A }},
			{"おくびょう", func(m NatureModifier) float64 { return m.S }, func(m NatureModifier) float64 { return m.A }},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				m, err := GetNatureModifiers(tt.name)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.upStat(m) != 1.1 {
					t.Errorf("up stat = %v, want 1.1", tt.upStat(m))
				}
				if tt.dnStat(m) != 0.9 {
					t.Errorf("down stat = %v, want 0.9", tt.dnStat(m))
				}
			})
		}
	})

	t.Run("H（HP）は常に 1.0", func(t *testing.T) {
		// HP に補正をかける性格は存在しない
		allNatures := []string{
			"がんばりや", "さみしがり", "いじっぱり", "やんちゃ", "ゆうかん",
			"ずぶとい", "わんぱく", "のうてんき", "のんき",
			"ひかえめ", "おっとり", "うっかりや", "れいせい",
			"おだやか", "おとなしい", "しんちょう", "なまいき",
			"おくびょう", "せっかち", "ようき", "むじゃき",
			"すなお", "てれや", "きまぐれ", "まじめ",
		}
		for _, name := range allNatures {
			m, err := GetNatureModifiers(name)
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", name, err)
			}
			if m.H != 1.0 {
				t.Errorf("%s: H = %v, want 1.0", name, m.H)
			}
		}
	})
}
