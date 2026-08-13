package rule

import (
	"math/rand"
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/vo"
)

// Confusion はこんらん状態を管理する。
// 行動時に 1/3 の確率で自分自身に攻撃してしまう状態異常。
// 指定ターン数が経過すると Resolve が cleared=true を返し、状態が解除される。
// 自傷判定は CheckSelfHit で行う。
type Confusion struct {
	kind      status.Rampage
	remaining vo.Count
}

func NewConfusion(turns int) *Confusion {
	return &Confusion{kind: "confusion", remaining: vo.NewCount(turns)}
}

func (c *Confusion) Resolve(_ other.OtherStatusContext) (bool, bool) {
	c.remaining = c.remaining.Decrement()
	if c.remaining.IsEmpty() {
		return true, false
	}
	return false, false
}

func (c *Confusion) Kind() status.OtherCondition { return status.OtherCondition(c.kind) }

// CheckSelfHit はこんらんによる自傷が発生するかを判定する。
// 1/3 の確率で true を返す。技使用前に呼び出し、true の場合は技をキャンセルして自傷ダメージを与える。
func (c *Confusion) CheckSelfHit() bool {
	return rand.Intn(3) == 0
}
