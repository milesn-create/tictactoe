package webhandler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"tictactoe/domain"
	"tictactoe/domain/domainmodel"
	"tictactoe/web/webmapper"
	"tictactoe/web/webmodel"

	"github.com/google/uuid"
)

type gameSaver interface {
	Save(ctx context.Context, game domainmodel.Game) error
}
type GameHandler struct {
	service domain.Service
	saver   gameSaver
}

func NewGameHandler(serv domain.Service, saver gameSaver) *GameHandler {
	return &GameHandler{
		service: serv,
		saver:   saver,
	}

}
func (h *GameHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	strID := r.PathValue("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}
	var webGame webmodel.Game
	if err := json.NewDecoder(r.Body).Decode(&webGame); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	game, err := webmapper.ToDomain(webGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if game.ID != id {
		http.Error(w, "game id mismatch", http.StatusBadRequest)
		return
	}
	if err := h.service.ValidateBoard(ctx, game); err != nil {
		if errors.Is(err, domain.ErrGameNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrBoardModified) || errors.Is(err, domain.ErrInvalidMovesCount) || errors.Is(err, domain.ErrGameOver) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	over, _ := h.service.IsGameOver(game)
	if !over {
		game, err = h.service.NextMove(game)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	if err := h.saver.Save(ctx, game); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	data, err := json.Marshal(webmapper.ToWeb(game))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
