package main

import (
	"log"
	"net/http"
	"tictactoe/datasource"
	"tictactoe/web"
)

func main() {
	storage := datasource.NewStorage()
	repo := datasource.NewRepository(storage)
	service := datasource.NewService(repo)
	handler := web.NewHandler(service, repo)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /game/{id}", handler.Handle)
	log.Fatal(http.ListenAndServe(":8080", mux))

}
