package status_test

import (
	"strings"
	"testing"

	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/vo"
)

func newMainTestStatus(t *testing.T, c status.Condition) status.MainStatus {
	t.Helper()
	m, err := status.NewMainStatus(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return m
}

func newMainTestSleep(t *testing.T, count int) status.MainStatus {
	t.Helper()
	m, err := status.NewSleep(vo.NewCount(count))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return m
}

func TestNewMainStatus(t *testing.T) {
	tests := []struct {
		name    string
		cond    status.Condition
		wantErr bool
	}{
		{"どく: 生成成功", status.Poison, false},
		{"まひ: 生成成功", status.Paralysis, false},
		{"やけど: 生成成功", status.Burn, false},
		{"こおり: 生成成功", status.Freeze, false},
		{"なし: 生成成功", status.None, false},
		{"ねむり: 専用コンストラクタが必要なためエラー", status.Sleep, true},
		{"もうどく: 専用コンストラクタが必要なためエラー", status.BadPoison, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := status.NewMainStatus(tt.cond)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if m.Condition() != tt.cond {
				t.Errorf("expected condition %v, got %v", tt.cond, m.Condition())
			}
			if m.Count().Value() != 0 {
				t.Errorf("expected count 0, got %d", m.Count().Value())
			}
		})
	}
}

func TestNewSleep(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{"count=2: 生成成功", 2, false},
		{"count=3: 生成成功", 3, false},
		{"count=0: エラー", 0, true},
		{"count=1: エラー", 1, true},
		{"count=4: エラー", 4, true},
		{"count=-1: エラー", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := status.NewSleep(vo.NewCount(tt.count))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if m.Condition() != status.Sleep {
				t.Errorf("expected Sleep, got %v", m.Condition())
			}
			if m.Count().Value() != tt.count {
				t.Errorf("expected count %d, got %d", tt.count, m.Count().Value())
			}
		})
	}
}

func TestNewBadPoison(t *testing.T) {
	t.Run("生成時: condition=BadPoison", func(t *testing.T) {
		m := status.NewBadPoison()
		if m.Condition() != status.BadPoison {
			t.Errorf("expected BadPoison, got %v", m.Condition())
		}
	})

	t.Run("生成時: count=1", func(t *testing.T) {
		m := status.NewBadPoison()
		if m.Count().Value() != 1 {
			t.Errorf("expected count 1, got %d", m.Count().Value())
		}
	})
}

func TestMainStatus_DecrementCount(t *testing.T) {
	t.Run("ねむり count=3: count=2", func(t *testing.T) {
		got := newMainTestSleep(t, 3).DecrementCount()
		if got.Count().Value() != 2 {
			t.Errorf("expected count 2, got %d", got.Count().Value())
		}
	})

	t.Run("ねむり count=3: condition は Sleep のまま", func(t *testing.T) {
		got := newMainTestSleep(t, 3).DecrementCount()
		if got.Condition() != status.Sleep {
			t.Errorf("expected Sleep, got %v", got.Condition())
		}
	})

	t.Run("2回減算 count=2: count=0", func(t *testing.T) {
		got := newMainTestSleep(t, 2).DecrementCount().DecrementCount()
		if !got.Count().IsEmpty() {
			t.Errorf("expected empty count, got %d", got.Count().Value())
		}
	})

	t.Run("count=0: 負にならず0のまま", func(t *testing.T) {
		got := newMainTestStatus(t, status.Paralysis).DecrementCount()
		if got.Count().Value() != 0 {
			t.Errorf("expected count 0, got %d", got.Count().Value())
		}
	})

	t.Run("元の値: 変更されない", func(t *testing.T) {
		m := newMainTestSleep(t, 3)
		_ = m.DecrementCount()
		if m.Count().Value() != 3 {
			t.Errorf("expected original count 3, got %d", m.Count().Value())
		}
	})
}

