package phase

import (
	"fmt"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/player"
)

type ForfeitPhaseHandler struct{}

func (f *ForfeitPhaseHandler) Handle(pl *player.Player, b *battle.Battle) string {
	b.SetWinner(b.Opponent(pl))
	return fmt.Sprintf("プレイヤー%sが降参しました。", pl.Id())
}
