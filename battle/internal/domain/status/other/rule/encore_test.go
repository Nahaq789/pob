package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other/rule"
)

func TestEncore_Resolve(t *testing.T) {
	t.Run("残りターンあり: cleared=false", func(t *testing.T) {
		e := rule.NewEncore(33)
		cleared, addConfusion, _ := e.Resolve(status.OtherStatusContext{})
		if cleared {
			t.Error("expected cleared=false")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})

	t.Run("残りターン0: cleared=true", func(t *testing.T) {
		e := rule.NewEncore(33)
		for range 2 {
			e.Resolve(status.OtherStatusContext{})
		}
		cleared, addConfusion, _ := e.Resolve(status.OtherStatusContext{})
		if !cleared {
			t.Error("expected cleared=true")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})

	t.Run("残りターン0: 解除メッセージあり", func(t *testing.T) {
		e := rule.NewEncore(33)
		for range 2 {
			e.Resolve(status.OtherStatusContext{})
		}
		_, _, message := e.Resolve(status.OtherStatusContext{ActorName: "ピカチュウ"})
		if message == "" {
			t.Error("expected non-empty message on cleared")
		}
	})

	t.Run("残りターンあり: メッセージなし", func(t *testing.T) {
		e := rule.NewEncore(33)
		_, _, message := e.Resolve(status.OtherStatusContext{ActorName: "ピカチュウ"})
		if message != "" {
			t.Errorf("expected empty message, got: %v", message)
		}
	})
}

func TestEncore_Kind(t *testing.T) {
	e := rule.NewEncore(33)
	if e.Kind() != status.OtherCondition("encore") {
		t.Errorf("unexpected kind: %v", e.Kind())
	}
}

func TestEncore_MoveId(t *testing.T) {
	e := rule.NewEncore(33)
	if e.MoveId() != 33 {
		t.Errorf("unexpected moveId: %v", e.MoveId())
	}
}

func TestEncore_Handle(t *testing.T) {
	t.Run("アンコールされた技を選択: blocked=false", func(t *testing.T) {
		e := rule.NewEncore(33)
		msg, blocked := e.Handle(33)
		if blocked {
			t.Error("expected blocked=false")
		}
		if msg != "" {
			t.Errorf("expected empty message, got: %v", msg)
		}
	})

	t.Run("別の技を選択: blocked=true", func(t *testing.T) {
		e := rule.NewEncore(33)
		msg, blocked := e.Handle(99)
		if !blocked {
			t.Error("expected blocked=true")
		}
		if msg == "" {
			t.Error("expected non-empty message")
		}
	})
}