func TestMainStatus_IncrementCount(t *testing.T) {
	t.Run("もうどく count=1: count=2", func(t *testing.T) {
		got := status.NewBadPoison().IncrementCount()
		if got.Count().Value() != 2 {
			t.Errorf("expected count 2, got %d", got.Count().Value())
		}
	})

	t.Run("もうどく: condition は BadPoison のまま", func(t *testing.T) {
		got := status.NewBadPoison().IncrementCount()
		if got.Condition() != status.BadPoison {
			t.Errorf("expected BadPoison, got %v", got.Condition())
		}
	})

	t.Run("3回加算: count=4", func(t *testing.T) {
		got := status.NewBadPoison().IncrementCount().IncrementCount().IncrementCount()
		if got.Count().Value() != 4 {
			t.Errorf("expected count 4, got %d", got.Count().Value())
		}
	})

	t.Run("元の値: 変更されない", func(t *testing.T) {
		m := status.NewBadPoison()
		_ = m.IncrementCount()
		if m.Count().Value() != 1 {
			t.Errorf("expected original count 1, got %d", m.Count().Value())
		}
	})
}

func TestMainStatus_IsNone(t *testing.T) {
	t.Run("nil レシーバー: true", func(t *testing.T) {
		var m *status.MainStatus
		if !m.IsNone() {
			t.Error("expected true")
		}
	})

	t.Run("condition=None: true", func(t *testing.T) {
		m := newMainTestStatus(t, status.None)
		if !m.IsNone() {
			t.Error("expected true")
		}
	})

	t.Run("condition=Burn: false", func(t *testing.T) {
		m := newMainTestStatus(t, status.Burn)
		if m.IsNone() {
			t.Error("expected false")
		}
	})

	t.Run("ねむり: false", func(t *testing.T) {
		m := newMainTestSleep(t, 2)
		if m.IsNone() {
			t.Error("expected false")
		}
	})
}

func TestMainStatus_IsSleep(t *testing.T) {
	t.Run("ねむり count=2: asleep=true", func(t *testing.T) {
		_, ok := newMainTestSleep(t, 2).IsSleep("ピカチュウ")
		if !ok {
			t.Error("expected asleep=true")
		}
	})

	t.Run("ねむり count=2: 眠っているメッセージ", func(t *testing.T) {
		msg, _ := newMainTestSleep(t, 2).IsSleep("ピカチュウ")
		if !strings.Contains(msg, "ピカチュウ") || !strings.Contains(msg, "眠っている") {
			t.Errorf("expected sleeping message, got: %v", msg)
		}
	})

	t.Run("ねむり count=0: asleep=false", func(t *testing.T) {
		_, ok := newMainTestSleep(t, 2).DecrementCount().DecrementCount().IsSleep("ピカチュウ")
		if ok {
			t.Error("expected asleep=false")
		}
	})

	t.Run("ねむり count=0: 目覚めたメッセージ", func(t *testing.T) {
		msg, _ := newMainTestSleep(t, 2).DecrementCount().DecrementCount().IsSleep("ピカチュウ")
		if !strings.Contains(msg, "ピカチュウ") || !strings.Contains(msg, "目覚めた") {
			t.Errorf("expected wake-up message, got: %v", msg)
		}
	})

	t.Run("ねむり以外: asleep=false", func(t *testing.T) {
		_, ok := newMainTestStatus(t, status.Paralysis).IsSleep("ピカチュウ")
		if ok {
			t.Error("expected asleep=false")
		}
	})
}

