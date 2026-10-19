package di

import (
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
func runServer(handler *web.Handler) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /game/{id}", handler.Handle)
	log.Fatal(http.ListenAndServe(":8080", mux))

}
