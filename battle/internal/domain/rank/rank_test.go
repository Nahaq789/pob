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
