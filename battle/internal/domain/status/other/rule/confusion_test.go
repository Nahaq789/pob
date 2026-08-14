package rule_test

import (
	"testing"

	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other/rule"
)

func TestConfusion_Resolve(t *testing.T) {
	t.Run("残りターンあり: cleared=false", func(t *testing.T) {
		c := rule.NewConfusion(2)
		cleared, addConfusion, _ := c.Resolve(status.OtherStatusContext{})
		if cleared {
			t.Error("expected cleared=false")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
	})

	t.Run("残りターン0: cleared=true かつ解除メッセージあり", func(t *testing.T) {
		c := rule.NewConfusion(1)
		cleared, addConfusion, message := c.Resolve(status.OtherStatusContext{ActorName: "ピカチュウ"})
		if !cleared {
			t.Error("expected cleared=true")
		}
		if addConfusion {
			t.Error("expected addConfusion=false")
		}
		if message == "" {
			t.Error("expected non-empty message on cleared")
		}
	})

	t.Run("残りターンあり: メッセージなし", func(t *testing.T) {
		c := rule.NewConfusion(2)
		_, _, message := c.Resolve(status.OtherStatusContext{ActorName: "ピカチュウ"})
		if message != "" {
			t.Errorf("expected empty message, got: %v", message)
		}
	})

	t.Run("複数ターン経過後に解除される", func(t *testing.T) {
		c := rule.NewConfusion(3)
		for i := range 2 {
			cleared, _, _ := c.Resolve(status.OtherStatusContext{ActorName: "ピカチュウ"})
			if cleared {
				t.Errorf("turn %d: expected cleared=false", i+1)
			}
		}
		cleared, _, _ := c.Resolve(status.OtherStatusContext{ActorName: "ピカチュウ"})
		if !cleared {
			t.Error("turn 3: expected cleared=true")
		}
	})
}

func TestConfusion_Kind(t *testing.T) {
	c := rule.NewConfusion(3)
	if c.Kind() != status.OtherCondition("confusion") {
		t.Errorf("unexpected kind: %v", c.Kind())
	}
}

func TestConfusion_CheckSelfHit(t *testing.T) {
	t.Run("true と false の両方が返る", func(t *testing.T) {
		c := rule.NewConfusion(1)
		gotTrue, gotFalse := false, false
		for range 300 {
			_, hit := c.CheckSelfHit("テスト")
			if hit {
				gotTrue = true
			} else {
				gotFalse = true
			}
			if gotTrue && gotFalse {
				break
			}
		}
		if !gotTrue {
			t.Error("CheckSelfHit never returned true in 300 trials")
		}
		if !gotFalse {
			t.Error("CheckSelfHit never returned false in 300 trials")
		}
	})

	t.Run("自傷率がおよそ 1/3 である", func(t *testing.T) {
		c := rule.NewConfusion(1)
		const trials = 3000
		hits := 0
		for range trials {
			if _, hit := c.CheckSelfHit("テスト"); hit {
				hits++
			}
		}
		rate := float64(hits) / trials
		if rate < 0.20 || rate > 0.47 {
			t.Errorf("unexpected self-hit rate: %.2f (expected ~0.33)", rate)
		}
	})
}
