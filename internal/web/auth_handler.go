package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"tictactoe/domain"
)

type AuthHandler struct {
	serv domain.AuthService
}

func NewAuthHandler(serv domain.AuthService) *AuthHandler {
	return &AuthHandler{serv: serv}
}
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var authModel SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&authModel); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	authDomain := signUpToDomain(authModel)
	err := h.serv.SignUp(ctx, authDomain)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
