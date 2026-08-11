package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/status/other/rule"
)

func TestClearedOnMoveTurnEnd_Resolve(t *testing.T) {
	c := rule.NewClearedOnMoveTurnEnd()
	cleared, addConfusion := c.Resolve(other.OtherStatusContext{})
	if !cleared {
		t.Error("expected cleared=true")
	}
	if addConfusion {
		t.Error("expected addConfusion=false")
	}
}
