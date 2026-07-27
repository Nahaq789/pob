package player

import (
	"fmt"
	"pob/battle/internal/domain/ground"
	"pob/battle/internal/domain/pokemon"
)

type Player struct {
	id            string
	party         [6]*pokemon.Pokemon // 手持ちポケモン
	selected      [3]*pokemon.Pokemon // 選出された3匹
	activeSlot    int                 // 現在場に出ているポケモンのindex
	grounds       []ground.State      // プレイヤーに影響するフィールド効果(ステロとか)
	pendingSwitch *SwitchRequest      // ポケモン交代
}

func NewPlayer(id string, party [6]*pokemon.Pokemon, grounds []ground.State) *Player {
	return &Player{
		id:            id,
		party:         party,
		selected:      [3]*pokemon.Pokemon{},
		activeSlot:    0,
		grounds:       grounds,
		pendingSwitch: nil,
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

// アクティブポケモンを交換する
// 実際、交換自体が複雑なので、ここでは交換意思を表明する
func (p *Player) Switch(index int) error {
	err := p.validateSlot(index)
	if err != nil {
		return err
	}

	incoming := p.selected[index]
	if incoming.IsFainted() {
		return fmt.Errorf("cannot switch to a fainted pokemon at slot: %d", index)
	}

	outgoing := p.Active()
	p.pendingSwitch = &SwitchRequest{Outgoing: outgoing, Incoming: incoming}
	return nil
}

func (p *Player) SetActiveSlot(index int) {
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

func (p *Player) PullPendingSwitch() *SwitchRequest {
	req := p.pendingSwitch
	p.pendingSwitch = nil
	return req
}
