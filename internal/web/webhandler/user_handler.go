package webhandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainservice"
	"tictactoe/web/webmapper"

	"github.com/google/uuid"
)

type UserHandler struct {
	service domainservice.UserService
}

func NewUserHandler(serv domainservice.UserService) *UserHandler {
	return &UserHandler{service: serv}
}
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDstr := r.PathValue("id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
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
	response := webmapper.ToUserResponse(domainUser)
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
