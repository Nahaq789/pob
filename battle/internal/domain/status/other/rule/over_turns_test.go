package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/status/other/rule"
)

func TestClearedOverTurns_Resolve(t *testing.T) {
	ctx := other.OtherStatusContext{}

	t.Run("turns=1 clears on first resolve", func(t *testing.T) {
		c := rule.NewClearedOverTurns(1)
		cleared, addConfusion := c.Resolve(ctx)
		if !cleared {
			t.Error("expected cleared=true")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})

	t.Run("turns=2 does not clear on first resolve", func(t *testing.T) {
		c := rule.NewClearedOverTurns(2)
		cleared, addConfusion := c.Resolve(ctx)
		if cleared {
			t.Error("expected cleared=false on first resolve")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}

		cleared, addConfusion = c.Resolve(ctx)
		if !cleared {
			t.Error("expected cleared=true on second resolve")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})
}
