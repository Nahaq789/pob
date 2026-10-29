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

func newCritPokemon(name string) *pokemon.Pokemon {
	a := ability.NewAbility(0, "")
	return pokemon.NewPokemon(
		1, name, 0, "",
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

// newCritBattle は "actor" を攻撃側プレイヤーとしてバトルを構築する。
func newCritBattle() *battle.Battle {
	p1 := player.NewPlayer("actor", "Player1",
		[6]*pokemon.Pokemon{newCritPokemon("こうげき1"), newCritPokemon("こうげき2"), newCritPokemon("こうげき3")}, nil)
	p2 := player.NewPlayer("opp", "Player2",
		[6]*pokemon.Pokemon{newCritPokemon("あいて1"), newCritPokemon("あいて2"), newCritPokemon("あいて3")}, nil)
	if err := p1.Select([3]int{0, 1, 2}); err != nil {
		panic(err)
	}
	if err := p2.Select([3]int{0, 1, 2}); err != nil {
		panic(err)
	}
	return battle.NewBattle(p1, p2)
}

func newCritContext(canCrit bool, bonus int) phase.MoveResolveContext {
	return phase.MoveResolveContext{
		Spec: damage.Spec{
			ActorId:   "actor",
			MoveId:    1,
			Type:      ptype.Normal,
			Power:     40,
			Category:  move.DamageClassPhysical,
			CanCrit:   canCrit,
			CritBonus: bonus,
		},
		Battle: newCritBattle(),
	}
}

func critRate(t *testing.T, bonus, trials int) float64 {
	t.Helper()
	h := rule.NewCritHandler()
	ctx := newCritContext(true, bonus)
	hits := 0
	for range trials {
		if h.Mod(ctx).Crit != nil {
			hits++
		}
	}
	return float64(hits) / float64(trials)
}

func TestCritHandler_Mod(t *testing.T) {
	t.Run("CanCrit=false: Crit が nil", func(t *testing.T) {
		h := rule.NewCritHandler()
		ctx := newCritContext(false, 3)
		for range 300 {
			if mod := h.Mod(ctx); mod.Crit != nil {
				t.Fatalf("expected nil Crit, got %+v", *mod.Crit)
			}
		}
	})

	t.Run("CritBonus=3: 常に急所", func(t *testing.T) {
		h := rule.NewCritHandler()
		ctx := newCritContext(true, 3)
		for range 300 {
			if h.Mod(ctx).Crit == nil {
				t.Fatal("expected crit with CritBonus=3")
			}
		}
	})

	t.Run("CritBonus=3: Crit.Value=1.5", func(t *testing.T) {
		mod := rule.NewCritHandler().Mod(newCritContext(true, 3))
		if mod.Crit == nil {
			t.Fatal("expected non-nil Crit")
		}
		if mod.Crit.Value != 1.5 {
			t.Errorf("expected 1.5, got %v", mod.Crit.Value)
		}
	})

	t.Run("CritBonus=3: Op=ModOpOverride", func(t *testing.T) {
		mod := rule.NewCritHandler().Mod(newCritContext(true, 3))
		if mod.Crit == nil {
			t.Fatal("expected non-nil Crit")
		}
		if mod.Crit.Op != damage.ModOpOverride {
			t.Errorf("expected ModOpOverride, got %v", mod.Crit.Op)
		}
	})

	t.Run("CritBonus=3: Crit 以外の補正は nil", func(t *testing.T) {
		mod := rule.NewCritHandler().Mod(newCritContext(true, 3))
		mod.Crit = nil
		if mod != (damage.DamageMod{}) {
			t.Errorf("expected only Crit to be set, got %+v", mod)
		}
	})

	t.Run("CritBonus=5: 上限クランプで常に急所", func(t *testing.T) {
		h := rule.NewCritHandler()
		// stage=0 + bonus=5 → 5 → clamp 3（確定急所）
		ctx := newCritContext(true, 5)
		for range 300 {
			if h.Mod(ctx).Crit == nil {
				t.Fatal("expected crit with CritBonus=5")
			}
		}
	})

	t.Run("CritBonus=0: 急所率がおよそ 1/24 である", func(t *testing.T) {
		rate := critRate(t, 0, 3000)
		if rate < 0.02 || rate > 0.07 {
			t.Errorf("unexpected crit rate: %.3f (expected ~0.042)", rate)
		}
	})

	t.Run("CritBonus=1: 急所率がおよそ 1/8 である", func(t *testing.T) {
		rate := critRate(t, 1, 3000)
		if rate < 0.08 || rate > 0.17 {
			t.Errorf("unexpected crit rate: %.3f (expected ~0.125)", rate)
		}
	})

	t.Run("CritBonus=2: 急所率がおよそ 1/2 である", func(t *testing.T) {
		rate := critRate(t, 2, 3000)
		if rate < 0.40 || rate > 0.60 {
			t.Errorf("unexpected crit rate: %.3f (expected ~0.5)", rate)
		}
	})
}
