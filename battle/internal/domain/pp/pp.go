package pp

import (
	"fmt"
	"pob/battle/internal/domain/vo"
)

type PP struct {
	current vo.Count
	max     int
}

func NewPP(c vo.Count, m int) PP {
	return PP{current: c, max: m}
}

func (p PP) ConsumeOne() PP {
	c := p.current.Decrement()
	return PP{current: c, max: p.max}
}

func (p PP) Consume(n int) PP {
	c := p.current.Consume(n)
	return PP{current: c, max: p.max}
}

func (p PP) Recover(n int) (PP, error) {
	if p.current.Value()+n > p.max {
		return PP{}, fmt.Errorf("recover amount %d exceeds max pp %d", n, p.max)
	}
	c := p.current.Increment(n)
	return PP{current: c, max: p.max}, nil
}

func (p PP) IsEmpty() bool {
	return p.current.Value() == 0
}
