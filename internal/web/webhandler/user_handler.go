package webhandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainservice"
	"tictactoe/web/webcontext"
	"tictactoe/web/webmapper"

	"github.com/google/uuid"
)

type UserHandler struct {
	service domainservice.UserService
}

func NewUserHandler(serv domainservice.UserService) *UserHandler {
	return &UserHandler{service: serv}
}
func (h *UserHandler) GetUserByIdFromPath(w http.ResponseWriter, r *http.Request) {
	userIDstr := r.PathValue("id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	h.getUserById(w, r, userID)
}

func (h *UserHandler) GetUserByIdFromCtx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := webcontext.UserID(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	h.getUserById(w, r, userID)

}
func (h *UserHandler) getUserById(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	ctx := r.Context()
	domainUser, err := h.service.GetByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, domainerror.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)

		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}
	response := webmapper.UserToWeb(domainUser)
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
