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
	"tictactoe/domain/domainrepository"
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
		newHandler,
		dsrepository.NewUserRepository,
		appservice.NewUserService,
		appservice.NewAuthService,
		webhandler.NewAuthHandler,
		webmiddleware.NewUserAuthenticator,
	),

	fx.Invoke(runServer),
)

func newHandler(serv domainservice.Service, repo domainrepository.GameRepository) *webhandler.GameHandler {
	return webhandler.NewGameHandler(serv, repo)
}
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
func runServer(lc fx.Lifecycle, handler *webhandler.GameHandler, autHandl *webhandler.AuthHandler, midd *webmiddleware.UserAuthenticator) {
	mux := http.NewServeMux()
	mux.Handle("POST /game/{id}", midd.Wrap(http.HandlerFunc(handler.Handle)))
	mux.HandleFunc("POST /signup", autHandl.SignUp)
	mux.HandleFunc("POST /signin", autHandl.SignIn)
	mux.Handle("POST /game", midd.Wrap(http.HandlerFunc(handler.Create)))
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
