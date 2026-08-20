package model

import "testing"

func TestGenderValid(t *testing.T) {
	tests := []struct {
		name string
		g    Gender
		want bool
	}{
		{"不明（0）: valid", GenderUnknown, true},
		{"オス（1）: valid", GenderMale, true},
		{"メス（2）: valid", GenderFemale, true},
		{"負値: invalid", Gender(-1), false},
		{"上限超過: invalid", Gender(3), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Valid(); got != tt.want {
				t.Errorf("Gender(%d).Valid() = %v, want %v", tt.g, got, tt.want)
			}
		})
	}
}
