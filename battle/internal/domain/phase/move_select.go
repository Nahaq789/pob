package phase

import (
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/player"
)

type MoveSelectPhaseHandler struct{}

func (m *MoveSelectPhaseHandler) Handle(req player.MoveRequest, b *battle.Battle) {
	b.PushPendingMove(req)
}
