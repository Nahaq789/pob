package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/status/other/rule"
)

func TestPersistUntilSwitch_Resolve(t *testing.T) {
	c := rule.NewPersistUntilSwitch()
	for i := 0; i < 5; i++ {
		cleared, addConfusion := c.Resolve(other.OtherStatusContext{})
		if cleared {
			t.Errorf("turn %d: expected cleared=false (only switch-out clears this)", i+1)
		}
		if addConfusion {
			t.Errorf("turn %d: expected addConfusion=false", i+1)
		}
	}
}
