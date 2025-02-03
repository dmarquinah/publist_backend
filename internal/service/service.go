package service

import (
	"context"

	"github.com/dmarquinah/publist_backend/internal/model"
	"github.com/dmarquinah/publist_backend/internal/repository"
)

type Service interface {
	CreateItem(ctx context.Context, item *model.Item) error
	GetItem(ctx context.Context, id string) (*model.Item, error)
	// Add more service methods as needed
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateItem(ctx context.Context, item *model.Item) error {
	// Add business logic here
	return s.repo.CreateItem(ctx, item)
}

func (s *service) GetItem(ctx context.Context, id string) (*model.Item, error) {
	// Add business logic here
	return s.repo.GetItem(ctx, id)
}
