package rule

import (
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
)

// ClearedOnOtherCondition はカウントダウン式の解除ルール。ClearedOverTurns と同型だが、
// 「別の特定条件が整うまで」という分類上の区別のために独立した構造体として定義する。
type ClearedOnOtherCondition struct {
	kind      status.ClearedOnOtherCondition
	remaining int
}

func NewClearedOnOtherCondition(turns int) *ClearedOnOtherCondition {
	return &ClearedOnOtherCondition{kind: "cleared_on_other_condition", remaining: turns}
}

func (c *ClearedOnOtherCondition) Resolve(_ other.OtherStatusContext) (bool, bool) {
	c.remaining--
	if c.remaining <= 0 {
		return true, false
	}
	return false, false
}

func (c *ClearedOnOtherCondition) Kind() status.OtherCondition {
	return status.OtherCondition(c.kind)
}
