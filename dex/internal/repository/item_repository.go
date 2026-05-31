package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pob/dex/internal/model"
	"pob/dex/internal/model/entity"
	"pob/dex/internal/shared"
	pkgredis "pob/pkg/redis"
)

const itemCacheTTL = 24 * time.Hour

type ItemRepository struct {
	db    *shared.DBClient
	redis *pkgredis.RedisClient
}

func NewItemRepository(db *shared.DBClient, redis *pkgredis.RedisClient) *ItemRepository {
	return &ItemRepository{db: db, redis: redis}
}

func (r *ItemRepository) cacheKey(id int) string {
	return fmt.Sprintf("item:%d", id)
}

func (r *ItemRepository) FindById(ctx context.Context, id int) (*model.Item, error) {
	key := r.cacheKey(id)

	if cached, err := r.redis.GetClient().Get(ctx, key).Bytes(); err == nil {
		var item model.Item
		if err := json.Unmarshal(cached, &item); err == nil {
			return &item, nil
		}
	}

	var e entity.Item
	if err := r.db.GetClient().WithContext(ctx).First(&e, id).Error; err != nil {
		return nil, err
	}

	item := &model.Item{
		Id:         e.Id,
		Name:       e.Name,
		Category:   e.Category,
		FlavorText: e.FlavorText,
	}

	if b, err := json.Marshal(item); err == nil {
		r.redis.GetClient().Set(ctx, key, b, itemCacheTTL)
	}

	return item, nil
}
