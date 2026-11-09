package rule_test

import (
	"testing"

	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/damage/rule"
	"pob/battle/internal/domain/ptype"
)

func TestStabHandler_Mod(t *testing.T) {
	single := func(tp ptype.Type) [2]ptype.Type { return [2]ptype.Type{tp, ptype.None} }

	tests := []struct {
		name          string
		moveType      ptype.Type
		attackerTypes [2]ptype.Type
		wantStab      bool
	}{
		{"単タイプで一致: 1.5倍", ptype.Fire, single(ptype.Fire), true},
		{"複合タイプの第1タイプと一致: 1.5倍", ptype.Water, [2]ptype.Type{ptype.Water, ptype.Ground}, true},
		{"複合タイプの第2タイプと一致: 1.5倍", ptype.Ground, [2]ptype.Type{ptype.Water, ptype.Ground}, true},
		{"不一致: 補正なし", ptype.Electric, [2]ptype.Type{ptype.Water, ptype.Ground}, false},
		// 単タイプの [1] は ptype.None のため、技タイプ None と誤って一致しないことを確認する。
		{"タイプなしの技を単タイプが使用: 補正なし", ptype.None, single(ptype.Normal), false},
		{"タイプなしの技を複合タイプが使用: 補正なし", ptype.None, [2]ptype.Type{ptype.Water, ptype.Ground}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTypeTestContext(tt.moveType, tt.attackerTypes, single(ptype.Normal))
			got := rule.NewStabHandler().Mod(ctx).Stab
			if !tt.wantStab {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("expected Stab to be set, got nil")
			}
			if got.Value != 1.5 {
				t.Errorf("expected 1.5, got %v", got.Value)
			}
		})
	}

	t.Run("一致: Op が Multiply", func(t *testing.T) {
		ctx := newTypeTestContext(ptype.Fire, single(ptype.Fire), single(ptype.Normal))
		got := rule.NewStabHandler().Mod(ctx).Stab
		if got == nil {
			t.Fatal("expected Stab to be set, got nil")
		}
		if got.Op != damage.ModOpMultiply {
			t.Errorf("expected ModOpMultiply, got %v", got.Op)
		}
	})

	t.Run("防御側のタイプ: 一致判定に影響しない", func(t *testing.T) {
		ctx := newTypeTestContext(ptype.Fire, single(ptype.Water), single(ptype.Fire))
		if got := rule.NewStabHandler().Mod(ctx).Stab; got != nil {
			t.Errorf("expected nil, got %+v", got)
		}
	})

	t.Run("Stab 以外のフィールド: nil", func(t *testing.T) {
		ctx := newTypeTestContext(ptype.Fire, single(ptype.Fire), single(ptype.Normal))
		mod := rule.NewStabHandler().Mod(ctx)
		mod.Stab = nil
		if mod != (damage.DamageMod{}) {
			t.Errorf("expected other fields nil, got %+v", mod)
		}
	})
}
