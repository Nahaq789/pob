package expiry

import (
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
)

type ClearedOnElectricMove struct {
	kind status.ClearedOnElectricMove
}

func NewClearedOnElectricMove() *ClearedOnElectricMove {
	return &ClearedOnElectricMove{kind: "cleared_on_electric_move"}
}

func (c *ClearedOnElectricMove) Resolve(ctx other.OtherStatusContext) (bool, bool) {
	if ctx.MoveType == ptype.Electric {
		return true, false
	}
	return false, false
}

func (c *ClearedOnElectricMove) Kind() status.OtherCondition {
	return status.OtherCondition(c.kind)
}
