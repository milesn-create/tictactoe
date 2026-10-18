package web

import (
	"net/http"
	"tictactoe/domain"
)

type gameSaver interface {
	Save(game domain.Game) error
}
type Handler struct {
	service domain.Service
	saver   gameSaver
}

func NewHandler(service domain.Service, saver gameSaver) *Handler {
	return &Handler{
		service: service,
		saver:   saver,
	}

}
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {

}
