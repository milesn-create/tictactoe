package datasource

import "tictactoe/domain"

type gameService struct {
	repo Repository
}

func NewService(repo Repository) domain.Service {
	return &gameService{repo: repo}

}

func (s *gameService) NextMove(game domain.Game) (domain.Game, error) {
	return domain.Game{}, nil

}

func (s *gameService) ValidateBoard(game domain.Game) error {
	return nil

}
func (s *gameService) IsGameOver(game domain.Game) (bool, domain.Cell) {
	return false, domain.Empty

}
