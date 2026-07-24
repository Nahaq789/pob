package phase

import (
	"errors"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/pokemon"
	"sort"
)

// EntryPhaseHandler はポケモンが場に出た際の一連の処理を担う。
// バトル開始時の同時入場（複数体）と、ターン中の交代による入場（1体、
// または双方同時交代時は2体）の両方をこのHandleで処理する。
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

// Handle は entered を素早さ順（降順）にソートしたうえで、
// 1体ずつ Entered() を呼び出し、発動した特性・道具のハンドラーを実行する。
// 戻り値はプレイヤーIDごとの発動メッセージ一覧。
func (e *EntryPhaseHandler) Handle(entered []EnteredPokemon, b *battle.Battle) (map[string][]string, error) {
	ordered := make([]EnteredPokemon, len(entered))
	copy(ordered, entered)

	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Pokemon.Speed() > ordered[j].Pokemon.Speed()
	})

	resultMessages := make(map[string][]string, 0)
	for _, ep := range ordered {
		// 同一プレイヤーの重複入場は呼び出し元の実装ミス以外では起こり得ないため異常系とする
		if _, exists := resultMessages[ep.PlayerId]; exists {
			return nil, errors.New("同一プレイヤーのポケモンが重複しています")
		}

		ep.Pokemon.Entered()
		messages, err := e.dispatch(ep, b)
		if err != nil {
			return nil, err
		}

		resultMessages[ep.PlayerId] = messages
	}
	return resultMessages, nil
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
		if handler, ok := e.registry.entryAbilityHandler[int(abilityId)]; ok {
			result := handler.Handle(ctx)
			if result.Err != nil {
				return nil, result.Err
			}
			messages = append(messages, result.Message)
		}

		// どうぐを持っていれば、どうぐのハンドラーも処理
		if item != nil {
			if handler, ok := e.registry.entryItemHandler[itemId]; ok {
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
