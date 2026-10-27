package webmiddleware

import (
	"errors"
	"log"
	"net/http"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainservice"
	"tictactoe/web/webcontext"
)

type UserAuthenticator struct {
	serv domainservice.AuthService
}

func NewUserAuthenticator(serv domainservice.AuthService) *UserAuthenticator {
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
		id, err := a.serv.Authorize(ctx, login, password)
		if err != nil {
			if errors.Is(err, domainerror.ErrInvalidCredentials) {
				w.Header().Set("WWW-Authenticate", `Basic realm="tictactoe"`)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		ctx = webcontext.WithUserID(ctx, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
