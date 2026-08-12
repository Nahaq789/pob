package rule

import (
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
)

type ClearedOnMoveTurnEnd struct {
	kind status.ClearedOnMoveTurnEnd
}

func NewClearedOnMoveTurnEnd() *ClearedOnMoveTurnEnd {
	return &ClearedOnMoveTurnEnd{kind: "cleared_on_move_turn_end"}
}

func (c *ClearedOnMoveTurnEnd) Resolve(_ other.OtherStatusContext) (bool, bool) {
	return true, false
}

func (c *ClearedOnMoveTurnEnd) Kind() status.OtherCondition {
	return status.OtherCondition(c.kind)
}
