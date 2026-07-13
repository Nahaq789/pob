package battle

import (
	"pob/battle/internal/domain/field"
	"pob/battle/internal/domain/player"
	"pob/battle/internal/domain/room"
	"pob/battle/internal/domain/weather"
)

type Battle struct {
	player1 *player.Player
	player2 *player.Player
	turn    int
	weather *weather.State
	field   *field.State
	room    *room.State
}

func NewBattle(p1, p2 *player.Player) *Battle {
	return &Battle{
		player1: p1,
		player2: p2,
		turn:    1,
		weather: nil,
		field:   nil,
		room:    nil,
	}
}

func (b *Battle) Opponent(p *player.Player) *player.Player {
	if p == b.player1 {
		return b.player2
	}
	return b.player1
}

func (b *Battle) NextTurn() {
	b.turn++
}

func (b *Battle) SetWeather(w *weather.State) {
	b.weather = w
}

func (b *Battle) SetField(f *field.State) {
	b.field = f
}

func (b *Battle) SetRoom(r *room.State) {
	b.room = r
}

func (b *Battle) PlayerById(id string) *player.Player {
	if b.player1.Id() == id {
		return b.player1
	}
	if b.player2.Id() == id {
		return b.player2
	}
	return nil
}
