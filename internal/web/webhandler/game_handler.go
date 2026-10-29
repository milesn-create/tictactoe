package webhandler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainservice"
	"tictactoe/web/webcontext"
	"tictactoe/web/webmapper"
	"tictactoe/web/webmodel"

	"github.com/google/uuid"
)

type GameHandler struct {
	service domainservice.Service
}

func NewGameHandler(serv domainservice.Service) *GameHandler {
	return &GameHandler{
		service: serv,
	}

}
func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	gameIDstr := r.PathValue("id")
	gameID, err := uuid.Parse(gameIDstr)
	if err != nil {
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}
	userID, ok := webcontext.UserID(ctx)
	if !ok {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println("failed get id from webcontext")
		return
	}
	var req webmodel.MakeMoveRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	board, err := webmapper.ToBoard(req.Board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gameResult, err := h.service.MakeMove(ctx, gameID, userID, board)
	if err != nil {
		switch {
		case errors.Is(err, domainerror.ErrGameNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, domainerror.ErrNotParticipant):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, domainerror.ErrNotYourTurn), errors.Is(err, domainerror.ErrGameNotStarted), errors.Is(err, domainerror.ErrGameOver):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, domainerror.ErrBoardModified), errors.Is(err, domainerror.ErrInvalidMovesCount):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}
	response := webmapper.ToResponse(gameResult)
	data, err := json.Marshal(response)
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
func (h *GameHandler) Join(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := webcontext.UserID(ctx)
	if !ok {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println("failed get id from webcontext")
		return
	}
	gameIDstr := r.PathValue("id")
	gameID, err := uuid.Parse(gameIDstr)
	if err != nil {
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}
	game, err := h.service.JoinGame(ctx, gameID, userID)
	if err != nil {
		if errors.Is(err, domainerror.ErrGameNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, domainerror.ErrGameNotJoinable) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, domainerror.ErrJoinOwnGame) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	response := webmapper.ToResponse(game)
	data, err := json.Marshal(response)
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
func (h *GameHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := webcontext.UserID(ctx)
	if !ok {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println("failed get id from webcontext")
		return
	}
	gameIDstr := r.PathValue("id")
	gameID, err := uuid.Parse(gameIDstr)
	if err != nil {
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}
	gameGet, err := h.service.GetGame(ctx, gameID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domainerror.ErrGameNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, domainerror.ErrNotParticipant):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}
	response := webmapper.ToResponse(gameGet)
	data, err := json.Marshal(response)
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
