package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/status/other/rule"
)

// TestClearedOnElectricMove_Resolve はスタブ状態のテスト。
// ptype の依存制約が解消されたら以下の意図したロジックで実装・テストを更新する:
//   - でんきタイプ技: cleared=true, addConfusion=false
//   - でんき以外の技: cleared=false, addConfusion=false
func TestClearedOnElectricMove_Resolve(t *testing.T) {
	c := rule.NewClearedOnElectricMove()
	cleared, addConfusion := c.Resolve(other.OtherStatusContext{MoveId: "some-move"})
	// スタブ: 常に (false, false)
	if cleared {
		t.Error("stub: expected cleared=false until ptype dependency is resolved")
	}
	if addConfusion {
		t.Error("expected addConfusion=false")
	}
}
