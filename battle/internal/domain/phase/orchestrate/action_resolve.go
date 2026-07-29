package orchestrate

import (
	"fmt"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/phase"
	"pob/battle/internal/domain/player"
	"sort"
)

// ActionResolveOrchestrator は各プレイヤーの保留行動を解決するオーケストレーター。
// 個々のフェーズハンドラー（ExitPhaseHandler・EntryPhaseHandler）を呼び出し、
// バトル1ターン分の行動解決フローを制御する。
type ActionResolveOrchestrator struct {
	exitHandler  *phase.ExitPhaseHandler
	entryHandler *phase.EntryPhaseHandler
}

func NewActionResolveOrchestrator(exit *phase.ExitPhaseHandler, entry *phase.EntryPhaseHandler) *ActionResolveOrchestrator {
	return &ActionResolveOrchestrator{exitHandler: exit, entryHandler: entry}
}

// pendingSwitch はプレイヤーと交代リクエストのペア。commitSwitch に渡す内部型。
type pendingSwitch struct {
	player *player.Player
	req    *player.SwitchRequest
}

// switchEntry は素早さ順ソートの対象となる交代エントリ。
// ExitPhase と commitSwitch それぞれに必要な情報をひとまとめにする。
type switchEntry struct {
	exiting phase.ExitingPokemon
	pending pendingSwitch
}

// Handle は各プレイヤーの保留行動（降参・交代・技）を解決する。
// 降参はバトルを即座に終了させる。交代は素早さ順で退場→入場を処理する。
// 技は PendingMove としてバトルに積み、後続フェーズで処理される。
// 戻り値はプレイヤーIDごとの発動メッセージ一覧。
func (ar *ActionResolveOrchestrator) Handle(b *battle.Battle) (map[string][]string, error) {
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
	entries := ar.collectPendingSwitches(players)
	if len(entries) > 0 {
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].exiting.Pokemon.Speed() > entries[j].exiting.Pokemon.Speed()
		})
		for _, entry := range entries {
			ep := entry.exiting
			// 退場処理
			result[ep.PlayerId] = append(result[ep.PlayerId], fmt.Sprintf("戻れ！%s！", ep.Pokemon.Name()))
			exitMsgs, err := ar.exitHandler.Handle(ep, b)
			if err != nil {
				return nil, err
			}
			result[ep.PlayerId] = append(result[ep.PlayerId], exitMsgs...)

			// 入場処理
			entered := ar.commitSwitch(entry.pending)
			result[ep.PlayerId] = append(result[ep.PlayerId], fmt.Sprintf("行け！%s！", entered.Pokemon.Name()))
			entryMsgs, err := ar.entryHandler.Handle(entered, b)
			if err != nil {
				return nil, err
			}
			result[ep.PlayerId] = append(result[ep.PlayerId], entryMsgs...)
		}
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
// 素早さ順ソートと commitSwitch に必要な情報をまとめた switchEntry スライスを返す。
func (ar *ActionResolveOrchestrator) collectPendingSwitches(players []*player.Player) []switchEntry {
	var entries []switchEntry

	for _, pl := range players {
		req := pl.PullPendingSwitch()
		if req == nil {
			continue
		}
		entries = append(entries, switchEntry{
			exiting: phase.ExitingPokemon{PlayerId: pl.Id(), Pokemon: req.Outgoing, Incoming: req.Incoming},
			pending: pendingSwitch{player: pl, req: req},
		})
	}
	return entries
}

// commitSwitch はプレイヤーのアクティブスロットを確定させ、EntryPhase に渡す EnteredPokemon を返す。
func (ar *ActionResolveOrchestrator) commitSwitch(pending pendingSwitch) phase.EnteredPokemon {
	pending.player.SetActiveSlot(pending.req.IncomingIndex)
	return phase.EnteredPokemon{
		PlayerId: pending.player.Id(),
		Pokemon:  pending.req.Incoming,
	}
}

// validatePendingActions は各プレイヤーが保留行動をちょうど1つ持っていることを検証する。
func (ar *ActionResolveOrchestrator) validatePendingActions(players []*player.Player) error {
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
