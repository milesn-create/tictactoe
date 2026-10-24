package webhandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"tictactoe/domain"
	"tictactoe/web/webmapper"
	"tictactoe/web/webmodel"
)

type AuthHandler struct {
	serv domain.AuthService
}

func NewAuthHandler(serv domain.AuthService) *AuthHandler {
	return &AuthHandler{serv: serv}
}
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var authModel webmodel.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&authModel); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	authDomain, err := webmapper.SignUpToDomain(authModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.serv.SignUp(ctx, authDomain)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			http.Error(w, domain.ErrUserExists.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	login, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := h.serv.Authorize(ctx, login, password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	data, err := json.Marshal(webmodel.SignInResponse{ID: id})
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
