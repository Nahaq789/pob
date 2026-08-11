package rule

import "pob/battle/internal/domain/status/other"

type ClearedOnNextAction struct{}

func NewClearedOnNextAction() *ClearedOnNextAction {
	return &ClearedOnNextAction{}
}

func (c *ClearedOnNextAction) Resolve(_ other.OtherStatusContext) (bool, bool) {
	return true, false
}
