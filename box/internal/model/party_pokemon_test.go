package model

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewPartyPokemon(t *testing.T) {
	t.Run("PartyPokemonId が non-nil UUID で生成される", func(t *testing.T) {
		pp := NewPartyPokemon(uuid.New(), uuid.New(), 1)
		if pp.PartyPokemonId == uuid.Nil {
			t.Error("expected non-nil PartyPokemonId")
		}
	})

	t.Run("PartyId が引数の値で設定される", func(t *testing.T) {
		partyId := uuid.New()
		pp := NewPartyPokemon(partyId, uuid.New(), 1)
		if pp.PartyId != partyId {
			t.Errorf("PartyId = %v, want %v", pp.PartyId, partyId)
		}
	})

	t.Run("BoxPokemonId が引数の値で設定される", func(t *testing.T) {
		boxPokemonId := uuid.New()
		pp := NewPartyPokemon(uuid.New(), boxPokemonId, 2)
		if pp.BoxPokemonId != boxPokemonId {
			t.Errorf("BoxPokemonId = %v, want %v", pp.BoxPokemonId, boxPokemonId)
		}
	})

	t.Run("Slot が引数の値で設定される", func(t *testing.T) {
		for slot := 1; slot <= 6; slot++ {
			pp := NewPartyPokemon(uuid.New(), uuid.New(), slot)
			if pp.Slot != slot {
				t.Errorf("Slot = %d, want %d", pp.Slot, slot)
			}
		}
	})

	t.Run("複数回呼び出すと異なる PartyPokemonId が生成される", func(t *testing.T) {
		pp1 := NewPartyPokemon(uuid.New(), uuid.New(), 1)
		pp2 := NewPartyPokemon(uuid.New(), uuid.New(), 1)
		if pp1.PartyPokemonId == pp2.PartyPokemonId {
			t.Error("expected different PartyPokemonIds")
		}
	})
}

func TestFromPartyPokemon(t *testing.T) {
	t.Run("全フィールドが引数の値で設定される", func(t *testing.T) {
		id := uuid.New()
		partyId := uuid.New()
		boxPokemonId := uuid.New()

		pp := FromPartyPokemon(id, partyId, boxPokemonId, 3)

		if pp.PartyPokemonId != id {
			t.Errorf("PartyPokemonId = %v, want %v", pp.PartyPokemonId, id)
		}
		if pp.PartyId != partyId {
			t.Errorf("PartyId = %v, want %v", pp.PartyId, partyId)
		}
		if pp.BoxPokemonId != boxPokemonId {
			t.Errorf("BoxPokemonId = %v, want %v", pp.BoxPokemonId, boxPokemonId)
		}
		if pp.Slot != 3 {
			t.Errorf("Slot = %d, want 3", pp.Slot)
		}
	})
}
