package rule

import (
	"fmt"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/status"
	statusother "pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/vo"
)

// MoveDisabled はかなしばり状態を管理する。
// 直前に使用した技を指定ターン数の間使用不能にする状態異常。
// remaining が 0 になると Resolve が cleared=true を返し、状態が解除される。
// 封じられた技の選択可否は Handle で判定する。
type MoveDisabled struct {
	kind      status.ClearedOverTurns
	remaining vo.Count
	moveId    int
}

func NewMoveDisabled(turns int, moveId int) *MoveDisabled {
	return &MoveDisabled{
		kind:      status.MoveDisabled,
		remaining: vo.NewCount(turns),
		moveId:    moveId,
	}
}

func (m *MoveDisabled) Resolve(ctx status.OtherStatusContext) (bool, bool, string) {
	m.remaining = m.remaining.Decrement()
	if m.remaining.IsEmpty() {
		return true, false, fmt.Sprintf("%sへのかなしばりが解除された", ctx.ActorName)
	}
	return false, false, ""
}

func (m *MoveDisabled) Kind() status.OtherCondition {
	return status.OtherCondition(m.kind)
}

func (m *MoveDisabled) MoveId() int { return m.moveId }

func (m *MoveDisabled) BlocksMoveId(moveId int) bool { return m.moveId == moveId }

// Handle は技選択時にかなしばりによる使用制限を判定する。
// 封じた技がアンコールで強制されている場合はわるあがき以外を弾く。
// それ以外は封じられた技の選択を弾く。
func (m *MoveDisabled) Handle(selectedMoveId int, others []status.OtherStatus) (string, bool) {
	if isDisabledMoveForced(others, m.moveId) {
		if selectedMoveId != move.StruggleId {
			return "アンコールされた技がかなしばりで出せない", true
		}
		return "", false
	}
	if selectedMoveId == m.moveId {
		return "かなしばりで技が出せない", true
	}
	return "", false
}

func isDisabledMoveForced(others []status.OtherStatus, moveId int) bool {
	for _, o := range others {
		if forcer, ok := o.(statusother.MoveForcer); ok && forcer.ForcesMoveId() == moveId {
			return true
		}
	}
	return false
}
