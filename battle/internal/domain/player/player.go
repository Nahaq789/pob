package player

import (
	"fmt"
	"pob/battle/internal/domain/ground"
	"pob/battle/internal/domain/pokemon"
)

type Player struct {
	id         string
	party      [6]*pokemon.Pokemon // 手持ちポケモン
	selected   [3]*pokemon.Pokemon // 選出された3匹
	activeSlot int                 // 現在場に出ているポケモンのindex
	grounds    []ground.State      // プレイヤーに影響するフィールド効果(ステロとか)
}

func NewPlayer(id string, party [6]*pokemon.Pokemon, grounds []ground.State) *Player {
	return &Player{
		id:         id,
		party:      party,
		selected:   [3]*pokemon.Pokemon{},
		activeSlot: 0,
		grounds:    grounds,
	}
}

func (p *Player) Active() *pokemon.Pokemon {
	return p.selected[p.activeSlot]
}

func (p *Player) Select(indices [3]int) error {
	if p.selected[0] != nil {
		return fmt.Errorf("player has already selected pokemon")
	}

	seen := map[int]bool{}
	for _, i := range indices {
		if i < 0 || i > 5 {
			return fmt.Errorf("invalid party index: %d", i)
		}
		if seen[i] {
			return fmt.Errorf("duplicate party index: %d", i)
		}
		seen[i] = true
	}

	for i, idx := range indices {
		p.selected[i] = p.party[idx]
	}

	return nil
}

func (p *Player) Switch(index int) error {
	// アクティブポケモンを変更する
	err := p.validateSlot(index)
	if err != nil {
		return err
	}

	if p.selected[index].IsFainted() {
		return fmt.Errorf("cannot switch to a fainted pokemon at slot: %d", index)
	}

	p.setActiveSlot(index)
	p.Active().ResetOnSwitchOut()
	return nil
}

func (p *Player) setActiveSlot(index int) {
	p.activeSlot = index
}

func (p *Player) validateSlot(index int) error {
	if p.activeSlot == index {
		return fmt.Errorf("pokemon at slot %d is already active", index)
	}
	if index < 0 || index > 2 {
		return fmt.Errorf("invalid selected slot: %d", index)
	}
	return nil
}

func (p *Player) Id() string {
	return p.id
}
