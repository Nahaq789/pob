package phase

import (
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/player"
)

type PreMoveSelectPhaseHandler struct {
	exitHandler  *ExitPhaseHandler
	entryHandler *EntryPhaseHandler
}

func NewPreMoveSelectPhaseHandler(exit *ExitPhaseHandler, entry *EntryPhaseHandler) *PreMoveSelectPhaseHandler {
	return &PreMoveSelectPhaseHandler{exitHandler: exit, entryHandler: entry}
}

type pendingSwitch struct {
	player *player.Player
	req    *player.SwitchRequest
}

func (p *PreMoveSelectPhaseHandler) Handle(b *battle.Battle) (map[string][]string, error) {
	players := []*player.Player{b.Player1(), b.Player2()}
	exiting, pendings := collectPendingSwitches(players)
	if len(exiting) == 0 {
		return nil, nil
	}

	result := make(map[string][]string)

	exitMsgs, err := p.exitHandler.Handle(exiting, b)
	if err != nil {
		return nil, err
	}
	for k, v := range exitMsgs {
		result[k] = append(result[k], v...)
	}

	entered := commitSwitches(pendings)
	entryMsgs, err := p.entryHandler.Handle(entered, b)
	if err != nil {
		return nil, err
	}
	for k, v := range entryMsgs {
		result[k] = append(result[k], v...)
	}

	return result, nil
}

// collectPendingSwitches は各プレイヤーの交代リクエストを取り出し、
// ExitPhase に渡す exiting リストと、スロット確定用の pendings を返す。
func collectPendingSwitches(players []*player.Player) ([]ExitingPokemon, []pendingSwitch) {
	var exiting []ExitingPokemon
	var pendings []pendingSwitch
	for _, pl := range players {
		req := pl.PullPendingSwitch()
		if req == nil {
			continue
		}
		pendings = append(pendings, pendingSwitch{pl, req})
		exiting = append(exiting, ExitingPokemon{PlayerId: pl.Id(), Pokemon: req.Outgoing, Incoming: req.Incoming})
	}
	return exiting, pendings
}

// commitSwitches は各プレイヤーのアクティブスロットを確定させ、
// EntryPhase に渡す entered リストを返す。
func commitSwitches(pendings []pendingSwitch) []EnteredPokemon {
	entered := make([]EnteredPokemon, 0, len(pendings))
	for _, pd := range pendings {
		pd.player.SetActiveSlot(pd.req.IncomingIndex)
		entered = append(entered, EnteredPokemon{
			PlayerId: pd.player.Id(),
			Pokemon:  pd.req.Incoming,
		})
	}
	return entered
}
