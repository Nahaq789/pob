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
