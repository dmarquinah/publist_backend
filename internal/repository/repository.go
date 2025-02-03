package repository

import (
	"context"
	"sync"

	"github.com/dmarquinah/publist_backend/internal/errors"
	"github.com/dmarquinah/publist_backend/internal/model"
)

type Repository interface {
	CreateItem(ctx context.Context, item *model.Item) error
	GetItem(ctx context.Context, id string) (*model.Item, error)
	// Add more repository methods as needed
}

type repository struct {
	items map[string]*model.Item
	mu    sync.RWMutex
}

func NewRepository() Repository {
	return &repository{
		items: make(map[string]*model.Item),
	}
}

func (r *repository) CreateItem(ctx context.Context, item *model.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[item.ID] = item
	return nil
}

func (r *repository) GetItem(ctx context.Context, id string) (*model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, errors.ErrItemNotFound
	}
	return item, nil
}
