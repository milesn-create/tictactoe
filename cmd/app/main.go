package main

import (
	"tictactoe/di"

	"go.uber.org/fx"
)

func main() {
	fx.New(di.Module).Run()
	//storage := datasource.NewStorage()
	//repo := datasource.NewRepository(storage)
	//service := datasource.NewService(repo)
	//handler := web.NewHandler(service, repo)
	//mux := http.NewServeMux()
	//mux.HandleFunc("POST /game/{id}", handler.Handle)
	//log.Fatal(http.ListenAndServe(":8080", mux))

}
