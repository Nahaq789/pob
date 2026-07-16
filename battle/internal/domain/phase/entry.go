package phase

import (
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/pokemon"
	"sort"
)

type EntryPhaseHandler struct {
	registry *Registry
}

func NewEntryPhaseHandler(r *Registry) *EntryPhaseHandler {
	return &EntryPhaseHandler{
		registry: r,
	}
}

func (e *EntryPhaseHandler) Handle(entered []*pokemon.Pokemon, b *battle.Battle) error {
	ordered := make([]*pokemon.Pokemon, len(entered))
	copy(ordered, entered)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Speed() > ordered[j].Speed()
	})

	for _, p := range ordered {
		p.Entered()
		e.dispatch(p, b)
	}
	return nil
}

func (e *EntryPhaseHandler) dispatch(p *pokemon.Pokemon, b *battle.Battle) error {
	events := p.PullEvent()
	for _, event := range events {
		if event.Kind != pokemon.EventEntered {
			continue
		}
		if handler, ok := e.registry.entryAbilityHandler[int(p.Ability().GetCurrentId())]; ok {
			abilityId := p.Ability().GetCurrentId()
			itemId := p.HeldItem().Id()
			ctx := NewEntryContext(
				int(p.Id()),
				int(abilityId),
				int(itemId),
				b,
			)
			handler.Handle(nil)
		}
	}

	return nil
}
