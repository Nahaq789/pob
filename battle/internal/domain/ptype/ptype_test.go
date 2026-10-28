package ptype_test

import (
	"testing"

	"pob/battle/internal/domain/ptype"
)

var allTypes = [...]ptype.Type{
	ptype.Normal, ptype.Fighting, ptype.Flying, ptype.Poison, ptype.Ground, ptype.Rock,
	ptype.Bug, ptype.Ghost, ptype.Steel, ptype.Fire, ptype.Water, ptype.Grass,
	ptype.Electric, ptype.Psychic, ptype.Ice, ptype.Dragon, ptype.Dark, ptype.Fairy,
}

func TestEffectiveness_SingleType(t *testing.T) {
	tests := []struct {
		name   string
		attack ptype.Type
		defend ptype.Type
		want   float64
	}{
		{"ほのお→くさ: 2倍", ptype.Fire, ptype.Grass, 2},
		{"みず→ほのお: 2倍", ptype.Water, ptype.Fire, 2},
		{"でんき→みず: 2倍", ptype.Electric, ptype.Water, 2},
		{"かくとう→あく: 2倍", ptype.Fighting, ptype.Dark, 2},
		{"ドラゴン→ドラゴン: 2倍", ptype.Dragon, ptype.Dragon, 2},
		{"フェアリー→ドラゴン: 2倍", ptype.Fairy, ptype.Dragon, 2},
		{"ほのお→みず: 0.5倍", ptype.Fire, ptype.Water, 0.5},
		{"くさ→はがね: 0.5倍", ptype.Grass, ptype.Steel, 0.5},
		{"はがね→はがね: 0.5倍", ptype.Steel, ptype.Steel, 0.5},
		{"ノーマル→いわ: 0.5倍", ptype.Normal, ptype.Rock, 0.5},
		{"ノーマル→ゴースト: 0倍", ptype.Normal, ptype.Ghost, 0},
		{"かくとう→ゴースト: 0倍", ptype.Fighting, ptype.Ghost, 0},
		{"ゴースト→ノーマル: 0倍", ptype.Ghost, ptype.Normal, 0},
		{"でんき→じめん: 0倍", ptype.Electric, ptype.Ground, 0},
		{"じめん→ひこう: 0倍", ptype.Ground, ptype.Flying, 0},
		{"エスパー→あく: 0倍", ptype.Psychic, ptype.Dark, 0},
		{"どく→はがね: 0倍", ptype.Poison, ptype.Steel, 0},
		{"ドラゴン→フェアリー: 0倍", ptype.Dragon, ptype.Fairy, 0},
		{"ノーマル→ノーマル: 1倍", ptype.Normal, ptype.Normal, 1},
		{"みず→でんき: 1倍", ptype.Water, ptype.Electric, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ptype.Effectiveness(tt.attack, [2]ptype.Type{tt.defend, ptype.None})
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestEffectiveness_DualType(t *testing.T) {
	tests := []struct {
		name   string
		attack ptype.Type
		defend [2]ptype.Type
		want   float64
	}{
		{"こおり→ドラゴン/ひこう: 4倍", ptype.Ice, [2]ptype.Type{ptype.Dragon, ptype.Flying}, 4},
		{"ほのお→くさ/むし: 4倍", ptype.Fire, [2]ptype.Type{ptype.Grass, ptype.Bug}, 4},
		{"いわ→ほのお/ひこう: 4倍", ptype.Rock, [2]ptype.Type{ptype.Fire, ptype.Flying}, 4},
		{"くさ→ほのお/ドラゴン: 0.25倍", ptype.Grass, [2]ptype.Type{ptype.Fire, ptype.Dragon}, 0.25},
		{"かくとう→どく/ひこう: 0.25倍", ptype.Fighting, [2]ptype.Type{ptype.Poison, ptype.Flying}, 0.25},
		{"ほのお→みず/くさ: 1倍", ptype.Fire, [2]ptype.Type{ptype.Water, ptype.Grass}, 1},
		{"でんき→みず/じめん: 0倍", ptype.Electric, [2]ptype.Type{ptype.Water, ptype.Ground}, 0},
		{"じめん→でんき/ひこう: 0倍", ptype.Ground, [2]ptype.Type{ptype.Electric, ptype.Flying}, 0},
		{"かくとう→ノーマル/ゴースト: 0倍", ptype.Fighting, [2]ptype.Type{ptype.Normal, ptype.Ghost}, 0},
		{"みず→ほのお/じめん: 4倍", ptype.Water, [2]ptype.Type{ptype.Fire, ptype.Ground}, 4},
		{"みず→じめん/ほのお: 4倍", ptype.Water, [2]ptype.Type{ptype.Ground, ptype.Fire}, 4},
		{"ほのお→みず/ドラゴン: 0.25倍", ptype.Fire, [2]ptype.Type{ptype.Water, ptype.Dragon}, 0.25},
		{"ほのお→ドラゴン/みず: 0.25倍", ptype.Fire, [2]ptype.Type{ptype.Dragon, ptype.Water}, 0.25},
		{"エスパー→ノーマル/ひこう: 1倍", ptype.Psychic, [2]ptype.Type{ptype.Normal, ptype.Flying}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ptype.Effectiveness(tt.attack, tt.defend)
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestEffectiveness_None(t *testing.T) {
	t.Run("第1タイプがNone: 第2タイプのみで計算", func(t *testing.T) {
		got := ptype.Effectiveness(ptype.Fire, [2]ptype.Type{ptype.None, ptype.Grass})
		if got != 2 {
			t.Errorf("expected 2, got %v", got)
		}
	})

	t.Run("両タイプがNone: 1倍", func(t *testing.T) {
		got := ptype.Effectiveness(ptype.Fire, [2]ptype.Type{ptype.None, ptype.None})
		if got != 1 {
			t.Errorf("expected 1, got %v", got)
		}
	})
}

func TestEffectiveness_FullChart(t *testing.T) {
	notNeutral := map[ptype.Type]map[ptype.Type]float64{
		ptype.Normal:   {ptype.Rock: 0.5, ptype.Ghost: 0, ptype.Steel: 0.5},
		ptype.Fire:     {ptype.Fire: 0.5, ptype.Water: 0.5, ptype.Grass: 2, ptype.Ice: 2, ptype.Bug: 2, ptype.Rock: 0.5, ptype.Dragon: 0.5, ptype.Steel: 2},
		ptype.Water:    {ptype.Fire: 2, ptype.Water: 0.5, ptype.Grass: 0.5, ptype.Ground: 2, ptype.Rock: 2, ptype.Dragon: 0.5},
		ptype.Electric: {ptype.Water: 2, ptype.Electric: 0.5, ptype.Grass: 0.5, ptype.Ground: 0, ptype.Flying: 2, ptype.Dragon: 0.5},
		ptype.Grass:    {ptype.Fire: 0.5, ptype.Water: 2, ptype.Grass: 0.5, ptype.Poison: 0.5, ptype.Ground: 2, ptype.Flying: 0.5, ptype.Bug: 0.5, ptype.Rock: 2, ptype.Dragon: 0.5, ptype.Steel: 0.5},
		ptype.Ice:      {ptype.Fire: 0.5, ptype.Water: 0.5, ptype.Grass: 2, ptype.Ice: 0.5, ptype.Ground: 2, ptype.Flying: 2, ptype.Dragon: 2, ptype.Steel: 0.5},
		ptype.Fighting: {ptype.Normal: 2, ptype.Ice: 2, ptype.Poison: 0.5, ptype.Flying: 0.5, ptype.Psychic: 0.5, ptype.Bug: 0.5, ptype.Rock: 2, ptype.Ghost: 0, ptype.Dark: 2, ptype.Steel: 2, ptype.Fairy: 0.5},
		ptype.Poison:   {ptype.Grass: 2, ptype.Poison: 0.5, ptype.Ground: 0.5, ptype.Rock: 0.5, ptype.Ghost: 0.5, ptype.Steel: 0, ptype.Fairy: 2},
		ptype.Ground:   {ptype.Fire: 2, ptype.Electric: 2, ptype.Grass: 0.5, ptype.Poison: 2, ptype.Flying: 0, ptype.Bug: 0.5, ptype.Rock: 2, ptype.Steel: 2},
		ptype.Flying:   {ptype.Electric: 0.5, ptype.Grass: 2, ptype.Fighting: 2, ptype.Bug: 2, ptype.Rock: 0.5, ptype.Steel: 0.5},
		ptype.Psychic:  {ptype.Fighting: 2, ptype.Poison: 2, ptype.Psychic: 0.5, ptype.Dark: 0, ptype.Steel: 0.5},
		ptype.Bug:      {ptype.Fire: 0.5, ptype.Grass: 2, ptype.Fighting: 0.5, ptype.Poison: 0.5, ptype.Flying: 0.5, ptype.Psychic: 2, ptype.Ghost: 0.5, ptype.Dark: 2, ptype.Steel: 0.5, ptype.Fairy: 0.5},
		ptype.Rock:     {ptype.Fire: 2, ptype.Ice: 2, ptype.Fighting: 0.5, ptype.Ground: 0.5, ptype.Flying: 2, ptype.Bug: 2, ptype.Steel: 0.5},
		ptype.Ghost:    {ptype.Normal: 0, ptype.Psychic: 2, ptype.Ghost: 2, ptype.Dark: 0.5},
		ptype.Dragon:   {ptype.Dragon: 2, ptype.Steel: 0.5, ptype.Fairy: 0},
		ptype.Dark:     {ptype.Fighting: 0.5, ptype.Psychic: 2, ptype.Ghost: 2, ptype.Dark: 0.5, ptype.Fairy: 0.5},
		ptype.Steel:    {ptype.Fire: 0.5, ptype.Water: 0.5, ptype.Electric: 0.5, ptype.Ice: 2, ptype.Rock: 2, ptype.Steel: 0.5, ptype.Fairy: 2},
		ptype.Fairy:    {ptype.Fire: 0.5, ptype.Fighting: 2, ptype.Poison: 0.5, ptype.Dragon: 2, ptype.Dark: 2, ptype.Steel: 0.5},
	}
	for _, atk := range allTypes {
		for _, def := range allTypes {
			want, ok := notNeutral[atk][def]
			if !ok {
				want = 1
			}
			got := ptype.Effectiveness(atk, [2]ptype.Type{def, ptype.None})
			if got != want {
				t.Errorf("attack=%d defend=%d: expected %v, got %v", atk, def, want, got)
			}
		}
	}
}
