package model

import "time"

type Playlist struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	HostID      string    `json:"host_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsModerated bool      `json:"is_moderated"`
}

type Track struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Artist   string    `json:"artist"`
	Duration int       `json:"duration"` // in seconds
	AddedAt  time.Time `json:"added_at"`
}

type Playlist_Track struct {
	ID         string    `json:"id"`
	PlaylistID string    `json:"playlist_id"`
	Title      string    `json:"title"`
	Artist     string    `json:"artist"`
	Duration   int       `json:"duration"` // in seconds
	Position   int       `json:"position"` // order in playlist
	AddedAt    time.Time `json:"added_at"`
	IsPlaying  bool      `json:"is_playing"`
}

type Host struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}
