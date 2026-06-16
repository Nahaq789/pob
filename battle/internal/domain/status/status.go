package status

import (
	"fmt"
	"math/rand/v2"
	"pob/battle/internal/domain/vo"
)

var overridableStatusMoves = map[int]string{
	173: "ねむる",
}

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

func (s *Status) OverrideMainStatus(m *MainStatus, moveId int) error {
	_, ok := overridableStatusMoves[moveId]
	if !ok {
		return fmt.Errorf("move %d is not overridable", moveId)
	}
	s.main = m
	return nil
}

func RollSleepCount() vo.Count {
	return vo.NewCount(2 + rand.IntN(2))
}
