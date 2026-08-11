package phase

import "pob/battle/internal/domain/battle"

type PreDamagePhaseHandler struct {
}

func NewPreDamagePhaseHandler() *PreDamagePhaseHandler {
	return &PreDamagePhaseHandler{}
}

func (pre *PreDamagePhaseHandler) Handle(b *battle.Battle) Result {
	return Result{}
}
