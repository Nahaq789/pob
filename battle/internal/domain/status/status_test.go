package status_test

import (
	"testing"

	"pob/battle/internal/domain/status"
)

type fakeOtherStatus struct {
	kind status.OtherCondition
}

func (f fakeOtherStatus) Resolve(status.OtherStatusContext) (bool, bool, string) {
	return false, false, ""
}

func (f fakeOtherStatus) Kind() status.OtherCondition { return f.kind }

func newStatusTestOther(kind status.OtherCondition) status.OtherStatus {
	return fakeOtherStatus{kind: kind}
}

func newStatusTestMain(t *testing.T, c status.Condition) *status.MainStatus {
	t.Helper()
	m, err := status.NewMainStatus(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return &m
}

func kindsOf(others []status.OtherStatus) []status.OtherCondition {
	kinds := make([]status.OtherCondition, 0, len(others))
	for _, o := range others {
		kinds = append(kinds, o.Kind())
	}
	return kinds
}

const (
	kindFlinch    = status.OtherCondition(status.Flinch)
	kindConfusion = status.OtherCondition(status.Confusion)
	kindEncore    = status.OtherCondition(status.Encore)
)

func TestNewStatus(t *testing.T) {
	t.Run("生成時: Main は nil", func(t *testing.T) {
		s := status.NewStatus()
		if s.Main() != nil {
			t.Errorf("expected nil main, got %v", s.Main())
		}
	})

	t.Run("生成時: Others は空かつ非nil", func(t *testing.T) {
		s := status.NewStatus()
		if s.Others() == nil {
			t.Fatal("expected non-nil others")
		}
		if len(s.Others()) != 0 {
			t.Errorf("expected 0 others, got %d", len(s.Others()))
		}
	})
}

func TestNewStatusWith(t *testing.T) {
	t.Run("main 指定: Main が同じポインタを返す", func(t *testing.T) {
		m := newStatusTestMain(t, status.Burn)
		s := status.NewStatusWith(m, nil)
		if s.Main() != m {
			t.Errorf("expected %p, got %p", m, s.Main())
		}
	})

	t.Run("others=nil: 空かつ非nilのスライス", func(t *testing.T) {
		s := status.NewStatusWith(nil, nil)
		if s.Others() == nil {
			t.Fatal("expected non-nil others")
		}
		if len(s.Others()) != 0 {
			t.Errorf("expected 0 others, got %d", len(s.Others()))
		}
	})

	t.Run("others 指定: 順序を保持する", func(t *testing.T) {
		s := status.NewStatusWith(nil, []status.OtherStatus{
			newStatusTestOther(kindFlinch),
			newStatusTestOther(kindConfusion),
		})
		got := kindsOf(s.Others())
		if len(got) != 2 || got[0] != kindFlinch || got[1] != kindConfusion {
			t.Errorf("expected [flinch confusion], got %v", got)
		}
	})
}

func TestStatus_OtherMap(t *testing.T) {
	t.Run("空: 空マップ", func(t *testing.T) {
		s := status.NewStatus()
		if len(s.OtherMap()) != 0 {
			t.Errorf("expected empty map, got %v", s.OtherMap())
		}
	})

	t.Run("複数登録: Kind をキーに引ける", func(t *testing.T) {
		s := status.NewStatusWith(nil, []status.OtherStatus{
			newStatusTestOther(kindFlinch),
			newStatusTestOther(kindEncore),
		})
		m := s.OtherMap()
		if len(m) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(m))
		}
		if _, ok := m[kindFlinch]; !ok {
			t.Error("expected flinch entry")
		}
		if _, ok := m[kindEncore]; !ok {
			t.Error("expected encore entry")
		}
	})

	t.Run("未登録の Kind: 存在しない", func(t *testing.T) {
		s := status.NewStatusWith(nil, []status.OtherStatus{newStatusTestOther(kindFlinch)})
		if _, ok := s.OtherMap()[kindConfusion]; ok {
			t.Error("expected no confusion entry")
		}
	})
}

func TestStatus_SetMainStatus(t *testing.T) {
	t.Run("未設定: 設定される", func(t *testing.T) {
		s := status.NewStatus()
		m := newStatusTestMain(t, status.Burn)
		s.SetMainStatus(m)
		if s.Main() != m {
			t.Errorf("expected %v, got %v", m, s.Main())
		}
	})

	t.Run("設定済み: 上書きされない", func(t *testing.T) {
		first := newStatusTestMain(t, status.Burn)
		s := status.NewStatusWith(first, nil)
		s.SetMainStatus(newStatusTestMain(t, status.Paralysis))
		if s.Main() != first {
			t.Errorf("expected Burn to remain, got %v", s.Main().Condition())
		}
	})
}

func TestStatus_ForceSetMainStatus(t *testing.T) {
	t.Run("設定済み: 上書きされる", func(t *testing.T) {
		s := status.NewStatusWith(newStatusTestMain(t, status.Burn), nil)
		m := newStatusTestMain(t, status.Paralysis)
		s.ForceSetMainStatus(m)
		if s.Main() != m {
			t.Errorf("expected Paralysis, got %v", s.Main().Condition())
		}
	})

	t.Run("nil 指定: 解除される", func(t *testing.T) {
		s := status.NewStatusWith(newStatusTestMain(t, status.Burn), nil)
		s.ForceSetMainStatus(nil)
		if s.Main() != nil {
			t.Errorf("expected nil, got %v", s.Main())
		}
	})
}

func TestStatus_AddOtherStatus(t *testing.T) {
	t.Run("追加: 末尾に追加される", func(t *testing.T) {
		s := status.NewStatus()
		s.AddOtherStatus(newStatusTestOther(kindFlinch))
		s.AddOtherStatus(newStatusTestOther(kindConfusion))
		got := kindsOf(s.Others())
		if len(got) != 2 || got[0] != kindFlinch || got[1] != kindConfusion {
			t.Errorf("expected [flinch confusion], got %v", got)
		}
	})
}

func TestStatus_RemoveOtherStatus(t *testing.T) {
	tests := []struct {
		name   string
		others []status.OtherCondition
		remove status.OtherCondition
		want   []status.OtherCondition
	}{
		{"先頭を削除: 残りの順序を保持", []status.OtherCondition{kindFlinch, kindConfusion, kindEncore}, kindFlinch, []status.OtherCondition{kindConfusion, kindEncore}},
		{"中間を削除: 残りの順序を保持", []status.OtherCondition{kindFlinch, kindConfusion, kindEncore}, kindConfusion, []status.OtherCondition{kindFlinch, kindEncore}},
		{"末尾を削除: 残りの順序を保持", []status.OtherCondition{kindFlinch, kindConfusion, kindEncore}, kindEncore, []status.OtherCondition{kindFlinch, kindConfusion}},
		{"未登録の Kind: 変化なし", []status.OtherCondition{kindFlinch}, kindConfusion, []status.OtherCondition{kindFlinch}},
		{"同じ Kind が複数: すべて削除", []status.OtherCondition{kindFlinch, kindEncore, kindFlinch}, kindFlinch, []status.OtherCondition{kindEncore}},
		{"空: 空のまま", nil, kindFlinch, []status.OtherCondition{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := status.NewStatus()
			for _, k := range tt.others {
				s.AddOtherStatus(newStatusTestOther(k))
			}
			s.RemoveOtherStatus(tt.remove)
			got := kindsOf(s.Others())
			if len(got) != len(tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("expected %v, got %v", tt.want, got)
					break
				}
			}
		})
	}
}
