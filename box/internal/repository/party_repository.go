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

type PartyRepository struct {
	db *shared.DBClient
}

func NewPartyRepository(db *shared.DBClient) *PartyRepository {
	return &PartyRepository{db: db}
}

func (r *PartyRepository) Create(ctx context.Context, party model.Party) error {
	query := `INSERT INTO parties (id, user_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.GetClient().Exec(ctx, query, party.PartyId, party.UserId, party.Name, party.CreatedAt, party.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create party", slog.Any("error", err))
		return err
	}
	return nil
}

func (r *PartyRepository) FindById(ctx context.Context, id uuid.UUID) (*model.Party, error) {
	query := `SELECT id, user_id, name, created_at, updated_at FROM parties WHERE id = $1`
	row := r.db.GetClient().QueryRow(ctx, query, id)

	var party model.Party
	err := row.Scan(&party.PartyId, &party.UserId, &party.Name, &party.CreatedAt, &party.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		slog.ErrorContext(ctx, "failed to find party by id", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	return &party, nil
}

func (r *PartyRepository) FindByUserId(ctx context.Context, userId uuid.UUID) ([]model.Party, error) {
	query := `SELECT id, user_id, name, created_at, updated_at FROM parties WHERE user_id = $1`
	rows, err := r.db.GetClient().Query(ctx, query, userId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find parties by user_id", slog.String("user_id", userId.String()), slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()

	var parties []model.Party
	for rows.Next() {
		var party model.Party
		if err := rows.Scan(&party.PartyId, &party.UserId, &party.Name, &party.CreatedAt, &party.UpdatedAt); err != nil {
			return nil, err
		}
		parties = append(parties, party)
	}
	return parties, rows.Err()
}

func (r *PartyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM parties WHERE id = $1`
	_, err := r.db.GetClient().Exec(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete party", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	return nil
}
