package rank

import "testing"

func TestNewBasicRank(t *testing.T) {
	r := NewBasicRank()
	if r.stage != 0 {
		t.Errorf("stage: got %d, want 0", r.stage)
	}
	if r.value != 1.0 {
		t.Errorf("value: got %f, want 1.0", r.value)
	}
}

func TestUp(t *testing.T) {
	tests := []struct {
		name         string
		initialStage int
		up           int
		wantStage    int
		wantValue    float64
		wantMsg      string
	}{
		{
			name: "1段上昇",
			up:   1, wantStage: 1, wantValue: 1.5, wantMsg: "が あがった",
		},
		{
			name: "2段上昇",
			up:   2, wantStage: 2, wantValue: 2.0, wantMsg: "が ぐーんとあがった",
		},
		{
			name: "3段上昇",
			up:   3, wantStage: 3, wantValue: 2.5, wantMsg: "が ぐぐーんとあがった",
		},
		{
			name: "4段上昇",
			up:   4, wantStage: 4, wantValue: 3.0, wantMsg: "が ぐぐーんとあがった",
		},
		{
			// stage=5 + up=3 → 8 → clamp 6, delta=1
			name:         "上限クランプ: delta=1",
			initialStage: 5, up: 3, wantStage: 6, wantValue: 4.0, wantMsg: "が あがった",
		},
		{
			// stage=4 + up=4 → 8 → clamp 6, delta=2
			name:         "上限クランプ: delta=2",
			initialStage: 4, up: 4, wantStage: 6, wantValue: 4.0, wantMsg: "が ぐーんとあがった",
		},
		{
			// stage=3 + up=4 → 7 → clamp 6, delta=3
			name:         "上限クランプ: delta=3以上",
			initialStage: 3, up: 4, wantStage: 6, wantValue: 4.0, wantMsg: "が ぐぐーんとあがった",
		},
		{
			name: "MaxStage: stage=0から強制最大",
			up:   BasicMaxStage, wantStage: 6, wantValue: 4.0, wantMsg: "が さいだいまであがった",
		},
		{
			name:         "MaxStage: 途中から強制最大",
			initialStage: 3, up: BasicMaxStage, wantStage: 6, wantValue: 4.0, wantMsg: "が さいだいまであがった",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := BasicRank{stage: tt.initialStage, value: basicRankMap[tt.initialStage]}
			got, msg, err := r.Up(tt.up)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.stage != tt.wantStage {
				t.Errorf("stage: got %d, want %d", got.stage, tt.wantStage)
			}
			if got.value != tt.wantValue {
				t.Errorf("value: got %v, want %v", got.value, tt.wantValue)
			}
			if msg != tt.wantMsg {
				t.Errorf("message: got %q, want %q", msg, tt.wantMsg)
			}
		})
	}
}

func TestDown(t *testing.T) {
	tests := []struct {
		name         string
		initialStage int
		down         int
		wantStage    int
		wantValue    float64
		wantMsg      string
	}{
		{
			name: "1段下降",
			down: 1, wantStage: -1, wantValue: 0.667, wantMsg: "が さがった",
		},
		{
			name: "2段下降",
			down: 2, wantStage: -2, wantValue: 0.5, wantMsg: "が がくっとさがった",
		},
		{
			name: "3段下降",
			down: 3, wantStage: -3, wantValue: 0.4, wantMsg: "が がくーんとさがった",
		},
		{
			name: "4段下降",
			down: 4, wantStage: -4, wantValue: 0.33, wantMsg: "が がくーんとさがった",
		},
		{
			// stage=-5 + down=3 → -8 → clamp -6, delta=-1
			name:         "下限クランプ: delta=-1",
			initialStage: -5, down: 3, wantStage: -6, wantValue: 0.25, wantMsg: "が さがった",
		},
		{
			// stage=-4 + down=4 → -8 → clamp -6, delta=-2
			name:         "下限クランプ: delta=-2",
			initialStage: -4, down: 4, wantStage: -6, wantValue: 0.25, wantMsg: "が がくっとさがった",
		},
		{
			// stage=-3 + down=4 → -7 → clamp -6, delta=-3
			name:         "下限クランプ: delta=-3以下",
			initialStage: -3, down: 4, wantStage: -6, wantValue: 0.25, wantMsg: "が がくーんとさがった",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := BasicRank{stage: tt.initialStage, value: basicRankMap[tt.initialStage]}
			got, msg, err := r.Down(tt.down)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.stage != tt.wantStage {
				t.Errorf("stage: got %d, want %d", got.stage, tt.wantStage)
			}
			if got.value != tt.wantValue {
				t.Errorf("value: got %v, want %v", got.value, tt.wantValue)
			}
			if msg != tt.wantMsg {
				t.Errorf("message: got %q, want %q", msg, tt.wantMsg)
			}
		})
	}
}

func TestReset(t *testing.T) {
	r := BasicRank{stage: 4, value: basicRankMap[4]}
	got := r.Reset()
	if got.stage != 0 {
		t.Errorf("stage: got %d, want 0", got.stage)
	}
	if got.value != 1.0 {
		t.Errorf("value: got %v, want 1.0", got.value)
	}
}
