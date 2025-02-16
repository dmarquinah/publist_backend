package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	errorsmsg "github.com/dmarquinah/publist_backend/internal/errors"
	"github.com/dmarquinah/publist_backend/internal/model"
	"github.com/dmarquinah/publist_backend/internal/repository"
)

type PlaylistService interface {
	CreatePlaylist(ctx context.Context, playlist *model.Playlist) error
	GetPlaylist(ctx context.Context, id string) (*model.Playlist, error)
	UpdatePlaylist(ctx context.Context, playlist *model.Playlist) error
	DeletePlaylist(ctx context.Context, id string, userID string, isAdmin bool) error
	AddTrack(ctx context.Context, track *model.Playlist_Track, userID string) error
	RemoveTrack(ctx context.Context, playlistID, trackID string, userID string) error
	ReorderTrack(ctx context.Context, playlistID, trackID string, newPosition int, userID string) error
	GetCurrentTrack(ctx context.Context, playlistID string) (*model.Playlist_Track, error)
	GetPlaylistTracks(ctx context.Context, playlistID string) ([]*model.Playlist_Track, error)
	ModeratePlaylist(ctx context.Context, playlistID string, isModerated bool) error
	GetPlaylistsByHost(ctx context.Context, hostID string) ([]*model.Playlist, error)
}

type playlistService struct {
	repo repository.PlaylistRepository
}

func NewPlaylistService(repo repository.PlaylistRepository) PlaylistService {
	return &playlistService{repo: repo}
}

func (s *playlistService) CreatePlaylist(ctx context.Context, playlist *model.Playlist) error {
	if err := s.validatePlaylist(playlist); err != nil {
		return fmt.Errorf("validating playlist: %w", err)
	}

	playlist.CreatedAt = time.Now()
	playlist.UpdatedAt = time.Now()
	playlist.IsModerated = false

	return s.repo.CreatePlaylist(ctx, playlist)
}

func (s *playlistService) GetPlaylist(ctx context.Context, id string) (*model.Playlist, error) {
	playlist, err := s.repo.GetPlaylist(ctx, id)
	if err != nil {
		if errors.Is(err, errorsmsg.ErrPlaylistNotFound) {
			return nil, errorsmsg.ErrPlaylistNotFound
		}
		return nil, fmt.Errorf("fetching playlist: %w", err)
	}
	return playlist, nil
}

func (s *playlistService) UpdatePlaylist(ctx context.Context, playlist *model.Playlist) error {
	existing, err := s.repo.GetPlaylist(ctx, playlist.ID)
	if err != nil {
		if errors.Is(err, errorsmsg.ErrPlaylistNotFound) {
			return errorsmsg.ErrPlaylistNotFound
		}
		return fmt.Errorf("fetching playlist: %w", err)
	}

	if err := s.validatePlaylist(playlist); err != nil {
		return fmt.Errorf("validating playlist: %w", err)
	}

	// Preserve immutable fields
	playlist.HostID = existing.HostID
	playlist.CreatedAt = existing.CreatedAt
	playlist.UpdatedAt = time.Now()
	playlist.IsModerated = existing.IsModerated

	return s.repo.UpdatePlaylist(ctx, playlist)
}

func (s *playlistService) DeletePlaylist(ctx context.Context, id string, userID string, isAdmin bool) error {
	existing, err := s.repo.GetPlaylist(ctx, id)
	if err != nil {
		if errors.Is(err, errorsmsg.ErrPlaylistNotFound) {
			return errorsmsg.ErrPlaylistNotFound
		}
		return fmt.Errorf("fetching playlist: %w", err)
	}

	if !isAdmin && existing.HostID != userID {
		return errorsmsg.ErrUnauthorized
	}

	return s.repo.DeletePlaylist(ctx, id)
}

func (s *playlistService) GetPlaylistsByHost(ctx context.Context, hostID string) ([]*model.Playlist, error) {
	playlists, err := s.repo.GetPlaylistsByHost(ctx, hostID)
	if err != nil {
		return nil, fmt.Errorf("fetching host playlists: %w", err)
	}
	return playlists, nil
}

