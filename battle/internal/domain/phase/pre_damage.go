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
	var messages []string

	// ひるみ
	// アンコール
	// 判定も追加する必要あり

	// かなしばり判定
	for _, s := range activeP.Status().Others() {
		if md, ok := s.(*rule.MoveDisabled); ok {
			if message, same := md.Handle(ctx.MoveId); same {
				return Result{
					Messages:  []string{message},
					NextPhase: PhaseEnd,
				}
			}
		}
	}

	// 状態異常判定
	main := activeP.Status().Main()
	switch main.Condition() {
	case status.Sleep:
		activeP.DecrementMainStatusCount()
		message, ok := main.IsSleep(activeP.Name())
		if ok {
			return Result{
				Messages:  []string{message},
				NextPhase: PhaseEnd,
			}
		}
		messages = append(messages, message)
	case status.Freeze:
		activeP.DecrementMainStatusCount()
		message, ok := main.IsFreeze(activeP.Name())
		if ok {
			return Result{
				Messages:  []string{message},
				NextPhase: PhaseEnd,
			}
		}
		messages = append(messages, message)
	}

	// 混乱判定
	// todo
	for _, s := range activeP.Status().Others() {
		confuse, ok := s.(*rule.Confusion)
		if !ok {
			return Result{
				Messages:  messages,
				NextPhase: PhaseDamage,
			}
		}

		if ok := confuse.CheckSelfHit(); ok {

		}

	}

	return Result{
		Messages:  messages,
		NextPhase: PhaseDamage,
	}
}
