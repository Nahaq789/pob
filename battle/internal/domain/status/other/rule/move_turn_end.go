package rule

import "pob/battle/internal/domain/status/other"

type ClearedOnMoveTurnEnd struct{}

func NewClearedOnMoveTurnEnd() *ClearedOnMoveTurnEnd {
	return &ClearedOnMoveTurnEnd{}
}

func (c *ClearedOnMoveTurnEnd) Resolve(_ other.OtherStatusContext) (bool, bool) {
	return true, false
}
