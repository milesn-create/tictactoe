package webhandler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"
	"tictactoe/domain/domainservice"
	"tictactoe/web/webmapper"
	"tictactoe/web/webmodel"
)

type AuthHandler struct {
	serv domainservice.AuthService
}

func NewAuthHandler(serv domainservice.AuthService) *AuthHandler {
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
		if errors.Is(err, domainerror.ErrUserExists) {
			http.Error(w, domainerror.ErrUserExists.Error(), http.StatusConflict)
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
	var signInModel webmodel.JwtRequest
	if err := json.NewDecoder(r.Body).Decode(&signInModel); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	signInDomain, err := webmapper.SignInToDomain(signInModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jwtResponseDomain, err := h.serv.Authorize(ctx, signInDomain)
	if err != nil {
		if errors.Is(err, domainerror.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	jwtResponseWeb := webmapper.JwtToWeb(jwtResponseDomain)

	data, err := json.Marshal(jwtResponseWeb)
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

func (h *AuthHandler) UpdateAccessToken(w http.ResponseWriter, r *http.Request) {
	h.serveRefresh(w, r, h.serv.UpdateAccessToken)
}

func (h *AuthHandler) UpdateRefreshToken(w http.ResponseWriter, r *http.Request) {
	h.serveRefresh(w, r, h.serv.UpdateRefreshToken)
}

func (h *AuthHandler) serveRefresh(w http.ResponseWriter, r *http.Request, update func(context.Context, domainmodel.RefreshJwtRequest) (domainmodel.JwtResponse, error)) {
	var jwtRequest webmodel.RefreshJwtRequest
	err := json.NewDecoder(r.Body).Decode(&jwtRequest)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	domainjwtRequest, err := webmapper.RefreshToDomain(jwtRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	response, err := update(ctx, domainjwtRequest)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Bearer`)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	responseWeb := webmapper.JwtToWeb(response)
	data, err := json.Marshal(responseWeb)
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
