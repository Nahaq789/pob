package phase

import (
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other/rule"
)

type PreDamagePhaseHandler struct {
}

func NewPreDamagePhaseHandler() *PreDamagePhaseHandler {
	return &PreDamagePhaseHandler{}
}

func (pre *PreDamagePhaseHandler) Handle(ctx PreDamageContext) Result {
	actor := ctx.Battle.PlayerById(ctx.ActorId)
	activeP := actor.Active()

	// かなしばり判定
	for _, s := range activeP.Status().Others() {
		if md, ok := s.(*rule.MoveDisabled); ok {
			if message, same := md.Handle(ctx.MoveId); same {
				return Result{
					Message:   message,
					NextPhase: PhaseEnd,
					Err:       nil,
				}
			}
		}
	}

	// 状態異常判定
	main := activeP.Status().Main()
	switch main.Condition() {
	case status.Sleep:
		activeP.DecrementMainStatusCount()
		if message, ok := main.IsSleep(activeP.Name()); ok {
			return Result{
				Message:   message,
				NextPhase: PhaseEnd,
				Err:       nil,
			}
		}
	}

	return Result{}
}
