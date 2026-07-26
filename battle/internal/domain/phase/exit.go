package phase

import (
	"errors"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/pokemon"
	"sort"
)

type ExitPhaseHandler struct {
	registry *Registry
}

func NewExitPhaseHandler(r *Registry) *ExitPhaseHandler {
	return &ExitPhaseHandler{registry: r}
}

type ExitingPokemon struct {
	PlayerId string
	Pokemon  *pokemon.Pokemon
	Incoming *pokemon.Pokemon
}

func (e *ExitPhaseHandler) Handle(exiting []ExitingPokemon, b *battle.Battle) (map[string][]string, error) {
	ordered := make([]ExitingPokemon, len(exiting))
	copy(ordered, exiting)

	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Pokemon.Speed() > ordered[j].Pokemon.Speed()
	})

	resultMessages := make(map[string][]string, 0)
	for _, ep := range ordered {
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

func (e *ExitPhaseHandler) dispatch(ep ExitingPokemon, b *battle.Battle) ([]string, error) {
	if ep.Pokemon == nil {
		return nil, errors.New("exiting pokemon is nil")
	}

	p := ep.Pokemon
	events := p.PullEvents()
	messages := make([]string, 0)

	for _, event := range events {
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

		if handler, ok := e.registry.exitAbilityHandlers[int(abilityId)]; ok {
			result := handler.Handle(ctx)
			if result.Err != nil {
				return nil, result.Err
			}
			messages = append(messages, result.Message)
		}

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
