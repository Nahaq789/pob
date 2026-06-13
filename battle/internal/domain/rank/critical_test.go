package rank

import "testing"

func TestNewCriticalRank(t *testing.T) {
	r := NewCriticalRank()
	if r.stage != 0 {
		t.Errorf("stage: got %d, want 0", r.stage)
	}
	if r.value != 0.417 {
		t.Errorf("value: got %v, want 0.417", r.value)
	}
}

func TestCriticalUp(t *testing.T) {
	tests := []struct {
		name         string
		initialStage int
		up           int
		wantStage    int
		wantValue    float64
	}{
		{
			name: "1段上昇",
			up:   1, wantStage: 1, wantValue: 0.125,
		},
		{
			name: "2段上昇",
			up:   2, wantStage: 2, wantValue: 0.5,
		},
		{
			name: "3段上昇",
			up:   3, wantStage: 3, wantValue: 1.0,
		},
		{
			// stage=2 + up=2 → 4 → clamp 3
			name:         "上限クランプ",
			initialStage: 2, up: 2, wantStage: 3, wantValue: 1.0,
		},
		{
			// stage=1 + up=5 → 6 → clamp 3
			name:         "上限クランプ: 大きな値",
			initialStage: 1, up: 5, wantStage: 3, wantValue: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CriticalRank{stage: tt.initialStage, value: criticalRankMap[tt.initialStage]}
			got, err := r.Up(tt.up)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.stage != tt.wantStage {
				t.Errorf("stage: got %d, want %d", got.stage, tt.wantStage)
			}
			if got.value != tt.wantValue {
				t.Errorf("value: got %v, want %v", got.value, tt.wantValue)
			}
		})
	}
}

func TestCriticalDown(t *testing.T) {
	tests := []struct {
		name         string
		initialStage int
		down         int
		wantStage    int
		wantValue    float64
	}{
		{
			name:         "1段下降",
			initialStage: 3, down: 1, wantStage: 2, wantValue: 0.5,
		},
		{
			name:         "2段下降",
			initialStage: 3, down: 2, wantStage: 1, wantValue: 0.125,
		},
		{
			name:         "3段下降",
			initialStage: 3, down: 3, wantStage: 0, wantValue: 0.417,
		},
		{
			// stage=1 + down=2 → -1 → clamp 0
			name:         "下限クランプ",
			initialStage: 1, down: 2, wantStage: 0, wantValue: 0.417,
		},
		{
			// stage=0 + down=1 → -1 → clamp 0
			name:         "下限クランプ: stage=0から",
			initialStage: 0, down: 1, wantStage: 0, wantValue: 0.417,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CriticalRank{stage: tt.initialStage, value: criticalRankMap[tt.initialStage]}
			got, err := r.Down(tt.down)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.stage != tt.wantStage {
				t.Errorf("stage: got %d, want %d", got.stage, tt.wantStage)
			}
			if got.value != tt.wantValue {
				t.Errorf("value: got %v, want %v", got.value, tt.wantValue)
			}
		})
	}
}

func TestCriticalReset(t *testing.T) {
	r := CriticalRank{stage: 2, value: criticalRankMap[2]}
	got := r.Reset()
	if got.stage != 0 {
		t.Errorf("stage: got %d, want 0", got.stage)
	}
	if got.value != 0.417 {
		t.Errorf("value: got %v, want 0.417", got.value)
	}
}
