package rule

import (
	"fmt"
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/vo"
)

// Encore はアンコール状態を管理する。
// 直前に使用した技を3ターンの間強制する状態異常。
// remaining が 0 になると Resolve が cleared=true を返し、状態が解除される。
// 強制技の選択可否は Handle で判定する。
type Encore struct {
	kind      status.ClearedOverTurns
	remaining vo.Count
	moveId    int
}

func NewEncore(moveId int) *Encore {
	return &Encore{
		kind:      status.Encore,
		remaining: vo.NewCount(3),
		moveId:    moveId,
	}
}

func (e *Encore) Resolve(ctx status.OtherStatusContext) (bool, bool, string) {
	e.remaining = e.remaining.Decrement()
	if e.remaining.IsEmpty() {
		return true, false, fmt.Sprintf("%sのアンコールが解除された", ctx.ActorName)
	}
	return false, false, ""
}

func (e *Encore) Kind() status.OtherCondition {
	return status.OtherCondition(e.kind)
}

func (e *Encore) MoveId() int { return e.moveId }

// Handle は技選択時にアンコールによる使用制限を判定する。
// selectedMoveId がアンコールされた技と一致しない場合、メッセージと true を返す。
func (e *Encore) Handle(selectedMoveId int) (string, bool) {
	if selectedMoveId != e.moveId {
		return "アンコールでその技は出せない", true
	}
	return "", false
}
