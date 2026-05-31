package repository

import (
	"context"
	"log/slog"

	"pob/box/internal/model"
	"pob/box/internal/shared"

	"github.com/google/uuid"
)

type PartyPokemonRepository struct {
	db *shared.DBClient
}

func NewPartyPokemonRepository(db *shared.DBClient) *PartyPokemonRepository {
	return &PartyPokemonRepository{db: db}
}

func (r *PartyPokemonRepository) AddPokemon(ctx context.Context, pp model.PartyPokemon) error {
	query := `INSERT INTO party_pokemon (id, party_id, box_pokemon_id, slot) VALUES ($1, $2, $3, $4)`
	_, err := r.db.GetClient().Exec(ctx, query, pp.PartyPokemonId, pp.PartyId, pp.BoxPokemonId, pp.Slot)
	if err != nil {
		slog.ErrorContext(ctx, "failed to add party_pokemon", slog.Any("error", err))
		return err
	}
	return nil
}

func (r *PartyPokemonRepository) FindByPartyId(ctx context.Context, partyId uuid.UUID) ([]model.PartyPokemon, error) {
	query := `SELECT id, party_id, box_pokemon_id, slot FROM party_pokemon WHERE party_id = $1 ORDER BY slot ASC`
	rows, err := r.db.GetClient().Query(ctx, query, partyId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find party_pokemon by party_id", slog.String("party_id", partyId.String()), slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()

	var result []model.PartyPokemon
	for rows.Next() {
		var pp model.PartyPokemon
		if err := rows.Scan(&pp.PartyPokemonId, &pp.PartyId, &pp.BoxPokemonId, &pp.Slot); err != nil {
			return nil, err
		}
		result = append(result, pp)
	}
	return result, rows.Err()
}

func (r *PartyPokemonRepository) RemovePokemon(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM party_pokemon WHERE id = $1`
	_, err := r.db.GetClient().Exec(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to remove party_pokemon", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	return nil
}

func (r *PartyPokemonRepository) UpdateSlot(ctx context.Context, id uuid.UUID, slot int) error {
	query := `UPDATE party_pokemon SET slot = $2 WHERE id = $1`
	_, err := r.db.GetClient().Exec(ctx, query, id, slot)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update party_pokemon slot", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	return nil
}
