package phase

import (
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/status"
	statusother "pob/battle/internal/domain/status/other"
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
	if main != nil {
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
	}

	otherMap := activeP.Status().OtherMap()

	// ひるみ判定
	if fl, ok := otherMap[status.OtherCondition(status.Flinch)].(statusother.Flincher); ok {
		messages = append(messages, fl.Handle(activeP.Name()))
		return Result{Messages: messages, NextPhase: PhaseEnd}
	}

	// アンコール判定
	if enc, ok := otherMap[status.OtherCondition(status.Encore)].(statusother.PreMoveChecker); ok {
		if message, blocked := enc.Handle(ctx.MoveId, activeP.Status().Others()); blocked {
			messages = append(messages, message)
			return Result{Messages: messages, NextPhase: PhaseEnd}
		}
	}

	// かなしばり判定
	if md, ok := otherMap[status.OtherCondition(status.MoveDisabled)].(statusother.PreMoveChecker); ok {
		if message, blocked := md.Handle(ctx.MoveId, activeP.Status().Others()); blocked {
			messages = append(messages, message)
			return Result{Messages: messages, NextPhase: PhaseEnd}
		}
	}

	// 混乱判定
	if confuseStatus, ok := otherMap[status.OtherCondition(status.Confusion)]; ok {
		cleared, _, message := confuseStatus.Resolve(status.OtherStatusContext{ActorName: activeP.Name()})
		if cleared {
			activeP.RemoveOtherStatus(status.OtherCondition(status.Confusion))
			if message != "" {
				messages = append(messages, message)
			}
		}
		if !cleared {
			if selfHitter, ok := confuseStatus.(statusother.SelfHitter); ok {
				if msg, hit := selfHitter.CheckSelfHit(activeP.Name()); hit {
					messages = append(messages, msg)
					spec := &damage.Spec{
						ActorId:    ctx.ActorId,
						Type:       ptype.None,
						Power:      40,
						Category:   move.DamageClassPhysical,
						MustHit:    true,
						CanCrit:    false,
						TargetSelf: true,
					}
					return Result{Messages: messages, NextPhase: PhaseDamage, DamageContext: spec}
				}
			}
		}
	}

	// まひ状態であれば、このタイミングで麻痺の判定をおこなう
	if main != nil && main.IsParalysis() {
		if message, p := main.CheckParalysis(activeP.Name()); p {
			messages = append(messages, message)
			return Result{
				Messages:  messages,
				NextPhase: PhaseEnd,
			}
		}
	}

	m, err := activeP.MoveById(ctx.MoveId)
	if err != nil {
		return Result{Err: err}
	}
	spec := &damage.Spec{
		ActorId:  ctx.ActorId,
		MoveId:   ctx.MoveId,
		Type:     m.Type(),
		Power:    m.Power(),
		Category: m.DamageClass(),
		MustHit:  m.Accuracy() == 0,
		CanCrit:  true,
	}
	return Result{
		Messages:      messages,
		NextPhase:     PhaseDamage,
		DamageContext: spec,
	}
}
