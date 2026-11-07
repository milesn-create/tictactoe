package webmiddleware

import (
	"net/http"
	"strings"
	"tictactoe/domain/domainservice"
	"tictactoe/web/webcontext"
)

type UserAuthenticator struct {
	jwtProvider domainservice.JwtProvider
}

func NewUserAuthenticator(jwt domainservice.JwtProvider) *UserAuthenticator {
	return &UserAuthenticator{jwtProvider: jwt}
}
func (a *UserAuthenticator) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			w.Header().Set("WWW-Authenticate", `Bearer`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		err := a.jwtProvider.ValidateAccessToken(token)
		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		userID, err := a.jwtProvider.GetUUIDFromToken(token)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Bearer`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = webcontext.WithUserID(ctx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
