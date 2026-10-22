package di

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"tictactoe/datasource/postgres"
	"tictactoe/datasource/repository"
	"tictactoe/datasource/service"
	"tictactoe/domain"
	"tictactoe/web"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		newPool,
		repository.NewRepository,
		service.NewService,
		newHandler),
	fx.Invoke(runServer),
)

func newHandler(serv domain.Service, repo repository.Repository) *web.Handler {
	return web.NewHandler(serv, repo)
}
func newPool(lc fx.Lifecycle) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	pool, err := postgres.NewPool(ctx, dbURL)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
		pool.Close()
		return nil
	}})
	return pool, nil

}
func runServer(lc fx.Lifecycle, handler *web.Handler) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /game/{id}", handler.Handle)
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
