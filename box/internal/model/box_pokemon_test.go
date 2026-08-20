package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewBoxPokemon(t *testing.T) {
	t.Run("BoxPokemonId が non-nil UUID で生成される", func(t *testing.T) {
		bp := NewBoxPokemon(uuid.New(), 1, 1, "ようき", GenderMale)
		if bp.BoxPokemonId == uuid.Nil {
			t.Error("expected non-nil BoxPokemonId")
		}
	})

	t.Run("BoxId が引数の値で設定される", func(t *testing.T) {
		boxId := uuid.New()
		bp := NewBoxPokemon(boxId, 1, 1, "ようき", GenderMale)
		if bp.BoxId != boxId {
			t.Errorf("BoxId = %v, want %v", bp.BoxId, boxId)
		}
	})

	t.Run("PokemonId・AbilityId・Nature・Gender が引数の値で設定される", func(t *testing.T) {
		bp := NewBoxPokemon(uuid.New(), 25, 65, "いじっぱり", GenderFemale)
		if bp.PokemonId != 25 {
			t.Errorf("PokemonId = %d, want 25", bp.PokemonId)
		}
		if bp.AbilityId != 65 {
			t.Errorf("AbilityId = %d, want 65", bp.AbilityId)
		}
		if bp.Nature != "いじっぱり" {
			t.Errorf("Nature = %q, want %q", bp.Nature, "いじっぱり")
		}
		if bp.Gender != GenderFemale {
			t.Errorf("Gender = %v, want GenderFemale", bp.Gender)
		}
	})

	t.Run("IV が全て 31 で初期化される", func(t *testing.T) {
		bp := NewBoxPokemon(uuid.New(), 1, 1, "ようき", GenderMale)
		ivs := []struct {
			name string
			val  int
		}{
			{"IvHp", bp.IvHp},
			{"IvAttack", bp.IvAttack},
			{"IvDefense", bp.IvDefense},
			{"IvSpAttack", bp.IvSpAttack},
			{"IvSpDefense", bp.IvSpDefense},
			{"IvSpeed", bp.IvSpeed},
		}
		for _, iv := range ivs {
			if iv.val != 31 {
				t.Errorf("%s = %d, want 31", iv.name, iv.val)
			}
		}
	})

	t.Run("EV が全て 0 で初期化される", func(t *testing.T) {
		bp := NewBoxPokemon(uuid.New(), 1, 1, "ようき", GenderMale)
		evs := []struct {
			name string
			val  int
		}{
			{"EvHp", bp.EvHp},
			{"EvAttack", bp.EvAttack},
			{"EvDefense", bp.EvDefense},
			{"EvSpAttack", bp.EvSpAttack},
			{"EvSpDefense", bp.EvSpDefense},
			{"EvSpeed", bp.EvSpeed},
		}
		for _, ev := range evs {
			if ev.val != 0 {
				t.Errorf("%s = %d, want 0", ev.name, ev.val)
			}
		}
	})

	t.Run("Nickname・HeldItemId・Move*Id が nil で初期化される", func(t *testing.T) {
		bp := NewBoxPokemon(uuid.New(), 1, 1, "ようき", GenderMale)
		if bp.Nickname != nil {
			t.Errorf("Nickname = %v, want nil", bp.Nickname)
		}
		if bp.HeldItemId != nil {
			t.Errorf("HeldItemId = %v, want nil", bp.HeldItemId)
		}
		if bp.Move1Id != nil || bp.Move2Id != nil || bp.Move3Id != nil || bp.Move4Id != nil {
			t.Error("expected all MoveIds to be nil")
		}
	})

	t.Run("CreatedAt と UpdatedAt がゼロ値でない", func(t *testing.T) {
		bp := NewBoxPokemon(uuid.New(), 1, 1, "ようき", GenderMale)
		if bp.CreatedAt.IsZero() {
			t.Error("expected non-zero CreatedAt")
		}
		if bp.UpdatedAt.IsZero() {
			t.Error("expected non-zero UpdatedAt")
		}
	})

	t.Run("複数回呼び出すと異なる BoxPokemonId が生成される", func(t *testing.T) {
		boxId := uuid.New()
		bp1 := NewBoxPokemon(boxId, 1, 1, "ようき", GenderMale)
		bp2 := NewBoxPokemon(boxId, 1, 1, "ようき", GenderMale)
		if bp1.BoxPokemonId == bp2.BoxPokemonId {
			t.Error("expected different BoxPokemonIds")
		}
	})
}

