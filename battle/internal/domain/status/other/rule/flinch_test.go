package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other/rule"
)

func TestFlinch_Resolve(t *testing.T) {
	t.Run("常にcleared=true", func(t *testing.T) {
		f := rule.NewFlinch()
		cleared, addConfusion, _ := f.Resolve(status.OtherStatusContext{})
		if !cleared {
			t.Error("expected cleared=true")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})

	t.Run("メッセージなし", func(t *testing.T) {
		f := rule.NewFlinch()
		_, _, message := f.Resolve(status.OtherStatusContext{})
		if message != "" {
			t.Errorf("expected empty message, got: %v", message)
		}
	})
}

func TestFlinch_Kind(t *testing.T) {
	f := rule.NewFlinch()
	if f.Kind() != status.OtherCondition("flinch") {
		t.Errorf("unexpected kind: %v", f.Kind())
	}
}

func TestFlinch_Handle(t *testing.T) {
	f := rule.NewFlinch()
	message := f.Handle("ピカチュウ")
	if message == "" {
		t.Error("expected non-empty message")
	}
}
