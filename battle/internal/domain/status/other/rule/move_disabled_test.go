package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/status/other/rule"
)

func TestMoveDisabled_Resolve(t *testing.T) {
	t.Run("残りターンあり: cleared=false", func(t *testing.T) {
		m := rule.NewMoveDisabled(2, 33)
		cleared, addConfusion := m.Resolve(other.OtherStatusContext{})
		if cleared {
			t.Error("expected cleared=false")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})

	t.Run("残りターン0: cleared=true", func(t *testing.T) {
		m := rule.NewMoveDisabled(1, 33)
		cleared, addConfusion := m.Resolve(other.OtherStatusContext{})
		if !cleared {
			t.Error("expected cleared=true")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})
}

func TestMoveDisabled_Kind(t *testing.T) {
	m := rule.NewMoveDisabled(3, 33)
	if m.Kind() != status.OtherCondition("move_disabled") {
		t.Errorf("unexpected kind: %v", m.Kind())
	}
}

func TestMoveDisabled_MoveId(t *testing.T) {
	m := rule.NewMoveDisabled(3, 33)
	if m.MoveId() != 33 {
		t.Errorf("unexpected moveId: %v", m.MoveId())
	}
}

func TestMoveDisabled_Handle(t *testing.T) {
	t.Run("封じられた技を選択: blocked=true", func(t *testing.T) {
		m := rule.NewMoveDisabled(3, 33)
		msg, blocked := m.Handle(33)
		if !blocked {
			t.Error("expected blocked=true")
		}
		if msg == "" {
			t.Error("expected non-empty message")
		}
	})

	t.Run("別の技を選択: blocked=false", func(t *testing.T) {
		m := rule.NewMoveDisabled(3, 33)
		msg, blocked := m.Handle(99)
		if blocked {
			t.Error("expected blocked=false")
		}
		if msg != "" {
			t.Errorf("expected empty message, got: %v", msg)
		}
	})
}
