package ground

import "pob/battle/internal/domain/vo"

type condition int

const (
	Gravity condition = iota

	StealthRock
)

var GlobalConditions = map[condition]StateHandler{
	// Gravity: "じゅうりょく",
}

var PerSideConditions = map[condition]StateHandler{
	// StealthRock: "ステルスロック",
}

type State struct {
	c     condition
	count vo.Count
}

func NewState(c condition, count vo.Count) State {
	return State{c: c, count: count}
}

func (g State) C() condition {
	return g.c
}

func (g State) Count() vo.Count {
	return g.count
}

func (g State) Decrement() State {
	return State{
		c:     g.c,
		count: g.count.Decrement(),
	}
}

// 仮
type StateHandler interface {
	Execute() error
}
