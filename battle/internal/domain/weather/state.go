package weather

import "pob/battle/internal/domain/vo"

type condition int

const (
	Sunny condition = iota
	Rain
	Sandstorm
	Hail
	Fog
	Snow
)

type State struct {
	c     condition
	count vo.Count
}

func NewState(c condition, count vo.Count) State {
	return State{c: c, count: count}
}

func (w State) C() condition {
	return w.c
}

func (w State) Count() vo.Count {
	return w.count
}

func (w State) Decrement() State {
	return State{
		c:     w.c,
		count: w.count.Decrement(),
	}
}
