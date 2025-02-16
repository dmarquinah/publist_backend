package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/dmarquinah/publist_backend/internal/auth"
	errorsmsg "github.com/dmarquinah/publist_backend/internal/errors"
	"github.com/dmarquinah/publist_backend/internal/model"
	"github.com/dmarquinah/publist_backend/internal/service"
	"github.com/google/uuid"
)

type PlaylistHandler struct {
	svc service.PlaylistService
}

func NewPlaylistHandler(svc service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{
		svc: svc,
	}
}

func (h *PlaylistHandler) RegisterRoutes(mux *http.ServeMux) {
	// Public endpoints
	mux.HandleFunc("GET /playlists/{id}", h.GetPlaylist)
	mux.HandleFunc("GET /playlists/{id}/current", h.GetCurrentTrack)
	mux.HandleFunc("GET /playlists/{id}/tracks", h.GetPlaylistTracks)

	// Host endpoints
	/* mux.HandleFunc("POST /host/playlists", h.requireRole("host", h.CreatePlaylist))
	mux.HandleFunc("PUT /host/playlists/{id}", h.requireRole("host", h.UpdatePlaylist))
	mux.HandleFunc("DELETE /host/playlists/{id}", h.requireRole("host", h.DeletePlaylist))
	mux.HandleFunc("GET /host/playlists", h.requireRole("host", h.GetHostPlaylists))
	mux.HandleFunc("POST /host/playlists/{id}/tracks", h.requireRole("host", h.AddTrack))
	mux.HandleFunc("DELETE /host/playlists/{id}/tracks/{trackId}", h.requireRole("host", h.RemoveTrack))
	mux.HandleFunc("PUT /host/playlists/{id}/tracks/{trackId}/position", h.requireRole("host", h.ReorderTrack)) */

	// Admin endpoints
	/* mux.HandleFunc("PUT /admin/playlists/{id}/moderate", h.requireRole("admin", h.ModeratePlaylist)) */
}

func (h *PlaylistHandler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	playlist, err := h.svc.GetPlaylist(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errorsmsg.ErrPlaylistNotFound):
			http.Error(w, "Playlist not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	respondJSON(w, http.StatusOK, playlist)
}

func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)

	var playlist model.Playlist
	if err := json.NewDecoder(r.Body).Decode(&playlist); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	playlist.ID = uuid.New().String()
	playlist.HostID = claims.UserID

	if err := h.svc.CreatePlaylist(r.Context(), &playlist); err != nil {
		switch {
		case errors.Is(err, errorsmsg.ErrInvalidName):
			http.Error(w, "Invalid playlist name", http.StatusBadRequest)
		case errors.Is(err, errorsmsg.ErrNameTooLong):
			http.Error(w, "Playlist name too long", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	respondJSON(w, http.StatusCreated, playlist)
}

func (h *PlaylistHandler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)
	id := r.PathValue("id")

	var playlist model.Playlist
	if err := json.NewDecoder(r.Body).Decode(&playlist); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	playlist.ID = id
	playlist.HostID = claims.UserID // We set the HostID from the claims

	if err := h.svc.UpdatePlaylist(r.Context(), &playlist); err != nil {
		switch {
		case errors.Is(err, errorsmsg.ErrUnauthorized):
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		case errors.Is(err, errorsmsg.ErrPlaylistNotFound):
			http.Error(w, "Playlist not found", http.StatusNotFound)
		case errors.Is(err, errorsmsg.ErrInvalidName):
			http.Error(w, "Invalid playlist name", http.StatusBadRequest)
		case errors.Is(err, errorsmsg.ErrNameTooLong):
			http.Error(w, "Playlist name too long", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	respondJSON(w, http.StatusOK, playlist)
}

func (h *PlaylistHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)
	id := r.PathValue("id")

	if err := h.svc.DeletePlaylist(r.Context(), id, claims.UserID, claims.Role == "admin"); err != nil {
		switch {
		case errors.Is(err, errorsmsg.ErrUnauthorized):
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		case errors.Is(err, errorsmsg.ErrPlaylistNotFound):
			http.Error(w, "Playlist not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlaylistHandler) GetHostPlaylists(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)

	playlists, err := h.svc.GetPlaylistsByHost(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, playlists)
}

func (h *PlaylistHandler) AddTrack(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)
	playlistID := r.PathValue("id")

	var track model.Playlist_Track
	if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	track.ID = uuid.New().String()
	track.ID = playlistID

	if err := h.svc.AddTrack(r.Context(), &track, claims.UserID); err != nil {
		switch {
		case errors.Is(err, errorsmsg.ErrUnauthorized):
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		case errors.Is(err, errorsmsg.ErrPlaylistNotFound):
			http.Error(w, "Playlist not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	respondJSON(w, http.StatusCreated, track)
}

func (h *PlaylistHandler) RemoveTrack(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)
	playlistID := r.PathValue("id")
	trackID := r.PathValue("trackId")

	if err := h.svc.RemoveTrack(r.Context(), playlistID, trackID, claims.UserID); err != nil {
		switch {
		case errors.Is(err, errorsmsg.ErrUnauthorized):
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		case errors.Is(err, errorsmsg.ErrPlaylistNotFound):
			http.Error(w, "Playlist not found", http.StatusNotFound)
		case errors.Is(err, errorsmsg.ErrTrackNotFound):
			http.Error(w, "Track not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlaylistHandler) ReorderTrack(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)
	playlistID := r.PathValue("id")
	trackID := r.PathValue("trackId")

	var body struct {
		Position int `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.ReorderTrack(r.Context(), playlistID, trackID, body.Position, claims.UserID); err != nil {
		switch {
		case errors.Is(err, errorsmsg.ErrUnauthorized):
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		case errors.Is(err, errorsmsg.ErrPlaylistNotFound):
			http.Error(w, "Playlist not found", http.StatusNotFound)
		case errors.Is(err, errorsmsg.ErrTrackNotFound):
			http.Error(w, "Track not found", http.StatusNotFound)
		case errors.Is(err, errorsmsg.ErrInvalidPosition):
			http.Error(w, "Invalid position", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PlaylistHandler) GetCurrentTrack(w http.ResponseWriter, r *http.Request) {
	playlistID := r.PathValue("id")

	track, err := h.svc.GetCurrentTrack(r.Context(), playlistID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if track == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	respondJSON(w, http.StatusOK, track)
}

func (h *PlaylistHandler) GetPlaylistTracks(w http.ResponseWriter, r *http.Request) {
	playlistID := r.PathValue("id")

	tracks, err := h.svc.GetPlaylistTracks(r.Context(), playlistID)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, tracks)
}

func (h *PlaylistHandler) ModeratePlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := r.PathValue("id")

	var body struct {
		IsModerated bool `json:"is_moderated"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.ModeratePlaylist(r.Context(), playlistID, body.IsModerated); err != nil {
		switch {
		case errors.Is(err, errorsmsg.ErrPlaylistNotFound):
			http.Error(w, "Playlist not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Middleware for role checking
func (h *PlaylistHandler) requireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*auth.Claims)
		if !ok || claims.Role != role {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// Helper function to send JSON responses
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
