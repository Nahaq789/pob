package pokemon_test

import (
	"testing"

	"pob/battle/internal/domain/ability"
	"pob/battle/internal/domain/hp"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/nature"
	"pob/battle/internal/domain/pokemon"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/rank"
	"pob/battle/internal/domain/status"
)

func newCritTestPokemon() *pokemon.Pokemon {
	a := ability.NewAbility(0, "")
	return pokemon.NewPokemon(
		1, "テスト", 0, "",
		[2]ptype.Type{ptype.Normal, ptype.Normal},
		pokemon.BaseStats{},
		pokemon.RealStats{Speed: 100},
		nature.Nature{},
		&a,
		[4]*move.Move{},
		hp.NewHP(100),
		rank.NewRank(),
		status.NewStatus(),
		nil, nil,
		false,
	)
}

func TestPokemon_CritStage(t *testing.T) {
	t.Run("NewRank で生成: stage=0", func(t *testing.T) {
		p := newCritTestPokemon()
		if got := p.CritStage(); got != 0 {
			t.Errorf("expected 0, got %d", got)
		}
	})
}

func TestPokemon_RollCritical(t *testing.T) {
	t.Run("bonus=3: 常に急所", func(t *testing.T) {
		p := newCritTestPokemon()
		for range 300 {
			crit, err := p.RollCritical(3)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !crit {
				t.Fatal("expected crit=true with bonus=3")
			}
		}
	})

	t.Run("bonus=0: true と false の両方が返る", func(t *testing.T) {
		p := newCritTestPokemon()
		gotTrue, gotFalse := false, false
		// 発生率 1/24 のため 300 試行では true が出ない可能性があり 3000 試行とする
		for range 3000 {
			crit, err := p.RollCritical(0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if crit {
				gotTrue = true
			} else {
				gotFalse = true
			}
			if gotTrue && gotFalse {
				break
			}
		}
		if !gotTrue {
			t.Error("3000試行で急所（true）が一度も発生しなかった")
		}
		if !gotFalse {
			t.Error("3000試行で急所なし（false）が一度も発生しなかった")
		}
	})
}

func newHasTypeTestPokemon(types [2]ptype.Type) *pokemon.Pokemon {
	a := ability.NewAbility(0, "")
	return pokemon.NewPokemon(
		1, "テスト", 0, "",
		types,
		pokemon.BaseStats{},
		pokemon.RealStats{Speed: 100},
		nature.Nature{},
		&a,
		[4]*move.Move{},
		hp.NewHP(100),
		rank.NewRank(),
		status.NewStatus(),
		nil, nil,
		false,
	)
}

func TestPokemon_HasType(t *testing.T) {
	tests := []struct {
		name  string
		types [2]ptype.Type
		arg   ptype.Type
		want  bool
	}{
		{"単タイプで一致: true", [2]ptype.Type{ptype.Fire, ptype.None}, ptype.Fire, true},
		{"複合タイプの第1タイプと一致: true", [2]ptype.Type{ptype.Water, ptype.Ground}, ptype.Water, true},
		{"複合タイプの第2タイプと一致: true", [2]ptype.Type{ptype.Water, ptype.Ground}, ptype.Ground, true},
		{"不一致: false", [2]ptype.Type{ptype.Water, ptype.Ground}, ptype.Fire, false},
		// 単タイプの [1] は ptype.None だが、None は常に false を返す。
		{"単タイプに None を指定: false", [2]ptype.Type{ptype.Fire, ptype.None}, ptype.None, false},
		{"複合タイプに None を指定: false", [2]ptype.Type{ptype.Water, ptype.Ground}, ptype.None, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newHasTypeTestPokemon(tt.types)
			if got := p.HasType(tt.arg); got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
