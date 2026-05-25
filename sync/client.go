package sync

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ApiClient struct {
	baseURL *url.URL
	client  *http.Client
}

func NewApiClient(baseURL string) (*ApiClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 10
	t.MaxConnsPerHost = 10
	t.MaxIdleConnsPerHost = 10
	c := &ApiClient{
		baseURL: u,
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		},
	}
	return c, nil
}

func (a *ApiClient) GetClient() *http.Client {
	return a.client
}

type DbClient struct {
	client *pgxpool.Pool
}

func NewDbClient(ctx context.Context, connStr string) (*DbClient, error) {
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

	return &DbClient{client: pool}, nil
}

func (d *DbClient) Close() {
	d.client.Close()
}

func (d *DbClient) GetClient() *pgxpool.Pool {
	return d.client
}
