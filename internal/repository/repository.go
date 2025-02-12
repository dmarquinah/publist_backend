package repository

import (
	"database/sql"
	"sync"

	"github.com/dmarquinah/publist_backend/internal/model"
)

type Repository interface {
	GetPlaylistRepository() PlaylistRepository
	PlaylistRepository
}

func NewRepository(db *sql.DB) Repository {
	playlistRepository := NewPlaylistRepository(db)
	return &repository{
		items:              make(map[string]*model.Item),
		PlaylistRepository: playlistRepository,
		mu:                 &sync.RWMutex{},
	}
}

type repository struct {
	items map[string]*model.Item
	mu    *sync.RWMutex
	PlaylistRepository
}

func (r *repository) GetPlaylistRepository() PlaylistRepository {
	return r.PlaylistRepository
}
