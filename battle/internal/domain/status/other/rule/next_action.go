package rule

import (
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
)

type ClearedOnNextAction struct {
	kind status.ClearedOnNextAction
}

func NewClearedOnNextAction() *ClearedOnNextAction {
	return &ClearedOnNextAction{kind: "cleared_on_next_action"}
}

func (c *ClearedOnNextAction) Resolve(_ other.OtherStatusContext) (bool, bool) {
	return true, false
}

func (c *ClearedOnNextAction) Kind() status.OtherCondition {
	return status.OtherCondition(c.kind)
}
