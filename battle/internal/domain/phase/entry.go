package phase

import (
	"errors"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/pokemon"
	"sort"
)

// ポケモンが場に出たときのフェーズ
type EntryPhaseHandler struct {
	registry *Registry
}

func NewEntryPhaseHandler(r *Registry) *EntryPhaseHandler {
	return &EntryPhaseHandler{
		registry: r,
	}
}

type EnteredPokemon struct {
	PlayerId string
	Pokemon  *pokemon.Pokemon
}

func (e *EntryPhaseHandler) Handle(entered []EnteredPokemon, b *battle.Battle) (map[string][]string, error) {
	ordered := make([]EnteredPokemon, len(entered))
	copy(ordered, entered)

	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Pokemon.Speed() > ordered[j].Pokemon.Speed()
	})

	resultMessages := make(map[string][]string, 0)
	for _, ep := range ordered {
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
