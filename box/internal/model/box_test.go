package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewBox(t *testing.T) {
	t.Run("BoxId が non-nil UUID で生成される", func(t *testing.T) {
		uid := uuid.New()
		b := NewBox(uid, "テストボックス")
		if b.BoxId == uuid.Nil {
			t.Error("expected non-nil BoxId")
		}
	})

	t.Run("UserId が引数の値で設定される", func(t *testing.T) {
		uid := uuid.New()
		b := NewBox(uid, "テストボックス")
		if b.UserId != uid {
			t.Errorf("UserId = %v, want %v", b.UserId, uid)
		}
	})

	t.Run("Name が引数の値で設定される", func(t *testing.T) {
		b := NewBox(uuid.New(), "マイボックス")
		if b.Name != "マイボックス" {
			t.Errorf("Name = %q, want %q", b.Name, "マイボックス")
		}
	})

	t.Run("CreatedAt と UpdatedAt がゼロ値でない", func(t *testing.T) {
		b := NewBox(uuid.New(), "テストボックス")
		if b.CreatedAt.IsZero() {
			t.Error("expected non-zero CreatedAt")
		}
		if b.UpdatedAt.IsZero() {
			t.Error("expected non-zero UpdatedAt")
		}
	})

	t.Run("複数回呼び出すと異なる BoxId が生成される", func(t *testing.T) {
		uid := uuid.New()
		b1 := NewBox(uid, "box1")
		b2 := NewBox(uid, "box2")
		if b1.BoxId == b2.BoxId {
			t.Error("expected different BoxIds")
		}
	})
}

func TestFromBox(t *testing.T) {
	t.Run("全フィールドが引数の値で設定される", func(t *testing.T) {
		id := uuid.New()
		uid := uuid.New()
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

		b := FromBox(id, uid, "ボックスA", createdAt, updatedAt)

		if b.BoxId != id {
			t.Errorf("BoxId = %v, want %v", b.BoxId, id)
		}
		if b.UserId != uid {
			t.Errorf("UserId = %v, want %v", b.UserId, uid)
		}
		if b.Name != "ボックスA" {
			t.Errorf("Name = %q, want %q", b.Name, "ボックスA")
		}
		if !b.CreatedAt.Equal(createdAt) {
			t.Errorf("CreatedAt = %v, want %v", b.CreatedAt, createdAt)
		}
		if !b.UpdatedAt.Equal(updatedAt) {
			t.Errorf("UpdatedAt = %v, want %v", b.UpdatedAt, updatedAt)
		}
	})
}
