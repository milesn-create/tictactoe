package di

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"tictactoe/application/appservice"
	"tictactoe/datasource/dspostgres"
	"tictactoe/datasource/dsrepository"
	"tictactoe/domain/domainservice"
	"tictactoe/web/webhandler"
	"tictactoe/web/webmiddleware"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		newPool,
		dsrepository.NewGameRepository,
		appservice.NewGameService,
		webhandler.NewGameHandler,
		dsrepository.NewUserRepository,
		appservice.NewUserService,
		appservice.NewAuthService,
		webhandler.NewAuthHandler,
		webhandler.NewUserHandler,
		webmiddleware.NewUserAuthenticator,
		newJwtProviderFromEnv,
		dsrepository.NewRefreshTokenRepository,
		appservice.NewRefreshTokenService,
	),

	fx.Invoke(runServer),
)

func newPool(lc fx.Lifecycle) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	pool, err := dspostgres.NewPool(ctx, dbURL)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
		pool.Close()
		return nil
	}})
	return pool, nil

}
func newJwtProviderFromEnv() (domainservice.JwtProvider, error) {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		return nil, errors.New("JWT_ACCESS_SECRET is not set")
	}
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		return nil, errors.New("JWT_REFRESH_SECRET is not set")
	}
	return appservice.NewJwtProvider([]byte(accessSecret), []byte(refreshSecret), (15 * time.Minute), (720 * time.Hour)), nil
}

func runServer(lc fx.Lifecycle, handler *webhandler.GameHandler, autHandl *webhandler.AuthHandler, userHandl *webhandler.UserHandler, midd *webmiddleware.UserAuthenticator) {
	mux := http.NewServeMux()
	mux.Handle("POST /game/{id}", midd.Wrap(http.HandlerFunc(handler.MakeMove)))
	mux.HandleFunc("POST /signup", autHandl.SignUp)
	mux.HandleFunc("POST /signin", autHandl.SignIn)
	mux.HandleFunc("POST /update/access", autHandl.UpdateAccessToken)
	mux.HandleFunc("POST /update/refresh", autHandl.UpdateRefreshToken)
	mux.Handle("POST /game", midd.Wrap(http.HandlerFunc(handler.Create)))
	mux.Handle("GET /game", midd.Wrap(http.HandlerFunc(handler.Available)))
	mux.Handle("POST /game/{id}/join", midd.Wrap(http.HandlerFunc(handler.Join)))
	mux.Handle("GET /game/completed", midd.Wrap(http.HandlerFunc(handler.Completed)))
	mux.Handle("GET /game/{id}", midd.Wrap(http.HandlerFunc(handler.Get)))
	mux.Handle("GET /user/{id}", midd.Wrap(http.HandlerFunc(userHandl.GetUserByIdFromPath)))
	mux.Handle("GET /user/me", midd.Wrap(http.HandlerFunc(userHandl.GetUserByIdFromCtx)))
	mux.Handle("GET /game/leaderboard", midd.Wrap(http.HandlerFunc(handler.Leaderboard)))

	//log.Fatal(http.ListenAndServe(":8080", mux))
	server := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			go func() {
				err := server.Serve(ln)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Println("server error:", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
