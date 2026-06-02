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

type BoxRepository struct {
	db *shared.DBClient
}

func NewBoxRepository(db *shared.DBClient) *BoxRepository {
	return &BoxRepository{db: db}
}

func (r *BoxRepository) Create(ctx context.Context, box model.Box) error {
	query := `INSERT INTO boxes (id, user_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.GetClient().Exec(ctx, query, box.BoxId, box.UserId, box.Name, box.CreatedAt, box.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create box", slog.Any("error", err))
		return err
	}
	return nil
}

func (r *BoxRepository) FindById(ctx context.Context, id uuid.UUID) (*model.Box, error) {
	query := `SELECT id, user_id, name, created_at, updated_at FROM boxes WHERE id = $1`
	row := r.db.GetClient().QueryRow(ctx, query, id)

	var box model.Box
	err := row.Scan(&box.BoxId, &box.UserId, &box.Name, &box.CreatedAt, &box.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		slog.ErrorContext(ctx, "failed to find box by id", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	return &box, nil
}

func (r *BoxRepository) FindByUserId(ctx context.Context, userId uuid.UUID) ([]model.Box, error) {
	query := `SELECT id, user_id, name, created_at, updated_at FROM boxes WHERE user_id = $1`
	rows, err := r.db.GetClient().Query(ctx, query, userId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find boxes by user_id", slog.String("user_id", userId.String()), slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()

	var boxes []model.Box
	for rows.Next() {
		var box model.Box
		if err := rows.Scan(&box.BoxId, &box.UserId, &box.Name, &box.CreatedAt, &box.UpdatedAt); err != nil {
			return nil, err
		}
		boxes = append(boxes, box)
	}
	return boxes, rows.Err()
}

func (r *BoxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM boxes WHERE id = $1`
	_, err := r.db.GetClient().Exec(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete box", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	return nil
}

func (r *BoxRepository) UpdateName(ctx context.Context, id uuid.UUID, name string) error {
    query := `UPDATE boxes SET name = $2, updated_at = NOW() WHERE id = $1`
    _, err := r.db.GetClient().Exec(ctx, query, id, name)
    if err != nil {
        slog.ErrorContext(ctx, "failed to update box name", slog.String("id", id.String()), slog.Any("error", err))
        return err
    }
    return nil
}