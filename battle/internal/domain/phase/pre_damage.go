package phase

import (
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other/rule"
)

type PreDamagePhaseHandler struct{}

func NewPreDamagePhaseHandler() *PreDamagePhaseHandler {
	return &PreDamagePhaseHandler{}
}

func (pre *PreDamagePhaseHandler) Handle(ctx PreDamageContext) Result {
	actor := ctx.Battle.PlayerById(ctx.ActorId)
	activeP := actor.Active()
	var messages []string

	// 状態異常判定
	main := activeP.Status().Main()
	switch main.Condition() {
	case status.Sleep:
		activeP.DecrementMainStatusCount()
		message, ok := main.IsSleep(activeP.Name())
		messages = append(messages, message)
		if ok {
			return Result{
				Messages:  messages,
				NextPhase: PhaseEnd,
			}
		}
	case status.Freeze:
		activeP.DecrementMainStatusCount()
		message, ok := main.IsFreeze(activeP.Name())
		messages = append(messages, message)
		if ok {
			return Result{
				Messages:  messages,
				NextPhase: PhaseEnd,
			}
		}
	default:
	}

	otherMap := activeP.Status().OtherMap()

	// ひるみ判定
	if fl, ok := otherMap[status.OtherCondition(status.Flinch)].(*rule.Flinch); ok {
		messages = append(messages, fl.Handle(activeP.Name()))
		return Result{Messages: messages, NextPhase: PhaseEnd}
	}

	// アンコール判定
	if enc, ok := otherMap[status.OtherCondition(status.Encore)].(*rule.Encore); ok {
		if message, blocked := enc.Handle(ctx.MoveId, activeP.Status().Others()); blocked {
			messages = append(messages, message)
			return Result{Messages: messages, NextPhase: PhaseEnd}
		}
	}

	// かなしばり判定
	if md, ok := otherMap[status.OtherCondition(status.MoveDisabled)].(*rule.MoveDisabled); ok {
		if message, blocked := md.Handle(ctx.MoveId, activeP.Status().Others()); blocked {
			messages = append(messages, message)
			return Result{Messages: messages, NextPhase: PhaseEnd}
		}
	}

	// 混乱判定
	if confuse, ok := otherMap[status.OtherCondition(status.Confusion)].(*rule.Confusion); ok {
		cleared, _, message := confuse.Resolve(status.OtherStatusContext{ActorName: activeP.Name()})
		if cleared {
			activeP.RemoveOtherStatus(status.OtherCondition(status.Confusion))
			if message != "" {
				messages = append(messages, message)
			}
		}
		if !cleared {
			if msg, hit := confuse.CheckSelfHit(activeP.Name()); hit {
				messages = append(messages, msg)
				return Result{Messages: messages, NextPhase: PhaseEnd}
			}
		}
	}

	// まひ状態であれば、このタイミングで麻痺の判定をおこなう
	if main.IsParalysis() {
		if message, p := main.CheckParalysis(activeP.Name()); p {
			messages = append(messages, message)
			return Result{
				Messages:  messages,
				NextPhase: PhaseEnd,
			}
		}
	}

	return Result{
		Messages:  messages,
		NextPhase: PhaseDamage,
	}
}
