package expiry

import (
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
)

type ClearedOverTurns struct {
	kind      status.ClearedOverTurns
	remaining int
}

func NewClearedOverTurns(turns int) *ClearedOverTurns {
	return &ClearedOverTurns{kind: "cleared_over_turns", remaining: turns}
}

func (c *ClearedOverTurns) Resolve(_ other.OtherStatusContext) (bool, bool) {
	c.remaining--
	if c.remaining <= 0 {
		return true, false
	}
	return false, false
}

func (c *ClearedOverTurns) Kind() status.OtherCondition {
	return status.OtherCondition(c.kind)
}
