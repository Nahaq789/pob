package ground

import "pob/battle/internal/domain/vo"

type condition int

const (
	Gravity condition = iota

	StealthRock
)

var GlobalConditions = map[condition]GroundHandler{
	// Gravity: "じゅうりょく",
}

var PerSideConditions = map[condition]GroundHandler{
	// StealthRock: "ステルスロック",
}

type Ground struct {
	c     condition
	count vo.Count
}

func NewGround(c condition, count vo.Count) Ground {
	return Ground{c: c, count: count}
}

func (g Ground) C() condition {
	return g.c
}

func (g Ground) Count() vo.Count {
	return g.count
}

func (g Ground) Decrement() Ground {
	return Ground{
		c:     g.c,
		count: g.count.Decrement(),
	}
}

// 仮
type GroundHandler interface {
	Execute() error
}
