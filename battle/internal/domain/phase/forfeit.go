package phase

import (
	"fmt"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/player"
)

type ForfeitPhaseHandler struct{}

func (f *ForfeitPhaseHandler) Handle(pl *player.Player, b *battle.Battle) map[string]string {
	opponent := b.Opponent(pl)
	b.SetWinner(opponent)
	result := make(map[string]string)
	result[pl.Id()] = fmt.Sprintf("プレイヤー%sが降参しました。", pl.Name())
	result[opponent.Id()] = fmt.Sprintf("プレイヤー%sが降参しました。%sの勝利です！", pl.Name(), opponent.Name())
	return result
}

func (f *ForfeitPhaseHandler) HandleDraw(players []*player.Player, b *battle.Battle) map[string]string {
	b.Draw()
	msg := "両者が降参しました。引き分けです。"
	result := make(map[string]string)
	for _, pl := range players {
		result[pl.Id()] = msg
	}
	return result
}
