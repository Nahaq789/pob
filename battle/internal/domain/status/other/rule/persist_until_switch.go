package rule

import "pob/battle/internal/domain/status/other"

// PersistUntilSwitch は Resolve 経由では解除されない。
// 交代時の解除は Exit フェーズ側の責務。
type PersistUntilSwitch struct{}

func NewPersistUntilSwitch() *PersistUntilSwitch {
	return &PersistUntilSwitch{}
}

func (c *PersistUntilSwitch) Resolve(_ other.OtherStatusContext) (bool, bool) {
	return false, false
}
