package priority

import (
	"fmt"
	"pob/battle/internal/domain/vo"
)

type Priority struct {
	value vo.Count
}

func NewPriority(p int) (Priority, error) {
	if p >= -7 && p <= 5 {
		return Priority{value: vo.NewCount(p)}, nil
	}

	return Priority{}, fmt.Errorf("")
}

func (p Priority) Value() vo.Count {
	return p.value
}

func (p Priority) Add(n int) Priority {
	v := p.value.Value() + n
	return Priority{value: vo.NewCount(v)}
}

func (p Priority) Higher(op Priority) bool {
	return p.value.Value() > op.value.Value()
}

func (p Priority) Same(op Priority) bool {
	return p.value.Value() == op.value.Value()
}
