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
	default:
	}

	// 混乱判定
	for _, s := range activeP.Status().Others() {
		if confuse, ok := s.(*rule.Confusion); ok {
			if message, h := confuse.CheckSelfHit(activeP.Name()); h {
				return Result{
					Messages:  []string{message},
					NextPhase: PhaseEnd,
				}
			}
		}
	}

	// こんらんを通り抜けた場合
	// まひ状態であれば、このタイミングで麻痺の判定をおこなう
	if main.IsParalysis() {
		if message, p := main.CheckParalysis(activeP.Name()); p {
			return Result{
				Messages:  []string{message},
				NextPhase: PhaseEnd,
			}
		}
	}

	return Result{
		Messages:  messages,
		NextPhase: PhaseDamage,
	}
}
