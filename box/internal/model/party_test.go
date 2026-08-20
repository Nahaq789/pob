package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewParty(t *testing.T) {
	t.Run("PartyId が non-nil UUID で生成される", func(t *testing.T) {
		uid := uuid.New()
		p := NewParty(uid, "テストパーティ")
		if p.PartyId == uuid.Nil {
			t.Error("expected non-nil PartyId")
		}
	})

	t.Run("UserId が引数の値で設定される", func(t *testing.T) {
		uid := uuid.New()
		p := NewParty(uid, "テストパーティ")
		if p.UserId != uid {
			t.Errorf("UserId = %v, want %v", p.UserId, uid)
		}
	})

	t.Run("Name が引数の値で設定される", func(t *testing.T) {
		p := NewParty(uuid.New(), "Aパーティ")
		if p.Name != "Aパーティ" {
			t.Errorf("Name = %q, want %q", p.Name, "Aパーティ")
		}
	})

	t.Run("CreatedAt と UpdatedAt がゼロ値でない", func(t *testing.T) {
		p := NewParty(uuid.New(), "テストパーティ")
		if p.CreatedAt.IsZero() {
			t.Error("expected non-zero CreatedAt")
		}
		if p.UpdatedAt.IsZero() {
			t.Error("expected non-zero UpdatedAt")
		}
	})

	t.Run("複数回呼び出すと異なる PartyId が生成される", func(t *testing.T) {
		uid := uuid.New()
		p1 := NewParty(uid, "party1")
		p2 := NewParty(uid, "party2")
		if p1.PartyId == p2.PartyId {
			t.Error("expected different PartyIds")
		}
	})
}

func TestFromParty(t *testing.T) {
	t.Run("全フィールドが引数の値で設定される", func(t *testing.T) {
		id := uuid.New()
		uid := uuid.New()
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

		p := FromParty(id, uid, "パーティB", createdAt, updatedAt)

		if p.PartyId != id {
			t.Errorf("PartyId = %v, want %v", p.PartyId, id)
		}
		if p.UserId != uid {
			t.Errorf("UserId = %v, want %v", p.UserId, uid)
		}
		if p.Name != "パーティB" {
			t.Errorf("Name = %q, want %q", p.Name, "パーティB")
		}
		if !p.CreatedAt.Equal(createdAt) {
			t.Errorf("CreatedAt = %v, want %v", p.CreatedAt, createdAt)
		}
		if !p.UpdatedAt.Equal(updatedAt) {
			t.Errorf("UpdatedAt = %v, want %v", p.UpdatedAt, updatedAt)
		}
	})
}
