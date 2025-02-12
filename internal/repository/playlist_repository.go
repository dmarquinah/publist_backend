package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/dmarquinah/publist_backend/internal/errors"
	"github.com/dmarquinah/publist_backend/internal/model"
)

type PlaylistRepository interface {
	CreatePlaylist(ctx context.Context, playlist *model.Playlist) error
	GetPlaylist(ctx context.Context, id string) (*model.Playlist, error)
	UpdatePlaylist(ctx context.Context, playlist *model.Playlist) error
	DeletePlaylist(ctx context.Context, id string) error
	GetPlaylistsByHost(ctx context.Context, hostID string) ([]*model.Playlist, error)
	AddTrack(ctx context.Context, track *model.Track) error
	RemoveTrack(ctx context.Context, playlistID, trackID string) error
	UpdateTrackPosition(ctx context.Context, playlistID, trackID string, newPosition int) error
	GetCurrentTrack(ctx context.Context, playlistID string) (*model.Track, error)
	GetPlaylistTracks(ctx context.Context, playlistID string) ([]*model.Track, error)
}

type playlistRepository struct {
	db *sql.DB
}

func NewPlaylistRepository(db *sql.DB) PlaylistRepository {
	return &playlistRepository{db: db}
}

func (r *playlistRepository) CreatePlaylist(ctx context.Context, playlist *model.Playlist) error {
	query := `
		INSERT INTO playlists (id, name, host_id, created_at, updated_at, is_moderated)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		playlist.ID,
		playlist.Name,
		playlist.HostID,
		time.Now(),
		time.Now(),
		true, // Create new playlist as able to be moderated
	)
	return err
}

func (r *playlistRepository) GetPlaylist(ctx context.Context, id string) (*model.Playlist, error) {
	query := `
		SELECT id, name, host_id, created_at, updated_at, is_moderated
		FROM playlists
		WHERE id = ?
	`
	playlist := &model.Playlist{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.HostID,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
		&playlist.IsModerated,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrPlaylistNotFound
	}
	return playlist, err
}

func (r *playlistRepository) UpdatePlaylist(ctx context.Context, playlist *model.Playlist) error {
	query := `
		UPDATE playlists
		SET name = ?, updated_at = ?, is_moderated = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query,
		playlist.Name,
		time.Now(),
		playlist.IsModerated,
		playlist.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrPlaylistNotFound
	}
	return nil
}

func (r *playlistRepository) DeletePlaylist(ctx context.Context, id string) error {
	query := `DELETE FROM playlists WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrPlaylistNotFound
	}
	return nil
}

func (r *playlistRepository) GetPlaylistsByHost(ctx context.Context, hostID string) ([]*model.Playlist, error) {
	query := `
		SELECT id, name, host_id, created_at, updated_at, is_moderated
		FROM playlists
		WHERE host_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, hostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []*model.Playlist
	for rows.Next() {
		playlist := &model.Playlist{}
		err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.HostID,
			&playlist.CreatedAt,
			&playlist.UpdatedAt,
			&playlist.IsModerated,
		)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, rows.Err()
}

func (r *playlistRepository) AddTrack(ctx context.Context, track *model.Track) error {
	query := `
		INSERT INTO tracks (id, playlist_id, title, artist, duration, position, added_at, is_playing)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		track.ID,
		track.PlaylistID,
		track.Title,
		track.Artist,
		track.Duration,
		track.Position,
		track.AddedAt,
		track.IsPlaying,
	)
	return err
}

func (r *playlistRepository) RemoveTrack(ctx context.Context, playlistID, trackID string) error {
	query := `DELETE FROM tracks WHERE playlist_id = ? AND id = ?`
	result, err := r.db.ExecContext(ctx, query, playlistID, trackID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrTrackNotFound
	}
	return nil
}

func (r *playlistRepository) UpdateTrackPosition(ctx context.Context, playlistID, trackID string, newPosition int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get current position
	var currentPos int
	err = tx.QueryRowContext(ctx,
		"SELECT position FROM tracks WHERE playlist_id = ? AND id = ?",
		playlistID, trackID).Scan(&currentPos)
	if err != nil {
		return err
	}

	// Update positions of other tracks
	if currentPos < newPosition {
		_, err = tx.ExecContext(ctx,
			`UPDATE tracks 
			SET position = position - 1 
			WHERE playlist_id = ? AND position > ? AND position <= ?`,
			playlistID, currentPos, newPosition)
	} else {
		_, err = tx.ExecContext(ctx,
			`UPDATE tracks 
			SET position = position + 1 
			WHERE playlist_id = ? AND position >= ? AND position < ?`,
			playlistID, newPosition, currentPos)
	}
	if err != nil {
		return err
	}

	// Update position of target track
	_, err = tx.ExecContext(ctx,
		"UPDATE tracks SET position = ? WHERE playlist_id = ? AND id = ?",
		newPosition, playlistID, trackID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *playlistRepository) GetCurrentTrack(ctx context.Context, playlistID string) (*model.Track, error) {
	query := `
		SELECT id, playlist_id, title, artist, duration, position, added_at, is_playing
		FROM tracks
		WHERE playlist_id = ? AND is_playing = true
		LIMIT 1
	`
	track := &model.Track{}
	err := r.db.QueryRowContext(ctx, query, playlistID).Scan(
		&track.ID,
		&track.PlaylistID,
		&track.Title,
		&track.Artist,
		&track.Duration,
		&track.Position,
		&track.AddedAt,
		&track.IsPlaying,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return track, err
}

func (r *playlistRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]*model.Track, error) {
	query := `
		SELECT id, playlist_id, title, artist, duration, position, added_at, is_playing
		FROM tracks
		WHERE playlist_id = ?
		ORDER BY position
	`
	rows, err := r.db.QueryContext(ctx, query, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*model.Track
	for rows.Next() {
		track := &model.Track{}
		err := rows.Scan(
			&track.ID,
			&track.PlaylistID,
			&track.Title,
			&track.Artist,
			&track.Duration,
			&track.Position,
			&track.AddedAt,
			&track.IsPlaying,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	return tracks, rows.Err()
}
