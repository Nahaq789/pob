package rule

import (
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
)

// Rampage はあばれるけい（Gen5以降ルール）の状態変化を管理する。
// turns は技側が乱数(2-3等)で決定した値を渡す。
type Rampage struct {
	kind      status.Rampage
	remaining int
}

func NewRampage(turns int) *Rampage {
	return &Rampage{kind: "rampage", remaining: turns}
}

func (r *Rampage) Resolve(ctx other.OtherStatusContext) (bool, bool) {
	if ctx.ActionSucceeded {
		r.remaining--
		if r.remaining <= 0 {
			// 規定ターン分暴れ終わった → 解除 + こんらん付与
			return true, true
		}
		return false, false
	}

	// 行動失敗: 即座に解除。残りターンが1だった(＝本来の最終ターン)場合はこんらん付与。
	addConfusion := r.remaining == 1
	return true, addConfusion
}

func (r *Rampage) Kind() status.OtherCondition {
	return status.OtherCondition(r.kind)
}
