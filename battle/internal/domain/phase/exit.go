package phase

import (
	"errors"
	"fmt"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/pokemon"
)

// ExitPhaseHandler はポケモンが場から退いた際の処理を担う。
// 複数体の同時退場は呼び出し元（orchestrate.ActionResolveOrchestrator）が素早さ順にソートして
// 1体ずつ Handle を呼び出す責務を持つ。
type ExitPhaseHandler struct {
	registry *Registry
}

func NewExitPhaseHandler(r *Registry) *ExitPhaseHandler {
	return &ExitPhaseHandler{registry: r}
}

// ExitingPokemon は場から退くポケモンと、そのポケモンを操作するプレイヤーのIDを紐づける。
// Incoming は退場後に入場する予定のポケモンで、退場時の特性・道具ハンドラーに渡される。
type ExitingPokemon struct {
	PlayerId string
	Pokemon  *pokemon.Pokemon
	Incoming *pokemon.Pokemon
}

// Handle は1体のポケモンの退場処理を行い、発動した特性・道具のハンドラーを実行する。
// 退場メッセージ（自分視点・相手視点）および特性・道具効果メッセージを
// プレイヤーIDごとにまとめて返す。
func (e *ExitPhaseHandler) Handle(exiting ExitingPokemon, b *battle.Battle) (map[string][]string, error) {
	self := b.PlayerById(exiting.PlayerId)
	opponent := b.Opponent(self)

	exiting.Pokemon.Exited()
	messages, err := e.dispatch(exiting, b)
	if err != nil {
		return nil, err
	}

	result := map[string][]string{
		self.Id():     {fmt.Sprintf("戻れ！%s！", exiting.Pokemon.Name())},
		opponent.Id(): {fmt.Sprintf("相手は%sを引っ込めた！", exiting.Pokemon.Name())},
	}
	for _, msg := range messages {
		result[self.Id()] = append(result[self.Id()], msg)
		result[opponent.Id()] = append(result[opponent.Id()], msg)
	}
	return result, nil
}

// dispatch は1体のポケモンについて発行されたイベントを取り出し、
// EventExited のイベントに対応する特性・道具ハンドラーを実行する。
// 特性ハンドラーを先に処理し、その後に道具ハンドラーを処理する
// （Gen1範囲では道具が特性より先に発動するケースは存在しないため）。
func (e *ExitPhaseHandler) dispatch(ep ExitingPokemon, b *battle.Battle) ([]string, error) {
	if ep.Pokemon == nil {
		return nil, errors.New("exiting pokemon is nil")
	}

	p := ep.Pokemon
	events := p.PullEvents()
	messages := make([]string, 0)

	for _, event := range events {
		// ポケモンを引っ込めたときのイベントのみループを続ける
		if event.Kind != pokemon.EventExited {
			continue
		}
		abilityId := p.Ability().GetCurrentId()
		item := p.HeldItem()
		var itemId int
		if item != nil {
			itemId = int(item.Id())
		}

		ctx := NewExitContext(ep.PlayerId, int(abilityId), itemId, ep.Incoming, b)

		// 先に特性ハンドラーを処理
		if handler, ok := e.registry.exitAbilityHandlers[int(abilityId)]; ok {
			result := handler.Handle(ctx)
			if result.Err != nil {
				return nil, result.Err
			}
			messages = append(messages, result.Message)
		}

		// どうぐを持っていれば、どうぐのハンドラーも処理
		if item != nil {
			if handler, ok := e.registry.exitItemHandlers[itemId]; ok {
				result := handler.Handle(ctx)
				if result.Err != nil {
					return nil, result.Err
				}

				messages = append(messages, result.Message)
			}
		}

	}

	return messages, nil
}
