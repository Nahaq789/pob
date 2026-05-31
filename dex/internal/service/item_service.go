package service

import (
	"context"

	"pob/dex/internal/model"
	"pob/dex/internal/repository"
)

type ItemService struct {
	itemRepo *repository.ItemRepository
}

func NewItemService(itemRepo *repository.ItemRepository) *ItemService {
	return &ItemService{itemRepo: itemRepo}
}

func (s *ItemService) GetItem(ctx context.Context, id int) (*model.Item, error) {
	return s.itemRepo.FindById(ctx, id)
}
