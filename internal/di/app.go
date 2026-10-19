package di

import (
	"context"
	"errors"
	"log"
	"net/http"
	"tictactoe/datasource"
	"tictactoe/domain"
	"tictactoe/web"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		datasource.NewStorage,
		datasource.NewRepository,
		datasource.NewService,
		newHandler),
	fx.Invoke(runServer),
)

func newHandler(service domain.Service, repo datasource.Repository) *web.Handler {
	return web.NewHandler(service, repo)
}
func runServer(lc fx.Lifecycle, handler *web.Handler) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /game/{id}", handler.Handle)
	//log.Fatal(http.ListenAndServe(":8080", mux))
	server := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := server.ListenAndServe()
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
