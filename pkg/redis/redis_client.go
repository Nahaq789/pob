package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(ctx context.Context, addr string) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:           addr,
		MaxActiveConns: 10,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{
		client: rdb,
	}, nil
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

func (r *RedisClient) Close() error {
	if err := r.client.Close(); err != nil {
		return err
	}
	return nil
}
