package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/status/other/rule"
)

func TestClearedOnWakeUp_Resolve(t *testing.T) {
	t.Run("ねむり中: cleared=false", func(t *testing.T) {
		c := rule.NewClearedOnWakeUp()
		cleared, addConfusion := c.Resolve(other.OtherStatusContext{MainCondition: status.Sleep})
		if cleared {
			t.Error("expected cleared=false while sleeping")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})

	t.Run("目覚め済み: cleared=true", func(t *testing.T) {
		c := rule.NewClearedOnWakeUp()
		cleared, addConfusion := c.Resolve(other.OtherStatusContext{MainCondition: status.None})
		if !cleared {
			t.Error("expected cleared=true after waking up")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})
}
