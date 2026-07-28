package phase

import (
	"fmt"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/player"
)

type ActionResolvePhaseHandler struct {
	exitHandler  *ExitPhaseHandler
	entryHandler *EntryPhaseHandler
}

func NewActionResolvePhaseHandler(exit *ExitPhaseHandler, entry *EntryPhaseHandler) *ActionResolvePhaseHandler {
	return &ActionResolvePhaseHandler{exitHandler: exit, entryHandler: entry}
}

type pendingSwitch struct {
	player *player.Player
	req    *player.SwitchRequest
}

func (ar *ActionResolvePhaseHandler) Handle(b *battle.Battle) (map[string][]string, error) {
	players := []*player.Player{b.Player1(), b.Player2()}

	// プレイヤーの選択した行動チェック
	if err := ar.validatePendingActions(players); err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	for _, pl := range players {
		if pl.HasPendingForfeit() {
			pl.PullPendingForfeit()
			b.SetWinner(b.Opponent(pl))
			result[pl.Id()] = append(result[pl.Id()], fmt.Sprintf("プレイヤー%sが降参しました。", pl.Id()))
			return result, nil
		}
	}

	// 交代
	exiting, pendings := ar.collectPendingSwitches(players)
	if len(exiting) > 0 {
		exitMsgs, err := ar.exitHandler.Handle(exiting, b)
		if err != nil {
			return nil, err
		}
		ar.mergeMessages(result, exitMsgs)

		entered := ar.commitSwitches(pendings)
		entryMsgs, err := ar.entryHandler.Handle(entered, b)
		if err != nil {
			return nil, err
		}

		ar.mergeMessages(result, entryMsgs)
	}

	// 技
	for _, pl := range players {
		if req := pl.PullPendingMove(); req != nil {
			b.PushPendingMove(*req)
		}
	}

	return result, nil
}

// collectPendingSwitches は各プレイヤーの交代リクエストを取り出し、
// ExitPhase に渡す exiting リストと、スロット確定用の pendings を返す。
func (ar *ActionResolvePhaseHandler) collectPendingSwitches(players []*player.Player) ([]ExitingPokemon, []pendingSwitch) {
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
func (ar *ActionResolvePhaseHandler) commitSwitches(pendings []pendingSwitch) []EnteredPokemon {
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

func (ar *ActionResolvePhaseHandler) validatePendingActions(players []*player.Player) error {
	for _, pl := range players {
		count := 0
		if pl.HasPendingSwitch() {
			count++
		}

		if pl.HasPendingMove() {
			count++
		}

		if pl.HasPendingForfeit() {
			count++
		}

		switch count {
		case 0:
			return fmt.Errorf("player %s has no pending action", pl.Id())
		case 1:
		default:
			return fmt.Errorf("player %s has multiple pending actions", pl.Id())
		}
	}
	return nil
}

func (ar *ActionResolvePhaseHandler) mergeMessages(dst, src map[string][]string) {
	for k, v := range src {
		dst[k] = append(dst[k], v...)
	}
}
