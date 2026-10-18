package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"tictactoe/domain"

	"github.com/google/uuid"
)

type gameSaver interface {
	Save(game domain.Game) error
}
type Handler struct {
	service domain.Service
	saver   gameSaver
}

func NewHandler(service domain.Service, saver gameSaver) *Handler {
	return &Handler{
		service: service,
		saver:   saver,
	}

}
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}
	var webGame Game
	if err := json.NewDecoder(r.Body).Decode(&webGame); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	game, err := toDomain(webGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if game.ID != id {
		http.Error(w, "game id mismatch", http.StatusBadRequest)
		return
	}
	if err := h.service.ValidateBoard(game); err != nil {
		if errors.Is(err, domain.ErrGameNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	over, _ := h.service.IsGameOver(game)
	if !over {
		game, err = h.service.NextMove(game)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := h.saver.Save(game); err != nil {
			http.Error(w, "failed to save game", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(toWeb(game))
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
