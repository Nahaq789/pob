package status

import (
	"math/rand/v2"
	"pob/battle/internal/domain/vo"
)

type Status struct {
	main  *MainStatus
	other []OtherStatus
}

func NewStatus() Status {
	return Status{main: nil, other: nil}
}

func (s *Status) SetMainStatus(m *MainStatus) {
	if s.main != nil {
		return
	}
	s.main = m
}

func (s *Status) ForceSetMainStatus(m *MainStatus) {
	s.main = m
}

// 仮 pokemon集約側に移動する
func RollSleepCount() vo.Count {
	return vo.NewCount(2 + rand.IntN(2))
}
