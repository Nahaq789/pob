package phase

import (
	"errors"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/pokemon"
)

// EntryPhaseHandler はポケモンが場に出た際の処理を担う。
// 複数体の同時入場は呼び出し元（ActionResolvePhaseHandler）が素早さ順にソートして
// 1体ずつ Handle を呼び出す責務を持つ。
type EntryPhaseHandler struct {
	registry *Registry
}

func NewEntryPhaseHandler(r *Registry) *EntryPhaseHandler {
	return &EntryPhaseHandler{
		registry: r,
	}
}

// EnteredPokemon は場に出たポケモンと、そのポケモンを操作するプレイヤーのIDを紐づける。
// Pokemon集約自身は自分がどのプレイヤーに属するか知らないため、
// phase層でこの対応関係を保持する。
type EnteredPokemon struct {
	PlayerId string
	Pokemon  *pokemon.Pokemon
}

// Handle は1体のポケモンの入場処理を行い、発動した特性・道具のハンドラーを実行する。
// 戻り値は発動メッセージ一覧。
func (e *EntryPhaseHandler) Handle(entered EnteredPokemon, b *battle.Battle) ([]string, error) {
	entered.Pokemon.Entered()
	messages, err := e.dispatch(entered, b)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// dispatch は1体のポケモンについて発行されたイベントを取り出し、
// EventEntered のイベントに対応する特性・道具ハンドラーを実行する。
// 特性ハンドラーを先に処理し、その後に道具ハンドラーを処理する
// （Gen1範囲では道具が特性より先に発動するケースは存在しないため）。
func (e *EntryPhaseHandler) dispatch(ep EnteredPokemon, b *battle.Battle) ([]string, error) {
	if ep.Pokemon == nil {
		return nil, errors.New("entered pokemon is nil")
	}

	p := ep.Pokemon
	events := p.PullEvents()
	messages := make([]string, 0)
	for _, event := range events {
		// ポケモンを出したときのイベントのみループを続ける
		if event.Kind != pokemon.EventEntered {
			continue
		}
		abilityId := p.Ability().GetCurrentId()
		item := p.HeldItem()
		var itemId int
		if item != nil {
			itemId = int(item.Id())
		}

		ctx := NewEntryContext(ep.PlayerId, int(abilityId), itemId, b)

		// 先に特性ハンドラーを処理
		if handler, ok := e.registry.entryAbilityHandlers[int(abilityId)]; ok {
			result := handler.Handle(ctx)
			if result.Err != nil {
				return nil, result.Err
			}
			messages = append(messages, result.Message)
		}

		// どうぐを持っていれば、どうぐのハンドラーも処理
		if item != nil {
			if handler, ok := e.registry.entryItemHandlers[itemId]; ok {
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
