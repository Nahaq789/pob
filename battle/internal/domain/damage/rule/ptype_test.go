package rule_test

import (
	"testing"

	"pob/battle/internal/domain/ability"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/damage/rule"
	"pob/battle/internal/domain/hp"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/nature"
	"pob/battle/internal/domain/phase"
	"pob/battle/internal/domain/player"
	"pob/battle/internal/domain/pokemon"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/rank"
	"pob/battle/internal/domain/status"
)

func newTypeTestPokemon(name string, types [2]ptype.Type) *pokemon.Pokemon {
	a := ability.NewAbility(0, "")
	return pokemon.NewPokemon(
		1, name, 0, "",
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

// newTypeTestContext は "actor" 側の先頭に attackerTypes、相手側の先頭に defenderTypes のポケモンを置いたコンテキストを返す。
func newTypeTestContext(moveType ptype.Type, attackerTypes, defenderTypes [2]ptype.Type) phase.MoveResolveContext {
	p1 := player.NewPlayer("actor", "Player1",
		[6]*pokemon.Pokemon{
			newTypeTestPokemon("こうげき1", attackerTypes),
			newTypeTestPokemon("こうげき2", attackerTypes),
			newTypeTestPokemon("こうげき3", attackerTypes),
		}, nil)
	p2 := player.NewPlayer("opp", "Player2",
		[6]*pokemon.Pokemon{
			newTypeTestPokemon("あいて1", defenderTypes),
			newTypeTestPokemon("あいて2", defenderTypes),
			newTypeTestPokemon("あいて3", defenderTypes),
		}, nil)
	if err := p1.Select([3]int{0, 1, 2}); err != nil {
		panic(err)
	}
	if err := p2.Select([3]int{0, 1, 2}); err != nil {
		panic(err)
	}
	return phase.MoveResolveContext{
		Spec: damage.Spec{
			ActorId:  "actor",
			MoveId:   1,
			Type:     moveType,
			Power:    40,
			Category: move.DamageClassPhysical,
		},
		Battle: battle.NewBattle(p1, p2),
	}
}

func TestPtypeHandler_Mod(t *testing.T) {
	single := func(tp ptype.Type) [2]ptype.Type { return [2]ptype.Type{tp, ptype.None} }

	tests := []struct {
		name          string
		moveType      ptype.Type
		defenderTypes [2]ptype.Type
		want          float64
	}{
		{"ほのお→くさ: 2倍", ptype.Fire, single(ptype.Grass), 2},
		{"ほのお→みず: 0.5倍", ptype.Fire, single(ptype.Water), 0.5},
		{"ノーマル→ゴースト: 0倍", ptype.Normal, single(ptype.Ghost), 0},
		{"ノーマル→ノーマル: 1倍", ptype.Normal, single(ptype.Normal), 1},
		{"みず→ほのお/じめん: 4倍", ptype.Water, [2]ptype.Type{ptype.Fire, ptype.Ground}, 4},
		{"ほのお→みず/ほのお: 0.25倍", ptype.Fire, [2]ptype.Type{ptype.Water, ptype.Fire}, 0.25},
		{"でんき→みず/じめん: 0倍", ptype.Electric, [2]ptype.Type{ptype.Water, ptype.Ground}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTypeTestContext(tt.moveType, single(ptype.Normal), tt.defenderTypes)
			got := rule.NewPtypeHandler().Mod(ctx).TypeEff
			if got == nil {
				t.Fatal("expected TypeEff to be set, got nil")
			}
			if got.Value != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got.Value)
			}
		})
	}

	t.Run("補正あり: Op が Multiply", func(t *testing.T) {
		ctx := newTypeTestContext(ptype.Fire, single(ptype.Normal), single(ptype.Grass))
		got := rule.NewPtypeHandler().Mod(ctx).TypeEff
		if got == nil {
			t.Fatal("expected TypeEff to be set, got nil")
		}
		if got.Op != damage.ModOpMultiply {
			t.Errorf("expected ModOpMultiply, got %v", got.Op)
		}
	})

	t.Run("攻撃側のタイプ: 相性に影響しない", func(t *testing.T) {
		ctx := newTypeTestContext(ptype.Fire, single(ptype.Water), single(ptype.Grass))
		got := rule.NewPtypeHandler().Mod(ctx).TypeEff
		if got == nil {
			t.Fatal("expected TypeEff to be set, got nil")
		}
		if got.Value != 2 {
			t.Errorf("expected 2, got %v", got.Value)
		}
	})

	t.Run("TypeEff 以外のフィールド: nil", func(t *testing.T) {
		ctx := newTypeTestContext(ptype.Fire, single(ptype.Normal), single(ptype.Grass))
		mod := rule.NewPtypeHandler().Mod(ctx)
		mod.TypeEff = nil
		if mod != (damage.DamageMod{}) {
			t.Errorf("expected other fields nil, got %+v", mod)
		}
	})
}
