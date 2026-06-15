package status

import (
	"fmt"
	"pob/battle/internal/domain/vo"
)

type MainStatus struct {
	condition Condition
	count     vo.Count
}

func NewMainStatus(c Condition) (MainStatus, error) {
	if c == Sleep || c == BadPoison {
		return MainStatus{}, fmt.Errorf("%s requires a dedicated constructor", c)
	}
	return MainStatus{
		condition: c,
		count:     vo.NewCount(0),
	}, nil
}

func NewSleep(c vo.Count) (MainStatus, error) {
	v := c.Value()
	if v != 2 && v != 3 {
		return MainStatus{}, fmt.Errorf("sleep count must be 2 or 3, got %d", v)
	}
	return MainStatus{
		condition: Sleep,
		count:     c,
	}, nil
}

func NewBadPoison() MainStatus {
	return MainStatus{
		condition: BadPoison,
		count:     vo.NewCount(1),
	}
}

// ねむり：技選択前フェーズで呼ぶ。count-- した新しい MainStatus を返す。
func (m MainStatus) DecrementCount() MainStatus {
	if m.count.IsEmpty() {
		return m
	}
	return MainStatus{
		condition: m.condition,
		count:     m.count.Decrement(),
	}
}

// もうどく：ターン終了フェーズで呼ぶ。count++ した新しい MainStatus を返す。
func (m MainStatus) IncrementCount() MainStatus {
	return MainStatus{
		condition: m.condition,
		count:     m.count.Increment(),
	}
}

func (m MainStatus) Condition() Condition {
	return m.condition
}

func (m MainStatus) Count() vo.Count {
	return m.count
}
