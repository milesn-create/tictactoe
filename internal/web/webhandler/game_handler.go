package webhandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"
	"tictactoe/domain/domainservice"
	"tictactoe/web/webcontext"
	"tictactoe/web/webmapper"
	"tictactoe/web/webmodel"

	"github.com/google/uuid"
)

type gameSaver interface {
	Save(ctx context.Context, game domainmodel.Game) error
}
type GameHandler struct {
	service domainservice.Service
	saver   gameSaver
}

func NewGameHandler(serv domainservice.Service, saver gameSaver) *GameHandler {
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
		if errors.Is(err, domainerror.ErrGameNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, domainerror.ErrBoardModified) || errors.Is(err, domainerror.ErrInvalidMovesCount) || errors.Is(err, domainerror.ErrGameOver) {
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
func (h *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idUser, ok := webcontext.UserID(ctx)
	if !ok {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println("failed get id from webcontext")
		return
	}
	var webgame webmodel.CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&webgame); err != nil {
		if !errors.Is(err, io.EOF) {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			log.Println(err)
			return
		}
	}
	game, err := h.service.CreateGame(ctx, idUser, webgame.VsComputer)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	gameResp := webmapper.ToResponse(game)
	data, err := json.Marshal(gameResp)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
func (h *GameHandler) Available(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := webcontext.UserID(ctx)
	if !ok {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println("failed get id from webcontext")
		return
	}
	games, err := h.service.GetAvailable(ctx, userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	responseGames := webmapper.ToResponseList(games)
	data, err := json.Marshal(responseGames)
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
