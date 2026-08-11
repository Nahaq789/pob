package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/status/other/rule"
)

func TestClearedOnOtherCondition_Resolve(t *testing.T) {
	ctx := other.OtherStatusContext{}

	t.Run("turns=1 clears on first resolve", func(t *testing.T) {
		c := rule.NewClearedOnOtherCondition(1)
		cleared, addConfusion := c.Resolve(ctx)
		if !cleared {
			t.Error("expected cleared=true")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})

	t.Run("turns=2 does not clear on first resolve", func(t *testing.T) {
		c := rule.NewClearedOnOtherCondition(2)
		cleared, _ := c.Resolve(ctx)
		if cleared {
			t.Error("expected cleared=false on first resolve")
		}
		cleared, addConfusion := c.Resolve(ctx)
		if !cleared {
			t.Error("expected cleared=true on second resolve")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})
}