func TestMainStatus_IsFreeze(t *testing.T) {
	// NewMainStatus(Freeze) は count=0 で生成されるため IsFreeze が即時 false を返し、解凍扱いとなる。
	t.Run("NewMainStatus(Freeze): frozen=false", func(t *testing.T) {
		_, ok := newMainTestStatus(t, status.Freeze).IsFreeze("ピカチュウ")
		if ok {
			t.Error("expected frozen=false")
		}
	})

	t.Run("NewMainStatus(Freeze): こおりが溶けたメッセージ", func(t *testing.T) {
		msg, _ := newMainTestStatus(t, status.Freeze).IsFreeze("ピカチュウ")
		if !strings.Contains(msg, "ピカチュウ") || !strings.Contains(msg, "こおりが溶けた") {
			t.Errorf("expected thaw message, got: %v", msg)
		}
	})

	t.Run("こおり以外 count>0: frozen=false", func(t *testing.T) {
		_, ok := newMainTestSleep(t, 2).IsFreeze("ピカチュウ")
		if ok {
			t.Error("expected frozen=false")
		}
	})

	t.Run("こおり count>0: frozen=true", func(t *testing.T) {
		m := newMainTestStatus(t, status.Freeze).IncrementCount()
		_, ok := m.IsFreeze("ピカチュウ")
		if !ok {
			t.Error("expected frozen=true")
		}
	})

	t.Run("こおり count>0: こおって動けないメッセージ", func(t *testing.T) {
		m := newMainTestStatus(t, status.Freeze).IncrementCount()
		msg, _ := m.IsFreeze("ピカチュウ")
		if !strings.Contains(msg, "こおって動くことができない") {
			t.Errorf("expected frozen message, got: %v", msg)
		}
	})
}

func TestMainStatus_IsParalysis(t *testing.T) {
	tests := []struct {
		name string
		cond status.Condition
		want bool
	}{
		{"まひ: true", status.Paralysis, true},
		{"どく: false", status.Poison, false},
		{"やけど: false", status.Burn, false},
		{"こおり: false", status.Freeze, false},
		{"なし: false", status.None, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMainTestStatus(t, tt.cond)
			if got := m.IsParalysis(); got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestMainStatus_CheckParalysis(t *testing.T) {
	t.Run("true と false（行動不能/行動可）の両方が返る", func(t *testing.T) {
		m := newMainTestStatus(t, status.Paralysis)
		gotTrue, gotFalse := false, false
		for range 300 {
			_, p := m.CheckParalysis("ピカチュウ")
			if p {
				gotTrue = true
			} else {
				gotFalse = true
			}
			if gotTrue && gotFalse {
				break
			}
		}
		if !gotTrue {
			t.Error("300試行で行動不能（true）が一度も発生しなかった")
		}
		if !gotFalse {
			t.Error("300試行で行動成功（false）が一度も発生しなかった")
		}
	})

	t.Run("行動不能時: しびれメッセージあり", func(t *testing.T) {
		m := newMainTestStatus(t, status.Paralysis)
		for range 300 {
			msg, p := m.CheckParalysis("ピカチュウ")
			if !p {
				continue
			}
			if !strings.Contains(msg, "ピカチュウ") || !strings.Contains(msg, "しびれて") {
				t.Errorf("expected paralysis message, got: %v", msg)
			}
			return
		}
		t.Error("300試行で行動不能が一度も発生しなかった")
	})

	t.Run("行動可能時: メッセージなし", func(t *testing.T) {
		m := newMainTestStatus(t, status.Paralysis)
		for range 300 {
			msg, p := m.CheckParalysis("ピカチュウ")
			if p {
				continue
			}
			if msg != "" {
				t.Errorf("expected empty message, got: %v", msg)
			}
			return
		}
		t.Error("300試行で行動成功が一度も発生しなかった")
	})

	t.Run("行動不能率がおよそ 1/8 である", func(t *testing.T) {
		m := newMainTestStatus(t, status.Paralysis)
		const trials = 3000
		hits := 0
		for range trials {
			if _, p := m.CheckParalysis("ピカチュウ"); p {
				hits++
			}
		}
		rate := float64(hits) / trials
		if rate < 0.08 || rate > 0.17 {
			t.Errorf("unexpected paralysis rate: %.3f (expected ~0.125)", rate)
		}
	})
}