func TestFromBoxPokemon(t *testing.T) {
	t.Run("全フィールドが引数の値で設定される", func(t *testing.T) {
		id := uuid.New()
		boxId := uuid.New()
		nickname := "ピカ"
		heldItem := 1
		move1, move2, move3, move4 := 10, 20, 30, 40
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

		bp := FromBoxPokemon(
			id, boxId,
			25, &nickname,
			65, "いじっぱり", GenderFemale, &heldItem,
			31, 30, 29, 28, 27, 26,
			0, 252, 0, 0, 4, 0,
			&move1, &move2, &move3, &move4,
			createdAt, updatedAt,
		)

		if bp.BoxPokemonId != id {
			t.Errorf("BoxPokemonId = %v, want %v", bp.BoxPokemonId, id)
		}
		if bp.BoxId != boxId {
			t.Errorf("BoxId = %v, want %v", bp.BoxId, boxId)
		}
		if bp.PokemonId != 25 {
			t.Errorf("PokemonId = %d, want 25", bp.PokemonId)
		}
		if bp.Nickname == nil || *bp.Nickname != "ピカ" {
			t.Errorf("Nickname = %v, want %q", bp.Nickname, "ピカ")
		}
		if bp.AbilityId != 65 {
			t.Errorf("AbilityId = %d, want 65", bp.AbilityId)
		}
		if bp.Nature != "いじっぱり" {
			t.Errorf("Nature = %q, want %q", bp.Nature, "いじっぱり")
		}
		if bp.Gender != GenderFemale {
			t.Errorf("Gender = %v, want GenderFemale", bp.Gender)
		}
		if bp.HeldItemId == nil || *bp.HeldItemId != 1 {
			t.Errorf("HeldItemId = %v, want 1", bp.HeldItemId)
		}
		if bp.IvHp != 31 || bp.IvAttack != 30 || bp.IvDefense != 29 {
			t.Errorf("IVs not set correctly: hp=%d atk=%d def=%d", bp.IvHp, bp.IvAttack, bp.IvDefense)
		}
		if bp.EvAttack != 252 || bp.EvSpDefense != 4 {
			t.Errorf("EVs not set correctly: atk=%d spd=%d", bp.EvAttack, bp.EvSpDefense)
		}
		if bp.Move1Id == nil || *bp.Move1Id != 10 {
			t.Errorf("Move1Id = %v, want 10", bp.Move1Id)
		}
		if bp.Move4Id == nil || *bp.Move4Id != 40 {
			t.Errorf("Move4Id = %v, want 40", bp.Move4Id)
		}
		if !bp.CreatedAt.Equal(createdAt) {
			t.Errorf("CreatedAt = %v, want %v", bp.CreatedAt, createdAt)
		}
		if !bp.UpdatedAt.Equal(updatedAt) {
			t.Errorf("UpdatedAt = %v, want %v", bp.UpdatedAt, updatedAt)
		}
	})

	t.Run("nil ポインタフィールドがそのまま保持される", func(t *testing.T) {
		bp := FromBoxPokemon(
			uuid.New(), uuid.New(),
			1, nil,
			1, "ようき", GenderMale, nil,
			31, 31, 31, 31, 31, 31,
			0, 0, 0, 0, 0, 0,
			nil, nil, nil, nil,
			time.Now(), time.Now(),
		)
		if bp.Nickname != nil {
			t.Errorf("Nickname = %v, want nil", bp.Nickname)
		}
		if bp.HeldItemId != nil {
			t.Errorf("HeldItemId = %v, want nil", bp.HeldItemId)
		}
		if bp.Move1Id != nil || bp.Move2Id != nil || bp.Move3Id != nil || bp.Move4Id != nil {
			t.Error("expected all MoveIds to be nil")
		}
	})
}
