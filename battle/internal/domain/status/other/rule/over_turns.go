package rule

import "pob/battle/internal/domain/status/other"

type ClearedOverTurns struct {
	remaining int
}

func NewClearedOverTurns(turns int) *ClearedOverTurns {
	return &ClearedOverTurns{remaining: turns}
}

func (c *ClearedOverTurns) Resolve(_ other.OtherStatusContext) (bool, bool) {
	c.remaining--
	if c.remaining <= 0 {
		return true, false
	}
	return false, false
}
