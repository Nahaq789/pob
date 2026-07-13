package hp

import "pob/battle/internal/domain/vo"

type HP struct {
	current vo.Count
	max     int
}

func NewHP(max int) HP {
	return HP{current: vo.NewCount(max), max: max}
}

func (h HP) Damage(amount int) HP {
	h.current = h.current.Consume(amount)
	return h
}

func (h HP) Heal(amount int) HP {
	next := h.current.Recover(amount)
	if next.Value() > h.max {
		return HP{current: vo.NewCount(h.max), max: h.max}
	}
	h.current = next
	return h
}

func (h HP) IsEmpty() bool {
	return h.current.IsEmpty()
}

func (h HP) Max() int {
	return h.max
}