func (s *playlistService) AddTrack(ctx context.Context, track *model.Playlist_Track, userID string) error {
	playlist, err := s.repo.GetPlaylist(ctx, track.PlaylistID)
	if err != nil {
		if errors.Is(err, errorsmsg.ErrPlaylistNotFound) {
			return errorsmsg.ErrPlaylistNotFound
		}
		return fmt.Errorf("fetching playlist: %w", err)
	}

	if playlist.HostID != userID {
		return errorsmsg.ErrUnauthorized
	}

	tracks, err := s.repo.GetPlaylistTracks(ctx, track.PlaylistID)
	if err != nil {
		return fmt.Errorf("fetching playlist tracks: %w", err)
	}

	track.Position = len(tracks) + 1
	track.AddedAt = time.Now()
	track.IsPlaying = false

	return s.repo.AddTrack(ctx, track)
}

func (s *playlistService) RemoveTrack(ctx context.Context, playlistID, trackID string, userID string) error {
	playlist, err := s.repo.GetPlaylist(ctx, playlistID)
	if err != nil {
		if errors.Is(err, errorsmsg.ErrPlaylistNotFound) {
			return errorsmsg.ErrPlaylistNotFound
		}
		return fmt.Errorf("fetching playlist: %w", err)
	}

	if playlist.HostID != userID {
		return errorsmsg.ErrUnauthorized
	}

	if err := s.repo.RemoveTrack(ctx, playlistID, trackID); err != nil {
		if errors.Is(err, errorsmsg.ErrTrackNotFound) {
			return errorsmsg.ErrTrackNotFound
		}
		return fmt.Errorf("removing track: %w", err)
	}

	return nil
}

func (s *playlistService) ReorderTrack(ctx context.Context, playlistID, trackID string, newPosition int, userID string) error {
	playlist, err := s.repo.GetPlaylist(ctx, playlistID)
	if err != nil {
		if errors.Is(err, errorsmsg.ErrPlaylistNotFound) {
			return errorsmsg.ErrPlaylistNotFound
		}
		return fmt.Errorf("fetching playlist: %w", err)
	}

	if playlist.HostID != userID {
		return errorsmsg.ErrUnauthorized
	}

	tracks, err := s.repo.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("fetching playlist tracks: %w", err)
	}

	if newPosition < 1 || newPosition > len(tracks) {
		return errorsmsg.ErrInvalidPosition
	}

	return s.repo.UpdateTrackPosition(ctx, playlistID, trackID, newPosition)
}

func (s *playlistService) GetCurrentTrack(ctx context.Context, playlistID string) (*model.Playlist_Track, error) {
	track, err := s.repo.GetCurrentTrack(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("fetching current track: %w", err)
	}
	return track, nil
}

func (s *playlistService) GetPlaylistTracks(ctx context.Context, playlistID string) ([]*model.Playlist_Track, error) {
	tracks, err := s.repo.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("fetching playlist tracks: %w", err)
	}
	return tracks, nil
}

func (s *playlistService) ModeratePlaylist(ctx context.Context, playlistID string, isModerated bool) error {
	playlist, err := s.repo.GetPlaylist(ctx, playlistID)
	if err != nil {
		if errors.Is(err, errorsmsg.ErrPlaylistNotFound) {
			return errorsmsg.ErrPlaylistNotFound
		}
		return fmt.Errorf("fetching playlist: %w", err)
	}

	playlist.IsModerated = isModerated
	playlist.UpdatedAt = time.Now()

	return s.repo.UpdatePlaylist(ctx, playlist)
}

func (s *playlistService) validatePlaylist(p *model.Playlist) error {
	if p.Name == "" {
		return errorsmsg.ErrInvalidName
	}
	if len(p.Name) > 255 {
		return errorsmsg.ErrNameTooLong
	}
	return nil
}
