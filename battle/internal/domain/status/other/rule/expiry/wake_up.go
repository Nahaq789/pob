package expiry

import (
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
)

type ClearedOnWakeUp struct {
	kind status.ClearedOnWakeUp
}

func NewClearedOnWakeUp() *ClearedOnWakeUp {
	return &ClearedOnWakeUp{kind: "cleared_on_wake_up"}
}

func (c *ClearedOnWakeUp) Resolve(ctx other.OtherStatusContext) (bool, bool) {
	if ctx.MainCondition != status.Sleep {
		return true, false
	}
	return false, false
}

func (c *ClearedOnWakeUp) Kind() status.OtherCondition {
	return status.OtherCondition(c.kind)
}
