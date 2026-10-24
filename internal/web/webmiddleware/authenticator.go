package webmiddleware

import (
	"errors"
	"log"
	"net/http"
	"tictactoe/domain"
)

type UserAuthenticator struct {
	serv domain.AuthService
}

func NewUserAuthenticator(serv domain.AuthService) *UserAuthenticator {
	return &UserAuthenticator{serv: serv}
}
func (a *UserAuthenticator) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="tictactoe"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		_, err := a.serv.Authorize(ctx, login, password)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidCredentials) {
				w.Header().Set("WWW-Authenticate", `Basic realm="tictactoe"`)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		next.ServeHTTP(w, r)
	})

}
