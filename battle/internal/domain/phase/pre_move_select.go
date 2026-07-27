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

func (p *PreMoveSelectPhaseHandler) Handle(b *battle.Battle) (map[string][]string, error) {
	players := []*player.Player{b.Player1(), b.Player2()}
	var exiting []ExitingPokemon
	type pending struct {
		player *player.Player
		req    *player.SwitchRequest
	}
	var pendings []pending

	for _, player := range players {
		request := player.PullPendingSwitch()
		if request == nil {
			continue
		}

		pendings = append(pendings, pending{player, request})
		exiting = append(exiting, ExitingPokemon{PlayerId: player.Id(), Pokemon: request.Outgoing, Incoming: request.Incoming})
	}

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

	var entered []EnteredPokemon
	for _, pd := range pendings {
		pd.player.SetActiveSlot(pd.req.IncomingIndex)
		entered = append(entered, EnteredPokemon{
			PlayerId: pd.player.Id(),
			Pokemon:  pd.req.Incoming,
		})
	}

	entryMsgs, err := p.entryHandler.Handle(entered, b)
	if err != nil {
		return nil, err
	}

	for k, v := range entryMsgs {
		result[k] = append(result[k], v...)
	}

	return result, nil
}
