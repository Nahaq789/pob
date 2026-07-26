package phase

import (
	"errors"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/pokemon"
	"sort"
)

// ExitPhaseHandler はポケモンが場から退いた際の一連の処理を担う。
// 交代による退場（1体、または双方同時交代時は2体）をこのHandleで処理する。
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

// Handle は exiting を素早さ順（降順）にソートしたうえで、
// 1体ずつ Exited() を呼び出し、発動した特性・道具のハンドラーを実行する。
// 戻り値はプレイヤーIDごとの発動メッセージ一覧。
func (e *ExitPhaseHandler) Handle(exiting []ExitingPokemon, b *battle.Battle) (map[string][]string, error) {
	ordered := make([]ExitingPokemon, len(exiting))
	copy(ordered, exiting)

	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Pokemon.Speed() > ordered[j].Pokemon.Speed()
	})

	resultMessages := make(map[string][]string, 0)
	for _, ep := range ordered {
		// 同一プレイヤーの重複退場は呼び出し元の実装ミス以外では起こり得ないため異常系とする
		if _, exists := resultMessages[ep.PlayerId]; exists {
			return nil, errors.New("同一プレイヤーのポケモンが重複しています")
		}

		ep.Pokemon.Exited()
		messages, err := e.dispatch(ep, b)
		if err != nil {
			return nil, err
		}

		resultMessages[ep.PlayerId] = messages
	}
	return resultMessages, nil
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
