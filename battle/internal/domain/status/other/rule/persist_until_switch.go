package rule

import (
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
)

// PersistUntilSwitch は Resolve 経由では解除されない。
// 交代時の解除は Exit フェーズ側の責務。
type PersistUntilSwitch struct {
	kind status.PersistUntilSwitch
}

func NewPersistUntilSwitch() *PersistUntilSwitch {
	return &PersistUntilSwitch{kind: "persist_until_switch"}
}

func (c *PersistUntilSwitch) Resolve(_ other.OtherStatusContext) (bool, bool) {
	return false, false
}

func (c *PersistUntilSwitch) Kind() status.OtherCondition {
	return status.OtherCondition(c.kind)
}
