package rank

import "testing"

func TestNewAccuracyRank(t *testing.T) {
	tests := []struct {
		name      string
		h, e      int
		wantStage int
		wantValue [2]int
	}{
		{name: "h=e=0: stage=0", h: 0, e: 0, wantStage: 0, wantValue: [2]int{3, 3}},
		{name: "h>e: stage=1", h: 1, e: 0, wantStage: 1, wantValue: [2]int{4, 3}},
		{name: "h>e: stage=2", h: 2, e: 0, wantStage: 2, wantValue: [2]int{5, 3}},
		{name: "h>e: stage=3", h: 3, e: 0, wantStage: 3, wantValue: [2]int{6, 3}},
		{name: "h>e: stage=4", h: 4, e: 0, wantStage: 4, wantValue: [2]int{7, 3}},
		{name: "h>e: stage=5", h: 5, e: 0, wantStage: 5, wantValue: [2]int{8, 3}},
		{name: "h>e: stage=6", h: 6, e: 0, wantStage: 6, wantValue: [2]int{9, 3}},
		{name: "h<e: stage=-1", h: 0, e: 1, wantStage: -1, wantValue: [2]int{3, 4}},
		{name: "h<e: stage=-2", h: 0, e: 2, wantStage: -2, wantValue: [2]int{3, 5}},
		{name: "h<e: stage=-3", h: 0, e: 3, wantStage: -3, wantValue: [2]int{3, 6}},
		{name: "h<e: stage=-4", h: 0, e: 4, wantStage: -4, wantValue: [2]int{3, 7}},
		{name: "h<e: stage=-5", h: 0, e: 5, wantStage: -5, wantValue: [2]int{3, 8}},
		{name: "h<e: stage=-6", h: 0, e: 6, wantStage: -6, wantValue: [2]int{3, 9}},
		{
			// h-e=9 → clamp 6
			name: "上限クランプ",
			h: 6, e: -3, wantStage: 6, wantValue: [2]int{9, 3},
		},
		{
			// h-e=-9 → clamp -6
			name: "下限クランプ",
			h: -3, e: 6, wantStage: -6, wantValue: [2]int{3, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccuracyRank(tt.h, tt.e)
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

func TestAccuracyRankValue(t *testing.T) {
	// h=5, e=3 → stage=2, value={5,3}
	r, err := NewAccuracyRank(5, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := [2]int{5, 3}
	if got := r.Value(); got != want {
		t.Errorf("Value(): got %v, want %v", got, want)
	}
}
