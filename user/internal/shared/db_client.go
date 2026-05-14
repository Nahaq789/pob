package shared

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBClient struct {
	client *pgxpool.Pool
}

func InitDbClient(ctx context.Context, connStr string) (*DBClient, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &DBClient{
		client: pool,
	}, nil
}

func (d *DBClient) Close() {
	d.client.Close()
}

func (d *DBClient) GetClient() *pgxpool.Pool {
	return d.client
}
