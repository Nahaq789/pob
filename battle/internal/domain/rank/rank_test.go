package rank

import "testing"

func TestNewRank(t *testing.T) {
	r := NewRank()

	basic := NewBasicRank()
	for _, tc := range []struct {
		name string
		got  BasicRank
	}{
		{"Attack", r.Attack()},
		{"Defence", r.Defence()},
		{"SpAttack", r.SpAttack()},
		{"SpDefence", r.SpDefence()},
		{"Speed", r.Speed()},
		{"Accuracy", r.Accuracy()},
		{"Evasion", r.Evasion()},
	} {
		if tc.got != basic {
			t.Errorf("%s: got %+v, want %+v", tc.name, tc.got, basic)
		}
	}

	if r.Critical() != NewCriticalRank() {
		t.Errorf("Critical: got %+v, want %+v", r.Critical(), NewCriticalRank())
	}
}

func TestRank_RollCritical(t *testing.T) {
	t.Run("stage3(必中): 常にtrueを返す", func(t *testing.T) {
		r := Rank{critical: CriticalRank{stage: 3, value: criticalRankMap[3]}}
		for range 20 {
			got, err := r.RollCritical(0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got {
				t.Error("stage3は必ず急所に当たる")
			}
		}
	})

	t.Run("stage0(4.17%): 急所あり・なしの両方が発生する", func(t *testing.T) {
		r := NewRank()
		gotTrue, gotFalse := false, false
		for range 300 {
			got, err := r.RollCritical(0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got {
				gotTrue = true
			} else {
				gotFalse = true
			}
			if gotTrue && gotFalse {
				break
			}
		}
		if !gotTrue {
			t.Error("300試行で急所が一度も発生しなかった")
		}
		if !gotFalse {
			t.Error("300試行で非急所が一度も発生しなかった")
		}
	})

	t.Run("stage2(50%): 急所あり・なしの両方が発生する", func(t *testing.T) {
		r := Rank{critical: CriticalRank{stage: 2, value: criticalRankMap[2]}}
		gotTrue, gotFalse := false, false
		for range 300 {
			got, err := r.RollCritical(0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got {
				gotTrue = true
			} else {
				gotFalse = true
			}
			if gotTrue && gotFalse {
				break
			}
		}
		if !gotTrue {
			t.Error("300試行で急所が一度も発生しなかった")
		}
		if !gotFalse {
			t.Error("300試行で非急所が一度も発生しなかった")
		}
	})

	t.Run("moveCriticalStage加算: stage0+3で必中になる", func(t *testing.T) {
		r := NewRank()
		for range 20 {
			got, err := r.RollCritical(3)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got {
				t.Error("moveCriticalStage=3加算でstage3(必中)になる")
			}
		}
	})

	t.Run("stage0の急所率がおよそ4.17%である", func(t *testing.T) {
		const trials = 3000
		r := NewRank()
		hits := 0
		for range trials {
			got, err := r.RollCritical(0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got {
				hits++
			}
		}
		rate := float64(hits) / trials
		if rate < 0.02 || rate > 0.07 {
			t.Errorf("想定外の急所率: %.4f (期待値 ~0.0417)", rate)
		}
	})

	t.Run("RollCritical後もstageは変化しない", func(t *testing.T) {
		r := Rank{critical: CriticalRank{stage: 1, value: criticalRankMap[1]}}
		if _, err := r.RollCritical(1); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.critical.stage != 1 {
			t.Errorf("expected stage=1, got %d", r.critical.stage)
		}
	})
}

func TestRankAccuracyRank(t *testing.T) {
	tests := []struct {
		name            string
		selfAccStage    int
		opponentEvStage int
		wantValue       [2]int
	}{
		{
			name:      "命中・回避ともにstage=0: 補正なし",
			wantValue: [2]int{3, 3},
		},
		{
			name:            "命中+2, 回避±0: stage=2",
			selfAccStage:    2,
			wantValue:       [2]int{5, 3},
		},
		{
			name:            "命中±0, 回避+3: stage=-3",
			opponentEvStage: 3,
			wantValue:       [2]int{3, 6},
		},
		{
			name:            "命中+2, 回避+1: stage=1",
			selfAccStage:    2, opponentEvStage: 1,
			wantValue: [2]int{4, 3},
		},
		{
			// selfAcc=6, opponentEv=-3 → h-e=9 → clamp 6
			name:            "上限クランプ",
			selfAccStage:    6, opponentEvStage: -3,
			wantValue: [2]int{9, 3},
		},
		{
			// selfAcc=-3, opponentEv=6 → h-e=-9 → clamp -6
			name:            "下限クランプ",
			selfAccStage:    -3, opponentEvStage: 6,
			wantValue: [2]int{3, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := Rank{accuracy: BasicRank{stage: tt.selfAccStage, value: basicRankMap[tt.selfAccStage]}}
			oe := BasicRank{stage: tt.opponentEvStage, value: basicRankMap[tt.opponentEvStage]}

			got, err := self.AccuracyRank(oe)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Value() != tt.wantValue {
				t.Errorf("Value(): got %v, want %v", got.Value(), tt.wantValue)
			}
		})
	}
}
