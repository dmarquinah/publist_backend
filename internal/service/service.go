package service

import (
	"github.com/dmarquinah/publist_backend/internal/repository"
)

type Service interface {
	// Add more service methods as needed
	PlaylistService
}

type service struct {
	repo            repository.Repository
	PlaylistService // Add PlaylistService field
}

func NewService(repo repository.Repository) Service {
	playlistService := NewPlaylistService(repo.GetPlaylistRepository())
	return &service{
		repo:            repo,
		PlaylistService: playlistService, // Initialize PlaylistService
	}
}
