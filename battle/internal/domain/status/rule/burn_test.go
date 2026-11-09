package rule_test

import (
	"testing"

	"pob/battle/internal/domain/ability"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/hp"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/nature"
	"pob/battle/internal/domain/phase"
	"pob/battle/internal/domain/player"
	"pob/battle/internal/domain/pokemon"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/rank"
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/rule"
)

func newBurnPokemon(name string, st status.Status) *pokemon.Pokemon {
	a := ability.NewAbility(0, "")
	return pokemon.NewPokemon(
		1, name, 0, "",
		[2]ptype.Type{ptype.Normal, ptype.None},
		pokemon.BaseStats{},
		pokemon.RealStats{Speed: 100},
		nature.Nature{},
		&a,
		[4]*move.Move{},
		hp.NewHP(100),
		rank.NewRank(),
		st,
		nil, nil,
		false,
	)
}

// newBurnStatus は c の主状態異常を持つ Status を返す。c が status.None の場合は状態異常なし。
func newBurnStatus(t *testing.T, c status.Condition) status.Status {
	t.Helper()
	if c == status.None {
		return status.NewStatus()
	}
	m, err := status.NewMainStatus(c)
	if err != nil {
		t.Fatal(err)
	}
	return status.NewStatusWith(&m, nil)
}

// newBurnContext は "actor" 側の先頭に attackerStatus のポケモンを置いたコンテキストを返す。
func newBurnContext(t *testing.T, category move.DamageClass, attackerStatus status.Condition) phase.MoveResolveContext {
	t.Helper()
	p1 := player.NewPlayer("actor", "Player1",
		[6]*pokemon.Pokemon{
			newBurnPokemon("こうげき1", newBurnStatus(t, attackerStatus)),
			newBurnPokemon("こうげき2", status.NewStatus()),
			newBurnPokemon("こうげき3", status.NewStatus()),
		}, nil)
	p2 := player.NewPlayer("opp", "Player2",
		[6]*pokemon.Pokemon{
			newBurnPokemon("あいて1", status.NewStatus()),
			newBurnPokemon("あいて2", status.NewStatus()),
			newBurnPokemon("あいて3", status.NewStatus()),
		}, nil)
	if err := p1.Select([3]int{0, 1, 2}); err != nil {
		t.Fatal(err)
	}
	if err := p2.Select([3]int{0, 1, 2}); err != nil {
		t.Fatal(err)
	}
	return phase.MoveResolveContext{
		Spec: damage.Spec{
			ActorId:  "actor",
			MoveId:   1,
			Type:     ptype.Normal,
			Power:    40,
			Category: category,
		},
		Battle: battle.NewBattle(p1, p2),
	}
}

func TestBurnHandler_Mod(t *testing.T) {
	tests := []struct {
		name      string
		category  move.DamageClass
		condition status.Condition
		wantBurn  bool
	}{
		{"やけど+物理技: 0.5倍", move.DamageClassPhysical, status.Burn, true},
		{"やけど+特殊技: 補正なし", move.DamageClassSpecial, status.Burn, false},
		{"状態異常なし+物理技: 補正なし", move.DamageClassPhysical, status.None, false},
		{"まひ+物理技: 補正なし", move.DamageClassPhysical, status.Paralysis, false},
		{"どく+物理技: 補正なし", move.DamageClassPhysical, status.Poison, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newBurnContext(t, tt.category, tt.condition)
			got := rule.NewBurnHandler().Mod(ctx).Burn
			if !tt.wantBurn {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("expected Burn to be set, got nil")
			}
			if got.Value != 0.5 {
				t.Errorf("expected 0.5, got %v", got.Value)
			}
		})
	}

	t.Run("やけど+物理技: Op が Multiply", func(t *testing.T) {
		ctx := newBurnContext(t, move.DamageClassPhysical, status.Burn)
		got := rule.NewBurnHandler().Mod(ctx).Burn
		if got == nil {
			t.Fatal("expected Burn to be set, got nil")
		}
		if got.Op != damage.ModOpMultiply {
			t.Errorf("expected ModOpMultiply, got %v", got.Op)
		}
	})

	t.Run("Burn 以外のフィールド: nil", func(t *testing.T) {
		ctx := newBurnContext(t, move.DamageClassPhysical, status.Burn)
		mod := rule.NewBurnHandler().Mod(ctx)
		mod.Burn = nil
		if mod != (damage.DamageMod{}) {
			t.Errorf("expected other fields nil, got %+v", mod)
		}
	})
}
