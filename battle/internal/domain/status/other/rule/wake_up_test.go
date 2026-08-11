package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/status/other/rule"
)

// TestClearedOnWakeUp_Resolve はスタブ状態のテスト。
// pokemon.MainStatus() が公開されたら以下の意図したロジックで実装・テストを更新する:
//   - ねむり中: cleared=false
//   - 目覚め済み(ねむり以外): cleared=true, addConfusion=false
func TestClearedOnWakeUp_Resolve(t *testing.T) {
	c := rule.NewClearedOnWakeUp()
	cleared, addConfusion := c.Resolve(other.OtherStatusContext{})
	// スタブ: 常に (false, false)
	if cleared {
		t.Error("stub: expected cleared=false until MainStatus() is available")
	}
	if addConfusion {
		t.Error("expected addConfusion=false")
	}
}
