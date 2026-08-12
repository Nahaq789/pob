package rule_test

import (
	"testing"

	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/status/other/rule"
)

func TestClearedOnElectricMove_Resolve(t *testing.T) {
	t.Run("でんきタイプ技: cleared=true", func(t *testing.T) {
		c := rule.NewClearedOnElectricMove()
		cleared, addConfusion := c.Resolve(other.OtherStatusContext{MoveType: ptype.Electric})
		if !cleared {
			t.Error("expected cleared=true for electric move")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})

	t.Run("でんき以外の技: cleared=false", func(t *testing.T) {
		c := rule.NewClearedOnElectricMove()
		cleared, addConfusion := c.Resolve(other.OtherStatusContext{MoveType: ptype.Fire})
		if cleared {
			t.Error("expected cleared=false for non-electric move")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})
}
