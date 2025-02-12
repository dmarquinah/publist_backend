package handler

import (
	"net/http"

	"github.com/dmarquinah/publist_backend/internal/service"
)

type Handler struct {
	svc             service.Service
	playlistHandler *PlaylistHandler
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{
		svc:             svc,
		playlistHandler: NewPlaylistHandler(svc),
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	h.playlistHandler.RegisterRoutes(mux)
}
