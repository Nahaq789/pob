package repository

import (
	"context"
	"errors"
	"log/slog"

	"pob/box/internal/model"
	"pob/box/internal/shared"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type BoxPokemonRepository struct {
	db *shared.DBClient
}

func NewBoxPokemonRepository(db *shared.DBClient) *BoxPokemonRepository {
	return &BoxPokemonRepository{db: db}
}

func (r *BoxPokemonRepository) Create(ctx context.Context, bp model.BoxPokemon) error {
	query := `
		INSERT INTO box_pokemon (
			id, box_id, pokemon_id, nickname, ability_id, nature, held_item_id,
			iv_hp, iv_attack, iv_defense, iv_sp_attack, iv_sp_defense, iv_speed,
			ev_hp, ev_attack, ev_defense, ev_sp_attack, ev_sp_defense, ev_speed,
			move1_id, move2_id, move3_id, move4_id,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18, $19,
			$20, $21, $22, $23,
			$24, $25
		)`
	_, err := r.db.GetClient().Exec(ctx, query,
		bp.BoxPokemonId, bp.BoxId, bp.PokemonId, bp.Nickname, bp.AbilityId, bp.Nature, bp.HeldItemId,
		bp.IvHp, bp.IvAttack, bp.IvDefense, bp.IvSpAttack, bp.IvSpDefense, bp.IvSpeed,
		bp.EvHp, bp.EvAttack, bp.EvDefense, bp.EvSpAttack, bp.EvSpDefense, bp.EvSpeed,
		bp.Move1Id, bp.Move2Id, bp.Move3Id, bp.Move4Id,
		bp.CreatedAt, bp.UpdatedAt,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create box_pokemon", slog.Any("error", err))
		return err
	}
	return nil
}

func (r *BoxPokemonRepository) FindById(ctx context.Context, id uuid.UUID) (*model.BoxPokemon, error) {
	query := `
		SELECT
			id, box_id, pokemon_id, nickname, ability_id, nature, held_item_id,
			iv_hp, iv_attack, iv_defense, iv_sp_attack, iv_sp_defense, iv_speed,
			ev_hp, ev_attack, ev_defense, ev_sp_attack, ev_sp_defense, ev_speed,
			move1_id, move2_id, move3_id, move4_id,
			created_at, updated_at
		FROM box_pokemon WHERE id = $1`
	row := r.db.GetClient().QueryRow(ctx, query, id)
	bp, err := scanBoxPokemon(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		slog.ErrorContext(ctx, "failed to find box_pokemon by id", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	return bp, nil
}

func (r *BoxPokemonRepository) FindByBoxId(ctx context.Context, boxId uuid.UUID) ([]model.BoxPokemon, error) {
	query := `
		SELECT
			id, box_id, pokemon_id, nickname, ability_id, nature, held_item_id,
			iv_hp, iv_attack, iv_defense, iv_sp_attack, iv_sp_defense, iv_speed,
			ev_hp, ev_attack, ev_defense, ev_sp_attack, ev_sp_defense, ev_speed,
			move1_id, move2_id, move3_id, move4_id,
			created_at, updated_at
		FROM box_pokemon WHERE box_id = $1`
	rows, err := r.db.GetClient().Query(ctx, query, boxId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find box_pokemon by box_id", slog.String("box_id", boxId.String()), slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()

	var result []model.BoxPokemon
	for rows.Next() {
		bp, err := scanBoxPokemon(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *bp)
	}
	return result, rows.Err()
}

func (r *BoxPokemonRepository) Update(ctx context.Context, bp model.BoxPokemon) error {
	query := `
		UPDATE box_pokemon SET
			nickname = $2, ability_id = $3, nature = $4, held_item_id = $5,
			iv_hp = $6, iv_attack = $7, iv_defense = $8, iv_sp_attack = $9, iv_sp_defense = $10, iv_speed = $11,
			ev_hp = $12, ev_attack = $13, ev_defense = $14, ev_sp_attack = $15, ev_sp_defense = $16, ev_speed = $17,
			move1_id = $18, move2_id = $19, move3_id = $20, move4_id = $21,
			updated_at = $22
		WHERE id = $1`
	_, err := r.db.GetClient().Exec(ctx, query,
		bp.BoxPokemonId, bp.Nickname, bp.AbilityId, bp.Nature, bp.HeldItemId,
		bp.IvHp, bp.IvAttack, bp.IvDefense, bp.IvSpAttack, bp.IvSpDefense, bp.IvSpeed,
		bp.EvHp, bp.EvAttack, bp.EvDefense, bp.EvSpAttack, bp.EvSpDefense, bp.EvSpeed,
		bp.Move1Id, bp.Move2Id, bp.Move3Id, bp.Move4Id,
		bp.UpdatedAt,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update box_pokemon", slog.String("id", bp.BoxPokemonId.String()), slog.Any("error", err))
		return err
	}
	return nil
}

func (r *BoxPokemonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM box_pokemon WHERE id = $1`
	_, err := r.db.GetClient().Exec(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete box_pokemon", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanBoxPokemon(s scanner) (*model.BoxPokemon, error) {
	var bp model.BoxPokemon
	err := s.Scan(
		&bp.BoxPokemonId, &bp.BoxId, &bp.PokemonId, &bp.Nickname, &bp.AbilityId, &bp.Nature, &bp.HeldItemId,
		&bp.IvHp, &bp.IvAttack, &bp.IvDefense, &bp.IvSpAttack, &bp.IvSpDefense, &bp.IvSpeed,
		&bp.EvHp, &bp.EvAttack, &bp.EvDefense, &bp.EvSpAttack, &bp.EvSpDefense, &bp.EvSpeed,
		&bp.Move1Id, &bp.Move2Id, &bp.Move3Id, &bp.Move4Id,
		&bp.CreatedAt, &bp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &bp, nil
}
